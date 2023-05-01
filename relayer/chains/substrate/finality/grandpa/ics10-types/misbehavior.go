package ics10_types

import (
	exported "github.com/cosmos/ibc-go/v7/modules/core/exported"
)

var (
	_ exported.ClientMessage = &Misbehaviour{}
)

func (m Misbehaviour) ClientType() string {
	return ClientType
}

func (m Misbehaviour) ValidateBasic() error {
	return nil
}
