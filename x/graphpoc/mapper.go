package graphpoc

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	wire "github.com/cosmos/cosmos-sdk/wire"
)

var globalCounter = []byte("globalCounter")

// This CounterMapper encodes/decodes accounts using the
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

func (cm CounterMapper) NewCounterWithName(ctx sdk.Context, name string) GCI {
	counter := cm.proto()
	err := acc.SetNodeName(name)
	if err != nil {
		panic(err)
	}
	err = acc.SetGlobalCounter(am.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return counter
}

// Turn an address to key used to get it from the account store
func NodeNameStoreKey(name string) []byte {
	return append([]byte("name:"), name)
}

// Implements sdk.AccountMapper.
func (cm CounterMapper) GetAccount(ctx sdk.Context, name string) GCI {
	store := ctx.KVStore(cm.key)
	bz := store.Get(NodeNameStoreKey(name))
	if bz == nil {
		return nil
	}
	counter := cm.decodeCounter(bz)
	return counter
}

// Implements sdk.AccountMapper.
func (am AccountMapper) SetAccount(ctx sdk.Context, acc Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(am.key)
	bz := am.encodeAccount(acc)
	store.Set(AddressStoreKey(addr), bz)
}

// Implements sdk.AccountMapper.
func (am AccountMapper) IterateAccounts(ctx sdk.Context, process func(Account) (stop bool)) {
	store := ctx.KVStore(am.key)
	iter := sdk.KVStorePrefixIterator(store, []byte("account:"))
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := am.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// Returns the Sequence of the account at address
func (am AccountMapper) GetSequence(ctx sdk.Context, addr sdk.AccAddress) (int64, sdk.Error) {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdk.ErrUnknownAddress(addr.String())
	}
	return acc.GetSequence(), nil
}

func (am AccountMapper) setSequence(ctx sdk.Context, addr sdk.AccAddress, newSequence int64) sdk.Error {
	acc := am.GetAccount(ctx, addr)
	if acc == nil {
		return sdk.ErrUnknownAddress(addr.String())
	}
	err := acc.SetSequence(newSequence)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	am.SetAccount(ctx, acc)
	return nil
}

// Returns and increments the global account number counter
func (am AccountMapper) GetNextAccountNumber(ctx sdk.Context) int64 {
	var accNumber int64
	store := ctx.KVStore(am.key)
	bz := store.Get(globalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		err := am.cdc.UnmarshalBinary(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = am.cdc.MustMarshalBinary(accNumber + 1)
	store.Set(globalAccountNumberKey, bz)

	return accNumber
}

//----------------------------------------
// misc.

func (am AccountMapper) encodeAccount(acc Account) []byte {
	bz, err := am.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	return bz
}

func (am AccountMapper) decodeCounter(bz []byte) (acc Account) {
	err := am.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
