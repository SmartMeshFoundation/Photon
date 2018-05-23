package cases

//
//import (
//	"net/http"
//	"encoding/json"
//)
//
//func CSettleChannel(node *Node, no string) {
//	logger.Println("CSettleChannel prepare...")
//	var channels []NodeChannel
//	_, body, _ := DoReq( &ReqParams{
//		No:no,
//		ApiName:"QueryNodeAllChannels",
//		FullUrl:node.Host + "/api/1/channels",
//		Method:http.MethodGet,
//		Payload:"",
//	})
//	json.Unmarshal(body, &channels)
//	logger.Println("CSettleChannel start...")
//	reqParams := ReqParams{
//		No:no,
//		ApiName:"SettleChannel",
//		FullUrl:node.Host + "/api/1/channels/" + channels[0].ChannelAddress,
//		Method:http.MethodPatch,
//		Payload:"{\"state\":\"settled\"}",
//	}
//	status, _, _ := DoReq(&reqParams)
//	DealStatus(status, "CSettleChannel")
//
//}
