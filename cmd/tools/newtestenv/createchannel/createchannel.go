package createchannel

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/network/helper"

	"encoding/hex"

	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

//TransferTo ether to address
func TransferTo(conn *ethclient.Client, from *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	ctx := context.Background()
	auth := bind.NewKeyedTransactor(from)
	fromaddr := auth.From
	nonce, err := conn.NonceAt(ctx, fromaddr, nil)
	if err != nil {
		return err
	}
	msg := ethereum.CallMsg{From: fromaddr, To: &to, Value: amount, Data: nil}
	gasLimit, err := conn.EstimateGas(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to estimate gas needed: %v", err)
	}
	gasPrice, err := conn.SuggestGasPrice(ctx)
	if err != nil {
		return fmt.Errorf("failed to suggest gas price: %v", err)
	}
	rawTx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	// Create the transaction, sign it and schedule it for execution

	signedTx, err := auth.Signer(types.HomesteadSigner{}, auth.From, rawTx)
	if err != nil {
		return err
	}
	if err = conn.SendTransaction(ctx, signedTx); err != nil {
		return err
	}
	_, err = bind.WaitMined(ctx, conn, signedTx)
	if err != nil {
		return err
	}
	fmt.Printf("transfer from %s to %s amount=%s\n", fromaddr.String(), to.String(), amount)
	return nil
}

//CreatAChannelAndDeposit create a channel
func CreatAChannelAndDeposit(account1, account2 common.Address, key1, key2 *ecdsa.PrivateKey, amount *big.Int, tokenNetworkAddres common.Address, token *contracts.Token, conn *helper.SafeEthClient) {
	log.Printf("createchannel between %s-%s\n", utils.APex(account1), utils.APex(account2))
	auth1 := bind.NewKeyedTransactor(key1)
	auth2 := bind.NewKeyedTransactor(key2)
	tokenNetwork, err := contracts.NewTokenNetwork(tokenNetworkAddres, conn)
	if err != nil {
		log.Fatalf("newtoken network for %s ,err %s", tokenNetworkAddres.String(), err)
	}
	tx, err := tokenNetwork.OpenChannel(auth1, account1, account2, 100)
	if err != nil {
		log.Printf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String())
		return
	}
	ctx := context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to NewChannel when mining :%v", err)
	}
	channelID, _, _, _, _, err := tokenNetwork.GetChannelInfo(nil, account1, account2)
	log.Printf("create channel gas %s:%d,channel identifier=0x%s,tokennetworkaddress=%s\n", tx.Hash().String(), tx.Gas(), hex.EncodeToString(channelID[:]), tokenNetworkAddres.String())
	fmt.Printf("NewChannel complete...\n")
	//step 2 deopsit
	//step 2.1 aprove
	wg2 := sync.WaitGroup{}
	go func() {
		wg2.Add(1)
		defer wg2.Done()
		approve := new(big.Int)
		approve = approve.Mul(amount, big.NewInt(100)) //保证多个通道创建的时候不会因为approve冲突
		tx, err := token.Approve(auth1, tokenNetworkAddres, approve)
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		log.Printf("approve gas %s:%d\n", tx.Hash().String(), tx.Gas())
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		fmt.Printf("Approve complete...\n")
		tx, err = tokenNetwork.Deposit(auth1, account1, account2, amount)
		if err != nil {
			log.Fatalf("Failed to Deposit1: %v", err)
		}
		log.Printf("deposit gas %s:%d\n", tx.Hash().String(), tx.Gas())
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Deposit when mining :%v", err)
		}
		fmt.Printf("Deposit complete...\n")
	}()
	go func() {
		wg2.Add(1)
		defer wg2.Done()
		approve := new(big.Int)
		approve = approve.Mul(amount, big.NewInt(100)) //保证多个通道创建的时候不会因为approve冲突
		tx, err := token.Approve(auth2, tokenNetworkAddres, approve)
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		fmt.Printf("Approve complete...\n")
		tx, err = tokenNetwork.Deposit(auth2, account2, account1, amount)
		if err != nil {
			log.Fatalf("Failed to Deposit2: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Deposit when mining :%v", err)
		}
		fmt.Printf("Deposit complete...\n")
	}()
	time.Sleep(time.Millisecond * 10)
	wg2.Wait()
	return
}
