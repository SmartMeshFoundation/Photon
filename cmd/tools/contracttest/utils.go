package contracttest

import (
	"math/rand"
	"testing"

	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"math/big"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

// CoOperativeSettleForContracts : param for CoOperativeSettle
type CoOperativeSettleForContracts struct {
	Particiant1         common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
	ChannelIdentifier   contracts.ChannelIdentifier
	OpenBlockNumber     uint64
	TokenNetworkAddress common.Address
	ChainID             *big.Int
}

func (c *CoOperativeSettleForContracts) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(c.Particiant1[:])
	buf.Write(utils.BigIntTo32Bytes(c.Participant1Balance))
	buf.Write(c.Participant2[:])
	buf.Write(utils.BigIntTo32Bytes(c.Participant2Balance))
	buf.Write(c.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, c.OpenBlockNumber)
	buf.Write(utils.BigIntTo32Bytes(c.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}

func cooperativeSettleChannelIfExists(t *testing.T, a1 *Account, a2 *Account) {
	// get channel info between a1 and a2
	channelID, _, openBlockNumber, state, _, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	if state == ChannelStateSettledOrNotExist {
		return
	}
	t.Logf("channel between %s and %s already exists, close first ...", a1.Address.String(), a2.Address.String())
	depositA1, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
	if err != nil {
		panic(err)
	}
	depositA2, _, _, err := env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
	if err != nil {
		panic(err)
	}
	ChainID, err := env.TokenNetwork.Chain_id(nil)
	if err != nil {
		panic(err)
	}
	cs := &CoOperativeSettleForContracts{
		Particiant1:         a1.Address,
		Participant2:        a2.Address,
		Participant1Balance: depositA1,
		Participant2Balance: depositA2,
		ChannelIdentifier:   channelID,
		OpenBlockNumber:     openBlockNumber,
		ChainID:             ChainID,
		TokenNetworkAddress: env.TokenNetworkAddress,
	}
	tx, err := env.TokenNetwork.CooperativeSettle(
		a1.Auth,
		a1.Address, depositA1,
		a2.Address, depositA2,
		cs.sign(a1.Key),
		cs.sign(a2.Key))
	if err != nil {
		panic(err)
	}
	_, err = bind.WaitMined(context.Background(), env.Client, tx)
	if err != nil {
		panic(err)
	}
}

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

func (env *Env) getThreeRandomAccount(t *testing.T) (*Account, *Account, *Account) {
	var index1, index2, index3 int
	n := len(env.Accounts)
	index1 = rand.Intn(n)
	index2 = rand.Intn(n)
	index3 = rand.Intn(n)
	for index1 == index2 {
		index2 = rand.Intn(n)
	}
	for index3 == index1 || index3 == index2 {
		index3 = rand.Intn(n)
	}
	t.Logf("a1=%s a2=%s a3=%s", env.Accounts[index1].Address.String(), env.Accounts[index2].Address.String(), env.Accounts[index3].Address.String())
	return env.Accounts[index1], env.Accounts[index2], env.Accounts[index3]
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

func waitAndAssertSuccess(t *testing.T, tx *types.Transaction, err error) {
	assert.Empty(t, err)
	if tx != nil {
		_, err = bind.WaitMined(context.Background(), env.Client, tx)
		assert.Empty(t, err)
	}
}

func waitAndAssertError(t *testing.T, tx *types.Transaction, err error) {
	assert.NotEmpty(t, err)
	if tx != nil {
		_, err = bind.WaitMined(context.Background(), env.Client, tx)
		assert.NotEmpty(t, err)
	}
}
