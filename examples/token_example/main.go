package main

import (
	"fmt"
	"time"

	"github.com/luaxlou/getui-go-sdk"
)

func main() {
	fmt.Println("🔐 个推Go SDK GetToken示例")
	fmt.Println("============================")

	// 从.env文件加载配置
	config := getui.LoadConfigFromEnvFileOrDefault(".env")

	fmt.Printf("📱 AppID: %s\n", config.AppID)
	fmt.Printf("🔑 AppKey: %s\n", config.AppKey)
	fmt.Printf("🔐 MasterSecret: %s\n", maskSecret(config.MasterSecret))
	fmt.Printf("🌐 Domain: %s\n", config.Domain)
	fmt.Println()

	// 创建客户端
	client := getui.NewClient(config)

	// 测试时间戳生成
	timestamp := fmt.Sprintf("%d", time.Now().UnixNano()/1e6)
	fmt.Printf("📝 时间戳: %s\n", timestamp)
	fmt.Println()

	// 获取Token
	fmt.Println("🔄 正在获取Token...")
	start := time.Now()

	token, err := client.GetToken()
	duration := time.Since(start)

	if err != nil {
		if apiErr, ok := err.(*getui.APIError); ok {
			fmt.Printf("❌ API错误: code=%d, message=%s\n", apiErr.Code, apiErr.Message)
		} else if networkErr, ok := err.(*getui.NetworkError); ok {
			fmt.Printf("❌ 网络错误: %s\n", networkErr.Message)
		} else {
			fmt.Printf("❌ 其他错误: %v\n", err)
		}
	} else {
		fmt.Printf("✅ Token获取成功: %s\n", maskToken(token))
	}

	fmt.Printf("⏱️  耗时: %v\n", duration)
	fmt.Println()

	// 测试Token缓存
	fmt.Println("🔄 再次获取Token（测试缓存）...")
	start = time.Now()

	token2, err2 := client.GetToken()
	duration2 := time.Since(start)

	if err2 != nil {
		if apiErr, ok := err2.(*getui.APIError); ok {
			fmt.Printf("❌ 第二次获取API错误: code=%d, message=%s\n", apiErr.Code, apiErr.Message)
		} else {
			fmt.Printf("❌ 第二次获取其他错误: %v\n", err2)
		}
	} else {
		fmt.Printf("✅ 第二次Token获取成功: %s\n", maskToken(token2))

		// 验证缓存
		if token == token2 {
			fmt.Println("✅ Token缓存验证成功")
		} else {
			fmt.Println("❌ Token缓存验证失败")
		}
	}

	fmt.Printf("⏱️  第二次耗时: %v\n", duration2)
	fmt.Println()

	fmt.Println("✅ Token示例演示完成")
}

// maskSecret 隐藏敏感信息
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "****" + secret[len(secret)-4:]
}

// maskToken 隐藏Token信息
func maskToken(token string) string {
	if len(token) <= 10 {
		return "***"
	}
	return token[:6] + "****" + token[len(token)-6:]
}
