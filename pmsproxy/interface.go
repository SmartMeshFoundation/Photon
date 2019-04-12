package pmsproxy

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"math/big"

	"github.com/SmartMeshFoundation/Photon/channel/channeltype"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/rerr"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// DelegateUpdateTransfer :
type DelegateUpdateTransfer struct {
	Nonce               uint64      `json:"nonce"`
	TransferAmount      *big.Int    `json:"transfer_amount"`
	Locksroot           common.Hash `json:"locksroot"`
	ExtraHash           common.Hash `json:"extra_hash"`
	ClosingSignature    []byte      `json:"closing_signature"`
	NonClosingSignature []byte      `json:"non_closing_signature"`
}

// DelegateUnlock : 第三方服务也负责链上unlock
type DelegateUnlock struct {
	Lock        *mtree.Lock `json:"lock"`
	MerkleProof []byte      `json:"merkle_proof"`
	Signature   []byte      `json:"signature"`
}

//DelegatePunish 需要委托给第三方的 punish证据
// DelegatePunish Punish proof that is delegated to third-party.
type DelegatePunish struct {
	LockHash       common.Hash `json:"lock_hash"` //the whole lock's hash,not lock secret hash
	AdditionalHash common.Hash `json:"additional_hash"`
	Signature      []byte      `json:"signature"`
}

//DelegateForPms is for 3rd party to call update transfer
type DelegateForPms struct {
	ChannelIdentifier common.Hash            `json:"channel_identifier"`
	OpenBlockNumber   int64                  `json:"open_block_number"`
	TokenAddrss       common.Address         `json:"token_address"`
	PartnerAddress    common.Address         `json:"partner_address"`
	UpdateTransfer    DelegateUpdateTransfer `json:"update_transfer"`
	Unlocks           []*DelegateUnlock      `json:"unlocks"`
	Punishes          []*DelegatePunish      `json:"punishes"`
}

//SignBalanceProofFor3rd make sure PartnerBalanceProof is not nil
func SignBalanceProofFor3rd(c *channeltype.Serialization, privkey *ecdsa.PrivateKey) (sig []byte, err error) {
	if c.PartnerBalanceProof == nil {
		log.Error(fmt.Sprintf("PartnerBalanceProof is nil,must ber a error"))
		return nil, rerr.ErrChannelBalanceProofNil.Append("empty PartnerBalanceProof")
	}
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractBalanceProofDelegateMessageLength))
	_, err = buf.Write(utils.BigIntTo32Bytes(c.PartnerBalanceProof.TransferAmount))
	_, err = buf.Write(c.PartnerBalanceProof.LocksRoot[:])
	err = binary.Write(buf, binary.BigEndian, c.PartnerBalanceProof.Nonce)
	_, err = buf.Write(c.ChannelIdentifier.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, c.ChannelIdentifier.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.ChainID))
	if err != nil {
		log.Error(fmt.Sprintf("buf write error %s", err))
	}
	dataToSign := buf.Bytes()
	return utils.SignData(privkey, dataToSign)
}

// SignUnlockFor3rd :
func SignUnlockFor3rd(c *channeltype.Serialization, u *DelegateUnlock, thirdAddress common.Address, privkey *ecdsa.PrivateKey) (sig []byte, err error) {
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractUnlockDelegateProofMessageLength))
	_, err = buf.Write(thirdAddress[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(big.NewInt(u.Lock.Expiration)))
	_, err = buf.Write(utils.BigIntTo32Bytes(u.Lock.Amount))
	_, err = buf.Write(u.Lock.LockSecretHash[:])
	_, err = buf.Write(c.ChannelIdentifier.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, c.ChannelIdentifier.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.ChainID))
	if err != nil {
		log.Error(fmt.Sprintf("buf write error %s", err))
		return
	}
	dataToSign := buf.Bytes()
	return utils.SignData(privkey, dataToSign)
}

/*
PmsProxy :
api to call pfg server
*/
type PmsProxy interface {
	/*
		submit partner's balance proof to pfg
	*/
	SubmitDelegate(data *DelegateForPms) error
}
