package getui

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client 个推SDK客户端
type Client struct {
	config       *Config
	httpClient   *http.Client
	tokenManager *TokenManager

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

	httpClient := config.GetHTTPClient()
	client := &Client{
		config:       config,
		httpClient:   httpClient,
		tokenManager: NewTokenManager(config, httpClient),
	}

	// 初始化API接口
	client.PushAPI = &PushAPI{client: client}
	client.UserAPI = &UserAPI{client: client}
	client.StatisticAPI = &StatisticAPI{client: client}

	return client
}

// GetToken 获取认证token
func (c *Client) GetToken() (string, error) {
	return c.tokenManager.GetToken()
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

// GetConfig 获取客户端配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// GetTokenManager 获取令牌管理器（用于测试）
func (c *Client) GetTokenManager() *TokenManager {
	return c.tokenManager
}
