package getui

import (
	"os"
	"testing"
)

func TestLoadConfigFromEnvFile(t *testing.T) {
	// 创建临时.env文件
	content := `# 个推Go SDK 测试环境变量配置
# 这是一个测试配置文件

# 个推应用配置
GETUI_TEST_APP_ID=test_app_id_from_env_file
GETUI_TEST_APP_KEY=test_app_key_from_env_file
GETUI_TEST_MASTER_SECRET=test_master_secret_from_env_file
GETUI_TEST_DOMAIN=https://test.getui.com/v2
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_env_*.env")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试内容
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("写入测试内容失败: %v", err)
	}
	tmpFile.Close()

	// 测试加载配置
	config, err := LoadConfigFromEnvFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("加载.env文件失败: %v", err)
	}

	// 验证配置值
	expectedValues := map[string]string{
		"AppID":        "test_app_id_from_env_file",
		"AppKey":       "test_app_key_from_env_file",
		"MasterSecret": "test_master_secret_from_env_file",
		"Domain":       "https://test.getui.com/v2",
	}

	if config.AppID != expectedValues["AppID"] {
		t.Errorf("AppID不匹配，期望: %s, 实际: %s", expectedValues["AppID"], config.AppID)
	}

	if config.AppKey != expectedValues["AppKey"] {
		t.Errorf("AppKey不匹配，期望: %s, 实际: %s", expectedValues["AppKey"], config.AppKey)
	}

	if config.MasterSecret != expectedValues["MasterSecret"] {
		t.Errorf("MasterSecret不匹配，期望: %s, 实际: %s", expectedValues["MasterSecret"], config.MasterSecret)
	}

	if config.Domain != expectedValues["Domain"] {
		t.Errorf("Domain不匹配，期望: %s, 实际: %s", expectedValues["Domain"], config.Domain)
	}
}

func TestLoadConfigFromEnvFile_FileNotExists(t *testing.T) {
	// 测试文件不存在的情况
	_, err := LoadConfigFromEnvFile("non_existent_file.env")
	if err == nil {
		t.Error("期望返回错误，但实际没有错误")
	}
}

func TestLoadConfigFromEnvFileOrDefault_FileNotExists(t *testing.T) {
	// 测试文件不存在时返回默认配置
	config := LoadConfigFromEnvFileOrDefault("non_existent_file.env")

	// 验证返回的是默认配置
	defaultConfig := NewDefaultConfig()
	if config.AppID != defaultConfig.AppID {
		t.Errorf("AppID不匹配，期望: %s, 实际: %s", defaultConfig.AppID, config.AppID)
	}
}

func TestLoadConfigFromEnvFile_WithQuotes(t *testing.T) {
	// 测试带引号的值
	content := `# 个推Go SDK 测试环境变量配置

GETUI_TEST_APP_ID="test_app_id_with_quotes"
GETUI_TEST_APP_KEY="test_app_key_with_quotes"
GETUI_TEST_MASTER_SECRET="test_master_secret_with_quotes"
GETUI_TEST_DOMAIN="https://test.getui.com/v2"
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_env_quotes_*.env")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试内容
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("写入测试内容失败: %v", err)
	}
	tmpFile.Close()

	// 测试加载配置
	config, err := LoadConfigFromEnvFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("加载.env文件失败: %v", err)
	}

	// 验证引号被正确移除
	expectedValues := map[string]string{
		"AppID":        "test_app_id_with_quotes",
		"AppKey":       "test_app_key_with_quotes",
		"MasterSecret": "test_master_secret_with_quotes",
		"Domain":       "https://test.getui.com/v2",
	}

	if config.AppID != expectedValues["AppID"] {
		t.Errorf("AppID不匹配，期望: %s, 实际: %s", expectedValues["AppID"], config.AppID)
	}

	if config.AppKey != expectedValues["AppKey"] {
		t.Errorf("AppKey不匹配，期望: %s, 实际: %s", expectedValues["AppKey"], config.AppKey)
	}

	if config.MasterSecret != expectedValues["MasterSecret"] {
		t.Errorf("MasterSecret不匹配，期望: %s, 实际: %s", expectedValues["MasterSecret"], config.MasterSecret)
	}

	if config.Domain != expectedValues["Domain"] {
		t.Errorf("Domain不匹配，期望: %s, 实际: %s", expectedValues["Domain"], config.Domain)
	}
}

func TestLoadConfigFromEnvFile_IgnoreCommentsAndEmptyLines(t *testing.T) {
	// 测试忽略注释和空行
	content := `# 这是注释行

# 另一个注释行
GETUI_TEST_APP_ID=test_app_id

# 空行后跟着配置
GETUI_TEST_APP_KEY=test_app_key

GETUI_TEST_MASTER_SECRET=test_master_secret
# 行尾注释 GETUI_TEST_DOMAIN=https://test.getui.com/v2
GETUI_TEST_DOMAIN=https://test.getui.com/v2
`

	// 创建临时文件
	tmpFile, err := os.CreateTemp("", "test_env_comments_*.env")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// 写入测试内容
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("写入测试内容失败: %v", err)
	}
	tmpFile.Close()

	// 测试加载配置
	config, err := LoadConfigFromEnvFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("加载.env文件失败: %v", err)
	}

	// 验证只读取了有效的配置行
	expectedValues := map[string]string{
		"AppID":        "test_app_id",
		"AppKey":       "test_app_key",
		"MasterSecret": "test_master_secret",
		"Domain":       "https://test.getui.com/v2",
	}

	if config.AppID != expectedValues["AppID"] {
		t.Errorf("AppID不匹配，期望: %s, 实际: %s", expectedValues["AppID"], config.AppID)
	}

	if config.AppKey != expectedValues["AppKey"] {
		t.Errorf("AppKey不匹配，期望: %s, 实际: %s", expectedValues["AppKey"], config.AppKey)
	}

	if config.MasterSecret != expectedValues["MasterSecret"] {
		t.Errorf("MasterSecret不匹配，期望: %s, 实际: %s", expectedValues["MasterSecret"], config.MasterSecret)
	}

	if config.Domain != expectedValues["Domain"] {
		t.Errorf("Domain不匹配，期望: %s, 实际: %s", expectedValues["Domain"], config.Domain)
	}
}
