package bsc

import (
	"math/big"
	"time"
)

type Transaction struct {
	Time     time.Time
	GasPrice big.Int
	GasUsed  int
}
