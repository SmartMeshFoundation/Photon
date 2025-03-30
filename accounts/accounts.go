package accounts

import (
	"bytes"
	"fmt"
	"path/filepath"

	"io/ioutil"

	"strings"

	"errors"

	"os"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/howeyc/gopass"
)

var errNoSuchAddress = errors.New("can not found this address")

/*
AccountManager List All Accounts in directory KeyPath
*/
type AccountManager struct {
	KeyPath  string
	Accounts []accounts.Account
}

// NewAccountManager create account manager
func NewAccountManager(keyPath string) (mgr *AccountManager) {
	mgr = &AccountManager{
		KeyPath: keyPath,
	}
	ks := keystore.NewKeyStore(keyPath, keystore.StandardScryptN, keystore.StandardScryptP)
	mgr.Accounts = ks.Accounts()
	ks.Close()
	return
}

// AddressInKeyStore returns true if found this address
func (am *AccountManager) AddressInKeyStore(addr common.Address) bool {
	for _, acc := range am.Accounts {
		if bytes.Equal(acc.Address[:], addr[:]) {
			return true
		}
	}
	return false
}

/*
GetPrivateKey Find the keystore file for an account, unlock it and get the private key
   	addr: The Ethereum address for which to find the keyfile in the system
	password: Mostly for testing purposes. A password can be provided
			  as the function argument here. If it's not then the
              user is interactively queried for one.
    return The private key associated with the address
*/
func (am *AccountManager) GetPrivateKey(addr common.Address, password string) (privKeyBin []byte, err error) {
	if !am.AddressInKeyStore(addr) {
		err = errNoSuchAddress
		return
	}
	addrhex := strings.ToLower(addr.Hex())
	filename := fmt.Sprintf("UTC--*%s", addrhex[2:]) //skip 0x
	path := filepath.Join(am.KeyPath, filename)
	files, err := filepath.Glob(path)
	if err != nil {
		return
	}
	if len(files) != 1 {
		err = fmt.Errorf("cannot find %s's key file", addr.String())
		return
	}
	keyjson, err := ioutil.ReadFile(files[0])
	if err != nil {
		return
	}
	key, err := keystore.DecryptKey(keyjson, password)
	if err != nil {
		return
	}
	privKeyBin = crypto.FromECDSA(key.PrivateKey)
	return
}

// PromptAccount get account private key by input password or password stored in file
func PromptAccount(adviceAddress common.Address, keystorePath, passwordfile string) (addr common.Address, keybin []byte, err error) {
	am := NewAccountManager(keystorePath)
	if len(am.Accounts) == 0 {
		err = fmt.Errorf("No Ethereum accounts found in the directory %s", keystorePath)
		return
	}
	if !am.AddressInKeyStore(adviceAddress) {
		if adviceAddress != utils.EmptyAddress {
			err = fmt.Errorf("account %s could not be found on the sytstem. aborting", adviceAddress.String())
			return
		}
		shouldPromt := true
		fmt.Println("The following accounts were found in your machine:")
		for i := 0; i < len(am.Accounts); i++ {
			fmt.Printf("%3d -  %s\n", i, am.Accounts[i].Address.String())
		}
		fmt.Println("")
		for shouldPromt {
			fmt.Printf("Select one of them by index to continue:\n")
			idx := -1
			_, err = fmt.Scanf("%d", &idx)
			if err != nil {
				return
			}
			if idx >= 0 && idx < len(am.Accounts) {
				shouldPromt = false
				addr = am.Accounts[idx].Address
			} else {
				fmt.Printf("Error: Provided index %d is out of bounds", idx)
			}
		}
	} else {
		addr = adviceAddress
	}
	if len(passwordfile) > 0 {
		var data []byte
		//#nosec
		data, err = ioutil.ReadFile(passwordfile)
		if err != nil {
			data = []byte(passwordfile)
		}
		password := string(data)
		//log.Trace(fmt.Sprintf("password is %s", password))
		keybin, err = am.GetPrivateKey(addr, password)
		if err != nil {
			err = fmt.Errorf("Incorrect password for %s in file. Aborting ... %s", addr.String(), err)
			return
		}
	} else {
		for i := 0; i < 3; i++ {
			var pb []byte
			//retries three times
			pb, err = gopass.GetPasswdPrompt("Enter the password to unlock:", false, os.Stdin, os.Stdout)
			if err != nil {
				return
			}
			password := string(pb) // getpass.Prompt("Enter the password to unlock:")
			keybin, err = am.GetPrivateKey(addr, password)
			if err != nil && i == 3 {
				log.Error(fmt.Sprintf("Exhausted passphrase unlock attempts for %s. Aborting ...", addr))
				return
			}
			if err != nil {
				log.Error(fmt.Sprintf("password incorrect\n Please try again or kill the process to quit.\nUsually Ctrl-c."))
				continue
			}
			break
		}
	}
	return
}
