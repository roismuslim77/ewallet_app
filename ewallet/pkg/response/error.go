package response

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code       string      `json:"code"`
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	IsNoError  bool        `json:"-"`
}

func Error(code string) *ErrorResponse {
	msg := Code(code)

	res := &ErrorResponse{
		Code:       code,
		StatusCode: http.StatusBadRequest,
		Message:    msg.Name(),
	}

	return res
}

func (e *ErrorResponse) WithStatusCode(statusCode int) *ErrorResponse {
	e.StatusCode = statusCode
	return e
}

func (e *ErrorResponse) WithError(err string) *ErrorResponse {
	e.Message = err
	return e
}

func (e *ErrorResponse) WithArgsMessage(args ...interface{}) *ErrorResponse {
	e.Message = fmt.Sprintf(e.Message, args...)
	return e
}

func (e *ErrorResponse) WithData(data interface{}) *ErrorResponse {
	e.Data = data
	return e
}

func NotError() *ErrorResponse {
	return &ErrorResponse{
		IsNoError: true,
	}
}
