package contractstest

import (
	"context"
	"os"
	"testing"

	"math/big"

	"bytes"
	"crypto/ecdsa"

	"encoding/hex"

	"fmt"

	"time"

	"encoding/binary"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contractsv300/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
)

var client *ethclient.Client
var auth *bind.TransactOpts
var tokenNetworkAddress common.Address
var tokenNetwork *contracts.TokenNetwork
var ChainID *big.Int
var totalAmount int64 = 50
var tokenAddress common.Address
var token *contracts.Token
var totalLockNumber = 3
var settleTimeout = 30
var TestRPCEndpoint = os.Getenv("ETHRPCENDPOINT")

//应该作为参数出现,但是为了简单,做成全局变量.
var openBlockNumber uint64

//给 punish 操作留出的窗口时间.
var punishBlockNumber uint64

//TestPrivKey for test only
var TestPrivKey *ecdsa.PrivateKey

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
	keybin, err := hex.DecodeString(os.Getenv("KEY1"))
	if err != nil {
		log.Crit("err %s", err)
	}
	TestPrivKey, err = crypto.ToECDSA(keybin)
	if err != nil {
		log.Crit("err %s", err)
	}
	setup()
}

func setup() {
	var err error
	client, err = ethclient.Dial(TestRPCEndpoint)
	auth = bind.NewKeyedTransactor(TestPrivKey)
	if err != nil {
		panic(err)
	}
	tokenNetworkAddress = common.HexToAddress("0x13ccf341dccef573a75d46b1820cb3096151fee5")
	tokenNetwork, err = contracts.NewTokenNetwork(tokenNetworkAddress, client)
	if err != nil {
		panic(err)
	}
	ChainID, err = tokenNetwork.Chain_id(nil)
	if err != nil {
		panic(err)
	}
	tokenAddress, err = tokenNetwork.Token(nil)
	if err != nil {
		panic(err)
	}
	token, err = contracts.NewToken(tokenAddress, client)
	if err != nil {
		panic(err)
	}
	log.Info(fmt.Sprintf("tokenAddr=%s,tokenNetwork=%s", tokenAddress.String(), tokenNetworkAddress.String()))
	punishBlockNumber, err = tokenNetwork.Punish_block_number(nil)
	if err != nil {
		panic(err)
	}
}

//TransferTo ether to address
func TransferTo(conn *ethclient.Client, from *ecdsa.PrivateKey, to common.Address, amount *big.Int) error {
	ctx := context.Background()
	auth2 := bind.NewKeyedTransactor(from)
	fromaddr := auth2.From
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
	log.Info(fmt.Sprintf("gasLimit=%d,gasPrice=%s", gasLimit, gasPrice.String()))
	rawTx := types.NewTransaction(nonce, to, amount, gasLimit, gasPrice, nil)
	// Create the transaction, sign it and schedule it for execution

	signedTx, err := auth2.Signer(types.HomesteadSigner{}, auth2.From, rawTx)
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

/*
creatAChannelAndDeposit create a channel
1,2之间创建通道,总是都由1作为 tx 发起人
*/
func creatAChannelAndDeposit(account1, account2 common.Address, key1 *ecdsa.PrivateKey, amount int64, conn *ethclient.Client) error {
	log.Trace(fmt.Sprintf("createchannel between %s-%s,tokenNetwork=%s\n", account1.String(), account2.String(), tokenNetworkAddress.String()))
	auth1 := bind.NewKeyedTransactor(key1)
	tx, err := tokenNetwork.OpenChannel(auth1, account1, account2, uint64(settleTimeout))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String()))
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to NewChannel when mining :%v", err)
	}
	log.Info(fmt.Sprintf("OpenChannel gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	channelID, _, _, state, settleBlock, err := tokenNetwork.GetChannelInfo(nil, account1, account2)
	if err != nil {
		return fmt.Errorf("GetChannelInfo %s-%s err %s", utils.APex2(account1), utils.APex2(account2), err)
	}
	log.Trace(fmt.Sprintf("create channel gas %s:%d,channel identifier=0x%s,tokennetworkaddress=%s\n", tx.Hash().String(), tx.Gas(), hex.EncodeToString(channelID[:]), tokenNetworkAddress.String()))
	log.Info("NewChannel complete...\n")
	if settleBlock != uint64(settleTimeout) {
		return fmt.Errorf("settleBlock err expect=%d,got=%d", settleTimeout, settleBlock)
	}
	if state != 1 {
		return fmt.Errorf("")
	}
	tx, err = tokenNetwork.Deposit(auth1, account1, account2, big.NewInt(amount))
	if err != nil {
		return fmt.Errorf("Failed to Deposit1: %v", err)

	}
	log.Trace(fmt.Sprintf("deposit gas %s:%d\n", tx.Hash().String(), tx.Gas()))
	ctx = context.Background()
	_, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to Deposit when mining :%v", err)
	}
	log.Info("Deposit1 complete...\n")

	tx, err = tokenNetwork.Deposit(auth1, account2, account1, big.NewInt(amount))
	if err != nil {
		return fmt.Errorf("Failed to Deposit2: %v", err)

	}
	ctx = context.Background()
	r, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to Deposit when mining :%v", err)
	}
	log.Info(fmt.Sprintf("Deposit complete...,gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	return nil
}

/*
open and deposit 一起,希望能够节省一点 gas
*/
func creatAChannelAndDeposit2(account1, account2 common.Address, key1 *ecdsa.PrivateKey, amount int64, conn *ethclient.Client) error {
	log.Trace(fmt.Sprintf("createchannel between %s-%s,tokenNetwork=%s\n", account1.String(), account2.String(), tokenNetworkAddress.String()))
	auth1 := bind.NewKeyedTransactor(key1)
	tx, err := tokenNetwork.OpenChannelWithDeposit(auth1, account1, account2, uint64(settleTimeout), big.NewInt(amount))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String()))
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to NewChannel when mining :%v", err)
	}
	log.Info(fmt.Sprintf("OpenChannelWithDeposit gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	channelID, _, _, state, settleBlock, err := tokenNetwork.GetChannelInfo(nil, account1, account2)
	if err != nil {
		return fmt.Errorf("GetChannelInfo %s-%s err %s", utils.APex2(account1), utils.APex2(account2), err)
	}
	log.Trace(fmt.Sprintf("create channel gas %s:%d,channel identifier=0x%s,tokennetworkaddress=%s\n", tx.Hash().String(), tx.Gas(), hex.EncodeToString(channelID[:]), tokenNetworkAddress.String()))
	log.Info("NewChannel complete...\n")
	if settleBlock != uint64(settleTimeout) {
		return fmt.Errorf("settleBlock err expect=%d,got=%d", settleTimeout, settleBlock)
	}
	if state != 1 {
		return fmt.Errorf("")
	}
	tx, err = tokenNetwork.Deposit(auth1, account2, account1, big.NewInt(amount))
	if err != nil {
		return fmt.Errorf("Failed to Deposit2: %v", err)

	}
	ctx = context.Background()
	r, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to Deposit when mining :%v", err)
	}
	log.Info(fmt.Sprintf("Deposit2 complete...,gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	return nil
}
func testApprove(t *testing.T) {
	tx, err := token.Approve(auth, tokenNetworkAddress, big.NewInt(50000000))
	if err != nil {
		t.Error(err)
		return
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Error("receipt status error")
		return
	}
	t.Logf("%s approve token %s for %s", auth.From.String(), tokenAddress.String(), tokenNetworkAddress.String())
}

//跑一次就够了,这样后续创建通道就不用每次 appro
func TestApprove(t *testing.T) {
	testApprove(t)
}

func getTestOpenChannel(t *testing.T) (channelID contracts.ChannelIdentifier, partnerAddr common.Address, partnerKey *ecdsa.PrivateKey, err error) {
	var settleBlockNumber uint64
	var state uint8
	testApprove(t)
	//keybin, err := hex.DecodeString(os.Getenv("KEY2"))
	//if err != nil {
	//	log.Crit("err %s", err)
	//}
	//partnerKey, err = crypto.ToECDSA(keybin)
	//if err != nil {
	//	log.Crit("err %s", err)
	//}
	//partnerAddr = crypto.PubkeyToAddress(partnerKey.PublicKey)
	partnerKey, partnerAddr = utils.MakePrivateKeyAddress()
	err = creatAChannelAndDeposit(auth.From, partnerAddr, TestPrivKey, totalAmount/2, client)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, settleBlockNumber, openBlockNumber, state, _, err = tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateOpened {
		err = fmt.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateOpened, state)
		return
	}
	t.Logf("channelID=%s,settleblockNumber=%d,state=%d,err=%s", hex.EncodeToString(channelID[:]), settleBlockNumber, state, err)
	return
}
func TestOpenChannel(t *testing.T) {
	_, _, _, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
	}
}
func TestOpenChannel2(t *testing.T) {
	var settleBlockNumber uint64
	var state uint8
	testApprove(t)
	_, partnerAddr := utils.MakePrivateKeyAddress()
	err := creatAChannelAndDeposit2(auth.From, partnerAddr, TestPrivKey, totalAmount/2, client)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, settleBlockNumber, _, state, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateOpened {
		err = fmt.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateOpened, state)
		return
	}
	t.Logf("channelID=%s,settleblockNumber=%d,state=%d,err=%s", hex.EncodeToString(channelID[:]), settleBlockNumber, state, err)
	return
}
func TestCloseChannel1(t *testing.T) {
	_, partnerAddr, _, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}

	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, utils.BigInt0, utils.EmptyHash, 0, utils.EmptyHash, nil)
	if err != nil {
		t.Error(err)
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
	}
	log.Info(fmt.Sprintf("CloseChannel no evidence gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
}

//BalanceData of contract
type BalanceData struct {
	TransferAmount *big.Int
	LocksRoot      common.Hash
}

func (b *BalanceData) Hash() common.Hash {
	buf := new(bytes.Buffer)
	buf.Write(b.LocksRoot[:])
	buf.Write(utils.BigIntTo32Bytes(b.TransferAmount))
	return utils.Sha3(buf.Bytes())
}

//BalanceProofForContract for contract
type BalanceProofForContract struct {
	AdditionalHash      common.Hash
	ChannelIdentifier   contracts.ChannelIdentifier
	TokenNetworkAddress common.Address
	ChainID             *big.Int
	Signature           []byte
	OpenBlockNumber     uint64
	Nonce               uint64
	BalanceData
}

func createPartnerBalanceProof(key *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier) *BalanceProofForContract {
	bd := &BalanceData{
		TransferAmount: big.NewInt(10),
		LocksRoot:      utils.EmptyHash,
	}
	bp := &BalanceProofForContract{
		BalanceData:         *bd,
		OpenBlockNumber:     openBlockNumber,
		AdditionalHash:      utils.Sha3([]byte("123")),
		ChannelIdentifier:   channelID,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
		Nonce:               3,
	}
	bp.sign(key)
	return bp
}

func createPartnerBalanceProofWithLocks(key *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier, lockNumber int, expiredBlock int64) (bp *BalanceProofForContract, locks []*mtree.Lock, secrets []common.Hash) {
	for i := 0; i < lockNumber; i++ {
		secret := utils.Sha3([]byte(utils.RandomString(10)))
		secrets = append(secrets, secret)
		l := &mtree.Lock{
			Expiration:     expiredBlock,
			Amount:         big.NewInt(1),
			LockSecretHash: utils.Sha3(secret[:]),
		}
		locks = append(locks, l)
	}
	m := mtree.NewMerkleTree(locks)
	bd := &BalanceData{
		TransferAmount: big.NewInt(10),
		LocksRoot:      m.MerkleRoot(),
	}
	bp = &BalanceProofForContract{
		BalanceData:         *bd,
		OpenBlockNumber:     openBlockNumber,
		AdditionalHash:      utils.Sha3([]byte("123")),
		ChannelIdentifier:   channelID,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
		Nonce:               3,
	}
	bp.sign(key)
	return
}
func (b *BalanceProofForContract) sign(key *ecdsa.PrivateKey) {
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(b.TransferAmount))
	buf.Write(b.LocksRoot[:])
	binary.Write(buf, binary.BigEndian, b.Nonce)
	buf.Write(b.AdditionalHash[:])
	buf.Write(b.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, b.OpenBlockNumber)
	//buf.Write(b.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(b.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	b.Signature = sig
}
func TestCloseChannel2(t *testing.T) {
	_, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, _, _, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	if err != nil {
		t.Error(err)
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
	}
	log.Info(fmt.Sprintf("CloseChannel with evidence gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
}

type BalanceProofUpdateForContracts struct {
	BalanceProofForContract
	NonClosingSignature []byte
}

func NewBalanceProofUpdateForContracts(closingKey, nonClosingKey *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier) *BalanceProofUpdateForContracts {
	bp1 := createPartnerBalanceProof(closingKey, channelID)
	bp2 := &BalanceProofUpdateForContracts{
		BalanceProofForContract: *bp1,
	}
	bp2.sign(nonClosingKey)
	return bp2
}
func NewBalanceProofUpdateForContractsWithLocks(closingKey, nonClosingKey *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier, lockNumber int, expiredBlock int64) (bp *BalanceProofUpdateForContracts, locks []*mtree.Lock, secrets []common.Hash) {
	bp1, locks, secrets := createPartnerBalanceProofWithLocks(closingKey, channelID, lockNumber, expiredBlock)
	bp = &BalanceProofUpdateForContracts{
		BalanceProofForContract: *bp1,
	}
	bp.sign(nonClosingKey)
	return
}
func (b *BalanceProofUpdateForContracts) sign(key *ecdsa.PrivateKey) {
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(b.TransferAmount))
	buf.Write(b.LocksRoot[:])
	binary.Write(buf, binary.BigEndian, b.Nonce)
	buf.Write(b.AdditionalHash[:])
	buf.Write(b.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, b.OpenBlockNumber)
	//buf.Write(b.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(b.ChainID))
	buf.Write(b.Signature)
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	b.NonClosingSignature = sig
}
func TestCloseChannelAndUpdateBalanceProofDelegateAndSettle(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	if err != nil {
		t.Error(err)
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	/*
	   updatebalanceproof delegate 只能在结算时间的后半段
	*/
	_, settleBlockNumber, _, state, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateClosed {
		t.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateClosed, state)
		return
	}
	for {
		var h *types.Header
		h, err = client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		if h.Number.Int64() >= int64(settleBlockNumber-uint64(settleTimeout)/2) {
			//could updatebalance proof
			break
		}
		time.Sleep(time.Second)
	}
	bp2 := NewBalanceProofUpdateForContracts(TestPrivKey, partnerKey, channelID)
	log.Info(fmt.Sprintf("UpdateBalanceProofDelegate closing_participant=%s,\nnon_closing_participant=%s,\ntransferred_amount=%s,\nlocksroot=%s,\nnonce=%d,\nold_transferred_amount=%s,\nold_locksroot=%s,\nold_nonce=%d,\nadditional_hash=%s\n,closing_signature=%s\nnon_closing_signature=%s\n",
		auth.From.String(),
		partnerAddr.String(),
		bp2.TransferAmount.String(),
		bp2.LocksRoot.String(),
		bp2.Nonce,
		utils.BigInt0, utils.EmptyHash.String(), utils.BigInt0,
		bp2.AdditionalHash.String(),
		hex.EncodeToString(bp2.Signature),
		hex.EncodeToString(bp2.NonClosingSignature),
	))

	tx, err = tokenNetwork.UpdateBalanceProofDelegate(auth, auth.From, partnerAddr, bp2.TransferAmount, bp2.LocksRoot, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("UpdateBalanceProofDelegate gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))

	for {
		var h *types.Header
		h, err = client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		if h.Number.Int64() > int64(settleBlockNumber+punishBlockNumber) {
			//could settle
			break
		}
		time.Sleep(time.Second)
	}
	log.Info(fmt.Sprintf("SettleChannel arg p1=%s,p2=%s,p1.transferAmount=%s,"+
		"p2.transferAmount=%s,p1.locksroot=%s,p2.locksroot=%s,p1.nonce=%d,p2.nonce=%d,"+
		"bp1.balance_hash=%s,bp2.balance_hash=%s",
		partnerAddr.String(), auth.From.String(),
		bp.TransferAmount, bp2.TransferAmount,
		bp.LocksRoot.String(), bp2.LocksRoot.String(),
		bp.Nonce, bp2.Nonce,
		bp.BalanceData.Hash().String(), bp2.BalanceData.Hash().String(),
	))
	tx, err = tokenNetwork.SettleChannel(
		auth,
		partnerAddr,
		bp.TransferAmount,
		bp.LocksRoot,
		auth.From,
		bp2.TransferAmount,
		bp2.LocksRoot,
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("SettleChannel gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
}

func TestCloseChannelAndUpdateBalanceProofAndSettle(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	/*
		partneraddr 需要有 ether 作为 gas
	*/
	err = TransferTo(client, TestPrivKey, partnerAddr, big.NewInt(params.Ether))
	if err != nil {
		t.Error(err)
	}
	partnerAuth := bind.NewKeyedTransactor(partnerKey)
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	if err != nil {
		t.Error(err)
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Trace(fmt.Sprintf("bp=\n%s", utils.StringInterface(bp, 3)))
	bp2 := NewBalanceProofUpdateForContracts(TestPrivKey, partnerKey, channelID)
	log.Info(fmt.Sprintf("UpdateBalanceProof closing_participant=%s,\ntransferred_amount=%s,\nlocksroot=%s,\nnonce=%d,\nold_transferred_amount=%s,\nold_locksroot=%s,\nold_nonce=%d,\nadditional_hash=%s\n,closing_signature=%s\n,balance_hash=%s",
		auth.From.String(),
		bp2.TransferAmount.String(),
		bp2.LocksRoot.String(),
		bp2.Nonce,
		utils.BigInt0, utils.EmptyHash.String(), utils.BigInt0,
		bp2.AdditionalHash.String(),
		hex.EncodeToString(bp2.Signature),
		bp2.BalanceData.Hash().String(),
	))
	tx, err = tokenNetwork.UpdateBalanceProof(partnerAuth, auth.From, bp2.TransferAmount, bp2.LocksRoot, bp2.Nonce, bp2.AdditionalHash, bp2.Signature)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("UpdateBalanceProof gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	_, blokNumber, _, state, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateClosed {
		t.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateClosed, state)
		return
	}
	for {
		var h *types.Header
		h, err = client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		if h.Number.Int64() > int64(blokNumber+punishBlockNumber) {
			//could settle
			break
		}
		time.Sleep(time.Second)
	}
	log.Trace(fmt.Sprintf("SettleChannel arg,p1=%s,p1.amount=%s,p1.lock=%s,p1.nonce=%d,p2=%s,p2.amount=%s,p2.lock=%s,p2.nonce=%d",
		partnerAddr.String(), bp.TransferAmount, bp.LocksRoot.String(), bp.Nonce, auth.From.String(), bp2.TransferAmount, bp2.LocksRoot.String(), bp2.Nonce,
	))
	tx, err = tokenNetwork.SettleChannel(
		auth,
		partnerAddr,
		bp.TransferAmount,
		bp.LocksRoot,
		auth.From,
		bp2.TransferAmount,
		bp2.LocksRoot,
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("SettleChannel gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
}

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
	//buf.Write(c.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(c.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}
func TestCooperateSettleChannel(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	cs := &CoOperativeSettleForContracts{
		Particiant1:         auth.From,
		Participant2:        partnerAddr,
		Participant1Balance: big.NewInt(3),
		Participant2Balance: big.NewInt(totalAmount - 3),
		ChannelIdentifier:   channelID,
		OpenBlockNumber:     openBlockNumber,
		ChainID:             ChainID,
		TokenNetworkAddress: tokenNetworkAddress,
	}
	log.Trace(fmt.Sprintf("cs=\n%s", utils.StringInterface(cs, 3)))
	tx, err := tokenNetwork.CooperativeSettle(
		auth,
		cs.Particiant1,
		cs.Participant1Balance,
		cs.Participant2,
		cs.Participant2Balance,
		cs.sign(TestPrivKey),
		cs.sign(partnerKey),
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
	}
	log.Info(fmt.Sprintf("CooperativeSettle gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
}
func TestRegisterSecret(t *testing.T) {
	secretRegistryAddress, err := tokenNetwork.Secret_registry(nil)
	if err != nil {
		t.Error(err)
		return
	}
	secretRegistry, err := contracts.NewSecretRegistry(secretRegistryAddress, client)
	if err != nil {
		t.Error(err)
		return
	}
	secret := utils.Sha3([]byte("123"))
	t.Logf("secret=%s", secret.String())
	tx, err := secretRegistry.RegisterSecret(auth, secret)
	if err != nil {
		t.Error(err)
		return
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	block, err := secretRegistry.GetSecretRevealBlockHeight(nil, utils.Sha3(secret[:]))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("register block=%s", block)
}

func TestUnlock(t *testing.T) {
	_, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	err = TransferTo(client, TestPrivKey, partnerAddr, big.NewInt(params.Ether))
	if err != nil {
		t.Error(err)
		return
	}
	partnerAuth := bind.NewKeyedTransactor(partnerKey)
	partnerBalance, err := token.BalanceOf(nil, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("before settle partner balance=%s", partnerBalance))
	secretRegistAddress, err := tokenNetwork.Secret_registry(nil)
	if err != nil {
		t.Error(err)
		return
	}
	secretRegistry, err := contracts.NewSecretRegistry(secretRegistAddress, client)
	if err != nil {
		t.Error(err)
		return
	}
	h, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	expiredBlock := h.Number.Int64() + 40
	channelID, _, _, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	if err != nil {
		t.Error(err)
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("close channel successful,gasused=%d,gasLimit=%d", r.GasUsed, tx.Gas()))
	//锁最多是2两个,三个就会失败
	bp2, locks, secrets := NewBalanceProofUpdateForContractsWithLocks(TestPrivKey, partnerKey, channelID, totalLockNumber, expiredBlock)
	//注册密码
	maxLocks := 5
	for i := 0; i < len(secrets); i++ {
		s := secrets[i]
		if i > maxLocks { //最多注册前五个密码,否则太浪费时间了.
			break
		}
		tx, err = secretRegistry.RegisterSecret(auth, s)
		if err != nil {
			t.Error(err)
			return
		}
		r, err = bind.WaitMined(context.Background(), client, tx)
		if err != nil {
			t.Error(err)
			return
		}
		if r.Status != types.ReceiptStatusSuccessful {
			t.Errorf("receipient err ,r=%s", r)
			return
		}
	}
	if len(secrets) < maxLocks {
		maxLocks = len(secrets)
	}
	log.Info(fmt.Sprintf("locksroot=%s", bp2.BalanceData.LocksRoot.String()))
	//提交对方的证据,必须是通道参与方才可以
	tx, err = tokenNetwork.UpdateBalanceProof(partnerAuth,
		auth.From,
		bp2.TransferAmount,
		bp2.LocksRoot,
		bp2.Nonce,
		bp2.AdditionalHash,
		bp2.Signature)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("UpdateBalanceProof successful,gasused=%d,gasLimit=%d", r.GasUsed, tx.Gas()))
	_, blokNumber, _, state, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateClosed {
		t.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateClosed, state)
		return
	}
	mp := mtree.NewMerkleTree(locks)
	lock := locks[0]
	proof := mp.MakeProof(lock.Hash())
	log.Info(fmt.Sprintf("unlockarg,partnerAddr=%s,part2=%s,lock=%s,merkle_proof=%s", partnerAddr.String(),
		auth.From.String(), locks[0], hex.EncodeToString(mtree.Proof2Bytes(proof))))
	tx, err = tokenNetwork.Unlock(
		partnerAuth,
		auth.From,
		bp2.TransferAmount,
		big.NewInt(lock.Expiration),
		lock.Amount,
		lock.LockSecretHash,
		mtree.Proof2Bytes(proof),
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("unlock success,gasUsed=%d,gasLimit=%d,txhash=%s", r.GasUsed, tx.Gas(), tx.Hash().String()))
	partnerBalance, err = token.BalanceOf(nil, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("after unlock partner balance balance=%s", partnerBalance)
	log.Info("waiting settle...")
	for {
		var h *types.Header
		h, err = client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		if h.Number.Int64() > int64(blokNumber+punishBlockNumber) {
			//could settle
			break
		}
		time.Sleep(time.Second)
	}
	tx, err = tokenNetwork.SettleChannel(
		auth,
		partnerAddr,
		bp.TransferAmount,
		bp.LocksRoot,
		auth.From,
		//要用更新后的 transfer amount 了.
		bp2.TransferAmount.Add(bp2.TransferAmount, big.NewInt(int64(1))),
		bp2.LocksRoot,
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("settle channel complete ,gasused=%d,gasLimit=%d", r.GasUsed, tx.Gas()))
	partnerBalance, err = token.BalanceOf(nil, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("after settle partner balance=%s", partnerBalance)

}

type WithDraw1ForContract struct {
	Participant1         common.Address
	Participant2         common.Address
	Participant1Deposit  *big.Int
	Participant2Deposit  *big.Int
	Participant1Withdraw *big.Int
	ChannelIdentifier    contracts.ChannelIdentifier
	OpenBlockNumber      uint64
	TokenNetworkAddress  common.Address
	ChainID              *big.Int
}

func (w *WithDraw1ForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(w.Participant1[:])
	buf.Write(utils.BigIntTo32Bytes(w.Participant1Deposit))
	buf.Write(w.Participant2[:])
	buf.Write(utils.BigIntTo32Bytes(w.Participant2Deposit))
	buf.Write(utils.BigIntTo32Bytes(w.Participant1Withdraw))
	buf.Write(w.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, w.OpenBlockNumber)
	//buf.Write(w.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(w.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}

type WithDraw2ForContract struct {
	Participant1         common.Address
	Participant2         common.Address
	Participant1Deposit  *big.Int
	Participant2Deposit  *big.Int
	Participant1Withdraw *big.Int
	Participant2Withdraw *big.Int
	ChannelIdentifier    contracts.ChannelIdentifier
	OpenBlockNumber      uint64
	TokenNetworkAddress  common.Address
	ChainID              *big.Int
}

func (w *WithDraw2ForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(w.Participant1[:])
	buf.Write(utils.BigIntTo32Bytes(w.Participant1Deposit))
	buf.Write(w.Participant2[:])
	buf.Write(utils.BigIntTo32Bytes(w.Participant2Deposit))
	buf.Write(utils.BigIntTo32Bytes(w.Participant1Withdraw))
	buf.Write(utils.BigIntTo32Bytes(w.Participant2Withdraw))
	buf.Write(w.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, w.OpenBlockNumber)
	//buf.Write(w.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(w.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}
func TestWithdraw(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	w1 := &WithDraw1ForContract{
		Participant1:         auth.From,
		Participant2:         partnerAddr,
		Participant1Withdraw: big.NewInt(1),
		Participant1Deposit:  big.NewInt(totalAmount / 2),
		Participant2Deposit:  big.NewInt(totalAmount / 2),
		ChannelIdentifier:    channelID,
		OpenBlockNumber:      openBlockNumber,
		ChainID:              ChainID,
		TokenNetworkAddress:  tokenNetworkAddress,
	}
	w2 := &WithDraw2ForContract{
		Participant1:         auth.From,
		Participant2:         partnerAddr,
		Participant1Withdraw: big.NewInt(1),
		Participant2Withdraw: big.NewInt(1),
		Participant1Deposit:  big.NewInt(totalAmount / 2),
		Participant2Deposit:  big.NewInt(totalAmount / 2),
		ChannelIdentifier:    channelID,
		OpenBlockNumber:      openBlockNumber,
		ChainID:              ChainID,
		TokenNetworkAddress:  tokenNetworkAddress,
	}
	log.Trace(fmt.Sprintf("w1=\n%s", utils.StringInterface(w1, 3)))
	log.Trace(fmt.Sprintf("w2=\n%s", utils.StringInterface(w2, 3)))
	log.Trace(fmt.Sprintf("WithDraw call, participant1=%s,participant2=%s,"+
		"p1deposit=%s,p2deposit=%s,p1withdarw=%s,p2withdraw=%s,"+
		"p1sig=0x%s,p2sig=0x%s ",
		w2.Participant1.String(),
		w2.Participant2.String(),
		w2.Participant1Deposit,
		w2.Participant2Deposit,
		w2.Participant2Withdraw,
		w2.Participant2Withdraw,
		hex.EncodeToString(w1.sign(TestPrivKey)),
		hex.EncodeToString(w2.sign(partnerKey)),
	))
	tx, err := tokenNetwork.WithDraw(
		auth,
		w2.Participant1,
		w2.Participant1Deposit,
		w2.Participant1Withdraw,
		w2.Participant2,
		w2.Participant2Deposit,
		w2.Participant2Withdraw,
		w1.sign(TestPrivKey),
		w2.sign(partnerKey),
	)
	if err != nil {
		t.Errorf("WithDraw failed err %s", err)
		return
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
	}
	log.Info(fmt.Sprintf("WithDraw complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
}

func TestSignature(t *testing.T) {
	a := big.NewInt(1)
	data := utils.BigIntTo32Bytes(a)
	sig, err := utils.SignData(TestPrivKey, data[:])
	if err != nil {
		t.Error(err)
		return
	}
	log.Trace(fmt.Sprintf("sig=%v", sig))
	log.Trace(fmt.Sprintf("from=%s,sig=%s", auth.From.String(), hex.EncodeToString(sig)))
}

type unlockDelegateForContract struct {
	Agent             common.Address
	Expiraition       int64
	Amount            *big.Int
	SecretHash        common.Hash
	ChannelIdentifier contracts.ChannelIdentifier
	OpenBlockNumber   uint64
	MerkleProof       []byte
}

func (u *unlockDelegateForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(u.Agent[:])
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(u.Expiraition)))
	buf.Write(utils.BigIntTo32Bytes(u.Amount))
	buf.Write(u.SecretHash[:])
	buf.Write(u.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, u.OpenBlockNumber)
	//buf.Write(u.MerkleProof)
	buf.Write(utils.BigIntTo32Bytes(ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}

type ObseleteUnlockForContract struct {
	ChannelIdentifier            contracts.ChannelIdentifier
	BeneficiaryAddress           common.Address
	LockHash                     common.Hash
	BeneficiaryTransferredAmount *big.Int
	BeneficiaryNonce             *big.Int
	AdditionalHash               common.Hash
	TokenNetworkAddress          common.Address
	ChainID                      *big.Int
	OpenBlockNumber              uint64
	MerkleProof                  []byte
}

func (w *ObseleteUnlockForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(w.LockHash[:])
	buf.Write(w.ChannelIdentifier[:])
	binary.Write(buf, binary.BigEndian, w.OpenBlockNumber)
	//buf.Write(w.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(w.ChainID))
	buf.Write(w.AdditionalHash[:])
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}
func calcLockHash(lock *mtree.Lock) common.Hash {
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(big.NewInt(lock.Expiration)))
	buf.Write(utils.BigIntTo32Bytes(lock.Amount))
	buf.Write(lock.LockSecretHash[:])
	return utils.Sha3(buf.Bytes())
}
func TestPunishObsoleteUnlock(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	partnerBalance, err := token.BalanceOf(nil, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	err = TransferTo(client, TestPrivKey, partnerAddr, big.NewInt(params.Ether))
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("before settle partner balance=%s", partnerBalance))
	secretRegistAddress, err := tokenNetwork.Secret_registry(nil)
	if err != nil {
		t.Error(err)
		return
	}
	secretRegistry, err := contracts.NewSecretRegistry(secretRegistAddress, client)
	if err != nil {
		t.Error(err)
		return
	}
	h, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		t.Error(err)
		return
	}
	expiredBlock := h.Number.Int64() + 40
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
	if err != nil {
		t.Error(err)
	}
	r, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("close channel successful,gasused=%d,gasLimit=%d", r.GasUsed, tx.Gas()))

	//锁最多是2两个,三个就会失败
	bp2, locks, secrets := NewBalanceProofUpdateForContractsWithLocks(TestPrivKey, partnerKey, channelID, totalLockNumber, expiredBlock)
	//注册密码
	for i := 0; i < len(secrets); i++ {
		s := secrets[i]
		if i > 5 { //最多注册前五个密码,否则太浪费时间了.
			break
		}
		tx, err = secretRegistry.RegisterSecret(auth, s)
		if err != nil {
			t.Error(err)
			return
		}
		r, err = bind.WaitMined(context.Background(), client, tx)
		if err != nil {
			t.Error(err)
			return
		}
		if r.Status != types.ReceiptStatusSuccessful {
			t.Errorf("receipient err ,r=%s", r)
			return
		}
	}
	log.Info(fmt.Sprintf("locksroot=%s", bp2.LocksRoot.String()))
	//提交对方的证据
	tx, err = tokenNetwork.UpdateBalanceProof(bind.NewKeyedTransactor(partnerKey), auth.From, bp2.TransferAmount, bp2.LocksRoot, bp2.Nonce, bp2.AdditionalHash, bp2.Signature)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("UpdateBalanceProofDelegate successful,gasused=%d,gasLimit=%d", r.GasUsed, tx.Gas()))
	/*
			unlock delegate
		只解锁第一个锁
	*/
	lock := locks[0]
	m := mtree.NewMerkleTree(locks)
	uf := &unlockDelegateForContract{
		Agent:             auth.From,
		Expiraition:       lock.Expiration,
		Amount:            lock.Amount,
		SecretHash:        lock.LockSecretHash,
		ChannelIdentifier: channelID,
		OpenBlockNumber:   openBlockNumber,
		MerkleProof:       mtree.Proof2Bytes(m.MakeProof(lock.Hash())),
	}
	log.Info(fmt.Sprintf("UnlockDelegate arg ,participant=%s,delegator=%s,transferAmount=%s,expiration=%d,amount=%s,"+
		"locksecret=%s,merkleproof=%s,signature=%s",
		auth.From.String(),
		partnerAddr.String(),
		bp2.TransferAmount,
		lock.Expiration,
		lock.Amount,
		lock.LockSecretHash.String(),
		hex.EncodeToString(uf.MerkleProof),
		hex.EncodeToString(uf.sign(partnerKey)),
	))
	tx, err = tokenNetwork.UnlockDelegate(auth,
		auth.From,
		partnerAddr,
		bp2.TransferAmount,
		big.NewInt(lock.Expiration),
		lock.Amount,
		lock.LockSecretHash,
		uf.MerkleProof,
		uf.sign(partnerKey),
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("unlockdelegate gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	lockhash := calcLockHash(lock)
	ou := &ObseleteUnlockForContract{
		ChannelIdentifier:   channelID,
		OpenBlockNumber:     openBlockNumber,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
		BeneficiaryAddress:  auth.From,
		LockHash:            lockhash,
		AdditionalHash:      utils.EmptyHash,
		MerkleProof:         mtree.Proof2Bytes(m.MakeProof(lockhash)),
	}
	log.Info(fmt.Sprintf("PunishObsoleteUnlock,channelid=%s,partnerAddr=%s,part2=%s,locksroot=%s", common.Hash(channelID).String(), partnerAddr.String(),
		auth.From.String(), ou.LockHash.String()))
	tx, err = tokenNetwork.PunishObsoleteUnlock(
		auth,
		auth.From,
		partnerAddr,
		lockhash,
		ou.AdditionalHash,
		ou.sign(partnerKey),
	)
	if err != nil {
		t.Error(err)
		return
	}
	r, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if r.Status != types.ReceiptStatusSuccessful {
		t.Errorf("receipient err ,r=%s", r)
		return
	}
	log.Info(fmt.Sprintf("PunishObsoleteUnlock success,gasUsed=%d,gasLimit=%d,txhash=%s", r.GasUsed, tx.Gas(), tx.Hash().String()))
	deposit, balancehash, nonce, err := tokenNetwork.GetChannelParticipantInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("beneficiary deposit=%s,nonce=%d,balance_hash=%s", deposit, nonce, hex.EncodeToString(balancehash[:])))
}
