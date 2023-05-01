package ics10_types

import (
	"github.com/cosmos/ibc-go/v7/modules/core/exported"
)

var _ exported.ClientMessage = &Header{}

func (m Header) ClientType() string {
	return ClientType
}

func (m Header) ValidateBasic() error {
	return nil
}
