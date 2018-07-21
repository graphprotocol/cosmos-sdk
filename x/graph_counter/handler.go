package graph_counter

import (
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewHandler returns a handler for "counter" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCounter:
			return handleMsgCounter(ctx, k, msg)

		default:
			errMsg := "Unrecognized counter Msg type: " + reflect.TypeOf(msg).Name()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgCounter(ctx sdk.Context, k Keeper, msg MsgCounter) sdk.Result {

	k.SetCounter(ctx, msg.Address, msg.Counter)

	return sdk.Result{}
}
