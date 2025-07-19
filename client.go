package getui

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client 个推SDK客户端
type Client struct {
	config          *Config
	httpClient      *http.Client
	token           string
	tokenExpireTime time.Time

	// API接口
	PushAPI      *PushAPI
	UserAPI      *UserAPI
	StatisticAPI *StatisticAPI
}

// NewClient 创建新的客户端
func NewClient(config *Config) *Client {
	if config == nil {
		config = NewDefaultConfig()
	}

	// 验证配置
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("invalid config: %v", err))
	}

	client := &Client{
		config:     config,
		httpClient: config.GetHTTPClient(),
	}

	// 初始化API接口
	client.PushAPI = &PushAPI{client: client}
	client.UserAPI = &UserAPI{client: client}
	client.StatisticAPI = &StatisticAPI{client: client}

	return client
}

// GetToken 获取认证token
func (c *Client) GetToken() (string, error) {
	// 检查token是否过期
	if c.token != "" && time.Now().Before(c.tokenExpireTime) {
		return c.token, nil
	}

	// 生成新的token
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sign := c.generateSign(timestamp)

	authDTO := &AuthDTO{
		Sign:      sign,
		Timestamp: timestamp,
		AppKey:    c.config.AppKey,
	}

	url := fmt.Sprintf("%s/%s/auth", c.config.Domain, c.config.AppID)

	body, err := json.Marshal(authDTO)
	if err != nil {
		return "", &NetworkError{Message: "failed to marshal auth request", Cause: err}
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return "", &NetworkError{Message: "failed to create auth request", Cause: err}
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	resp, err := c.httpClient.Do(req)
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

	c.token = tokenData.Token
	c.tokenExpireTime = time.Now().Add(23 * time.Hour) // token有效期24小时，提前1小时刷新

	return c.token, nil
}

// generateSign 生成签名
func (c *Client) generateSign(timestamp string) string {
	// 签名算法：MD5(appkey + timestamp + master_secret)
	signStr := c.config.AppKey + timestamp + c.config.MasterSecret
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}

// DoRequest 执行HTTP请求
func (c *Client) DoRequest(method, uri string, body interface{}) (*ApiResult, error) {
	// 获取token
	token, err := c.GetToken()
	if err != nil {
		return nil, err
	}

	// 构建URL
	url := fmt.Sprintf("%s/%s%s", c.config.Domain, c.config.AppID, uri)

	// 准备请求体
	var reqBody []byte
	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, &NetworkError{Message: "failed to marshal request body", Cause: err}
		}
	}

	// 创建请求
	req, err := http.NewRequest(method, url, strings.NewReader(string(reqBody)))
	if err != nil {
		return nil, &NetworkError{Message: "failed to create request", Cause: err}
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json;charset=utf-8")
	req.Header.Set("token", token)

	// 设置自定义超时
	if customTimeout := c.config.GetCustomSocketTimeout(uri); customTimeout > 0 {
		c.httpClient.Timeout = time.Duration(customTimeout) * time.Millisecond
	}

	// 执行请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, &NetworkError{Message: "failed to send request", Cause: err}
	}
	defer resp.Body.Close()

	// 解析响应
	var result ApiResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, &NetworkError{Message: "failed to decode response", Cause: err}
	}

	return &result, nil
}

// GenerateRequestID 生成请求ID
func (c *Client) GenerateRequestID() string {
	return strconv.FormatInt(time.Now().UnixNano(), 10)
}
