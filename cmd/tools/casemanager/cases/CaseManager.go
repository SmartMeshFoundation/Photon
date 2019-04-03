package cases

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"errors"

	"github.com/SmartMeshFoundation/Photon/cmd/tools/casemanager/models"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
)

func init() {
	log.Root().SetHandler(log.LvlFilterHandler(log.LvlTrace, utils.MyStreamHandler(os.Stderr)))
}

// ErrorSkip 用于跳过某个case时返回,方便统计
var ErrorSkip = errors.New("skip")

// CaseManager include env and cases
type CaseManager struct {
	Cases                 map[string]reflect.Value
	FailedCaseNames       []string
	IsAutoRun             bool
	UseMatrix             bool
	RunSlow               bool //是否运行哪些长时间运行的case
	EthEndPoint           string
	LowWaitSeconds        int
	MediumWaitSeconds     int
	HighMediumWaitSeconds int
	MDNSLifeTime          time.Duration //多久认为一个不进行mdns广播的节点下线了.
}

// NewCaseManager constructor
func NewCaseManager(isAutoRun bool, useMatrix bool, ethEndPoint string, runSlow bool) (caseManager *CaseManager) {
	var err error
	caseManager = new(CaseManager)
	caseManager.IsAutoRun = isAutoRun
	caseManager.UseMatrix = useMatrix
	caseManager.EthEndPoint = ethEndPoint
	caseManager.LowWaitSeconds = 10
	caseManager.MediumWaitSeconds = 50
	caseManager.HighMediumWaitSeconds = 300
	caseManager.RunSlow = runSlow
	caseManager.Cases = make(map[string]reflect.Value)
	//
	if useMatrix {
		caseManager.LowWaitSeconds = 10 + 100
		caseManager.MediumWaitSeconds = 50 + 160 //config for settle time
		caseManager.HighMediumWaitSeconds = 300 + 100
	}
	//会通过启动参数来指定修改mdns的间隔时间 params.DefaultMDNSKeepalive修改为1秒,params.DefaultMDNSQueryInterval修改为50ms
	caseManager.MDNSLifeTime = time.Second + time.Millisecond*50*2 //实际上是params.DefaultMDNSKeepalive + 2*params.DefaultMDNSQueryInterval
	// use reflect to load all cases
	_, err = fmt.Println("load cases...")
	vf := reflect.ValueOf(caseManager)
	vft := vf.Type()
	for i := 0; i < vf.NumMethod(); i++ {
		mName := vft.Method(i).Name
		if strings.Contains(mName, "Case") {
			_, err = fmt.Println("CaseName:", mName)
			caseManager.Cases[mName] = vf.Method(i)
		}
	}
	_, err = fmt.Printf("Total %d cases load success\n", len(caseManager.Cases))
	_ = err
	return
}

// RunAll run all
func (c *CaseManager) RunAll(skip string) {
	var err error
	_, err = fmt.Println("Run all cases...")
	// 排序
	// sort
	var keys []string
	for k := range c.Cases {
		keys = append(keys, k)
	}
	eachUsed := make(map[string]int64)
	sort.Strings(keys)
	errorMsg := ""
	success := 0
	total := len(keys)
	for _, k := range keys {
		s := time.Now()
		v := c.Cases[k]
		rs := v.Call(nil)
		if rs[0].Interface() == nil {
			success++
		} else {
			err := rs[0].Interface().(error)
			if err == ErrorSkip {
				total--
				fmt.Printf("%s SKIP \n", k)
				continue
			}
			if err == nil {
				fmt.Printf("%s SUCCESS\n", k)
			} else {
				errorMsg = fmt.Sprintf("%s FAILED!!!,err=%s\n", k, err)
				fmt.Println(errorMsg)
				c.FailedCaseNames = append(c.FailedCaseNames, k)
				if skip != "true" {
					break
				}
			}
		}
		eachUsed[k] = time.Now().Unix() - s.Unix()
	}
	_, err = fmt.Println("Casemanager Result:")
	_, err = fmt.Printf("Cases num : %d,successed=%d\n", total, success)
	_, err = fmt.Printf("Fail num : %d :\n", len(c.FailedCaseNames))
	for _, v := range c.FailedCaseNames {
		_, err = fmt.Println(v)
	}
	_, err = fmt.Printf("Time used: \n")
	for k, u := range eachUsed {
		fmt.Printf("%d seconds : %s\n", u, k)
	}
	_, err = fmt.Println("Pelease check log in ./log")
	if errorMsg != "" && skip != "true" {
		panic(errorMsg)
	}
	_ = err
}

// RunOne run one
func (c *CaseManager) RunOne(caseName string) {
	var err error
	if v, ok := c.Cases[caseName]; ok {
		s := time.Now().Unix()
		_, err = fmt.Println("----------------------------->Start to run case " + caseName + "...")
		rs := v.Call(nil)
		if rs[0].Interface() == nil {
			_, err = fmt.Printf("%s SUCCESS\n", caseName)
		} else {
			err := rs[0].Interface().(error)
			if err == nil {
				_, err = fmt.Printf("%s SUCCESS\n", caseName)
			} else {
				_, err = fmt.Printf("%s FAILED!!! err=%s\n", caseName, err)
			}
		}
		fmt.Printf("Time used : %d seconds\n", time.Now().Unix()-s)
	} else {
		_, err = fmt.Printf("%s doesn't exist !!! \n", caseName)
	}
	_, err = fmt.Println("Please check log in ./log")
	_ = err
}

// caseFail :
func (c *CaseManager) caseFail(caseName string) error {
	models.Logger.Println(caseName + " END ====> FAILED")
	return fmt.Errorf("Case [%s] FAILED", caseName)
}

// caseFail :
func (c *CaseManager) caseFailWithWrongChannelData(caseName string, channelName string) error {
	models.Logger.Println(channelName + " data wrong !!!")
	models.Logger.Println(caseName + " END ====> FAILED")
	return fmt.Errorf("Case [%s] FAILED", caseName)
}

func (c *CaseManager) logSeparatorLine(s string) {
	models.Logger.Println("===============================================>")
	models.Logger.Println(s)
	models.Logger.Println("===============================================>")
}

func (c *CaseManager) checkNodesStartComplete(nodes []*models.PhotonNode) bool {
	for i := 0; i < len(nodes); i++ {
		if !nodes[i].IsRunning() {
			return false
		}
	}
	return true
}

func (c *CaseManager) startNodes(env *models.TestEnv, nodes ...*models.PhotonNode) {
	n := len(nodes)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(index int) {
			if index == 0 {
				nodes[index].SetDoPprof()
			}
			nodes[index].Start(env)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(c.MDNSLifeTime)
}

func (c *CaseManager) startNodesWithFee(env *models.TestEnv, nodes ...*models.PhotonNode) {
	n := len(nodes)
	wg := sync.WaitGroup{}
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(index int) {
			nodes[index].StartWithFeeAndPFS(env)
			wg.Done()
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second)
}

type repeatReturnNilSuccessFunc func() error

/*
在seconds秒内,如果f返回nil直接返回,
否则一直尝试执行f,
如果超过seconds秒则返回失败
*/
func (c *CaseManager) tryInSeconds(seconds int, f repeatReturnNilSuccessFunc) error {
	var err error
	var i = 0
	for i = 0; i < seconds; i++ {
		time.Sleep(time.Second)
		err = f()
		if err == nil {
			break
		}
	}
	if i == seconds {
		return err
	}
	return nil
}

//在seconds秒内结算通道
func (c *CaseManager) trySettleInSeconds(seconds int, node *models.PhotonNode, channelIdentifier string) error {
	needsettle := true
	return c.tryInSeconds(seconds, func() error {
		if needsettle {
			err := node.Settle(channelIdentifier)
			if err == nil { //只要error不为空,就表示settle没有成功
				needsettle = false
				err = errors.New("wait settled")
			}
			return err
		}
		//进入等待交易被打包状态
		_, err := node.SpecifiedChannel(channelIdentifier)
		if err != nil {
			return nil //这里应该检测结果,确定是channel不存在,这里简化一下
		}
		return errors.New("retry")
	})
}

func (c *CaseManager) nodesExcept(nodes []*models.PhotonNode, n *models.PhotonNode) []*models.PhotonNode {
	r := make([]*models.PhotonNode, 0, len(nodes))
	for _, n2 := range nodes {
		if n2 == n {
			continue
		}
		r = append(r, n2)
	}
	return r
}
