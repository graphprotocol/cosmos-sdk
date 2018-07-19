package graphpoc2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Keeper struct {
	cm CounterMapper
}

// NewKeeper returns a new Keeper
func NewKeeper(cm CounterMapper) Keeper {
	return Keeper{cm: cm}
}

// GetCoins returns the coins at the addr.
func (keeper Keeper) GetCounter(ctx sdk.Context, name string) int64 {
	return getCounter(ctx, keeper.cm, name)
}

// SetCoins sets the coins at the addr.
func (keeper Keeper) SetCounter(ctx sdk.Context, name string, counter int64) sdk.Error {
	return setCounter(ctx, keeper.cm, name, counter)
}

func (keeper Keeper) RegisterEvent(ctx sdk.Context, name string) error {
	return registerEvent(ctx, keeper.cm, name)
}

//______________________________________________________________________________________________

func getCounter(ctx sdk.Context, cm CounterMapper, name string) int64 {
	event := cm.GetEvent(ctx, name)
	if event == nil {
		return 0
	}
	return event.GetEventCounter()
}

//this is the ACT of storing INSIDE our store., calling upon mapper functions
func setCounter(ctx sdk.Context, cm CounterMapper, name string, counter int64) sdk.Error {
	event := cm.GetEvent(ctx, name) //getting the event object
	if event == nil {
		event = cm.NewEventObjectWithName(ctx, name)
	}
	err := event.SetEventCounter(counter) //updating the event object
	if err != nil {
		panic(err)
	}
	cm.SetEvent(ctx, event) //updating the store
	return nil
}

func registerEvent(ctx sdk.Context, cm CounterMapper, name string) sdk.Error {
	event := cm.NewEventObjectWithName(ctx, name)
	cm.SetEvent(ctx, event) //updating the store
	return nil
}
