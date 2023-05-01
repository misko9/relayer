package finality

import (
	"context"

	ibcexported "github.com/cosmos/ibc-go/v7/modules/core/exported"
	"github.com/cosmos/relayer/v2/relayer/provider"
)

const (
	PrefixParas           = "Paras"
	MethodHeads           = "Heads"
	MethodParachains      = "Parachains"
)

type FinalityGadget interface {
	QueryLatestHeight(context.Context) (int64, int64, error)
	QueryHeaderAt(latestRelayChainHeight uint64) (header ibcexported.Header, err error)
	QueryHeaderOverBlocks(finalizedBlockHeight, previouslyFinalizedBlockHeight uint64) (ibcexported.Header, error)
	IBCHeader(header ibcexported.Header) provider.IBCHeader
	ClientState(header provider.IBCHeader) (ibcexported.ClientState, error)
}
