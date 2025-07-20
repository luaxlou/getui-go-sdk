# GetToken示例

这个示例演示了如何使用个推Go SDK的GetToken功能。

## 功能说明

- 从.env文件加载配置
- 演示Token获取过程
- 测试Token缓存机制
- 显示错误处理
- 性能测试

## 使用方法

### 1. 准备.env文件

```bash
# 复制示例配置文件
cp ../../env.example .env

# 编辑.env文件，填入实际的个推配置
```

### 2. 运行示例

```bash
go run main.go
```

## 输出示例

```
🔐 个推Go SDK GetToken示例
============================
📱 AppID: your_app_id_here
🔑 AppKey: your_app_key_here
🔐 MasterSecret: your****cret
🌐 Domain: https://restapi.getui.com/v2

📝 时间戳: 1703123456789

🔄 正在获取Token...
❌ API错误: code=20001, message=sign is invalid
⏱️  耗时: 234.56ms

🔄 再次获取Token（测试缓存）...
❌ 第二次获取API错误: code=20001, message=sign is invalid
⏱️  第二次耗时: 0.12ms

✅ Token示例演示完成
```

## 特性演示

### 1. 配置加载
- 自动从.env文件加载配置
- 显示配置信息（隐藏敏感数据）

### 2. Token获取
- 演示完整的Token获取流程
- 显示请求耗时
- 错误处理和类型判断

### 3. 缓存机制
- 演示Token缓存功能
- 比较两次获取的耗时差异
- 验证缓存一致性

### 4. 错误处理
- API错误处理
- 网络错误处理
- 错误类型识别

## 技术要点

### Token获取流程
1. 检查缓存Token是否有效
2. 生成时间戳和签名
3. 发送认证请求
4. 解析响应获取Token
5. 缓存Token和过期时间

### 缓存机制
- Token有效期为24小时
- 提前1小时刷新
- 避免重复请求

### 签名算法
- MD5(appkey + timestamp + master_secret)
- 32位十六进制字符串
- 时间戳精确到毫秒

## 注意事项

- 需要有效的个推应用配置
- 网络连接正常
- Token有24小时有效期
- 缓存机制提高性能 