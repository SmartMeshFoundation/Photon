package rpc

import (
	"github.com/SmartMeshFoundation/raiden-network/params"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/node"
)

const key = `
{
  "address": "1a9ec3b0b807464e6d3398a59d6b0a369bf422fa",
  "crypto": {
    "cipher": "aes-128-ctr",
    "ciphertext": "a471054846fb03e3e271339204420806334d1f09d6da40605a1a152e0d8e35f3",
    "cipherparams": {
      "iv": "44c5095dc698392c55a65aae46e0b5d9"
    },
    "kdf": "scrypt",
    "kdfparams": {
      "dklen": 32,
      "n": 262144,
      "p": 1,
      "r": 8,
      "salt": "e0a5fbaecaa3e75e20bccf61ee175141f3597d3b1bae6a28fe09f3507e63545e"
    },
    "mac": "cb3f62975cf6e7dfb454c2973bdd4a59f87262956d5534cdc87fb35703364043"
  },
  "id": "e08301fb-a263-4643-9c2b-d28959f66d6a",
  "version": 3
}
`

func MakeTestBlockChainService() *BlockChainService {
	conn, err := ethclient.Dial(node.DefaultIPCEndpoint("geth"))
	if err != nil {
		log.Error("Failed to connect to the Ethereum client: ", err)
	}
	privkey, err := keystore.DecryptKey([]byte(key), "123")
	if err != nil {
		log.Error("Failed to create authorized transactor: ", err)
	}
	return NewBlockChainService(privkey.PrivateKey, params.ROPSTEN_REGISTRY_ADDRESS, conn)
}
