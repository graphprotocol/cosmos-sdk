package graphpoc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

var globalCounter = []byte("globalCounter")

// This CounterMapper encodes/decodes events using the
// go-amino (binary) encoding/decoding library.
type CounterMapper struct {

	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical GCI constructor.
	proto func() GCI

	// The wire codec for binary encoding/decoding of accounts.
	cdc *wire.Codec
}

// NewCounterMapper returns a new sdk.AccountMapper that
// uses go-amino to (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewCounterMapper(cdc *wire.Codec, key sdk.StoreKey, proto func() GCI) CounterMapper {
	return CounterMapper{
		key:   key,
		proto: proto,
		cdc:   cdc,
	}
}

func (cm CounterMapper) NewEventWithName(ctx sdk.Context, name string) GCI {
	event := cm.proto()
	err := event.SetEventName(name)
	if err != nil {
		panic(err)
	}
	cm.GetNextGlobalEventNumber(ctx)
	return event
}

// Turn an address to key used to get it from the account store
func NodeNameStoreKey(name string) []byte {
	return append([]byte("name:"), []byte(name)...)
}

// Implements sdk.AccountMapper.
func (cm CounterMapper) GetEvent(ctx sdk.Context, name string) GCI {
	store := ctx.KVStore(cm.key)
	bz := store.Get(NodeNameStoreKey(name))
	if bz == nil {
		return nil
	}
	counter := cm.decodeCounter(bz)
	return counter
}

// Here we store the GraphEventcounter in our store
func (cm CounterMapper) SetEvent(ctx sdk.Context, gc GCI) {
	eventName := gc.GetEventName()
	store := ctx.KVStore(cm.key)
	bz := cm.encodeCounter(gc)
	store.Set(NodeNameStoreKey(eventName), bz)
}

func (cm CounterMapper) IterateEvents(ctx sdk.Context, process func(GCI) (stop bool)) {
	store := ctx.KVStore(cm.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("name:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		gc := cm.decodeCounter(val)
		if process(gc) {
			return
		}
		iter.Next()
	}
}

// Returns the Sequence of the account at address
func (cm CounterMapper) GetCounter(ctx sdk.Context, name string) (int64, sdk.Error) {
	event := cm.GetEvent(ctx, name)
	if event == nil {
		return 0, sdk.ErrUnknownAddress(name) //TODO: make unique error
	}
	return event.GetEventCounter(), nil
}

func (cm CounterMapper) setCounter(ctx sdk.Context, name string, newCount int64) sdk.Error {
	event := cm.GetEvent(ctx, name)
	if event == nil {
		return sdk.ErrUnknownAddress(name) //TODO: make unique error
	}
	err := event.SetEventCounter(newCount)
	if err != nil {
		panic(err)
	}
	cm.SetEvent(ctx, event)
	return nil
}

// Returns and increments the global account number counter
func (cm CounterMapper) GetNextGlobalEventNumber(ctx sdk.Context) int64 {
	var globalEventNum int64
	store := ctx.KVStore(cm.key)
	bz := store.Get(globalCounter)
	if bz == nil {
		globalEventNum = 0
	} else {
		err := cm.cdc.UnmarshalBinary(bz, &globalEventNum)
		if err != nil {
			panic(err)
		}
	}

	bz = cm.cdc.MustMarshalBinary(globalEventNum + 1)
	store.Set(globalCounter, bz)

	return globalEventNum
}

//----------------------------------------
// misc.

func (cm CounterMapper) encodeCounter(gc GCI) []byte {
	bz, err := cm.cdc.MarshalBinaryBare(gc)
	if err != nil {
		panic(err)
	}
	return bz
}

func (cm CounterMapper) decodeCounter(bz []byte) (gc GCI) {
	err := cm.cdc.UnmarshalBinaryBare(bz, &gc)
	if err != nil {
		panic(err)
	}
	return
}
