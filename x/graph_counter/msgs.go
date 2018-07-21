package graph_counter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgCounter struct {
	Address sdk.AccAddress `json:"address"`
	Counter int64          `json:"counter"`
}

var _ sdk.Msg = MsgCounter{}

func NewMsgCounter(addr sdk.AccAddress, counter int64) MsgCounter {
	return MsgCounter{Address: addr, Counter: counter}
}

// Implements Msg.
func (msg MsgCounter) Type() string { return "counter" }

// Implements Msg.
func (msg MsgCounter) ValidateBasic() sdk.Error {
	return nil
}

// Implements Msg.
func (msg MsgCounter) GetSignBytes() []byte {
	b, err := msgCdc.MarshalJSON(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgCounter) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Address}
}
