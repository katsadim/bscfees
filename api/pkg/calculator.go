package pkg

import (
	"bsc-fees/pkg/bsc"
	"bsc-fees/pkg/config"
	"math/big"
)

// how much eth does eth have
const wei = 1e+18

type Calculator interface {
	CalculateFees(txs []bsc.Transaction) big.Float
}

type calculator struct {
	cfg config.CalculatorOptions
}

func NewCalculator(cfg config.CalculatorOptions) Calculator {
	return &calculator{
		cfg: cfg,
	}
}

func (c *calculator) CalculateFees(txs []bsc.Transaction) big.Float {
	sum := big.NewFloat(0)
	for _, tx := range txs {
		fee := calculateFee(tx)
		sum.Add(sum, fee)
	}
	return *sum
}

func calculateFee(tx bsc.Transaction) *big.Float {
	x := big.NewFloat(0).Mul(new(big.Float).SetInt(&tx.GasPrice), big.NewFloat(float64(tx.GasUsed)))
	return big.NewFloat(0).Quo(x, big.NewFloat(wei))
}
