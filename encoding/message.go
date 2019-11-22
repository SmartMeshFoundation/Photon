package encoding

import (
	"bytes"
	"encoding/binary"

	"crypto/ecdsa"

	"math/big"

	"errors"
	"fmt"

	"encoding/gob"

	"encoding/hex"

	"encoding/json"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/network/rpc/contracts"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/transfer/mtree"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// MessageVersionControlMap 保存每个消息支持的最低版本号
var MessageVersionControlMap = map[int16]int16{
	MediatedTransferCmdID: int16(1), // 2019-03 MediatedTransfer消息升级,带上了Path,不兼容verison<1的版本
}

//MessageType is the type of message for receive and send
type MessageType int

//消息设计原则: 只要引起balanceproof 变化,那么 nonce 就应该加1
// MessageDesign : as long as balanceproof changes, nonce should plus 1.
const (
	//AckCmdID id of Ack message
	AckCmdID = 0
	//PingCmdID id of ping message
	PingCmdID = iota*2 + 1
	//SecretRequestCmdID id of SecretRequest message
	SecretRequestCmdID
	//UnlockCmdID id of Secret message
	UnlockCmdID
	//DirectTransferCmdID id of DirectTransfer, it's now deprecated
	DirectTransferCmdID
	//MediatedTransferCmdID id of MediatedTransfer
	MediatedTransferCmdID
	//AnnounceDisposedTransferCmdID id of AnnounceDisposed message
	AnnounceDisposedTransferCmdID
	//RevealSecretCmdID id of RevealSecret message
	RevealSecretCmdID
	//RemoveExpiredLockCmdID id of RemoveExpiredHashlock message
	RemoveExpiredLockCmdID
	/*
		发起合作关闭通道请求
	*/
	// send CooperativeRequest
	SettleRequestCmdID
	/*
		合作关闭通道响应.
	*/
	SettleResponseCmdID
	/*
		发起提现请求
	*/
	// Send Withdraw Request
	WithdrawRequestCmdID
	/*
		提现响应
	*/
	// Respond Withdraw Request
	WithdrawResponseCmdID
	/*
		refund 响应,
	*/
	// Respond Refund
	AnnounceDisposedTransferResponseCmdID
	/*
		当消息接收方收到消息以后,可能会发现有问题呢,但是没有途径通知接收方除了问题,
		因此增加途径来做消息通知,因为消息会重复发送,如果没有处理妥当会造成错误消息也会反复发送.
		因此针对错误消息,我的想法是保存一个lru进行管理,错误通知多了应该也没什么严重的问题, 只是用户体验不好而已.
	*/
	ErrorNotifyCmdID
)

const signatureLength = 65

var errPacketLength = errors.New("packet length error")

//MessagePacker serialize of a message
type MessagePacker interface {
	//pack message to byte array
	Pack() []byte
}

//MessageUnpacker deserialize of message
type MessageUnpacker interface {
	//unpack message from byte array
	UnPack(data []byte) error
}

//MessagePackerUnpacker is packer and unpacker
type MessagePackerUnpacker interface {
	MessagePacker
	MessageUnpacker
}

//Messager interface for all  message type
type Messager interface {
	//Cmd id of message
	Cmd() int
	//Tag is used for save and restore
	Tag() interface{}
	//SetTag set tage
	SetTag(tag interface{})
	//Name of this message
	Name() string
	//String fmt.Stringer
	String() string
	MessagePackerUnpacker
}

//CmdStruct base of message
type CmdStruct struct {
	CmdID       int16
	Version     int16
	InternalTag interface{} //for save to database
}

//Cmd id of this message
func (cmd *CmdStruct) Cmd() int {
	return int(cmd.CmdID)
}

//Tag for internal state save
func (cmd *CmdStruct) Tag() interface{} {
	return cmd.InternalTag
}

//SetTag for internal state save
func (cmd *CmdStruct) SetTag(tag interface{}) {
	cmd.InternalTag = tag
}

//Name of this message
func (cmd *CmdStruct) Name() string {
	return MessageType(cmd.CmdID).String()
}

// WriteCmdStructToBuf write cmdID and version into buf
func (cmd *CmdStruct) WriteCmdStructToBuf(buf *bytes.Buffer) (err error) {
	err = binary.Write(buf, binary.LittleEndian, cmd.CmdID)
	err = binary.Write(buf, binary.LittleEndian, cmd.Version)
	return
}

// ReadCmdStructFromBuf read CmdStruct from buf
func (cmd *CmdStruct) ReadCmdStructFromBuf(buf *bytes.Buffer) (err error) {
	var t int16
	err = binary.Read(buf, binary.LittleEndian, &t)
	if err != nil {
		return
	}
	cmd.CmdID = t
	err = binary.Read(buf, binary.LittleEndian, &t)
	if err != nil {
		return
	}
	cmd.Version = t
	return
}

//SignedMessager interface of message that needs signed
type SignedMessager interface {
	Messager
	GetSender() common.Address
	Sign(priveKey *ecdsa.PrivateKey, pack MessagePacker) error
	verifySignature(data []byte) error
}

//EnvelopMessager is message contains new balance proof
type EnvelopMessager interface {
	SignedMessager
	//GetEnvelopMessage returns EnvelopMessage
	GetEnvelopMessage() *EnvelopMessage
}

// String return the string representation of message type.
func (t MessageType) String() string {
	switch t {
	case AckCmdID:
		return "Ack"
	case PingCmdID:
		return "Ping"
	case SecretRequestCmdID:
		return "SecretRequest"
	case UnlockCmdID:
		return "Unlock"
	case DirectTransferCmdID:
		return "DirectTransfer"
	case MediatedTransferCmdID:
		return "MediatedTransfer"
	case AnnounceDisposedTransferCmdID:
		return "AnnounceDisposed"
	case AnnounceDisposedTransferResponseCmdID:
		return "AnnounceDisposedResponse"
	case RevealSecretCmdID:
		return "RevealSecret"
	case RemoveExpiredLockCmdID:
		return "RemoveExpiredHashlock"
	case SettleRequestCmdID:
		return "SettleRequest"
	case SettleResponseCmdID:
		return "SettleResponse"
	case WithdrawRequestCmdID:
		return "WithdrawRequest"
	case WithdrawResponseCmdID:
		return "WithdrawResponse"
	case ErrorNotifyCmdID:
		return "ErrorNotify"
	default:
		return "<unknown>"
	}
}

/*
Ack All accepted messages should be confirmed by an `Ack` which echoes the
orginals Message hash.

We don'T Sign Acks because attack vector can be mitigated and to speed up
things.
*/
type Ack struct {
	CmdStruct
	Sender common.Address
	Echo   common.Hash
}

//NewAck create ack message
func NewAck(sender common.Address, echo common.Hash) *Ack {
	return &Ack{
		CmdStruct: CmdStruct{CmdID: AckCmdID},
		Sender:    sender,
		Echo:      echo,
	}
}

//Pack implements of MessagePacker
func (ack *Ack) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = ack.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, ack.CmdID)
	_, err = buf.Write(ack.Sender[:])
	_, err = buf.Write(ack.Echo[:])
	if err != nil {
		panic(fmt.Sprintf("Ack Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is implements of MessageUnpacker
func (ack *Ack) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = ack.ReadCmdStructFromBuf(buf)
	if AckCmdID != ack.CmdID {
		panic(fmt.Sprintf("Ack Unpack cmdid should be 0,but get %d", ack.CmdID))
	}
	_, err = buf.Read(ack.Sender[:])
	n, err := buf.Read(ack.Echo[:])
	if err != nil {
		return err
	}
	if n != len(ack.Echo) {
		return errPacketLength
	}
	return nil
}
func (ack *Ack) String() string {
	return fmt.Sprintf("Message{type=Ack sender=%s,echo=%s}", utils.APex2(ack.Sender), utils.HPex(ack.Echo))
}

//SignedMessage is corresponding of SignedMessager
type SignedMessage struct {
	CmdStruct
	Sender    common.Address
	Signature []byte
}

//GetSender returns the sender of this message
func (m *SignedMessage) GetSender() common.Address {
	return m.Sender
}

//Sign this message
func (m *SignedMessage) Sign(priveKey *ecdsa.PrivateKey, pack MessagePacker) error {
	if len(m.Signature) > 0 {
		log.Warn("duplicate Sign")
		return errors.New("duplicate Sign")
	}
	m.Signature = SignMessage(priveKey, pack)
	m.Sender = crypto.PubkeyToAddress(priveKey.PublicKey)
	return nil
}

//verifySignature returns error if is not a valid signature
func (m *SignedMessage) verifySignature(data []byte) error {
	sender, err := VerifyMessage(data)
	if err != nil {
		return err
	}
	m.Sender = sender
	return nil
}

//SignMessage signs a message
func SignMessage(privKey *ecdsa.PrivateKey, pack MessagePacker) []byte {
	data := pack.Pack()
	sig, err := utils.SignData(privKey, data)
	if err != nil {
		panic(fmt.Sprintf("SignMessage error %s", err))
	}
	return sig
}

//HashMessageWithoutSignature returns the raw hash of this message
func HashMessageWithoutSignature(pack MessagePacker) common.Hash {
	data := pack.Pack()
	if len(data) > signatureLength {
		data = data[:len(data)-signatureLength]
	}
	return utils.Sha3(data)
}

//VerifyMessage returns the sender of message if data is a valid SignedMessage
func VerifyMessage(data []byte) (sender common.Address, err error) {
	messageData := data[:len(data)-signatureLength]
	signature := make([]byte, signatureLength)
	copy(signature, data[len(data)-signatureLength:])
	hash := utils.Sha3(messageData)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return
	}
	sender = utils.PubkeyToAddress(pubkey)
	return
}

//Ping message
type Ping struct {
	SignedMessage
	Nonce int64
}

//NewPing create ping message
func NewPing(nonce int64) *Ping {
	p := &Ping{
		//SignedMessage:SignedMessage{CmdStruct: CmdStruct{CmdID: PingCmdID}},
		Nonce: nonce,
	}
	p.CmdID = PingCmdID
	return p
}

//Pack is MessagePacker
func (p *Ping) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = p.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, p.CmdID) //only one byte
	err = binary.Write(buf, binary.BigEndian, p.Nonce)
	_, err = buf.Write(p.Signature)
	if err != nil {
		panic(fmt.Sprintf("Ping Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (p *Ping) UnPack(data []byte) error {
	var err error
	if len(data) != 77 { //stun response here
		return errPacketLength
	}

	buf := bytes.NewBuffer(data)
	err = p.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if PingCmdID != p.CmdID {
		return fmt.Errorf("Ping Unpack cmdid should be  1,but get %d", p.CmdID)
	}
	err = binary.Read(buf, binary.BigEndian, &p.Nonce)
	p.Signature = make([]byte, signatureLength)
	_, err = buf.Read(p.Signature)
	err = p.SignedMessage.verifySignature(data)
	if err != nil {
		return err
	}
	return nil
}

//String is fmt.Stringer
func (p *Ping) String() string {
	return fmt.Sprintf("Message{type=Ping nonce=%d,sender=%s, has signature=%v}", p.Nonce, utils.APex2(p.Sender), len(p.Signature) != 0)
}

//ErrorNotifyType 来自对方的错误通知,类型
type ErrorNotifyType int16

const (
	//消息最长是1200,65-签名,2-cmdid,2-version,2-ErrorNotifyType,2-DataLength
	errorNotifyMaxRelatedDataLength = 1200 - 65 - 4 - 4

	//InvalidNonceErrorNotify 接收方收到了带有BalanceProof的消息,但是因为数据库不一致,导致Nonce错误
	InvalidNonceErrorNotify = iota
)

//ErrorNotify 发消息通知对方发生了错误
type ErrorNotify struct {
	SignedMessage
	ErrorNotifyType ErrorNotifyType
	RelatedData     []byte
}

//NewErrorNotify 错误通知
func NewErrorNotify(notifyType ErrorNotifyType, errorData []byte) *ErrorNotify {
	p := &ErrorNotify{
		ErrorNotifyType: notifyType,
		RelatedData:     errorData,
	}
	p.CmdID = ErrorNotifyCmdID
	return p
}

//Pack is MessagePacker
func (en *ErrorNotify) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = en.WriteCmdStructToBuf(buf)
	err = binary.Write(buf, binary.BigEndian, en.ErrorNotifyType)
	var rl = uint16(len(en.RelatedData))
	err = binary.Write(buf, binary.BigEndian, rl)
	if len(en.RelatedData) > errorNotifyMaxRelatedDataLength {
		panic("relateddata length error")
	}
	_, err = buf.Write(en.RelatedData)
	_, err = buf.Write(en.Signature)
	if err != nil {
		panic(fmt.Sprintf("ErrorNotify pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnpacker,注意可能包含有附加消息也可能没有包含附加消息
func (en *ErrorNotify) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = en.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if ErrorNotifyCmdID != en.CmdID {
		return fmt.Errorf("ErrorNotify Unpack cmdid should be  %d,but get %d", ErrorNotifyCmdID, en.CmdID)
	}
	err = binary.Read(buf, binary.BigEndian, &en.ErrorNotifyType)
	if err != nil {
		return err
	}
	var relatedDataLen uint16
	err = binary.Read(buf, binary.BigEndian, &relatedDataLen)
	if err != nil {
		return err
	}
	if relatedDataLen > errorNotifyMaxRelatedDataLength {
		return fmt.Errorf("relatedDataLen is too large,max=%d,got=%d", errorNotifyMaxRelatedDataLength, relatedDataLen)
	}
	en.RelatedData = make([]byte, relatedDataLen)
	_, err = buf.Read(en.RelatedData)
	l := buf.Len()
	if l != signatureLength {
		return fmt.Errorf("ErrorNotify ,leftLen=%d, not signature", l)
	}
	en.Signature = make([]byte, signatureLength)
	_, err = buf.Read(en.Signature)
	err = en.verifySignature(data)
	return err
}

//String is fmt.Stringer
func (en *ErrorNotify) String() string {
	return fmt.Sprintf("Message{type=ErrorNotify ErrorNotifyType=%d,errorDataLen=%d,sender=%s,has signature=%v}",
		en.ErrorNotifyType, len(en.RelatedData), utils.APex2(en.Sender), len(en.Signature) != 0)
}

//SecretRequest Requests the secret which unlocks a hashlock.
type SecretRequest struct {
	SignedMessage
	LockSecretHash common.Hash
	PaymentAmount  *big.Int
}

//NewSecretRequest create SecretRequest
func NewSecretRequest(lockSecretHash common.Hash, paymentAmount *big.Int) *SecretRequest {
	p := &SecretRequest{
		LockSecretHash: lockSecretHash,
		PaymentAmount:  new(big.Int).Set(paymentAmount),
	}
	p.CmdID = SecretRequestCmdID
	return p
}

//Pack is MessagePacker
func (sr *SecretRequest) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = sr.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, sr.CmdID) //only one byte..
	_, err = buf.Write(sr.LockSecretHash[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(sr.PaymentAmount))
	_, err = buf.Write(sr.Signature)
	if err != nil {
		panic(fmt.Sprintf("SecretRequest Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnpacker
func (sr *SecretRequest) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = sr.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if SecretRequestCmdID != sr.CmdID {
		return fmt.Errorf("SecretRequest Unpack cmdid should be  3,but get %d", sr.CmdID)
	}
	_, err = buf.Read(sr.LockSecretHash[:])
	sr.PaymentAmount = utils.ReadBigInt(buf)
	sr.Signature = make([]byte, signatureLength)
	n, err := buf.Read(sr.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errPacketLength
	}
	err = sr.verifySignature(data)
	return err
}

//String is fmt.Stringer
func (sr *SecretRequest) String() string {
	return fmt.Sprintf("Message{type=SecretRequest LockSecretHash=%s,paymentAmount=%s,sender=%s,has signature=%v}",
		utils.HPex(sr.LockSecretHash), sr.PaymentAmount.String(), utils.APex2(sr.Sender), len(sr.Signature) != 0)
}

/*
RevealSecret used to reveal a secret to party known to have interest in it.

This message is not sufficient for state changes in the Photon Channel, the
reason is that a node participating in split transfer or in both mediated
transfer for an exchange might can reveal the secret to it's partners, but
that must not update the internal channel state.
*/
type RevealSecret struct {
	SignedMessage
	LockSecret     common.Hash
	lockSecretHash common.Hash
	Data           []byte // used to transfer custom message, length should < 256
}

//NewRevealSecret create RevealSecret
func NewRevealSecret(lockSecret common.Hash) *RevealSecret {
	p := &RevealSecret{
		LockSecret: lockSecret,
	}
	p.CmdID = RevealSecretCmdID
	return p
}

//CloneRevealSecret clones a RevealSecret Message
func CloneRevealSecret(rs *RevealSecret) *RevealSecret {
	rs2 := *rs
	return &rs2
}

//LockSecretHash return hash of secret
func (rs *RevealSecret) LockSecretHash() common.Hash {
	if rs.lockSecretHash == utils.EmptyHash {
		rs.lockSecretHash = utils.ShaSecret(rs.LockSecret[:])
	}
	return rs.lockSecretHash
}

//Pack is MessagePacker
func (rs *RevealSecret) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	//err = binary.Write(buf, binary.LittleEndian, rs.CmdID) //only one byte.
	err = rs.WriteCmdStructToBuf(buf)
	_, err = buf.Write(rs.LockSecret[:])
	// write data
	dataLen := uint64(len(rs.Data))
	err = utils.WriteVarInt(buf, dataLen)
	if dataLen > 0 {
		_, err = buf.Write(rs.Data)
	}
	_, err = buf.Write(rs.Signature)
	if err != nil {
		panic(fmt.Sprintf("RevealSecret Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (rs *RevealSecret) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = rs.ReadCmdStructFromBuf(buf)
	if RevealSecretCmdID != rs.CmdID {
		return fmt.Errorf("RevealSecret Unpack cmdid should be  11,but get %d", rs.CmdID)
	}
	_, err = buf.Read(rs.LockSecret[:])
	// readData
	dataLen, err := utils.ReadVarInt(buf)
	if err != nil {
		return err
	}
	if dataLen > uint64(params.Cfg.UDPMaxMessageSize) {
		return fmt.Errorf("RevealSecret unpack data error, too large data, maby attack")
	}
	if dataLen > 0 {
		rs.Data = make([]byte, dataLen)
		err = binary.Read(buf, binary.LittleEndian, &rs.Data)
		if err != nil {
			return errors.New("RevealSecret unpack data error")
		}
	}
	rs.Signature = make([]byte, signatureLength)
	n, err := buf.Read(rs.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errPacketLength
	}
	return rs.verifySignature(data)
}

//String fmt.Stringer
func (rs *RevealSecret) String() string {
	return fmt.Sprintf("Message{type=RevealSecret,hashlock=%s,secret=%s,sender=%s,has signature=%v}", utils.HPex(rs.LockSecretHash()),
		utils.HPex(rs.LockSecret), utils.APex2(rs.Sender), len(rs.Signature) != 0)
}

//BalanceProof in the message ,not the same as data need by the contract
type BalanceProof struct {
	Nonce             uint64
	ChannelIdentifier common.Hash
	OpenBlockNumber   int64    //open blocknumber 和 channelIdentifier 一起作为通道的唯一标识	// the only tag for a channel = OpenBlockNumber + ChannelIdentifier
	TransferAmount    *big.Int //The number has been transferred to the other party
	Locksroot         common.Hash
}

//NewBalanceProof create a balance proof
func NewBalanceProof(nonce uint64, transferredAmount *big.Int, locksRoot common.Hash, channelID *contracts.ChannelUniqueID) *BalanceProof {
	return &BalanceProof{
		Nonce:             nonce,
		TransferAmount:    transferredAmount,
		Locksroot:         locksRoot,
		ChannelIdentifier: channelID.ChannelIdentifier,
		OpenBlockNumber:   channelID.OpenBlockNumber,
	}
}

//EnvelopMessage is general part of message that contains a new balanceproof
type EnvelopMessage struct {
	SignedMessage
	BalanceProof
}

//String is fmt.Stringer
func (m *EnvelopMessage) String() string {
	return fmt.Sprintf("EnvelopMessage{nonce=%d,Channel=%s,openBlockNumber=%d,TransferAmount=%s,Locksroot=%s, sender=%s,has signature=%v}", m.Nonce,
		utils.HPex(m.ChannelIdentifier), m.OpenBlockNumber, m.TransferAmount, utils.HPex(m.Locksroot), utils.APex2(m.Sender), len(m.Signature) != 0)
}
func (m *EnvelopMessage) signData(datahash common.Hash) []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractBalanceProofMessageLength))
	_, err = buf.Write(utils.BigIntTo32Bytes(m.TransferAmount))
	_, err = buf.Write(m.Locksroot[:])
	err = binary.Write(buf, binary.BigEndian, m.Nonce)
	_, err = buf.Write(datahash[:])
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.Cfg.ChainID))
	if err != nil {
		log.Error(fmt.Sprintf("signData err %s", err))
	}
	dataToSign := buf.Bytes()
	return dataToSign
}

/*
Sign data=(once+transferamount+locksroot+channel+hash(data))
*/
func (m *EnvelopMessage) Sign(privKey *ecdsa.PrivateKey, msg MessagePacker) error {
	data := msg.Pack() //before signed, Sign twice will be error
	datahash := utils.Sha3(data)
	//compute data to Sign
	dataToSign := m.signData(datahash)
	sig, err := utils.SignData(privKey, dataToSign)
	if err != nil {
		return err
	}
	m.Signature = sig
	m.Sender = crypto.PubkeyToAddress(privKey.PublicKey)
	return nil
}

//verifySignature returns error if is not a valid signature
func (m *EnvelopMessage) verifySignature(data []byte) error {
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	datatosign := m.signData(datahash)
	//should not change data's content,because its name is verify.
	var signature = make([]byte, signatureLength)
	copy(signature, data[len(data)-signatureLength:])
	hash := utils.Sha3(datatosign)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return err
	}
	m.Sender = utils.PubkeyToAddress(pubkey)
	return nil

}
func (m *EnvelopMessage) checkValid() error {
	if m.Nonce <= 0 {
		return fmt.Errorf("nonce must be positive %d", m.Nonce)
	}
	if !utils.IsValidUint256(m.TransferAmount) {
		return fmt.Errorf("transfer amount must not be negative %s", m.TransferAmount)
	}
	if m.OpenBlockNumber <= 0 || m.ChannelIdentifier == utils.EmptyHash {
		return fmt.Errorf("channel id error, openBlockNumber=%d,channelIdentifier=%s", m.OpenBlockNumber, utils.HPex(m.ChannelIdentifier))
	}
	return nil
}

//GetEnvelopMessage return EnvelopMessage
func (m *EnvelopMessage) GetEnvelopMessage() *EnvelopMessage {
	return m
}

func (m *EnvelopMessage) pack(buf *bytes.Buffer) {
	var err error
	err = binary.Write(buf, binary.BigEndian, m.Nonce)
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(m.TransferAmount))
	_, err = buf.Write(m.Locksroot[:])
	_, err = buf.Write(m.Signature)
	if err != nil {
		log.Error(fmt.Sprintf("EnvelopMessage pack err %s", err))
	}
}
func (m *EnvelopMessage) unpack(buf *bytes.Buffer) error {
	var err error
	err = binary.Read(buf, binary.BigEndian, &m.Nonce)
	_, err = buf.Read(m.ChannelIdentifier[:])
	err = binary.Read(buf, binary.BigEndian, &m.OpenBlockNumber)
	m.TransferAmount = utils.ReadBigInt(buf)
	_, err = buf.Read(m.Locksroot[:])
	m.Signature = make([]byte, signatureLength)
	n, err := buf.Read(m.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errors.New("packet length error")
	}
	if err := m.checkValid(); err != nil {
		return err
	}
	return nil
}

func (m *EnvelopMessage) fromBalanceProof(bp *BalanceProof) {
	m.Nonce = bp.Nonce
	m.ChannelIdentifier = bp.ChannelIdentifier
	m.OpenBlockNumber = bp.OpenBlockNumber
	m.TransferAmount = new(big.Int).Set(bp.TransferAmount)
	m.Locksroot = bp.Locksroot
}

/*
UnLock Message used to do state changes on a partner Photon Channel.

Locksroot changes need to be synchronized among both participants, the
protocol is for only the side unlocking to send the Secret message allowing
the other party to withdraw.
*/
type UnLock struct {
	EnvelopMessage
	LockSecret common.Hash
}

//LockSecretHash is Hash of secret
func (s *UnLock) LockSecretHash() common.Hash {
	return utils.ShaSecret(s.LockSecret[:])
}

//NewUnlock create Secret message
func NewUnlock(bp *BalanceProof, lockSecret common.Hash) *UnLock {
	p := &UnLock{
		LockSecret: lockSecret,
	}
	p.CmdID = UnlockCmdID
	p.EnvelopMessage.fromBalanceProof(bp)
	return p
}

//Pack is MessagePacker
func (s *UnLock) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = s.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, s.CmdID) //only one byte.
	_, err = buf.Write(s.LockSecret[:])
	s.EnvelopMessage.pack(buf)
	if err != nil {
		panic(fmt.Sprintf("UnLock Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (s *UnLock) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = s.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if UnlockCmdID != s.CmdID {
		return fmt.Errorf("Ack Secret cmdid should be  4,but get %d", s.CmdID)
	}
	_, err = buf.Read(s.LockSecret[:])
	err = s.EnvelopMessage.unpack(buf)
	if err != nil {
		return err
	}
	return s.EnvelopMessage.verifySignature(data)
}

//String is fmt.Stringer
func (s *UnLock) String() string {
	return fmt.Sprintf("Message{type=Unlock secret=%s,%s}", utils.HPex(s.LockSecret), s.EnvelopMessage.String())
}

/*
RemoveExpiredHashlockTransfer message from sender to receiver, notify to remove a expired hashlock, provide new blance proof.

Removes one lock that has expired. Used to trim the merkle tree and recover the locked capacity. This message is only valid if the corresponding lock expiration is lower than the latest block number for the corresponding blockchain.
Fields
Field Name 	     Field Type 	      Description
secrethash 	    bytes32 	          The secrethash to remove
balance_proof 	BalanceProof 	      The updated balance proof
signature 	    bytes 	              Elliptic Curve 256k1 signature
*/
type RemoveExpiredHashlockTransfer struct {
	EnvelopMessage
	LockSecretHash common.Hash
}

//NewRemoveExpiredHashlockTransfer create  RemoveExpiredHashlockTransfer
func NewRemoveExpiredHashlockTransfer(bp *BalanceProof, lockSecretHash common.Hash) *RemoveExpiredHashlockTransfer {
	p := &RemoveExpiredHashlockTransfer{
		LockSecretHash: lockSecretHash,
	}
	p.CmdID = RemoveExpiredLockCmdID
	p.EnvelopMessage.fromBalanceProof(bp)
	return p
}

//Pack is MessagePacker
func (m *RemoveExpiredHashlockTransfer) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID) //only one byte.
	_, err = buf.Write(m.LockSecretHash[:])
	m.EnvelopMessage.pack(buf)
	if err != nil {
		panic(fmt.Sprintf("RemoveExpiredHashlockTransfer Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *RemoveExpiredHashlockTransfer) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if RemoveExpiredLockCmdID != m.CmdID {
		return fmt.Errorf("Ack Secret cmdid should be  4,but get %d", m.CmdID)
	}
	_, err = buf.Read(m.LockSecretHash[:])
	err = m.EnvelopMessage.unpack(buf)
	if err != nil {
		return err
	}
	return m.EnvelopMessage.verifySignature(data)
}

//String is fmt.Stringer
func (m *RemoveExpiredHashlockTransfer) String() string {
	return fmt.Sprintf("Message{type=RemoveExpiredHashlockTransfer LockSecretHash=%s,%s}", utils.HPex(m.LockSecretHash), m.EnvelopMessage.String())
}

/*
DirectTransfer is a direct token exchange, used when both participants have a previously
opened channel.

Signs the unidirectional settled `balance` of `token` to `recipient` plus
locked transfers.

Settled refers to the inclusion of formerly locked amounts.
Locked amounts are not included in the balance yet, but represented
by the `locksroot`.

Args:
    nonce: A sequential nonce, used to protected against replay attacks and
        to give a total order for the messages. This nonce is per
        participant, not shared.
    token: The address of the token being exchanged in the channel.
    transferred_amount: The total amount of token that was transferred to
        the channel partner. This value is monotonically increasing and can
        be larger than a channels deposit, since the channels are
        bidirecional.
    recipient: The address of the Photon node participating in the channel.
    locksroot: The root of a merkle tree which records the current
        outstanding locks.
*/
type DirectTransfer struct {
	EnvelopMessage
	Data               []byte      // used to transfer custom message, length should < 256
	FakeLockSecretHash common.Hash // used when save transfer status to db, do not be used when message pack/unpack
}

//String is fmt.Stringer
func (m *DirectTransfer) String() string {
	return fmt.Sprintf("Message{type=DirectTransfer %s}", m.EnvelopMessage.String())
}

//NewDirectTransfer create DirectTransfer
func NewDirectTransfer(bp *BalanceProof) *DirectTransfer {
	p := &DirectTransfer{}
	p.CmdID = DirectTransferCmdID
	p.EnvelopMessage.fromBalanceProof(bp)
	return p
}

//Pack is MessagePacker
func (m *DirectTransfer) Pack() []byte {
	buf := new(bytes.Buffer)
	err := m.WriteCmdStructToBuf(buf)
	//err := binary.Write(buf, binary.LittleEndian, m.CmdID) //only one byte.
	if err != nil {
		panic(fmt.Sprintf("DirectTransfer Pack err %s", err))
	}
	// write data
	dataLen := uint64(len(m.Data))
	err = utils.WriteVarInt(buf, dataLen)
	if dataLen > 0 {
		_, err = buf.Write(m.Data)
	}
	if err != nil {
		panic(fmt.Sprintf("DirectTransfer Pack err %s", err))
	}
	m.EnvelopMessage.pack(buf)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *DirectTransfer) UnPack(data []byte) error {
	buf := bytes.NewBuffer(data)
	err := m.ReadCmdStructFromBuf(buf)
	//err := binary.Read(buf, binary.LittleEndian, &t)
	if err != nil {
		return err
	}
	if DirectTransferCmdID != m.CmdID {
		return errors.New("DirectTransfer unpack cmdid error")
	}
	// readData
	dataLen, err := utils.ReadVarInt(buf)
	if err != nil {
		return err
	}
	if dataLen > uint64(params.Cfg.UDPMaxMessageSize) {
		return fmt.Errorf("RevealSecret unpack data error, too large data, maby attack")
	}
	if dataLen > 0 {
		m.Data = make([]byte, dataLen)
		err = binary.Read(buf, binary.LittleEndian, &m.Data)
		if err != nil {
			return errors.New("DirectTransfer unpack data error")
		}
	}
	err = m.EnvelopMessage.unpack(buf)
	if err != nil {
		return err
	}
	return m.EnvelopMessage.verifySignature(data)
}

/*
MediatedTransfer has a `target` address to which a chain of transfers shall
be established. Here the `haslock` is mandatory.

`fee` is the remaining fee a recipient shall use to complete the mediated transfer.
The recipient can deduct his own fee from the amount and lower `fee` to the remaining fee.
Just as the recipient can fail to forward at all, or the assumed amount,
it can deduct a too high fee, but this would render completion of the transfer unlikely.

The initiator of a mediated transfer will calculate fees based on the likely fees along the
path. Note, it can not determine the path, as it does not know which nodes are available.

Initial `amount` should be expected received amount + fees.

Fees are always payable by the initiator.

`initiator` is the party that knows the secret to the `hashlock`
*/
type MediatedTransfer struct {
	EnvelopMessage
	Expiration     int64
	LockSecretHash common.Hash
	PaymentAmount  *big.Int //The number transferred to party
	Target         common.Address
	Initiator      common.Address
	Fee            *big.Int
	Path           []common.Address // 2019-03 消息升级后,带全路径信息
}

//String is fmt.Stringer
func (m *MediatedTransfer) String() string {
	return fmt.Sprintf("Message{type=MediatedTransfer expiration=%d,target=%s,initiator=%s,hashlock=%s,amount=%s,fee=%s,path=%s,%s}",
		m.Expiration, utils.APex2(m.Target), utils.APex2(m.Initiator),
		utils.HPex(m.LockSecretHash), m.PaymentAmount, m.Fee, m.GetPathStr(), m.EnvelopMessage.String())
}

//NewMediatedTransfer create MediatedTransfer
func NewMediatedTransfer(bp *BalanceProof, lock *mtree.Lock,
	target, initiator common.Address, fee *big.Int, path []common.Address) *MediatedTransfer {
	p := &MediatedTransfer{
		Target:         target,
		Initiator:      initiator,
		Fee:            new(big.Int).Set(fee),
		PaymentAmount:  lock.Amount,
		LockSecretHash: lock.LockSecretHash,
		Expiration:     lock.Expiration,
		Path:           path, // 2019-03 消息升级后,带全路径信息
	}
	p.CmdID = MediatedTransferCmdID
	p.Version = MessageVersionControlMap[MediatedTransferCmdID]
	p.EnvelopMessage.fromBalanceProof(bp)
	return p
}

//GetLock returns Lock of this Transfer
func (m *MediatedTransfer) GetLock() *mtree.Lock {
	return &mtree.Lock{
		Expiration:     m.Expiration,
		Amount:         m.PaymentAmount,
		LockSecretHash: m.LockSecretHash,
	}
}

//Pack is MessagePacker
func (m *MediatedTransfer) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID) //one byte
	//HTLC
	err = binary.Write(buf, binary.BigEndian, m.Expiration)
	_, err = buf.Write(m.LockSecretHash[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.PaymentAmount))

	_, err = buf.Write(m.Target[:])
	_, err = buf.Write(m.Initiator[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Fee))
	// 2019-03 消息升级,带全路径信息
	pathLen := int32(len(m.Path))
	err = binary.Write(buf, binary.BigEndian, pathLen)
	for _, addr := range m.Path {
		_, err = buf.Write(addr[:])
	}
	m.EnvelopMessage.pack(buf)
	if err != nil {
		panic(fmt.Sprintf("MediatedTransfer Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *MediatedTransfer) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if m.CmdID != MediatedTransferCmdID {
		return errors.New("MediatedTransfer unpack cmd error")
	}
	if m.Version < MessageVersionControlMap[MediatedTransferCmdID] {
		return fmt.Errorf("MediatedTransfer unpack cmd error, version too low ,expect version %d ,but got %d", MessageVersionControlMap[MediatedTransferCmdID], m.Version)
	}
	//HTLC
	err = binary.Read(buf, binary.BigEndian, &m.Expiration)
	_, err = buf.Read(m.LockSecretHash[:])
	m.PaymentAmount = utils.ReadBigInt(buf)

	_, err = buf.Read(m.Target[:])
	_, err = buf.Read(m.Initiator[:])
	m.Fee = utils.ReadBigInt(buf)
	// 2019-03 消息升级,带全路径信息
	var pathLen int32
	err = binary.Read(buf, binary.BigEndian, &pathLen)
	m.Path = []common.Address{}
	for i := int32(0); i < pathLen; i++ {
		var addr common.Address
		_, err = buf.Read(addr[:])
		m.Path = append(m.Path, addr)
	}
	err = m.EnvelopMessage.unpack(buf)
	if err != nil {
		return err
	}
	return m.verifySignature(data)
}

// GetPathStr get string of path to print
func (m *MediatedTransfer) GetPathStr() string {
	if m.Path == nil {
		return ""
	}
	buf, err := json.Marshal(m.Path)
	if err != nil {
		log.Error(fmt.Sprintf("MediatedTransfer.GetPathStr() err : %s", err.Error()))
		return ""
	}
	return string(buf)
}

//GetMtrFromLockedTransfer returns the MediatedTransfer ,the caller must maker sure this message is a  locked transfer
func GetMtrFromLockedTransfer(tr Messager) (mtr *MediatedTransfer) {
	if !(tr.Cmd() == MediatedTransferCmdID) {
		panic("getmtr should never panic")
	}
	mtr, ok := tr.(*MediatedTransfer)
	if !ok {
		panic(fmt.Sprintf("GetMtrFromLockedTransfer from message=%s", utils.StringInterface(tr, 3)))
	}
	return
}

//ChannelIDInMessage common part of message that don't have a balance proof
type ChannelIDInMessage struct {
	ChannelIdentifier common.Hash
	OpenBlockNumber   int64
}

/*
AnnounceDisposedProof is proof used for contracts
*/
type AnnounceDisposedProof struct {
	Lock *mtree.Lock
	ChannelIDInMessage
}

//AnnounceDisposed is used for a mediated node who cannot find any next node to send the mediatedTransfer
type AnnounceDisposed struct {
	SignedMessage
	AnnounceDisposedProof
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
}

//String is fmt.Stringer
func (m *AnnounceDisposed) String() string {
	return fmt.Sprintf("Message{type=AnnounceDisposed Lock=%s,"+
		"ChannelIdentifier=%s-%d ErrorCode=%d ErrorMsg=%s Signature=%s}",
		m.Lock,
		utils.HPex(m.ChannelIdentifier),
		m.OpenBlockNumber,
		m.ErrorCode,
		m.ErrorMsg,
		common.Bytes2Hex(m.Signature),
	)
}

//NewAnnounceDisposed create AnnounceDisposed
func NewAnnounceDisposed(rp *AnnounceDisposedProof, errorCode int, errMsg string) *AnnounceDisposed {
	p := &AnnounceDisposed{
		AnnounceDisposedProof: *rp,
		ErrorCode:             errorCode,
		ErrorMsg:              errMsg,
	}
	p.CmdID = AnnounceDisposedTransferCmdID
	return p
}

//Pack implemnts Messager interface
func (m *AnnounceDisposed) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID)
	_, err = buf.Write(m.Lock.AsBytes())
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	// 2019-03 添加错误码及错误信息
	errCode := int32(m.ErrorCode)
	err = binary.Write(buf, binary.BigEndian, errCode)
	errorMsgBytes := utils.StringToBytes(m.ErrorMsg)
	errorMsgBytesLen := int32(len(errorMsgBytes))
	err = binary.Write(buf, binary.BigEndian, errorMsgBytesLen)
	if errorMsgBytesLen > 0 {
		_, err = buf.Write(errorMsgBytes)
	}
	_, err = buf.Write(m.Signature)
	if err != nil {
		panic(fmt.Sprintf("pack AnnounceDisposed err %s", err))
	}
	return buf.Bytes()
}

//UnPack implements MessageUnPacker
func (m *AnnounceDisposed) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	err = m.ReadCmdStructFromBuf(buf)
	if m.CmdID != AnnounceDisposedTransferCmdID {
		return fmt.Errorf("AnnounceDisposed UnPack cmd error,expect=%d,got=%d", AnnounceDisposedTransferCmdID, m.CmdID)
	}
	m.Lock = new(mtree.Lock)
	err = m.Lock.FromReader(buf)
	_, err = buf.Read(m.ChannelIdentifier[:])
	err = binary.Read(buf, binary.BigEndian, &m.OpenBlockNumber)
	// 2019-03 添加错误码及错误信息
	var errCode int32
	err = binary.Read(buf, binary.BigEndian, &errCode)
	m.ErrorCode = int(errCode)
	var errorMsgBytesLen int32
	err = binary.Read(buf, binary.BigEndian, &errorMsgBytesLen)
	if errorMsgBytesLen > 0 {
		errorMsgBuf := make([]byte, errorMsgBytesLen)
		_, err = buf.Read(errorMsgBuf)
		m.ErrorMsg = utils.BytesToString(errorMsgBuf)
	}
	m.Signature = make([]byte, signatureLength)
	n, err := buf.Read(m.Signature)
	if err != nil || n != signatureLength {
		return fmt.Errorf("AnnounceDisposed UnPack err=%v,read length=%d", err, n)
	}
	return m.verifySignature(data)
}
func (m *AnnounceDisposed) signData(datahash common.Hash) []byte {
	var err error
	buf := new(bytes.Buffer)
	lockhash := m.Lock.Hash()
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractDisposedProofMessageLength))
	_, err = buf.Write(lockhash[:])
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	//// 2019-03 添加错误码及错误信息,这里不能添加,否则需要修改合约
	//err = binary.Write(buf, binary.BigEndian, m.ErrorCode)
	//errorMsgBytes := utils.StringToBytes(m.ErrorMsg)
	//err = binary.Write(buf, binary.BigEndian, len(errorMsgBytes))
	//if len(errorMsgBytes) > 0 {
	//	_, err = buf.Write(errorMsgBytes)
	//}
	_, err = buf.Write(utils.BigIntTo32Bytes(params.Cfg.ChainID))
	_, err = buf.Write(datahash[:])
	if err != nil {
		log.Error(fmt.Sprintf("signData err %s", err))
	}
	dataToSign := buf.Bytes()
	return dataToSign
}

//GetAdditionalHash return hash of this message
func (m *AnnounceDisposed) GetAdditionalHash() common.Hash {
	if m.GetSender() == utils.EmptyAddress {
		panic("should not happen")
	}
	data := m.Pack()
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	return datahash
}

/*
Sign data=(once+transferamount+locksroot+channel+hash(data))
*/
func (m *AnnounceDisposed) Sign(privKey *ecdsa.PrivateKey, msg MessagePacker) error {
	data := msg.Pack() //before signed, Sign twice will be error
	datahash := utils.Sha3(data)
	//compute data to Sign
	dataToSign := m.signData(datahash)
	sig, err := utils.SignData(privKey, dataToSign)
	if err != nil {
		return err
	}
	m.Signature = sig
	m.Sender = crypto.PubkeyToAddress(privKey.PublicKey)
	return nil
}

func (m *AnnounceDisposed) verifySignature(data []byte) error {
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	datatosign := m.signData(datahash)
	//should not change data's content,because its name is verify.
	var signature = make([]byte, signatureLength)
	copy(signature, data[len(data)-signatureLength:])
	hash := utils.Sha3(datatosign)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return err
	}
	m.Sender = utils.PubkeyToAddress(pubkey)
	return nil

}

/*
AnnounceDisposedResponse 收到 AnnounceDisposed为对方提供新的权益证明.
*/
type AnnounceDisposedResponse struct {
	EnvelopMessage
	LockSecretHash common.Hash
}

//NewAnnounceDisposedResponse create
func NewAnnounceDisposedResponse(bp *BalanceProof, lockSecretHash common.Hash) *AnnounceDisposedResponse {
	p := &AnnounceDisposedResponse{
		LockSecretHash: lockSecretHash,
	}
	p.CmdID = AnnounceDisposedTransferResponseCmdID
	p.EnvelopMessage.fromBalanceProof(bp)
	return p
}

//Pack is MessagePacker
func (m *AnnounceDisposedResponse) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID) //only one byte.
	_, err = buf.Write(m.LockSecretHash[:])
	m.EnvelopMessage.pack(buf)
	if err != nil {
		panic(fmt.Sprintf("RemoveExpiredHashlockTransfer Pack err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *AnnounceDisposedResponse) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if AnnounceDisposedTransferResponseCmdID != m.CmdID {
		return fmt.Errorf("AnnounceDisposedResponse UnPack cmdid should be  4,but get %d", m.CmdID)
	}
	_, err = buf.Read(m.LockSecretHash[:])
	err = m.EnvelopMessage.unpack(buf)
	if err != nil {
		return err
	}
	return m.EnvelopMessage.verifySignature(data)
}

//String is fmt.Stringer
func (m *AnnounceDisposedResponse) String() string {
	return fmt.Sprintf("Message{type=AnnounceDisposedResponse LockSecretHash=%s,%s}", utils.HPex(m.LockSecretHash), m.EnvelopMessage.String())
}

//SettleDataInMessage common part of settle request and response
type SettleDataInMessage struct {
	ChannelIDInMessage
	Participant1        common.Address
	Participant1Balance *big.Int
	Participant2        common.Address
	Participant2Balance *big.Int
}

//WithdrawRequestData for contract
type WithdrawRequestData struct {
	ChannelIDInMessage
	Participant1          common.Address
	Participant2          common.Address
	Participant1Balance   *big.Int
	Participant1Withdraw  *big.Int
	Participant1Signature []byte
}

/*
WithdrawRequest 向对方提出我要不关闭通道取现,节点应该标注通道不可用,然后再发送消息
*/
type WithdrawRequest struct {
	SignedMessage
	WithdrawRequestData
}

//NewWithdrawRequest create withdraw request from `WithdrawRequestData`
func NewWithdrawRequest(wd *WithdrawRequestData) *WithdrawRequest {
	m := &WithdrawRequest{
		WithdrawRequestData: *wd,
	}
	m.CmdID = WithdrawRequestCmdID
	return m
}
func (m *WithdrawRequest) String() string {
	return fmt.Sprintf("Message{type=WithdrawRequest Channel=%s-%d,Participant1=%s,Participant2=%s,"+
		"Participant1Balance=%s,Participant1Withdraw=%s}",
		utils.HPex(m.ChannelIdentifier), m.OpenBlockNumber,
		utils.APex2(m.Participant1), utils.APex2(m.Participant2), m.Participant1Balance, m.Participant1Withdraw,
	)
}

//Pack is MessagePacker
func (m *WithdrawRequest) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID)
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(m.Participant2[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Withdraw))
	_, err = buf.Write(m.Participant1Signature)
	_, err = buf.Write(m.Signature)
	if err != nil {
		panic(fmt.Sprintf("pack AnnounceDisposed err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *WithdrawRequest) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if WithdrawRequestCmdID != m.CmdID {
		return fmt.Errorf("WithdrawRequest UnPack cmdid expect=%d,got=%d", WithdrawRequestCmdID, m.CmdID)
	}
	_, err = buf.Read(m.ChannelIdentifier[:])
	err = binary.Read(buf, binary.BigEndian, &m.OpenBlockNumber)
	_, err = buf.Read(m.Participant1[:])
	_, err = buf.Read(m.Participant2[:])
	m.Participant1Balance = utils.ReadBigInt(buf)
	m.Participant1Withdraw = utils.ReadBigInt(buf)
	m.Participant1Signature = make([]byte, signatureLength)
	n, err := buf.Read(m.Participant1Signature)
	if err != nil || n != signatureLength {
		return fmt.Errorf("WithdrawRequest UnPack Participant1Signature err=%v,n=%d", err, n)
	}
	m.Signature = make([]byte, signatureLength)
	n, err = buf.Read(m.Signature)
	if err != nil || (n != signatureLength) {
		return fmt.Errorf("WithdrawRequest UnPack Signature err=%v,n=%d", err, n)
	}
	return m.verifySignature(data)
}
func (m *WithdrawRequest) verifySignature(data []byte) error {
	var err error
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	//should not change data's content,because its name is verify.
	var signature = data[len(dataWithoutSignature):]
	m.Sender, err = utils.Ecrecover(datahash, signature)
	if err != nil {
		return err
	}
	if m.Sender != m.Participant1 {
		return fmt.Errorf("WithdrawRequest signature error ,sender=%s,participant1=%s",
			utils.APex2(m.Sender), utils.APex2(m.Participant1))
	}
	fmt.Printf("vm=%s,sig=%s,p1sig=%s\n", m, hex.EncodeToString(m.Signature), hex.EncodeToString(m.Participant1Signature))
	datahash = utils.Sha3(m.signDataForContract())
	addr, err := utils.Ecrecover(datahash, m.Participant1Signature)
	if err != nil {
		return err
	}
	if m.Participant1 != addr {
		return fmt.Errorf("Participant1Signature err, Participant1=%s,but signed with other address=%s",
			utils.APex2(m.Participant1), addr.String())
	}
	return nil
}
func (m *WithdrawRequest) signDataForContract() []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractWithdrawProofMessageLength))
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Withdraw))
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.Cfg.ChainID))
	if err != nil {
		panic(fmt.Sprintf("signDataForContract err %s", err))
	}
	return buf.Bytes()
}

//Sign is SignedMessager
func (m *WithdrawRequest) Sign(key *ecdsa.PrivateKey, msg MessagePacker) (err error) {
	m.Participant1Signature, err = utils.SignData(key, m.signDataForContract())
	if err != nil {
		return
	}
	data := msg.Pack()
	m.Signature, err = utils.SignData(key, data)
	if err != nil {
		return
	}
	m.Sender = crypto.PubkeyToAddress(key.PublicKey)
	return
}

//WithdrawReponseData data for withdrawResponse
type WithdrawReponseData struct {
	ChannelIDInMessage
	Participant1          common.Address
	Participant2          common.Address
	Participant1Balance   *big.Int
	Participant1Withdraw  *big.Int
	Participant1Signature []byte
	Participant2Signature []byte
}

//WithdrawResponse is response for partner's withdraw request
type WithdrawResponse struct {
	SignedMessage
	WithdrawReponseData
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
}

//NewWithdrawResponse create withdraw response from `WithdrawReponseData`
func NewWithdrawResponse(wd *WithdrawReponseData, errorCode int, errorMsg string) *WithdrawResponse {
	m := &WithdrawResponse{
		WithdrawReponseData: *wd,
		ErrorCode:           errorCode,
		ErrorMsg:            errorMsg,
	}
	m.CmdID = WithdrawResponseCmdID
	return m
}

// NewErrorWithdrawResponseAndSign 创建返回错误信息的SettleResponse
func NewErrorWithdrawResponseAndSign(req *WithdrawRequest, privateKey *ecdsa.PrivateKey, errorCode int, errorMsg string) (res *WithdrawResponse) {
	res = &WithdrawResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	}
	// 这里填充无效数据,仅为了过验证
	res.CmdID = WithdrawResponseCmdID
	res.ChannelIdentifier = req.ChannelIdentifier
	res.ChannelIdentifier = req.ChannelIdentifier
	res.OpenBlockNumber = req.OpenBlockNumber
	res.Participant1 = utils.EmptyAddress
	res.Participant2 = crypto.PubkeyToAddress(privateKey.PublicKey)
	res.Participant1Balance = big.NewInt(0)
	res.Participant1Withdraw = big.NewInt(0)
	err2 := res.Sign(privateKey, res)
	if err2 != nil {
		panic(fmt.Sprintf("sign message for withdraw response err %s", err2))
	}
	return
}

func (m *WithdrawResponse) String() string {
	return fmt.Sprintf("Message{type=WithdrawResponse Channel=%s-%d,Participant1=%s,Participant2=%s,"+
		"Participant1Balance=%s,Participant1Withdraw=%s,ErrorCode=%d,ErrorMsg=%s",
		utils.HPex(m.ChannelIdentifier), m.OpenBlockNumber,
		utils.APex2(m.Participant1), utils.APex2(m.Participant2),
		m.Participant1Balance, m.Participant1Withdraw,
		m.ErrorCode, m.ErrorMsg,
	)
}

//Pack is MessagePacker
func (m *WithdrawResponse) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID)
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(m.Participant2[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Withdraw))
	_, err = buf.Write(m.Participant2Signature)
	// 2019-03 添加错误码及错误信息
	errCode := int32(m.ErrorCode)
	err = binary.Write(buf, binary.BigEndian, errCode)
	errorMsgBytes := utils.StringToBytes(m.ErrorMsg)
	errorMsgBytesLen := int32(len(errorMsgBytes))
	err = binary.Write(buf, binary.BigEndian, errorMsgBytesLen)
	if errorMsgBytesLen > 0 {
		_, err = buf.Write(errorMsgBytes)
	}
	_, err = buf.Write(m.Signature)
	if err != nil {
		panic(fmt.Sprintf("pack AnnounceDisposed err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *WithdrawResponse) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if WithdrawResponseCmdID != m.CmdID {
		return fmt.Errorf("WithdrawRequest UnPack cmdid expect=%d,got=%d", WithdrawRequestCmdID, m.CmdID)
	}
	_, err = buf.Read(m.ChannelIdentifier[:])
	err = binary.Read(buf, binary.BigEndian, &m.OpenBlockNumber)
	_, err = buf.Read(m.Participant1[:])
	_, err = buf.Read(m.Participant2[:])
	m.Participant1Balance = utils.ReadBigInt(buf)
	m.Participant1Withdraw = utils.ReadBigInt(buf)
	m.Participant2Signature = make([]byte, signatureLength)
	n, err := buf.Read(m.Participant2Signature)
	if err != nil || n != signatureLength {
		return fmt.Errorf("WithdrawRequest UnPack Participant1Signature err=%v,n=%d", err, n)
	}
	// 2019-03 添加错误码及错误信息
	var errCode int32
	err = binary.Read(buf, binary.BigEndian, &errCode)
	m.ErrorCode = int(errCode)
	var errorMsgBytesLen int32
	err = binary.Read(buf, binary.BigEndian, &errorMsgBytesLen)
	if errorMsgBytesLen > 0 {
		errorMsgBuf := make([]byte, errorMsgBytesLen)
		_, err = buf.Read(errorMsgBuf)
		m.ErrorMsg = utils.BytesToString(errorMsgBuf)
	}
	m.Signature = make([]byte, signatureLength)
	n, err = buf.Read(m.Signature)
	if err != nil || (n != signatureLength) {
		return fmt.Errorf("WithdrawRequest UnPack Signature err=%v,n=%d", err, n)
	}
	return m.verifySignature(data)
}
func (m *WithdrawResponse) verifySignature(data []byte) error {
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	//should not change data's content,because its name is verify.
	var signature = data[len(dataWithoutSignature):]
	addr, err := utils.Ecrecover(datahash, signature)
	if err != nil {
		return err
	}
	m.Sender = addr
	if m.Sender != m.Participant2 {
		return fmt.Errorf("WithdrawResponse signature error ,sender=%s,participant2=%s",
			utils.APex2(m.Sender), utils.APex2(m.Participant2))
	}
	datahash = utils.Sha3(m.signDataForContract())
	addr, err = utils.Ecrecover(datahash, m.Participant2Signature)
	if m.Participant2 != addr {
		return fmt.Errorf("Participant2Signature err, Participant2=%s,but signed with other address=%s",
			utils.APex2(m.Participant1), addr.String())
	}
	return nil
}
func (m *WithdrawResponse) signDataForContract() []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractWithdrawProofMessageLength))
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Withdraw))
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.Cfg.ChainID))
	if err != nil {
		panic(fmt.Sprintf("signDataForContract err %s", err))
	}
	return buf.Bytes()
}

//Sign is SignedMessager
func (m *WithdrawResponse) Sign(key *ecdsa.PrivateKey, msg MessagePacker) (err error) {
	m.Participant2Signature, err = utils.SignData(key, m.signDataForContract())
	if err != nil {
		return
	}
	data := msg.Pack()
	m.Signature, err = utils.SignData(key, data)
	m.Sender = crypto.PubkeyToAddress(key.PublicKey)
	return
}

//SettleRequestData for contract
type SettleRequestData struct {
	SettleDataInMessage
	Participant1Signature []byte
}

/*
SettleRequest 向对方提出我要合作关闭通道.
*/
type SettleRequest struct {
	SignedMessage
	SettleRequestData
}

//NewSettleRequest create  settle request from `SettleRequestData`
func NewSettleRequest(wd *SettleRequestData) *SettleRequest {
	m := &SettleRequest{
		SettleRequestData: *wd,
	}
	m.CmdID = SettleRequestCmdID
	return m
}
func (m *SettleRequest) String() string {
	return fmt.Sprintf("Message{type=SettleRequest Channel=%s-%d,Participant1=%s,Participant1Balance=%s,"+
		"Participant2=%s,Participant2Balance=%s}",
		utils.HPex(m.ChannelIdentifier), m.OpenBlockNumber,
		utils.APex2(m.Participant1), m.Participant1Balance,
		utils.APex2(m.Participant2), m.Participant2Balance,
	)
}

//Pack is MessagePacker
func (m *SettleRequest) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID)
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(m.Participant2[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant2Balance))
	_, err = buf.Write(m.Participant1Signature)
	_, err = buf.Write(m.Signature)
	if err != nil {
		panic(fmt.Sprintf("pack AnnounceDisposed err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *SettleRequest) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if SettleRequestCmdID != m.CmdID {
		return fmt.Errorf("SettleRequest UnPack cmdid expect=%d,got=%d", SettleRequestCmdID, m.CmdID)
	}
	_, err = buf.Read(m.ChannelIdentifier[:])
	err = binary.Read(buf, binary.BigEndian, &m.OpenBlockNumber)
	_, err = buf.Read(m.Participant1[:])
	m.Participant1Balance = utils.ReadBigInt(buf)
	_, err = buf.Read(m.Participant2[:])
	m.Participant2Balance = utils.ReadBigInt(buf)
	m.Participant1Signature = make([]byte, signatureLength)
	n, err := buf.Read(m.Participant1Signature)
	if err != nil || n != signatureLength {
		return fmt.Errorf("SettleRequest UnPack Participant1Signature err=%v,n=%d", err, n)
	}
	m.Signature = make([]byte, signatureLength)
	n, err = buf.Read(m.Signature)
	if err != nil || (n != signatureLength) {
		return fmt.Errorf("SettleRequest UnPack Signature err=%v,n=%d", err, n)
	}
	return m.verifySignature(data)
}
func (m *SettleRequest) verifySignature(data []byte) error {
	var err error
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	//should not change data's content,because its name is verify.
	var signature = data[len(dataWithoutSignature):]
	m.Sender, err = utils.Ecrecover(datahash, signature)
	if err != nil {
		return err
	}
	if m.Sender != m.Participant1 {
		return fmt.Errorf("SettleRequest signature error ,sender=%s,participant1=%s",
			utils.APex2(m.Sender), utils.APex2(m.Participant1))
	}
	fmt.Printf("vm=%s,sig=%s,p1sig=%s\n", m, hex.EncodeToString(m.Signature), hex.EncodeToString(m.Participant1Signature))
	datahash = utils.Sha3(m.SignDataForContract())
	addr, err := utils.Ecrecover(datahash, m.Participant1Signature)
	if err != nil {
		return err
	}
	if m.Participant1 != addr {
		return fmt.Errorf("Participant1Signature err, Participant1=%s,but signed with other address=%s",
			utils.APex2(m.Participant1), addr.String())
	}
	return nil
}

//SignDataForContract 生成合约调用签名数据
func (m *SettleRequest) SignDataForContract() []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractCooperativeSettleMessageLength))
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(m.Participant2[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant2Balance))
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.Cfg.ChainID))
	if err != nil {
		panic(fmt.Sprintf("signDataForContract err %s", err))
	}
	return buf.Bytes()
}

//Sign is SignedMessager
func (m *SettleRequest) Sign(key *ecdsa.PrivateKey, msg MessagePacker) (err error) {
	m.Participant1Signature, err = utils.SignData(key, m.SignDataForContract())
	if err != nil {
		return
	}
	data := msg.Pack()
	m.Signature, err = utils.SignData(key, data)
	if err != nil {
		return
	}
	m.Sender = crypto.PubkeyToAddress(key.PublicKey)
	return
}

//SettleResponseData for contract
type SettleResponseData struct {
	SettleDataInMessage
	Participant2Signature []byte
}

/*
SettleResponse 相应对方合作关闭通道要求
*/
type SettleResponse struct {
	SignedMessage
	SettleResponseData
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
}

//NewSettleResponse create settle response from `SettleResponseData`
func NewSettleResponse(wd *SettleResponseData, errorCode int, errorMsg string) *SettleResponse {
	m := &SettleResponse{
		SettleResponseData: *wd,
		ErrorCode:          errorCode,
		ErrorMsg:           errorMsg,
	}
	m.CmdID = SettleResponseCmdID
	return m
}

// NewErrorCooperativeSettleResponseAndSign 创建返回错误信息的SettleResponse
func NewErrorCooperativeSettleResponseAndSign(req *SettleRequest, privateKey *ecdsa.PrivateKey, errorCode int, errorMsg string) (res *SettleResponse) {
	res = &SettleResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
	}
	// 这里填充无效数据,仅为了过验证
	res.CmdID = SettleResponseCmdID
	res.ChannelIdentifier = req.ChannelIdentifier
	res.OpenBlockNumber = req.OpenBlockNumber
	res.Participant1 = utils.EmptyAddress
	res.Participant1Balance = big.NewInt(0)
	res.Participant2 = crypto.PubkeyToAddress(privateKey.PublicKey)
	res.Participant2Balance = big.NewInt(0)
	err2 := res.Sign(privateKey, res)
	if err2 != nil {
		panic(fmt.Sprintf("sign message for settle response err %s", err2))
	}
	return
}

func (m *SettleResponse) String() string {
	return fmt.Sprintf("Message{type=SettleResponse Channel=%s-%d,Participant1=%s,Participant1Balance=%s,"+
		"Participant2=%s,Participant2Balance=%s,ErrorCode=%d,ErrorMsg=%s}",
		utils.HPex(m.ChannelIdentifier), m.OpenBlockNumber,
		utils.APex2(m.Participant1), m.Participant1Balance,
		utils.APex2(m.Participant2), m.Participant2Balance,
		m.ErrorCode, m.ErrorMsg,
	)
}

//Pack is MessagePacker
func (m *SettleResponse) Pack() []byte {
	var err error
	buf := new(bytes.Buffer)
	err = m.WriteCmdStructToBuf(buf)
	//err = binary.Write(buf, binary.LittleEndian, m.CmdID)
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(m.Participant2[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant2Balance))
	_, err = buf.Write(m.Participant2Signature)
	// 2019-03 添加错误码及错误信息
	errCode := int32(m.ErrorCode)
	err = binary.Write(buf, binary.BigEndian, errCode)
	errorMsgBytes := utils.StringToBytes(m.ErrorMsg)
	errorMsgBytesLen := int32(len(errorMsgBytes))
	err = binary.Write(buf, binary.BigEndian, errorMsgBytesLen)
	if errorMsgBytesLen > 0 {
		_, err = buf.Write(errorMsgBytes)
	}
	_, err = buf.Write(m.Signature)
	if err != nil {
		panic(fmt.Sprintf("pack AnnounceDisposed err %s", err))
	}
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (m *SettleResponse) UnPack(data []byte) error {
	var err error
	buf := bytes.NewBuffer(data)
	err = m.ReadCmdStructFromBuf(buf)
	//err = binary.Read(buf, binary.LittleEndian, &t)
	if SettleResponseCmdID != m.CmdID {
		return fmt.Errorf("SettleResponse UnPack cmdid expect=%d,got=%d", SettleResponseCmdID, m.CmdID)
	}
	_, err = buf.Read(m.ChannelIdentifier[:])
	err = binary.Read(buf, binary.BigEndian, &m.OpenBlockNumber)
	_, err = buf.Read(m.Participant1[:])
	m.Participant1Balance = utils.ReadBigInt(buf)
	_, err = buf.Read(m.Participant2[:])
	m.Participant2Balance = utils.ReadBigInt(buf)
	m.Participant2Signature = make([]byte, signatureLength)
	n, err := buf.Read(m.Participant2Signature)
	if err != nil || n != signatureLength {
		return fmt.Errorf("SettleResponse UnPack Participant2Signature err=%v,n=%d", err, n)
	}
	// 2019-03 添加错误码及错误信息
	var errCode int32
	err = binary.Read(buf, binary.BigEndian, &errCode)
	m.ErrorCode = int(errCode)
	var errorMsgBytesLen int32
	err = binary.Read(buf, binary.BigEndian, &errorMsgBytesLen)
	if errorMsgBytesLen > 0 {
		errorMsgBuf := make([]byte, errorMsgBytesLen)
		_, err = buf.Read(errorMsgBuf)
		m.ErrorMsg = utils.BytesToString(errorMsgBuf)
	}
	m.Signature = make([]byte, signatureLength)
	n, err = buf.Read(m.Signature)
	if err != nil || (n != signatureLength) {
		return fmt.Errorf("SettleResponse UnPack Signature err=%v,n=%d", err, n)
	}
	return m.verifySignature(data)
}
func (m *SettleResponse) verifySignature(data []byte) error {
	var err error
	dataWithoutSignature := data[:len(data)-signatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	//should not change data's content,because its name is verify.
	var signature = data[len(dataWithoutSignature):]
	m.Sender, err = utils.Ecrecover(datahash, signature)
	if err != nil {
		return err
	}
	if m.Sender != m.Participant2 {
		return fmt.Errorf("SettleResponse signature error ,sender=%s,participant2=%s",
			utils.APex2(m.Sender), utils.APex2(m.Participant2))
	}
	fmt.Printf("vm=%s,sig=%s,p1sig=%s\n", m, hex.EncodeToString(m.Signature), hex.EncodeToString(m.Participant2Signature))
	datahash = utils.Sha3(m.SignDataForContract())
	addr, err := utils.Ecrecover(datahash, m.Participant2Signature)
	if err != nil {
		return err
	}
	if m.Participant2 != addr {
		return fmt.Errorf("Participant2Signature err, Participant2=%s,but signed with other address=%s",
			utils.APex2(m.Participant1), addr.String())
	}
	return nil
}

//SignDataForContract 生成合约调用数据
func (m *SettleResponse) SignDataForContract() []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(params.ContractSignaturePrefix)
	_, err = buf.Write([]byte(params.ContractCooperativeSettleMessageLength))
	_, err = buf.Write(m.Participant1[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant1Balance))
	_, err = buf.Write(m.Participant2[:])
	_, err = buf.Write(utils.BigIntTo32Bytes(m.Participant2Balance))
	_, err = buf.Write(m.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, m.OpenBlockNumber)
	_, err = buf.Write(utils.BigIntTo32Bytes(params.Cfg.ChainID))
	if err != nil {
		panic(fmt.Sprintf("signDataForContract err %s", err))
	}
	return buf.Bytes()
}

//Sign is SignedMessager
func (m *SettleResponse) Sign(key *ecdsa.PrivateKey, msg MessagePacker) (err error) {
	m.Participant2Signature, err = utils.SignData(key, m.SignDataForContract())
	if err != nil {
		return
	}
	data := msg.Pack()
	m.Signature, err = utils.SignData(key, data)
	if err != nil {
		return
	}
	m.Sender = crypto.PubkeyToAddress(key.PublicKey)
	return
}

//MessageMap contains all message can send and receive.
//DirectTransfer has been deprecated
var MessageMap = map[int]Messager{
	PingCmdID:                             new(Ping),
	AckCmdID:                              new(Ack),
	SecretRequestCmdID:                    new(SecretRequest),
	UnlockCmdID:                           new(UnLock),
	DirectTransferCmdID:                   new(DirectTransfer),
	RevealSecretCmdID:                     new(RevealSecret),
	MediatedTransferCmdID:                 new(MediatedTransfer),
	AnnounceDisposedTransferCmdID:         new(AnnounceDisposed),
	RemoveExpiredLockCmdID:                new(RemoveExpiredHashlockTransfer),
	AnnounceDisposedTransferResponseCmdID: new(AnnounceDisposedResponse),
	WithdrawRequestCmdID:                  new(WithdrawRequest),
	WithdrawResponseCmdID:                 new(WithdrawResponse),
	SettleRequestCmdID:                    new(SettleRequest),
	SettleResponseCmdID:                   new(SettleResponse),
	ErrorNotifyCmdID:                      new(ErrorNotify),
}

func init() {
	gob.Register(&Ack{})
	gob.Register(&CmdStruct{})
	gob.Register(&DirectTransfer{})
	gob.Register(&EnvelopMessage{})
	gob.Register(&MediatedTransfer{})
	gob.Register(&Ping{})
	gob.Register(&AnnounceDisposed{})
	gob.Register(&UnLock{})
	gob.Register(&SecretRequest{})
	gob.Register(&RemoveExpiredHashlockTransfer{})
	gob.Register(&AnnounceDisposedResponse{})
	gob.Register(&WithdrawRequest{})
	gob.Register(&WithdrawResponse{})
	gob.Register(&SettleRequest{})
	gob.Register(&SettleResponse{})
}
