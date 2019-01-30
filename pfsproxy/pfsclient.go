package pfsproxy

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"math/big"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
)

// ErrNotInit :
var ErrNotInit = errors.New("pfgClient not init")

/*
pfsClient :
*/
type pfsClient struct {
	host       string
	privateKey *ecdsa.PrivateKey
}

/*
NewPfsProxy :
*/
func NewPfsProxy(pfgHost string, privateKey *ecdsa.PrivateKey) (pfsProxy PfsProxy) {
	pfsProxy = &pfsClient{
		host:       pfgHost,
		privateKey: privateKey,
	}
	return
}

/*
example :
{
    "balance_proof": {
        "nonce": 8,
        "transfer_amount": 320,
        "locks_root": "0x0000000000000000000000000000000000000000000000000000000000000000",
        "channel_identifier": "0x0398beea63f098e2d3bb59884be79eda00cf042e39ad65e5c43a0a280f969f93",
        "open_block_number": 7228470,
        "addition_hash": "0x3189fad45065c5e505180de926dbf496ad8213e6137a711f72609c6241959718",
        "signature": "ovWIT4r48tXrKFeLK2WA93VRcciyIbi7rvycL5R9wxpl7ZQgOgU0QiK+BPFQDJPgHxTd5Lpjrf8mXLPa2fTtEhw="
    },
    "balance_signature": "yXKPCkGzvRsFrg51NXsxYZ1xkCRWOgWNUdxUkHGDJwQT0g0LKAN7tt7fzN9y1+5sKYWTSfs5zOSngO0SvjSxRRs=",
    "lock_amount": 0
}
*/
type submitBalancePayload struct {
	BalanceProof     *balanceProof `json:"balance_proof"`
	BalanceSignature []byte        `json:"balance_signature"`
	LockAmount       *big.Int      `json:"lock_amount"`
}

type balanceProof struct {
	Nonce             uint64      `json:"nonce"`
	TransferAmount    *big.Int    `json:"transfer_amount"`
	Locksroot         common.Hash `json:"locks_root"`
	ChannelIdentifier common.Hash `json:"channel_identifier"`
	OpenBlockNumber   int64       `json:"open_block_number"`
	AdditionHash      common.Hash `json:"addition_hash"`
	Signature         []byte      `json:"signature"`
}

func (p *submitBalancePayload) sign(key *ecdsa.PrivateKey) []byte {
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, p.BalanceProof.Nonce)
	_, err = buf.Write(utils.BigIntTo32Bytes(p.BalanceProof.TransferAmount))
	_, err = buf.Write(p.BalanceProof.Locksroot[:])
	_, err = buf.Write(p.BalanceProof.ChannelIdentifier[:])
	err = binary.Write(buf, binary.BigEndian, p.BalanceProof.OpenBlockNumber)
	_, err = buf.Write(p.BalanceProof.AdditionHash[:])
	_, err = buf.Write(p.BalanceProof.Signature)
	_, err = buf.Write(utils.BigIntTo32Bytes(p.LockAmount))
	if err != nil {
		log.Error(fmt.Sprintf("signData err %s", err))
	}
	p.BalanceSignature, err = utils.SignData(key, buf.Bytes())
	if err != nil {
		log.Crit(fmt.Sprintf("signDataFor submitBalancePayload err %s", err))
	}
	return p.BalanceSignature
}

/*
SubmitBalance :
*/
func (pfg *pfsClient) SubmitBalance(nonce uint64, transferAmount, lockAmount *big.Int, openBlockNumber int64, locksroot, channelIdentifier, additionHash common.Hash, signature []byte) (err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		return ErrNotInit
	}
	payload := &submitBalancePayload{
		BalanceProof: &balanceProof{
			Nonce:             nonce,
			TransferAmount:    transferAmount,
			Locksroot:         locksroot,
			ChannelIdentifier: channelIdentifier,
			OpenBlockNumber:   openBlockNumber,
			AdditionHash:      additionHash,
			Signature:         signature,
		},
		LockAmount: lockAmount,
	}
	payload.sign(pfg.privateKey)
	req := &req{
		FullURL: pfg.host + "/pfs/1/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String() + "/balance",
		Method:  http.MethodPut,
		Payload: marshal(payload),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI SubmitBalance %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI SubmitBalance %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	return nil
}

// findPathPayload :
type findPathPayload struct {
	PeerFrom          common.Address `json:"peer_from"`
	PeerTo            common.Address `json:"peer_to"`
	TokenAddress      common.Address `json:"token_address"`
	LimitPaths        int            `json:"limit_paths"`
	SendAmount        *big.Int       `json:"send_amount"`
	SortDemand        string         `json:"sort_demand"`
	Signature         []byte         `json:"signature"`
	PeerFromChargeFee bool           `json:"peer_from_charge_fee"`
}

func (p *findPathPayload) sign(key *ecdsa.PrivateKey) []byte {
	var err error
	buf := new(bytes.Buffer)
	_, err = buf.Write(p.PeerFrom[:])
	_, err = buf.Write(p.PeerTo[:])
	_, err = buf.Write(p.TokenAddress[:])
	err = binary.Write(buf, binary.BigEndian, p.LimitPaths)
	_, err = buf.Write(utils.BigIntTo32Bytes(p.SendAmount))
	_, err = buf.Write([]byte(p.SortDemand))
	if p.PeerFromChargeFee {
		_, err = buf.Write([]byte{byte(1)})
	} else {
		_, err = buf.Write([]byte{byte(0)})
	}

	if err != nil {
		log.Error(fmt.Sprintf("signData err %s", err))
	}
	p.Signature, err = utils.SignData(key, buf.Bytes())
	if err != nil {
		log.Crit(fmt.Sprintf("signDataFor FindPathPayload err %s", err))
	}
	return p.Signature
}

// FindPathResponse :
type FindPathResponse struct {
	PathID  int      `json:"path_id"`
	PathHop int      `json:"path_hop"`
	Fee     *big.Int `json:"fee"`
	Result  []string `json:"result"`
}

/*
FindPath : find path
*/
func (pfg *pfsClient) FindPath(peerFrom, peerTo, token common.Address, amount *big.Int, isInitiator bool) (resp []FindPathResponse, err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		err = ErrNotInit
		return
	}
	payload := &findPathPayload{
		PeerFrom:          peerFrom,
		PeerTo:            peerTo,
		TokenAddress:      token,
		LimitPaths:        1,
		SendAmount:        amount,
		SortDemand:        "",
		PeerFromChargeFee: !isInitiator,
	}
	payload.sign(pfg.privateKey)
	req := &req{
		FullURL: pfg.host + "/pfs/1/paths",
		Method:  http.MethodPost,
		Payload: marshal(payload),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI FindPath %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI FindPath %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}
	log.Trace(fmt.Sprintf("resp=%s", string(body)))
	return
}

// setFeePayload :
type setFeePayload struct {
	FeeConstant *big.Int `json:"fee_constant"`
	FeePercent  int64    `json:"fee_percent"`
	Signature   []byte   `json:"signature"`
}

func (p *setFeePayload) sign(key *ecdsa.PrivateKey) []byte {
	var err error
	buf := new(bytes.Buffer)
	err = binary.Write(buf, binary.BigEndian, p.FeePercent)
	_, err = buf.Write(utils.BigIntTo32Bytes(p.FeeConstant))
	if err != nil {
		log.Error(fmt.Sprintf("signData err %s", err))
	}
	p.Signature, err = utils.SignData(key, buf.Bytes())
	if err != nil {
		log.Crit(fmt.Sprintf("signDataFor SetFeeRatePayload err %s", err))
	}
	return p.Signature
}

// getFeeResponse :
type getFeeResponse struct {
	FeePolicy   int64    `json:"fee_policy"`
	FeeConstant *big.Int `json:"fee_constant"`
	FeePercent  int64    `json:"fee_percent"`
}

/*
SetFeePolicy :set fee rate by account
*/
func (pfg *pfsClient) SetFeePolicy(fp *models.FeePolicy) (err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		return ErrNotInit
	}
	fp.Sign(pfg.privateKey)
	req := &req{
		FullURL: pfg.host + "/pfs/1/feerate/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodPut,
		Payload: marshal(fp),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error(fmt.Sprintf("PfgAPI SetFeePolicy %s err :%s", req.FullURL, err))
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI SetFeePolicy %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	return nil
}

/*
SetAccountFeeRate :set fee rate by account
*/
func (pfg *pfsClient) SetAccountFee(feeConstant *big.Int, feePercent int64) (err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		return ErrNotInit
	}
	payload := &setFeePayload{
		FeeConstant: feeConstant,
		FeePercent:  feePercent,
	}
	payload.sign(pfg.privateKey)
	req := &req{
		FullURL: pfg.host + "/pfs/1/account_rate/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodPut,
		Payload: marshal(payload),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI SetAccountFeeRate %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI SetAccountFeeRate %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	return nil
}

/*
GetAccountFee : get fee rate by account
*/
func (pfg *pfsClient) GetAccountFee() (feeConstant *big.Int, feePercent int64, err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		err = ErrNotInit
		return
	}
	req := &req{
		FullURL: pfg.host + "/pfs/1/account_rate/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodGet,
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI GetAccountFee %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI GetAccountFee %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	var resp getFeeResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}
	return resp.FeeConstant, resp.FeePercent, nil
}

/*
SetTokenFee :set fee rate of a token
*/
func (pfg *pfsClient) SetTokenFee(feeConstant *big.Int, feePercent int64, tokenAddress common.Address) (err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		return ErrNotInit
	}
	payload := &setFeePayload{
		FeeConstant: feeConstant,
		FeePercent:  feePercent,
	}
	payload.sign(pfg.privateKey)
	req := &req{
		FullURL: pfg.host + "/pfs/1/token_rate/" + tokenAddress.String() + "/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodPut,
		Payload: marshal(payload),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI SetTokenFee %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI SetTokenFee %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	return nil
}

/*
GetTokenFee : get fee rate by token
*/
func (pfg *pfsClient) GetTokenFee(tokenAddress common.Address) (feeConstant *big.Int, feePercent int64, err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		err = ErrNotInit
		return
	}
	req := &req{
		FullURL: pfg.host + "/pfs/1/token_rate/" + tokenAddress.String() + "/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodGet,
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI GetTokenFee %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI GetTokenFee %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	var resp getFeeResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}
	return resp.FeeConstant, resp.FeePercent, nil
}

/*
SetChannelFee :set fee rate of a channel
*/
func (pfg *pfsClient) SetChannelFee(feeConstant *big.Int, feePercent int64, channelIdentifier common.Hash) (err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		return ErrNotInit
	}
	payload := &setFeePayload{
		FeeConstant: feeConstant,
		FeePercent:  feePercent,
	}
	payload.sign(pfg.privateKey)
	req := &req{
		FullURL: pfg.host + "/pfs/1/channel_rate/" + channelIdentifier.String() + "/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodPut,
		Payload: marshal(payload),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI SetChannelFee %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI SetChannelFee %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	return nil
}

/*
GetChannelFee : get fee rate by channel
*/
func (pfg *pfsClient) GetChannelFee(channelIdentifier common.Hash) (feeConstant *big.Int, feePercent int64, err error) {
	if pfg.host == "" || pfg.privateKey == nil {
		err = ErrNotInit
		return
	}
	req := &req{
		FullURL: pfg.host + "/pfs/1/channel_rate/" + channelIdentifier.String() + "/" + crypto.PubkeyToAddress(pfg.privateKey.PublicKey).String(),
		Method:  http.MethodGet,
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	log.Debug(req.ToString())
	if err != nil {
		log.Error("PfgAPI GetChannelFee %s err :%s", req.FullURL, err)
		return
	}
	if statusCode != 200 {
		err = fmt.Errorf("PfgAPI GetChannelFee %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		log.Error(err.Error())
		return
	}
	var resp getFeeResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		panic(err)
	}
	return resp.FeeConstant, resp.FeePercent, nil
}

func marshal(v interface{}) string {
	p, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(p)
}

func marshalIndent(v interface{}) string {
	p, err := json.MarshalIndent(v, "\t", "")
	if err != nil {
		panic(err)
	}
	return string(p)
}
