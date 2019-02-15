package v1

import (
	"fmt"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"github.com/SmartMeshFoundation/Photon/dto"
	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ant0ine/go-json-rest/rest"
)

/*
Tokens is api of /api/1/tokens
*/
func Tokens(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> Tokens ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	resp = dto.NewSuccessAPIResponse(API.GetTokenTokenNetorks())
}

/*
TokenPartners is api of /api/1/:token/:partner
*/
func TokenPartners(w rest.ResponseWriter, r *rest.Request) {
	var resp *dto.APIResponse
	defer func() {
		log.Trace(fmt.Sprintf("Restful Api Call ----> TokenPartners ,err=%s", resp.ToFormatString()))
		writejson(w, resp)
	}()
	type partnersDataResponse struct {
		PartnerAddress string `json:"partner_address"`
		Channel        string `json:"channel"`
	}
	tokenAddr, err := utils.HexToAddress(r.PathParam("token"))
	if err != nil {
		resp = dto.NewExceptionAPIResponse(rerr.ErrArgumentError.AppendError(err))
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
		resp = dto.NewExceptionAPIResponse(rerr.ErrTokenNotFound)
		return
	}
	chs, err := API.GetChannelList(tokenAddr, utils.EmptyAddress)
	if err != nil {
		resp = dto.NewExceptionAPIResponse(err)
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
	resp = dto.NewSuccessAPIResponse(datas)
}
