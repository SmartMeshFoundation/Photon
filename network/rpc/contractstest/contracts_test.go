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
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
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
var settleTimeout = 50
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
	tokenNetworkAddress = common.HexToAddress("0xf9b14f75fe8cf31cdae0f010f99e0a2c1e843ee3")
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
	log.Info(fmt.Sprintf("OpenChannel gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	channelID, settleBlock, state, err := tokenNetwork.GetChannelInfo(nil, account1, account2)
	if err != nil {
		return fmt.Errorf("GetChannelInfo %s-%s err %s", utils.APex2(account1), utils.APex2(account2), err)
	}
	log.Trace(fmt.Sprintf("create channel gas %s:%d,channel identifier=0x%s,tokennetworkaddress=%s\n", tx.Hash().String(), tx.Gas(), hex.EncodeToString(channelID[:]), tokenNetworkAddress.String()))
	log.Info("NewChannel complete...\n")
	if settleBlock.Int64() != int64(settleTimeout) {
		return fmt.Errorf("settleBlock err expect=%d,got=%s", settleTimeout, settleBlock)
	}
	if state != 1 {
		return fmt.Errorf("")
	}
	tx, err = tokenNetwork.SetTotalDeposit(auth1, account1, big.NewInt(amount), account2)
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

	tx, err = tokenNetwork.SetTotalDeposit(auth1, account2, big.NewInt(amount), account1)
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

//跑一次就够了,这样后续创建通道就不用每次 appro
func TestApprove(t *testing.T) {
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
func getTestOpenChannel(t *testing.T) (partnerAddr common.Address, partnerKey *ecdsa.PrivateKey, err error) {
	partnerKey, partnerAddr = utils.MakePrivateKeyAddress()
	err = creatAChannelAndDeposit(auth.From, partnerAddr, TestPrivKey, totalAmount/2, client)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, settleBlockNumber, state, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	if state != contracts.ChannelStateOpened {
		err = fmt.Errorf("channel state err expect=%d,got=%d", contracts.ChannelStateOpened, state)
		return
	}
	t.Logf("channelID=%s,settleblockNumber=%s,state=%d,err=%s", hex.EncodeToString(channelID[:]), settleBlockNumber, state, err)
	return
}
func TestOpenChannel(t *testing.T) {
	_, _, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
	}
}

func TestCloseChannel1(t *testing.T) {
	partnerAddr, _, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}

	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, utils.EmptyHash, utils.BigInt0, utils.EmptyHash, nil)
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
	LockedAmount   *big.Int
	LocksRoot      common.Hash
}

func (b *BalanceData) Hash() common.Hash {
	buf := new(bytes.Buffer)
	buf.Write(utils.BigIntTo32Bytes(b.TransferAmount))
	buf.Write(utils.BigIntTo32Bytes(b.LockedAmount))
	buf.Write(b.LocksRoot[:])
	return utils.Sha3(buf.Bytes())
}

//BalanceProofForContract for contract
type BalanceProofForContract struct {
	BalanceHash         common.Hash
	Nonce               *big.Int
	AdditionalHash      common.Hash
	ChannelIdentifier   contracts.ChannelIdentifier
	TokenNetworkAddress common.Address
	ChainID             *big.Int
	Signature           []byte
	BalanceData         *BalanceData
}

func createPartnerBalanceProof(key *ecdsa.PrivateKey, channelID contracts.ChannelIdentifier) *BalanceProofForContract {
	bd := &BalanceData{
		TransferAmount: big.NewInt(10),
		LockedAmount:   utils.BigInt0,
		LocksRoot:      utils.EmptyHash,
	}
	bp := &BalanceProofForContract{
		BalanceHash:         bd.Hash(),
		BalanceData:         bd,
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
	bd := &BalanceData{
		TransferAmount: big.NewInt(10),
		LockedAmount:   big.NewInt(int64(lockNumber)),
		LocksRoot:      m.MerkleRoot(),
	}
	bp = &BalanceProofForContract{
		BalanceHash:         bd.Hash(),
		BalanceData:         bd,
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
	buf.Write(b.BalanceHash[:])
	buf.Write(utils.BigIntTo32Bytes(b.Nonce))
	buf.Write(b.AdditionalHash[:])
	buf.Write(b.ChannelIdentifier[:])
	buf.Write(b.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(b.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	b.Signature = sig
}
func TestCloseChannel2(t *testing.T) {
	partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.BalanceHash, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	buf.Write(b.BalanceHash[:])
	buf.Write(utils.BigIntTo32Bytes(b.Nonce))
	buf.Write(b.AdditionalHash[:])
	buf.Write(b.ChannelIdentifier[:])
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
	partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.BalanceHash, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	bp2 := NewBalanceProofUpdateForContracts(TestPrivKey, partnerKey, channelID)
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, auth.From, partnerAddr, bp2.BalanceHash, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
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
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
	_, blokNumber, state, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
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
		bp.BalanceData.TransferAmount,
		bp.BalanceData.LockedAmount,
		bp.BalanceData.LocksRoot,
		auth.From,
		bp2.BalanceData.TransferAmount,
		bp2.BalanceData.LockedAmount,
		bp2.BalanceData.LocksRoot,
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
	buf.Write(c.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(c.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}
func TestCooperateSettleChannel(t *testing.T) {
	partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
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

type WithDrawForContract struct {
	ChannelIdentifier   contracts.ChannelIdentifier
	Participant         common.Address
	TotalWithDraw       *big.Int
	TokenNetworkAddress common.Address
	ChainID             *big.Int
}

func (w *WithDrawForContract) sign(key *ecdsa.PrivateKey) []byte {
	buf := new(bytes.Buffer)
	buf.Write(w.Participant[:])
	buf.Write(utils.BigIntTo32Bytes(w.TotalWithDraw))
	buf.Write(w.ChannelIdentifier[:])
	buf.Write(w.TokenNetworkAddress[:])
	buf.Write(utils.BigIntTo32Bytes(w.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	return sig
}
func TestSetTotalWithdraw(t *testing.T) {
	partnerAddr, partnerKey, err := getTestOpenChannel(t)
	if err != nil {
		t.Error(err)
		return
	}
	channelID, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	w := &WithDrawForContract{
		Participant:         auth.From,
		TotalWithDraw:       big.NewInt(totalAmount + 1),
		ChannelIdentifier:   channelID,
		ChainID:             ChainID,
		TokenNetworkAddress: tokenNetworkAddress,
	}
	log.Trace(fmt.Sprintf("w=\n%s", utils.StringInterface(w, 3)))
	tx, err := tokenNetwork.SetTotalWithdraw(
		auth,
		w.Participant,
		w.TotalWithDraw,
		partnerAddr,
		w.sign(TestPrivKey),
		w.sign(partnerKey),
	)
	if err == nil {
		t.Errorf("too big amount should faile")
		return
	}
	w.TotalWithDraw = big.NewInt(10)
	log.Trace(fmt.Sprintf("w=\n%s", utils.StringInterface(w, 3)))
	tx, err = tokenNetwork.SetTotalWithdraw(
		auth,
		w.Participant,
		w.TotalWithDraw,
		partnerAddr,
		w.sign(TestPrivKey),
		w.sign(partnerKey),
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
	log.Info(fmt.Sprintf("withdraw gasLimit=%d,gasUsed=%d", tx.Gas(), r.GasUsed))
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
	partnerAddr, partnerKey, err := getTestOpenChannel(t)
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
	channelID, _, _, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
	if err != nil {
		t.Error(err)
		return
	}
	bp := createPartnerBalanceProof(partnerKey, contracts.ChannelIdentifier(channelID))
	tx, err := tokenNetwork.CloseChannel(auth, partnerAddr, bp.BalanceHash, bp.Nonce, bp.AdditionalHash, bp.Signature)
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
	log.Info(fmt.Sprintf("locksroot=%s", bp2.BalanceData.LocksRoot.String()))
	//提交对方的证据
	tx, err = tokenNetwork.UpdateNonClosingBalanceProof(auth, auth.From, partnerAddr, bp2.BalanceHash, bp2.Nonce, bp2.AdditionalHash, bp2.Signature, bp2.NonClosingSignature)
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
	log.Info(fmt.Sprintf("UpdateNonClosingBalanceProof successful,gasused=%d,gasLimit=%d", r.GasUsed, tx.Gas()))
	_, blokNumber, state, err := tokenNetwork.GetChannelInfo(nil, auth.From, partnerAddr)
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
		bp.BalanceData.TransferAmount,
		bp.BalanceData.LockedAmount,
		bp.BalanceData.LocksRoot,
		auth.From,
		bp2.BalanceData.TransferAmount,
		bp2.BalanceData.LockedAmount,
		bp2.BalanceData.LocksRoot,
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
	lockbytes := encodingLocks(locks)
	log.Info(fmt.Sprintf("unlockarg,partnerAddr=%s,part2=%s,locks=%s", partnerAddr.String(),
		auth.From.String(), hex.EncodeToString(lockbytes)))
	tx, err = tokenNetwork.Unlock(auth, partnerAddr, auth.From, lockbytes)
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
}

//func TestUnlock2(t *testing.T) {
//	//locksroot:0x3895b20a2dda9ce0f42e6ed97fc7d9c75546862158dc2e18304947e4934f8391
//	partnerAddr := common.HexToAddress("0x3C20dd31a7e7E5703D08CBA501D315E33d9A0Ddf")
//	lockbytes, err := hex.DecodeString("00000000000000000000000000000000000000000000000000000000002453240000000000000000000000000000000000000000000000000000000000000001c5cce3a57696d6ef896e65f3229880fecf054f61ed5b6e4bd7120d91e66e973600000000000000000000000000000000000000000000000000000000002453240000000000000000000000000000000000000000000000000000000000000001dad1107dd0857192f3c8dc22b538a202a17d070b0a66f19461fa7fd1db0ad7fc0000000000000000000000000000000000000000000000000000000000245324000000000000000000000000000000000000000000000000000000000000000123aba73277e9f38c82311e5e1d542ceff6afdde11bb4bc438a4c712ad5f9be2f")
//	tx, err := tokenNetwork.Unlock(auth, partnerAddr, auth.From, lockbytes)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	r, err := bind.WaitMined(context.Background(), client, tx)
//	if err != nil {
//		t.Error(err)
//		return
//	}
//	if r.Status != types.ReceiptStatusSuccessful {
//		t.Errorf("receipient err ,r=%s", r)
//		return
//	}
//}
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
