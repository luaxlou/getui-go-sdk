package main

import (
	"fmt"
	"log"

	"github.com/luaxlou/getui-go-sdk"
)

func main() {
	// 创建配置
	config := &getui.Config{
		AppID:        "your_app_id",
		AppKey:       "your_app_key",
		MasterSecret: "your_master_secret",
		Domain:       "https://restapi.getui.com/v2",
	}

	// 创建客户端
	client := getui.NewClient(config)

	// 示例1: 单推
	singlePushExample(client)

	// 示例2: 批量推送
	batchPushExample(client)

	// 示例3: 群推
	allPushExample(client)

	// 示例4: 查询用户状态
	queryUserStatusExample(client)
}

// 单推示例
func singlePushExample(client *getui.Client) {
	fmt.Println("=== 单推示例 ===")

	// 创建推送消息
	pushMessage := &getui.PushMessage{
		Notification: &getui.Notification{
			Title:     "推送标题",
			Body:      "推送内容",
			ClickType: "url",
			URL:       "https://www.getui.com",
		},
	}

	// 创建推送请求
	pushDTO := &getui.PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		Audience: &getui.Audience{
			CIDs: []string{"target_cid_here"},
		},
	}

	// 执行单推
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)
	if err != nil {
		log.Printf("推送失败: %v", err)
		return
	}

	if result.IsSuccess() {
		log.Printf("推送成功: %+v", result.Data)
	} else {
		log.Printf("推送失败: code=%d, msg=%s", result.Code, result.Msg)
	}
}

// 批量推送示例
func batchPushExample(client *getui.Client) {
	fmt.Println("=== 批量推送示例 ===")

	// 创建推送消息
	pushMessage := &getui.PushMessage{
		Notification: &getui.Notification{
			Title:     "批量推送标题",
			Body:      "批量推送内容",
			ClickType: "intent",
			Intent:    "intent://com.example.app/main#Intent;scheme=example;launchFlags=0x10000000;end",
		},
	}

	// 创建批量推送请求
	batchDTO := &getui.PushBatchDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		Audience: &getui.Audience{
			CIDs: []string{"cid1", "cid2", "cid3"},
		},
	}

	// 执行批量推送
	result, err := client.PushAPI.PushBatchByCID(batchDTO)
	if err != nil {
		log.Printf("批量推送失败: %v", err)
		return
	}

	if result.IsSuccess() {
		log.Printf("批量推送成功: %+v", result.Data)
	} else {
		log.Printf("批量推送失败: code=%d, msg=%s", result.Code, result.Msg)
	}
}

// 群推示例
func allPushExample(client *getui.Client) {
	fmt.Println("=== 群推示例 ===")

	// 创建推送消息
	pushMessage := &getui.PushMessage{
		Notification: &getui.Notification{
			Title:     "群推标题",
			Body:      "群推内容",
			ClickType: "payload",
			Payload:   `{"key": "value", "action": "open_app"}`,
		},
	}

	// 创建群推请求
	pushDTO := &getui.PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		Audience:    "all", // 推送给所有用户
	}

	// 执行群推
	result, err := client.PushAPI.PushAll(pushDTO)
	if err != nil {
		log.Printf("群推失败: %v", err)
		return
	}

	if result.IsSuccess() {
		var taskID getui.TaskIDDTO
		if err := result.UnmarshalData(&taskID); err != nil {
			log.Printf("解析任务ID失败: %v", err)
			return
		}
		log.Printf("群推成功，任务ID: %s", taskID.TaskID)
	} else {
		log.Printf("群推失败: code=%d, msg=%s", result.Code, result.Msg)
	}
}

// 查询用户状态示例
func queryUserStatusExample(client *getui.Client) {
	fmt.Println("=== 查询用户状态示例 ===")

	cids := []string{"cid1", "cid2", "cid3"}

	// 查询用户状态
	result, err := client.UserAPI.QueryUserStatus(cids)
	if err != nil {
		log.Printf("查询用户状态失败: %v", err)
		return
	}

	if result.IsSuccess() {
		var userStatus map[string]getui.CidStatusDTO
		if err := result.UnmarshalData(&userStatus); err != nil {
			log.Printf("解析用户状态失败: %v", err)
			return
		}

		for cid, status := range userStatus {
			log.Printf("CID: %s, 状态: %s", cid, status.Status)
		}
	} else {
		log.Printf("查询用户状态失败: code=%d, msg=%s", result.Code, result.Msg)
	}
}

// 透传消息示例
func transmissionMessageExample(client *getui.Client) {
	fmt.Println("=== 透传消息示例 ===")

	// 创建透传消息
	pushMessage := &getui.PushMessage{
		Transmission: `{"type": "custom", "data": {"key": "value"}}`,
	}

	// 创建推送请求
	pushDTO := &getui.PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		Audience: &getui.Audience{
			CIDs: []string{"target_cid_here"},
		},
	}

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)
	if err != nil {
		log.Printf("透传消息推送失败: %v", err)
		return
	}

	if result.IsSuccess() {
		log.Printf("透传消息推送成功")
	} else {
		log.Printf("透传消息推送失败: code=%d, msg=%s", result.Code, result.Msg)
	}
}

// 厂商通道推送示例
func vendorChannelExample(client *getui.Client) {
	fmt.Println("=== 厂商通道推送示例 ===")

	// 创建推送消息
	pushMessage := &getui.PushMessage{
		Notification: &getui.Notification{
			Title:     "个推通道标题",
			Body:      "个推通道内容",
			ClickType: "url",
			URL:       "https://www.getui.com",
		},
	}

	// 创建推送通道配置
	pushChannel := &getui.PushChannel{
		Android: &getui.AndroidDTO{
			UPS: &getui.UPS{
				Notification: &getui.ThirdNotification{
					Title:     "厂商通道标题",
					Body:      "厂商通道内容",
					ClickType: "url",
					URL:       "https://www.getui.com",
				},
			},
		},
		IOS: &getui.IOSDTO{
			APNS: &getui.APNS{
				Alert: &getui.Alert{
					Title: "iOS通知标题",
					Body:  "iOS通知内容",
				},
				Badge: 1,
				Sound: "default",
			},
		},
	}

	// 创建推送请求
	pushDTO := &getui.PushDTO{
		RequestID:   client.GenerateRequestID(),
		PushMessage: pushMessage,
		PushChannel: pushChannel,
		Audience: &getui.Audience{
			CIDs: []string{"target_cid_here"},
		},
	}

	// 执行推送
	result, err := client.PushAPI.PushToSingleByCID(pushDTO)
	if err != nil {
		log.Printf("厂商通道推送失败: %v", err)
		return
	}

	if result.IsSuccess() {
		log.Printf("厂商通道推送成功")
	} else {
		log.Printf("厂商通道推送失败: code=%d, msg=%s", result.Code, result.Msg)
	}
}
