package contracttest

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/ethereum/go-ethereum/core/types"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

func assertSuccess(t *testing.T, count *int, err error) {
	*count++
	assert.Empty(t, err)
}

func assertFail(t *testing.T, count *int, err error) {
	*count++
	assert.NotEmpty(t, err)
}

func assertTxSuccess(t *testing.T, count *int, tx *types.Transaction, err error) {
	*count++
	assert.Empty(t, err)
	if tx != nil {
		_, err = bind.WaitMined(context.Background(), env.Client, tx)
		assert.Empty(t, err)
	}
}

func assertTxFail(t *testing.T, count *int, tx *types.Transaction, err error) {
	*count++
	assert.NotEmpty(t, err)
	if tx != nil {
		_, err = bind.WaitMined(context.Background(), env.Client, tx)
		assert.NotEmpty(t, err)
	}
}

func  assertEqual(t *testing.T, count *int, expect interface{}, actual interface{}) {
	*count++
	assert.Equal(t, expect, actual)
}