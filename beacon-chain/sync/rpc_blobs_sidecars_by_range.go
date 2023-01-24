package sync

import (
	"context"
	"time"

	libp2pcore "github.com/libp2p/go-libp2p/core"
	"github.com/pkg/errors"
	p2ptypes "github.com/prysmaticlabs/prysm/v3/beacon-chain/p2p/types"
	"github.com/prysmaticlabs/prysm/v3/cmd/beacon-chain/flags"
	"github.com/prysmaticlabs/prysm/v3/config/params"
	"github.com/prysmaticlabs/prysm/v3/encoding/bytesutil"
	"github.com/prysmaticlabs/prysm/v3/monitoring/tracing"
	pb "github.com/prysmaticlabs/prysm/v3/proto/prysm/v1alpha1"
	"go.opencensus.io/trace"
)

type BlobsSidecarProcessor func(sidecar *pb.BlobsSidecar) error

// blobsSidecarsByRangeRPCHandler looks up the request blobs from the database from a given start slot index
func (s *Service) blobsSidecarsByRangeRPCHandler(ctx context.Context, msg interface{}, stream libp2pcore.Stream) error {
	ctx, span := trace.StartSpan(ctx, "sync.BlobsSidecarsByRangeHandler")
	defer span.End()
	ctx, cancel := context.WithTimeout(ctx, respTimeout)
	defer cancel()
	SetRPCStreamDeadlines(stream)

	r, ok := msg.(*pb.BlobsSidecarsByRangeRequest)
	if !ok {
		return errors.New("message is not type *pb.BlobsSidecarsByRangeRequest")
	}

	startSlot := r.StartSlot
	count := r.Count
	endSlot := startSlot.Add(count)

	allowedBlocksPerSecond := uint64(flags.Get().BlockBatchLimit)
	var numBlobs uint64
	maxRequestBlobsSidecars := params.BeaconNetworkConfig().MaxRequestBlobsSidecars
	for slot := startSlot; slot < endSlot && numBlobs < maxRequestBlobsSidecars; slot = slot.Add(1) {
		if err := s.rateLimiter.validateRequest(stream, allowedBlocksPerSecond); err != nil {
			return err
		}

		sidecars, err := s.cfg.beaconDB.BlobsSidecarsBySlot(ctx, slot)
		if err != nil {
			s.writeErrorResponseToStream(responseCodeServerError, p2ptypes.ErrGeneric.Error(), stream)
			tracing.AnnotateError(span, err)
			return err
		}
		if len(sidecars) == 0 {
			continue
		}

		var respondedBlobs int64
		for _, sidecar := range sidecars {
			if bytesutil.ToBytes32(sidecar.BeaconBlockRoot) == params.BeaconConfig().ZeroHash {
				continue
			}
			SetStreamWriteDeadline(stream, defaultWriteDuration)
			if chunkErr := WriteBlobsSidecarChunk(stream, s.cfg.chain, s.cfg.p2p.Encoding(), sidecar); chunkErr != nil {
				log.WithError(chunkErr).Debug("Could not send a chunked response")
				s.writeErrorResponseToStream(responseCodeServerError, p2ptypes.ErrGeneric.Error(), stream)
				tracing.AnnotateError(span, chunkErr)
				return chunkErr
			}
			respondedBlobs++
		}
		numBlobs++
		s.rateLimiter.add(stream, respondedBlobs)

		// Short-circuit immediately once we've sent the last blob.
		if slot.Add(1) >= endSlot {
			break
		}

		key := stream.Conn().RemotePeer().String()
		sidecarLimiter, err := s.rateLimiter.topicCollector(string(stream.Protocol()))
		if err != nil {
			return err
		}
		// Throttling - wait until we have enough tokens to send the next blobs

		allowedBlocksBurst := int64(flags.Get().BlockBatchLimitBurstFactor * flags.Get().BlockBatchLimit)
		if sidecarLimiter.Remaining(key) < allowedBlocksBurst {
			timer := time.NewTimer(sidecarLimiter.TillEmpty(key))
			select {
			case <-ctx.Done():
				timer.Stop()
				return ctx.Err()
			case <-timer.C:
				timer.Stop()
			}
		}
	}

	closeStream(stream, log)
	return nil
}

func estimateBlobsSidecarCost(sidecar *pb.BlobsSidecar) int {
	// This represents the fixed cost (in bytes) of the beacon_block_root and beacon_block_slot fields in the sidecar
	const overheadCost = 32 + 8
	cost := overheadCost
	for _, blob := range sidecar.Blobs {
		cost += len(blob.Data)
	}
	return cost
}
