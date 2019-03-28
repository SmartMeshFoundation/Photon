package debug

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSetup(t *testing.T) {
	req := require.New(t)
	m1 := mac()
	req.EqualValues(m1, mac())
	req.EqualValues(m1, mac())
	req.EqualValues(m1, mac())
	req.EqualValues(m1, mac())
	req.EqualValues(m1, mac())
}
