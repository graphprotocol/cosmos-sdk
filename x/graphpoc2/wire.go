package graphpoc2

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// Register concrete types on wire codec for default AppAccount
func RegisterWire(cdc *wire.Codec) {
	cdc.RegisterInterface((*GCI)(nil), nil)
	cdc.RegisterConcrete(&GraphEvent{}, "graphpoc2/GCI", nil)
	cdc.RegisterConcrete(MsgRegisterEvent{}, "graphpoc2/MsgRegisterEvent", nil)
}

var msgCdc = wire.NewCodec()

func init() {
	RegisterWire(msgCdc)
	wire.RegisterCrypto(msgCdc)
}
