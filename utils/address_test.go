package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHex2Address(t *testing.T) {
	addressDefaultValidation = true
	cases := map[string]bool{
		//All caps
		"0x52908400098527886E0F7030069857D2E4169EE7": true,
		"0x8617E340B3D01FA5F11F306F4090FD50E238070D": true,
		//All Lower
		"0xde709f2102306220921060314715629080e2fb77": true,
		"0x27b1fdb04752bbc536007a920d24acb045561c26": true,
		//# Normal
		"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed": true,
		"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359": true,
		"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB": true,
		"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb": true,
		//modified
		//All caps
		"0x52908400098527886E0F7030069857D2E4169Ee7": false,
		"0x8617E340B3D01FA5F11F306F4090FD50E238070d": false,
		//All Lower
		"0xde709f2102306220921060314715629080e2fB77": false,
		"0x27B1fdb04752bbc536007a920d24acb045561c26": false,
		//# Normal
		"0x5aaeb6053F3E94C9b9A09f33669435E7Ef1BeAed": false,
		"0xfb6916095ca1df60bB79Ce92cE3Ea74c37c5d359": false,
		"0xdBf03B407c01E7cD3CBea99509d93f8DDDC8C6FB": false,
		"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDB": false,
	}
	for addr, isRight := range cases {
		_, err := HexToAddress(addr)
		if isRight {
			assert.EqualValues(t, err, nil)
		} else {
			assert.NotEqual(t, err, nil)
		}
	}
}
