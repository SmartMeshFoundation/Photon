package contracttest

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (env *Env) getTwoRandomAccount(t *testing.T) (*Account, *Account) {
	var index1, index2 int
	n := len(env.Accounts)
	index1 = rand.Intn(n)
	index2 = rand.Intn(n)
	for index1 == index2 {
		index2 = rand.Intn(n)
	}
	t.Logf("a1=%s a2=%s", env.Accounts[index1].Address.String(), env.Accounts[index2].Address.String())
	return env.Accounts[index1], env.Accounts[index2]
}

func assertError(t *testing.T, err error) {
	if err != nil {
		assert.NotEmpty(t, err, err.Error())
	}
}

func assertErrorWithMsg(t *testing.T, err error, msg string) {
	if err != nil {
		assert.NotEmpty(t, err, msg)
	}
}
