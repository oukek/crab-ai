package _type

import (
	"errors"
	"fmt"
)

// Error represents a standardized API error format.
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"` // e.g., "invalid_request_error", "api_error"
	Param   string `json:"param,omitempty"`
	Code    any    `json:"code,omitempty"` // Can be string or int depending on API
}

// Implement the error interface
func (e *Error) Error() string {
	return fmt.Sprintf("API Error: type=%s code=%v message=%s param=%s", e.Type, e.Code, e.Message, e.Param)
}

// NewError creates a new API Error
func NewError(param, message, errType string, code any) *Error {
	return &Error{
		Param:   param,
		Message: message,
		Type:    errType,
		Code:    code,
	}
}

// 429 error
var ErrResourceExhausted = errors.New("系统繁忙，请稍后再试")

// prompt被拦截
var ErrPromptBlocked = errors.New("系统检测到您的输入包含敏感内容，请重新输入")

// Method not implemented in provider
var ErrMethodNotImplemented = errors.New("方法未在provider中实现")
