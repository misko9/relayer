package grandpa

import (
	"fmt"

	ibcexported "github.com/cosmos/ibc-go/v5/modules/core/exported"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	rpctypes "github.com/ComposableFi/go-substrate-rpc-client/v4/types"
	"github.com/cosmos/relayer/v2/relayer/chains/substrate/finality/grandpa/ics10-types"
	codec "github.com/ComposableFi/go-substrate-rpc-client/v4/scale"
	rpcclienttypes "github.com/ComposableFi/go-substrate-rpc-client/v4/types"
)

type GrandpaIBCHeader struct {
	height       uint64
	SignedHeader *ParachainHeadersWithFinalityProof
	//SignedHeader *types.Header  //Contract should expect ParachainHeadersWithFinalityProof
}

func (h GrandpaIBCHeader) Height() uint64 {
	return h.height
}

func (h GrandpaIBCHeader) ConsensusState() ibcexported.ConsensusState {
	// todo: should this be the first parachain header in the list?
	//parachainHeader := h.SignedHeader.ParachainHeaders[0].ParachainHeader
	//timestamp, err := decodeExtrinsicTimestamp(parachainHeader.Extrinsic)
	//if err != nil {
	//	panic(err)
	//}

	// todo: construct the grandpa consensus state and wrap it in a Wasm consensus
	//return types.ConsensusState{Timestamp: timestamppb.New(timestamp)}
	//return nil
	panic("[ConsensusState] implement me")
}

type ParachainHeadersWithFinalityProof struct {
	FinalityProof    *FinalityProof
	ParachainHeaders []types.ParachainHeaderWithRelayHash
	current          int
}

func (p *ParachainHeadersWithFinalityProof) Reset() {
	p.FinalityProof = nil
	p.ParachainHeaders = nil
}

func (p *ParachainHeadersWithFinalityProof) String() string {
	return fmt.Sprintf("FinalityProof: %v, ParachainHeaders: %v", p.FinalityProof, p.ParachainHeaders)
}

func (p *ParachainHeadersWithFinalityProof) ProtoMessage() {
	fmt.Println("ProtoMessage")
}

func (p *ParachainHeadersWithFinalityProof) ClientType() string {
	return types.ClientType
}

func (p *ParachainHeadersWithFinalityProof) GetHeight() ibcexported.Height {
	fmt.Println("TODO: GetHeight")
	return clienttypes.Height{
		RevisionNumber: 0,
		RevisionHeight: 0,
	}
	//return p.ParachainHeaders[current].ParachainHeader.StateProof
}

func (p *ParachainHeadersWithFinalityProof) ValidateBasic() error {
	fmt.Println("TODO: ValidateBasic")
	return nil
}

/// Finality for block B is proved by providing:
/// 1) the justification for the descendant block F;
/// 2) headers sub-chain (B; F] if B != F;
type FinalityProof struct {
	/// The hash of block F for which justification is provided.
	Block rpctypes.Hash
	/// Justification of the block F.
	Justification []byte
	/// The set of headers in the range (B; F] that we believe are unknown to the caller. Ordered.
	UnknownHeaders []rpctypes.Header
}

func (fp FinalityProof) Encode(encoder codec.Encoder) error {
	var err error
	err = encoder.Write(fp.Block[:])
	err = encoder.Encode(fp.Justification)
	err = encoder.Encode(fp.UnknownHeaders)

	if err != nil {
		return err
	}
	return nil
}

func (fp *FinalityProof) Decode(decoder *codec.Decoder) error {
	var err error
	err = decoder.Decode(&fp.Block)
	if err != nil {
		return err
	}
	err = decoder.Decode(&fp.Justification)
	if err != nil {
		return err
	}
	err = decoder.Decode(&fp.UnknownHeaders)
	if err != nil {
		return err
	}
	_, err = decoder.ReadOneByte()
	if err == nil {
		return fmt.Errorf("expected end of stream")
	}
	return nil
}

type LocalHeader rpcclienttypes.Header

func (header *LocalHeader) Decode(decoder codec.Decoder) error {
	var err error
	var bn uint32
	err = decoder.Decode(&bn)
	if err != nil {
		return err
	}
	header.Number = rpcclienttypes.BlockNumber(bn)
	err = decoder.Decode(&header.ParentHash)
	if err != nil {
		return err
	}
	err = decoder.Decode(&header.StateRoot)
	if err != nil {
		return err
	}
	err = decoder.Decode(&header.ExtrinsicsRoot)
	if err != nil {
		return err
	}
	err = decoder.Decode(&header.Digest)
	if err != nil {
		return err
	}
	_, err = decoder.ReadOneByte()
	if err == nil {
		return fmt.Errorf("unexpected data after decoding header")
	}
	return nil
}
