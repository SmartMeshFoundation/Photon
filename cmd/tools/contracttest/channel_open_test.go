package contracttest

import (
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"math/big"
)

// TestOpenChannelFail :
func TestOpenChannelFail(t *testing.T) {
	t.Log("Test channel open to fail ...")
	InitEnv(t, "./env.INI")
	a1, a2 := env.getTwoRandomAccount()
	testSettleTimeout := TestSettleTimeoutMin + 5
	var err error
	// test cases 1
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, 0)
	assert.NotEmpty(t, err, err.Error())
	// test cases 2
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, common.StringToAddress("0x0"), a2.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 3
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, common.StringToAddress(""), a2.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 4
	_, err = env.TokenNetwork.OpenChannel(a2.Auth, FakeAccountAddress, a2.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 5
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, common.StringToAddress("0x0"), testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 6
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, common.StringToAddress(""), testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 7
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, FakeAccountAddress, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 8
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, EmptyAccountAddress, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 9
	_, err = env.TokenNetwork.OpenChannel(a2.Auth, EmptyAccountAddress, a2.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 10
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a1.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 11
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMin-1)
	assert.NotEmpty(t, err, err.Error())
	// test cases 12
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, TestSettleTimeoutMax+1)
	assert.NotEmpty(t, err, err.Error())
	t.Log("Test done SUCCESS")
}

// TestOpenChannelState :
func TestOpenChannelState(t *testing.T) {
	t.Log("Test open channel state ...")
	InitEnv(t, "./env.INI")
	a1, a2 := env.getTwoRandomAccount()
	testSettleTimeout := TestSettleTimeoutMin + 10
	// test cases 1
	_, err := env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assert.Empty(t, err)
	// test cases 2
	_, _, _, state, _, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	assert.Empty(t, err)
	assert.Equal(t, contracts.ChannelStateOpened, state)
	// test cases 3
	deposit, balanceHash, nonce, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assert.Empty(t, err)
	assert.Equal(t, big.NewInt(0), *deposit)
	assert.Equal(t, uint(0), nonce)
	assert.Equal(t, nil, balanceHash)
	// test cases 4
	deposit, balanceHash, nonce, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assert.Empty(t, err)
	assert.Equal(t, 0, deposit)
	assert.Equal(t, 0, nonce)
	assert.Equal(t, nil, balanceHash)
	t.Log("Test done SUCCESS")
}

// TestOpenChannelRepeat :
func TestOpenChannelRepeat(t *testing.T) {
	t.Log("Test open repeat channel ...")
	InitEnv(t, "./env.INI")
	a1, a2 := env.getTwoRandomAccount()
	testSettleTimeout := TestSettleTimeoutMin + 10
	env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)

	var err error
	// test cases 1
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	// test cases 2
	_, err = env.TokenNetwork.OpenChannel(a1.Auth, a2.Address, a1.Address, testSettleTimeout)
	assert.NotEmpty(t, err, err.Error())
	t.Log("Test done SUCCESS")
}
