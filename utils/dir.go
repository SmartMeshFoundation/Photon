package utils

import (
	"encoding/hex"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/node"
)

//DefaultDataDir default work directory
func DefaultDataDir() string {
	// Try to place the data folder in the user's home dir
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, "Library", "photon")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "photon")
		} else {
			return filepath.Join(home, ".photon")
		}
	}
	// As we cannot guess a stable location, return empty and handle later
	return ""
}

/*
BuildPhotonDbPath 根据datadir和运行photon的账号构建数据库目录
*/
func BuildPhotonDbPath(dataDir string, myAddress common.Address) (dbPath string, err error) {
	if !Exists(dataDir) {
		err = os.MkdirAll(dataDir, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", dataDir, err)
			return
		}
	}
	userDbPath := hex.EncodeToString(myAddress[:])
	userDbPath = userDbPath[:8]
	userDbPath = filepath.Join(dataDir, userDbPath)
	if !Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", dataDir, err)
			return
		}
	}
	dbPath = filepath.Join(userDbPath, "log.db")
	return
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}
	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}
	return ""
}

//DefaultKeyStoreDir keystore path of ethereum
func DefaultKeyStoreDir() string {
	return filepath.Join(node.DefaultDataDir(), "keystore")
}
