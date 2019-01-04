package v1

import (
	"fmt"
	"net/http"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
Tokens is api of /api/1/tokens
*/
func Tokens(w rest.ResponseWriter, r *rest.Request) {
	err := w.WriteJson(API.GetTokenTokenNetorks())
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}

/*
TokenPartners is api of /api/1/:token/:partner
*/
func TokenPartners(w rest.ResponseWriter, r *rest.Request) {
	type partnersDataResponse struct {
		PartnerAddress string `json:"partner_address"`
		Channel        string `json:"channel"`
	}
	tokenAddr, err := utils.HexToAddress(r.PathParam("token"))
	if err != nil {
		log.Error(err.Error())
		rest.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Trace(fmt.Sprintf("TokenPartners tokenAddr=%s", utils.APex(tokenAddr)))
	tokens := API.GetTokenList()
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
	chs, err := API.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var datas []*partnersDataResponse
	for _, c := range chs {
		d := &partnersDataResponse{
			PartnerAddress: c.PartnerAddress().String(),
			Channel:        "api/1/channles/" + c.ChannelIdentifier.ChannelIdentifier.String(),
		}
		datas = append(datas, d)
	}
	err = w.WriteJson(datas)
	if err != nil {
		log.Warn(fmt.Sprintf("writejson err %s", err))
	}
}
