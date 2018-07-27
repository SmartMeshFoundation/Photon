package xmpppass

import (
	"crypto/ecdsa"

	"time"

	"encoding/hex"

	"errors"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ethereum/go-ethereum/crypto"
)

//#nosec
const passwordFormat = "2006-01-02"

//CreatePassword is helper function for login to xmpp server
func CreatePassword(privKey *ecdsa.PrivateKey) (sig string, err error) {
	t := time.Now().UTC()
	data := []byte(t.Format(passwordFormat))
	hash := crypto.Keccak256Hash(data)
	signature, err := crypto.Sign(hash[:], privKey)
	if err == nil {
		sig = hex.EncodeToString(signature)
	}
	return
}

//VerifySignature verify user,password is right or not
func VerifySignature(addr, signature string) (err error) {
	t := time.Now().UTC()
	data := []byte(t.Format(passwordFormat))
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
