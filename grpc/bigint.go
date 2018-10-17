package grpc

import (
	"math/big"
)

// BigInt converts this BigInt to a big.Int.
func (i *BigInt) BigInt() *big.Int {
	if i == nil {
		return nil
	}
	bi := new(big.Int)
	bi.SetBytes(i.Bytes)
	if i.Negative {
		bi.Neg(bi)
	}
	return bi
}

// NewBigInt converts a big.Int to BigInt.
func NewBigInt(bi *big.Int) *BigInt {
	if bi == nil {
		return nil
	}

	return &BigInt{
		Bytes:    bi.Bytes(),
		Negative: bi.Sign() < 0,
	}
}
