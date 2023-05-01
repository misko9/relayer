package grandpa

import (
	"context"
	"fmt"
	"github.com/ChainSafe/chaindb"
	rpcclient "github.com/misko9/go-substrate-rpc-client/v4"
	rpctypes "github.com/misko9/go-substrate-rpc-client/v4/types"
	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/cosmos/relayer/v2/relayer/chains/substrate/finality"
	"github.com/cosmos/relayer/v2/relayer/chains/substrate/finality/grandpa/ics10-types"
	"github.com/cosmos/relayer/v2/relayer/provider"
	//"unsafe"
)

/*
#cgo LDFLAGS: -L"${SRCDIR}/../../../../lib" -lgo_export
#include "../../../../lib/go-export.h"
*/
//SAM127import "C"

// todo
//   - implement condition for when previouslyFinalizedHeight isn't passed as an argument
//   - implement finality methods
//   - construct grandpa consensus state
//   - implement parsing methods from grandpa to wasm client types.
//   - check if validation data can be used to fetch relaychain data from the parachain
//   - make justification subscription in chain processor generic. It only subscribes to beefy justifications currently
//   - write tests for grandpa construct methods
var _ finality.FinalityGadget = &Grandpa{}

const GrandpaFinalityGadget = "grandpa"

type Grandpa struct {
	parachainClient  *rpcclient.SubstrateAPI
	relayChainClient *rpcclient.SubstrateAPI
	paraID           uint32
	relayChain       int32
	memDB            *chaindb.BadgerDB
	relayMetadata    *rpctypes.Metadata
	paraMetadata     *rpctypes.Metadata
	clientState      *types.ClientState
}

func NewGrandpa(
	parachainClient,
	relayChainClient *rpcclient.SubstrateAPI,
	paraID uint32,
	relayChain int32,
	memDB *chaindb.BadgerDB,
	clientState *types.ClientState,
) *Grandpa {
	paraMetadata, err := parachainClient.RPC.State.GetMetadataLatest()
	relayMetadata, err := relayChainClient.RPC.State.GetMetadataLatest()
	if err != nil {
		panic(err)
	}
	return &Grandpa{
		parachainClient,
		relayChainClient,
		paraID,
		relayChain,
		memDB,
		relayMetadata,
		paraMetadata,
		clientState,
	}
}

func (g *Grandpa) QueryLatestHeight(ctx context.Context) (int64, int64, error) {
	paraHeight, err := g.parachainClient.RPC.IBC.QueryLatestHeight(ctx)
	if err != nil {
		return 0, 0, err
	}
	//paraHeight = int64(block.Block.Header.Number)

	relayBlock, err := g.relayChainClient.RPC.Chain.GetBlockLatest()
	if err != nil {
		return 0, 0, err
	}
	relayChainHeight := int64(relayBlock.Block.Header.Number)
	return int64(paraHeight), relayChainHeight, nil
}

func (g *Grandpa) QueryHeaderAt(latestRelayChainHeight uint64) (header ibcexported.Header, err error) {
	fmt.Println("Querying header at height", latestRelayChainHeight)
	//hash, err := g.relayChainClient.RPC.Chain.GetBlockHash(latestRelayChainHeight)
	//header, err = g.relayChainClient.RPC.Chain.GetHeader(hash)
	//header, err = g.parachainClient.RPC.Chain.GetHeader(hash)
	header, err = g.queryFinalizedParachainHeadersWithProof(g.clientState, 1, nil)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (g *Grandpa) QueryHeaderOverBlocks(finalizedBlockHeight, previouslyFinalizedBlockHeight uint64) (ibcexported.Header, error) {
	fmt.Println("Querying header over blocks", finalizedBlockHeight, previouslyFinalizedBlockHeight)
	header, err := g.queryFinalizedParachainHeadersWithProof(g.clientState, uint32(previouslyFinalizedBlockHeight), nil)
	if err != nil {
		return nil, err
	}
	return header, nil
}

func (g *Grandpa) IBCHeader(header ibcexported.Header) provider.IBCHeader {
	fmt.Println("Converting header to IBC header")
	signedHeader := header.(*ParachainHeadersWithFinalityProof)
	return GrandpaIBCHeader{
		height:       signedHeader.GetHeight().GetRevisionHeight(),
		SignedHeader: signedHeader,
	}
}

func (g *Grandpa) ClientState(header provider.IBCHeader) (ibcexported.ClientState, error) {
	fmt.Println("Querying client state")
	grandpaHeader, ok := header.(GrandpaIBCHeader)
	if !ok {
		return nil, fmt.Errorf("got data of type %T but wanted  finality.GrandpaIBCHeader \n", header)
	}
	
	currentAuthorities, err := g.getCurrentAuthorities()
	if err != nil {
		return nil, err
	}
	
	blockHash, err := g.relayChainClient.RPC.Chain.GetBlockHash(grandpaHeader.height)
	if err != nil {
		return nil, err
	}
	
	currentSetId, err := g.getCurrentSetId(blockHash)
	if err != nil {
		return nil, err
	}
	
	latestRelayHash, err := g.relayChainClient.RPC.Chain.GetFinalizedHead()
	if err != nil {
		return nil, err
	}
	
	latestRelayheader, err := g.relayChainClient.RPC.Chain.GetHeader(latestRelayHash)
	if err != nil {
		return nil, err
	}
	
	paraHeader, err := g.getParachainHeader(latestRelayHash)
	if err != nil {
		return nil, err
	}
	
	var relayChain types.RelayChain
	switch types.RelayChain(g.relayChain) {
	case types.RelayChain_POLKADOT:
		relayChain = types.RelayChain_POLKADOT
	case types.RelayChain_KUSAMA:
		relayChain = types.RelayChain_KUSAMA
	case types.RelayChain_ROCOCO:
		relayChain = types.RelayChain_ROCOCO
	}

	return &types.ClientState{
		ParaId:             g.paraID,
		CurrentSetId:       currentSetId,
		CurrentAuthorities: currentAuthorities,
		LatestRelayHash:    latestRelayHash[:],
		LatestRelayHeight:  uint32(latestRelayheader.Number),
		LatestParaHeight:   uint32(paraHeader.Number),
		RelayChain:         relayChain,
		//XFrozenHeight: &types.ClientState_FrozenHeight{FrozenHeight: 0},
	}, nil
}

func (g *Grandpa) ConsensusState() (types.ConsensusState, error) {
	fmt.Println("Querying consensus state")
	// todo
	return types.ConsensusState{}, nil
}

func (g *Grandpa) fetchTimestampExtrinsicWithProof(blockHash rpctypes.Hash) (TimeStampExtWithProof, error) {
	var extWithProof TimeStampExtWithProof
	var block *rpctypes.SignedBlock
	var err error
	block, err = g.parachainClient.RPC.Chain.GetBlock(blockHash)
	if err != nil {
		return extWithProof, err
	}
	if block == nil {
		return extWithProof, fmt.Errorf("block not found")
	}
	extrinsics := block.Block.Extrinsics
	if len(extrinsics) == 0 {
		return extWithProof, fmt.Errorf("block has no extrinsics")
	}
	extWithProof.Ext, err = rpctypes.Encode(&extrinsics[0])
	if err != nil {
		return extWithProof, err
	}

	/*db := NewTrieDBMut()
	defer db.Free()
	for i, ext := range extrinsics {
		key, err := rpctypes.Encode(rpctypes.NewUCompactFromUInt(uint64(i)))
		if err != nil {
			return extWithProof, err
		}
		extEncoded, err := rpctypes.Encode(&ext)
		if err != nil {
			return extWithProof, err
		}
		db.Insert(key, extEncoded)
	}
	root := db.Root()
	key, err := rpctypes.Encode(rpctypes.NewUCompactFromUInt(0))
	extWithProof.Proof = db.GenerateTrieProof(root, key)*/
	return extWithProof, nil
}

