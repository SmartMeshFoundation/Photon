package pmsproxy

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

// ErrConnect :
var ErrConnect = errors.New("pmsClient connect to pms error")

type pmsClient struct {
	host        string
	selfAddress common.Address
	pmsAddress  common.Address
}

// NewPmsProxy init
func NewPmsProxy(pmsHost string, selfAddress common.Address, pmsAddress common.Address) PmsProxy {
	return &pmsClient{
		host:        pmsHost,
		selfAddress: selfAddress,
		pmsAddress:  pmsAddress,
	}
}

// SubmitDelegate 向pfs提交一个通道的所有委托,包含UpdateBalanceProof,Unlock及Punish三种
func (c *pmsClient) SubmitDelegate(data *DelegateForPms) (err error) {
	if data == nil {
		return
	}
	req := &utils.Req{
		FullURL: c.host + "/delegate/" + c.selfAddress.String(),
		Method:  http.MethodPost,
		Payload: utils.Marshal(data),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		return ErrConnect
	}
	if statusCode != 200 {
		err = fmt.Errorf("PmsAPI SubmitDelegate %s err : http status=%d body=%s", req.FullURL, statusCode, string(body))
		return
	}
	log.Info(fmt.Sprintf("PmsAPI SubmitDelegate of channel %s SUCCESS", data.ChannelIdentifier.String()))
	return nil
}
