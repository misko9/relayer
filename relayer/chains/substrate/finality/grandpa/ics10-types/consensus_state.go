package ics10_types

import (
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

var _ exported.ConsensusState = (*ConsensusState)(nil)

func (ConsensusState) ClientType() string {
	return ClientType
}

func (cs ConsensusState) GetTimestamp() uint64 {
	return uint64(cs.Timestamp.Nanosecond())
}

func (cs ConsensusState) ValidateBasic() error {
	if len(cs.Root) == 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "root cannot be empty")
	}

	if cs.GetTimestamp() <= 0 {
		return sdkerrors.Wrap(clienttypes.ErrInvalidConsensus, "timestamp must be a positive Unix time")
	}
	return nil
}
