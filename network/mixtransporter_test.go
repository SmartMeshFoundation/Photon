package network

import (
	"testing"
)

func TestNewMixDiscovery(t *testing.T) {

	NewMixTranspoter("test", "127.0.0.0.1", 5001, nil, &dummyPolicy{})

}
