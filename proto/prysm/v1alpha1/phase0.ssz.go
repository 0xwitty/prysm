package eth

import (
	"fmt"
	ssz "github.com/prysmaticlabs/fastssz"
	go_bitfield "github.com/prysmaticlabs/go-bitfield"
)

func (c *AggregateAttestationAndProof) SizeSSZ() int {
	size := 108
	if c.Aggregate == nil {
		c.Aggregate = new(Attestation)
	}
	size += c.Aggregate.SizeSSZ()
	return size
}

func (c *AggregateAttestationAndProof) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *AggregateAttestationAndProof) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 108

	// Field 0: AggregatorIndex
	if dst, err = c.AggregatorIndex.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Aggregate
	if c.Aggregate == nil {
		c.Aggregate = new(Attestation)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Aggregate.SizeSSZ()

	// Field 2: SelectionProof
	if len(c.SelectionProof) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.SelectionProof...)

	// Field 1: Aggregate
	if dst, err = c.Aggregate.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}

func (c *AggregateAttestationAndProof) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 108 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]    // c.AggregatorIndex
	sszSlice2 := buf[12:108] // c.SelectionProof

	sszVarOffset1 := ssz.ReadOffset(buf[8:12]) // c.Aggregate
	if sszVarOffset1 < 108 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset1 > size {
		return ssz.ErrOffset
	}
	sszSlice1 := buf[sszVarOffset1:] // c.Aggregate

	// Field 0: AggregatorIndex
	if err = c.AggregatorIndex.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Aggregate
	c.Aggregate = new(Attestation)
	if err = c.Aggregate.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: SelectionProof
	c.SelectionProof = make([]byte, 0, 96)
	c.SelectionProof = append(c.SelectionProof, sszSlice2...)
	return err
}

func (c *AggregateAttestationAndProof) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *AggregateAttestationAndProof) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: AggregatorIndex
	if hash, err := c.AggregatorIndex.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 1: Aggregate
	if err := c.Aggregate.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 2: SelectionProof
	if len(c.SelectionProof) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.SelectionProof)
	hh.Merkleize(indx)
	return nil
}

func (c *Attestation) SizeSSZ() int {
	size := 228
	size += len(c.AggregationBits)
	return size
}

func (c *Attestation) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Attestation) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 228

	// Field 0: AggregationBits
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.AggregationBits)

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(AttestationData)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: AggregationBits
	if len(c.AggregationBits) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	dst = append(dst, c.AggregationBits...)
	return dst, err
}

func (c *Attestation) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 228 {
		return ssz.ErrSize
	}

	sszSlice1 := buf[4:132]   // c.Data
	sszSlice2 := buf[132:228] // c.Signature

	sszVarOffset0 := ssz.ReadOffset(buf[0:4]) // c.AggregationBits
	if sszVarOffset0 < 228 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset0 > size {
		return ssz.ErrOffset
	}
	sszSlice0 := buf[sszVarOffset0:] // c.AggregationBits

	// Field 0: AggregationBits
	if err = ssz.ValidateBitlist(sszSlice0, 2048); err != nil {
		return err
	}
	c.AggregationBits = append([]byte{}, go_bitfield.Bitlist(sszSlice0)...)

	// Field 1: Data
	c.Data = new(AttestationData)
	if err = c.Data.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice2...)
	return err
}

func (c *Attestation) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Attestation) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: AggregationBits
	if len(c.AggregationBits) == 0 {
		return ssz.ErrEmptyBitlist
	}
	hh.PutBitlist(c.AggregationBits, 2048)
	// Field 1: Data
	if err := c.Data.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 2: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *AttestationData) SizeSSZ() int {
	size := 128

	return size
}

func (c *AttestationData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *AttestationData) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Slot
	if dst, err = c.Slot.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: CommitteeIndex
	if dst, err = c.CommitteeIndex.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: BeaconBlockRoot
	if len(c.BeaconBlockRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BeaconBlockRoot...)

	// Field 3: Source
	if c.Source == nil {
		c.Source = new(Checkpoint)
	}
	if dst, err = c.Source.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 4: Target
	if c.Target == nil {
		c.Target = new(Checkpoint)
	}
	if dst, err = c.Target.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *AttestationData) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 128 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]    // c.Slot
	sszSlice1 := buf[8:16]   // c.CommitteeIndex
	sszSlice2 := buf[16:48]  // c.BeaconBlockRoot
	sszSlice3 := buf[48:88]  // c.Source
	sszSlice4 := buf[88:128] // c.Target

	// Field 0: Slot
	if err = c.Slot.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: CommitteeIndex
	if err = c.CommitteeIndex.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: BeaconBlockRoot
	c.BeaconBlockRoot = make([]byte, 0, 32)
	c.BeaconBlockRoot = append(c.BeaconBlockRoot, sszSlice2...)

	// Field 3: Source
	c.Source = new(Checkpoint)
	if err = c.Source.UnmarshalSSZ(sszSlice3); err != nil {
		return err
	}

	// Field 4: Target
	c.Target = new(Checkpoint)
	if err = c.Target.UnmarshalSSZ(sszSlice4); err != nil {
		return err
	}
	return err
}

func (c *AttestationData) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *AttestationData) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Slot
	if hash, err := c.Slot.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 1: CommitteeIndex
	if hash, err := c.CommitteeIndex.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 2: BeaconBlockRoot
	if len(c.BeaconBlockRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.BeaconBlockRoot)
	// Field 3: Source
	if err := c.Source.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 4: Target
	if err := c.Target.HashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}

func (c *AttesterSlashing) SizeSSZ() int {
	size := 8
	if c.Attestation_1 == nil {
		c.Attestation_1 = new(IndexedAttestation)
	}
	size += c.Attestation_1.SizeSSZ()
	if c.Attestation_2 == nil {
		c.Attestation_2 = new(IndexedAttestation)
	}
	size += c.Attestation_2.SizeSSZ()
	return size
}

func (c *AttesterSlashing) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *AttesterSlashing) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 8

	// Field 0: Attestation_1
	if c.Attestation_1 == nil {
		c.Attestation_1 = new(IndexedAttestation)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Attestation_1.SizeSSZ()

	// Field 1: Attestation_2
	if c.Attestation_2 == nil {
		c.Attestation_2 = new(IndexedAttestation)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Attestation_2.SizeSSZ()

	// Field 0: Attestation_1
	if dst, err = c.Attestation_1.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Attestation_2
	if dst, err = c.Attestation_2.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}

func (c *AttesterSlashing) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 8 {
		return ssz.ErrSize
	}

	sszVarOffset0 := ssz.ReadOffset(buf[0:4]) // c.Attestation_1
	if sszVarOffset0 < 8 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset0 > size {
		return ssz.ErrOffset
	}
	sszVarOffset1 := ssz.ReadOffset(buf[4:8]) // c.Attestation_2
	if sszVarOffset1 > size || sszVarOffset1 < sszVarOffset0 {
		return ssz.ErrOffset
	}
	sszSlice0 := buf[sszVarOffset0:sszVarOffset1] // c.Attestation_1
	sszSlice1 := buf[sszVarOffset1:]              // c.Attestation_2

	// Field 0: Attestation_1
	c.Attestation_1 = new(IndexedAttestation)
	if err = c.Attestation_1.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Attestation_2
	c.Attestation_2 = new(IndexedAttestation)
	if err = c.Attestation_2.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}
	return err
}

func (c *AttesterSlashing) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *AttesterSlashing) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Attestation_1
	if err := c.Attestation_1.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 1: Attestation_2
	if err := c.Attestation_2.HashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}

func (c *BeaconBlock) SizeSSZ() int {
	size := 84
	if c.Body == nil {
		c.Body = new(BeaconBlockBody)
	}
	size += c.Body.SizeSSZ()
	return size
}

func (c *BeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *BeaconBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 84

	// Field 0: Slot
	if dst, err = c.Slot.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: ProposerIndex
	if dst, err = c.ProposerIndex.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: ParentRoot
	if len(c.ParentRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ParentRoot...)

	// Field 3: StateRoot
	if len(c.StateRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.StateRoot...)

	// Field 4: Body
	if c.Body == nil {
		c.Body = new(BeaconBlockBody)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Body.SizeSSZ()

	// Field 4: Body
	if dst, err = c.Body.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}

func (c *BeaconBlock) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 84 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]   // c.Slot
	sszSlice1 := buf[8:16]  // c.ProposerIndex
	sszSlice2 := buf[16:48] // c.ParentRoot
	sszSlice3 := buf[48:80] // c.StateRoot

	sszVarOffset4 := ssz.ReadOffset(buf[80:84]) // c.Body
	if sszVarOffset4 < 84 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset4 > size {
		return ssz.ErrOffset
	}
	sszSlice4 := buf[sszVarOffset4:] // c.Body

	// Field 0: Slot
	if err = c.Slot.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: ProposerIndex
	if err = c.ProposerIndex.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: ParentRoot
	c.ParentRoot = make([]byte, 0, 32)
	c.ParentRoot = append(c.ParentRoot, sszSlice2...)

	// Field 3: StateRoot
	c.StateRoot = make([]byte, 0, 32)
	c.StateRoot = append(c.StateRoot, sszSlice3...)

	// Field 4: Body
	c.Body = new(BeaconBlockBody)
	if err = c.Body.UnmarshalSSZ(sszSlice4); err != nil {
		return err
	}
	return err
}

func (c *BeaconBlock) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconBlock) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Slot
	if hash, err := c.Slot.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 1: ProposerIndex
	if hash, err := c.ProposerIndex.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 2: ParentRoot
	if len(c.ParentRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ParentRoot)
	// Field 3: StateRoot
	if len(c.StateRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.StateRoot)
	// Field 4: Body
	if err := c.Body.HashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}

func (c *BeaconBlockBody) SizeSSZ() int {
	size := 220
	size += len(c.ProposerSlashings) * 416
	size += func() int {
		s := 0
		for _, o := range c.AttesterSlashings {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()
	size += func() int {
		s := 0
		for _, o := range c.Attestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()
	size += len(c.Deposits) * 1240
	size += len(c.VoluntaryExits) * 112
	return size
}

func (c *BeaconBlockBody) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *BeaconBlockBody) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 220

	// Field 0: RandaoReveal
	if len(c.RandaoReveal) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.RandaoReveal...)

	// Field 1: Eth1Data
	if c.Eth1Data == nil {
		c.Eth1Data = new(Eth1Data)
	}
	if dst, err = c.Eth1Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: Graffiti
	if len(c.Graffiti) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Graffiti...)

	// Field 3: ProposerSlashings
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.ProposerSlashings) * 416

	// Field 4: AttesterSlashings
	dst = ssz.WriteOffset(dst, offset)
	offset += func() int {
		s := 0
		for _, o := range c.AttesterSlashings {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()

	// Field 5: Attestations
	dst = ssz.WriteOffset(dst, offset)
	offset += func() int {
		s := 0
		for _, o := range c.Attestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()

	// Field 6: Deposits
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Deposits) * 1240

	// Field 7: VoluntaryExits
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.VoluntaryExits) * 112

	// Field 3: ProposerSlashings
	if len(c.ProposerSlashings) > 16 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.ProposerSlashings {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 4: AttesterSlashings
	if len(c.AttesterSlashings) > 2 {
		return nil, ssz.ErrListTooBig
	}
	{
		offset = 4 * len(c.AttesterSlashings)
		for _, o := range c.AttesterSlashings {
			dst = ssz.WriteOffset(dst, offset)
			offset += o.SizeSSZ()
		}
	}
	for _, o := range c.AttesterSlashings {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 5: Attestations
	if len(c.Attestations) > 128 {
		return nil, ssz.ErrListTooBig
	}
	{
		offset = 4 * len(c.Attestations)
		for _, o := range c.Attestations {
			dst = ssz.WriteOffset(dst, offset)
			offset += o.SizeSSZ()
		}
	}
	for _, o := range c.Attestations {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 6: Deposits
	if len(c.Deposits) > 16 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Deposits {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 7: VoluntaryExits
	if len(c.VoluntaryExits) > 16 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.VoluntaryExits {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}
	return dst, err
}

func (c *BeaconBlockBody) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 220 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:96]    // c.RandaoReveal
	sszSlice1 := buf[96:168]  // c.Eth1Data
	sszSlice2 := buf[168:200] // c.Graffiti

	sszVarOffset3 := ssz.ReadOffset(buf[200:204]) // c.ProposerSlashings
	if sszVarOffset3 < 220 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset3 > size {
		return ssz.ErrOffset
	}
	sszVarOffset4 := ssz.ReadOffset(buf[204:208]) // c.AttesterSlashings
	if sszVarOffset4 > size || sszVarOffset4 < sszVarOffset3 {
		return ssz.ErrOffset
	}
	sszVarOffset5 := ssz.ReadOffset(buf[208:212]) // c.Attestations
	if sszVarOffset5 > size || sszVarOffset5 < sszVarOffset4 {
		return ssz.ErrOffset
	}
	sszVarOffset6 := ssz.ReadOffset(buf[212:216]) // c.Deposits
	if sszVarOffset6 > size || sszVarOffset6 < sszVarOffset5 {
		return ssz.ErrOffset
	}
	sszVarOffset7 := ssz.ReadOffset(buf[216:220]) // c.VoluntaryExits
	if sszVarOffset7 > size || sszVarOffset7 < sszVarOffset6 {
		return ssz.ErrOffset
	}
	sszSlice3 := buf[sszVarOffset3:sszVarOffset4] // c.ProposerSlashings
	sszSlice4 := buf[sszVarOffset4:sszVarOffset5] // c.AttesterSlashings
	sszSlice5 := buf[sszVarOffset5:sszVarOffset6] // c.Attestations
	sszSlice6 := buf[sszVarOffset6:sszVarOffset7] // c.Deposits
	sszSlice7 := buf[sszVarOffset7:]              // c.VoluntaryExits

	// Field 0: RandaoReveal
	c.RandaoReveal = make([]byte, 0, 96)
	c.RandaoReveal = append(c.RandaoReveal, sszSlice0...)

	// Field 1: Eth1Data
	c.Eth1Data = new(Eth1Data)
	if err = c.Eth1Data.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: Graffiti
	c.Graffiti = make([]byte, 0, 32)
	c.Graffiti = append(c.Graffiti, sszSlice2...)

	// Field 3: ProposerSlashings
	{
		if len(sszSlice3)%416 != 0 {
			return fmt.Errorf("misaligned bytes: c.ProposerSlashings length is %d, which is not a multiple of 416", len(sszSlice3))
		}
		numElem := len(sszSlice3) / 416
		if numElem > 16 {
			return fmt.Errorf("ssz-max exceeded: c.ProposerSlashings has %d elements, ssz-max is 16", numElem)
		}
		c.ProposerSlashings = make([]*ProposerSlashing, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *ProposerSlashing
			tmp = new(ProposerSlashing)
			tmpSlice := sszSlice3[i*416 : (1+i)*416]
			if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.ProposerSlashings[i] = tmp
		}
	}

	// Field 4: AttesterSlashings
	{
		// empty lists are zero length, so make sure there is room for an offset
		// before attempting to unmarshal it
		if len(sszSlice4) > 3 {
			firstOffset := ssz.ReadOffset(sszSlice4[0:4])
			if firstOffset%4 != 0 {
				return fmt.Errorf("misaligned list bytes: when decoding c.AttesterSlashings, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
			}
			listLen := firstOffset / 4
			if listLen > 2 {
				return fmt.Errorf("ssz-max exceeded: c.AttesterSlashings has %d elements, ssz-max is 2", listLen)
			}
			listOffsets := make([]uint64, listLen)
			for i := 0; uint64(i) < listLen; i++ {
				listOffsets[i] = ssz.ReadOffset(sszSlice4[i*4 : (i+1)*4])
			}
			c.AttesterSlashings = make([]*AttesterSlashing, len(listOffsets))
			for i := 0; i < len(listOffsets); i++ {
				var tmp *AttesterSlashing
				tmp = new(AttesterSlashing)
				var tmpSlice []byte
				if i+1 == len(listOffsets) {
					tmpSlice = sszSlice4[listOffsets[i]:]
				} else {
					tmpSlice = sszSlice4[listOffsets[i]:listOffsets[i+1]]
				}
				if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
					return err
				}
				c.AttesterSlashings[i] = tmp
			}
		}
	}

	// Field 5: Attestations
	{
		// empty lists are zero length, so make sure there is room for an offset
		// before attempting to unmarshal it
		if len(sszSlice5) > 3 {
			firstOffset := ssz.ReadOffset(sszSlice5[0:4])
			if firstOffset%4 != 0 {
				return fmt.Errorf("misaligned list bytes: when decoding c.Attestations, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
			}
			listLen := firstOffset / 4
			if listLen > 128 {
				return fmt.Errorf("ssz-max exceeded: c.Attestations has %d elements, ssz-max is 128", listLen)
			}
			listOffsets := make([]uint64, listLen)
			for i := 0; uint64(i) < listLen; i++ {
				listOffsets[i] = ssz.ReadOffset(sszSlice5[i*4 : (i+1)*4])
			}
			c.Attestations = make([]*Attestation, len(listOffsets))
			for i := 0; i < len(listOffsets); i++ {
				var tmp *Attestation
				tmp = new(Attestation)
				var tmpSlice []byte
				if i+1 == len(listOffsets) {
					tmpSlice = sszSlice5[listOffsets[i]:]
				} else {
					tmpSlice = sszSlice5[listOffsets[i]:listOffsets[i+1]]
				}
				if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
					return err
				}
				c.Attestations[i] = tmp
			}
		}
	}

	// Field 6: Deposits
	{
		if len(sszSlice6)%1240 != 0 {
			return fmt.Errorf("misaligned bytes: c.Deposits length is %d, which is not a multiple of 1240", len(sszSlice6))
		}
		numElem := len(sszSlice6) / 1240
		if numElem > 16 {
			return fmt.Errorf("ssz-max exceeded: c.Deposits has %d elements, ssz-max is 16", numElem)
		}
		c.Deposits = make([]*Deposit, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *Deposit
			tmp = new(Deposit)
			tmpSlice := sszSlice6[i*1240 : (1+i)*1240]
			if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Deposits[i] = tmp
		}
	}

	// Field 7: VoluntaryExits
	{
		if len(sszSlice7)%112 != 0 {
			return fmt.Errorf("misaligned bytes: c.VoluntaryExits length is %d, which is not a multiple of 112", len(sszSlice7))
		}
		numElem := len(sszSlice7) / 112
		if numElem > 16 {
			return fmt.Errorf("ssz-max exceeded: c.VoluntaryExits has %d elements, ssz-max is 16", numElem)
		}
		c.VoluntaryExits = make([]*SignedVoluntaryExit, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *SignedVoluntaryExit
			tmp = new(SignedVoluntaryExit)
			tmpSlice := sszSlice7[i*112 : (1+i)*112]
			if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.VoluntaryExits[i] = tmp
		}
	}
	return err
}

func (c *BeaconBlockBody) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconBlockBody) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: RandaoReveal
	if len(c.RandaoReveal) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.RandaoReveal)
	// Field 1: Eth1Data
	if err := c.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 2: Graffiti
	if len(c.Graffiti) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Graffiti)
	// Field 3: ProposerSlashings
	{
		if len(c.ProposerSlashings) > 16 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.ProposerSlashings {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.ProposerSlashings)), 16)
	}
	// Field 4: AttesterSlashings
	{
		if len(c.AttesterSlashings) > 2 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.AttesterSlashings {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.AttesterSlashings)), 2)
	}
	// Field 5: Attestations
	{
		if len(c.Attestations) > 128 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Attestations {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Attestations)), 128)
	}
	// Field 6: Deposits
	{
		if len(c.Deposits) > 16 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Deposits {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Deposits)), 16)
	}
	// Field 7: VoluntaryExits
	{
		if len(c.VoluntaryExits) > 16 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.VoluntaryExits {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.VoluntaryExits)), 16)
	}
	hh.Merkleize(indx)
	return nil
}

func (c *BeaconBlockHeader) SizeSSZ() int {
	size := 112

	return size
}

func (c *BeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *BeaconBlockHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Slot
	if dst, err = c.Slot.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: ProposerIndex
	if dst, err = c.ProposerIndex.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: ParentRoot
	if len(c.ParentRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ParentRoot...)

	// Field 3: StateRoot
	if len(c.StateRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.StateRoot...)

	// Field 4: BodyRoot
	if len(c.BodyRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BodyRoot...)

	return dst, err
}

func (c *BeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 112 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]    // c.Slot
	sszSlice1 := buf[8:16]   // c.ProposerIndex
	sszSlice2 := buf[16:48]  // c.ParentRoot
	sszSlice3 := buf[48:80]  // c.StateRoot
	sszSlice4 := buf[80:112] // c.BodyRoot

	// Field 0: Slot
	if err = c.Slot.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: ProposerIndex
	if err = c.ProposerIndex.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: ParentRoot
	c.ParentRoot = make([]byte, 0, 32)
	c.ParentRoot = append(c.ParentRoot, sszSlice2...)

	// Field 3: StateRoot
	c.StateRoot = make([]byte, 0, 32)
	c.StateRoot = append(c.StateRoot, sszSlice3...)

	// Field 4: BodyRoot
	c.BodyRoot = make([]byte, 0, 32)
	c.BodyRoot = append(c.BodyRoot, sszSlice4...)
	return err
}

func (c *BeaconBlockHeader) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconBlockHeader) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Slot
	if hash, err := c.Slot.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 1: ProposerIndex
	if hash, err := c.ProposerIndex.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 2: ParentRoot
	if len(c.ParentRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ParentRoot)
	// Field 3: StateRoot
	if len(c.StateRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.StateRoot)
	// Field 4: BodyRoot
	if len(c.BodyRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.BodyRoot)
	hh.Merkleize(indx)
	return nil
}

func (c *BeaconState) SizeSSZ() int {
	size := 2687377
	size += len(c.HistoricalRoots) * 32
	size += len(c.Eth1DataVotes) * 72
	size += len(c.Validators) * 121
	size += len(c.Balances) * 8
	size += func() int {
		s := 0
		for _, o := range c.PreviousEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()
	size += func() int {
		s := 0
		for _, o := range c.CurrentEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()
	return size
}

func (c *BeaconState) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *BeaconState) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 2687377

	// Field 0: GenesisTime
	dst = ssz.MarshalUint64(dst, c.GenesisTime)

	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.GenesisValidatorsRoot...)

	// Field 2: Slot
	if dst, err = c.Slot.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 3: Fork
	if c.Fork == nil {
		c.Fork = new(Fork)
	}
	if dst, err = c.Fork.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 4: LatestBlockHeader
	if c.LatestBlockHeader == nil {
		c.LatestBlockHeader = new(BeaconBlockHeader)
	}
	if dst, err = c.LatestBlockHeader.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 5: BlockRoots
	if len(c.BlockRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.BlockRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 6: StateRoots
	if len(c.StateRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.StateRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 7: HistoricalRoots
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.HistoricalRoots) * 32

	// Field 8: Eth1Data
	if c.Eth1Data == nil {
		c.Eth1Data = new(Eth1Data)
	}
	if dst, err = c.Eth1Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 9: Eth1DataVotes
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Eth1DataVotes) * 72

	// Field 10: Eth1DepositIndex
	dst = ssz.MarshalUint64(dst, c.Eth1DepositIndex)

	// Field 11: Validators
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Validators) * 121

	// Field 12: Balances
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.Balances) * 8

	// Field 13: RandaoMixes
	if len(c.RandaoMixes) != 65536 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.RandaoMixes {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 14: Slashings
	if len(c.Slashings) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.Slashings {
		dst = ssz.MarshalUint64(dst, o)
	}

	// Field 15: PreviousEpochAttestations
	dst = ssz.WriteOffset(dst, offset)
	offset += func() int {
		s := 0
		for _, o := range c.PreviousEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()

	// Field 16: CurrentEpochAttestations
	dst = ssz.WriteOffset(dst, offset)
	offset += func() int {
		s := 0
		for _, o := range c.CurrentEpochAttestations {
			s += 4
			s += o.SizeSSZ()
		}
		return s
	}()

	// Field 17: JustificationBits
	if len([]byte(c.JustificationBits)) != 1 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, []byte(c.JustificationBits)...)

	// Field 18: PreviousJustifiedCheckpoint
	if c.PreviousJustifiedCheckpoint == nil {
		c.PreviousJustifiedCheckpoint = new(Checkpoint)
	}
	if dst, err = c.PreviousJustifiedCheckpoint.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 19: CurrentJustifiedCheckpoint
	if c.CurrentJustifiedCheckpoint == nil {
		c.CurrentJustifiedCheckpoint = new(Checkpoint)
	}
	if dst, err = c.CurrentJustifiedCheckpoint.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 20: FinalizedCheckpoint
	if c.FinalizedCheckpoint == nil {
		c.FinalizedCheckpoint = new(Checkpoint)
	}
	if dst, err = c.FinalizedCheckpoint.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 7: HistoricalRoots
	if len(c.HistoricalRoots) > 16777216 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.HistoricalRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 9: Eth1DataVotes
	if len(c.Eth1DataVotes) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Eth1DataVotes {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 11: Validators
	if len(c.Validators) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Validators {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 12: Balances
	if len(c.Balances) > 1099511627776 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.Balances {
		dst = ssz.MarshalUint64(dst, o)
	}

	// Field 15: PreviousEpochAttestations
	if len(c.PreviousEpochAttestations) > 4096 {
		return nil, ssz.ErrListTooBig
	}
	{
		offset = 4 * len(c.PreviousEpochAttestations)
		for _, o := range c.PreviousEpochAttestations {
			dst = ssz.WriteOffset(dst, offset)
			offset += o.SizeSSZ()
		}
	}
	for _, o := range c.PreviousEpochAttestations {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}

	// Field 16: CurrentEpochAttestations
	if len(c.CurrentEpochAttestations) > 4096 {
		return nil, ssz.ErrListTooBig
	}
	{
		offset = 4 * len(c.CurrentEpochAttestations)
		for _, o := range c.CurrentEpochAttestations {
			dst = ssz.WriteOffset(dst, offset)
			offset += o.SizeSSZ()
		}
	}
	for _, o := range c.CurrentEpochAttestations {
		if dst, err = o.MarshalSSZTo(dst); err != nil {
			return nil, err
		}
	}
	return dst, err
}

func (c *BeaconState) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 2687377 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]              // c.GenesisTime
	sszSlice1 := buf[8:40]             // c.GenesisValidatorsRoot
	sszSlice2 := buf[40:48]            // c.Slot
	sszSlice3 := buf[48:64]            // c.Fork
	sszSlice4 := buf[64:176]           // c.LatestBlockHeader
	sszSlice5 := buf[176:262320]       // c.BlockRoots
	sszSlice6 := buf[262320:524464]    // c.StateRoots
	sszSlice8 := buf[524468:524540]    // c.Eth1Data
	sszSlice10 := buf[524544:524552]   // c.Eth1DepositIndex
	sszSlice13 := buf[524560:2621712]  // c.RandaoMixes
	sszSlice14 := buf[2621712:2687248] // c.Slashings
	sszSlice17 := buf[2687256:2687257] // c.JustificationBits
	sszSlice18 := buf[2687257:2687297] // c.PreviousJustifiedCheckpoint
	sszSlice19 := buf[2687297:2687337] // c.CurrentJustifiedCheckpoint
	sszSlice20 := buf[2687337:2687377] // c.FinalizedCheckpoint

	sszVarOffset7 := ssz.ReadOffset(buf[524464:524468]) // c.HistoricalRoots
	if sszVarOffset7 < 2687377 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset7 > size {
		return ssz.ErrOffset
	}
	sszVarOffset9 := ssz.ReadOffset(buf[524540:524544]) // c.Eth1DataVotes
	if sszVarOffset9 > size || sszVarOffset9 < sszVarOffset7 {
		return ssz.ErrOffset
	}
	sszVarOffset11 := ssz.ReadOffset(buf[524552:524556]) // c.Validators
	if sszVarOffset11 > size || sszVarOffset11 < sszVarOffset9 {
		return ssz.ErrOffset
	}
	sszVarOffset12 := ssz.ReadOffset(buf[524556:524560]) // c.Balances
	if sszVarOffset12 > size || sszVarOffset12 < sszVarOffset11 {
		return ssz.ErrOffset
	}
	sszVarOffset15 := ssz.ReadOffset(buf[2687248:2687252]) // c.PreviousEpochAttestations
	if sszVarOffset15 > size || sszVarOffset15 < sszVarOffset12 {
		return ssz.ErrOffset
	}
	sszVarOffset16 := ssz.ReadOffset(buf[2687252:2687256]) // c.CurrentEpochAttestations
	if sszVarOffset16 > size || sszVarOffset16 < sszVarOffset15 {
		return ssz.ErrOffset
	}
	sszSlice7 := buf[sszVarOffset7:sszVarOffset9]    // c.HistoricalRoots
	sszSlice9 := buf[sszVarOffset9:sszVarOffset11]   // c.Eth1DataVotes
	sszSlice11 := buf[sszVarOffset11:sszVarOffset12] // c.Validators
	sszSlice12 := buf[sszVarOffset12:sszVarOffset15] // c.Balances
	sszSlice15 := buf[sszVarOffset15:sszVarOffset16] // c.PreviousEpochAttestations
	sszSlice16 := buf[sszVarOffset16:]               // c.CurrentEpochAttestations

	// Field 0: GenesisTime
	c.GenesisTime = ssz.UnmarshallUint64(sszSlice0)

	// Field 1: GenesisValidatorsRoot
	c.GenesisValidatorsRoot = make([]byte, 0, 32)
	c.GenesisValidatorsRoot = append(c.GenesisValidatorsRoot, sszSlice1...)

	// Field 2: Slot
	if err = c.Slot.UnmarshalSSZ(sszSlice2); err != nil {
		return err
	}

	// Field 3: Fork
	c.Fork = new(Fork)
	if err = c.Fork.UnmarshalSSZ(sszSlice3); err != nil {
		return err
	}

	// Field 4: LatestBlockHeader
	c.LatestBlockHeader = new(BeaconBlockHeader)
	if err = c.LatestBlockHeader.UnmarshalSSZ(sszSlice4); err != nil {
		return err
	}

	// Field 5: BlockRoots
	{
		var tmp []byte
		c.BlockRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := sszSlice5[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.BlockRoots[i] = tmp
		}
	}

	// Field 6: StateRoots
	{
		var tmp []byte
		c.StateRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := sszSlice6[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.StateRoots[i] = tmp
		}
	}

	// Field 7: HistoricalRoots
	{
		if len(sszSlice7)%32 != 0 {
			return fmt.Errorf("misaligned bytes: c.HistoricalRoots length is %d, which is not a multiple of 32", len(sszSlice7))
		}
		numElem := len(sszSlice7) / 32
		if numElem > 16777216 {
			return fmt.Errorf("ssz-max exceeded: c.HistoricalRoots has %d elements, ssz-max is 16777216", numElem)
		}
		c.HistoricalRoots = make([][]byte, numElem)
		for i := 0; i < numElem; i++ {
			var tmp []byte

			tmpSlice := sszSlice7[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.HistoricalRoots[i] = tmp
		}
	}

	// Field 8: Eth1Data
	c.Eth1Data = new(Eth1Data)
	if err = c.Eth1Data.UnmarshalSSZ(sszSlice8); err != nil {
		return err
	}

	// Field 9: Eth1DataVotes
	{
		if len(sszSlice9)%72 != 0 {
			return fmt.Errorf("misaligned bytes: c.Eth1DataVotes length is %d, which is not a multiple of 72", len(sszSlice9))
		}
		numElem := len(sszSlice9) / 72
		if numElem > 2048 {
			return fmt.Errorf("ssz-max exceeded: c.Eth1DataVotes has %d elements, ssz-max is 2048", numElem)
		}
		c.Eth1DataVotes = make([]*Eth1Data, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *Eth1Data
			tmp = new(Eth1Data)
			tmpSlice := sszSlice9[i*72 : (1+i)*72]
			if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Eth1DataVotes[i] = tmp
		}
	}

	// Field 10: Eth1DepositIndex
	c.Eth1DepositIndex = ssz.UnmarshallUint64(sszSlice10)

	// Field 11: Validators
	{
		if len(sszSlice11)%121 != 0 {
			return fmt.Errorf("misaligned bytes: c.Validators length is %d, which is not a multiple of 121", len(sszSlice11))
		}
		numElem := len(sszSlice11) / 121
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.Validators has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.Validators = make([]*Validator, numElem)
		for i := 0; i < numElem; i++ {
			var tmp *Validator
			tmp = new(Validator)
			tmpSlice := sszSlice11[i*121 : (1+i)*121]
			if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
				return err
			}
			c.Validators[i] = tmp
		}
	}

	// Field 12: Balances
	{
		if len(sszSlice12)%8 != 0 {
			return fmt.Errorf("misaligned bytes: c.Balances length is %d, which is not a multiple of 8", len(sszSlice12))
		}
		numElem := len(sszSlice12) / 8
		if numElem > 1099511627776 {
			return fmt.Errorf("ssz-max exceeded: c.Balances has %d elements, ssz-max is 1099511627776", numElem)
		}
		c.Balances = make([]uint64, numElem)
		for i := 0; i < numElem; i++ {
			var tmp uint64

			tmpSlice := sszSlice12[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.Balances[i] = tmp
		}
	}

	// Field 13: RandaoMixes
	{
		var tmp []byte
		c.RandaoMixes = make([][]byte, 65536)
		for i := 0; i < 65536; i++ {
			tmpSlice := sszSlice13[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.RandaoMixes[i] = tmp
		}
	}

	// Field 14: Slashings
	{
		var tmp uint64
		c.Slashings = make([]uint64, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := sszSlice14[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.Slashings[i] = tmp
		}
	}

	// Field 15: PreviousEpochAttestations
	{
		// empty lists are zero length, so make sure there is room for an offset
		// before attempting to unmarshal it
		if len(sszSlice15) > 3 {
			firstOffset := ssz.ReadOffset(sszSlice15[0:4])
			if firstOffset%4 != 0 {
				return fmt.Errorf("misaligned list bytes: when decoding c.PreviousEpochAttestations, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
			}
			listLen := firstOffset / 4
			if listLen > 4096 {
				return fmt.Errorf("ssz-max exceeded: c.PreviousEpochAttestations has %d elements, ssz-max is 4096", listLen)
			}
			listOffsets := make([]uint64, listLen)
			for i := 0; uint64(i) < listLen; i++ {
				listOffsets[i] = ssz.ReadOffset(sszSlice15[i*4 : (i+1)*4])
			}
			c.PreviousEpochAttestations = make([]*PendingAttestation, len(listOffsets))
			for i := 0; i < len(listOffsets); i++ {
				var tmp *PendingAttestation
				tmp = new(PendingAttestation)
				var tmpSlice []byte
				if i+1 == len(listOffsets) {
					tmpSlice = sszSlice15[listOffsets[i]:]
				} else {
					tmpSlice = sszSlice15[listOffsets[i]:listOffsets[i+1]]
				}
				if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
					return err
				}
				c.PreviousEpochAttestations[i] = tmp
			}
		}
	}

	// Field 16: CurrentEpochAttestations
	{
		// empty lists are zero length, so make sure there is room for an offset
		// before attempting to unmarshal it
		if len(sszSlice16) > 3 {
			firstOffset := ssz.ReadOffset(sszSlice16[0:4])
			if firstOffset%4 != 0 {
				return fmt.Errorf("misaligned list bytes: when decoding c.CurrentEpochAttestations, end-of-list offset is %d, which is not a multiple of 4 (offset size)", firstOffset)
			}
			listLen := firstOffset / 4
			if listLen > 4096 {
				return fmt.Errorf("ssz-max exceeded: c.CurrentEpochAttestations has %d elements, ssz-max is 4096", listLen)
			}
			listOffsets := make([]uint64, listLen)
			for i := 0; uint64(i) < listLen; i++ {
				listOffsets[i] = ssz.ReadOffset(sszSlice16[i*4 : (i+1)*4])
			}
			c.CurrentEpochAttestations = make([]*PendingAttestation, len(listOffsets))
			for i := 0; i < len(listOffsets); i++ {
				var tmp *PendingAttestation
				tmp = new(PendingAttestation)
				var tmpSlice []byte
				if i+1 == len(listOffsets) {
					tmpSlice = sszSlice16[listOffsets[i]:]
				} else {
					tmpSlice = sszSlice16[listOffsets[i]:listOffsets[i+1]]
				}
				if err = tmp.UnmarshalSSZ(tmpSlice); err != nil {
					return err
				}
				c.CurrentEpochAttestations[i] = tmp
			}
		}
	}

	// Field 17: JustificationBits
	c.JustificationBits = make([]byte, 0, 1)
	c.JustificationBits = append(c.JustificationBits, go_bitfield.Bitvector4(sszSlice17)...)

	// Field 18: PreviousJustifiedCheckpoint
	c.PreviousJustifiedCheckpoint = new(Checkpoint)
	if err = c.PreviousJustifiedCheckpoint.UnmarshalSSZ(sszSlice18); err != nil {
		return err
	}

	// Field 19: CurrentJustifiedCheckpoint
	c.CurrentJustifiedCheckpoint = new(Checkpoint)
	if err = c.CurrentJustifiedCheckpoint.UnmarshalSSZ(sszSlice19); err != nil {
		return err
	}

	// Field 20: FinalizedCheckpoint
	c.FinalizedCheckpoint = new(Checkpoint)
	if err = c.FinalizedCheckpoint.UnmarshalSSZ(sszSlice20); err != nil {
		return err
	}
	return err
}

func (c *BeaconState) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *BeaconState) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: GenesisTime
	hh.PutUint64(c.GenesisTime)
	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.GenesisValidatorsRoot)
	// Field 2: Slot
	if hash, err := c.Slot.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 3: Fork
	if err := c.Fork.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 4: LatestBlockHeader
	if err := c.LatestBlockHeader.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 5: BlockRoots
	{
		if len(c.BlockRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.BlockRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 6: StateRoots
	{
		if len(c.StateRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.StateRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 7: HistoricalRoots
	{
		if len(c.HistoricalRoots) > 16777216 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.HistoricalRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.HistoricalRoots)), 16777216)
	}
	// Field 8: Eth1Data
	if err := c.Eth1Data.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 9: Eth1DataVotes
	{
		if len(c.Eth1DataVotes) > 2048 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Eth1DataVotes {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Eth1DataVotes)), 2048)
	}
	// Field 10: Eth1DepositIndex
	hh.PutUint64(c.Eth1DepositIndex)
	// Field 11: Validators
	{
		if len(c.Validators) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Validators {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.Validators)), 1099511627776)
	}
	// Field 12: Balances
	{
		if len(c.Balances) > 1099511627776 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.Balances {
			hh.AppendUint64(o)
		}
		hh.FillUpTo32()
		numItems := uint64(len(c.Balances))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(1099511627776, numItems, 8))
	}
	// Field 13: RandaoMixes
	{
		if len(c.RandaoMixes) != 65536 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.RandaoMixes {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 14: Slashings
	{
		if len(c.Slashings) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.Slashings {
			hh.AppendUint64(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 15: PreviousEpochAttestations
	{
		if len(c.PreviousEpochAttestations) > 4096 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.PreviousEpochAttestations {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.PreviousEpochAttestations)), 4096)
	}
	// Field 16: CurrentEpochAttestations
	{
		if len(c.CurrentEpochAttestations) > 4096 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.CurrentEpochAttestations {
			if err := o.HashTreeRootWith(hh); err != nil {
				return err
			}
		}
		hh.MerkleizeWithMixin(subIndx, uint64(len(c.CurrentEpochAttestations)), 4096)
	}
	// Field 17: JustificationBits
	if len([]byte(c.JustificationBits)) != 1 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes([]byte(c.JustificationBits))
	// Field 18: PreviousJustifiedCheckpoint
	if err := c.PreviousJustifiedCheckpoint.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 19: CurrentJustifiedCheckpoint
	if err := c.CurrentJustifiedCheckpoint.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 20: FinalizedCheckpoint
	if err := c.FinalizedCheckpoint.HashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}

func (c *Checkpoint) SizeSSZ() int {
	size := 40

	return size
}

func (c *Checkpoint) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Checkpoint) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Epoch
	if dst, err = c.Epoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Root
	if len(c.Root) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Root...)

	return dst, err
}

func (c *Checkpoint) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 40 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]  // c.Epoch
	sszSlice1 := buf[8:40] // c.Root

	// Field 0: Epoch
	if err = c.Epoch.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Root
	c.Root = make([]byte, 0, 32)
	c.Root = append(c.Root, sszSlice1...)
	return err
}

func (c *Checkpoint) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Checkpoint) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Epoch
	if hash, err := c.Epoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 1: Root
	if len(c.Root) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Root)
	hh.Merkleize(indx)
	return nil
}

func (c *Deposit) SizeSSZ() int {
	size := 1240

	return size
}

func (c *Deposit) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Deposit) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Proof
	if len(c.Proof) != 33 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.Proof {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(Deposit_Data)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *Deposit) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 1240 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:1056]    // c.Proof
	sszSlice1 := buf[1056:1240] // c.Data

	// Field 0: Proof
	{
		var tmp []byte
		c.Proof = make([][]byte, 33)
		for i := 0; i < 33; i++ {
			tmpSlice := sszSlice0[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.Proof[i] = tmp
		}
	}

	// Field 1: Data
	c.Data = new(Deposit_Data)
	if err = c.Data.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}
	return err
}

func (c *Deposit) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Deposit) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Proof
	{
		if len(c.Proof) != 33 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.Proof {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 1: Data
	if err := c.Data.HashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}

func (c *Deposit_Data) SizeSSZ() int {
	size := 184

	return size
}

func (c *Deposit_Data) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Deposit_Data) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.PublicKey...)

	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.WithdrawalCredentials...)

	// Field 2: Amount
	dst = ssz.MarshalUint64(dst, c.Amount)

	// Field 3: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	return dst, err
}

func (c *Deposit_Data) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 184 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:48]   // c.PublicKey
	sszSlice1 := buf[48:80]  // c.WithdrawalCredentials
	sszSlice2 := buf[80:88]  // c.Amount
	sszSlice3 := buf[88:184] // c.Signature

	// Field 0: PublicKey
	c.PublicKey = make([]byte, 0, 48)
	c.PublicKey = append(c.PublicKey, sszSlice0...)

	// Field 1: WithdrawalCredentials
	c.WithdrawalCredentials = make([]byte, 0, 32)
	c.WithdrawalCredentials = append(c.WithdrawalCredentials, sszSlice1...)

	// Field 2: Amount
	c.Amount = ssz.UnmarshallUint64(sszSlice2)

	// Field 3: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice3...)
	return err
}

func (c *Deposit_Data) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Deposit_Data) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PublicKey)
	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.WithdrawalCredentials)
	// Field 2: Amount
	hh.PutUint64(c.Amount)
	// Field 3: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *DepositMessage) SizeSSZ() int {
	size := 88

	return size
}

func (c *DepositMessage) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *DepositMessage) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.PublicKey...)

	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.WithdrawalCredentials...)

	// Field 2: Amount
	dst = ssz.MarshalUint64(dst, c.Amount)

	return dst, err
}

func (c *DepositMessage) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 88 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:48]  // c.PublicKey
	sszSlice1 := buf[48:80] // c.WithdrawalCredentials
	sszSlice2 := buf[80:88] // c.Amount

	// Field 0: PublicKey
	c.PublicKey = make([]byte, 0, 48)
	c.PublicKey = append(c.PublicKey, sszSlice0...)

	// Field 1: WithdrawalCredentials
	c.WithdrawalCredentials = make([]byte, 0, 32)
	c.WithdrawalCredentials = append(c.WithdrawalCredentials, sszSlice1...)

	// Field 2: Amount
	c.Amount = ssz.UnmarshallUint64(sszSlice2)
	return err
}

func (c *DepositMessage) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *DepositMessage) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PublicKey)
	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.WithdrawalCredentials)
	// Field 2: Amount
	hh.PutUint64(c.Amount)
	hh.Merkleize(indx)
	return nil
}

func (c *ENRForkID) SizeSSZ() int {
	size := 16

	return size
}

func (c *ENRForkID) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *ENRForkID) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: CurrentForkDigest
	if len(c.CurrentForkDigest) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.CurrentForkDigest...)

	// Field 1: NextForkVersion
	if len(c.NextForkVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.NextForkVersion...)

	// Field 2: NextForkEpoch
	if dst, err = c.NextForkEpoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *ENRForkID) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:4]  // c.CurrentForkDigest
	sszSlice1 := buf[4:8]  // c.NextForkVersion
	sszSlice2 := buf[8:16] // c.NextForkEpoch

	// Field 0: CurrentForkDigest
	c.CurrentForkDigest = make([]byte, 0, 4)
	c.CurrentForkDigest = append(c.CurrentForkDigest, sszSlice0...)

	// Field 1: NextForkVersion
	c.NextForkVersion = make([]byte, 0, 4)
	c.NextForkVersion = append(c.NextForkVersion, sszSlice1...)

	// Field 2: NextForkEpoch
	if err = c.NextForkEpoch.UnmarshalSSZ(sszSlice2); err != nil {
		return err
	}
	return err
}

func (c *ENRForkID) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *ENRForkID) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: CurrentForkDigest
	if len(c.CurrentForkDigest) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentForkDigest)
	// Field 1: NextForkVersion
	if len(c.NextForkVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.NextForkVersion)
	// Field 2: NextForkEpoch
	if hash, err := c.NextForkEpoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	hh.Merkleize(indx)
	return nil
}

func (c *Eth1Data) SizeSSZ() int {
	size := 72

	return size
}

func (c *Eth1Data) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Eth1Data) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: DepositRoot
	if len(c.DepositRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.DepositRoot...)

	// Field 1: DepositCount
	dst = ssz.MarshalUint64(dst, c.DepositCount)

	// Field 2: BlockHash
	if len(c.BlockHash) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BlockHash...)

	return dst, err
}

func (c *Eth1Data) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 72 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:32]  // c.DepositRoot
	sszSlice1 := buf[32:40] // c.DepositCount
	sszSlice2 := buf[40:72] // c.BlockHash

	// Field 0: DepositRoot
	c.DepositRoot = make([]byte, 0, 32)
	c.DepositRoot = append(c.DepositRoot, sszSlice0...)

	// Field 1: DepositCount
	c.DepositCount = ssz.UnmarshallUint64(sszSlice1)

	// Field 2: BlockHash
	c.BlockHash = make([]byte, 0, 32)
	c.BlockHash = append(c.BlockHash, sszSlice2...)
	return err
}

func (c *Eth1Data) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Eth1Data) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: DepositRoot
	if len(c.DepositRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.DepositRoot)
	// Field 1: DepositCount
	hh.PutUint64(c.DepositCount)
	// Field 2: BlockHash
	if len(c.BlockHash) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.BlockHash)
	hh.Merkleize(indx)
	return nil
}

func (c *Fork) SizeSSZ() int {
	size := 16

	return size
}

func (c *Fork) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Fork) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: PreviousVersion
	if len(c.PreviousVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.PreviousVersion...)

	// Field 1: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.CurrentVersion...)

	// Field 2: Epoch
	if dst, err = c.Epoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *Fork) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:4]  // c.PreviousVersion
	sszSlice1 := buf[4:8]  // c.CurrentVersion
	sszSlice2 := buf[8:16] // c.Epoch

	// Field 0: PreviousVersion
	c.PreviousVersion = make([]byte, 0, 4)
	c.PreviousVersion = append(c.PreviousVersion, sszSlice0...)

	// Field 1: CurrentVersion
	c.CurrentVersion = make([]byte, 0, 4)
	c.CurrentVersion = append(c.CurrentVersion, sszSlice1...)

	// Field 2: Epoch
	if err = c.Epoch.UnmarshalSSZ(sszSlice2); err != nil {
		return err
	}
	return err
}

func (c *Fork) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Fork) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: PreviousVersion
	if len(c.PreviousVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PreviousVersion)
	// Field 1: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentVersion)
	// Field 2: Epoch
	if hash, err := c.Epoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	hh.Merkleize(indx)
	return nil
}

func (c *ForkData) SizeSSZ() int {
	size := 36

	return size
}

func (c *ForkData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *ForkData) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.CurrentVersion...)

	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.GenesisValidatorsRoot...)

	return dst, err
}

func (c *ForkData) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 36 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:4]  // c.CurrentVersion
	sszSlice1 := buf[4:36] // c.GenesisValidatorsRoot

	// Field 0: CurrentVersion
	c.CurrentVersion = make([]byte, 0, 4)
	c.CurrentVersion = append(c.CurrentVersion, sszSlice0...)

	// Field 1: GenesisValidatorsRoot
	c.GenesisValidatorsRoot = make([]byte, 0, 32)
	c.GenesisValidatorsRoot = append(c.GenesisValidatorsRoot, sszSlice1...)
	return err
}

func (c *ForkData) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *ForkData) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: CurrentVersion
	if len(c.CurrentVersion) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.CurrentVersion)
	// Field 1: GenesisValidatorsRoot
	if len(c.GenesisValidatorsRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.GenesisValidatorsRoot)
	hh.Merkleize(indx)
	return nil
}

func (c *HistoricalBatch) SizeSSZ() int {
	size := 524288

	return size
}

func (c *HistoricalBatch) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *HistoricalBatch) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: BlockRoots
	if len(c.BlockRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.BlockRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	// Field 1: StateRoots
	if len(c.StateRoots) != 8192 {
		return nil, ssz.ErrBytesLength
	}
	for _, o := range c.StateRoots {
		if len(o) != 32 {
			return nil, ssz.ErrBytesLength
		}
		dst = append(dst, o...)
	}

	return dst, err
}

func (c *HistoricalBatch) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 524288 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:262144]      // c.BlockRoots
	sszSlice1 := buf[262144:524288] // c.StateRoots

	// Field 0: BlockRoots
	{
		var tmp []byte
		c.BlockRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := sszSlice0[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.BlockRoots[i] = tmp
		}
	}

	// Field 1: StateRoots
	{
		var tmp []byte
		c.StateRoots = make([][]byte, 8192)
		for i := 0; i < 8192; i++ {
			tmpSlice := sszSlice1[i*32 : (1+i)*32]
			tmp = make([]byte, 0, 32)
			tmp = append(tmp, tmpSlice...)
			c.StateRoots[i] = tmp
		}
	}
	return err
}

func (c *HistoricalBatch) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *HistoricalBatch) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: BlockRoots
	{
		if len(c.BlockRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.BlockRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	// Field 1: StateRoots
	{
		if len(c.StateRoots) != 8192 {
			return ssz.ErrVectorLength
		}
		subIndx := hh.Index()
		for _, o := range c.StateRoots {
			if len(o) != 32 {
				return ssz.ErrBytesLength
			}
			hh.Append(o)
		}
		hh.Merkleize(subIndx)
	}
	hh.Merkleize(indx)
	return nil
}

func (c *IndexedAttestation) SizeSSZ() int {
	size := 228
	size += len(c.AttestingIndices) * 8
	return size
}

func (c *IndexedAttestation) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *IndexedAttestation) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 228

	// Field 0: AttestingIndices
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.AttestingIndices) * 8

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(AttestationData)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: AttestingIndices
	if len(c.AttestingIndices) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	for _, o := range c.AttestingIndices {
		dst = ssz.MarshalUint64(dst, o)
	}
	return dst, err
}

func (c *IndexedAttestation) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 228 {
		return ssz.ErrSize
	}

	sszSlice1 := buf[4:132]   // c.Data
	sszSlice2 := buf[132:228] // c.Signature

	sszVarOffset0 := ssz.ReadOffset(buf[0:4]) // c.AttestingIndices
	if sszVarOffset0 < 228 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset0 > size {
		return ssz.ErrOffset
	}
	sszSlice0 := buf[sszVarOffset0:] // c.AttestingIndices

	// Field 0: AttestingIndices
	{
		if len(sszSlice0)%8 != 0 {
			return fmt.Errorf("misaligned bytes: c.AttestingIndices length is %d, which is not a multiple of 8", len(sszSlice0))
		}
		numElem := len(sszSlice0) / 8
		if numElem > 2048 {
			return fmt.Errorf("ssz-max exceeded: c.AttestingIndices has %d elements, ssz-max is 2048", numElem)
		}
		c.AttestingIndices = make([]uint64, numElem)
		for i := 0; i < numElem; i++ {
			var tmp uint64

			tmpSlice := sszSlice0[i*8 : (1+i)*8]
			tmp = ssz.UnmarshallUint64(tmpSlice)
			c.AttestingIndices[i] = tmp
		}
	}

	// Field 1: Data
	c.Data = new(AttestationData)
	if err = c.Data.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice2...)
	return err
}

func (c *IndexedAttestation) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *IndexedAttestation) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: AttestingIndices
	{
		if len(c.AttestingIndices) > 2048 {
			return ssz.ErrListTooBig
		}
		subIndx := hh.Index()
		for _, o := range c.AttestingIndices {
			hh.AppendUint64(o)
		}
		hh.FillUpTo32()
		numItems := uint64(len(c.AttestingIndices))
		hh.MerkleizeWithMixin(subIndx, numItems, ssz.CalculateLimit(2048, numItems, 8))
	}
	// Field 1: Data
	if err := c.Data.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 2: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *PendingAttestation) SizeSSZ() int {
	size := 148
	size += len(c.AggregationBits)
	return size
}

func (c *PendingAttestation) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *PendingAttestation) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 148

	// Field 0: AggregationBits
	dst = ssz.WriteOffset(dst, offset)
	offset += len(c.AggregationBits)

	// Field 1: Data
	if c.Data == nil {
		c.Data = new(AttestationData)
	}
	if dst, err = c.Data.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 2: InclusionDelay
	if dst, err = c.InclusionDelay.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 3: ProposerIndex
	if dst, err = c.ProposerIndex.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 0: AggregationBits
	if len(c.AggregationBits) > 2048 {
		return nil, ssz.ErrListTooBig
	}
	dst = append(dst, c.AggregationBits...)
	return dst, err
}

func (c *PendingAttestation) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 148 {
		return ssz.ErrSize
	}

	sszSlice1 := buf[4:132]   // c.Data
	sszSlice2 := buf[132:140] // c.InclusionDelay
	sszSlice3 := buf[140:148] // c.ProposerIndex

	sszVarOffset0 := ssz.ReadOffset(buf[0:4]) // c.AggregationBits
	if sszVarOffset0 < 148 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset0 > size {
		return ssz.ErrOffset
	}
	sszSlice0 := buf[sszVarOffset0:] // c.AggregationBits

	// Field 0: AggregationBits
	if err = ssz.ValidateBitlist(sszSlice0, 2048); err != nil {
		return err
	}
	c.AggregationBits = append([]byte{}, go_bitfield.Bitlist(sszSlice0)...)

	// Field 1: Data
	c.Data = new(AttestationData)
	if err = c.Data.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}

	// Field 2: InclusionDelay
	if err = c.InclusionDelay.UnmarshalSSZ(sszSlice2); err != nil {
		return err
	}

	// Field 3: ProposerIndex
	if err = c.ProposerIndex.UnmarshalSSZ(sszSlice3); err != nil {
		return err
	}
	return err
}

func (c *PendingAttestation) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *PendingAttestation) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: AggregationBits
	if len(c.AggregationBits) == 0 {
		return ssz.ErrEmptyBitlist
	}
	hh.PutBitlist(c.AggregationBits, 2048)
	// Field 1: Data
	if err := c.Data.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 2: InclusionDelay
	if hash, err := c.InclusionDelay.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 3: ProposerIndex
	if hash, err := c.ProposerIndex.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	hh.Merkleize(indx)
	return nil
}

func (c *PowBlock) SizeSSZ() int {
	size := 96

	return size
}

func (c *PowBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *PowBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: BlockHash
	if len(c.BlockHash) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.BlockHash...)

	// Field 1: ParentHash
	if len(c.ParentHash) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ParentHash...)

	// Field 2: TotalDifficulty
	if len(c.TotalDifficulty) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.TotalDifficulty...)

	return dst, err
}

func (c *PowBlock) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 96 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:32]  // c.BlockHash
	sszSlice1 := buf[32:64] // c.ParentHash
	sszSlice2 := buf[64:96] // c.TotalDifficulty

	// Field 0: BlockHash
	c.BlockHash = make([]byte, 0, 32)
	c.BlockHash = append(c.BlockHash, sszSlice0...)

	// Field 1: ParentHash
	c.ParentHash = make([]byte, 0, 32)
	c.ParentHash = append(c.ParentHash, sszSlice1...)

	// Field 2: TotalDifficulty
	c.TotalDifficulty = make([]byte, 0, 32)
	c.TotalDifficulty = append(c.TotalDifficulty, sszSlice2...)
	return err
}

func (c *PowBlock) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *PowBlock) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: BlockHash
	if len(c.BlockHash) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.BlockHash)
	// Field 1: ParentHash
	if len(c.ParentHash) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ParentHash)
	// Field 2: TotalDifficulty
	if len(c.TotalDifficulty) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.TotalDifficulty)
	hh.Merkleize(indx)
	return nil
}

func (c *ProposerSlashing) SizeSSZ() int {
	size := 416

	return size
}

func (c *ProposerSlashing) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *ProposerSlashing) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Header_1
	if c.Header_1 == nil {
		c.Header_1 = new(SignedBeaconBlockHeader)
	}
	if dst, err = c.Header_1.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Header_2
	if c.Header_2 == nil {
		c.Header_2 = new(SignedBeaconBlockHeader)
	}
	if dst, err = c.Header_2.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *ProposerSlashing) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 416 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:208]   // c.Header_1
	sszSlice1 := buf[208:416] // c.Header_2

	// Field 0: Header_1
	c.Header_1 = new(SignedBeaconBlockHeader)
	if err = c.Header_1.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Header_2
	c.Header_2 = new(SignedBeaconBlockHeader)
	if err = c.Header_2.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}
	return err
}

func (c *ProposerSlashing) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *ProposerSlashing) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Header_1
	if err := c.Header_1.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 1: Header_2
	if err := c.Header_2.HashTreeRootWith(hh); err != nil {
		return err
	}
	hh.Merkleize(indx)
	return nil
}

func (c *SignedAggregateAttestationAndProof) SizeSSZ() int {
	size := 100
	if c.Message == nil {
		c.Message = new(AggregateAttestationAndProof)
	}
	size += c.Message.SizeSSZ()
	return size
}

func (c *SignedAggregateAttestationAndProof) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *SignedAggregateAttestationAndProof) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 100

	// Field 0: Message
	if c.Message == nil {
		c.Message = new(AggregateAttestationAndProof)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Message.SizeSSZ()

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: Message
	if dst, err = c.Message.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}

func (c *SignedAggregateAttestationAndProof) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 100 {
		return ssz.ErrSize
	}

	sszSlice1 := buf[4:100] // c.Signature

	sszVarOffset0 := ssz.ReadOffset(buf[0:4]) // c.Message
	if sszVarOffset0 < 100 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset0 > size {
		return ssz.ErrOffset
	}
	sszSlice0 := buf[sszVarOffset0:] // c.Message

	// Field 0: Message
	c.Message = new(AggregateAttestationAndProof)
	if err = c.Message.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice1...)
	return err
}

func (c *SignedAggregateAttestationAndProof) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SignedAggregateAttestationAndProof) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Message
	if err := c.Message.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 1: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *SignedBeaconBlock) SizeSSZ() int {
	size := 100
	if c.Block == nil {
		c.Block = new(BeaconBlock)
	}
	size += c.Block.SizeSSZ()
	return size
}

func (c *SignedBeaconBlock) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *SignedBeaconBlock) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error
	offset := 100

	// Field 0: Block
	if c.Block == nil {
		c.Block = new(BeaconBlock)
	}
	dst = ssz.WriteOffset(dst, offset)
	offset += c.Block.SizeSSZ()

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	// Field 0: Block
	if dst, err = c.Block.MarshalSSZTo(dst); err != nil {
		return nil, err
	}
	return dst, err
}

func (c *SignedBeaconBlock) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 100 {
		return ssz.ErrSize
	}

	sszSlice1 := buf[4:100] // c.Signature

	sszVarOffset0 := ssz.ReadOffset(buf[0:4]) // c.Block
	if sszVarOffset0 < 100 {
		return ssz.ErrInvalidVariableOffset
	}
	if sszVarOffset0 > size {
		return ssz.ErrOffset
	}
	sszSlice0 := buf[sszVarOffset0:] // c.Block

	// Field 0: Block
	c.Block = new(BeaconBlock)
	if err = c.Block.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice1...)
	return err
}

func (c *SignedBeaconBlock) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SignedBeaconBlock) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Block
	if err := c.Block.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 1: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *SignedBeaconBlockHeader) SizeSSZ() int {
	size := 208

	return size
}

func (c *SignedBeaconBlockHeader) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *SignedBeaconBlockHeader) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Header
	if c.Header == nil {
		c.Header = new(BeaconBlockHeader)
	}
	if dst, err = c.Header.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	return dst, err
}

func (c *SignedBeaconBlockHeader) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 208 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:112]   // c.Header
	sszSlice1 := buf[112:208] // c.Signature

	// Field 0: Header
	c.Header = new(BeaconBlockHeader)
	if err = c.Header.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice1...)
	return err
}

func (c *SignedBeaconBlockHeader) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SignedBeaconBlockHeader) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Header
	if err := c.Header.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 1: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *SignedVoluntaryExit) SizeSSZ() int {
	size := 112

	return size
}

func (c *SignedVoluntaryExit) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *SignedVoluntaryExit) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Exit
	if c.Exit == nil {
		c.Exit = new(VoluntaryExit)
	}
	if dst, err = c.Exit.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: Signature
	if len(c.Signature) != 96 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Signature...)

	return dst, err
}

func (c *SignedVoluntaryExit) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 112 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:16]   // c.Exit
	sszSlice1 := buf[16:112] // c.Signature

	// Field 0: Exit
	c.Exit = new(VoluntaryExit)
	if err = c.Exit.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: Signature
	c.Signature = make([]byte, 0, 96)
	c.Signature = append(c.Signature, sszSlice1...)
	return err
}

func (c *SignedVoluntaryExit) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SignedVoluntaryExit) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Exit
	if err := c.Exit.HashTreeRootWith(hh); err != nil {
		return err
	}
	// Field 1: Signature
	if len(c.Signature) != 96 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Signature)
	hh.Merkleize(indx)
	return nil
}

func (c *SigningData) SizeSSZ() int {
	size := 64

	return size
}

func (c *SigningData) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *SigningData) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: ObjectRoot
	if len(c.ObjectRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ObjectRoot...)

	// Field 1: Domain
	if len(c.Domain) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.Domain...)

	return dst, err
}

func (c *SigningData) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 64 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:32]  // c.ObjectRoot
	sszSlice1 := buf[32:64] // c.Domain

	// Field 0: ObjectRoot
	c.ObjectRoot = make([]byte, 0, 32)
	c.ObjectRoot = append(c.ObjectRoot, sszSlice0...)

	// Field 1: Domain
	c.Domain = make([]byte, 0, 32)
	c.Domain = append(c.Domain, sszSlice1...)
	return err
}

func (c *SigningData) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *SigningData) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: ObjectRoot
	if len(c.ObjectRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ObjectRoot)
	// Field 1: Domain
	if len(c.Domain) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.Domain)
	hh.Merkleize(indx)
	return nil
}

func (c *Status) SizeSSZ() int {
	size := 84

	return size
}

func (c *Status) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Status) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: ForkDigest
	if len(c.ForkDigest) != 4 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.ForkDigest...)

	// Field 1: FinalizedRoot
	if len(c.FinalizedRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.FinalizedRoot...)

	// Field 2: FinalizedEpoch
	if dst, err = c.FinalizedEpoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 3: HeadRoot
	if len(c.HeadRoot) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.HeadRoot...)

	// Field 4: HeadSlot
	if dst, err = c.HeadSlot.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *Status) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 84 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:4]   // c.ForkDigest
	sszSlice1 := buf[4:36]  // c.FinalizedRoot
	sszSlice2 := buf[36:44] // c.FinalizedEpoch
	sszSlice3 := buf[44:76] // c.HeadRoot
	sszSlice4 := buf[76:84] // c.HeadSlot

	// Field 0: ForkDigest
	c.ForkDigest = make([]byte, 0, 4)
	c.ForkDigest = append(c.ForkDigest, sszSlice0...)

	// Field 1: FinalizedRoot
	c.FinalizedRoot = make([]byte, 0, 32)
	c.FinalizedRoot = append(c.FinalizedRoot, sszSlice1...)

	// Field 2: FinalizedEpoch
	if err = c.FinalizedEpoch.UnmarshalSSZ(sszSlice2); err != nil {
		return err
	}

	// Field 3: HeadRoot
	c.HeadRoot = make([]byte, 0, 32)
	c.HeadRoot = append(c.HeadRoot, sszSlice3...)

	// Field 4: HeadSlot
	if err = c.HeadSlot.UnmarshalSSZ(sszSlice4); err != nil {
		return err
	}
	return err
}

func (c *Status) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Status) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: ForkDigest
	if len(c.ForkDigest) != 4 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.ForkDigest)
	// Field 1: FinalizedRoot
	if len(c.FinalizedRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.FinalizedRoot)
	// Field 2: FinalizedEpoch
	if hash, err := c.FinalizedEpoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 3: HeadRoot
	if len(c.HeadRoot) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.HeadRoot)
	// Field 4: HeadSlot
	if hash, err := c.HeadSlot.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	hh.Merkleize(indx)
	return nil
}

func (c *Validator) SizeSSZ() int {
	size := 121

	return size
}

func (c *Validator) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *Validator) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.PublicKey...)

	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return nil, ssz.ErrBytesLength
	}
	dst = append(dst, c.WithdrawalCredentials...)

	// Field 2: EffectiveBalance
	dst = ssz.MarshalUint64(dst, c.EffectiveBalance)

	// Field 3: Slashed
	dst = ssz.MarshalBool(dst, c.Slashed)

	// Field 4: ActivationEligibilityEpoch
	if dst, err = c.ActivationEligibilityEpoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 5: ActivationEpoch
	if dst, err = c.ActivationEpoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 6: ExitEpoch
	if dst, err = c.ExitEpoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 7: WithdrawableEpoch
	if dst, err = c.WithdrawableEpoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *Validator) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 121 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:48]    // c.PublicKey
	sszSlice1 := buf[48:80]   // c.WithdrawalCredentials
	sszSlice2 := buf[80:88]   // c.EffectiveBalance
	sszSlice3 := buf[88:89]   // c.Slashed
	sszSlice4 := buf[89:97]   // c.ActivationEligibilityEpoch
	sszSlice5 := buf[97:105]  // c.ActivationEpoch
	sszSlice6 := buf[105:113] // c.ExitEpoch
	sszSlice7 := buf[113:121] // c.WithdrawableEpoch

	// Field 0: PublicKey
	c.PublicKey = make([]byte, 0, 48)
	c.PublicKey = append(c.PublicKey, sszSlice0...)

	// Field 1: WithdrawalCredentials
	c.WithdrawalCredentials = make([]byte, 0, 32)
	c.WithdrawalCredentials = append(c.WithdrawalCredentials, sszSlice1...)

	// Field 2: EffectiveBalance
	c.EffectiveBalance = ssz.UnmarshallUint64(sszSlice2)

	// Field 3: Slashed
	c.Slashed = ssz.UnmarshalBool(sszSlice3)

	// Field 4: ActivationEligibilityEpoch
	if err = c.ActivationEligibilityEpoch.UnmarshalSSZ(sszSlice4); err != nil {
		return err
	}

	// Field 5: ActivationEpoch
	if err = c.ActivationEpoch.UnmarshalSSZ(sszSlice5); err != nil {
		return err
	}

	// Field 6: ExitEpoch
	if err = c.ExitEpoch.UnmarshalSSZ(sszSlice6); err != nil {
		return err
	}

	// Field 7: WithdrawableEpoch
	if err = c.WithdrawableEpoch.UnmarshalSSZ(sszSlice7); err != nil {
		return err
	}
	return err
}

func (c *Validator) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *Validator) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: PublicKey
	if len(c.PublicKey) != 48 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.PublicKey)
	// Field 1: WithdrawalCredentials
	if len(c.WithdrawalCredentials) != 32 {
		return ssz.ErrBytesLength
	}
	hh.PutBytes(c.WithdrawalCredentials)
	// Field 2: EffectiveBalance
	hh.PutUint64(c.EffectiveBalance)
	// Field 3: Slashed
	hh.PutBool(c.Slashed)
	// Field 4: ActivationEligibilityEpoch
	if hash, err := c.ActivationEligibilityEpoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 5: ActivationEpoch
	if hash, err := c.ActivationEpoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 6: ExitEpoch
	if hash, err := c.ExitEpoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 7: WithdrawableEpoch
	if hash, err := c.WithdrawableEpoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	hh.Merkleize(indx)
	return nil
}

func (c *VoluntaryExit) SizeSSZ() int {
	size := 16

	return size
}

func (c *VoluntaryExit) MarshalSSZ() ([]byte, error) {
	buf := make([]byte, c.SizeSSZ())
	return c.MarshalSSZTo(buf[:0])
}

func (c *VoluntaryExit) MarshalSSZTo(dst []byte) ([]byte, error) {
	var err error

	// Field 0: Epoch
	if dst, err = c.Epoch.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	// Field 1: ValidatorIndex
	if dst, err = c.ValidatorIndex.MarshalSSZTo(dst); err != nil {
		return nil, err
	}

	return dst, err
}

func (c *VoluntaryExit) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 16 {
		return ssz.ErrSize
	}

	sszSlice0 := buf[0:8]  // c.Epoch
	sszSlice1 := buf[8:16] // c.ValidatorIndex

	// Field 0: Epoch
	if err = c.Epoch.UnmarshalSSZ(sszSlice0); err != nil {
		return err
	}

	// Field 1: ValidatorIndex
	if err = c.ValidatorIndex.UnmarshalSSZ(sszSlice1); err != nil {
		return err
	}
	return err
}

func (c *VoluntaryExit) HashTreeRoot() ([32]byte, error) {
	hh := ssz.DefaultHasherPool.Get()
	if err := c.HashTreeRootWith(hh); err != nil {
		ssz.DefaultHasherPool.Put(hh)
		return [32]byte{}, err
	}
	root, err := hh.HashRoot()
	ssz.DefaultHasherPool.Put(hh)
	return root, err
}

func (c *VoluntaryExit) HashTreeRootWith(hh *ssz.Hasher) (err error) {
	indx := hh.Index()
	// Field 0: Epoch
	if hash, err := c.Epoch.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	// Field 1: ValidatorIndex
	if hash, err := c.ValidatorIndex.HashTreeRoot(); err != nil {
		return err
	} else {
		hh.AppendBytes32(hash[:])
	}
	hh.Merkleize(indx)
	return nil
}
