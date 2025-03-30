package photon

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SmartMeshFoundation/Photon/log"
	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/params"
	"github.com/SmartMeshFoundation/Photon/utils"
)

/*
ReportService 负责photon对外部系统的通知
*/
type ReportService struct {
	photon *Service
}

/*
NewReportService constructor
*/
func NewReportService(photon *Service) *ReportService {
	return &ReportService{
		photon: photon,
	}
}

/*
Start 主循环,目前就一个付费上网的上报需求,对接的系统尚未开发,先采用主动上报的模式
后续应改为开放一个websocket端口,走订阅模式
*/
func (s *ReportService) Start() {
	if params.Cfg.ReceivedTransferReportURL == "" || s.photon == nil {
		return
	}
	go func() {
		log.Info(" ReportService Start...")
		receivedTransferChan := s.photon.NotifyHandler.GetReceivedTransferChan()
		for {
			select {
			case rt, ok := <-receivedTransferChan:
				if !ok {
					// 该通道关闭说明photon stop了
					log.Info("ReportService exit")
					return
				}
				s.reportReceiveTransfer(rt)
			}
		}
	}()
}

func (s *ReportService) reportReceiveTransfer(rt *models.ReceivedTransfer) {
	if params.Cfg.ReceivedTransferReportURL == "" {
		return
	}
	url := params.Cfg.ReceivedTransferReportURL
	req := &utils.Req{
		FullURL: url,
		Method:  http.MethodPost,
		Payload: utils.Marshal(rt),
		Timeout: time.Second * 10,
	}
	statusCode, body, err := req.Invoke()
	if err != nil {
		log.Error(fmt.Sprintf("ReportService report received transfer to %s err :%s", url, err))
		return
	}
	if statusCode != 200 {
		log.Error(fmt.Sprintf("ReportService report received transfer to %s err : http status=%d body=%s", url, statusCode, string(body)))
		return
	}
	log.Info(fmt.Sprintf("ReportService eport received transfer to %s SUCCESS", url))
}
