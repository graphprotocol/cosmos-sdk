package types

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGasMeter(t *testing.T) {
	for i := 0; i < 20; i++ {
		limit := Gas(rand.Int63n(10000))
		meter := NewGasMeter(limit)
		used := int64(0)

		for {
			gas := Gas(rand.Int63n(limit))
			used += gas
			if used > limit {
				require.Panics(t, func() { meter.ConsumeGas(gas, "") })
				break
			}
			require.NotPanics(t, func() { meter.ConsumeGas(gas, "") })
			require.Equal(t, used, meter.GasConsumed())
		}
	}
}
