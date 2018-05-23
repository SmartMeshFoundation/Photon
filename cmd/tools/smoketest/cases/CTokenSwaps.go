package cases

//
//import (
//	"encoding/json"
//	"math/rand"
//	"net/http"
//	"strconv"
//	"time"
//)
//
////For TokenSwaps API  http body
//type TokenSwapsPayload struct {
//	Role            string `json:"role"`
//	SendingAmount   int32  `json:"sending_amount"`
//	SendingToken    string `json:"sending_token"`
//	ReceivingAmount int32  `json:"receiving_amount"`
//	ReceivingToken  string `json:"receiving_token"`
//}
//
//func CTokenSwaps(node1 *Node, node2 *Node, no string, sendNum int32, recvNum int32, role string) {
//	logger.Println("CTokenSwaps prepare...")
//	// get id
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	id := r.Int63n(9223372036854775807)
//
//	// get test token
//	_, body, _ := DoReq( &ReqParams{
//		No:no,
//		ApiName:"QueryRegisteredTokens",
//		FullUrl:node1.Host + "/api/1/tokens",
//		Method:http.MethodGet,
//		Payload:"",
//	})
//	var tokens []string
//	json.Unmarshal(body, &tokens)
//
//	// build test payload
//	var payload TokenSwapsPayload
//	payload.Role = role
//	payload.SendingAmount = sendNum
//	if (role == "maker") {
//		payload.SendingToken = tokens[0]
//		payload.ReceivingToken = tokens[1]
//	} else {
//		payload.SendingToken = tokens[1]
//		payload.ReceivingToken = tokens[0]
//	}
//	payload.ReceivingAmount = recvNum
//	p, _ := json.Marshal(payload)
//
//	// get test target address
//	_, body, _ = DoReq(&ReqParams{
//		No:no,
//		ApiName:"QueryNodeAddress",
//		FullUrl:node2.Host + "/api/1/address",
//		Method:http.MethodGet,
//		Payload:"",
//	})
//	var targetAddress struct{
//		Our_address string `json:"our_address"`
//	}
//	json.Unmarshal(body, &targetAddress)
//
//	logger.Println("CTokenSwaps start...")
//	reqParams := ReqParams{
//		No:no,
//		ApiName:"TokenSwaps",
//		FullUrl:node1.Host + "/api/1/token_swaps/" + targetAddress.Our_address + "/" + strconv.FormatInt(id, 10),
//		Method:http.MethodPut,
//		Payload:string(p),
//	}
//	status, _, _ := DoReq(&reqParams)
//	DealStatus(status, "CTokenSwaps")
//}
