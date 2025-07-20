package getui

import (
	"testing"
	"time"
)

// 从.env文件获取测试配置
func getTestConfig() *Config {
	return LoadConfigFromEnvFileOrDefault(".env")
}

func TestNewClient(t *testing.T) {
	// 测试创建客户端
	config := getTestConfig()

	client := NewClient(config)
	if client == nil {
		t.Fatal("client should not be nil")
	}

	if client.config != config {
		t.Error("client config should match input config")
	}

	if client.PushAPI == nil {
		t.Error("PushAPI should not be nil")
	}

	if client.UserAPI == nil {
		t.Error("UserAPI should not be nil")
	}

	if client.StatisticAPI == nil {
		t.Error("StatisticAPI should not be nil")
	}
}

func TestNewDefaultConfig(t *testing.T) {
	config := NewDefaultConfig()
	if config == nil {
		t.Fatal("default config should not be nil")
	}

	if config.Domain != "https://restapi.getui.com/v2" {
		t.Errorf("expected domain %s, got %s", "https://restapi.getui.com/v2", config.Domain)
	}

	if config.SocketTimeout != 30000 {
		t.Errorf("expected socket timeout %d, got %d", 30000, config.SocketTimeout)
	}

	if config.ConnectTimeout != 10000 {
		t.Errorf("expected connect timeout %d, got %d", 10000, config.ConnectTimeout)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr error
	}{
		{
			name: "valid config",
			config: &Config{
				AppID:        "test_app_id",
				AppKey:       "test_app_key",
				MasterSecret: "test_master_secret",
				Domain:       "https://restapi.getui.com/v2",
			},
			wantErr: nil,
		},
		{
			name: "missing app_id",
			config: &Config{
				AppKey:       "test_app_key",
				MasterSecret: "test_master_secret",
				Domain:       "https://restapi.getui.com/v2",
			},
			wantErr: ErrAppIDRequired,
		},
		{
			name: "missing app_key",
			config: &Config{
				AppID:        "test_app_id",
				MasterSecret: "test_master_secret",
				Domain:       "https://restapi.getui.com/v2",
			},
			wantErr: ErrAppKeyRequired,
		},
		{
			name: "missing master_secret",
			config: &Config{
				AppID:  "test_app_id",
				AppKey: "test_app_key",
				Domain: "https://restapi.getui.com/v2",
			},
			wantErr: ErrMasterSecretRequired,
		},
		{
			name: "missing domain",
			config: &Config{
				AppID:        "test_app_id",
				AppKey:       "test_app_key",
				MasterSecret: "test_master_secret",
			},
			wantErr: ErrDomainRequired,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if err != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGenerateRequestID(t *testing.T) {
	config := getTestConfig()
	client := NewClient(config)

	requestID1 := client.GenerateRequestID()

	// 添加小延迟确保时间戳不同
	time.Sleep(1 * time.Millisecond)

	requestID2 := client.GenerateRequestID()

	if requestID1 == "" {
		t.Error("request ID should not be empty")
	}

	if requestID1 == requestID2 {
		t.Error("request IDs should be different")
	}
}

func TestGetCustomSocketTimeout(t *testing.T) {
	config := NewDefaultConfig()
	config.URIToSocketTimeoutMap["/test/uri"] = 5000

	timeout := config.GetCustomSocketTimeout("/test/uri")
	if timeout != 5000 {
		t.Errorf("expected timeout %d, got %d", 5000, timeout)
	}

	timeout = config.GetCustomSocketTimeout("/unknown/uri")
	if timeout != config.SocketTimeout {
		t.Errorf("expected timeout %d, got %d", config.SocketTimeout, timeout)
	}
}

func TestApiResultIsSuccess(t *testing.T) {
	tests := []struct {
		name      string
		apiResult *ApiResult
		expected  bool
	}{
		{
			name: "success result",
			apiResult: &ApiResult{
				Code: 0,
				Msg:  "success",
			},
			expected: true,
		},
		{
			name: "failure result",
			apiResult: &ApiResult{
				Code: 1001,
				Msg:  "error",
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.apiResult.IsSuccess()
			if result != tt.expected {
				t.Errorf("ApiResult.IsSuccess() = %v, want %v", result, tt.expected)
			}
		})
	}
}
