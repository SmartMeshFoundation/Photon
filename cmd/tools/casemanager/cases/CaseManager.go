package cases

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/casemanager/models"
)

// CaseManager include env and cases
type CaseManager struct {
	Cases           map[string]reflect.Value
	FailedCaseNames []string
}

// NewCaseManager constructor
func NewCaseManager() (caseManager *CaseManager) {
	caseManager = new(CaseManager)
	caseManager.Cases = make(map[string]reflect.Value)
	// use reflect to load all cases
	fmt.Println("load cases...")
	vf := reflect.ValueOf(caseManager)
	vft := vf.Type()
	for i := 0; i < vf.NumMethod(); i++ {
		mName := vft.Method(i).Name
		if strings.Contains(mName, "Case") {
			fmt.Println("CaseName:", mName)
			caseManager.Cases[mName] = vf.Method(i)
		}
	}
	fmt.Printf("Total %d cases load success\n", len(caseManager.Cases))
	return
}

// RunAll run all
func (c *CaseManager) RunAll(skip string) {
	fmt.Println("Run all cases...")
	// 排序
	// sort
	var keys []string
	for k := range c.Cases {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	errorMsg := ""
	for _, k := range keys {
		v := c.Cases[k]
		rs := v.Call(nil)
		fmt.Println("----------------------------->Start to run case " + k + "...")
		if rs[0].Interface() == nil {
			fmt.Printf("%s SUCCESS\n", k)
		} else {
			err := rs[0].Interface().(error)
			if err == nil {
				fmt.Printf("%s SUCCESS\n", k)
			} else {
				errorMsg = fmt.Sprintf("%s FAILED!!!\n", k)
				fmt.Printf(errorMsg)
				c.FailedCaseNames = append(c.FailedCaseNames, k)
				if skip != "true" {
					break
				}
			}
		}
	}
	fmt.Println("Casemanager Result:")
	fmt.Printf("Cases num : %d\n", len(keys))
	fmt.Printf("Fail num : %d :\n", len(c.FailedCaseNames))
	for _, v := range c.FailedCaseNames {
		fmt.Println(v)
	}
	fmt.Println("Pelease check log in ./log")
	if errorMsg != "" {
		panic(errorMsg)
	}
}

// RunOne run one
func (c *CaseManager) RunOne(caseName string) {
	if v, ok := c.Cases[caseName]; ok {
		fmt.Println("----------------------------->Start to run case " + caseName + "...")
		rs := v.Call(nil)
		if rs[0].Interface() == nil {
			fmt.Printf("%s SUCCESS\n", caseName)
		} else {
			err := rs[0].Interface().(error)
			if err == nil {
				fmt.Printf("%s SUCCESS\n", caseName)
			} else {
				fmt.Printf("%s FAILED!!!\n", caseName)
			}
		}
	} else {
		fmt.Printf("%s doesn't exist !!! \n", caseName)
	}
	fmt.Println("Pelease check log in ./log")
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
