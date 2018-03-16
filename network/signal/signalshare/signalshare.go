package signalshare

import (
	"crypto/ecdsa"

	"time"

	"encoding/hex"

	"github.com/SmartMeshFoundation/raiden-network/utils"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-errors/errors"
)

const PasswordFormat = "2006-01-02"

func CreatePassword(privKey *ecdsa.PrivateKey) (sig string, err error) {
	t := time.Now()
	data := []byte(t.Format(PasswordFormat))
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash[:], privKey)
	if err == nil {
		sig = hex.EncodeToString(signature)
	}
	return
}

func VerifySignature(addr, signature string) (err error) {
	t := time.Now()
	data := []byte(t.Format(PasswordFormat))
	sig, err := hex.DecodeString(signature)
	if err != nil {
		return err
	}
	hash := crypto.Keccak256Hash(data)
	pubkey, err := crypto.Ecrecover(hash[:], sig)
	if err != nil {
		return
	}
	sender := utils.PubkeyToAddress(pubkey)
	if addr != sender.String() {
		return errors.New("not match")
	}
	return nil
}
