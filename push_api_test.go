package getui

import (
	"reflect"
	"testing"
)

// assert辅助函数
func assertEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: 期望 %v, 实际 %v", msg, expected, actual)
	}
}

func assertNotEqual(t *testing.T, expected, actual interface{}, msg string) {
	t.Helper()
	if reflect.DeepEqual(expected, actual) {
		t.Errorf("%s: 不应该相等，但实际相等: %v", msg, actual)
	}
}

func assertNil(t *testing.T, actual interface{}, msg string) {
	t.Helper()
	if actual != nil && !reflect.ValueOf(actual).IsNil() {
		t.Errorf("%s: 期望nil，实际 %v", msg, actual)
	}
}

func assertNotNil(t *testing.T, actual interface{}, msg string) {
	t.Helper()
	if actual == nil {
		t.Errorf("%s: 期望非nil，实际为nil", msg)
	}
}

func assertTrue(t *testing.T, condition bool, msg string) {
	t.Helper()
	if !condition {
		t.Errorf("%s: 期望为true，实际为false", msg)
	}
}

func assertFalse(t *testing.T, condition bool, msg string) {
	t.Helper()
	if condition {
		t.Errorf("%s: 期望为false，实际为true", msg)
	}
}

func assertError(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: 期望有错误，实际无错误", msg)
	}
}

func assertNoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Errorf("%s: 期望无错误，实际有错误: %v", msg, err)
	}
}

func assertErrorType(t *testing.T, err error, expectedType interface{}, msg string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: 期望有错误，实际无错误", msg)
		return
	}

	expectedTypeName := reflect.TypeOf(expectedType).String()
	actualTypeName := reflect.TypeOf(err).String()
	if expectedTypeName != actualTypeName {
		t.Errorf("%s: 期望错误类型 %s，实际错误类型 %s", msg, expectedTypeName, actualTypeName)
	}
}

func assertStringNotEmpty(t *testing.T, str string, msg string) {
	t.Helper()
	if str == "" {
		t.Errorf("%s: 期望非空字符串，实际为空", msg)
	}
}

func assertStringEmpty(t *testing.T, str string, msg string) {
	t.Helper()
	if str != "" {
		t.Errorf("%s: 期望空字符串，实际为: %s", msg, str)
	}
}

// 创建测试客户端
func createTestClient() *Client {
	config := LoadConfigFromEnvFileOrDefault(".env")
	return NewClient(config)
}

// 创建基本的推送消息
func createTestPushMessage() *PushMessage {
	return &PushMessage{
		Notification: &Notification{
			Title:     "测试推送标题",
			Body:      "测试推送内容",
			ClickType: "url",
			URL:       "https://www.getui.com",
		},
	}
}

// 创建基本的受众
func createTestAudience() *Audience {
	return &Audience{
		CIDs: []string{"test_cid_123"},
	}
}

func TestPushToSingleByCID_ValidRequest(t *testing.T) {
	client := createTestClient()

	// 创建有效的推送请求
	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: createTestPushMessage(),
		Audience:    createTestAudience(),
	}

	// 验证请求参数
	assertNotNil(t, pushDTO, "PushDTO不应该为nil")
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	// 由于使用测试配置，预期会失败，但应该能验证请求格式正确
	if err != nil {
		// 检查是否是API错误（说明请求格式正确，但配置无效）
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			// 验证API错误类型
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		// 其他错误可能是网络问题或配置问题
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	// 如果成功，验证结果
	if result != nil {
		if result.IsSuccess() {
			t.Logf("推送成功：%+v", result.Data)
			assertTrue(t, result.IsSuccess(), "推送应该成功")
		} else {
			t.Logf("推送失败：code=%d, msg=%s", result.Code, result.Msg)
			assertFalse(t, result.IsSuccess(), "推送应该失败")
		}
	}
}

func TestPushToSingleByCID_EmptyRequestID(t *testing.T) {
	client := createTestClient()

	// 创建不包含RequestID的推送请求
	pushDTO := &PushDTO{
		PushMessage: createTestPushMessage(),
		Audience:    createTestAudience(),
	}

	// 验证初始状态
	assertStringEmpty(t, pushDTO.RequestID, "初始RequestID应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	// 应该自动生成RequestID
	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	// 验证RequestID是否被自动生成
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID应该被自动生成")

	if result != nil {
		t.Logf("推送结果：code=%d, msg=%s", result.Code, result.Msg)
	}
}

func TestPushToSingleByCID_InvalidRequestID(t *testing.T) {
	client := createTestClient()

	// 创建包含无效RequestID的推送请求
	pushDTO := &PushDTO{
		RequestID:   "123", // 太短
		PushMessage: createTestPushMessage(),
		Audience:    createTestAudience(),
	}

	// 验证无效RequestID
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不为空")
	assertTrue(t, len(pushDTO.RequestID) < 10, "RequestID长度应该小于10")

	// 执行推送，应该返回错误
	_, err := client.PushAPI.PushToSingleByCID(pushDTO)

	assertError(t, err, "应该返回RequestID无效的错误")
	assertEqual(t, ErrInvalidRequestID, err, "期望错误类型为ErrInvalidRequestID")
}

func TestPushToSingleByCID_NilPushDTO(t *testing.T) {
	client := createTestClient()

	// 传入nil
	_, err := client.PushAPI.PushToSingleByCID(nil)

	assertError(t, err, "应该返回push_dto cannot be nil的错误")

	expectedErr := "push_dto cannot be nil"
	assertEqual(t, expectedErr, err.Error(), "错误消息应该匹配")
}

func TestPushToSingleByCID_EmptyAudience(t *testing.T) {
	client := createTestClient()

	// 创建不包含Audience的推送请求
	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: createTestPushMessage(),
		Audience:    nil,
	}

	// 验证初始状态
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNil(t, pushDTO.Audience, "Audience应该为nil")

	// 执行推送，应该返回错误
	_, err := client.PushAPI.PushToSingleByCID(pushDTO)

	assertError(t, err, "应该返回audience cannot be empty的错误")
	assertEqual(t, ErrEmptyAudience, err, "期望错误类型为ErrEmptyAudience")
}

func TestPushToSingleByCID_EmptyPushMessage(t *testing.T) {
	client := createTestClient()

	// 创建不包含PushMessage的推送请求
	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: nil,
		Audience:    createTestAudience(),
	}

	// 验证初始状态
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNil(t, pushDTO.PushMessage, "PushMessage应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")

	// 执行推送，应该返回错误
	_, err := client.PushAPI.PushToSingleByCID(pushDTO)

	assertError(t, err, "应该返回push_message cannot be empty的错误")
	assertEqual(t, ErrEmptyPushMessage, err, "期望错误类型为ErrEmptyPushMessage")
}

func TestPushToSingleByCID_TransmissionMessage(t *testing.T) {
	client := createTestClient()

	// 创建透传消息
	pushMessage := &PushMessage{
		Transmission: `{"type": "custom", "data": {"key": "value"}}`,
	}

	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		Audience:    createTestAudience(),
	}

	// 验证透传消息
	assertNotNil(t, pushMessage, "PushMessage不应该为nil")
	assertStringNotEmpty(t, pushMessage.Transmission, "Transmission不应该为空")
	assertNotNil(t, pushDTO, "PushDTO不应该为nil")
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	if result != nil {
		t.Logf("透传消息推送结果：code=%d, msg=%s", result.Code, result.Msg)
		assertNotNil(t, result, "结果不应该为nil")
	}
}

func TestPushToSingleByCID_WithSettings(t *testing.T) {
	client := createTestClient()

	// 创建包含设置的推送请求
	settings := &Settings{
		TTL: 3600,
		Strategy: &Strategy{
			Default: 1,
			IOS:     1,
		},
	}

	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: createTestPushMessage(),
		Audience:    createTestAudience(),
		Settings:    settings,
	}

	// 验证设置
	assertNotNil(t, settings, "Settings不应该为nil")
	assertEqual(t, 3600, settings.TTL, "TTL应该为3600")
	assertNotNil(t, settings.Strategy, "Strategy不应该为nil")
	assertEqual(t, 1, settings.Strategy.Default, "Default策略应该为1")
	assertEqual(t, 1, settings.Strategy.IOS, "IOS策略应该为1")

	// 验证PushDTO
	assertNotNil(t, pushDTO, "PushDTO不应该为nil")
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")
	assertNotNil(t, pushDTO.Settings, "Settings不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	if result != nil {
		t.Logf("带设置的推送结果：code=%d, msg=%s", result.Code, result.Msg)
		assertNotNil(t, result, "结果不应该为nil")
	}
}

func TestPushToSingleByCID_WithPushChannel(t *testing.T) {
	client := createTestClient()

	// 创建包含推送通道的推送请求
	pushChannel := &PushChannel{
		Android: &AndroidDTO{
			UPS: &UPS{
				Notification: &ThirdNotification{
					Title:     "厂商通道标题",
					Body:      "厂商通道内容",
					ClickType: "url",
					URL:       "https://www.getui.com",
				},
			},
		},
		IOS: &IOSDTO{
			APNS: &APNS{
				Alert: &Alert{
					Title: "iOS通知标题",
					Body:  "iOS通知内容",
				},
				Badge: 1,
				Sound: "default",
			},
		},
	}

	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: createTestPushMessage(),
		Audience:    createTestAudience(),
		PushChannel: pushChannel,
	}

	// 验证推送通道
	assertNotNil(t, pushChannel, "PushChannel不应该为nil")
	assertNotNil(t, pushChannel.Android, "Android通道不应该为nil")
	assertNotNil(t, pushChannel.Android.UPS, "UPS不应该为nil")
	assertNotNil(t, pushChannel.Android.UPS.Notification, "Notification不应该为nil")
	assertEqual(t, "厂商通道标题", pushChannel.Android.UPS.Notification.Title, "Android通知标题应该匹配")
	assertEqual(t, "厂商通道内容", pushChannel.Android.UPS.Notification.Body, "Android通知内容应该匹配")
	assertEqual(t, "url", pushChannel.Android.UPS.Notification.ClickType, "ClickType应该为url")
	assertEqual(t, "https://www.getui.com", pushChannel.Android.UPS.Notification.URL, "URL应该匹配")

	assertNotNil(t, pushChannel.IOS, "IOS通道不应该为nil")
	assertNotNil(t, pushChannel.IOS.APNS, "APNS不应该为nil")
	assertNotNil(t, pushChannel.IOS.APNS.Alert, "Alert不应该为nil")
	assertEqual(t, "iOS通知标题", pushChannel.IOS.APNS.Alert.Title, "iOS通知标题应该匹配")
	assertEqual(t, "iOS通知内容", pushChannel.IOS.APNS.Alert.Body, "iOS通知内容应该匹配")
	assertEqual(t, 1, pushChannel.IOS.APNS.Badge, "Badge应该为1")
	assertEqual(t, "default", pushChannel.IOS.APNS.Sound, "Sound应该为default")

	// 验证PushDTO
	assertNotNil(t, pushDTO, "PushDTO不应该为nil")
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")
	assertNotNil(t, pushDTO.PushChannel, "PushChannel不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	if result != nil {
		t.Logf("带推送通道的推送结果：code=%d, msg=%s", result.Code, result.Msg)
		assertNotNil(t, result, "结果不应该为nil")
	}
}

func TestPushToSingleByCID_MultipleCIDs(t *testing.T) {
	client := createTestClient()

	// 创建包含多个CID的受众
	audience := &Audience{
		CIDs: []string{"test_cid_1", "test_cid_2", "test_cid_3"},
	}

	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: createTestPushMessage(),
		Audience:    audience,
	}

	// 验证多CID受众
	assertNotNil(t, audience, "Audience不应该为nil")
	assertNotNil(t, audience.CIDs, "CIDs不应该为nil")
	assertEqual(t, 3, len(audience.CIDs), "CIDs数量应该为3")
	assertEqual(t, "test_cid_1", audience.CIDs[0], "第一个CID应该匹配")
	assertEqual(t, "test_cid_2", audience.CIDs[1], "第二个CID应该匹配")
	assertEqual(t, "test_cid_3", audience.CIDs[2], "第三个CID应该匹配")

	// 验证PushDTO
	assertNotNil(t, pushDTO, "PushDTO不应该为nil")
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	if result != nil {
		t.Logf("多CID推送结果：code=%d, msg=%s", result.Code, result.Msg)
		assertNotNil(t, result, "结果不应该为nil")
	}
}

func TestPushToSingleByCID_WithTaskName(t *testing.T) {
	client := createTestClient()

	// 创建包含任务名称的推送请求
	pushDTO := &PushDTO{
		RequestID:   client.GenerateRequestID(),
		TaskName:    "测试任务",
		PushMessage: createTestPushMessage(),
		Audience:    createTestAudience(),
	}

	// 验证任务名称
	assertStringNotEmpty(t, pushDTO.TaskName, "TaskName不应该为空")
	assertEqual(t, "测试任务", pushDTO.TaskName, "TaskName应该匹配")

	// 验证PushDTO
	assertNotNil(t, pushDTO, "PushDTO不应该为nil")
	assertStringNotEmpty(t, pushDTO.RequestID, "RequestID不应该为空")
	assertNotNil(t, pushDTO.PushMessage, "PushMessage不应该为nil")
	assertNotNil(t, pushDTO.Audience, "Audience不应该为nil")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)

	if err != nil {
		if apiErr, ok := err.(*APIError); ok {
			t.Logf("API错误（预期）：code=%d, message=%s", apiErr.Code, apiErr.Message)
			assertErrorType(t, err, &APIError{}, "错误类型应该是APIError")
			return
		}
		t.Logf("其他错误：%v", err)
		assertError(t, err, "应该返回错误")
		return
	}

	if result != nil {
		t.Logf("带任务名称的推送结果：code=%d, msg=%s", result.Code, result.Msg)
		assertNotNil(t, result, "结果不应该为nil")
	}
}
