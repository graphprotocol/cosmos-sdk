package graphpoc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
)

// GCI is the graph-counter-interface.

type GCI interface {
	GetSequence() int64
	SetSequence(int64) error

	GetNodeName() sdk.Coins
	SetNodeName(sdk.Coins) error
}

// AccountDecoder unmarshals account bytes
type AccountDecoder func(accountBytes []byte) (Account, error)

//-----------------------------------------------------------
// BaseAccount

var _ GC = (*GraphCounter)(nil) //CANT REMEMBER IF I NEED THIS

type GraphCounter struct {
	EventName         string `json:"eventname"`
	EventContractAddr string `json:"eventcontract"`
	EventData         string `json:"eventdata"`
	Counter           int64  `json:"sequence"`
}

// Prototype function for BaseAccount
func ProtoBaseGC() GC {
	return &GraphCounter{}
}

func NewNode(name string, addr string, data string) GraphCounter {
	return GraphCounter{
		EventName:         name,
		EventContractAddr: addr,
		EventData:         data,
	}
}

func (gc *GraphCounter) GetEventName() string {
	return gc.NodeName
}

func (gc *GraphCounter) SetEventName(name string) error {
	gc.NodeName = name
	return nil
}

func (gc *GraphCounter) GetEventContractAddr() string {
	return gc.NodeName
}

func (gc *GraphCounter) SetEventContractAddr(addr string) error {
	gc.NodeName = name
	return nil
}

func (gc *GraphCounter) GetEventData() string {
	return gc.NodeName
}

func (gc *GraphCounter) SetEventData(name string) error {
	gc.NodeName = name
	return nil
}

// Implements sdk.Account.
func (gc *GraphCounter) GetGraphCounter() int64 {
	return gc.Counter
}

// Implements sdk.Account.
func (gc *GraphCounter) SetGraphCounter(counter int64) error {
	gc.Counter = counter
	return nil
}

//----------------------------------------
// Wire

// Most users shouldn't use this, but this comes handy for tests.
func RegisterBaseAccount(cdc *wire.Codec) {
	cdc.RegisterInterface((*Account)(nil), nil)
	cdc.RegisterConcrete(&BaseAccount{}, "cosmos-sdk/BaseAccount", nil)
	wire.RegisterCrypto(cdc)
}
