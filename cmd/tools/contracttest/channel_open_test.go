package contracttest

import (
	"testing"

	"encoding/hex"

	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

// TestOpenChannelFail :
func TestOpenChannelFail(t *testing.T) {
	InitEnv(t, "./env.INI")
	t.Log("Test channel open to fail ...")
	a1, a2 := env.getTwoRandomAccount(t)
	testSettleTimeout := TestSettleTimeoutMin + 5
	var err error
	// test cases 1
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, 0)
	assertError(t, err)
	// test cases 2
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, common.StringToAddress("0x0"), a2.Address, testSettleTimeout)
	assertError(t, err)
	// test cases 3
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, common.StringToAddress(""), a2.Address, testSettleTimeout)
	assertError(t, err)
	// test cases 4
	_, err = env.TokenNetwork.OpenChannel(a2.Auth, FakeAccountAddress, a2.Address, testSettleTimeout)
	assertError(t, err)
	// test cases 5
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, common.StringToAddress("0x0"), testSettleTimeout)
	assertError(t, err)
	// test cases 6
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, common.StringToAddress(""), testSettleTimeout)
	assertError(t, err)
	// test cases 7
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, FakeAccountAddress, testSettleTimeout)
	assertError(t, err)
	// test cases 8
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, EmptyAccountAddress, testSettleTimeout)
	assertError(t, err)
	// test cases 9
	_, err = env.TokenNetwork.OpenChannel(a2.Auth, EmptyAccountAddress, a2.Address, testSettleTimeout)
	assertError(t, err)
	// test cases 10
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a1.Address, testSettleTimeout)
	assertError(t, err)
	// test cases 11
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMin-1)
	assertError(t, err)
	// test cases 12
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMax+1)
	assertError(t, err)
	t.Log("Test channel open to fail done")
}

// TestOpenChannelState :
func TestOpenChannelState(t *testing.T) {
	InitEnv(t, "./env.INI")
	t.Log("Test open channel state ...")
	a1, a2 := env.getTwoRandomAccount(t)
	testSettleTimeout := TestSettleTimeoutMin + 10
	// test cases 1
	_, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	time.Sleep(time.Second * 2)
	// test cases 2
	_, _, _, state, _, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	assert.Empty(t, err)
	assert.Equal(t, ChannelStateOpened, state)
	// test cases 3
	deposit, balanceHash, nonce, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assert.Empty(t, err)
	assert.Equal(t, int64(0), deposit.Int64())
	assert.Equal(t, uint64(0), nonce)
	assert.Equal(t, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))
	// test cases 4
	deposit, balanceHash, nonce, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assert.Empty(t, err)
	assert.Equal(t, int64(0), deposit.Int64())
	assert.Equal(t, uint64(0), nonce)
	assert.Equal(t, EmptyBalanceHash, hex.EncodeToString(balanceHash[:]))
	t.Log("Test open channel state done")
}

// TestOpenChannelRepeat :
func TestOpenChannelRepeat(t *testing.T) {
	InitEnv(t, "./env.INI")
	t.Log("Test open repeat channel ...")
	a1, a2 := env.getTwoRandomAccount(t)
	testSettleTimeout := TestSettleTimeoutMin + 10
	env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	var err error
	// test cases 1
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assertError(t, err)
	// test cases 2
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a2.Address, a1.Address, testSettleTimeout)
	assertError(t, err)
	t.Log("Test open repeat channel down")
}
