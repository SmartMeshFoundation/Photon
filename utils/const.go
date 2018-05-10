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
var MaxBigUInt256, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0)

func IsValidPositiveInt256(i *big.Int) bool {
	if i.Cmp(BigInt0) <= 0 {
		return false
	}
	if i.Cmp(MaxBigUInt256) > 0 {
		return false
	}
	return true
}

func IsValidUint256(i *big.Int) bool {
	if i.Cmp(BigInt0) < 0 {
		return false
	}
	if i.Cmp(MaxBigUInt256) > 0 {
		return false
	}
	return true
}
