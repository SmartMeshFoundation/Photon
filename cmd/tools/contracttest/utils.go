package contracttest

import (
	"math/rand"
	"testing"

	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/binary"
	"math/big"

	"fmt"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
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

func cooperativeSettleChannelIfExists(a1 *Account, a2 *Account) {
	cs := getCooperativeSettleParams(a1, a2, big.NewInt(0), big.NewInt(0))
	if cs == nil {
		return
	}
	tx, err := env.TokenNetwork.CooperativeSettle(
		a1.Auth,
		a1.Address, cs.Participant1Balance,
		a2.Address, cs.Participant2Balance,
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
	return env.Accounts[index1], env.Accounts[index2], env.Accounts[index3]
}

func getCooperativeSettleParams(a1, a2 *Account, balanceA1, balanceA2 *big.Int) *CoOperativeSettleForContracts {
	channelID, _, openBlockNumber, state, _, err := env.TokenNetwork.GetChannelInfo(nil, a1.Address, a2.Address)
	if err != nil {
		panic(err)
	}
	if state == ChannelStateSettledOrNotExist {
		return nil
	}
	if state == ChannelStateOpened {
		balanceA1, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a1.Address, a2.Address)
		if err != nil {
			panic(err)
		}
		balanceA2, _, _, err = env.TokenNetwork.GetChannelParticipantInfo(nil, a2.Address, a1.Address)
		if err != nil {
			panic(err)
		}
	}
	ChainID, err := env.TokenNetwork.Chain_id(nil)
	if err != nil {
		panic(err)
	}
	return &CoOperativeSettleForContracts{
		Particiant1:         a1.Address,
		Participant2:        a2.Address,
		Participant1Balance: balanceA1,
		Participant2Balance: balanceA2,
		ChannelIdentifier:   channelID,
		OpenBlockNumber:     openBlockNumber,
		ChainID:             ChainID,
		TokenNetworkAddress: env.TokenNetworkAddress,
	}
}

func openChannelAndDeposit(a1, a2 *Account, depositA1, depositA2 *big.Int, settleTimeout uint64) {
	cooperativeSettleChannelIfExists(a1, a2)
	var tx *types.Transaction
	var err error
	if depositA1.Int64() > 0 {
		tx, err = env.TokenNetwork.OpenChannelWithDeposit(a1.Auth, a1.Address, a2.Address, settleTimeout, depositA1)
	} else {
		tx, err = env.TokenNetwork.OpenChannel(a1.Auth, a1.Address, a2.Address, settleTimeout)
	}
	if err != nil {
		panic(err)
	}
	_, err = bind.WaitMined(context.Background(), env.Client, tx)
	if err != nil {
		panic(err)
	}
	if depositA2.Int64() > 0 {
		tx, err = env.TokenNetwork.Deposit(a2.Auth, a2.Address, a1.Address, depositA2)
		if err != nil {
			panic(err)
		}
		_, err = bind.WaitMined(context.Background(), env.Client, tx)
		if err != nil {
			panic(err)
		}
	}
}

func endMsg(name string, count int, accounts ...*Account) string {
	msg := fmt.Sprintf("%s完成 CaseNum=%d", name, count)
	if accounts != nil && len(accounts) > 0 {
		msg = msg + " 使用账号:"
	}
	for index, account := range accounts {
		msg = msg + fmt.Sprintf("a%d=%s ", index+1, account.Address.String())
	}
	return msg
}
