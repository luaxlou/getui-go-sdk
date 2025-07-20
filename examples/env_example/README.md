# .env文件配置示例

这个示例演示了如何使用.env文件来配置个推Go SDK。

## 功能说明

- 从.env文件自动加载配置
- 显示配置信息（隐藏敏感数据）
- 执行测试推送
- 处理API错误

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
🚀 个推Go SDK .env文件配置示例
==================================
📱 AppID: your_app_id_here
🔑 AppKey: your_app_key_here
🔐 MasterSecret: your****cret
🌐 Domain: https://restapi.getui.com/v2

📤 发送测试推送...
❌ API错误: code=20001, message=sign is invalid
```

## 配置说明

.env文件支持以下格式：

```bash
# 注释行
GETUI_TEST_APP_ID=your_app_id
GETUI_TEST_APP_KEY="your_app_key_with_quotes"
GETUI_TEST_MASTER_SECRET=your_master_secret
GETUI_TEST_DOMAIN=https://restapi.getui.com/v2
```

## 特性

- ✅ 自动加载.env文件
- ✅ 支持注释和空行
- ✅ 支持带引号的值
- ✅ 容错处理（文件不存在时使用默认配置）
- ✅ 敏感信息隐藏显示 