package lansenger

import (
	"errors"
	"strconv"
)

type LansengerError struct {
	Message   string
	ErrCode   int
	Retryable bool
}

func (e *LansengerError) Error() string {
	if e.ErrCode != 0 {
		return e.Message + " (errCode: " + strconv.Itoa(e.ErrCode) + ")"
	}
	return e.Message
}

// Code 暴露服务端 errCode，供 ErrCodeOf 通过 errors.As 读取（含嵌入类型）。
func (e *LansengerError) Code() int {
	return e.ErrCode
}

type AuthError struct {
	LansengerError
}

type ConfigError struct {
	LansengerError
}

type APIError struct {
	LansengerError
}

type NetworkError struct {
	LansengerError
}

type FileError struct {
	LansengerError
}

func NewAuthError(message string) *AuthError {
	return &AuthError{LansengerError{Message: message, Retryable: false}}
}

func NewConfigError(message string) *ConfigError {
	return &ConfigError{LansengerError{Message: message, Retryable: false}}
}

func NewAPIError(message string, errCode int) *APIError {
	retryable := errCode != 0
	return &APIError{LansengerError{Message: message, ErrCode: errCode, Retryable: retryable}}
}

func NewNetworkError(message string) *NetworkError {
	return &NetworkError{LansengerError{Message: message, Retryable: true}}
}

func NewFileError(message string) *FileError {
	return &FileError{LansengerError{Message: message, Retryable: false}}
}

// ErrCodeOf returns the server errCode carried by err, or 0 when err carries no
// server code (local validation, transport failure, request shape errors).
//
// Result values already mirror this in their ErrCode field; use this helper when
// building a result from an error.
func ErrCodeOf(err error) int {
	if err == nil {
		return 0
	}
	var coder interface{ Code() int }
	if errors.As(err, &coder) {
		return coder.Code()
	}
	return 0
}
