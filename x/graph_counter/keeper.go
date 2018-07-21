package graph_counter

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

// Keeper manages transfers between accounts
type Keeper struct {
	am auth.AccountMapper
}

// NewKeeper returns a new Keeper
func NewKeeper(am auth.AccountMapper) Keeper {
	return Keeper{am: am}
}

// GetCounter returns the counter at the addr.
func (keeper Keeper) GetCounter(ctx sdk.Context, addr sdk.AccAddress) int64 {
	return getCounter(ctx, keeper.am, addr)
}

// SetCounter sets the counter at the addr.
func (keeper Keeper) SetCounter(ctx sdk.Context, addr sdk.AccAddress, counter int64) sdk.Error {
	return setCounter(ctx, keeper.am, addr, counter)
}

func getCounter(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress) int64 {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return 0
	}
	return acc.GetCounter()
}

func setCounter(ctx sdk.Context, am auth.AccountMapper, addr sdk.AccAddress, counter int64) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		acc = am.NewAccountWithAddress(ctx, addr)
	}
	err := acc.SetCounter(counter)
	if err != nil {
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}
