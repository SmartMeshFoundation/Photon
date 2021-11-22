package supernode

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/SmartMeshFoundation/Photon/dto"
)

// Req a photon api http request
type Req struct {
	APIName string        `json:"api_name"`
	FullURL string        `json:"url"`
	Method  string        `json:"method"`
	Payload string        `json:"payload"`
	Timeout time.Duration `json:"timeout"`
}

// GetReq : get *http.Request
func (r *Req) GetReq() *http.Request {
	var reqBody io.Reader
	if r.Payload == "" {
		reqBody = nil
	} else {
		reqBody = strings.NewReader(r.Payload)
	}
	req, err := http.NewRequest(r.Method, r.FullURL, reqBody)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "name=anny")
	return req
}

// Invoke : send a http request
func (r *Req) Invoke() ([]byte, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if r.Timeout == 0 {
					r.Timeout = time.Second * 180 // default timeout 3 min
				}
				c, err := net.DialTimeout(netw, addr, 5*time.Second) // default dial timeout 5 seconds
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(r.Timeout))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}
	req := r.GetReq()
	resp, err := client.Do(req)
	defer func() {
		var err2 error
		if req.Body != nil {
			err2 = req.Body.Close()
		}
		if resp != nil && resp.Body != nil {
			err2 = resp.Body.Close()
		}
		_ = err2
	}()
	if err != nil {
		return nil, err
	}

	var buf [4096 * 1024]byte
	n := 0
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}

	var res dto.APIResponse
	err = json.Unmarshal([]byte(buf[:n]), &res)
	if err != nil {
		panic(err)
	}
	if res.ErrorCode != dto.SUCCESS {
		err = errors.New(res.ErrorMsg)
		return nil, err
	}
	return res.Data, nil
}

// InvokeWithoutErrorCode : send a http request
func (r *Req) InvokeWithoutErrorCode() ([]byte, error) {
	client := http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				if r.Timeout == 0 {
					r.Timeout = time.Second * 180 // default timeout 3 min
				}
				c, err := net.DialTimeout(netw, addr, 5*time.Second) // default dial timeout 5 seconds
				if err != nil {
					return nil, err
				}
				err = c.SetDeadline(time.Now().Add(r.Timeout))
				if err != nil {
					return nil, err
				}
				return c, nil
			},
		},
	}
	req := r.GetReq()
	resp, err := client.Do(req)
	defer func() {
		var err2 error
		if req.Body != nil {
			err2 = req.Body.Close()
		}
		if resp != nil && resp.Body != nil {
			err2 = resp.Body.Close()
		}
		_ = err2
	}()
	if err != nil {
		return nil, err
	}
	var buf [4096 * 1024]byte
	n := 0
	n, err = resp.Body.Read(buf[:])
	if err != nil && err.Error() == "EOF" {
		err = nil
	}
	if resp.StatusCode != http.StatusOK {
		err = errors.New("not 200")
	}
	return []byte(buf[:n]), err
}

// ToString : get json of this
func (r *Req) ToString() string {
	buf, err := json.MarshalIndent(r, "", "\t")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("\n%s\n", string(buf))
}
