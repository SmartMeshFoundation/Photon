package contracttest

import (
	"math/big"
	"testing"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
)

// TestChannelWithdrawRight : 正确调用测试
func TestChannelWithdrawRight(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	third := env.getRandomAccountExcept(t, self, partner)
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// get the data before test
	tokenBalanceSelf, depositSelf, tokenBalancePartner, depositPartner := checkStateAfterWithdraw(t, &count, self, nil, depositSelf, big.NewInt(0), partner, nil, depositPartner, big.NewInt(0))

	// self withdraw
	withdrawSelf, withdrawPartner := big.NewInt(1), big.NewInt(0)
	wpSelf := createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err := env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxSuccess(t, &count, tx, err)
	// check state
	tokenBalanceSelf, depositSelf, tokenBalancePartner, depositPartner = checkStateAfterWithdraw(t, &count, self, tokenBalanceSelf, depositSelf, withdrawSelf, partner, tokenBalancePartner, depositPartner, withdrawPartner)

	// partner withdraw
	withdrawSelf, withdrawPartner = big.NewInt(0), big.NewInt(1)
	wpSelf = createWithdrawParam(self, depositSelf, big.NewInt(0), partner, depositPartner, big.NewInt(1))
	tx, err = env.TokenNetwork.WithDraw(
		partner.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxSuccess(t, &count, tx, err)
	// check state
	tokenBalanceSelf, depositSelf, tokenBalancePartner, depositPartner = checkStateAfterWithdraw(t, &count, self, tokenBalanceSelf, depositSelf, withdrawSelf, partner, tokenBalancePartner, depositPartner, withdrawPartner)

	// third withdraw
	withdrawSelf, withdrawPartner = big.NewInt(2), big.NewInt(2)
	wpSelf = createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		third.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxSuccess(t, &count, tx, err)
	// check state
	tokenBalanceSelf, depositSelf, tokenBalancePartner, depositPartner = checkStateAfterWithdraw(t, &count, self, tokenBalanceSelf, depositSelf, withdrawSelf, partner, tokenBalancePartner, depositPartner, withdrawPartner)

	t.Log(endMsg("ChannelWithdraw 正确调用测试", count))
}

// TestChannelWithdrawException : 异常调用测试
func TestChannelWithdrawException(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// with draw when channel close
	tx, err := env.TokenNetwork.CloseChannel(self.Auth, partner.Address, big.NewInt(0), utils.EmptyHash, 0, utils.EmptyHash, nil)
	assertTxSuccess(t, nil, tx, err)
	withdrawSelf, withdrawPartner := big.NewInt(1), big.NewInt(0)
	wpSelf := createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	// with draw when channel settled
	waitToSettle(self, partner)
	tx, err = env.TokenNetwork.SettleChannel(self.Auth, self.Address, big.NewInt(0), utils.EmptyHash, partner.Address, big.NewInt(0), utils.EmptyHash)
	assertTxSuccess(t, nil, tx, err)
	withdrawSelf, withdrawPartner = big.NewInt(1), big.NewInt(0)
	wpSelf = createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("ChannelWithdraw 异常调用测试", count))

}

// TestChannelWithdrawEdge : 边界测试
func TestChannelWithdrawEdge(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	third := env.getRandomAccountExcept(t, self, partner)
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	// create param
	withdrawSelf, withdrawPartner := big.NewInt(1), big.NewInt(0)
	wpSelf := createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)

	// withdraw with data changed
	tx, err := env.TokenNetwork.WithDraw(
		self.Auth,
		EmptyAccountAddress,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		FakeAccountAddress,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		big.NewInt(0),
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		big.NewInt(0),
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	// withdraw with wrong sig
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		nil,
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		nil,
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(partner.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(self.Key),
	)
	assertTxFail(t, &count, tx, err)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(third.Key),
		wpSelf.sign(third.Key),
	)
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("ChannelWithdraw 边界测试", count))
}

// TestChannelWithdrawAttack : 恶意调用测试
func TestChannelWithdrawAttack(t *testing.T) {
	InitEnv(t, "./env.INI")
	count := 0
	// prepare
	testSettleTimeout := TestSettleTimeoutMin + 1
	self, partner := env.getTwoAccountWithoutChannelClose(t)
	depositSelf := big.NewInt(25)
	depositPartner := big.NewInt(20)
	// open channel
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)

	// withdraw when one's withdraw > deposit
	withdrawSelf, withdrawPartner := big.NewInt(26), big.NewInt(0)
	wpSelf := createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err := env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	// withdraw when one's deposit > true deposit
	withdrawSelf, withdrawPartner = big.NewInt(1), big.NewInt(0)
	wpSelf = createWithdrawParam(self, big.NewInt(depositSelf.Int64()+1), withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	wpSelf = createWithdrawParam(self, depositSelf, withdrawSelf, partner, big.NewInt(depositPartner.Int64()+1), withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	// withdraw when depositA1 + depositA2 != total deposit
	withdrawSelf, withdrawPartner = big.NewInt(1), big.NewInt(0)
	wpSelf = createWithdrawParam(self, big.NewInt(depositSelf.Int64()-1), withdrawSelf, partner, depositPartner, withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)
	wpSelf = createWithdrawParam(self, depositSelf, withdrawSelf, partner, big.NewInt(depositPartner.Int64()-1), withdrawPartner)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	// withdraw on reopen channel with old param
	withdrawSelf, withdrawPartner = big.NewInt(1), big.NewInt(0)
	wpSelf = createWithdrawParam(self, depositSelf, withdrawSelf, partner, depositPartner, withdrawPartner)
	cooperativeSettleChannelIfExists(self, partner)
	openChannelAndDeposit(self, partner, depositSelf, depositPartner, testSettleTimeout)
	tx, err = env.TokenNetwork.WithDraw(
		self.Auth,
		wpSelf.Participant1,
		wpSelf.Participant2,
		wpSelf.Participant1Deposit,
		wpSelf.Participant1Withdraw,
		wpSelf.sign(self.Key),
		wpSelf.sign(partner.Key),
	)
	assertTxFail(t, &count, tx, err)

	t.Log(endMsg("ChannelWithdraw 恶意调用测试", count))
}

func checkStateAfterWithdraw(
	t *testing.T,
	count *int,
	a1 *Account,
	tokenBalanceA1 *big.Int,
	depositA1 *big.Int,
	withdrawA1 *big.Int,
	a2 *Account,
	tokenBalanceA2 *big.Int,
	depositA2 *big.Int,
	withdrawA2 *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int) {
	// check a1's token
	tokenBalanceA1Now, err := env.Token.BalanceOf(nil, a1.Address)
	assertSuccess(t, nil, err)
	if tokenBalanceA1 != nil {
		assertEqual(t, count, tokenBalanceA1.Add(tokenBalanceA1, withdrawA1), tokenBalanceA1Now)
	}

	// check a2's token
	tokenBalanceA2Now, err := env.Token.BalanceOf(nil, a2.Address)
	assertSuccess(t, nil, err)
	if tokenBalanceA2 != nil {
		assertEqual(t, count, tokenBalanceA2.Add(tokenBalanceA2, withdrawA2), tokenBalanceA2Now)
	}

	// check a1's deposit
	depositA1Now, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, count, depositA1.Sub(depositA1, withdrawA1), depositA1Now)

	// get a2's deposit
	depositA2Now, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	assertSuccess(t, nil, err)
	assertEqual(t, count, depositA2.Sub(depositA2, withdrawA2), depositA2Now)

	// return new
	return tokenBalanceA1Now, depositA1Now, tokenBalanceA2Now, depositA2Now
}
