package dto

import (
	"encoding/json"
	"errors"

	"github.com/SmartMeshFoundation/Photon/rerr"

	"fmt"

	"github.com/SmartMeshFoundation/Photon/log"
)

const (
	// SUCCESS 成功通用Code
	SUCCESS = 0
	// UNKNOWNERROR 未知错误
	UNKNOWNERROR = -1
)

// APIResponse 接口统一返回格式
type APIResponse struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_message"`
	//CallID    string          `json:"call_id,omitempty"` // for async
	Data json.RawMessage `json:"data,omitempty"`
}

//API2JSON helper function
func API2JSON(a *APIResponse) string {
	b, err := json.Marshal(a)
	if err != nil {
		panic(err)
	}
	return string(b)
}

//NewAPIResponse 辅助http接口创建
func NewAPIResponse(err error, data interface{}) *APIResponse {
	if err != nil {
		return NewExceptionAPIResponse(err)
	}
	return NewSuccessAPIResponse(data)
}

//NewMobileResponse mobile接口,err为空认为成功,否则认为失败
func NewMobileResponse(err error, data interface{}) string {
	return API2JSON(NewAPIResponse(err, data))
}

//NewSuccessMobileResponse 直接序列化为string,方便处理
func NewSuccessMobileResponse(data interface{}) string {
	a := NewSuccessAPIResponse(data)
	return API2JSON(a)
}

//NewErrorMobileResponse 直接序列化为string,方便处理
func NewErrorMobileResponse(err error) string {
	a := NewExceptionAPIResponse(err)
	return API2JSON(a)
}

// NewSuccessAPIResponse  http接口所用,创建成功的返回
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

// NewExceptionAPIResponse http接口用,创建失败的返回,如果err为空会被认为是成功
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

// String fmt.Formater
func (r *APIResponse) String() string {
	buf, err := json.Marshal(r)
	if err != nil {
		log.Error(fmt.Sprintf("APIResponse marshal err = %s", err.Error()))
		return ""
	}
	return string(buf)
}

// ToFormatString 打印格式化的json结构体
func (r *APIResponse) ToFormatString() string {
	buf, err := json.MarshalIndent(r, "\t", "")
	if err != nil {
		log.Error(fmt.Sprintf("APIResponse marshal err = %s", err.Error()))
		return ""
	}
	return string(buf)
}

//ParseResult helper function
func ParseResult(result string, output interface{}) (err error) {
	var res APIResponse
	err = json.Unmarshal([]byte(result), &res)
	if err != nil {
		panic(err)
	}
	if res.ErrorCode != SUCCESS {
		return errors.New(res.ErrorMsg)
	}
	if output != nil {
		err = json.Unmarshal([]byte(res.Data), output)
	}
	return nil
}
