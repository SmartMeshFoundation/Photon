package pfsproxy

import (
	"bytes"
	"crypto/ecdsa"
	"os"
	"testing"

	"math/big"

	"encoding/binary"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/codefortest"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

var testPfgHost = "http://transport01.smartmesh.cn:7002"

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

func TestPfsClient_SubmitBalance(t *testing.T) {
	if testing.Short() {
		return
	}
	params.ChainID = big.NewInt(8888)
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	bob, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x201B20123b3C489b47Fde27ce5b451a0fA55FD60"))
	if err != nil {
		t.Error(err)
		return
	}
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)

	//nonce := big.NewInt(6)
	//transferAmount := big.NewInt(150)
	//lockAmount := big.NewInt(0)
	//openBlockNumber := big.NewInt(7218036)
	//locksroot := utils.EmptyHash
	//channelIdentifier := common.HexToHash("0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a")
	//additionHash := common.HexToHash("0x3f8cd08a836d1993c597016d323a23d7b64b9277b2d35ab81c06fc1fabbfb961")
	//signature := common.Hex2Bytes("33UcHuv4YMjMeDhLla7KQOkpmASA44HgOmRQSMHFUOwXG9apEdRjrWQ1seua4FlE+tsTsPsZZgWdvBxxOvjrchw=")
	//err = c.SubmitBalance(nonce, transferAmount, lockAmount, openBlockNumber, locksroot, channelIdentifier, additionHash, signature[:])
	//if err != nil {
	//	t.Error(err)
	//}
	nonce := big.NewInt(10)
	transferAmount := big.NewInt(210)
	lockAmount := big.NewInt(0)
	openBlockNumber := big.NewInt(7218036)
	locksroot := utils.EmptyHash
	channelIdentifier := common.HexToHash("0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a")
	additionHash := common.HexToHash("0x2ce9a9b8ee85aae05d1847d69573029919131635ad62e8322432dde43c07075c")
	bp := createPartnerBalanceProof(bob, transferAmount, locksroot, additionHash, nonce.Uint64(), openBlockNumber, channelIdentifier)

	err = c.SubmitBalance(nonce.Uint64(), transferAmount, lockAmount, openBlockNumber.Int64(), locksroot, channelIdentifier, additionHash, bp.Signature)
	if err != nil {
		t.Error(err)
	}
}

//BalanceProofForContract for contract
type BalanceProofForContract struct {
	AdditionalHash      common.Hash
	ChannelIdentifier   common.Hash
	TokenNetworkAddress common.Address
	ChainID             *big.Int
	Signature           []byte
	OpenBlockNumber     uint64
	Nonce               uint64
	TransferAmount      *big.Int
	LocksRoot           common.Hash
}

func (b *BalanceProofForContract) sign(key *ecdsa.PrivateKey) {
	buf := new(bytes.Buffer)
	_, err := buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte("176"))
	_, err = buf.Write(utils.BigIntTo32Bytes(b.TransferAmount))
	_, err = buf.Write(b.LocksRoot[:])
	err = binary.Write(buf, binary.BigEndian, b.Nonce)
	_, err = buf.Write(b.AdditionalHash[:])
	_, err = buf.Write(b.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, b.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(b.ChainID))
	sig, err := utils.SignData(key, buf.Bytes())
	if err != nil {
		panic(err)
	}
	b.Signature = sig
}

func createPartnerBalanceProof(
	partner codefortest.TestAccount,
	transferAmount *big.Int, locksroot common.Hash, additionalHash common.Hash, nonce uint64, openBlockNumber *big.Int, channelID common.Hash) *BalanceProofForContract {
	bp := &BalanceProofForContract{
		OpenBlockNumber:   openBlockNumber.Uint64(),
		AdditionalHash:    additionalHash,
		ChannelIdentifier: channelID,
		ChainID:           params.ChainID,
		Nonce:             nonce,
		TransferAmount:    transferAmount,
		LocksRoot:         locksroot,
	}
	bp.sign(partner.PrivateKey)
	return bp
}

func TestPfsClient_FindPath(t *testing.T) {
	if testing.Short() {
		return
	}
	tokenAddress := common.HexToAddress("0x76fCe6fF759B208D27E4D48828F820d79d1719f3")
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	bob, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x201B20123b3C489b47Fde27ce5b451a0fA55FD60"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	routes, err := c.FindPath(alice.Address, bob.Address, tokenAddress, big.NewInt(20), true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(routes)
}

func TestPfsClient_SetAccountFee(t *testing.T) {
	if testing.Short() {
		return
	}
	feeConstant := big.NewInt(5)
	feePercent := int64(10000)
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	err = c.SetAccountFee(feeConstant, feePercent)
	if err != nil {
		t.Error(err)
	}
}

func TestPfsClient_GetAccountFee(t *testing.T) {
	if testing.Short() {
		return
	}
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	//channelIdentifier := common.HexToHash("0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a")
	feeConstant, feePercent, err := c.GetAccountFee()
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(feeConstant, feePercent)
}

func TestPfsClient_SetTokenFee(t *testing.T) {
	if testing.Short() {
		return
	}
	feeConstant := big.NewInt(6)
	feePercent := int64(30000)
	tokenAddress := common.HexToAddress("0x76fCe6fF759B208D27E4D48828F820d79d1719f3")
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	err = c.SetTokenFee(feeConstant, feePercent, tokenAddress)
	if err != nil {
		t.Error(err)
	}
}

func TestPfsClient_GetTokenFee(t *testing.T) {
	if testing.Short() {
		return
	}
	tokenAddress := common.HexToAddress("0x76fCe6fF759B208D27E4D48828F820d79d1719f3")
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	feeConstant, feePercent, err := c.GetTokenFee(tokenAddress)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(feeConstant, feePercent)
}

func TestPfsClient_SetChannelFee(t *testing.T) {
	if testing.Short() {
		return
	}
	feeConstant := big.NewInt(5)
	feePercent := int64(20000)
	channelIdentifier := common.HexToHash("0x640b3a6c160eadc37f133400b6a6be62d4d8a2b7ccd67beb04426e84251455ea")
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	err = c.SetChannelFee(feeConstant, feePercent, channelIdentifier)
	if err != nil {
		t.Error(err)
	}
}

func TestPfsClient_GetChannelFee(t *testing.T) {
	if testing.Short() {
		return
	}
	channelIdentifier := common.HexToHash("0x640b3a6c160eadc37f133400b6a6be62d4d8a2b7ccd67beb04426e84251455ea")
	alice, err := codefortest.GetAccountsByAddress(common.HexToAddress("0x10b256b3C83904D524210958FA4E7F9cAFFB76c6"))
	c := NewPfsProxy(testPfgHost, alice.PrivateKey)
	//channelIdentifier := common.HexToHash("0x622924d11071238ac70c39b508c37216d1a392097a80b26f5299a8d8f4bc0b7a")
	feeConstant, feePercent, err := c.GetChannelFee(channelIdentifier)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Println(feeConstant, feePercent)
}
