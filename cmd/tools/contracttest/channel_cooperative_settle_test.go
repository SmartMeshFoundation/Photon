package contracttest

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ethereum/go-ethereum/common"
)

// TestCooperativeSettleFail :
func TestCooperativeSettleByThird(t *testing.T) {
	InitEnv(t, "./env.INI")
	t.Log("Test cooperative settle, expect fail ...")
	a1, a2, a3 := env.getThreeRandomAccount(t)
	cooperativeSettleChannelIfExists(t, a1, a2)
	testSettleTimeout := TestSettleTimeoutMin + 10
	depositA1 := big.NewInt(20)
	depositA2 := big.NewInt(10)
	balanceA1 := big.NewInt(4)
	balanceA2 := big.NewInt(26)
	// open channel and deposit for test
	tx, err := env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, testSettleTimeout, depositA1)
	waitAndAssertSuccess(t, tx, err)
	tx, err = env.TokenNetwork.Deposit(a2.Auth, a2.Address, a1.Address, depositA2)
	waitAndAssertSuccess(t, tx, err)
	// get the sign for a1,a2,a3
	channelID, _, openBlockNumber, _, _, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	assert.Empty(t, err)
	ChainID, err := env.TokenNetwork.Chain_id(nil)
	assert.Empty(t, err)
	cs := &CoOperativeSettleForContracts{
		Particiant1:         a1.Address,
		Participant2:        a2.Address,
		Participant1Balance: balanceA1,
		Participant2Balance: balanceA2,
		ChannelIdentifier:   channelID,
		OpenBlockNumber:     openBlockNumber,
		ChainID:             ChainID,
		TokenNetworkAddress: env.TokenNetworkAddress,
	}
	signA1 := cs.sign(a1.Key)
	signA2 := cs.sign(a2.Key)
	// test case 1
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, big.NewInt(-1), a2.Address, balanceA2, signA1, signA2)
	waitAndAssertError(t, tx, err)
	// test case 2
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, big.NewInt(-1), signA1, signA2)
	waitAndAssertError(t, tx, err)
	// test case 3
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, common.HexToAddress("0x0"), balanceA1, a2.Address, balanceA2, signA1, signA2)
	waitAndAssertError(t, tx, err)
	// test case 4
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, common.HexToAddress("0x0"), balanceA2, signA1, signA2)
	waitAndAssertError(t, tx, err)
	// test case 5
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, balanceA2, common.Hex2Bytes("0x0"), signA2)
	waitAndAssertError(t, tx, err)
	// test case 6
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, balanceA2, signA1, common.Hex2Bytes("0x0"))
	waitAndAssertError(t, tx, err)
	// test case 7
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, EmptyAccountAddress, balanceA1, a2.Address, balanceA2, signA1, signA2)
	waitAndAssertError(t, tx, err)
	// test case 8
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, EmptyAccountAddress, balanceA2, signA1, signA2)
	waitAndAssertError(t, tx, err)
	// test case 9
	tx, err = env.TokenNetwork.CooperativeSettle(a3.Auth, a1.Address, balanceA1, a2.Address, balanceA2, signA1, signA2)
	waitAndAssertSuccess(t, tx, err)
}
