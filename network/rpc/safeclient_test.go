package rpc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/SmartMeshFoundation/Photon/models"
	"github.com/SmartMeshFoundation/Photon/utils"

	"github.com/davecgh/go-spew/spew"
)

func TestBrokenClient(t *testing.T) {
	bcs := MakeTestBlockChainService()
	_, err := bcs.Client.BalanceAt(context.Background(), bcs.NodeAddress, nil)
	if err != nil {
		t.Error(err)
		return

	}
	fmt.Println("shutdown geth now...")
	time.Sleep(5 * time.Second)
	fmt.Println("try  operation on broken connection")
	_, err = bcs.Client.BalanceAt(context.Background(), bcs.NodeAddress, nil)
	if err != nil {
		spew.Dump(err)
		t.Error(err)
	}

}

//0xc479184abeb8c508ee96e4c093ee47af2256cbbf registry合约地址
//公链地址: http://transport01.smartmesh.cn:17888
func TestCallOperativeSettle(t *testing.T) {
	if testing.Short() {
		return
	}
	var txparams models.ChannelCooperativeSettleTXParams
	s := `{"token_address":"0xf0123c3267af5cbbfab985d39171f5f5758c0900","p1_address":"0xc0692451a239c44d596867db56a54bf6cd61a10c","p1_balance":18997400000000000000,"p2_address":"0xbbad695e60d8c3b50bafd78bd9400522ff14c95d","p2_balance":81002600000000000000,"p1_signature":"BimOlYEdp4axVgm09KsJDUCEXx0QxEMDV0AQUnX7XrFs47/7UbhlyA5BEN2Twx3skDmFvGzLOggQtDfHyVXicRw=","p2_signature":"UAB1p4vA33rjITKeaAaMT8SukZ0udh3CR3+E6G64vY8SqftgbK0l+dfyR6IidDMY3ofv263IlWw2w/Kv+EyhFxw="}`
	err := json.Unmarshal([]byte(s), &txparams)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tx=%s", utils.StringInterface(txparams, 7))
	bcs := MakeTestBlockChainService()
	tn, err := bcs.TokenNetwork(bcs.tokenNetworkAddress)
	if err != nil {
		t.Error(err)
		return
	}
	err = tn.CooperativeSettle(txparams.P1Address, txparams.P2Address, txparams.P1Balance, txparams.P2Balance,
		txparams.P1Signature, txparams.P2Signature)
	if err != nil {
		t.Error(err)
	}
}
