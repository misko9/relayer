package ics10_types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

var _ exported.ClientState = (*ClientState)(nil)

const ClientType = "ics10-grandpa"

func (c ClientState) ClientType() string {
	return ClientType
}

func (c ClientState) GetLatestHeight() exported.Height {
	//TODO implement me
	panic("implement me return latest height")
}

func (c ClientState) Validate() error {
	//TODO implement me
	panic("implement me validate client state")
}

func (c ClientState) Status(ctx sdk.Context, store sdk.KVStore, cdc codec.BinaryCodec) exported.Status {
	//TODO implement me
	panic("implement me status")
}

func (c ClientState) ExportMetadata(store sdk.KVStore) []exported.GenesisMetadata {
	//TODO implement me
	panic("implement me export metadata")
}

func (c ClientState) ZeroCustomFields() exported.ClientState {
	//TODO implement me
	panic("implement me zero custom fields")
}

func (c ClientState) GetTimestampAtHeight(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
) (uint64, error) {
	//TODO implement me
	panic("implement me GetTimestampAtHeight")
}

func (c ClientState) Initialize(context sdk.Context, marshaler codec.BinaryCodec, store sdk.KVStore, state exported.ConsensusState) error {
	//TODO implement me
	panic("implement me initialize client state")
}

func (c ClientState) VerifyMembership(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
	value []byte,
) error {
	//TODO implement me
	panic("implement me VerifyMembership")
}

func (c ClientState) VerifyNonMembership(
	ctx sdk.Context,
	clientStore sdk.KVStore,
	cdc codec.BinaryCodec,
	height exported.Height,
	delayTimePeriod uint64,
	delayBlockPeriod uint64,
	proof []byte,
	path exported.Path,
) error {
	//TODO implement me
	panic("implement me VerifyNonMembership")
}

// VerifyClientMessage must verify a ClientMessage. A ClientMessage could be a Header, Misbehaviour, or batch update.
// It must handle each type of ClientMessage appropriately. Calls to CheckForMisbehaviour, UpdateState, and UpdateStateOnMisbehaviour
// will assume that the content of the ClientMessage has been verified and can be trusted. An error should be returned
// if the ClientMessage fails to verify.
func (c ClientState) VerifyClientMessage(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, clientMsg exported.ClientMessage) error {
	//TODO implement me
	panic("implement me VerifyClientMessage")
}

func (c ClientState) CheckForMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, msg exported.ClientMessage) bool {
	//TODO implement me
	panic("implement me CheckForMisbehaviour")
}

// UpdateStateOnMisbehaviour should perform appropriate state changes on a client state given that misbehaviour has been detected and verified
func (c ClientState) UpdateStateOnMisbehaviour(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, clientMsg exported.ClientMessage) {
	//TODO implement me
	panic("implement me UpdateStateOnMisbehaviour")

}

func (c ClientState) UpdateState(ctx sdk.Context, cdc codec.BinaryCodec, clientStore sdk.KVStore, clientMsg exported.ClientMessage) []exported.Height {
	//TODO implement me
	panic("implement me UpdateState")
}

func (c ClientState) CheckSubstituteAndUpdateState(
	ctx sdk.Context, cdc codec.BinaryCodec, subjectClientStore,
	substituteClientStore sdk.KVStore, substituteClient exported.ClientState,
) error {
	//TODO implement me
	panic("implement me CheckSubstituteAndUpdateState")
}

func (c ClientState) VerifyUpgradeAndUpdateState(
	ctx sdk.Context,
	cdc codec.BinaryCodec,
	store sdk.KVStore,
	newClient exported.ClientState,
	newConsState exported.ConsensusState,
	proofUpgradeClient,
	proofUpgradeConsState []byte,
) error {
	//TODO implement me
	panic("implement me VerifyUpgradeAndUpdateState")
}

// NewClientState creates a new ClientState instance.
func NewClientState(latestSequence uint64, consensusState *ConsensusState) *ClientState {
	//TODO implement me
	panic("implement me NewClientState")
}
