package cli

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/graphprotocol/cosmos-sdk/x/graphpoc"
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
func GetEventDecoder(cdc *wire.Codec) graphpoc.CounterDecoder {
	return func(eventBytes []byte) (graphpoc.GCI, error) {
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
func GetEventCmd(storeName string, cdc *wire.Codec, decoder graphpoc.CounterDecoder) *cobra.Command {
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
			res, err := ctx.QueryStore(graphpoc.EventNameStoreKey(eventName), storeName)
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

func RegisterEventCounter(storeName string, cdc *wire.Codec, decoder graphpoc.CounterDecoder) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "Register the event",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.NewCoreContextFromViper()

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
			msg := graphpoc.NewRegisterEvent(eventName, from)

			err = ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, []sdk.Msg{msg}, cdc)
			if err != nil {
				return err
			}
			return nil

		},
	}

	cmd.Flags().String(flagEventName, "", "Event name to update")

	return cmd
}
