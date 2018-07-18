package graphpoc

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// msg type for registering events
type MsgRegisterEvent struct {
	From      sdk.AccAddress
	EventName string `json:"eventname"`
	// EventContractAddr string `json:"eventcontract"`
	// EventData         string `json:"eventdata"`
}

var _ sdk.Msg = MsgRegisterEvent{}

func NewRegisterEvent(eventName string, from sdk.AccAddress) *MsgRegisterEvent {
	return &MsgRegisterEvent{
		EventName: eventName,
		From:      from,
	}
}

//nolint
func (msg MsgRegisterEvent) Type() string { return "RegisterEvent" }

func (msg MsgRegisterEvent) ValidateBasic() sdk.Error {
	if len(msg.From) == 0 {
		return sdk.ErrInvalidAddress("From address is empty")
	}
	if len(msg.EventName) == 0 {
		return sdk.ErrInvalidAddress("Event Name is empty")
	}
	return nil
}

// Implements Msg. JSON encode the message.
func (msg MsgRegisterEvent) GetSignBytes() []byte {
	bz, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(bz)
}

// Implements Msg. Return the signer.
func (msg MsgRegisterEvent) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

//------------------------------------------------------------------
// Tx

// Simple tx to wrap the Msg.
type graphTx struct {
	MsgRegisterEvent
}

// This tx only has one Msg.
func (tx graphTx) GetMsgs() []sdk.Msg {
	return []sdk.Msg{tx.MsgRegisterEvent}
}

// JSON decode MsgSend.
func txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx graphTx
	err := json.Unmarshal(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode(err.Error())
	}
	return tx, nil
}
