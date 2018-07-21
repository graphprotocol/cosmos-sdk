package cli

import (
	"github.com/pkg/errors"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	gc "github.com/cosmos/cosmos-sdk/x/graph_counter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flagCounter = "counter"
)

// SendTxCmd will create a send tx and sign it with the given key
func CountCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "count",
		Short: "Add a count to the counter",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from/to address
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			fromAcc, err := ctx.QueryStore(auth.AddressStoreKey(from), ctx.AccountStore)
			if err != nil {
				return err
			}

			// Check if account was found
			if fromAcc == nil {
				return errors.Errorf("No account with address %s was found in the state.\nAre you sure there has been a transaction involving it?", from)
			}

			counter := viper.GetInt64(flagCounter)

			msg := gc.NewMsgCounter(from, counter)

			err = ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, []sdk.Msg{msg}, cdc)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().String(flagCounter, "", "Counter value to set")

	return cmd
}
