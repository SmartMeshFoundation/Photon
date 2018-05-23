package cases

//
//import (
//	"encoding/json"
//	"math/rand"
//	"net/http"
//	"time"
//)
//
////For  InitiatingTransfer API  http body
//type TransferRequest struct {
//	Amount     int32 `json:"amount"`
//	Identifier int64 `json:"identifier"`
//}
//
//func CInitiatingTransfer(node1 *Node, node2 *Node, no string, amount int32) {
//	logger.Println("CInitiatingTransfer prepare...")
//	// build test payload
//	var payload TransferRequest
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	Intsn := r.Int63n(9223372036854775807)
//	payload.Amount = amount
//	payload.Identifier = Intsn
//	p, _ := json.Marshal(payload)
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
//	logger.Println("CInitiatingTransfer start...")
//	reqParams := ReqParams{
//		No:no,
//		ApiName:"InitiatingTransfer",
//		FullUrl:node1.Host + "/api/1/transfers/" + tokens[0] + "/" + targetAddress.Our_address,
//		Method:http.MethodPost,
//		Payload:string(p),
//	}
//	status, _, _ := DoReq(&reqParams)
//	DealStatus(status, "CInitiatingTransfer")
//}
