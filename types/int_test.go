package types

import (
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFromInt64(t *testing.T) {
	for n := 0; n < 20; n++ {
		r := rand.Int63()
		require.Equal(t, r, NewInt(r).Int64())
	}
}

func TestInt(t *testing.T) {
	// Max Int = 2^255-1 = 5.789e+76
	// Min Int = -(2^255-1) = -5.789e+76
	require.NotPanics(t, func() { NewIntWithDecimal(1, 76) })
	i1 := NewIntWithDecimal(1, 76)
	require.NotPanics(t, func() { NewIntWithDecimal(2, 76) })
	i2 := NewIntWithDecimal(2, 76)
	require.NotPanics(t, func() { NewIntWithDecimal(3, 76) })
	i3 := NewIntWithDecimal(3, 76)

	require.Panics(t, func() { NewIntWithDecimal(6, 76) })
	require.Panics(t, func() { NewIntWithDecimal(9, 80) })

	// Overflow check
	require.NotPanics(t, func() { i1.Add(i1) })
	require.NotPanics(t, func() { i2.Add(i2) })
	require.Panics(t, func() { i3.Add(i3) })

	require.NotPanics(t, func() { i1.Sub(i1.Neg()) })
	require.NotPanics(t, func() { i2.Sub(i2.Neg()) })
	require.Panics(t, func() { i3.Sub(i3.Neg()) })

	require.Panics(t, func() { i1.Mul(i1) })
	require.Panics(t, func() { i2.Mul(i2) })
	require.Panics(t, func() { i3.Mul(i3) })

	require.Panics(t, func() { i1.Neg().Mul(i1.Neg()) })
	require.Panics(t, func() { i2.Neg().Mul(i2.Neg()) })
	require.Panics(t, func() { i3.Neg().Mul(i3.Neg()) })

	// Underflow check
	i3n := i3.Neg()
	require.NotPanics(t, func() { i3n.Sub(i1) })
	require.NotPanics(t, func() { i3n.Sub(i2) })
	require.Panics(t, func() { i3n.Sub(i3) })

	require.NotPanics(t, func() { i3n.Add(i1.Neg()) })
	require.NotPanics(t, func() { i3n.Add(i2.Neg()) })
	require.Panics(t, func() { i3n.Add(i3.Neg()) })

	require.Panics(t, func() { i1.Mul(i1.Neg()) })
	require.Panics(t, func() { i2.Mul(i2.Neg()) })
	require.Panics(t, func() { i3.Mul(i3.Neg()) })

	// Bound check
	intmax := NewIntFromBigInt(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(255), nil), big.NewInt(1)))
	intmin := intmax.Neg()
	require.NotPanics(t, func() { intmax.Add(ZeroInt()) })
	require.NotPanics(t, func() { intmin.Sub(ZeroInt()) })
	require.Panics(t, func() { intmax.Add(OneInt()) })
	require.Panics(t, func() { intmin.Sub(OneInt()) })

	// Division-by-zero check
	require.Panics(t, func() { i1.Div(NewInt(0)) })
}

func TestUint(t *testing.T) {
	// Max Uint = 1.15e+77
	// Min Uint = 0
	require.NotPanics(t, func() { NewUintWithDecimal(5, 76) })
	i1 := NewUintWithDecimal(5, 76)
	require.NotPanics(t, func() { NewUintWithDecimal(10, 76) })
	i2 := NewUintWithDecimal(10, 76)
	require.NotPanics(t, func() { NewUintWithDecimal(11, 76) })
	i3 := NewUintWithDecimal(11, 76)

	require.Panics(t, func() { NewUintWithDecimal(12, 76) })
	require.Panics(t, func() { NewUintWithDecimal(1, 80) })

	// Overflow check
	require.NotPanics(t, func() { i1.Add(i1) })
	require.Panics(t, func() { i2.Add(i2) })
	require.Panics(t, func() { i3.Add(i3) })

	require.Panics(t, func() { i1.Mul(i1) })
	require.Panics(t, func() { i2.Mul(i2) })
	require.Panics(t, func() { i3.Mul(i3) })

	// Underflow check
	require.NotPanics(t, func() { i2.Sub(i1) })
	require.NotPanics(t, func() { i2.Sub(i2) })
	require.Panics(t, func() { i2.Sub(i3) })

	// Bound check
	uintmax := NewUintFromBigInt(new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1)))
	uintmin := NewUint(0)
	require.NotPanics(t, func() { uintmax.Add(ZeroUint()) })
	require.NotPanics(t, func() { uintmin.Sub(ZeroUint()) })
	require.Panics(t, func() { uintmax.Add(OneUint()) })
	require.Panics(t, func() { uintmin.Sub(OneUint()) })

	// Division-by-zero check
	require.Panics(t, func() { i1.Div(uintmin) })
}

// Tests below uses randomness
// Since we are using *big.Int as underlying value
// and (U/)Int is immutable value(see TestImmutability(U/)Int)
// it is save to use randomness in the tests

func TestIdentInt(t *testing.T) {
	for d := 0; d < 20; d++ {
		n := rand.Int63()
		i := NewInt(n)

		ifromstr, ok := NewIntFromString(strconv.FormatInt(n, 10))
		require.True(t, ok)

		cases := []int64{
			i.Int64(),
			i.BigInt().Int64(),
			ifromstr.Int64(),
			NewIntFromBigInt(big.NewInt(n)).Int64(),
			NewIntWithDecimal(n, 0).Int64(),
		}

		for _, tc := range cases {
			require.Equal(t, n, tc)
		}
	}
}

func minint(i1, i2 int64) int64 {
	if i1 < i2 {
		return i1
	}
	return i2
}

func TestArithInt(t *testing.T) {
	for d := 0; d < 20; d++ {
		n1 := int64(rand.Int31())
		i1 := NewInt(n1)
		n2 := int64(rand.Int31())
		i2 := NewInt(n2)

		cases := []struct {
			ires Int
			nres int64
		}{
			{i1.Add(i2), n1 + n2},
			{i1.Sub(i2), n1 - n2},
			{i1.Mul(i2), n1 * n2},
			{i1.Div(i2), n1 / n2},
			{i1.AddRaw(n2), n1 + n2},
			{i1.SubRaw(n2), n1 - n2},
			{i1.MulRaw(n2), n1 * n2},
			{i1.DivRaw(n2), n1 / n2},
			{MinInt(i1, i2), minint(n1, n2)},
			{i1.Neg(), -n1},
		}

		for _, tc := range cases {
			require.Equal(t, tc.nres, tc.ires.Int64())
		}
	}

}

func TestCompInt(t *testing.T) {
	for d := 0; d < 20; d++ {
		n1 := int64(rand.Int31())
		i1 := NewInt(n1)
		n2 := int64(rand.Int31())
		i2 := NewInt(n2)

		cases := []struct {
			ires bool
			nres bool
		}{
			{i1.Equal(i2), n1 == n2},
			{i1.GT(i2), n1 > n2},
			{i1.LT(i2), n1 < n2},
		}

		for _, tc := range cases {
			require.Equal(t, tc.nres, tc.ires)
		}
	}
}

func TestIdentUint(t *testing.T) {
	for d := 0; d < 20; d++ {
		n := rand.Uint64()
		i := NewUint(n)

		ifromstr, ok := NewUintFromString(strconv.FormatUint(n, 10))
		require.True(t, ok)

		cases := []uint64{
			i.Uint64(),
			i.BigInt().Uint64(),
			ifromstr.Uint64(),
			NewUintFromBigInt(new(big.Int).SetUint64(n)).Uint64(),
			NewUintWithDecimal(n, 0).Uint64(),
		}

		for _, tc := range cases {
			require.Equal(t, n, tc)
		}
	}
}

func minuint(i1, i2 uint64) uint64 {
	if i1 < i2 {
		return i1
	}
	return i2
}

func TestArithUint(t *testing.T) {
	for d := 0; d < 20; d++ {
		n1 := uint64(rand.Uint32())
		i1 := NewUint(n1)
		n2 := uint64(rand.Uint32())
		i2 := NewUint(n2)

		cases := []struct {
			ires Uint
			nres uint64
		}{
			{i1.Add(i2), n1 + n2},
			{i1.Mul(i2), n1 * n2},
			{i1.Div(i2), n1 / n2},
			{i1.AddRaw(n2), n1 + n2},
			{i1.MulRaw(n2), n1 * n2},
			{i1.DivRaw(n2), n1 / n2},
			{MinUint(i1, i2), minuint(n1, n2)},
		}

		for _, tc := range cases {
			require.Equal(t, tc.nres, tc.ires.Uint64())
		}

		if n2 > n1 {
			continue
		}

		subs := []struct {
			ires Uint
			nres uint64
		}{
			{i1.Sub(i2), n1 - n2},
			{i1.SubRaw(n2), n1 - n2},
		}

		for _, tc := range subs {
			require.Equal(t, tc.nres, tc.ires.Uint64())
		}
	}
}

func TestCompUint(t *testing.T) {
	for d := 0; d < 20; d++ {
		n1 := rand.Uint64()
		i1 := NewUint(n1)
		n2 := rand.Uint64()
		i2 := NewUint(n2)

		cases := []struct {
			ires bool
			nres bool
		}{
			{i1.Equal(i2), n1 == n2},
			{i1.GT(i2), n1 > n2},
			{i1.LT(i2), n1 < n2},
		}

		for _, tc := range cases {
			require.Equal(t, tc.nres, tc.ires)
		}
	}
}

func randint() Int {
	return NewInt(rand.Int63())
}

func TestImmutabilityInt(t *testing.T) {
	ops := []func(*Int){
		func(i *Int) { _ = i.Add(randint()) },
		func(i *Int) { _ = i.Sub(randint()) },
		func(i *Int) { _ = i.Mul(randint()) },
		func(i *Int) { _ = i.Div(randint()) },
		func(i *Int) { _ = i.AddRaw(rand.Int63()) },
		func(i *Int) { _ = i.SubRaw(rand.Int63()) },
		func(i *Int) { _ = i.MulRaw(rand.Int63()) },
		func(i *Int) { _ = i.DivRaw(rand.Int63()) },
		func(i *Int) { _ = i.Neg() },
		func(i *Int) { _ = i.IsZero() },
		func(i *Int) { _ = i.Sign() },
		func(i *Int) { _ = i.Equal(randint()) },
		func(i *Int) { _ = i.GT(randint()) },
		func(i *Int) { _ = i.LT(randint()) },
		func(i *Int) { _ = i.String() },
	}

	for i := 0; i < 20; i++ {
		n := rand.Int63()
		ni := NewInt(n)

		for _, op := range ops {
			op(&ni)

			require.Equal(t, n, ni.Int64())
			require.Equal(t, NewInt(n), ni)
		}
	}
}

func TestImmutabilityUint(t *testing.T) {
	ops := []func(*Uint){
		func(i *Uint) { _ = i.Add(randuint()) },
		func(i *Uint) { _ = i.Sub(NewUint(rand.Uint64() % i.Uint64())) },
		func(i *Uint) { _ = i.Mul(randuint()) },
		func(i *Uint) { _ = i.Div(randuint()) },
		func(i *Uint) { _ = i.AddRaw(rand.Uint64()) },
		func(i *Uint) { _ = i.SubRaw(rand.Uint64() % i.Uint64()) },
		func(i *Uint) { _ = i.MulRaw(rand.Uint64()) },
		func(i *Uint) { _ = i.DivRaw(rand.Uint64()) },
		func(i *Uint) { _ = i.IsZero() },
		func(i *Uint) { _ = i.Sign() },
		func(i *Uint) { _ = i.Equal(randuint()) },
		func(i *Uint) { _ = i.GT(randuint()) },
		func(i *Uint) { _ = i.LT(randuint()) },
		func(i *Uint) { _ = i.String() },
	}

	for i := 0; i < 20; i++ {
		n := rand.Uint64()
		ni := NewUint(n)

		for _, op := range ops {
			op(&ni)

			require.Equal(t, n, ni.Uint64())
			require.Equal(t, NewUint(n), ni)
		}
	}

}

func randuint() Uint {
	return NewUint(rand.Uint64())
}

func TestEncodingRandom(t *testing.T) {
	for i := 0; i < 20; i++ {
		n := rand.Int63()
		ni := NewInt(n)
		var ri Int

		str, err := ni.MarshalAmino()
		require.Nil(t, err)
		err = (&ri).UnmarshalAmino(str)
		require.Nil(t, err)

		require.Equal(t, ni, ri)

		bz, err := ni.MarshalJSON()
		require.Nil(t, err)
		err = (&ri).UnmarshalJSON(bz)
		require.Nil(t, err)

		require.Equal(t, ni, ri)
	}

	for i := 0; i < 20; i++ {
		n := rand.Uint64()
		ni := NewUint(n)
		var ri Uint

		str, err := ni.MarshalAmino()
		require.Nil(t, err)
		err = (&ri).UnmarshalAmino(str)
		require.Nil(t, err)

		require.Equal(t, ni, ri)

		bz, err := ni.MarshalJSON()
		require.Nil(t, err)
		err = (&ri).UnmarshalJSON(bz)
		require.Nil(t, err)

		require.Equal(t, ni, ri)
	}
}

func TestEncodingTableInt(t *testing.T) {
	var i Int

	cases := []struct {
		i   Int
		bz  []byte
		str string
	}{
		{NewInt(0), []byte("\"0\""), "0"},
		{NewInt(100), []byte("\"100\""), "100"},
		{NewInt(51842), []byte("\"51842\""), "51842"},
		{NewInt(19513368), []byte("\"19513368\""), "19513368"},
		{NewInt(999999999999), []byte("\"999999999999\""), "999999999999"},
	}

	for _, tc := range cases {
		bz, err := tc.i.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, tc.bz, bz)
		require.Nil(t, (&i).UnmarshalJSON(bz))
		require.Equal(t, tc.i, i)

		str, err := tc.i.MarshalAmino()
		require.Nil(t, err)
		require.Equal(t, tc.str, str)
		require.Nil(t, (&i).UnmarshalAmino(str))
		require.Equal(t, tc.i, i)
	}
}

func TestEncodingTableUint(t *testing.T) {
	var i Uint

	cases := []struct {
		i   Uint
		bz  []byte
		str string
	}{
		{NewUint(0), []byte("\"0\""), "0"},
		{NewUint(100), []byte("\"100\""), "100"},
		{NewUint(51842), []byte("\"51842\""), "51842"},
		{NewUint(19513368), []byte("\"19513368\""), "19513368"},
		{NewUint(999999999999), []byte("\"999999999999\""), "999999999999"},
	}

	for _, tc := range cases {
		bz, err := tc.i.MarshalJSON()
		require.Nil(t, err)
		require.Equal(t, tc.bz, bz)
		require.Nil(t, (&i).UnmarshalJSON(tc.bz))
		require.Equal(t, tc.i, i)

		str, err := tc.i.MarshalAmino()
		require.Nil(t, err)
		require.Equal(t, tc.str, str)
		require.Nil(t, (&i).UnmarshalAmino(str))
		require.Equal(t, tc.i, i)

	}
}
