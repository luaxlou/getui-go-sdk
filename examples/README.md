# 个推Go SDK 示例

本目录包含了使用个推Go SDK的各种示例代码。

## 基本使用示例

### 运行示例

1. 首先确保你已经有个推开发者账号和应用配置
2. 修改 `basic_usage.go` 中的配置信息：
   ```go
   config := &getui.Config{
       AppID:        "your_app_id",        // 替换为你的应用ID
       AppKey:       "your_app_key",       // 替换为你的应用Key
       MasterSecret: "your_master_secret", // 替换为你的主密钥
       Domain:       "https://restapi.getui.com/v2",
   }
   ```
3. 运行示例：
   ```bash
   go run basic_usage.go
   ```

### 示例功能

`basic_usage.go` 包含了以下示例：

1. **单推示例** - 向单个用户推送消息
2. **批量推送示例** - 向多个用户批量推送消息
3. **群推示例** - 向所有用户推送消息
4. **查询用户状态示例** - 查询用户在线状态
5. **透传消息示例** - 发送透传消息
6. **厂商通道推送示例** - 使用厂商通道推送

### 注意事项

- 示例中的CID、别名等需要替换为实际的值
- 推送前请确保目标用户已经注册到个推平台
- 建议在测试环境中先验证功能
- 生产环境中请注意保护应用密钥的安全性

## 更多示例

更多使用示例请参考主README文件中的API文档。 