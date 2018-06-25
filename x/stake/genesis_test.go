package stake

import (
	"testing"

	"github.com/stretchr/testify/require"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keep "github.com/cosmos/cosmos-sdk/x/stake/keeper"
	"github.com/cosmos/cosmos-sdk/x/stake/types"
)

func TestInitGenesis(t *testing.T) {
	ctx, _, keeper := keep.CreateTestInput(t, false, 1000)

	pool := keeper.GetPool(ctx)
	pool.UnbondedTokens = 1
	pool.UnbondedShares = sdk.OneRat()

	params := keeper.GetParams(ctx)
	var delegations []Delegation

	validators := []Validator{
		NewValidator(keep.Addrs[0], keep.PKs[0], Description{Moniker: "hoop"}),
	}

	// Setting status to bonded since for sdk.ABCIValidator
	// dosen't affect on actual execution
	validators[0].PoolShares.Status = sdk.Bonded
	validators[0].PoolShares.Amount = sdk.OneRat()
	validators[0].DelegatorShares = sdk.OneRat()

	genesisState := types.NewGenesisState(pool, params, validators, delegations)
	vals, err := InitGenesis(ctx, keeper, genesisState)
	require.NoError(t, err)

	abcivals := make([]abci.Validator, len(vals))
	for i, val := range validators {
		abcivals[i] = sdk.ABCIValidator(val)
	}

	require.Equal(t, abcivals, vals)
}
