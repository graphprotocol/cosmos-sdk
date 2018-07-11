package types

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var codeTypes = []CodeType{
	CodeInternal,
	CodeTxDecode,
	CodeInvalidSequence,
	CodeUnauthorized,
	CodeInsufficientFunds,
	CodeUnknownRequest,
	CodeInvalidAddress,
	CodeInvalidPubKey,
	CodeUnknownAddress,
	CodeInsufficientCoins,
	CodeInvalidCoins,
	CodeOutOfGas,
	CodeMemoTooLarge,
}

type errFn func(msg string) Error

var errFns = []errFn{
	ErrInternal,
	ErrTxDecode,
	ErrInvalidSequence,
	ErrUnauthorized,
	ErrInsufficientFunds,
	ErrUnknownRequest,
	ErrInvalidAddress,
	ErrInvalidPubKey,
	ErrUnknownAddress,
	ErrInsufficientCoins,
	ErrInvalidCoins,
	ErrOutOfGas,
	ErrMemoTooLarge,
}

func TestCodeType(t *testing.T) {
	require.True(t, ABCICodeOK.IsOK())

	for _, c := range codeTypes {
		msg := CodeToDefaultMsg(c)
		require.False(t, strings.HasPrefix(msg, "unknown code"))
	}

	msg := CodeToDefaultMsg(CodeOK)
	require.True(t, strings.HasPrefix(msg, "unknown code"))
}

func TestErrFn(t *testing.T) {
	for i, errFn := range errFns {
		err := errFn("")
		codeType := codeTypes[i]
		require.Equal(t, err.Code(), codeType)
		require.Equal(t, err.Result().Code, ToABCICode(CodespaceRoot, codeType))
	}

	require.Equal(t, ABCICodeOK, ToABCICode(CodespaceRoot, CodeOK))
}
