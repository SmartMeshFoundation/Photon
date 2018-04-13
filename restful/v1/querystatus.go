package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/SmartRaiden/utils"
	"github.com/ant0ine/go-json-rest/rest"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
)

type DataMap map[string]interface{}

func Address(w rest.ResponseWriter, r *rest.Request) {
	data := make(DataMap)
	data["our_address"] = RaidenApi.Raiden.NodeAddress.String()
	w.WriteJson(data)
}

func Tokens(w rest.ResponseWriter, r *rest.Request) {
	w.WriteJson(RaidenApi.GetTokenList())
}

type partnersData struct {
	PartnerAddress string `json:"partner_address"`
	Channel        string `json:"channel"`
}

func TokenPartners(w rest.ResponseWriter, r *rest.Request) {
	tokenAddr := common.HexToAddress(r.PathParam("token"))
	log.Trace(fmt.Sprintf("TokenPartners tokenAddr=%s", utils.APex(tokenAddr)))
	tokens:= RaidenApi.GetTokenList()
	found:=false
	for _,t:=range tokens{
		if t==tokenAddr{
			found=true
			break
		}
	}
	if !found{
		rest.Error(w,"token doesn't exist",http.StatusNotFound)
		return
	}
	chs, err := RaidenApi.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var datas []*partnersData
	for _, c := range chs {
		d := &partnersData{
			PartnerAddress: c.PartnerAddress.String(),
			Channel:        "api/1/channles/" + c.ChannelAddress.String(),
		}
		datas = append(datas, d)
	}
	w.WriteJson(datas)
}
