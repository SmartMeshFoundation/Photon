package contracttest

import (
	"context"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func assertSuccess(t *testing.T, count *int, err error) {
	if count != nil {
		*count++
	}
	assert.Empty(t, err)
}

func assertFail(t *testing.T, count *int, err error) {
	if count != nil {
		*count++
	}
	assert.NotEmpty(t, err)
}

func assertTxSuccess(t *testing.T, count *int, tx *types.Transaction, err error) {
	if count != nil {
		*count++
	}
	assert.Empty(t, err)
	if tx != nil {
		r, err := bind.WaitMined(context.Background(), env.Client, tx)
		assert.Empty(t, err)
		assert.EqualValues(t, 1, r.Status)
	}
}

func assertTxFail(t *testing.T, count *int, tx *types.Transaction, err error) {
	if count != nil {
		*count++
	}
	assert.NotEmpty(t, err)
	if tx != nil {
		r, err := bind.WaitMined(context.Background(), env.Client, tx)
		assert.Empty(t, err)
		assert.EqualValues(t, 0, r.Status)
	}
}

func assertEqual(t *testing.T, count *int, expect interface{}, actual interface{}) {
	if count != nil {
		*count++
	}
	assert.EqualValues(t, expect, actual)
}
