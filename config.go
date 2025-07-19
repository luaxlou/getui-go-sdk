package getui

import (
	"crypto/tls"
	"net/http"
	"time"
)

// Config 个推SDK配置
type Config struct {
	// 应用配置
	AppID        string `json:"app_id"`
	AppKey       string `json:"app_key"`
	MasterSecret string `json:"master_secret"`
	Domain       string `json:"domain"`

	// HTTP配置
	SocketTimeout            int  `json:"socket_timeout"`             // HTTP读取超时时间(ms)
	ConnectTimeout           int  `json:"connect_timeout"`            // HTTP连接超时时间(ms)
	ConnectionRequestTimeout int  `json:"connection_request_timeout"` // 从连接池获取连接超时时间(ms)
	MaxHTTPTryTime           int  `json:"max_http_try_time"`          // HTTP重试次数
	TrustSSL                 bool `json:"trust_ssl"`                  // 是否信任SSL证书

	// 域名检测配置
	OpenAnalyseStableDomain     bool          `json:"open_analyse_stable_domain"`     // 是否开启稳定域名检测
	AnalyseStableDomainInterval time.Duration `json:"analyse_stable_domain_interval"` // 检测稳定域名时间间隔
	MaxFailedNum                int           `json:"max_failed_num"`                 // 最大失败次数阈值
	ContinuousFailedNum         int           `json:"continuous_failed_num"`          // 连续失败次数阈值
	CheckMaxFailedNumInterval   time.Duration `json:"check_max_failed_num_interval"`  // 重置最大失败次数的时间间隔
	HTTPCheckTimeout            int           `json:"http_check_timeout"`             // 域名检测超时时间(ms)

	// 健康检测配置
	OpenCheckHealthDataSwitch bool          `json:"open_check_health_data_switch"` // 是否开启健康检测
	CheckHealthInterval       time.Duration `json:"check_health_interval"`         // 健康检测时间间隔

	// 代理配置
	ProxyConfig *HTTPProxyConfig `json:"proxy_config"`

	// 自定义超时配置
	URIToSocketTimeoutMap map[string]int `json:"uri_to_socket_timeout_map"` // URI到超时时间的映射
}

// HTTPProxyConfig HTTP代理配置
type HTTPProxyConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewDefaultConfig 创建默认配置
func NewDefaultConfig() *Config {
	return &Config{
		Domain:                      "https://restapi.getui.com/v2",
		SocketTimeout:               30000,
		ConnectTimeout:              10000,
		ConnectionRequestTimeout:    0,
		MaxHTTPTryTime:              1,
		TrustSSL:                    false,
		OpenAnalyseStableDomain:     true,
		AnalyseStableDomainInterval: 2 * time.Minute,
		MaxFailedNum:                10,
		ContinuousFailedNum:         3,
		CheckMaxFailedNumInterval:   3 * time.Second,
		HTTPCheckTimeout:            100,
		OpenCheckHealthDataSwitch:   false,
		CheckHealthInterval:         30 * time.Second,
		URIToSocketTimeoutMap:       make(map[string]int),
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.AppID == "" {
		return ErrAppIDRequired
	}
	if c.AppKey == "" {
		return ErrAppKeyRequired
	}
	if c.MasterSecret == "" {
		return ErrMasterSecretRequired
	}
	if c.Domain == "" {
		return ErrDomainRequired
	}
	return nil
}

// GetHTTPClient 获取HTTP客户端
func (c *Config) GetHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.TrustSSL,
		},
		DisableKeepAlives: false,
		IdleConnTimeout:   30 * time.Second,
	}

	// 设置代理
	if c.ProxyConfig != nil {
		// 这里可以添加代理配置逻辑
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(c.SocketTimeout) * time.Millisecond,
	}

	return client
}

// GetCustomSocketTimeout 获取自定义超时时间
func (c *Config) GetCustomSocketTimeout(uri string) int {
	if timeout, exists := c.URIToSocketTimeoutMap[uri]; exists {
		return timeout
	}
	return c.SocketTimeout
}
