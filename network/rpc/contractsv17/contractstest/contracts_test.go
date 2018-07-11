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

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contractsv17/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer/mtree"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
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
	tokenNetworkAddress = common.HexToAddress("0x00e72c290ff8429ca2f188cbf192ee811d3cab83")
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
}

/*s
creatAChannelAndDeposit create a channel
1,2之间创建通道,总是都由1作为 tx 发起人
*/
func creatAChannelAndDeposit(account1, account2 common.Address, key1 *ecdsa.PrivateKey, amount int64, conn *ethclient.Client) error {
	log.Trace(fmt.Sprintf("createchannel between %s-%s,tokenNetwork=%s\n", account1.String(), account2.String(), tokenNetworkAddress.String()))
	auth1 := bind.NewKeyedTransactor(key1)
	tx, err := tokenNetwork.OpenChannel(auth1, account1, account2, big.NewInt(int64(settleTimeout)))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("Failed to NewChannel: %v,%s,%s", err, auth1.From.String(), account2.String()))
	}
	ctx := context.Background()
	r, err := bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to NewChannel when mining :%v", err)
	}
	log.Info(fmt.Sprintf("OpenChannel complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
	channelID, settleBlock, _, state, err := tokenNetwork.GetChannelInfo(nil, account1, account2)
	if err != nil {
		return fmt.Errorf("GetChannelInfo %s-%s err %s", utils.APex2(account1), utils.APex2(account2), err)
	}
	log.Trace(fmt.Sprintf("create channel gas %s:%d,tokennetworkaddress=%s\n", tx.Hash().String(), tx.Gas(), tokenNetworkAddress.String()))
	log.Info("NewChannel complete...\n")
	if settleBlock.Int64() != int64(settleTimeout) {
		return fmt.Errorf("settleBlock err expect=%d,got=%s", settleTimeout, settleBlock)
	}
	if state != 1 {
		return fmt.Errorf("")
	}
	tx, err = tokenNetwork.SetTotalDeposit(auth1, channelID, account1, big.NewInt(amount))
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

	tx, err = tokenNetwork.SetTotalDeposit(auth1, channelID, account2, big.NewInt(amount))
	if err != nil {
		return fmt.Errorf("Failed to Deposit2: %v", err)

	}
	ctx = context.Background()
	r, err = bind.WaitMined(ctx, conn, tx)
	if err != nil {
		return fmt.Errorf("failed to Deposit when mining :%v", err)
	}
	log.Info(fmt.Sprintf("Deposit2 complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
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
func getTestOpenChannel(t *testing.T) (channelID *big.Int, partnerAddr common.Address, partnerKey *ecdsa.PrivateKey, err error) {
	testApprove(t)
	partnerKey, partnerAddr = utils.MakePrivateKeyAddress()
	err = creatAChannelAndDeposit(auth.From, partnerAddr, TestPrivKey, totalAmount/2, client)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, settleBlockNumber, _, state, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateOpened {
		err = fmt.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateOpened, state)
		return
	}
	t.Logf("channelID=%s,settleblockNumber=%s,state=%d,err=%v", channelID.String(), settleBlockNumber, state, err)
	return
}
func TestOpenChannel(t *testing.T) {
	_, _, _, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
	}
}

func TestCloseChannel1(t *testing.T) {
	channelID, _, _, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}

	tx, err := tokenNetwork.CloseChannel(auth, channelID, utils.BigInt0, utils.EmptyHash, utils.BigInt0, utils.EmptyHash, nil)
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
	log.Info(fmt.Sprintf("CloseChannel no evidence complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
}

//BalanceProofForContract for contract
type BalanceProofForContract struct {
	TransferAmount      *big.Int
	LocksRoot           common.Hash
	Nonce               *big.Int
	AdditionalHash      common.Hash
	ChannelIdentifier   contracts.ChannelIdentifier
	TokenNetworkAddress common.Address
	ChainID             *big.Int
	Signature           []byte
}

func createPartnerBalanceProof(key *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier) *BalanceProofForContract {
	bp := &BalanceProofForContract{
		TransferAmount:      big.NewInt(10),
		LocksRoot:           utils.EmptyHash,
		Nonce:               big.NewInt(3),
		AdditionalHash:      utils.Sha3([]byte("123")),
		ChannelIdentifier:   channelID,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
	}
	bp.sign(key)
	return bp
}

func encodingLocks(locks []*mtree.Lock) []byte {
	buf := new(bytes.Buffer)
	for _, l := range locks {
		buf.Write(utils.BigIntTo32Bytes(big.NewInt(l.Expiration)))
		buf.Write(utils.BigIntTo32Bytes(l.Amount))
		buf.Write(l.LockSecretHash[:])
	}
	return buf.Bytes()
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
	bp = &BalanceProofForContract{
		TransferAmount:      big.NewInt(10),
		LocksRoot:           m.MerkleRoot(),
		Nonce:               big.NewInt(3),
		AdditionalHash:      utils.Sha3([]byte("123")),
		ChannelIdentifier:   channelID,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
	}
	bp.sign(key)
	return
}
func (b *BalanceProofForContract) sign(key *ecdsa.PrivateKey) {
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(b.TransferAmount))
	buf.Write(b.LocksRoot[:])
	buf.Write(utils.BigIntTo32Bytes(b.Nonce))
	buf.Write(b.AdditionalHash[:])
	buf.Write(utils.BigIntTo32Bytes(b.ChannelIdentifier))
	buf.Write(b.TokenNetworkAddress[:])
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
	channelID, _, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, channelID, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	log.Info(fmt.Sprintf("CloseChannel with evidence  complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
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
func createPartnerBalanceProof2(key *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier) *BalanceProofForContract {
	bp := &BalanceProofForContract{
		TransferAmount:      big.NewInt(10),
		LocksRoot:           utils.EmptyHash,
		Nonce:               big.NewInt(5),
		AdditionalHash:      utils.Sha3([]byte("123")),
		ChannelIdentifier:   channelID,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
	}
	bp.sign(key)
	return bp
}

func NewBalanceProofUpdateForContracts2(closingKey, nonClosingKey *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier) *BalanceProofUpdateForContracts {
	bp1 := createPartnerBalanceProof2(closingKey, channelID)
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
	buf.Write(utils.BigIntTo32Bytes(b.Nonce))
	buf.Write(b.AdditionalHash[:])
	buf.Write(utils.BigIntTo32Bytes(b.ChannelIdentifier))
	buf.Write(b.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(b.ChainID))
	buf.Write(b.Signature)
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	b.NonClosingSignature = sig
}
func TestCloseChannelAndUpdateNonClosingAndSettle(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, channelID, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	log.Info(fmt.Sprintf("CloseChannel   complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
	bp2 := NewBalanceProofUpdateForContracts(TestPrivKey, partnerKey, channelID)
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, channelID, partnerAddr, bp2.LocksRoot, bp2.TransferAmount, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
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
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof   complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
	_, blokNumber, _, state, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
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
		if h.Number.Cmp(blokNumber) > 0 {
			//could settle
			break
		}
		time.Sleep(time.Second)
	}
	tx, err = tokenNetwork.SettleChannel(
		auth,
		partnerAddr,
		auth.From,
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
	log.Info(fmt.Sprintf("SettleChannel  complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
}

//可以 updateNonClosingProof 多次,只要合法,并且 nonce 较大即可.
func TestCloseChannelAndUpdateNonClosingAndSettle2(t *testing.T) {
	channelID, partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, channelID, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	log.Info(fmt.Sprintf("CloseChannel   complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
	bp2 := NewBalanceProofUpdateForContracts(TestPrivKey, partnerKey, channelID)
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, channelID, partnerAddr, bp2.LocksRoot, bp2.TransferAmount, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
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
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof   complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
	//同样的数据在 update, 应该失败.
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, channelID, partnerAddr, bp2.LocksRoot, bp2.TransferAmount, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
	if err == nil {
		t.Error("should fail, nonce only once")
		return
	}
	bp3 := NewBalanceProofUpdateForContracts2(TestPrivKey, partnerKey, channelID)
	//第二次 update, 应该也要成功
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, channelID, partnerAddr, bp3.LocksRoot, bp3.TransferAmount, bp3.Nonce, bp3.AdditionalHash, bp3.Signature, bp3.NonClosingSignature)
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
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof2   complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
}

type CoOperativeSettleForContracts struct {
	Particiant1         common.Address
	Participant2        common.Address
	Participant1Balance *big.Int
	Participant2Balance *big.Int
	ChannelIdentifier   contracts.ChannelIdentifier
	TokenNetworkAddress common.Address
	ChainID             *big.Int
}

func (c *CoOperativeSettleForContracts) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(c.Particiant1[:])
	buf.Write(utils.BigIntTo32Bytes(c.Participant1Balance))
	buf.Write(c.Participant2[:])
	buf.Write(utils.BigIntTo32Bytes(c.Participant2Balance))
	buf.Write(utils.BigIntTo32Bytes(c.ChannelIdentifier))
	buf.Write(c.TokenNetworkAddress[:])
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
		Participant1Balance: big.NewInt(totalAmount / 2),
		Participant2Balance: big.NewInt(totalAmount / 2),
		ChannelIdentifier:   channelID,
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
	log.Info(fmt.Sprintf("CooperativeSettle   complete.. gasLimit=%d,gasUsed=%d", r.GasUsed, tx.Gas()))
}

type WithDraw1ForContract struct {
	Participant1         common.Address
	Participant2         common.Address
	Participant1Deposit  *big.Int
	Participant2Deposit  *big.Int
	Participant1Withdraw *big.Int
	ChannelIdentifier    *big.Int
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
	buf.Write(utils.BigIntTo32Bytes(w.ChannelIdentifier))
	buf.Write(w.TokenNetworkAddress[:])
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
	ChannelIdentifier    *big.Int
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
	buf.Write(utils.BigIntTo32Bytes(w.ChannelIdentifier))
	buf.Write(w.TokenNetworkAddress[:])
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
		ChainID:              ChainID,
		TokenNetworkAddress:  tokenNetworkAddress,
	}
	log.Trace(fmt.Sprintf("w=\n%s", utils.StringInterface(w1, 3)))
	log.Trace(fmt.Sprintf("WithDraw call, participant1=%s,participant2=%s,"+
		"p1deposit=%s,p2deposit=%s,p1withdarw=%s,p2withdraw=%s,"+
		"p1sig=0x%s,p2sig=0x%s,channelid=%s",
		w2.Participant1.String(),
		w2.Participant2.String(),
		w2.Participant1Deposit,
		w2.Participant2Deposit,
		w2.Participant2Withdraw,
		w2.Participant2Withdraw,
		hex.EncodeToString(w1.sign(TestPrivKey)),
		hex.EncodeToString(w2.sign(partnerKey)),
		w2.ChannelIdentifier,
	))
	tx, err := tokenNetwork.WithDraw(
		auth,
		w2.Participant1,
		w2.Participant2,
		w2.Participant1Deposit,
		w2.Participant2Deposit,
		w2.Participant1Withdraw,
		w2.Participant2Withdraw,
		w1.sign(TestPrivKey),
		w2.sign(partnerKey),
		w2.ChannelIdentifier,
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
	log.Info(fmt.Sprintf("RegisterSecret success.. gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	block, err := secretRegistry.GetSecretRevealBlockHeight(nil, utils.Sha3(secret[:]))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("register block=%s", block)
}

func TestUnlock(t *testing.T) {
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
	tx, err := tokenNetwork.CloseChannel(auth, channelID, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, channelID, partnerAddr, bp2.LocksRoot, bp2.TransferAmount, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
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
	_, _, locksroot, transferAmount, _, err := tokenNetwork.GetChannelParticipantInfo(nil, channelID, auth.From)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof successful,gasused=%d,gasLimit=%d,locksroot=%s,transferamount=%s", r.GasUsed, tx.Gas(), hex.EncodeToString(locksroot[:]), transferAmount))
	lockbytes := encodingLocks(locks)
	log.Info(fmt.Sprintf("unlockarg,channelid=%s,partnerAddr=%s,part2=%s,locks=%s", channelID, partnerAddr.String(),
		auth.From.String(), hex.EncodeToString(lockbytes)))
	tx, err = tokenNetwork.Unlock(auth, channelID, auth.From, lockbytes)
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
	_, _, locksroot, transferAmount, _, err = tokenNetwork.GetChannelParticipantInfo(nil, channelID, auth.From)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("after unlock locksroot=%s,transferamount=%s", hex.EncodeToString(locksroot[:]), transferAmount))
	_, blokNumber, _, state, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateClosed {
		t.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateClosed, state)
		return
	}
	log.Info("waiting settle...")
	for {
		var h *types.Header
		h, err = client.HeaderByNumber(context.Background(), nil)
		if err != nil {
			t.Error(err)
			return
		}
		if h.Number.Cmp(blokNumber) > 0 {
			//could settle
			break
		}
		time.Sleep(time.Second)
	}
	tx, err = tokenNetwork.SettleChannel(
		auth,
		partnerAddr,
		auth.From,
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

type ObseleteUnlockForContract struct {
	ChannelIdentifier   *big.Int
	BeneficiaryAddress  common.Address
	LockHash            common.Hash
	AdditionalHash      common.Hash
	TokenNetworkAddress common.Address
	ChainID             *big.Int
	MerkleProof         []byte
}

func (w *ObseleteUnlockForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(w.LockHash[:])
	buf.Write(utils.BigIntTo32Bytes(w.ChannelIdentifier))
	buf.Write(w.TokenNetworkAddress[:])
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
	tx, err := tokenNetwork.CloseChannel(auth, channelID, bp.TransferAmount, bp.LocksRoot, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, channelID, partnerAddr, bp2.LocksRoot, bp2.TransferAmount, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
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
	_, _, locksroot, transferAmount, _, err := tokenNetwork.GetChannelParticipantInfo(nil, channelID, auth.From)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof successful,gasused=%d,gasLimit=%d,locksroot=%s,transferamount=%s", r.GasUsed, tx.Gas(), hex.EncodeToString(locksroot[:]), transferAmount))
	m := mtree.NewMerkleTree(locks)
	lockhash := calcLockHash(locks[0])
	ou := &ObseleteUnlockForContract{
		ChannelIdentifier:   channelID,
		TokenNetworkAddress: tokenNetworkAddress,
		ChainID:             ChainID,
		BeneficiaryAddress:  auth.From,
		LockHash:            lockhash,
		AdditionalHash:      utils.EmptyHash,
		MerkleProof:         mtree.Proof2Bytes(m.MakeProof(lockhash)),
	}
	log.Info(fmt.Sprintf("unlockarg,channelid=%s,partnerAddr=%s,part2=%s,proof=%s", channelID, partnerAddr.String(),
		auth.From.String(), hex.EncodeToString(ou.MerkleProof)))
	tx, err = tokenNetwork.PunishObsoleteUnlock(auth, channelID, auth.From, lockhash, ou.AdditionalHash,
		ou.sign(partnerKey), ou.MerkleProof,
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
	_, _, locksroot, transferAmount, _, err = tokenNetwork.GetChannelParticipantInfo(nil, channelID, auth.From)
	if err != nil {
		t.Error(err)
		return
	}
	log.Info(fmt.Sprintf("after PunishObsoleteUnlock locksroot=%s,transferamount=%s", hex.EncodeToString(locksroot[:]), transferAmount))
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
