package cases

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/SmartMeshFoundation/SmartRaiden/cmd/tools/smoketest/models"
)

// LogPath : api test case log path
var LogPath string

// Logger : api test case logger
var Logger *log.Logger

// FailCases : failed case names
var FailCases []string

// GlobalPassword : default password
var GlobalPassword = "123"

// InitCaseLogger : init logger
func InitCaseLogger(logPath string) {
	LogPath = logPath
	caseLogFile, err := os.Create(LogPath)
	//defer caseLogFile.Close()
	if err != nil {
		log.Fatalln("Create smoketest.log file error !")
	}
	Logger = log.New(caseLogFile, "", log.LstdFlags|log.Lshortfile)
	log.Println("Case log file : " + LogPath)
}

// APITestCase : api test case
type APITestCase struct {
	Req              *models.Req `json:"req"`
	CaseName         string      `json:"case_name"`
	TargetStatusCode int         `json:"target_status_code"`
	TargetBody       interface{} `json:"target_body,omitempty"`
	AllowFail        bool        `json:"allow_fail"`
}

// Run : run api test case
func (c *APITestCase) Run() {
	Logger.Printf("Test case [%s] START...", c.CaseName)
	if c.TargetStatusCode != 0 {
		Logger.Printf("Expect response http code : [%d]", c.TargetStatusCode)
	}
	if c.TargetBody != nil && c.TargetBody != "" {
		bodyStr, _ := json.MarshalIndent(c.TargetBody, "", "\t")
		Logger.Printf("Expect response http body : \n%s\n", bodyStr)
	}
	Logger.SetFlags(0)
	Logger.Printf("----->SEND : %s", c.Req.ToString())
	startTime := time.Now()
	statusCode, body, err := c.Req.Invoke()
	duration := time.Since(startTime)
	Logger.Printf("----->RECEIVE %d in %d ms : ", statusCode, duration.Nanoseconds()/1e6)
	if body != nil && len(body) > 0 {
		Logger.Printf("%s\n", string(body))
	}
	Logger.SetFlags(log.LstdFlags | log.Lshortfile)
	if err != nil {
		if !c.AllowFail {
			Logger.Printf("Test case [%s] FAILED !!!", c.CaseName)
			Logger.Println("allowFail = false,exit")
			panic(err)
		}
	}
	Logger.Printf("Expect [%d] and Get [%d]", c.TargetStatusCode, statusCode)
	if statusCode != c.TargetStatusCode {
		FailCases = append(FailCases, c.CaseName)
		if statusCode == 0 {
			log.Printf("Case [%-40s] TIMEOUT !!!", c.CaseName)
			Logger.Printf("Test case [%s] TIMEOUT !!!", c.CaseName)
		} else {
			log.Printf("Case [%-40s] FAILED !!!", c.CaseName)
			Logger.Printf("Test case [%s] FAILED !!!", c.CaseName)
		}
		if !c.AllowFail {
			log.Println("AllowFail = false,exit")
			Logger.Println("allowFail = false,exit")
			panic(string(body))
		}
	} else {
		log.Printf("Case [%-40s] SUCCESS", c.CaseName)
		Logger.Printf("Test case [%s] SUCCESS", c.CaseName)
	}
	Logger.Println("==================================================")
}
