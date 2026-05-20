package lansenger

import "strconv"

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
