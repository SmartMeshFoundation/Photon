package dto

import (
	"encoding/json"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
)

// APIErrorCode :
type APIErrorCode int

const (
	// SUCCESS :
	SUCCESS = 0
	// EXCEPTION :
	UNKNOWNERROR = -1
)

// APIResponse :
// 接口统一返回格式
type APIResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
	//CallID    string          `json:"call_id,omitempty"` // for async
	Data json.RawMessage `json:"data,omitempty"`
}

func apitojson(a *APIResponse) string {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(b)
}

//NewSuccessMobileResponse 直接序列化为string,方便处理
func NewSuccessMobileResponse(data interface{}) string {
	a := NewSuccessAPIResponse(data)
	return apitojson(a)
}

//NewErrorMobileResponse 直接序列化为string,方便处理
func NewErrorMobileResponse(err error) string {
	a := NewExceptionAPIResponse(err)
	return apitojson(a)
}

// NewSuccessAPIResponse :
func NewSuccessAPIResponse(data interface{}) *APIResponse {
	d, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return &APIResponse{
		ErrorCode: SUCCESS,
		ErrorMsg:  "SUCCESS",
		Data:      json.RawMessage(d),
	}
}

// NewExceptionAPIResponse :
func NewExceptionAPIResponse(err error) *APIResponse {
	if err == nil {
		return NewSuccessAPIResponse(nil)
	}
	e1, ok := err.(rerr.StandardError)
	if ok {
		return &APIResponse{
			ErrorCode: e1.ErrorCode,
			ErrorMsg:  e1.ErrorMsg,
		}
	}
	e2, ok := err.(rerr.StandardDataError)
	if ok {
		return &APIResponse{
			ErrorCode: e2.ErrorCode,
			ErrorMsg:  e2.ErrorMsg,
			Data:      e2.Data,
		}
	}
	return &APIResponse{
		ErrorCode: UNKNOWNERROR,
		ErrorMsg:  err.Error(),
	}
}

// ToString :
func (r *APIResponse) ToString() string {
	buf, err := json.Marshal(r)
	if err != nil {
		log.Error(fmt.Sprintf("APIResponse marshal err = %s", err.Error()))
		return ""
	}
	return string(buf)
}

// ToFormatString :
func (r *APIResponse) ToFormatString() string {
	buf, err := json.MarshalIndent(r, "\t", "")
	if err != nil {
		log.Error(fmt.Sprintf("APIResponse marshal err = %s", err.Error()))
		return ""
	}
	return string(buf)
}
