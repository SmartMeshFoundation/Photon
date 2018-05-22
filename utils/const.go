package utils

import "math/big"

//BigInt0 as name
var BigInt0 = big.NewInt(0)

//MaxBigUInt256 as name
var MaxBigUInt256, _ = new(big.Int).SetString("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff", 0)

//IsValidPositiveInt256 returns true if i is a valid positive int256
func IsValidPositiveInt256(i *big.Int) bool {
	if i.Cmp(BigInt0) <= 0 {
		return false
	}
	if i.Cmp(MaxBigUInt256) > 0 {
		return false
	}
	return true
}

//IsValidUint256 retuns true if i is a valid uint256
func IsValidUint256(i *big.Int) bool {
	if i.Cmp(BigInt0) < 0 {
		return false
	}
	if i.Cmp(MaxBigUInt256) > 0 {
		return false
	}
	return true
}
