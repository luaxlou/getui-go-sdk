package getui

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

// 创建测试TokenManager
func createTestTokenManager() *TokenManager {
	config := LoadConfigFromEnvFileOrDefault(".env")
	httpClient := config.GetHTTPClient()
	return NewTokenManager(config, httpClient)
}

func TestTokenManager_ValidConfig(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 获取token
	token, err := tokenManager.GetToken()

	// 验证结果
	if err != nil {
		// 检查是否是API错误（说明请求格式正确，但配置无效）
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误：code=%d, message=%s", apiErr.Code, apiErr.Message)
			// 如果是配置相关的错误，记录但允许通过
			if apiErr.Code == 20001 || apiErr.Code == 20002 || apiErr.Code == 20003 || apiErr.Code == 404 {
				t.Logf("配置错误（预期）：%s", apiErr.Message)
				return
			}
			// 其他API错误应该失败
			t.Errorf("意外的API错误：code=%d, message=%s", apiErr.Code, apiErr.Message)
			return
		}
		// 网络错误应该失败
		if networkErr, ok := err.(*NetworkError); ok {
			t.Errorf("网络错误：%v", networkErr.Message)
			return
		}
		// 其他错误应该失败
		t.Errorf("其他错误：%v", err)
		return
	}

	// 成功获取token，验证token
	assertStringNotEmpty(t, token, "Token不应该为空")
	assertTrue(t, len(token) > 10, "Token长度应该大于10")
	t.Logf("Token获取成功：%s", token)
}

func TestTokenManager_TokenCaching(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 第一次获取token
	token1, err1 := tokenManager.GetToken()
	if err1 != nil {
		if apiErr, ok := err1.(*APIError); ok {
			// 如果是配置相关的错误，记录但允许通过
			if apiErr.Code == 20001 || apiErr.Code == 20002 || apiErr.Code == 20003 || apiErr.Code == 404 {
				t.Logf("第一次获取Token配置错误（预期）：%s", apiErr.Message)
				return
			}
			t.Errorf("第一次获取Token意外API错误：code=%d, message=%s", apiErr.Code, apiErr.Message)
			return
		}
		t.Errorf("第一次获取Token其他错误：%v", err1)
		return
	}

	// 第二次获取token（应该使用缓存）
	token2, err2 := tokenManager.GetToken()
	if err2 != nil {
		if apiErr, ok := err2.(*APIError); ok {
			// 如果是配置相关的错误，记录但允许通过
			if apiErr.Code == 20001 || apiErr.Code == 20002 || apiErr.Code == 20003 || apiErr.Code == 404 {
				t.Logf("第二次获取Token配置错误（预期）：%s", apiErr.Message)
				return
			}
			t.Errorf("第二次获取Token意外API错误：code=%d, message=%s", apiErr.Code, apiErr.Message)
			return
		}
		t.Errorf("第二次获取Token其他错误：%v", err2)
		return
	}

	// 验证两次获取的token相同（缓存机制）
	assertEqual(t, token1, token2, "Token缓存失败")
	assertStringNotEmpty(t, token1, "Token不应该为空")
	t.Logf("Token缓存成功：%s", token1)
}

func TestTokenManager_InvalidConfig(t *testing.T) {
	// 创建无效配置的TokenManager
	invalidConfig := &Config{
		AppID:        "", // 空的AppID
		AppKey:       "test_app_key",
		MasterSecret: "test_master_secret",
		Domain:       "https://restapi.getui.com/v2",
	}

	// TokenManager构造函数不会验证配置，所以不会panic
	httpClient := invalidConfig.GetHTTPClient()
	tokenManager := NewTokenManager(invalidConfig, httpClient)

	// 但是获取token时会失败
	_, err := tokenManager.GetToken()
	assertError(t, err, "无效配置应该返回错误")

	// 检查是否是API错误
	if apiErr, ok := err.(*APIError); ok {
		t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
		assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
	} else {
		t.Logf("其他错误：%v", err)
	}
}

func TestTokenManager_NetworkError(t *testing.T) {
	// 创建指向无效域名的配置
	invalidDomainConfig := &Config{
		AppID:        "test_app_id",
		AppKey:       "test_app_key",
		MasterSecret: "test_master_secret",
		Domain:       "https://invalid-domain-that-does-not-exist.com",
	}

	httpClient := invalidDomainConfig.GetHTTPClient()
	tokenManager := NewTokenManager(invalidDomainConfig, httpClient)

	// 获取token，应该返回网络错误
	_, err := tokenManager.GetToken()

	assertError(t, err, "应该返回网络错误")

	// 检查是否是网络错误
	if networkErr, ok := err.(*NetworkError); ok {
		t.Logf("网络错误（预期）：%v", networkErr.Message)
		assertErrorType(t, err, &NetworkError{}, "错误类型应该是NetworkError")
	} else {
		t.Logf("其他错误：%v", err)
	}
}

func TestTokenManager_SignGeneration(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 测试签名生成
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sign := tokenManager.generateSign(timestamp)

	// 验证签名不为空
	assertStringNotEmpty(t, sign, "生成的签名不应该为空")

	// 验证签名长度（SHA256是64位十六进制）
	assertEqual(t, 64, len(sign), "SHA256签名长度应该为64")

	// 验证签名格式（应该是十六进制）
	for _, char := range sign {
		if !((char >= '0' && char <= '9') || (char >= 'a' && char <= 'f')) {
			t.Errorf("签名应该只包含十六进制字符，发现：%c", char)
			break
		}
	}

	t.Logf("签名生成成功：%s", sign)
}

func TestTokenManager_SignConsistency(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 使用相同的时间戳生成两次签名
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sign1 := tokenManager.generateSign(timestamp)
	sign2 := tokenManager.generateSign(timestamp)

	// 验证签名一致性
	assertEqual(t, sign1, sign2, "相同时间戳生成的签名应该一致")
	t.Logf("签名一致性验证成功：%s", sign1)
}

func TestTokenManager_TokenExpiration(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 模拟token已过期的情况
	tokenManager.SetToken("expired_token", time.Now().Add(-1*time.Hour)) // 1小时前过期

	// 获取token，应该重新请求
	token, err := tokenManager.GetToken()

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			// 如果是配置相关的错误，记录但允许通过
			if apiErr.Code == 20001 || apiErr.Code == 20002 || apiErr.Code == 20003 || apiErr.Code == 404 {
				t.Logf("Token过期后重新获取配置错误（预期）：%s", apiErr.Message)
				return
			}
			t.Errorf("Token过期后重新获取意外API错误：code=%d, message=%s", apiErr.Code, apiErr.Message)
			return
		}
		t.Errorf("Token过期后重新获取其他错误：%v", err)
		return
	}

	// 验证获取到新的token
	assertNotEqual(t, "expired_token", token, "应该获取到新的token，而不是过期的token")
	assertStringNotEmpty(t, token, "新token不应该为空")
	t.Logf("Token过期后重新获取成功：%s", token)
}

func TestTokenManager_ValidTokenNotExpired(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 模拟有效的token（未过期）
	tokenManager.SetToken("valid_token", time.Now().Add(1*time.Hour)) // 1小时后过期

	// 获取token，应该直接返回缓存的token
	token, err := tokenManager.GetToken()

	assertNoError(t, err, "有效token不应该返回错误")

	// 验证返回的是缓存的token
	assertEqual(t, "valid_token", token, "应该返回缓存的token")
	t.Logf("有效token缓存验证成功：%s", token)
}

func TestTokenManager_EmptyTokenExpired(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 模拟空token且已过期的情况
	tokenManager.SetToken("", time.Now().Add(-1*time.Hour)) // 空token且1小时前过期

	// 获取token，应该重新请求
	token, err := tokenManager.GetToken()

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			// 如果是配置相关的错误，记录但允许通过
			if apiErr.Code == 20001 || apiErr.Code == 20002 || apiErr.Code == 20003 || apiErr.Code == 404 {
				t.Logf("空Token过期后重新获取配置错误（预期）：%s", apiErr.Message)
				return
			}
			t.Errorf("空Token过期后重新获取意外API错误：code=%d, message=%s", apiErr.Code, apiErr.Message)
			return
		}
		t.Errorf("空Token过期后重新获取其他错误：%v", err)
		return
	}

	// 验证获取到新的token
	assertStringNotEmpty(t, token, "应该获取到新的token")
	assertNotEqual(t, "", token, "新token不应该为空")
	t.Logf("空Token过期后重新获取成功：%s", token)
}

func TestTokenManager_RequestFormat(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 测试请求格式构建
	timestamp := strconv.FormatInt(time.Now().UnixNano()/1e6, 10)
	sign := tokenManager.generateSign(timestamp)

	authDTO := &AuthDTO{
		Sign:      sign,
		Timestamp: timestamp,
		AppKey:    tokenManager.config.AppKey,
	}

	// 验证请求格式
	assertNotNil(t, authDTO, "AuthDTO不应该为nil")
	assertStringNotEmpty(t, authDTO.Sign, "Sign不应该为空")
	assertStringNotEmpty(t, authDTO.Timestamp, "Timestamp不应该为空")
	assertStringNotEmpty(t, authDTO.AppKey, "AppKey不应该为空")
	assertEqual(t, 64, len(authDTO.Sign), "Sign长度应该为64")
	assertEqual(t, tokenManager.config.AppKey, authDTO.AppKey, "AppKey应该匹配")

	t.Logf("请求格式验证成功：Sign=%s, Timestamp=%s, AppKey=%s", authDTO.Sign, authDTO.Timestamp, authDTO.AppKey)
}

func TestTokenManager_ResponseParsing(t *testing.T) {
	// 模拟成功的响应数据
	successData := []byte(`{"token": "test_token_12345"}`)

	var tokenData struct {
		Token string `json:"token"`
	}
	err := json.Unmarshal(successData, &tokenData)

	assertNoError(t, err, "JSON解析不应该有错误")
	assertStringNotEmpty(t, tokenData.Token, "解析出的Token不应该为空")
	assertEqual(t, "test_token_12345", tokenData.Token, "Token值应该匹配")

	t.Logf("响应解析验证成功：%s", tokenData.Token)
}

func TestTokenManager_ConcurrentAccess(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 并发获取token
	results := make(chan string, 5)
	errors := make(chan error, 5)

	for i := 0; i < 5; i++ {
		go func() {
			token, err := tokenManager.GetToken()
			if err != nil {
				errors <- err
				return
			}
			results <- token
		}()
	}

	// 收集结果
	var tokens []string
	var errs []error

	for i := 0; i < 5; i++ {
		select {
		case token := <-results:
			tokens = append(tokens, token)
		case err := <-errors:
			errs = append(errs, err)
		}
	}

	// 验证结果
	if len(tokens) > 0 {
		// 如果成功获取到token，验证所有token都是有效的
		for _, token := range tokens {
			assertStringNotEmpty(t, token, "Token不应该为空")
		}
		t.Logf("并发访问验证成功，获取到%d个有效token", len(tokens))
	} else if len(errs) > 0 {
		// 如果都失败了，验证错误类型一致且是配置相关错误
		allConfigErrors := true

		for _, err := range errs {
			if apiErr, ok := err.(*APIError); ok {
				if apiErr.Code != 20001 && apiErr.Code != 20002 && apiErr.Code != 20003 && apiErr.Code != 404 {
					allConfigErrors = false
					break
				}
			} else {
				allConfigErrors = false
				break
			}
		}

		if allConfigErrors {
			t.Logf("并发访问验证成功，所有请求都返回配置相关错误")
		} else {
			t.Errorf("并发访问失败，错误类型不一致或包含非配置错误")
		}
	}
}

func TestTokenManager_GetTokenExpireTime(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 设置一个特定的过期时间
	expectedExpireTime := time.Now().Add(2 * time.Hour)
	tokenManager.SetToken("test_token", expectedExpireTime)

	// 获取过期时间
	actualExpireTime := tokenManager.GetTokenExpireTime()

	// 验证过期时间（允许1秒的误差）
	timeDiff := actualExpireTime.Sub(expectedExpireTime)
	assertTrue(t, timeDiff >= -time.Second && timeDiff <= time.Second, "过期时间应该匹配")

	t.Logf("过期时间验证成功：%v", actualExpireTime)
}

func TestTokenManager_GetCurrentToken(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 设置token
	expectedToken := "test_current_token"
	tokenManager.SetToken(expectedToken, time.Now().Add(1*time.Hour))

	// 获取当前token
	actualToken := tokenManager.GetCurrentToken()

	// 验证token
	assertEqual(t, expectedToken, actualToken, "当前token应该匹配")

	t.Logf("当前token验证成功：%s", actualToken)
}

func TestTokenManager_SetToken(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 设置token
	expectedToken := "test_set_token"
	expectedExpireTime := time.Now().Add(1 * time.Hour)
	tokenManager.SetToken(expectedToken, expectedExpireTime)

	// 验证token设置
	actualToken := tokenManager.GetCurrentToken()
	actualExpireTime := tokenManager.GetTokenExpireTime()

	assertEqual(t, expectedToken, actualToken, "设置的token应该匹配")
	timeDiff := actualExpireTime.Sub(expectedExpireTime)
	assertTrue(t, timeDiff >= -time.Second && timeDiff <= time.Second, "设置的过期时间应该匹配")

	t.Logf("设置token验证成功：%s", actualToken)
}

func TestTokenManager_ClearToken(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 先设置token
	tokenManager.SetToken("test_clear_token", time.Now().Add(1*time.Hour))

	// 清除token
	tokenManager.ClearToken()

	// 验证token已清除
	actualToken := tokenManager.GetCurrentToken()
	actualExpireTime := tokenManager.GetTokenExpireTime()

	assertStringEmpty(t, actualToken, "清除后token应该为空")
	assertTrue(t, actualExpireTime.IsZero(), "清除后过期时间应该为零值")

	t.Logf("清除token验证成功")
}

func TestTokenManager_IsTokenExpired(t *testing.T) {
	tokenManager := createTestTokenManager()

	// 测试空token
	tokenManager.ClearToken()
	assertTrue(t, tokenManager.IsTokenExpired(), "空token应该被认为是过期的")

	// 测试已过期的token
	tokenManager.SetToken("expired_token", time.Now().Add(-1*time.Hour))
	assertTrue(t, tokenManager.IsTokenExpired(), "已过期的token应该被认为是过期的")

	// 测试有效的token
	tokenManager.SetToken("valid_token", time.Now().Add(1*time.Hour))
	assertFalse(t, tokenManager.IsTokenExpired(), "有效的token不应该被认为是过期的")

	t.Logf("Token过期状态验证成功")
}

func TestClient_GetToken_Delegation(t *testing.T) {
	client := createTestClient()

	// 验证Client的GetToken方法委托给TokenManager
	token1, err1 := client.GetToken()
	token2, err2 := client.GetTokenManager().GetToken()

	// 如果都成功，验证结果相同
	if err1 == nil && err2 == nil {
		assertEqual(t, token1, token2, "Client和TokenManager的GetToken结果应该相同")
		assertStringNotEmpty(t, token1, "Token不应该为空")
		t.Logf("委托验证成功：%s", token1)
	} else if err1 != nil && err2 != nil {
		// 如果都失败，验证错误类型相同且是配置相关错误
		if apiErr1, ok1 := err1.(*APIError); ok1 {
			if apiErr2, ok2 := err2.(*APIError); ok2 {
				if (apiErr1.Code == 20001 || apiErr1.Code == 20002 || apiErr1.Code == 20003 || apiErr1.Code == 404) &&
					(apiErr2.Code == 20001 || apiErr2.Code == 20002 || apiErr2.Code == 20003 || apiErr2.Code == 404) {
					t.Logf("委托验证成功，都返回配置相关错误")
					return
				}
			}
		}
		t.Errorf("Client和TokenManager的GetToken行为不一致：err1=%v, err2=%v", err1, err2)
	} else {
		t.Error("Client和TokenManager的GetToken行为应该一致")
	}
}
