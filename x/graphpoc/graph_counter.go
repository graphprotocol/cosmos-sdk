package graphpoc

import (
	"github.com/cosmos/cosmos-sdk/wire"
)

// GCI is the graph-counter-interface.

type GCI interface {
	GetEventCounter() int64
	SetEventCounter(int64) error

	GetEventName() string
	SetEventName(string) error
}

// AccountDecoder unmarshals account bytes
type CounterDecoder func(counterBytes []byte) (GCI, error)

//-----------------------------------------------------------
// BaseAccount

var _ GCI = (*GraphEvent)(nil) //CANT REMEMBER IF I NEED THIS

type GraphEvent struct {
	EventName         string `json:"eventname"`
	EventContractAddr string `json:"eventcontract"`
	EventData         string `json:"eventdata"`
	Counter           int64  `json:"sequence"`
}

// Prototype function for BaseAccount
func ProtoBaseGC() GCI {
	return &GraphEvent{}
}

func NewEvent(name string, addr string, data string) GraphEvent {
	return GraphEvent{
		EventName:         name,
		EventContractAddr: addr,
		EventData:         data,
	}
}

func (gc *GraphEvent) GetEventName() string {
	return gc.EventName
}

func (gc *GraphEvent) SetEventName(name string) error {
	gc.EventName = name
	return nil
}

func (gc *GraphEvent) GetEventContractAddr() string {
	return gc.EventContractAddr
}

func (gc *GraphEvent) SetEventContractAddr(addr string) error {
	gc.EventContractAddr = addr
	return nil
}

func (gc *GraphEvent) GetEventData() string {
	return gc.EventData
}

func (gc *GraphEvent) SetEventData(name string) error {
	gc.EventData = name
	return nil
}

// Implements sdk.Account.
func (gc *GraphEvent) GetEventCounter() int64 {
	return gc.Counter
}

// Implements sdk.Account.
func (gc *GraphEvent) SetEventCounter(counter int64) error {
	gc.Counter = counter
	return nil
}

//----------------------------------------
// Wire

// Most users shouldn't use this, but this comes handy for tests.
func RegisterBaseAccount(cdc *wire.Codec) {
	cdc.RegisterInterface((*GCI)(nil), nil)
	cdc.RegisterConcrete(&GraphEvent{}, "cosmos-sdk/GraphEventCounter", nil)
	wire.RegisterCrypto(cdc)
}
