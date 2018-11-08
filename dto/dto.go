package dto

import (
	"encoding/json"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
)

// APIErrorCode :
type APIErrorCode int

const (
	// SUCCESS :
	SUCCESS = 0
	// EXCEPTION :
	EXCEPTION = 99
)

// APIResponse :
// 接口统一返回格式
type APIResponse struct {
	ErrorCode APIErrorCode `json:"error_code"`
	ErrorMsg  string       `json:"error_msg"`
	CallID    string       `json:"call_id,omitempty"` // for async
	Data      interface{}  `json:"data,omitempty"`
}

// NewAPIResponse :
func NewAPIResponse(errorCode APIErrorCode, errorMsg string, data interface{}) *APIResponse {
	return &APIResponse{
		ErrorCode: errorCode,
		ErrorMsg:  errorMsg,
		Data:      data,
	}
}

// NewSuccessAPIResponse :
func NewSuccessAPIResponse(data interface{}) *APIResponse {
	return &APIResponse{
		ErrorCode: SUCCESS,
		ErrorMsg:  "SUCCESS",
		Data:      data,
	}
}

// NewExceptionAPIResponse :
func NewExceptionAPIResponse(err error) *APIResponse {
	return &APIResponse{
		ErrorCode: EXCEPTION,
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
