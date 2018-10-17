package mainimpl

import (
	"fmt"
	"os"
	"testing"
	"time"

	"path/filepath"

	smartraiden "github.com/SmartMeshFoundation/SmartRaiden"
	"github.com/SmartMeshFoundation/SmartRaiden/accounts"
	"github.com/SmartMeshFoundation/SmartRaiden/codefortest"
	"github.com/SmartMeshFoundation/SmartRaiden/params"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/stretchr/testify/assert"
)

func TestPromptAccount(t *testing.T) {
	accounts.PromptAccount(utils.EmptyAddress, `../../../testdata/keystore`, "123")
}
func panicOnNullValue() {
	var c []int
	c[0] = 0
}

func TestPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			//t.Error(err)
		} else {
			t.Error("should panic")
		}
	}()
	panicOnNullValue()
}

type T struct {
	a int
}

func TestStruct(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			//t.Error(err)
		} else {
			t.Error("should panic")
		}
	}()
	var a *T
	t.Logf("a.a=%d", a.a)
}

func TestStart(t *testing.T) {
	os.Args = make([]string, 0, 20)
	os.Args = append(os.Args, "smartraiden")
	os.Args = append(os.Args, fmt.Sprintf("--address=%s", "0x1a9ec3b0b807464e6d3398a59d6b0a369bf422fa"))
	os.Args = append(os.Args, fmt.Sprintf("--keystore-path=%s", "../../../testdata/keystore"))
	os.Args = append(os.Args, fmt.Sprintf("--eth-rpc-endpoint=%s", os.Getenv("ETHRPCENDPOINT")))
	os.Args = append(os.Args, fmt.Sprintf("--datadir=%s", ".smartraiden"))
	os.Args = append(os.Args, fmt.Sprintf("--password-file=%s", "../../../testdata/keystore/pass"))
	os.Args = append(os.Args, fmt.Sprintf("--api-address=%s", "127.0.0.1:2000"))
	os.Args = append(os.Args, fmt.Sprintf("--listen-address=%s", "127.0.0.1:20000"))
	os.Args = append(os.Args, fmt.Sprintf("--verbosity=5"))
	os.Args = append(os.Args, fmt.Sprintf("--debug"))
	params.MobileMode = true

	var api *smartraiden.RaidenAPI
	var err error
	// 1. 无公链第一次启动,must fail
	clearData(".smartraiden")
	os.Args[3] = fmt.Sprintf("--eth-rpc-endpoint=%s", "ws://127.0.0.1:9999")
	api, err = StartMain()
	assert.Error(t, err)
	time.Sleep(5 * time.Second)
	api = nil
	err = nil
	// 2. 有公链第一次启动,must success
	clearData(".smartraiden")
	os.Args[3] = fmt.Sprintf("--eth-rpc-endpoint=%s", os.Getenv("ETHRPCENDPOINT"))
	api, err = StartMain()
	assert.Empty(t, err)
	time.Sleep(5 * time.Second)
	api.Stop()
	// 3. 无公链非第一次启动,must success
	os.Args[3] = fmt.Sprintf("--eth-rpc-endpoint=%s", "ws://127.0.0.1:9999")
	api, err = StartMain()
	assert.Empty(t, err)
	time.Sleep(5 * time.Second)
	api.Stop()
	// 4. matrix启动, must success
	os.Args[3] = fmt.Sprintf("--eth-rpc-endpoint=%s", os.Getenv("ETHRPCENDPOINT"))
	os.Args = append(os.Args, fmt.Sprintf("--matrix"))
	os.Args = append(os.Args, fmt.Sprintf("--matrix-server=%s", "transport01.smartmesh.cn"))
	api, err = StartMain()
	assert.Empty(t, err)
	time.Sleep(5 * time.Second)
	api.Stop()
	// 5. nonetwork启动, must success
	os.Args[3] = fmt.Sprintf("--eth-rpc-endpoint=%s", os.Getenv("ETHRPCENDPOINT"))
	os.Args[len(os.Args)-2] = fmt.Sprintf("--nonetwork")
	os.Args[len(os.Args)-1] = fmt.Sprintf("")
	api, err = StartMain()
	assert.Empty(t, err)
	time.Sleep(5 * time.Second)
	api.Stop()
}

func TestMeshBoxStart(t *testing.T) {
	if os.Getenv("IS_MESH_BOX") != "true" {
		return
	}
	os.Args = make([]string, 0, 20)
	os.Args = append(os.Args, "smartraiden")
	os.Args = append(os.Args, fmt.Sprintf("--eth-rpc-endpoint=%s", os.Getenv("ETHRPCENDPOINT")))
	os.Args = append(os.Args, fmt.Sprintf("--datadir=%s", ".smartraiden"))
	os.Args = append(os.Args, fmt.Sprintf("--api-address=%s", "127.0.0.1:2000"))
	os.Args = append(os.Args, fmt.Sprintf("--listen-address=%s", "127.0.0.1:20000"))
	os.Args = append(os.Args, fmt.Sprintf("--verbosity=5"))
	os.Args = append(os.Args, fmt.Sprintf("--debug"))
	params.MobileMode = true

	var api *smartraiden.RaidenAPI
	var err error
	// 1. 无公链第一次启动,must fail
	clearData(".smartraiden")
	api, err = StartMain()
	if err != nil {
		panic(err)
	}
	time.Sleep(5 * time.Second)
	api.Stop()
}

func clearData(dataPath string) {
	filepath.Walk(dataPath, func(path string, fi os.FileInfo, err error) error {
		if nil == fi {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		name := fi.Name()

		if name == ".smartraiden" {
			err := os.RemoveAll(path)
			if err != nil {
				fmt.Println("delet dir error:", err)
			}
		}
		return nil
	})
}

func TestVerifyContractCode(t *testing.T) {
	client, err := codefortest.GetEthClient()
	if err != nil {
		t.Error(err.Error())
	}
	registryAddress, _, _, err := codefortest.DeployRegistryContract()
	if err != nil {
		t.Error(err.Error())
	}
	err = verifyContractCode(registryAddress, client)
	if err != nil {
		t.Error(err.Error())
	}
}

func TestChan(t *testing.T) {
	c := make(chan int)
	ok := false
	select {
	case _, ok2 := <-c:
		ok = ok2
	default:
	}
	if ok {
		close(c)
	} else {
		fmt.Println("already close")
	}
}
