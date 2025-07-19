package getui

import (
	"errors"
	"fmt"
)

// 配置相关错误
var (
	ErrAppIDRequired        = errors.New("app_id is required")
	ErrAppKeyRequired       = errors.New("app_key is required")
	ErrMasterSecretRequired = errors.New("master_secret is required")
	ErrDomainRequired       = errors.New("domain is required")
)

// API相关错误
var (
	ErrInvalidRequestID = errors.New("request_id must be between 10-32 characters")
	ErrEmptyAudience    = errors.New("audience cannot be empty")
	ErrEmptyPushMessage = errors.New("push_message cannot be empty")
	ErrInvalidCID       = errors.New("cid cannot be empty")
	ErrInvalidAlias     = errors.New("alias cannot be empty")
)

// HTTP相关错误
var (
	ErrHTTPRequestFailed = errors.New("http request failed")
	ErrInvalidResponse   = errors.New("invalid response")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrRateLimited       = errors.New("rate limited")
)

// 认证相关错误
var (
	ErrTokenExpired = errors.New("token expired")
	ErrInvalidToken = errors.New("invalid token")
)

// 自定义错误类型
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("API error: code=%d, message=%s", e.Code, e.Message)
}

// 网络错误
type NetworkError struct {
	Message string
	Cause   error
}

func (e *NetworkError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("network error: %s, cause: %v", e.Message, e.Cause)
	}
	return fmt.Sprintf("network error: %s", e.Message)
}

func (e *NetworkError) Unwrap() error {
	return e.Cause
}

// 配置错误
type ConfigError struct {
	Field   string
	Message string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("config error: field=%s, message=%s", e.Field, e.Message)
}
