package pkg

import (
	"bsc-fees/pkg/bsc"
	"bsc-fees/pkg/config"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tx1 = bsc.Transaction{GasPrice: *big.NewInt(10000000000), GasUsed: 121823}
var tx2 = bsc.Transaction{GasPrice: *big.NewInt(11000000000), GasUsed: 21000}

func Test_calculator_CalculateFees(t *testing.T) {

	tests := []struct {
		name string
		txs  []bsc.Transaction
		want big.Float
	}{
		{"No transactions", []bsc.Transaction{}, *big.NewFloat(0)},
		{"One transaction", []bsc.Transaction{tx1}, *big.NewFloat(0.00121823)},
		{"Two transactions", []bsc.Transaction{tx1, tx2}, *big.NewFloat(0.00144923)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCalculator(config.CalculatorOptions{})

			got := c.CalculateFees(tt.txs)

			assert.Equal(t, got.String(), tt.want.String())
		})
	}
}
