package mobile

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/SmartMeshFoundation/Photon/dto"

	"github.com/SmartMeshFoundation/Photon/models/stormdb"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"
	"github.com/ethereum/go-ethereum/common"
)

//SimpleAPI 不启动photon就可以查询通道信息
type SimpleAPI struct {
	dao models.Dao
}

//NewSimpleAPI 创建数据库访问接口
func NewSimpleAPI(datadir, address string) (api *SimpleAPI, err error) {
	addr := common.HexToAddress(address)
	userDbPath := hex.EncodeToString(addr[:])
	userDbPath = userDbPath[:8]
	userDbPath = filepath.Join(datadir, userDbPath)
	if !utils.Exists(userDbPath) {
		err = os.MkdirAll(userDbPath, os.ModePerm)
		if err != nil {
			err = fmt.Errorf("datadir:%s doesn't exist and cannot create %v", userDbPath, err)
			return
		}
	}
	databasePath := filepath.Join(userDbPath, "log.db")
	dao, err := stormdb.OpenDb(databasePath)
	if err != nil {
		err = fmt.Errorf("open db error %s", err)
		return
	}
	api = &SimpleAPI{
		dao: dao,
	}
	return
}

//Stop 关闭数据库
func (a *SimpleAPI) Stop() {
	a.dao.CloseDB()
}

//BalanceAvailabelOnPhoton 查询某个token在整个token上的可用金额
func (a *SimpleAPI) BalanceAvailabelOnPhoton(token string) (result string) {
	tokenAddress := common.HexToAddress(token)

	channels, err := a.dao.GetChannelList(tokenAddress, utils.EmptyAddress)
	if err != nil {
		dto.NewErrorMobileResponse(err)
		return
	}
	v := big.NewInt(0)

	for _, channel := range channels {
		v = v.Add(v, channel.OurBalance())
	}
	return dto.NewSuccessMobileResponse(v)
}
