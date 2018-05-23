package cases

//
//import (
//	"net/http"
//	"encoding/json"
//)
//
//func CCloseChannel(node *Node, no string) {
//	logger.Println("CCloseChannel prepare...")
//	var channels []NodeChannel
//	_, body, _ := DoReq( &ReqParams{
//		No:no,
//		ApiName:"QueryNodeAllChannels",
//		FullUrl:node.Host + "/api/1/channels",
//		Method:http.MethodGet,
//		Payload:"",
//	})
//	json.Unmarshal(body, &channels)
//	logger.Println("CCloseChannel start...")
//	reqParams := ReqParams{
//		No:no,
//		ApiName:"CloseChannel",
//		FullUrl:node.Host + "/api/1/channels/" + channels[0].ChannelAddress,
//		Method:http.MethodPatch,
//		Payload:"{\"state\":\"closed\"}",
//	}
//	status, _, _ := DoReq(&reqParams)
//	DealStatus(status, "CCloseChannel")
//}
