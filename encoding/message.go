package encoding

import (
	"bytes"
	"encoding/binary"

	"crypto/ecdsa"

	"math/big"

	"io"

	"errors"
	"fmt"

	"encoding/gob"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type MessageType int

const ACK_CMDID = 0
const PING_CMDID = 1
const SECRETREQUEST_CMDID = 3
const SECRET_CMDID = 4
const DIRECTTRANSFER_CMDID = 5
const MEDIATEDTRANSFER_CMDID = 7
const REFUNDTRANSFER_CMDID = 8
const REVEALSECRET_CMDID = 11

const SignatureLength = 65
const TokenLength = 20

var errPacketLength = errors.New("packet length error")

type MessagePacker interface {
	//pack message to byte array
	Pack() []byte
}
type MessageUnpacker interface {
	//unpack message from byte array
	UnPack(data []byte) error
}

//type MessageVerifier interface {
//	//verify signature is valid ,if valid return true and sender's address
//	VerifySignature() (bool, common.Address)
//}
type MessagePackerUnpacker interface {
	MessagePacker
	MessageUnpacker
}

type Messager interface {
	Cmd() int
	Tag() interface{}
	SetTag(tag interface{})
	Name() string
	String() string
	MessagePackerUnpacker
}

type CmdStruct struct {
	CmdId       int32
	InternalTag interface{} //for save to database
}

func (this *CmdStruct) Cmd() int {
	return int(this.CmdId)
}
func (this *CmdStruct) Tag() interface{} {
	return this.InternalTag
}
func (this *CmdStruct) SetTag(tag interface{}) {
	this.InternalTag = tag
}
func (this *CmdStruct) Name() string {
	return MessageType(this.CmdId).String()
}

type SignedMessager interface {
	Messager
	GetSender() common.Address
	Sign(priveKey *ecdsa.PrivateKey, pack MessagePacker) error
	VerifySignature(data []byte) error
	GetSignature() []byte
}
type EnvelopMessager interface {
	SignedMessager
	GetEnvelopMessage() *EnvelopMessage
}

// String return the string representation of message type.
func (t MessageType) String() string {
	switch t {
	case ACK_CMDID:
		return "Ack"
	case PING_CMDID:
		return "Ping"
	case SECRETREQUEST_CMDID:
		return "SecretRequest"
	case SECRET_CMDID:
		return "Secret"
	case DIRECTTRANSFER_CMDID:
		return "DirectTransfer"
	case MEDIATEDTRANSFER_CMDID:
		return "MediatedTransfer"
	case REFUNDTRANSFER_CMDID:
		return "RefundTransfer"
	case REVEALSECRET_CMDID:
		return "RevealSecret"
	default:
		return "<unknown>"
	}
}

/*All accepted messages should be confirmed by an `Ack` which echoes the
orginals Message hash.

We don'T sign Acks because attack vector can be mitigated and to speed up
things.
*/
type Ack struct {
	CmdStruct
	Sender common.Address
	Echo   common.Hash
}

func NewAck(sender common.Address, echo common.Hash) *Ack {
	return &Ack{
		CmdStruct: CmdStruct{CmdId: ACK_CMDID},
		Sender:    sender,
		Echo:      echo,
	}
}

func (this *Ack) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId)
	buf.Write(this.Sender[:])
	buf.Write(this.Echo[:])
	return buf.Bytes()
}

func (this *Ack) UnPack(data []byte) error {
	var t int32
	this.CmdId = ACK_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != this.CmdId {
		panic(fmt.Sprint("Ack Unpack cmdid should be 0,but get %d", t))
	}
	buf.Read(this.Sender[:])
	n, err := buf.Read(this.Echo[:])
	if err != nil {
		return err
	}
	if n != len(this.Echo) {
		return errPacketLength
	}
	return nil
}
func (this *Ack) String() string {
	return fmt.Sprintf("Message{type=Ack sender=%s,echo=%s}", utils.APex2(this.Sender), utils.HPex(this.Echo))
}

type SignedMessage struct {
	CmdStruct
	Sender    common.Address
	Signature []byte
}

func (this *SignedMessage) GetSender() common.Address {
	return this.Sender
}
func (this *SignedMessage) GetSignature() []byte {
	return this.Signature
}
func (this *SignedMessage) Sign(priveKey *ecdsa.PrivateKey, pack MessagePacker) error {
	if len(this.Signature) > 0 {
		log.Warn("duplicate sign")
		return errors.New("duplicate sign")
	}
	this.Signature = SignMessage(priveKey, pack)
	this.Sender = crypto.PubkeyToAddress(priveKey.PublicKey)
	return nil
}
func (this *SignedMessage) VerifySignature(data []byte) error {
	sender, err := VerifyMessage(data)
	if err != nil {
		return err
	}
	this.Sender = sender
	return nil
}
func SignMessage(privKey *ecdsa.PrivateKey, pack MessagePacker) []byte {
	data := pack.Pack()
	sig, err := utils.SignData(privKey, data)
	if err != nil {
		panic(fmt.Sprintf("SignMessage error %s", err))
	}
	return sig
}

func HashMessage(pack MessagePacker) common.Hash {
	return utils.Sha3(pack.Pack())
}
func HashMessageWithoutSignature(pack MessagePacker) common.Hash {
	data := pack.Pack()
	if len(data) > SignatureLength {
		data = data[:len(data)-SignatureLength]
	}
	return utils.Sha3(data)
}

func VerifyMessage(data []byte) (sender common.Address, err error) {
	messageData := data[:len(data)-SignatureLength]
	signature := make([]byte, SignatureLength)
	copy(signature, data[len(data)-SignatureLength:])
	hash := utils.Sha3(messageData)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return
	}
	sender = utils.PubkeyToAddress(pubkey)
	return
}

type Ping struct {
	SignedMessage
	Nonce int64
}

func NewPing(nonce int64) *Ping {
	p := &Ping{
		//SignedMessage:SignedMessage{CmdStruct: CmdStruct{CmdId: PING_CMDID}},
		Nonce: nonce,
	}
	p.CmdId = PING_CMDID
	return p
}

func (this *Ping) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId) //只有一个字节..
	binary.Write(buf, binary.BigEndian, this.Nonce)
	buf.Write(this.Signature)
	return buf.Bytes()
}

func (this *Ping) UnPack(data []byte) error {
	var t int32
	this.CmdId = PING_CMDID
	if len(data) != 77 { //stun response here
		return errPacketLength
	}

	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != this.CmdId {
		return fmt.Errorf("Ping Unpack cmdid should be  1,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &this.Nonce)
	this.Signature = make([]byte, SignatureLength)
	buf.Read(this.Signature)
	err := this.SignedMessage.VerifySignature(data)
	if err != nil {
		return err
	}
	return nil
}

func (this *Ping) String() string {
	return fmt.Sprintf("Message{type=Ping Nonce=%d,sender=%s, has signature=%v}", this.Nonce, utils.APex2(this.Sender), len(this.Signature) != 0)
}

//Requests the secret which unlocks a hashlock.
type SecretRequest struct {
	SignedMessage
	Identifier uint64
	HashLock   common.Hash
	Amount     *big.Int
}

func NewSecretRequest(Identifier uint64, hashLock common.Hash, amount *big.Int) *SecretRequest {
	p := &SecretRequest{
		Identifier: Identifier,
		HashLock:   hashLock,
		Amount:     new(big.Int).Set(amount),
	}
	p.CmdId = SECRETREQUEST_CMDID
	return p
}

func (this *SecretRequest) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId) //只有一个字节..
	binary.Write(buf, binary.BigEndian, this.Identifier)
	buf.Write(this.HashLock[:])
	buf.Write(utils.BigIntTo32Bytes(this.Amount))
	buf.Write(this.Signature)
	return buf.Bytes()
}
func readBigInt(reader io.Reader) *big.Int {
	bi := new(big.Int)
	tmpbuf := make([]byte, 32)
	reader.Read(tmpbuf)
	bi.SetBytes(tmpbuf)
	return bi
}
func (this *SecretRequest) UnPack(data []byte) error {
	var t int32
	this.CmdId = SECRETREQUEST_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != this.CmdId {
		return fmt.Errorf("SecretRequest Unpack cmdid should be  3,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &this.Identifier)
	buf.Read(this.HashLock[:])
	this.Amount = readBigInt(buf)
	this.Signature = make([]byte, SignatureLength)
	n, err := buf.Read(this.Signature)
	if err != nil {
		return err
	}
	if n != SignatureLength {
		return errPacketLength
	}
	err = this.VerifySignature(data)
	return err
}

func (this *SecretRequest) String() string {
	return fmt.Sprintf("Message{type=SecretRequest identifier=%d,hashlock=%s,amount=%s,sender=%s,has signature=%v}", this.Identifier,
		utils.HPex(this.HashLock), this.Amount.String(), utils.APex2(this.Sender), len(this.Signature) != 0)
}

/*
   Message used to reveal a secret to party known to have interest in it.

   This message is not sufficient for state changes in the raiden Channel, the
   reason is that a node participating in split transfer or in both mediated
   transfer for an exchange might can reveal the secret to it's partners, but
   that must not update the internal channel state.
*/
type RevealSecret struct {
	SignedMessage
	Secret   common.Hash
	hashLock common.Hash
}

func NewRevealSecret(secret common.Hash) *RevealSecret {
	p := &RevealSecret{
		Secret: secret,
	}
	p.CmdId = REVEALSECRET_CMDID
	return p
}
func CloneRevealSecret(rs *RevealSecret) *RevealSecret {
	rs2 := *rs
	return &rs2
}
func (this *RevealSecret) HashLock() common.Hash {
	if bytes.Equal(this.hashLock[:], utils.EmptyHash[:]) {
		this.hashLock = utils.Sha3(this.Secret[:])
	}
	return this.hashLock
}
func (this *RevealSecret) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId) //只有一个字节..
	buf.Write(this.Secret[:])
	buf.Write(this.Signature)
	return buf.Bytes()
}
func (this *RevealSecret) UnPack(data []byte) error {
	var t int32
	this.CmdId = REVEALSECRET_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != this.CmdId {
		return fmt.Errorf("RevealSecret Unpack cmdid should be  11,but get %d", t)
	}
	buf.Read(this.Secret[:])
	this.Signature = make([]byte, SignatureLength)
	n, err := buf.Read(this.Signature)
	if err != nil {
		return err
	}
	if n != SignatureLength {
		return errPacketLength
	}
	return this.VerifySignature(data)
}
func (this *RevealSecret) String() string {
	return fmt.Sprintf("Message{type=RevealSecret,hashlock=%s,secret=%s,sender=%s,has signature=%v}", utils.HPex(this.hashLock),
		utils.HPex(this.Secret), utils.APex2(this.Sender), len(this.Signature) != 0)
}

type EnvelopMessage struct {
	SignedMessage
	Nonce          int64
	Channel        common.Address
	TransferAmount *big.Int //已经转给对方的数量(确认过了,随时可以到链上提现的)
	Locksroot      common.Hash
	Identifier     uint64
}

func (this *EnvelopMessage) String() string {
	return fmt.Sprintf("EnvelopMessage{Nonce=%d,Channel=%s,TransferAmount=%s,Locksroot=%s,Identifier=%d sender=%s,has signature=%v}", this.Nonce,
		utils.APex2(this.Channel), this.TransferAmount, utils.HPex(this.Locksroot), this.Identifier, utils.APex2(this.Sender), len(this.Signature) != 0)
}
func (this *EnvelopMessage) signData(datahash common.Hash) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, this.Nonce)
	buf.Write(utils.BigIntTo32Bytes(this.TransferAmount))
	buf.Write(this.Locksroot[:])
	buf.Write(this.Channel[:])
	buf.Write(datahash[:])
	dataToSign := buf.Bytes()
	return dataToSign
}

/*
to sign data(once+transferamount+locksroot+channel+hash(data))
*/
func (this *EnvelopMessage) Sign(privKey *ecdsa.PrivateKey, msg MessagePacker) error {
	data := msg.Pack() //before signed, sign twice will be error
	datahash := utils.Sha3(data)
	//compute data to sign
	dataToSign := this.signData(datahash)
	sig, err := utils.SignData(privKey, dataToSign)
	if err != nil {
		return err
	}
	this.Signature = sig
	this.Sender = crypto.PubkeyToAddress(privKey.PublicKey)
	return nil
}
func (this *EnvelopMessage) VerifySignature(data []byte) error {
	dataWithoutSignature := data[:len(data)-SignatureLength]
	datahash := utils.Sha3(dataWithoutSignature)
	datatosign := this.signData(datahash)
	//should not change data's content,because its name is verify.
	var signature = make([]byte, SignatureLength)
	copy(signature, data[len(data)-SignatureLength:])
	hash := utils.Sha3(datatosign)
	signature[len(signature)-1] -= 27 //why?
	pubkey, err := crypto.Ecrecover(hash[:], signature)
	if err != nil {
		return err
	}
	this.Sender = utils.PubkeyToAddress(pubkey)
	return nil

}
func (this *EnvelopMessage) GetEnvelopMessage() *EnvelopMessage {
	return this
}

/*
Message used to do state changes on a partner Raiden Channel.

Locksroot changes need to be synchronized among both participants, the
protocol is for only the side unlocking to send the Secret message allowing
the other party to withdraw.
*/
type Secret struct {
	EnvelopMessage
	Secret common.Hash
}

func (this *Secret) HashLock() common.Hash {
	return utils.Sha3(this.Secret[:])
}
func NewSecret(Identifier uint64, nonce int64, channel common.Address,
	transferamount *big.Int, locksroot common.Hash, secret common.Hash) *Secret {
	p := &Secret{
		Secret: secret,
	}
	p.Identifier = Identifier
	p.CmdId = SECRET_CMDID
	p.Nonce = nonce
	p.Channel = channel
	p.TransferAmount = new(big.Int).Set(transferamount)
	p.Locksroot = locksroot
	return p
}

func (this *Secret) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId) //只有一个字节..
	binary.Write(buf, binary.BigEndian, this.Identifier)
	buf.Write(this.Secret[:])
	binary.Write(buf, binary.BigEndian, this.Nonce)
	buf.Write(this.Channel[:])
	buf.Write(utils.BigIntTo32Bytes(this.TransferAmount))
	buf.Write(this.Locksroot[:])
	buf.Write(this.Signature)
	return buf.Bytes()
}
func (this *Secret) UnPack(data []byte) error {
	var t int32
	this.CmdId = SECRET_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != this.CmdId {
		return fmt.Errorf("Ack Secret cmdid should be  4,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &this.Identifier)
	buf.Read(this.Secret[:])

	binary.Read(buf, binary.BigEndian, &this.Nonce)
	buf.Read(this.Channel[:])
	this.TransferAmount = readBigInt(buf)
	buf.Read(this.Locksroot[:])
	this.Signature = make([]byte, SignatureLength)
	n, err := buf.Read(this.Signature)
	if err != nil {
		return err
	}
	if n != SignatureLength {
		return errors.New("packet length error")
	}
	return this.EnvelopMessage.VerifySignature(data)
}
func (this *Secret) String() string {
	return fmt.Sprintf("Message{type=Secret secret=%s,%s}", utils.HPex(this.Secret), this.EnvelopMessage.String())
}

/*
""" A direct token exchange, used when both participants have a previously
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
        recipient: The address of the raiden node participating in the channel.
        locksroot: The root of a merkle tree which records the current
            outstanding locks.
    """
*/
type DirectTransfer struct {
	EnvelopMessage
	Token     common.Address //20bytes
	Recipient common.Address //20bytes
}

func (this *DirectTransfer) String() string {
	return fmt.Sprintf("Message{type=DirectTransfer token=%s,recipient=%s,%s}", utils.APex2(this.Token),
		utils.APex2(this.Recipient), this.EnvelopMessage.String())
}
func NewDirectTransfer(identifier uint64, nonce int64, token common.Address,
	channel common.Address, transferAmount *big.Int,
	recipient common.Address, locksroot common.Hash) *DirectTransfer {
	p := &DirectTransfer{
		Token:     token,
		Recipient: recipient,
	}
	p.Identifier = identifier
	p.CmdId = DIRECTTRANSFER_CMDID
	p.Nonce = nonce
	p.Channel = channel
	p.TransferAmount = new(big.Int).Set(transferAmount)
	p.Locksroot = locksroot
	return p
}

func (this *DirectTransfer) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId) //只有一个字节..
	binary.Write(buf, binary.BigEndian, this.Nonce)
	binary.Write(buf, binary.BigEndian, this.Identifier)
	buf.Write(this.Token[:])
	buf.Write(this.Channel[:])
	buf.Write(this.Recipient[:])
	buf.Write(utils.BigIntTo32Bytes(this.TransferAmount))
	buf.Write(this.Locksroot[:]) //todo locksroot pack unpack maybe error
	buf.Write(this.Signature)
	return buf.Bytes()
}
func (this *DirectTransfer) UnPack(data []byte) error {
	var t int32
	this.CmdId = DIRECTTRANSFER_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != this.CmdId {
		return errors.New("DirectTransfer unpack cmdid error")
	}
	binary.Read(buf, binary.BigEndian, &this.Nonce)
	binary.Read(buf, binary.BigEndian, &this.Identifier)
	buf.Read(this.Token[:])
	buf.Read(this.Channel[:])
	buf.Read(this.Recipient[:])
	this.TransferAmount = readBigInt(buf)
	buf.Read(this.Locksroot[:])
	this.Signature = make([]byte, SignatureLength)
	n, err := buf.Read(this.Signature)
	if err != nil {
		return err
	}
	if n != SignatureLength {
		return errPacketLength
	}
	return this.EnvelopMessage.VerifySignature(data)
}

type Lock struct {
	Expiration int64 //expiration block number
	Amount     *big.Int
	HashLock   common.Hash
}

func (this *Lock) AsBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, this.Expiration)
	buf.Write(utils.BigIntTo32Bytes(this.Amount))
	buf.Write(this.HashLock[:])
	return buf.Bytes()
}
func (this *Lock) FromBytes(locksencoded []byte) {
	buf := bytes.NewBuffer(locksencoded)
	binary.Read(buf, binary.BigEndian, &this.Expiration)
	this.Amount = readBigInt(buf)
	buf.Read(this.HashLock[:])
}

/*
"""
    A MediatedTransfer has a `target` address to which a chain of transfers shall
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
    """
*/
type MediatedTransfer struct {
	EnvelopMessage
	Expiration int64
	Token      common.Address
	Recipient  common.Address
	Target     common.Address
	Initiator  common.Address
	HashLock   common.Hash
	Amount     *big.Int //此次转给对方的数量
	Fee        *big.Int
}

func (this *MediatedTransfer) String() string {
	return fmt.Sprintf("Message{type=MediatedTransfer expiration=%d,token=%s,recipient=%s,target=%s,initiator=%s,hashlock=%s,amount=%s,fee=%s,%s}",
		this.Expiration, utils.APex2(this.Token), utils.APex2(this.Recipient), utils.APex2(this.Target), utils.APex2(this.Initiator),
		utils.HPex(this.HashLock), this.Amount, this.Fee, this.EnvelopMessage.String())
}
func NewMediatedTransfer(identifier uint64, nonce int64, token common.Address,
	channel common.Address, transferAmount *big.Int,
	recipient common.Address, locksroot common.Hash, lock *Lock,
	target common.Address, initiator common.Address, fee *big.Int) *MediatedTransfer {
	p := &MediatedTransfer{
		Token:     token,
		Recipient: recipient,
		Target:    target,
		Initiator: initiator,
		Fee:       new(big.Int).Set(fee),
	}
	p.Identifier = identifier
	p.Nonce = nonce
	p.TransferAmount = new(big.Int).Set(transferAmount)
	p.Locksroot = locksroot //包含此次未完全完成交易的merkletree root
	p.CmdId = MEDIATEDTRANSFER_CMDID
	p.Channel = channel
	p.Expiration = lock.Expiration
	p.HashLock = lock.HashLock
	p.Amount = new(big.Int).Set(lock.Amount)
	return p
}

func (this *MediatedTransfer) GetLock() *Lock {
	return &Lock{
		Expiration: this.Expiration,
		Amount:     this.Amount,
		HashLock:   this.HashLock,
	}
}

func (this *MediatedTransfer) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, this.CmdId) //只有一个字节..
	binary.Write(buf, binary.BigEndian, this.Nonce)
	binary.Write(buf, binary.BigEndian, this.Identifier)
	binary.Write(buf, binary.BigEndian, this.Expiration)
	buf.Write(this.Token[:])
	buf.Write(this.Channel[:])
	buf.Write(this.Recipient[:])
	buf.Write(this.Target[:])
	buf.Write(this.Initiator[:])
	buf.Write(this.Locksroot[:])
	buf.Write(this.HashLock[:])
	buf.Write(utils.BigIntTo32Bytes(this.TransferAmount))
	buf.Write(utils.BigIntTo32Bytes(this.Amount))
	buf.Write(utils.BigIntTo32Bytes(this.Fee))
	buf.Write(this.Signature)
	return buf.Bytes()
}
func (this *MediatedTransfer) UnPack(data []byte) error {
	var t int32
	//this.CmdId = MEDIATEDTRANSFER_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	this.CmdId = t
	if this.CmdId != MEDIATEDTRANSFER_CMDID && this.CmdId != REFUNDTRANSFER_CMDID {
		return errors.New("MediatedTransfer unpack cmd error")
	}
	binary.Read(buf, binary.BigEndian, &this.Nonce)
	binary.Read(buf, binary.BigEndian, &this.Identifier)
	binary.Read(buf, binary.BigEndian, &this.Expiration)
	buf.Read(this.Token[:])
	buf.Read(this.Channel[:])
	buf.Read(this.Recipient[:])
	buf.Read(this.Target[:])
	buf.Read(this.Initiator[:])
	buf.Read(this.Locksroot[:])
	buf.Read(this.HashLock[:])
	this.TransferAmount = readBigInt(buf)
	this.Amount = readBigInt(buf)
	this.Fee = readBigInt(buf)
	this.Signature = make([]byte, SignatureLength)
	n, err := buf.Read(this.Signature)
	if err != nil {
		return err
	}
	if n != SignatureLength {
		return errPacketLength
	}
	return this.VerifySignature(data)
}

type RefundTransfer struct {
	MediatedTransfer
}

func (this *RefundTransfer) String() string {
	return fmt.Sprintf("Message{type=RefundTransfer expiration=%d,token=%s,recipient=%s,target=%s,initiator=%s,hashlock=%s,amount=%s,fee=%s,%s}",
		this.Expiration, utils.APex2(this.Token), utils.APex2(this.Recipient), utils.APex2(this.Target), utils.APex2(this.Initiator),
		utils.HPex(this.HashLock), this.Amount, this.Fee, this.EnvelopMessage.String())
}
func NewRefundTransfer(identifier uint64, nonce int64, token common.Address,
	channel common.Address, transferAmount *big.Int,
	recipient common.Address, locksroot common.Hash, lock *Lock,
	target common.Address, initiator common.Address, fee *big.Int) *RefundTransfer {
	p := &RefundTransfer{}
	p.MediatedTransfer = *(NewMediatedTransfer(identifier, nonce, token, channel, transferAmount, recipient,
		locksroot, lock, target, initiator, fee))
	p.CmdId = REFUNDTRANSFER_CMDID
	return p
}
func NewRefundTransferFromMediatedTransfer(mtr *MediatedTransfer) *RefundTransfer {
	p := &RefundTransfer{}
	p.MediatedTransfer = *mtr
	p.CmdId = REFUNDTRANSFER_CMDID
	return p
}
func IsLockedTransfer(msg Messager) bool {
	return msg.Cmd() == REFUNDTRANSFER_CMDID || msg.Cmd() == MEDIATEDTRANSFER_CMDID
}

//make sure tr is locked transfer
func GetMtrFromLockedTransfer(tr Messager) (mtr *MediatedTransfer) {
	if !IsLockedTransfer(tr) {
		panic("getmtr should never panic")
	}
	mtr, ok := tr.(*MediatedTransfer)
	if !ok {
		rtr, ok := tr.(*RefundTransfer)
		if ok {
			mtr = &rtr.MediatedTransfer
		}
	}
	return
}

var MessageMap = map[int]Messager{
	PING_CMDID:             new(Ping),
	ACK_CMDID:              new(Ack),
	SECRETREQUEST_CMDID:    new(SecretRequest),
	SECRET_CMDID:           new(Secret),
	DIRECTTRANSFER_CMDID:   new(DirectTransfer),
	REVEALSECRET_CMDID:     new(RevealSecret),
	MEDIATEDTRANSFER_CMDID: new(MediatedTransfer),
	REFUNDTRANSFER_CMDID:   new(RefundTransfer),
}

func init() {
	gob.Register(&Ack{})
	gob.Register(&CmdStruct{})
	gob.Register(&DirectTransfer{})
	gob.Register(&EnvelopMessage{})
	gob.Register(&Lock{})
	gob.Register(&MediatedTransfer{})
	gob.Register(&Ping{})
	gob.Register(&RefundTransfer{})
	gob.Register(&Secret{})
	gob.Register(&SecretRequest{})
}
