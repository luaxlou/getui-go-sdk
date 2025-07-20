package main

import (
	"fmt"

	"github.com/luaxlou/getui-go-sdk"
)

func main() {
	fmt.Println("🚀 个推Go SDK .env文件配置示例")
	fmt.Println("==================================")

	// 从.env文件加载配置
	config := getui.LoadConfigFromEnvFileOrDefault(".env")

	fmt.Printf("📱 AppID: %s\n", config.AppID)
	fmt.Printf("🔑 AppKey: %s\n", config.AppKey)
	fmt.Printf("🔐 MasterSecret: %s\n", maskSecret(config.MasterSecret))
	fmt.Printf("🌐 Domain: %s\n", config.Domain)
	fmt.Println()

	// 创建客户端
	client := getui.NewClient(config)

	// 创建测试推送消息
	pushMessage := &getui.PushMessage{
		Notification: &getui.Notification{
			Title:     "测试推送",
			Body:      "这是一个测试推送消息",
			ClickType: "url",
			URL:       "https://www.getui.com",
		},
	}

	// 创建推送请求
	pushDTO := &getui.PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		Audience: &getui.Audience{
			CIDs: []string{"test_cid_123"},
		},
	}

	fmt.Println("📤 发送测试推送...")

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)
	if err != nil {
		if apiErr, ok := err.(*getui.APIError); ok {
			fmt.Printf("❌ API错误: code=%d, message=%s\n", apiErr.Code, apiErr.Message)
		} else {
			fmt.Printf("❌ 其他错误: %v\n", err)
		}
		return
	}

	if result.IsSuccess() {
		fmt.Printf("✅ 推送成功: %+v\n", result.Data)
	} else {
		fmt.Printf("❌ 推送失败: code=%d, msg=%s\n", result.Code, result.Msg)
	}
}

// maskSecret 隐藏敏感信息
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}
