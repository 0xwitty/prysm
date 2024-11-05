package validator

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/epbs"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/core/helpers"
	"github.com/prysmaticlabs/prysm/v5/beacon-chain/state"
	"github.com/prysmaticlabs/prysm/v5/config/params"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/blocks"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/interfaces"
	"github.com/prysmaticlabs/prysm/v5/consensus-types/primitives"
	"github.com/prysmaticlabs/prysm/v5/encoding/ssz"
	enginev1 "github.com/prysmaticlabs/prysm/v5/proto/engine/v1"
	eth "github.com/prysmaticlabs/prysm/v5/proto/prysm/v1alpha1"
	"github.com/prysmaticlabs/prysm/v5/runtime/version"
	"github.com/prysmaticlabs/prysm/v5/time/slots"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

// SubmitSignedExecutionPayloadEnvelope submits a signed execution payload envelope to the validator client.
func (vs *Server) SubmitSignedExecutionPayloadEnvelope(ctx context.Context, env *enginev1.SignedExecutionPayloadEnvelope) (*emptypb.Empty, error) {
	if err := vs.P2P.Broadcast(ctx, env); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to broadcast signed execution payload envelope: %v", err)
	}

	m, err := blocks.WrappedROExecutionPayloadEnvelope(env.Message)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed to wrap execution payload envelope: %v", err)
	}

	if err := vs.ExecutionPayloadReceiver.ReceiveExecutionPayloadEnvelope(ctx, m, nil); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to receive execution payload envelope: %v", err)
	}

	return nil, nil
}

// GetExecutionPayloadEnvelope returns the execution payload envelope for a given slot.
func (vs *Server) GetExecutionPayloadEnvelope(ctx context.Context, req *eth.PayloadEnvelopeRequest) (*enginev1.ExecutionPayloadEnvelope, error) {
	if vs.payloadEnvelope == nil {
		return nil, status.Error(codes.NotFound, "No execution payload response available")
	}
	if req.ProposerIndex != vs.payloadEnvelope.BuilderIndex {
		return nil, status.Errorf(codes.InvalidArgument, "proposer index mismatch: expected %d, got %d", vs.payloadEnvelope.BuilderIndex, req.ProposerIndex)
	}
	if req.Slot != vs.TimeFetcher.CurrentSlot() {
		return nil, status.Errorf(codes.InvalidArgument, "current slot mismatch: expected %d, got %d", vs.TimeFetcher.CurrentSlot(), req.Slot)
	}

	_, r := vs.ForkchoiceFetcher.HighestReceivedBlockSlotRoot()
	payloadStatus := vs.ForkchoiceFetcher.GetPTCVote(r)

	if payloadStatus == primitives.PAYLOAD_WITHHELD {
		return &enginev1.ExecutionPayloadEnvelope{
			Payload:            nil, // TODO: I'm not sure if I need to pass in and hydrate a empty payload here.
			BuilderIndex:       req.ProposerIndex,
			BeaconBlockRoot:    r[:],
			BlobKzgCommitments: [][]byte{},
			PayloadWithheld:    true,
			StateRoot:          []byte{},
		}, nil
	}

	// TODO: calculate state root
	var stateRoot []byte
	vs.payloadEnvelope.StateRoot = stateRoot

	return vs.payloadEnvelope, nil
}

func (vs *Server) SubmitSignedExecutionPayloadHeader(ctx context.Context, h *enginev1.SignedExecutionPayloadHeader) (*emptypb.Empty, error) {
	if vs.TimeFetcher.CurrentSlot() != h.Message.Slot {
		return nil, status.Errorf(codes.InvalidArgument, "current slot mismatch: expected %d, got %d", vs.TimeFetcher.CurrentSlot(), h.Message.Slot)
	}

	vs.signedExecutionPayloadHeader = h
	log.Info("Cached signed execution payload header ", h.Message.Slot, h.Message.BuilderIndex)

	headState, _, err := vs.getParentState(ctx, h.Message.Slot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve head state: %v", err)
	}

	proposerIndex, err := helpers.BeaconProposerIndexAtSlot(ctx, headState, h.Message.Slot)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to retrieve proposer index: %v", err)
	}
	if proposerIndex != h.Message.BuilderIndex {
		if err := vs.P2P.Broadcast(ctx, h); err != nil {
			return nil, status.Errorf(codes.Internal, "failed to broadcast signed execution payload header: %v", err)
		}
	}

	return nil, nil
}

// computePostPayloadStateRoot computes the state root after an execution
// payload envelope has been processed through a state transition and
// returns it to the validator client.
func (vs *Server) computePostPayloadStateRoot(ctx context.Context, envelope interfaces.ROExecutionPayloadEnvelope) ([]byte, error) {
	beaconState, err := vs.StateGen.StateByRoot(ctx, envelope.BeaconBlockRoot())
	if err != nil {
		return nil, errors.Wrap(err, "could not retrieve beacon state")
	}
	beaconState = beaconState.Copy()
	err = epbs.ProcessPayloadStateTransition(
		ctx,
		beaconState,
		envelope,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "could not calculate post payload state root at slot %d", beaconState.Slot())
	}

	root, err := beaconState.HashTreeRoot(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "could not calculate post payload state root at slot %d", beaconState.Slot())
	}
	log.WithField("beaconStateRoot", fmt.Sprintf("%#x", root)).Debugf("Computed state root at execution stage")
	return root[:], nil
}

// GetLocalHeader returns the local header for a given slot and proposer index.
func (vs *Server) GetLocalHeader(ctx context.Context, req *eth.HeaderRequest) (*enginev1.ExecutionPayloadHeaderEPBS, error) {
	if vs.SyncChecker.Syncing() {
		return nil, status.Error(codes.FailedPrecondition, "Syncing to latest head, not ready to respond")
	}

	if err := vs.optimisticStatus(ctx); err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, "Validator is not ready to propose: %v", err)
	}

	slot := req.Slot
	epoch := slots.ToEpoch(slot)
	if params.BeaconConfig().EPBSForkEpoch > epoch {
		return nil, status.Errorf(codes.FailedPrecondition, "EPBS fork has not occurred yet")
	}
	if slot != vs.TimeFetcher.CurrentSlot() {
		return nil, status.Errorf(codes.InvalidArgument, "current slot mismatch: expected %d, got %d", vs.TimeFetcher.CurrentSlot(), slot)
	}

	st, parentRoot, err := vs.getParentState(ctx, slot)
	if err != nil {
		return nil, err
	}

	proposerIndex := req.ProposerIndex
	localPayload, err := vs.getLocalPayloadFromEngine(ctx, st, parentRoot, slot, proposerIndex)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get local payload: %v", err)
	}
	electraPayload, ok := localPayload.ExecutionData.Proto().(*enginev1.ExecutionPayloadElectra)
	if !ok {
		return nil, status.Error(codes.Internal, "Could not get electra payload")
	}
	vs.payloadEnvelope = &enginev1.ExecutionPayloadEnvelope{
		Payload:      electraPayload,
		BuilderIndex: proposerIndex,
	}
	vs.blobsBundle = localPayload.BlobsBundle

	kzgRoot, err := ssz.KzgCommitmentsRoot(localPayload.BlobsBundle.KzgCommitments)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could not get kzg commitments root: %v", err)
	}

	return &enginev1.ExecutionPayloadHeaderEPBS{
		ParentBlockHash:        localPayload.ExecutionData.ParentHash(),
		ParentBlockRoot:        parentRoot[:],
		BlockHash:              localPayload.ExecutionData.BlockHash(),
		GasLimit:               localPayload.ExecutionData.GasLimit(),
		BuilderIndex:           proposerIndex,
		Slot:                   slot,
		Value:                  0,
		BlobKzgCommitmentsRoot: kzgRoot[:],
	}, nil
}

// Set payload attestations for the block if it meets the following conditions:
// - Block version is ePBS or higher
// - Parent block's slot is exactly one less than the current block's slot
// - Local chain view is respected to determine whether the parent slot is full or empty:
//   - If the parent slot is empty, pack a "WITHHELD" status
//   - If the parent slot is full, pack a "PRESENT" status
func (vs *Server) setPayloadAttestations(block interfaces.SignedBeaconBlock, state state.BeaconState) error {
	if block.Version() < version.EPBS {
		return nil
	}

	parentRoot := block.Block().ParentRoot()
	parentSlot, err := vs.ForkchoiceFetcher.RecentBlockSlot(parentRoot)
	if err != nil {
		return err
	}

	if parentSlot+1 != block.Block().Slot() {
		return nil
	}

	lastFullSlot, err := state.LatestFullSlot()
	if err != nil {
		return err
	}
	payloadIsPresent := lastFullSlot+1 == block.Block().Slot()
	status := primitives.PAYLOAD_WITHHELD
	if payloadIsPresent {
		status = primitives.PAYLOAD_PRESENT
	}

	payloadAttestation := vs.PayloadAttestationCache.Get(parentRoot, status)
	return block.SetPayloadAttestations([]*eth.PayloadAttestation{payloadAttestation})
}
