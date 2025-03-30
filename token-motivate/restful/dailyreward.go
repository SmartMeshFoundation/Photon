package restful

import (
	"math/big"
	"time"

	"fmt"

	"encoding/base64"
	"strings"

	"go.cryptoscope.co/ssb/restful/params"
)

var rewardPeriod = time.Second * 90

/*
RewardDailyTask
巡查数据库usertaskcollect
rewardPeriod处理一次，发送成功则记录，发送不成功则一直重试

	处理：
		1、登录	即时处理，发送激励，每日不超过MaxDailyRewarding
		2、发帖	即时处理，发送激励，每日不超过MaxDailyRewarding
		3、评论	即时处理，发送激励，每日不超过MaxDailyRewarding
		4、NFT	即时处理，发送激励，每日不超过MaxDailyRewarding
	另外：
		5、注册，即时回复，延后激励，条件：在线（重试40分钟）发送激励MLT，接着链上发送SMT
		6、我点赞别人，即时处理，发送激励，每日不超过MaxDailyRewarding
		7、我的受赞，暂时由supernode处理
		4、举报，pub提供接口，发送激励
*/
func RewardProcess() {
	//1、处理usertaskcollect,1-登录 2-发帖 3-评论 4-铸造NFT
	//2、处理未发成功激励的事件
	/*var starttime = req.StartTime
	var endtime = req.EndTime
	taskcollctions, err := likeDB.GetUserTaskCollect(author, msgtype, starttime, endtime)*/
}

//func RecordRewarding2Db()

// ExceedRewardLimit
func ExceedRewardLimit(clientID, rewardType string, msgTime int64, thisAmount int64) bool {
	var starttime = msgTime - time.Hour.Milliseconds()*24
	var endtime = msgTime

	if HealthCheckClientID(clientID) == false {
		fmt.Println(fmt.Errorf("base64 decode of remoteKey failed: %s", clientID))
		return true
	}
	//如果存在未发送成功的记录,也记为本次比较的数量，因为延后会继续处理未成功的事件
	if rewardType != SignUp {
		rewardType = ""
	}
	num, err := likeDB.SelectHistoryReward(clientID, rewardType, starttime, endtime)
	if err != nil {
		fmt.Println(fmt.Errorf("ExceedRewardLimit SelectHistoryReward err =%v", err))
		return true
	}
	historyTokens := new(big.Int).Add(num, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(thisAmount)))
	maxRewardTokes := big.NewInt(0)
	if rewardType == SignUp {
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MaxSignupReward)))
	} else {
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(params.MaxDailyRewarding)))
	}

	if historyTokens.Cmp(maxRewardTokes) == -1 || historyTokens.Cmp(maxRewardTokes) == 0 {
		return false
	} else {
		return true
	}
	/*taskcollctions, err := likeDB.GetUserTaskCollect(clientID, rewardType, starttime, endtime)
	if err != nil {
		fmt.Println(fmt.Sprintf("do ExceedRewardLimit GetUserTaskCollect err=%s", err))
		return false
	}
	var historyNum = 0
	for _, task := range taskcollctions {
		task.ClientEthAddress
	}*/
	return true
}

func HealthCheckClientID(rk string) bool {
	if !strings.Contains(rk, ".ed25519") {
		return false
	}
	if !strings.HasPrefix(rk, "@") {
		return false
	}
	rk = strings.TrimSuffix(rk, ".ed25519")
	rk = strings.TrimPrefix(rk, "@")
	_, err := base64.StdEncoding.DecodeString(rk)
	if err != nil {
		return false
		//return nil, fmt.Errorf("init: base64 decode of --remoteKey failed: %w", err)
	}
	return true
}

func PassRule202302(rewardObject, rewardReason string, msgTime, thisAmount int64) error {
	//登录60分钟内只认为1个有效
	//发贴5分钟内只认为1个有效
	//点赞5分钟内只认为1个有效
	//评论5分钟内只认为1个有效
	//报告垃圾帖子的处理最多10次
	var starttime = msgTime - time.Hour.Milliseconds()*24
	var endtime = msgTime
	num, err := likeDB.SelectHistoryReward(rewardObject, rewardReason, starttime, endtime)
	if err != nil {
		return fmt.Errorf("Rule202302 SelectHistoryReward err =%v", err)
	}
	historyTokens := new(big.Int).Add(num, new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(thisAmount)))
	maxRewardTokes := big.NewInt(0)

	switch rewardReason {
	case PostMessage:
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(10*params.RewardOfPostMessage)))
	case PostComment:
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(10*params.RewardOfPostComment)))
	case DailyLogin:
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(3*params.RewardOfDailyLogin)))
	case LikePost:
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(5*params.RewardOfLikePost)))
	case ReportProblematicPost:
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(1*params.RewardOfReportProblematicPost)))
	default:
		maxRewardTokes = new(big.Int).Mul(big.NewInt(params.Ether), big.NewInt(int64(5*params.RewardOfPostMessage)))
	}

	if historyTokens.Cmp(maxRewardTokes) == -1 || historyTokens.Cmp(maxRewardTokes) == 0 {
		return nil
	} else {
		return fmt.Errorf("Check Rule202302: Not Passed")
	}
}
