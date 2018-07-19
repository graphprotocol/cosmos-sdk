package cli

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	graphpoc9 "github.com/graphprotocol/cosmos-sdk/x/graphpoc"
	graphpoc "github.com/graphprotocol/cosmos-sdk/x/graphpoc2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
)

const (
	flagEventName = "event"
)

// // GetAccountCmd for the auth.BaseAccount type
// func GetAccountCmdDefault(storeName string, cdc *wire.Codec) *cobra.Command {
// 	return GetAccountCmd(storeName, cdc, GetAccountDecoder(cdc))
// }

// Get account decoder for auth.DefaultAccount
func GetEventDecoder(cdc *wire.Codec) graphpoc9.CounterDecoder {
	return func(eventBytes []byte) (graphpoc9.GCI, error) {
		// acct := new(auth.BaseAccount)
		gc := new(graphpoc.GraphEvent)
		err := cdc.UnmarshalBinaryBare(eventBytes, &gc)
		if err != nil {
			panic(err)
		}
		return gc, err
	}
}

// GetAccountCmd returns a query account that will display the
// state of the account at a given address
func GetEventCmd(storeName string, cdc *wire.Codec, decoder graphpoc9.CounterDecoder) *cobra.Command {
	return &cobra.Command{
		Use:   "event [name]",
		Short: "Query event counter",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			// find the key to look up the account
			eventName := args[0]

			// key, err := sdk.AccAddressFromBech32(addr)
			// if err != nil {
			// 	return err
			// }

			// perform query
			ctx := context.NewCoreContextFromViper()
			res, err := ctx.QueryStore(graphpoc9.EventNameStoreKey(eventName), storeName)
			if err != nil {
				return err
			}

			// Check if event was found
			if res == nil {
				return errors.New("No event with name " + eventName +
					" was found in the state.\nAre you sure it exists?")
			}

			// decode the value
			event, err := decoder(res)
			if err != nil {
				return err
			}

			// print out whole account
			output, err := wire.MarshalJSONIndent(cdc, event)
			if err != nil {
				return err
			}
			fmt.Println(string(output))
			return nil
		},
	}
}

func RegisterEventCounter(storeName string, cdc *wire.Codec, decoder graphpoc9.CounterDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register the event",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper().WithDecoder(authcmd.GetAccountDecoder(cdc))

			// get the from/to address
			from, err := ctx.GetFromAddress()
			if err != nil {
				return err
			}

			// fromAcc, err := ctx.QueryStore(auth.AddressStoreKey(from), ctx.AccountStore)
			// if err != nil {
			// 	return err
			// }

			eventName := viper.GetString(flagEventName)

			// // perform query
			// res, err := ctx.QueryStore(graphpoc.EventNameStoreKey(eventName), storeName)
			// if err != nil {
			// 	return err
			// }

			// // Check if event was found
			// if res == nil {
			// 	return errors.New("No event with name " + eventName +
			// 		" was found in the state.\nAre you sure it exists?")
			// }

			// // decode the value
			// event, err := decoder(res)
			// if err != nil {
			// 	return err
			// }

			// build and sign the transaction, then broadcast to Tendermint
			msg := graphpoc9.NewRegisterEvent(eventName, from)
			fmt.Println("MSG: ", msg)
			err = ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, []sdk.Msg{msg}, cdc)
			if err != nil {
				fmt.Println("AHAHA")
				return err
			}
			return nil

		},
	}

	cmd.Flags().String(flagEventName, "", "Event name to update")

	return cmd
}
