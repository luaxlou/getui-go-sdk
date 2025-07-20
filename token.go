package getui

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// TokenManager 令牌管理器
type TokenManager struct {
	config          *Config
	httpClient      *http.Client
	token           string
	tokenExpireTime time.Time
}

// NewTokenManager 创建新的令牌管理器
func NewTokenManager(config *Config, httpClient *http.Client) *TokenManager {
	return &TokenManager{
		config:     config,
		httpClient: httpClient,
	}
}

// GetToken 获取认证token
func (tm *TokenManager) GetToken() (string, error) {
	// 检查token是否过期
	if tm.token != "" && time.Now().Before(tm.tokenExpireTime) {
		return tm.token, nil
	}

	// 生成新的token
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sign := tm.generateSign(timestamp)

	authDTO := &AuthDTO{
		Sign:      sign,
		Timestamp: timestamp,
		AppKey:    tm.config.AppKey,
	}

	url := fmt.Sprintf("%s/%s/auth", tm.config.Domain, tm.config.AppID)

	body, err := json.Marshal(authDTO)
	if err != nil {
		return "", &NetworkError{Message: "failed to marshal auth request", Cause: err}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return "", &NetworkError{Message: "failed to create auth request", Cause: err}
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := tm.httpClient.Do(req)
	if err != nil {
		return "", &NetworkError{Message: "failed to send auth request", Cause: err}
	}
	defer resp.Body.Close()

	var result ApiResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", &NetworkError{Message: "failed to decode auth response", Cause: err}
	}

	if !result.IsSuccess() {
		return "", &APIError{Code: result.Code, Message: result.Msg}
	}

	// 解析token
	var tokenData struct {
		Token string `json:"token"`
	}
	if err := json.Unmarshal(result.Data, &tokenData); err != nil {
		return "", &NetworkError{Message: "failed to parse token", Cause: err}
	}

	tm.token = tokenData.Token
	tm.tokenExpireTime = time.Now().Add(23 * time.Hour) // token有效期24小时，提前1小时刷新

	return tm.token, nil
}

// generateSign 生成签名
func (tm *TokenManager) generateSign(timestamp string) string {
	// 签名算法：SHA256(appkey + timestamp + master_secret)
	signStr := tm.config.AppKey + timestamp + tm.config.MasterSecret
	return fmt.Sprintf("%x", sha256.Sum256([]byte(signStr)))
}

// GetTokenExpireTime 获取token过期时间
func (tm *TokenManager) GetTokenExpireTime() time.Time {
	return tm.tokenExpireTime
}

// GetCurrentToken 获取当前token（不检查过期）
func (tm *TokenManager) GetCurrentToken() string {
	return tm.token
}

// SetToken 设置token（用于测试）
func (tm *TokenManager) SetToken(token string, expireTime time.Time) {
	tm.token = token
	tm.tokenExpireTime = expireTime
}

// ClearToken 清除token（用于测试）
func (tm *TokenManager) ClearToken() {
	tm.token = ""
	tm.tokenExpireTime = time.Time{}
}

// IsTokenExpired 检查token是否过期
func (tm *TokenManager) IsTokenExpired() bool {
	return tm.token == "" || time.Now().After(tm.tokenExpireTime)
}
