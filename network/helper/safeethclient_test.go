package helper

import (
	"testing"
	"time"
)

func TestNewSafeClient_HTTP_Right(t *testing.T) {
	url := "http://192.168.124.13:5554"
	c, _ := NewSafeClient(url)
	for {
		err := checkConnectStatus(c.Client)
		if err != nil {
			c.RecoverDisconnect()
		}
		time.Sleep(1 * time.Second)
	}
}

func TestNewSafeClient_HTTP_Wrong(t *testing.T) {
	url := "http://192.168.124.13:5552"
	NewSafeClient(url)
	time.Sleep(100 * time.Second)
}

func TestNewSafeClient_WS_Right(t *testing.T) {
	url := "ws://192.168.124.13:5555"
	c, _ := NewSafeClient(url)
	for {
		err := checkConnectStatus(c.Client)
		if err != nil {
			c.RecoverDisconnect()
		}
		time.Sleep(1 * time.Second)
	}
}
