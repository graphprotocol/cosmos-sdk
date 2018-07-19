package graphpoc2

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

//------------------------------------------------------------------
// Handler for the message

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgRegisterEvent:
			return handleMsgEventRegister(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("invalid message parse in graphpoc module").Result()
		}
	}
}

func handleMsgEventRegister(ctx sdk.Context, msg MsgRegisterEvent, k Keeper) sdk.Result {

	k.RegisterEvent(ctx, msg.EventName)

	return sdk.Result{}

	// sendMsg, ok := msg.(MsgRegisterEvent)
	// if !ok {
	// 	// Create custom error message and return result
	// 	// Note: Using unreserved error codespace
	// 	return sdk.NewError(2, 1, "MsgRegisterEvent is malformed").Result()
	// }

	// // Load the store.
	// store := ctx.KVStore(key)
	// from := sendMsg.From

	// // Get sender account from the store.
	// eventBytes := store.Get(from)
	// if eventBytes == nil {
	// 	return sdk.NewError(2, 101, "Event is not stored").Result()
	// }

	// // Unmarshal the JSON account bytes.
	// var event GraphEvent
	// err := json.Unmarshal(eventBytes, &event)
	// if err != nil {
	// 	// InternalError
	// 	return sdk.ErrInternal("Error when deserializing event").Result()
	// }

	// event.SetEventCounter(event.Counter + 1)

	// // Encode sender account.
	// eventBytes, err = json.Marshal(event)
	// if err != nil {
	// 	return sdk.ErrInternal("Event encoding error").Result()
	// }

	// // Update store with updated sender account
	// store.Set(from, eventBytes)

	// // Return a success (Code 0).
	// // Add list of key-value pair descriptors ("tags").
	// return sdk.Result{
	// 	Tags: sendMsg.Tags(),
	// }
}

// Returns the sdk.Tags for the message
// func (msg MsgRegisterEvent) Tags() sdk.Tags {
// 	return sdk.NewTags("event", []byte(msg.From.String())).AppendTag("name", []byte(msg.EventName))
// }
