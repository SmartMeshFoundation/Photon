package utils

import "math/big"

const (
	IntSize  = 32 << (^uint(0) >> 63)
	UintSize = 32 << (^uint(0) >> 63)
)

const (
	MaxInt  = 1<<(IntSize-1) - 1
	MinInt  = -1 << (IntSize - 1)
	MaxUint = 1<<UintSize - 1
)

var BigInt0 = big.NewInt(0)
