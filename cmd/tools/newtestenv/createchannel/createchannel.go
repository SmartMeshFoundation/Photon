package createchannel

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

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
func CreatAChannelAndDeposit(account1, account2 common.Address, key1, key2 *ecdsa.PrivateKey, amount int64, manager *contracts.ChannelManagerContract, token *contracts.Token, conn *ethclient.Client) (channelAddress common.Address) {
	log.Printf("createchannel between %s-%s\n", utils.APex(account1), utils.APex(account2))
	auth1 := bind.NewKeyedTransactor(key1)
	//auth1.GasLimit = uint64(params.GasLimit)
	//auth1.GasPrice = big.NewInt(params.GasPrice)
	callAuth1 := &bind.CallOpts{
		Pending: false,
		From:    account1,
		Context: context.Background(),
	}
	auth2 := bind.NewKeyedTransactor(key2)
	//auth2.GasLimit = uint64(params.GasLimit)
	//auth2.GasPrice = big.NewInt(params.GasPrice)
	tx, err := manager.NewChannel(auth1, account2, big.NewInt(40))
	if err != nil {
		log.Printf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String())
		return
	}
	log.Printf("create channel gas %s:%d,channel address=%s\n", tx.Hash().String(), tx.Gas(), tx.To().String())
	ctx := context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		log.Fatalf("failed to NewChannel when mining :%v", err)
	}
	fmt.Printf("NewChannel complete...\n")
	//step 2 deopsit
	//step 2.1 aprove
	channelAddress, err = manager.GetChannelWith(callAuth1, account2)
	if err != nil {
		log.Fatalf("failed to get channel %s", err)
		return
	}
	channel, err := contracts.NewNettingChannelContract(channelAddress, conn)
	if err != nil {
		log.Fatalf("NewNettingChannelContract err%s", err)
	}
	wg2 := sync.WaitGroup{}
	go func() {
		wg2.Add(1)
		defer wg2.Done()
		tx, err := token.Approve(auth1, channelAddress, big.NewInt(amount))
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
		tx, err = channel.Deposit(auth1, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Deposit: %v", err)
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
		tx, err := token.Approve(auth2, channelAddress, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Approve: %v", err)
		}
		ctx = context.Background()
		_, err = bind.WaitMined(ctx, conn, tx)
		if err != nil {
			log.Fatalf("failed to Approve when mining :%v", err)
		}
		fmt.Printf("Approve complete...\n")
		tx, err = channel.Deposit(auth2, big.NewInt(amount))
		if err != nil {
			log.Fatalf("Failed to Deposit: %v", err)
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
