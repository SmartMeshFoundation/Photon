package mobile

import (
	"testing"

	"github.com/SmartMeshFoundation/Photon/dto"

	"github.com/stretchr/testify/assert"
)

var testETHENDPOINT = "http://192.168.124.13:5554"
var testCONTRACT = "0xeb312827167065654509758956EEc44f4FCA058C"
var testAddress = "0x292650fee408320d888e06ed89d938294ea42f99"
var testKeystore = "../testdata/mykeystore"

func TestSNM(t *testing.T) {
	if testing.Short() {
		return
	}
	ast := assert.New(t)
	s, err := NewSNM(testAddress, testKeystore, testETHENDPOINT, "123", testCONTRACT, "")
	ast.Nil(err)
	result := s.AddFunds("10000000")
	err = dto.ParseResult(result, nil)
	ast.Nil(err)
	result = s.TryStopContract()
	err = dto.ParseResult(result, nil)
	ast.NotNil(err)

}
