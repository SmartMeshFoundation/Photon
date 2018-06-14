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

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

//MessageType is the type of message for receive and send
type MessageType int

//AckCmdID id of Ack message
const (
	AckCmdID = 0
	//PingCmdID id of ping message
	PingCmdID = 1
	//SecretRequestCmdID id of SecretRequest message
	SecretRequestCmdID = 3
	//SecretCmdID id of Secret message
	SecretCmdID = 4
	//DirectTransferCmdID id of DirectTransfer, it's now deprecated
	DirectTransferCmdID = 5
	//MediatedTransferCmdID id of MediatedTransfer
	MediatedTransferCmdID = 7
	//RefundTransferCmdID id of RefundTransfer message
	RefundTransferCmdID = 8
	//RevealSecretCmdID id of RevealSecret message
	RevealSecretCmdID = 11
	//RemoveExpiredHashLockCmdID id of RemoveExpiredHashlock message
	RemoveExpiredHashLockCmdID = 13
)

const signatureLength = 65
const tokenLength = 20

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
	CmdID       int32
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
	case SecretCmdID:
		return "Secret"
	case DirectTransferCmdID:
		return "DirectTransfer"
	case MediatedTransferCmdID:
		return "MediatedTransfer"
	case RefundTransferCmdID:
		return "RefundTransfer"
	case RevealSecretCmdID:
		return "RevealSecret"
	case RemoveExpiredHashLockCmdID:
		return "RemoveExpiredHashlock"
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, ack.CmdID)
	buf.Write(ack.Sender[:])
	buf.Write(ack.Echo[:])
	return buf.Bytes()
}

//UnPack is implements of MessageUnpacker
func (ack *Ack) UnPack(data []byte) error {
	var t int32
	ack.CmdID = AckCmdID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != ack.CmdID {
		panic(fmt.Sprintf("Ack Unpack cmdid should be 0,but get %d", t))
	}
	buf.Read(ack.Sender[:])
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
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, p.CmdID) //only one byte
	binary.Write(buf, binary.BigEndian, p.Nonce)
	buf.Write(p.Signature)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (p *Ping) UnPack(data []byte) error {
	var t int32
	p.CmdID = PingCmdID
	if len(data) != 77 { //stun response here
		return errPacketLength
	}

	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != p.CmdID {
		return fmt.Errorf("Ping Unpack cmdid should be  1,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &p.Nonce)
	p.Signature = make([]byte, signatureLength)
	buf.Read(p.Signature)
	err := p.SignedMessage.verifySignature(data)
	if err != nil {
		return err
	}
	return nil
}

//String is fmt.Stringer
func (p *Ping) String() string {
	return fmt.Sprintf("Message{type=Ping nonce=%d,sender=%s, has signature=%v}", p.Nonce, utils.APex2(p.Sender), len(p.Signature) != 0)
}

//SecretRequest Requests the secret which unlocks a hashlock.
type SecretRequest struct {
	SignedMessage
	Identifier uint64
	HashLock   common.Hash
	Amount     *big.Int
}

//NewSecretRequest create SecretRequest
func NewSecretRequest(Identifier uint64, hashLock common.Hash, amount *big.Int) *SecretRequest {
	p := &SecretRequest{
		Identifier: Identifier,
		HashLock:   hashLock,
		Amount:     new(big.Int).Set(amount),
	}
	p.CmdID = SecretRequestCmdID
	return p
}

//Pack is MessagePacker
func (sr *SecretRequest) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, sr.CmdID) //only one byte..
	binary.Write(buf, binary.BigEndian, sr.Identifier)
	buf.Write(sr.HashLock[:])
	buf.Write(utils.BigIntTo32Bytes(sr.Amount))
	buf.Write(sr.Signature)
	return buf.Bytes()
}
func readBigInt(reader io.Reader) *big.Int {
	bi := new(big.Int)
	tmpbuf := make([]byte, 32)
	reader.Read(tmpbuf)
	bi.SetBytes(tmpbuf)
	return bi
}

//UnPack is MessageUnpacker
func (sr *SecretRequest) UnPack(data []byte) error {
	var t int32
	sr.CmdID = SecretRequestCmdID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != sr.CmdID {
		return fmt.Errorf("SecretRequest Unpack cmdid should be  3,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &sr.Identifier)
	buf.Read(sr.HashLock[:])
	sr.Amount = readBigInt(buf)
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
	return fmt.Sprintf("Message{type=SecretRequest identifier=%d,hashlock=%s,amount=%s,sender=%s,has signature=%v}", sr.Identifier,
		utils.HPex(sr.HashLock), sr.Amount.String(), utils.APex2(sr.Sender), len(sr.Signature) != 0)
}

/*
RevealSecret used to reveal a secret to party known to have interest in it.

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

//NewRevealSecret create RevealSecret
func NewRevealSecret(secret common.Hash) *RevealSecret {
	p := &RevealSecret{
		Secret: secret,
	}
	p.CmdID = RevealSecretCmdID
	return p
}

//CloneRevealSecret clones a RevealSecret Message
func CloneRevealSecret(rs *RevealSecret) *RevealSecret {
	rs2 := *rs
	return &rs2
}

//HashLock return hash of secret
func (rs *RevealSecret) HashLock() common.Hash {
	if bytes.Equal(rs.hashLock[:], utils.EmptyHash[:]) {
		rs.hashLock = utils.Sha3(rs.Secret[:])
	}
	return rs.hashLock
}

//Pack is MessagePacker
func (rs *RevealSecret) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rs.CmdID) //only one byte.
	buf.Write(rs.Secret[:])
	buf.Write(rs.Signature)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (rs *RevealSecret) UnPack(data []byte) error {
	var t int32
	rs.CmdID = RevealSecretCmdID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != rs.CmdID {
		return fmt.Errorf("RevealSecret Unpack cmdid should be  11,but get %d", t)
	}
	buf.Read(rs.Secret[:])
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
	return fmt.Sprintf("Message{type=RevealSecret,hashlock=%s,secret=%s,sender=%s,has signature=%v}", utils.HPex(rs.hashLock),
		utils.HPex(rs.Secret), utils.APex2(rs.Sender), len(rs.Signature) != 0)
}

//EnvelopMessage is general part of message that contains a new balanceproof
type EnvelopMessage struct {
	SignedMessage
	Nonce          int64
	Channel        common.Address
	TransferAmount *big.Int //The number has been transferred to the other party
	Locksroot      common.Hash
	Identifier     uint64
}

//String is fmt.Stringer
func (m *EnvelopMessage) String() string {
	return fmt.Sprintf("EnvelopMessage{nonce=%d,Channel=%s,TransferAmount=%s,Locksroot=%s,Identifier=%d sender=%s,has signature=%v}", m.Nonce,
		utils.APex2(m.Channel), m.TransferAmount, utils.HPex(m.Locksroot), m.Identifier, utils.APex2(m.Sender), len(m.Signature) != 0)
}
func (m *EnvelopMessage) signData(datahash common.Hash) []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, m.Nonce)
	buf.Write(utils.BigIntTo32Bytes(m.TransferAmount))
	buf.Write(m.Locksroot[:])
	buf.Write(m.Channel[:])
	buf.Write(datahash[:])
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
	return nil
}

//GetEnvelopMessage return EnvelopMessage
func (m *EnvelopMessage) GetEnvelopMessage() *EnvelopMessage {
	return m
}

/*
Secret Message used to do state changes on a partner Raiden Channel.

Locksroot changes need to be synchronized among both participants, the
protocol is for only the side unlocking to send the Secret message allowing
the other party to withdraw.
*/
type Secret struct {
	EnvelopMessage
	Secret common.Hash
}

//HashLock is Hash of secret
func (s *Secret) HashLock() common.Hash {
	return utils.Sha3(s.Secret[:])
}

//NewSecret create Secret message
func NewSecret(Identifier uint64, nonce int64, channel common.Address,
	transferamount *big.Int, locksroot common.Hash, secret common.Hash) *Secret {
	p := &Secret{
		Secret: secret,
	}
	p.Identifier = Identifier
	p.CmdID = SecretCmdID
	p.Nonce = nonce
	p.Channel = channel
	p.TransferAmount = new(big.Int).Set(transferamount)
	p.Locksroot = locksroot
	return p
}

//Pack is MessagePacker
func (s *Secret) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, s.CmdID) //only one byte.
	binary.Write(buf, binary.BigEndian, s.Identifier)
	buf.Write(s.Secret[:])
	binary.Write(buf, binary.BigEndian, s.Nonce)
	buf.Write(s.Channel[:])
	buf.Write(utils.BigIntTo32Bytes(s.TransferAmount))
	buf.Write(s.Locksroot[:])
	buf.Write(s.Signature)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (s *Secret) UnPack(data []byte) error {
	var t int32
	s.CmdID = SecretCmdID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != s.CmdID {
		return fmt.Errorf("Ack Secret cmdid should be  4,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &s.Identifier)
	buf.Read(s.Secret[:])

	binary.Read(buf, binary.BigEndian, &s.Nonce)
	buf.Read(s.Channel[:])
	s.TransferAmount = readBigInt(buf)
	buf.Read(s.Locksroot[:])
	s.Signature = make([]byte, signatureLength)
	n, err := buf.Read(s.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errors.New("packet length error")
	}
	if err := s.checkValid(); err != nil {
		return err
	}
	return s.EnvelopMessage.verifySignature(data)
}

//String is fmt.Stringer
func (s *Secret) String() string {
	return fmt.Sprintf("Message{type=Secret secret=%s,%s}", utils.HPex(s.Secret), s.EnvelopMessage.String())
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
	HashLock common.Hash
}

//NewRemoveExpiredHashlockTransfer create  RemoveExpiredHashlockTransfer
func NewRemoveExpiredHashlockTransfer(Identifier uint64, nonce int64, channel common.Address,
	transferamount *big.Int, locksroot common.Hash, hashlock common.Hash) *RemoveExpiredHashlockTransfer {
	p := &RemoveExpiredHashlockTransfer{
		HashLock: hashlock,
	}
	if Identifier != 0 {
		panic("identifier is useless")
	}
	p.Identifier = Identifier
	p.CmdID = RemoveExpiredHashLockCmdID
	p.Nonce = nonce
	p.Channel = channel
	p.TransferAmount = new(big.Int).Set(transferamount)
	p.Locksroot = locksroot
	return p
}

//Pack is MessagePacker
func (reht *RemoveExpiredHashlockTransfer) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, reht.CmdID) //only one byte.
	binary.Write(buf, binary.BigEndian, reht.Identifier)
	buf.Write(reht.HashLock[:])
	binary.Write(buf, binary.BigEndian, reht.Nonce)
	buf.Write(reht.Channel[:])
	buf.Write(utils.BigIntTo32Bytes(reht.TransferAmount))
	buf.Write(reht.Locksroot[:])
	buf.Write(reht.Signature)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (reht *RemoveExpiredHashlockTransfer) UnPack(data []byte) error {
	var t int32
	reht.CmdID = RemoveExpiredHashLockCmdID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != reht.CmdID {
		return fmt.Errorf("Ack Secret cmdid should be  4,but get %d", t)
	}
	binary.Read(buf, binary.BigEndian, &reht.Identifier)
	if reht.Identifier != 0 {
		panic("identifier should be 0")
	}
	buf.Read(reht.HashLock[:])

	binary.Read(buf, binary.BigEndian, &reht.Nonce)
	buf.Read(reht.Channel[:])
	reht.TransferAmount = readBigInt(buf)
	buf.Read(reht.Locksroot[:])
	reht.Signature = make([]byte, signatureLength)
	n, err := buf.Read(reht.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errors.New("packet length error")
	}
	if err := reht.checkValid(); err != nil {
		return err
	}
	return reht.EnvelopMessage.verifySignature(data)
}

//String is fmt.Stringer
func (reht *RemoveExpiredHashlockTransfer) String() string {
	return fmt.Sprintf("Message{type=RemoveExpiredHashlockTransfer secret=%s,%s}", utils.HPex(reht.HashLock), reht.EnvelopMessage.String())
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
    recipient: The address of the raiden node participating in the channel.
    locksroot: The root of a merkle tree which records the current
        outstanding locks.
*/
type DirectTransfer struct {
	EnvelopMessage
	Token     common.Address //20bytes
	Recipient common.Address //20bytes
}

//String is fmt.Stringer
func (dt *DirectTransfer) String() string {
	return fmt.Sprintf("Message{type=DirectTransfer token=%s,recipient=%s,%s}", utils.APex2(dt.Token),
		utils.APex2(dt.Recipient), dt.EnvelopMessage.String())
}

//NewDirectTransfer create DirectTransfer
func NewDirectTransfer(identifier uint64, nonce int64, token common.Address,
	channel common.Address, transferAmount *big.Int,
	recipient common.Address, locksroot common.Hash) *DirectTransfer {
	p := &DirectTransfer{
		Token:     token,
		Recipient: recipient,
	}
	p.Identifier = identifier
	p.CmdID = DirectTransferCmdID
	p.Nonce = nonce
	p.Channel = channel
	p.TransferAmount = new(big.Int).Set(transferAmount)
	p.Locksroot = locksroot
	return p
}

//Pack is MessagePacker
func (dt *DirectTransfer) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, dt.CmdID) //only one byte.
	binary.Write(buf, binary.BigEndian, dt.Nonce)
	binary.Write(buf, binary.BigEndian, dt.Identifier)
	buf.Write(dt.Token[:])
	buf.Write(dt.Channel[:])
	buf.Write(dt.Recipient[:])
	buf.Write(utils.BigIntTo32Bytes(dt.TransferAmount))
	buf.Write(dt.Locksroot[:]) //todo locksroot pack unpack maybe error
	buf.Write(dt.Signature)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (dt *DirectTransfer) UnPack(data []byte) error {
	var t int32
	dt.CmdID = DirectTransferCmdID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	if t != dt.CmdID {
		return errors.New("DirectTransfer unpack cmdid error")
	}
	binary.Read(buf, binary.BigEndian, &dt.Nonce)
	binary.Read(buf, binary.BigEndian, &dt.Identifier)
	buf.Read(dt.Token[:])
	buf.Read(dt.Channel[:])
	buf.Read(dt.Recipient[:])
	dt.TransferAmount = readBigInt(buf)
	buf.Read(dt.Locksroot[:])
	dt.Signature = make([]byte, signatureLength)
	n, err := buf.Read(dt.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errPacketLength
	}
	if err := dt.checkValid(); err != nil {
		return err
	}
	return dt.EnvelopMessage.verifySignature(data)
}

//Lock of HTLC
type Lock struct {
	Expiration int64 //expiration block number
	Amount     *big.Int
	HashLock   common.Hash
}

//AsBytes serialize Lock
func (l *Lock) AsBytes() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, l.Expiration)
	buf.Write(utils.BigIntTo32Bytes(l.Amount))
	buf.Write(l.HashLock[:])
	return buf.Bytes()
}

//FromBytes deserialize Lock
func (l *Lock) FromBytes(locksencoded []byte) {
	buf := bytes.NewBuffer(locksencoded)
	binary.Read(buf, binary.BigEndian, &l.Expiration)
	l.Amount = readBigInt(buf)
	buf.Read(l.HashLock[:])
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
	Expiration int64
	Token      common.Address
	Recipient  common.Address
	Target     common.Address
	Initiator  common.Address
	HashLock   common.Hash
	Amount     *big.Int //The number transferred to party
	Fee        *big.Int
}

//String is fmt.Stringer
func (mt *MediatedTransfer) String() string {
	return fmt.Sprintf("Message{type=MediatedTransfer expiration=%d,token=%s,recipient=%s,target=%s,initiator=%s,hashlock=%s,amount=%s,fee=%s,%s}",
		mt.Expiration, utils.APex2(mt.Token), utils.APex2(mt.Recipient), utils.APex2(mt.Target), utils.APex2(mt.Initiator),
		utils.HPex(mt.HashLock), mt.Amount, mt.Fee, mt.EnvelopMessage.String())
}

//NewMediatedTransfer create MediatedTransfer
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
	p.Locksroot = locksroot //Including the merkletree root of the incomplete  transaction
	p.CmdID = MediatedTransferCmdID
	p.Channel = channel
	p.Expiration = lock.Expiration
	p.HashLock = lock.HashLock
	p.Amount = new(big.Int).Set(lock.Amount)
	return p
}

//GetLock returns Lock of this Transfer
func (mt *MediatedTransfer) GetLock() *Lock {
	return &Lock{
		Expiration: mt.Expiration,
		Amount:     mt.Amount,
		HashLock:   mt.HashLock,
	}
}

//Pack is MessagePacker
func (mt *MediatedTransfer) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, mt.CmdID) //one byte
	binary.Write(buf, binary.BigEndian, mt.Nonce)
	binary.Write(buf, binary.BigEndian, mt.Identifier)
	binary.Write(buf, binary.BigEndian, mt.Expiration)
	buf.Write(mt.Token[:])
	buf.Write(mt.Channel[:])
	buf.Write(mt.Recipient[:])
	buf.Write(mt.Target[:])
	buf.Write(mt.Initiator[:])
	buf.Write(mt.Locksroot[:])
	buf.Write(mt.HashLock[:])
	buf.Write(utils.BigIntTo32Bytes(mt.TransferAmount))
	buf.Write(utils.BigIntTo32Bytes(mt.Amount))
	buf.Write(utils.BigIntTo32Bytes(mt.Fee))
	buf.Write(mt.Signature)
	return buf.Bytes()
}

//UnPack is MessageUnPacker
func (mt *MediatedTransfer) UnPack(data []byte) error {
	var t int32
	//mt.CmdID = MEDIATEDTRANSFER_CMDID
	buf := bytes.NewBuffer(data)
	binary.Read(buf, binary.LittleEndian, &t)
	mt.CmdID = t
	if mt.CmdID != MediatedTransferCmdID && mt.CmdID != RefundTransferCmdID {
		return errors.New("MediatedTransfer unpack cmd error")
	}
	binary.Read(buf, binary.BigEndian, &mt.Nonce)
	binary.Read(buf, binary.BigEndian, &mt.Identifier)
	binary.Read(buf, binary.BigEndian, &mt.Expiration)
	buf.Read(mt.Token[:])
	buf.Read(mt.Channel[:])
	buf.Read(mt.Recipient[:])
	buf.Read(mt.Target[:])
	buf.Read(mt.Initiator[:])
	buf.Read(mt.Locksroot[:])
	buf.Read(mt.HashLock[:])
	mt.TransferAmount = readBigInt(buf)
	mt.Amount = readBigInt(buf)
	mt.Fee = readBigInt(buf)
	mt.Signature = make([]byte, signatureLength)
	n, err := buf.Read(mt.Signature)
	if err != nil {
		return err
	}
	if n != signatureLength {
		return errPacketLength
	}
	if err := mt.checkValid(); err != nil {
		return err
	}
	if mt.Expiration <= 0 {
		return fmt.Errorf("expiration must be positive %d", mt.Expiration)
	}
	if !utils.IsValidPositiveInt256(mt.Amount) {
		return fmt.Errorf("amount must be positive %s", mt.Amount)
	}
	return mt.verifySignature(data)
}

//RefundTransfer is used for a mediated node who cannot find any next node to send the mediatedTransfer
type RefundTransfer struct {
	MediatedTransfer
}

//String is fmt.Stringer
func (rt *RefundTransfer) String() string {
	return fmt.Sprintf("Message{type=RefundTransfer expiration=%d,token=%s,recipient=%s,target=%s,initiator=%s,hashlock=%s,amount=%s,fee=%s,%s}",
		rt.Expiration, utils.APex2(rt.Token), utils.APex2(rt.Recipient), utils.APex2(rt.Target), utils.APex2(rt.Initiator),
		utils.HPex(rt.HashLock), rt.Amount, rt.Fee, rt.EnvelopMessage.String())
}

//NewRefundTransfer create RefundTransfer
func NewRefundTransfer(identifier uint64, nonce int64, token common.Address,
	channel common.Address, transferAmount *big.Int,
	recipient common.Address, locksroot common.Hash, lock *Lock,
	target common.Address, initiator common.Address, fee *big.Int) *RefundTransfer {
	p := &RefundTransfer{}
	p.MediatedTransfer = *(NewMediatedTransfer(identifier, nonce, token, channel, transferAmount, recipient,
		locksroot, lock, target, initiator, fee))
	p.CmdID = RefundTransferCmdID
	return p
}

//NewRefundTransferFromMediatedTransfer create  RefundTransfer from a MediatedTransfer
func NewRefundTransferFromMediatedTransfer(mtr *MediatedTransfer) *RefundTransfer {
	p := &RefundTransfer{}
	p.MediatedTransfer = *mtr
	p.CmdID = RefundTransferCmdID
	return p
}

//IsLockedTransfer return true when this message is a RefundTransfer or MediatedTransfer
func IsLockedTransfer(msg Messager) bool {
	return msg.Cmd() == RefundTransferCmdID || msg.Cmd() == MediatedTransferCmdID
}

//GetMtrFromLockedTransfer returns the MediatedTransfer ,the caller must maker sure this message is a  locked transfer
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

//MessageMap contains all message can send and receive.
//DirectTransfer has been deprecated
var MessageMap = map[int]Messager{
	PingCmdID:                  new(Ping),
	AckCmdID:                   new(Ack),
	SecretRequestCmdID:         new(SecretRequest),
	SecretCmdID:                new(Secret),
	DirectTransferCmdID:        new(DirectTransfer),
	RevealSecretCmdID:          new(RevealSecret),
	MediatedTransferCmdID:      new(MediatedTransfer),
	RefundTransferCmdID:        new(RefundTransfer),
	RemoveExpiredHashLockCmdID: new(RemoveExpiredHashlockTransfer),
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
	gob.Register(&RemoveExpiredHashlockTransfer{})
}
