package restful

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"strings"

	"net"

	"io"
	"mime/multipart"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"sort"

	"io/ioutil"

	"crypto/aes"

	"encoding/base64"

	"crypto/cipher"

	"encoding/json"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ip2location/ip2location-go/v9"
	"go.cryptoscope.co/ssb/restful/params"
	"go.cryptoscope.co/ssb/restful/wordcount"
)

// clientPublicIP
func clientPublicIP(r *http.Request) string {
	var ip string
	for _, ip = range strings.Split(r.Header.Get("X-Forwarded-For"), ",") {
		ip = strings.TrimSpace(ip)
		if ip != "" && !HasLocalIPddr(ip) {
			return ip
		}
	}
	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" && !HasLocalIPddr(ip) {
		return ip
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		if !HasLocalIPddr(ip) {
			return ip
		}
	}
	return ""
}

// HasLocalIPddr
func HasLocalIPddr(ip string) bool {
	return HasLocalIPAddr(ip)
}

// HasLocalIPAddr
func HasLocalIPAddr(ip string) bool {
	return HasLocalIP(net.ParseIP(ip))
}

// HasLocalIP
func HasLocalIP(ip net.IP) bool {
	if ip.IsLoopback() {
		return true
	}

	ip4 := ip.To4()
	if ip4 == nil {
		return false
	}

	return ip4[0] == 10 || // 10.0.0.0/8
		(ip4[0] == 172 && ip4[1] >= 16 && ip4[1] <= 31) || // 172.16.0.0/12
		(ip4[0] == 169 && ip4[1] == 254) || // 169.254.0.0/16
		(ip4[0] == 192 && ip4[1] == 168) // 192.168.0.0/16
}

// GetPublicIPLocation
func GetPublicIPLocation(w rest.ResponseWriter, r *rest.Request) {
	clientpublicip := clientPublicIP(r.Request)
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetPublicIpLocation ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	/*var req IPLoacation
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var ip = req.PublicIp*/
	if clientpublicip == "" {
		clientpublicip = strings.Split(params.InviteCodeOfPub2, ":")[0]
	}
	var ip = clientpublicip

	db, err := ip2location.OpenDB(params.Ip2LocationLiteDbPath)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusNotImplemented)
		return
	}
	result, err := db.Get_all(ip)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	countryLong := result.Country_long

	pbi := &PubInfoByIP{}
	pbi.ReqPublicIP = clientpublicip
	pbi.ContryShort = result.Country_short
	pbi.ContryLong = countryLong
	pbi.Region = result.Region
	pbi.City = result.City

	if countryLong == "China" {
		pbi.FirstChoicePubHost = fmt.Sprintf("%s:%d", strings.Split(params.InviteCodeOfPub1, ":")[0], params.ServePort)
		pbi.FirstChoicePubInviteCode = params.InviteCodeOfPub1
		pbi.SecondChoicePubHost = fmt.Sprintf("%s:%d", strings.Split(params.InviteCodeOfPub2, ":")[0], params.ServePort)
		pbi.SecondChoicePubInviteCode = params.InviteCodeOfPub2
	} else {
		pbi.FirstChoicePubHost = fmt.Sprintf("%s:%d", strings.Split(params.InviteCodeOfPub2, ":")[0], params.ServePort)
		pbi.FirstChoicePubInviteCode = params.InviteCodeOfPub2
		pbi.SecondChoicePubHost = fmt.Sprintf("%s:%d", strings.Split(params.InviteCodeOfPub1, ":")[0], params.ServePort)
		pbi.SecondChoicePubInviteCode = params.InviteCodeOfPub1
	}
	resp = NewAPIResponse(err, pbi)
}

// TippedOff
func UserFeedBack(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UserFeedBack ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var req UserFeedBackStu
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userssbid = req.UserSsbId
	var submittype = req.SubmitType
	var content = req.Content
	var useremail = req.UserEmail

	var submittime = time.Now().UnixNano() / 1e6
	_, err = likeDB.RecordUserSubmit(userssbid, submittype, content, submittime, useremail, 0)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp = NewAPIResponse(err, "Success")
}

// GetUserFeedBack
func GetUserFeedBack(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetUserFeedBack ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req UserFeedBackStu
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userid = req.UserSsbId
	submits, err := likeDB.SelectUserSubmit(userid)
	resp = NewAPIResponse(err, submits)
}

// GetSomeoneLike
func GetRewardInfo(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetRewardInfo ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req RewardingReq
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clientid = req.ClientID
	//var grandsuccess = req.GrantSuccess
	//var rewardreason = req.RewardReason
	var timefrom = req.TimeFrom
	var timeTo = req.TimeTo

	rresult, err := likeDB.SelectRewardResult(clientid, timefrom, timeTo)
	resp = NewAPIResponse(err, rresult)
}

func GetRewardSubtotals(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetRewardSubtotals ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req RewardingReq
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clientid = req.ClientID
	var grandsuccess = req.GrantSuccess
	//var rewardreason = req.RewardReason
	var timefrom = req.TimeFrom
	var timeTo = req.TimeTo

	rresult, err := likeDB.SelectRewardSum(clientid, grandsuccess, timefrom, timeTo)

	resp = NewAPIResponse(err, rresult)
}

func GetRewardSummary(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetRewardSummary ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req RewardingReq
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clientid = req.ClientID
	var grandsuccess = req.GrantSuccess
	var timefrom = req.TimeFrom
	var timeTo = req.TimeTo

	rresult, err := likeDB.SelectRewardSummary(clientid, grandsuccess, timefrom, timeTo)

	resp = NewAPIResponse(err, rresult)
}

// GetAllSetLikes
func GetAllSetLikes(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetAllSetLikes ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	setlikes, err := likeDB.SelectUserSetLikeInfo("")
	resp = NewAPIResponse(err, setlikes)
}

// GetSomeoneLike
func GetSomeoneSetLikes(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetSomeoneSetLikes ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req Name2ProfileReponse
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var cid = req.ID
	setlikes, err := likeDB.SelectUserSetLikeInfo(cid)
	resp = NewAPIResponse(err, setlikes)
}

// NotifyCreatedNFT
func NotifyCreatedNFT(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> NotifyCreatedNFT ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req ReqCreatedNFT
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cid = req.ClientID
	var ctime = req.NftCreatedTime
	var tx = req.NfttxHash
	var tokenid = req.NftTokenId
	var storeurl = req.NftStoredUrl
	_, err = likeDB.InsertUserTaskCollect(params.PubID, cid, "", "4", "", ctime, tx, tokenid, storeurl)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	{ //发送激励
		name2addr, err := GetNodeProfile(cid)
		if err != nil || len(name2addr) != 1 {
			fmt.Println(fmt.Errorf(MintNft+" Reward %s ethereum address failed, err= not found or %s", cid, err))
		} else {
			ehtAddr := name2addr[0].EthAddress
			go PubRewardToken(ehtAddr, int64(params.RewardOfMintNft), cid, MintNft, "", time.Now().UnixNano()/1e6)
		}
	}

	resp = NewAPIResponse(err, "Success")
}

// NotifyUserLogin
func NotifyUserLogin(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> NotifyUserLogin ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req ReqUserLoginApp
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cid = req.ClientID
	var logintime = req.LoginTime
	_, err = likeDB.InsertUserTaskCollect(params.PubID, cid, "", "1", "", logintime, "", "", "")
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	{ //发送激励
		name2addr, err := GetNodeProfile(cid)
		if err != nil || len(name2addr) != 1 {
			fmt.Println(fmt.Errorf(DailyLogin+" Reward %s ethereum address failed, err= not found or %s", cid, err))
		} else {
			ehtAddr := name2addr[0].EthAddress
			go PubRewardToken(ehtAddr, int64(params.RewardOfDailyLogin), cid, DailyLogin, "", time.Now().UnixNano()/1e6)
		}
	}
	resp = NewAPIResponse(err, "Success")
}

// GetUserDailyTasks
func GetUserDailyTasks(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetUserDailyTasks ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req ReqUserTask
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var author = req.Author
	var msgtype = req.MessageType
	var starttime = req.StartTime
	var endtime = req.EndTime

	taskcollctions, err := likeDB.GetUserTaskCollect(author, msgtype, starttime, endtime)
	resp = NewAPIResponse(err, taskcollctions)
}

// GetEventSensitiveWord
func GetEventSensitiveWord(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetEventSensitiveWord ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req EventSensitive
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var tag = req.DealTag
	senvents, err := likeDB.SelectSensitiveWordRecord(tag)
	resp = NewAPIResponse(err, senvents)
}

// DealSensitiveWord
func DealSensitiveWord(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> DealSensitiveWord ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()

	var req EventSensitive
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var msgkey = req.MessageKey
	var dealtag = req.DealTag
	var dealtime = time.Now().UnixNano() / 1e6
	var author = req.MessageAuthor
	_, err = likeDB.UpdateSensitiveWordRecord(dealtag, dealtime, msgkey)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.DealTag == "1" { ////for table sensitivewordrecord, dealtag=0初始化  =1属实 =2否定
		// block 'the author who publish sensitive word' ONCE
		err = contactSomeone(nil, author, false, true)
		if err != nil {
			resp = NewAPIResponse(err, fmt.Sprintf("block %s failed", author))
			return
		}
		fmt.Println(fmt.Sprintf(PrintTime()+"Success to block %s", author))
	}
	resp = NewAPIResponse(err, "success")
}

// TippedOff
func TippedOff(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> TippedWhoOff ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var req TippedOffStu
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var plaintiff = req.Plaintiff
	var defendant = req.Defendant
	var mkey = req.MessageKey
	var reasons = req.Reasons

	if defendant == params.PubID {
		resp = NewAPIResponse(err, fmt.Sprintf("Permission denied, from pub : %s", params.PubID))
		return
	}
	var recordtime = time.Now().UnixNano() / 1e6
	lstid, err := likeDB.InsertViolation(recordtime, plaintiff, defendant, mkey, reasons)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if lstid == -1 {
		resp = NewAPIResponse(err, "You've already reported it, thank your again👍")
		return
	}

	resp = NewAPIResponse(err, "Success, the pub administrator will verify as soon as possible, thank you for your report👍")
}

// TippedOffInfo get infos
func GetTippedOffInfo(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetTippedOffInfo ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req TippedOffStu
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	datas, err := likeDB.SelectViolationByWhere(req.Plaintiff, req.Defendant, req.MessageKey, req.Reasons, req.DealTag)

	resp = NewAPIResponse(err, datas)
}

// DealTippedOff
func DealTippedOff(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> DealTippedOff ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var req TippedOffStu
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dtime = time.Now().UnixNano() / 1e6
	_, err = likeDB.UpdateViolation(req.DealTag, dtime, req.Dealreward, req.Plaintiff, req.Defendant, req.MessageKey)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.DealTag == "1" { ////for table violationrecord, dealtag=0举报 =1属实 =2事实不清,不予处理
		//1 unfollow and block 'the defendant' and sign him to blacklist
		err = contactSomeone(nil, req.Defendant, false, true)
		if err != nil {
			resp = NewAPIResponse(err, fmt.Sprintf("Unfollow and block %s failed, err=%s", req.Defendant, err))
			return
		}
		fmt.Println(fmt.Sprintf(PrintTime()+"Success to Unfollow and block %s", req.Defendant))

		{ //发送激励
			name2addr, err := GetNodeProfile(req.Plaintiff)
			if err != nil || len(name2addr) != 1 {
				fmt.Println(fmt.Errorf(ReportProblematicPost+" Reward %s ethereum address failed, err= not found or %s", req.Plaintiff, err))
			} else {
				ehtAddr := name2addr[0].EthAddress
				go PubRewardToken(ehtAddr, int64(params.RewardOfReportProblematicPost), req.Plaintiff, ReportProblematicPost, req.MessageKey, dtime)
			}
		}

		_, err = likeDB.UpdateViolation(req.DealTag, dtime, fmt.Sprintf("%d%s", params.RewardOfReportProblematicPost, "e18-"), req.Plaintiff, req.Defendant, req.MessageKey)
		if err != nil {
			rest.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		resp = NewAPIResponse(err, fmt.Sprintf("success, [%s] has been block by [pub administrator], and pub send award token to [%s]", req.Defendant, req.Plaintiff))
		return
	}
	resp = NewAPIResponse(err, "success")
}

// GetPubWhoami
func GetPubWhoami(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetPubWhoami ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()

	pinfo := &Whoami{}
	pinfo.Pub_Id = params.PubID
	pinfo.Pub_Eth_Address = params.PhotonAddress
	resp = NewAPIResponse(nil, pinfo)
	return
}

// clientid2Profile
func clientid2Profiles(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> node-infos ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	name2addr, err := GetAllNodesProfile()
	resp = NewAPIResponse(err, name2addr)
	return
}

// clientid2Profile
func clientid2Profile(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> node-infos ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req Name2ProfileReponse
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cid = req.ID
	name2addr, err := GetNodeProfile(cid)
	resp = NewAPIResponse(err, name2addr)
}

// rewardForSomeReason
func rewardForSomeReason(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> rewardForSomeReason ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req Name2ProfileReponse
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cid = req.ID
	name2addr, err := GetNodeProfile(cid)
	if err != nil || len(name2addr) != 1 {
		fmt.Println(fmt.Errorf(InviteEarn+" Reward(DealInviteEarn) %s ethereum address failed, err= not found or %s", cid, err))
	} else {
		ehtAddr := name2addr[0].EthAddress
		go PubRewardToken(ehtAddr, int64(params.RewardOfInvite), cid, InviteEarn, "", time.Now().UnixNano()/1e6)
	}
	resp = NewAPIResponse(err, true)
}

// UpdateEthAddr
func UpdateEthAddr(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UpdateEthAddr ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var req = &Name2ProfileReponse{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	/*//此处跳过校验，前端不好处理
	_, err = HexToAddress(req.EthAddress)
	if err != nil {
		resp = NewAPIResponse(err, nil)
		return
	}*/

	//在一个小时内不允许同一个IP地址注册
	clientpublicip := clientPublicIP(r.Request)
	nowT := time.Now().UnixNano() / 1e6
	fmt.Printf("now check-ip=%v", params.CheckIP)
	if params.CheckIP {

		if clientpublicip == "" {
			rest.Error(w, fmt.Errorf("Unknown network IP address").Error(), http.StatusBadRequest)
			return
		}
		if t, ok := RegisterSourceMap[clientpublicip]; ok {
			fmt.Printf("now ip=%v,regtime=%v", clientpublicip, t)
			ca := nowT - t
			if ca <= 1*time.Hour.Milliseconds() {
				rest.Error(w, fmt.Errorf("It is illegal to register with the same IP(%s) address within 1 Hour", clientpublicip).Error(), http.StatusBadRequest)
				return
			}
		}
	}

	ethAddress := common.HexToAddress(req.EthAddress)
	appVersion := req.AppVersion
	personalInviteEncCode := req.PersonalInviteCode
	if HealthCheckClientID(req.ID) == false {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var personalInviteDecCodeStr = ""
	var personalInviteDecCodeBool = false
	if personalInviteEncCode != "" {
		//var personalInviteEncCodeBytes = []byte(personalInviteEncCode)
		var personalInviteEncCodeBytes, err = base64.StdEncoding.DecodeString(personalInviteEncCode)
		if err != nil {
			personalInviteDecCodeStr = err.Error()
		}
		personalInviteDecCodeStr, err = AesDecrypt(personalInviteEncCodeBytes, req.ID)
		if err != nil {
			//如果邀请码有问题，允许继续注册的
			personalInviteDecCodeStr = err.Error()
		} else {
			personalInviteDecCodeBool = true
		}
	}

	_, err = likeDB.UpdateUserProfile(0, req.ID, req.Name, ethAddress.String(), appVersion, personalInviteDecCodeStr)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Println("personalInviteDecCodeBool")
		fmt.Println(personalInviteDecCodeBool)
		go NewChannelDeal(ethAddress.String(), req.ID, time.Now().UnixNano()/1e6, personalInviteDecCodeBool)
		resp = NewAPIResponse(err, "success")
		RegisterSourceMap[clientpublicip] = nowT

		//对邀请人发放
		if personalInviteDecCodeBool {
			inviter := strings.Split(personalInviteDecCodeStr, "|")[0]
			name2addr, err := GetNodeProfile(inviter)
			if err != nil || len(name2addr) != 1 {
				fmt.Println(fmt.Errorf(InviteEarn+" Reward(DealInviteEarn) %s ethereum address failed, err= not found or %s", inviter, err))
			} else {
				ehtAddr := name2addr[0].EthAddress
				go PubRewardToken(ehtAddr, int64(params.RewardOfInvite), inviter, InviteEarn, "", nowT)
			}

		}
	}

}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData) // 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func AesDecrypt(cipherText []byte, inviteWho string) (string, error) {
	key := []byte("A27F8f580C01Db06")
	iv := []byte("A27F8f580C01Db06")

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainText := make([]byte, len(cipherText))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plainText, cipherText)
	plainText = PKCS5UnPadding(plainText)

	//check personal invite code: (ssbId|platform|version|deviceId|platform|version)
	//出示邀请码的人的ssbid|他的手机型号|他的APP版本号|注册人的手机设备ID|注册人的手机型号|注册人的APP版本号
	//ret := fmt.Sprintf("%s", plainText)
	ret := string(plainText[:])
	fmt.Println("ret:" + ret)

	//检查邀请码格式合法性
	if strings.Count(ret, "|") != 5 {
		return "", errors.New("formart error when decrypt the personal invite code")
	}

	cssbid := strings.Split(ret, "|")[0]

	//检查邀请人的ssbid合法性
	if HealthCheckClientID(cssbid) == false {
		return "", errors.New("check inviter's ssb-id error when decrypt the personal invite code")
	}
	//自己不能邀请自己
	if cssbid == inviteWho {
		return "", errors.New("you cannot invite yourself")
	}

	//检查受邀人的deviceID是否存在过
	userInfos, err := GetNodeProfile("")
	if err != nil {
		return "", errors.New("check if the device-id exists error(internal)")
	}
	deviceIDExist := false
	registerDeviceID := strings.Split(ret, "|")[3]
	for _, u := range userInfos {
		var inviteCode = u.PersonalInviteCode
		if strings.Count(inviteCode, "|") != 5 {
			continue
		}
		tmpdevid := strings.Split(inviteCode, "|")[3]
		if tmpdevid == registerDeviceID {
			deviceIDExist = true
			break
		}

	}
	if deviceIDExist {
		return "", errors.New("check if the device-id exists error")
	}

	//检查邀请人是否在系统中注册过
	// todo(2023-09-24)Tag:增加邀请人是否在其他所有pub上注册过，通过grpc增加一条：pubX的注册触发pubY的邀请人激励
	chkOnAllPub := CheckNodeProfileOnAllPubs(cssbid)
	if chkOnAllPub == false { //CheckNodeProfileOnAllPubs如果没有error已经触发被邀请人所在pub的激励动作
		return "", errors.New("check if the inviter exists error(internal)")
	}
	/*if len(userInfo) != 1 || userInfo[0].ID != cssbid {
		return "", errors.New("check if the inviter exists error")
	}*/

	return ret, nil
}

// NotifyForceCloseChannel
func NotifyForceCloseChannel(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> NotifyForceCloseChannel ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()

	var req = &CloseChannelReq{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Signature != "/ewqDIECew3434q5dsEuhjyhgyuut" {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ethAddress := common.HexToAddress(req.PartnerAddress)
	err = pubNode.Close(params.TokenAddress, ethAddress.String())
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		resp = NewAPIResponse(err, "close channel success")
	}

}

// GetAllNodesProfile
func GetAllNodesProfile() (datas []*Name2ProfileReponse, err error) {
	profiles, err := likeDB.SelectUserProfile("")
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Failed to db-SelectUserProfileAll", err))
		return
	}
	datas = profiles
	return
}

// GetNodeProfile
func GetNodeProfile(cid string) (datas []*Name2ProfileReponse, err error) {
	profile, err := likeDB.SelectUserProfile(cid)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Failed to db-SelectUserEthAddrAll", err))
		return
	}
	datas = profile
	return
}

// CheckNodeProfileOnAllPubs
func CheckNodeProfileOnAllPubs(cid string) bool {
	ret := false
	profile, err := likeDB.SelectUserProfile(cid)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Failed to db-SelectUserEthAddrAll-local-pub", err))
		//return ret, err
	}
	if len(profile) == 1 {
		ret = true
		return ret
	}

	//继续检查其他的pub数据库
	if checkUserinfoOnOtherPub(cid) == true {
		ret = true
		return ret
	}
	return ret
}

// checkUserinfoOnOtherPub :
func checkUserinfoOnOtherPub(cssid string) bool {
	p, err := json.Marshal(Name2ProfileReponse{
		ID: cssid,
	})
	req := &Req{
		FullURL: "http://" + params.AnotherServe + "/ssb/api/node-info",
		Method:  http.MethodPost,
		Payload: string(p),
		Timeout: time.Second * 20,
	}
	body, err := req.Invoke()
	if err != nil {
		fmt.Println("1111111" + req.FullURL)
		return false
	}

	var userInfos []*Name2ProfileReponse
	err = json.Unmarshal(body, &userInfos)
	if err != nil {
		fmt.Println(fmt.Errorf("bodylen=%d,body=%s", len(body), string(body)))
		fmt.Println("2222222")
		return false
	}
	if len(userInfos) == 1 {
		//通知对方pub发送
		pl, _ := json.Marshal(Name2ProfileReponse{
			ID: cssid,
		})
		reqc := &Req{
			FullURL: "http://" + params.AnotherServe + "/ssb/api/internal/reward-for-some-reason",
			Method:  http.MethodPost,
			Payload: string(pl),
			Timeout: time.Second * 20,
		}
		_, err = reqc.Invoke()

		return true
	}
	fmt.Println("3333333:" + req.FullURL)
	fmt.Println(len(userInfos))
	fmt.Println("3333333:" + cssid)
	return false
}

// GetAllLikes
func GetAllLikes(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetAllLikes ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	likes, err := CalcGetLikeSum("")

	resp = NewAPIResponse(err, likes)
}

// GetSomeoneLike
func GetSomeoneLike(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetSomeoneLike ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req Name2ProfileReponse
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var cid = req.ID
	like, err := CalcGetLikeSum(cid)
	resp = NewAPIResponse(err, like)
}

// GetAllNodesProfile
func CalcGetLikeSum(someoneOrAll string) (datas map[string]*LasterNumLikes, err error) {
	likes, err := likeDB.SelectLikeSum(someoneOrAll)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Failed to db-SelectLikeSum", err))
		return
	}
	datas = likes
	return
}

func PostWordCountBigThan10(words string) bool {
	/*wlen := 0
	wordsSlice := strings.Split(words, " ")
	for _, word := range wordsSlice {
		if word != "" {
			wlen++
		}
	}
	return wlen > 10*/
	counter := &wordcount.WordCounter{}
	counter.Stat(words)
	return counter.Words >= 10
}

// --------------------------------------------------------
// --------------------------------------------------------
// LoadGameInfo
func LoadGameInfo(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> LoadGameInfo ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	ginfos, err := likeDB.SelectGameInfo()
	resp = NewAPIResponse(err, ginfos)
	return
}

// UploadGamePlay
func UploadGamePlay(w rest.ResponseWriter, r *rest.Request) {
	r.ParseMultipartForm(32 << 20)
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UploadGamePlay ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	file, handler, err := r.FormFile("player_play_photos")
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"UploadGamePlay err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	defer file.Close()

	ethAddr := r.PostFormValue("wallet_address")
	ssbID := r.FormValue("ssb_id")
	submitType := r.FormValue("submit_type")
	gameId := r.FormValue("game_id")
	playerId := r.FormValue("player_id")
	playerName := r.FormValue("player_name")
	playerMark := r.FormValue("player_mark")
	if ethAddr == "" || submitType == "" || gameId == "" || playerId == "" || playerMark == "" {
		w.WriteJson(errors.New("incomplete game voucher data"))
		return
	}

	/*photos, err := NewZipHandle(file, handler.Size)
	photos.OnHandle(CreateFile)
	err = photos.Handle()
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"dbUpLoad Handle err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}*/
	var dtime = time.Now().UnixNano() / 1e6
	playResPath, err := CreatePicFile(params.GameUserFilePath, ethAddr, strconv.FormatInt(dtime, 10)+strings.ToLower(path.Ext(handler.Filename)), file, handler)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"CreatePicFile thumb err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	fmt.Println("UploadGamePlay:" + playResPath)
	playResPath1 := "/" + ethAddr + "/" + strconv.FormatInt(dtime, 10)
	ginfos, err := likeDB.InsertPlayHistory(ethAddr, submitType, gameId, playerId, playerName, playerMark, playResPath1, dtime, ssbID)
	resp = NewAPIResponse(err, ginfos)
	fmt.Fprintln(w.(http.ResponseWriter), "success")

	return
}

// GetGamePlay
func GetGamePlay(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> LoadGameInfo ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req = &PlayHistory{}
	err := r.DecodeJsonPayload(req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	ssbid := req.Ssbid
	gameid := req.GameId
	ginfos, err := likeDB.SelectPlayHistory(ssbid, gameid)
	resp = NewAPIResponse(err, ginfos)
	return
}

// UploadGameInfo
func UploadGameInfo(w rest.ResponseWriter, r *rest.Request) {
	r.ParseMultipartForm(32 << 20)
	//w.Header().Set("Content-Type", "multipart/form-data")
	//w.Header().Set("Content-Type", "application/json")
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UploadGameInfo ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	//r.Header.Set("Content-Type", "application/x-www-form-urlencoded ; charset=UTF-8")

	coverPhoto, handler1, err := r.FormFile("game_cover_photo")
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Upload coverPhoto err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	defer coverPhoto.Close()

	thumbnail, handler2, err := r.FormFile("game_thumbnail")
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Upload thumbnail err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	defer thumbnail.Close()

	baners, handler3, err := r.FormFile("game_banner_photos")
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"Upload game_banner_photos err= [%s]", err))
		w.WriteJson(err.Error())
		return
	}
	defer baners.Close()

	gameName := r.PostFormValue("game_name")
	gameVersion := r.FormValue("game_version")
	gameType := r.FormValue("game_type")
	gameIntro := r.FormValue("game_introduction")
	gamePlay := r.FormValue("game_play")
	downloadlink := r.FormValue("resource_download")
	if gameName == "" || gameVersion == "" || gameType == "" || gameIntro == "" || gamePlay == "" || downloadlink == "" {
		//w.WriteJson(errors.New("incomplete game info data"))
		resp = NewAPIResponse(errors.New("incomplete game info data"), "failed")
		return
	}

	coverpath, err := CreatePicFile(params.GameResourcePath, gameName, "cover"+strings.ToLower(path.Ext(handler1.Filename)), coverPhoto, handler1)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"CreatePicFile cover err= [%s]", err))
		resp = NewAPIResponse(err, "CreatePicFile cover err")
		return
	}
	thumbpath, err := CreatePicFile(params.GameResourcePath, gameName, "thumb"+strings.ToLower(path.Ext(handler2.Filename)), thumbnail, handler2)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"CreatePicFile thumb err= [%s]", err))
		resp = NewAPIResponse(err, "CreatePicFile thumb err")
		return
	}
	bannerspath, err := CreatePicFile(params.GameResourcePath, gameName, "banner"+strings.ToLower(path.Ext(handler3.Filename)), baners, handler3)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"CreatePicFile banner err= [%s]", err))
		resp = NewAPIResponse(err, "CreatePicFile banner err")
		return
	}

	/*coverpath1 := strings.Replace(coverpath, params.GameResourcePath, "", 1)
	thumbpath1 := strings.Replace(thumbpath, params.GameResourcePath, "", 1)
	bannerspath1 := strings.Replace(bannerspath, params.GameResourcePath, "", 1)*/
	coverpath1 := "/" + gameName + "/cover"
	thumbpath1 := "/" + gameName + "/thumb"
	bannerspath1 := "/" + gameName + "/banner"

	var dtime = time.Now().UnixNano() / 1e6 //strconv.FormatInt(dtime, 10)
	_, err = likeDB.InsertGameInfo(gameName, gameName, gameVersion, coverpath1, thumbpath1, gameType,
		gameIntro, gamePlay, bannerspath1, downloadlink, dtime)
	if err != nil {
		fmt.Println(fmt.Sprintf(PrintTime()+"UploadGameInfo Handle err= [%s]", err))
		w.WriteJson(err.Error())
		return

	}
	resp = NewAPIResponse(nil, "success")
	//fmt.Fprintln(w.(http.ResponseWriter), "success")
	fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UploadGameInfo\t,file = [%s]", coverpath))
	fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UploadGameInfo\t,file = [%s]", thumbpath))
	fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> UploadGameInfo\t,file = [%s]", bannerspath))

	//return
}

func GetResource(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetResource ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	filename := fmt.Sprintf("%s.jpg", r.PathParam("resourcename"))

	file, err := os.Open(filepath.Join(params.GameResourcePath, r.PathParam("gamename"), filename))
	//fmt.Println(file.Name())
	if err != nil {
		resp = NewAPIResponse(err, "get resource failed")
		return
	}
	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		resp = NewAPIResponse(err, "get resource failed")
		return
	}
	h := w.Header()
	h.Set("Content-type", "application/octet-stream")
	h.Set("Content-Disposition", "attachment;filename="+filename)
	resp = NewAPIResponse(err, buff)
	return
}

/*
	func GetUserPhoto(w rest.ResponseWriter, r *rest.Request) {
		var resp *APIResponse
		defer func() {
			fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetUserPhoto ,err=%s", resp.ErrorMsg))
		}()
		ethAddr := r.PathParam("eth-address")
		picName := r.PathParam("pic-name")

		userPhotoDir := filepath.Join(params.GameUserFilePath, ethAddr)
		files, err := ioutil.ReadDir(userPhotoDir)
		userPhotoName := ""
		for _, file := range files {
			filenameWithSuffix := file.Name()
			fileSuffix := path.Ext(filenameWithSuffix)
			filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
			if filenameOnly == picName {
				userPhotoName = filenameWithSuffix
				break
			}
		}

		file, err := os.Open(filepath.Join(params.GameUserFilePath, ethAddr, userPhotoName))
		fmt.Println("download UserPhoto:" + file.Name())
		if err != nil {
			resp = NewAPIResponse(err, "get resource failed")
			return
		}
		defer file.Close()
		buff, err := ioutil.ReadAll(file)
		if err != nil {
			resp = NewAPIResponse(err, "get user photo failed")
			return
		}
		h := w.Header()
		h.Set("Content-type", "application/octet-stream")
		h.Set("Content-Disposition", "attachment;filename="+userPhotoName)
		w.WriteJson(buff)
		return
	}
*/
func GetUserPhoto(w http.ResponseWriter, r *http.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetUserPhoto ,err=%s", resp.ErrorMsg))
	}()
	ethAddr := strings.Split(r.URL.Path, "/")[7]
	picName := strings.Split(r.URL.Path, "/")[8]
	fmt.Println("ethAddr:" + ethAddr)
	fmt.Println("picName:" + picName)

	userPhotoDir := filepath.Join(params.GameUserFilePath, ethAddr)
	files, err := ioutil.ReadDir(userPhotoDir)
	userPhotoName := ""
	for _, file := range files {
		filenameWithSuffix := file.Name()
		fileSuffix := path.Ext(filenameWithSuffix)
		filenameOnly := strings.TrimSuffix(filenameWithSuffix, fileSuffix)
		if filenameOnly == picName {
			userPhotoName = filenameWithSuffix
			break
		}
	}

	file, err := os.Open(filepath.Join(params.GameUserFilePath, ethAddr, userPhotoName))
	fmt.Println("download UserPhoto:" + file.Name())
	if err != nil {
		resp = NewAPIResponse(err, "get resource failed")
		return
	}
	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		resp = NewAPIResponse(err, "get user photo failed")
		return
	}
	h := w.Header()
	h.Set("Content-type", "image/png")
	h.Set("Content-Disposition", "attachment;filename="+userPhotoName)
	w.Write(buff)
	return
}

func IsPic(picFile *multipart.FileHeader) (ok bool) {
	extName := strings.ToLower(path.Ext(picFile.Filename))
	allowExtMap := map[string]bool{
		".jpg":  true,
		".png":  true,
		".jpeg": true,
		".zip":  true,
	}
	_, ok = allowExtMap[extName]
	return
}

func CreatePicName(prefix string) string {
	return prefix + "-" + strconv.FormatInt(time.Now().UnixNano(), 10)

}

func CreatePicFile(path, subpath, filename string, file2 multipart.File, head *multipart.FileHeader) (filePath string, err error) {

	if !IsPic(head) {
		return "", errors.New("the format of image resources only supports jpg")
	}
	ppath := filepath.Join(path, subpath)
	if !IsFileExist(ppath) {
		os.MkdirAll(ppath, os.ModePerm)
	}
	f, err := os.Create(filepath.Join(ppath, filename))
	defer f.Close()
	if err != nil {
		return "", err
	}
	_, err = io.Copy(f, file2)
	if err != nil {
		return "", err
	}
	return filepath.Join(ppath, filename), nil
}

type NodeProfileWrapper struct {
	nodeinfo []*Name2ProfileReponse
	by       func(p, q *Name2ProfileReponse) bool
}

type SortBy func(p, q *Name2ProfileReponse) bool

func (np NodeProfileWrapper) Len() int {
	return len(np.nodeinfo)
}
func (np NodeProfileWrapper) Swap(i, j int) {
	np.nodeinfo[i], np.nodeinfo[j] = np.nodeinfo[j], np.nodeinfo[i]
}
func (np NodeProfileWrapper) Less(i, j int) bool {
	return np.by(np.nodeinfo[i], np.nodeinfo[j])
}
func sortNodeProfiles(nodes []*Name2ProfileReponse, by SortBy) {
	sort.Sort(NodeProfileWrapper{nodes, by})
}

type FriendSort struct {
	SortType   string `json:"sortby"`
	SortNumber int    `json:"number"`
}

// GetFriendMaybe
func GetFriendMaybe(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetFriendMaybe ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()

	var req FriendSort
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var stype = req.SortType
	var snum = req.SortNumber
	if stype != "reg-time" && stype != "like-time" {
		resp = NewAPIResponse(err, "field 'sortby' as 'reg-time' or 'like-time'")
		return
	}

	name2addr, err := GetAllNodesProfile()
	if stype == "reg-time" {
		sortNodeProfiles(name2addr, func(p, q *Name2ProfileReponse) bool {
			return q.RegisteTime < p.RegisteTime
		})
	} else {
		sortNodeProfiles(name2addr, func(p, q *Name2ProfileReponse) bool {
			return q.LastactiveTime < p.LastactiveTime
		})
	}

	if snum > len(name2addr) {
		snum = len(name2addr)
	}

	resp = NewAPIResponse(err, name2addr[0:snum])

	return
}

// DealGameEarn
func DealGameEarn(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> DealGameEarn ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	var req PlayHistory
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var dtime = time.Now().UnixNano() / 1e6
	var grantToken int64
	var submitTime = req.SubmitTime
	if req.DealTag == "1" {
		grantToken = int64(params.RewardOfGameEarn)
	}
	_, err = likeDB.UpdateGameEarn(req.DealTag, req.Ssbid, submitTime, dtime, grantToken)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if req.DealTag == "1" { ////for table violationrecord, dealtag=0提交 =1 通过 =2不予通过

		{ //发送激励
			name2addr, err := GetNodeProfile(req.Ssbid)
			if err != nil || len(name2addr) != 1 {
				fmt.Println(fmt.Errorf(ReportProblematicPost+" Reward(DealGameEarn) %s ethereum address failed, err= not found or %s", req.Ssbid, err))
			} else {
				ehtAddr := name2addr[0].EthAddress
				go PubRewardToken(ehtAddr, int64(params.RewardOfGameEarn), req.Ssbid, GameEarn, "", submitTime)
			}
		}

		resp = NewAPIResponse(err, fmt.Sprintf("success, pub send award token to [%s]", req.Ssbid))
		return
	}
	resp = NewAPIResponse(err, "success")
}

// GetPlayEarn
func GetPlayEarn(w rest.ResponseWriter, r *rest.Request) {
	var resp *APIResponse
	defer func() {
		fmt.Println(fmt.Sprintf(PrintTime()+"Restful Api Call ----> GetPlayEarn ,err=%s", resp.ErrorMsg))
		writejson(w, resp)
	}()
	var req RewardingReq
	err := r.DecodeJsonPayload(&req)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var clientid = req.ClientID
	var timefrom = req.TimeFrom
	var timeTo = req.TimeTo

	rresult, err := likeDB.SelectPlayEarn(clientid, timefrom, timeTo)
	resp = NewAPIResponse(err, rresult)
}
