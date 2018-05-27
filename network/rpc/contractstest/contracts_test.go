package contractstest

import (
	"os"
	"testing"

	"math/big"

	"bytes"
	"crypto/ecdsa"
	"encoding/binary"

	"context"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/encoding"
	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc"
	"github.com/SmartMeshFoundation/SmartRaiden/network/rpc/contracts"
	"github.com/SmartMeshFoundation/SmartRaiden/transfer"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

/*
CHANNEL	0x97370d43844b11a659de7dccce9cd086a3ab0262f4a91ff679951cd734e3ddcb
REGISTRY	0xBBD9D2A904a7863e044E942663DCF02A2E9696FB
DISCOVERY	0x8cA6Ca4139909F69b053126d0818A8C9BD1e0573
ETHRPCENDPOINT	ws://127.0.0.1:8546
TOKEN	0x6b0aDe8E98E73fC65AA28A0C88683ed91f2994eD
MANAGER	0x280A6d22a1dF783EB002C5EC7400Bdf592eeEe93
*/

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}
func TestUpdateTransferDelegate(t *testing.T) {
	var tx *types.Transaction
	token := common.HexToAddress(os.Getenv("TOKEN"))
	client, err := ethclient.Dial(rpc.TestRPCEndpoint)
	auth := bind.NewKeyedTransactor(rpc.TestPrivKey)
	if err != nil {
		t.Error(err)
		return
	}
	key1, addr1 := utils.MakePrivateKeyAddress()
	key2, addr2 := utils.MakePrivateKeyAddress()
	chaddr, err := rpc.CreateChannelBetweenAddress(client, addr1, addr2, key1, key2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("new channel addr=%s", chaddr.String())
	//transfer from 1 to 2
	dt := encoding.NewDirectTransfer(1, 1, token, chaddr, big.NewInt(30), addr2, utils.EmptyHash)
	err = dt.Sign(key1, dt)
	if err != nil {
		t.Error(err)
		return
	}
	assert.EqualValues(t, dt.Sender, addr1)
	bp := transfer.NewBalanceProofStateFromEnvelopMessage(dt)
	nonclosingSignature, err := signFor3rd(bp.Nonce, bp.TransferAmount, bp.LocksRoot, bp.ChannelAddress, bp.MessageHash, bp.Signature, auth.From, key2)
	if err != nil {
		t.Error(err)
		return
	}
	ch, err := contracts.NewNettingChannelContract(chaddr, client)
	if err != nil {
		t.Error(err)
	}
	tx, err = ch.Close(bind.NewKeyedTransactor(key1), 0, big.NewInt(0), [32]byte{}, [32]byte{}, make([]byte, 0))
	if err != nil {
		t.Error(err)
		return
	}
	_, err = bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	settletime, err := ch.SettleTimeout(nil)
	if err != nil {
		t.Error(err)
		return
	}
	closed, _ := ch.Closed(nil)
	t.Logf("closed=%s,settle=%s,now=%s", closed, settletime, time.Now())
	time.Sleep(time.Second * 2 * time.Duration(settletime.Int64()/2+2))
	t.Logf("now=%s,nonclosing=%d", time.Now(), len(nonclosingSignature))
	tx, err = ch.UpdateTransferDelegate(auth, uint64(bp.Nonce), bp.TransferAmount, bp.LocksRoot, bp.MessageHash, bp.Signature, nonclosingSignature)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("UpdateTransferDelegate tx hash=%s", tx.Hash().String())
	receipts, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if receipts.Status != 1 {
		t.Errorf("UpdateTransferDelegate finished err %s", receipts.String())
		return
	}
	time.Sleep(time.Second * 2 * time.Duration(settletime.Int64()/2+2))
	tp, _ := contracts.NewToken(token, client)
	balance1, err := tp.BalanceOf(nil, addr1)
	if err != nil {
		t.Error(err)
		return
	}
	balance2, err := tp.BalanceOf(nil, addr2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("old balance1=%s,balance2=%s", balance1, balance2)

	tx, err = ch.Settle(bind.NewKeyedTransactor(key1))
	if err != nil {
		t.Error(err)
		return
	}
	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		t.Error(err)
		return
	}
	if receipt.Status != 1 {
		t.Errorf("receipt err %s", receipt)
		return
	}
	time.Sleep(time.Second * 2)
	balance1, err = tp.BalanceOf(nil, addr1)
	if err != nil {
		t.Error(err)
		return
	}
	balance2, err = tp.BalanceOf(nil, addr2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("new balance1=%s,balance2=%s", balance1, balance2)
}

//make sure PartnerBalanceProof is not nil
func signFor3rd(Nonce int64, transferAmount *big.Int, LocksRoot common.Hash, channelAddress common.Address, MessageHash common.Hash, closingSig []byte, thirdAddr common.Address, privkey *ecdsa.PrivateKey) (sig []byte, err error) {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, Nonce)
	buf.Write(utils.BigIntTo32Bytes(transferAmount))
	buf.Write(LocksRoot[:])
	buf.Write(channelAddress[:])
	buf.Write(MessageHash[:])
	buf.Write(closingSig)
	buf.Write(thirdAddr[:])
	dataToSign := buf.Bytes()
	//fmt.Printf("datatosigh=%d, \n%s", len(dataToSign), hex.Dump(dataToSign))
	return utils.SignData(privkey, dataToSign)
}
