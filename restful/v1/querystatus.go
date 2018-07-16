package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/log"
	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
)

type dataMap map[string]interface{}

/*
Address is api of /api/1/address
*/
func Address(w rest.ResponseWriter, r *rest.Request) {
	data := make(dataMap)
	data["our_address"] = RaidenAPI.Raiden.NodeAddress.String()
	w.WriteJson(data)
}

/*
Tokens is api of /api/1/tokens
*/
func Tokens(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(RaidenAPI.GetTokenList())
}

type partnersData struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

/*
TokenPartners is api of /api/1/:token/:partner
*/
func TokenPartners(w rest.ResponseWriter, r *rest.Request) {
	tokenAddr := common.HexToAddress(r.PathParam("token"))
	log.Trace(fmt.Sprintf("TokenPartners tokenAddr=%s", utils.APex(tokenAddr)))
	tokens := RaidenAPI.GetTokenList()
	found := false
	for _, t := range tokens {
		if t == tokenAddr {
			found = true
			break
		}
	}
	if !found {
		rest.Error(w, "token doesn't exist", http.StatusNotFound)
		return
	}
	chs, err := RaidenAPI.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var datas []*partnersData
	for _, c := range chs {
		d := &partnersData{
			PartnerAddress: c.PartnerAddress.String(),
			Channel:        "api/1/channles/" + c.ChannelIdentifier.ChannelIdentifier.String(),
		}
		datas = append(datas, d)
	}
	w.WriteJson(datas)
}
