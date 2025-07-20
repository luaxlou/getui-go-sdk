# 个推签名问题解决方案

## 问题描述

在使用个推Go SDK时，Token获取一直返回 "sign is invalid" 错误，即使配置参数正确。

## 问题分析

通过详细的调试和测试，发现问题的根本原因是：

### 1. 签名算法错误
**错误实现**: `MD5(appkey + timestamp + master_secret)`
**正确实现**: `SHA256(appkey + timestamp + master_secret)`

### 2. 签名长度不匹配
- MD5签名: 32位十六进制字符串
- SHA256签名: 64位十六进制字符串

### 3. 时间戳格式正确
- 个推要求毫秒时间戳
- 当前实现使用 `time.Now().UnixNano()/1e6` 是正确的

## 解决方案

### 1. 更新签名算法

**修改前** (`token.go`):
```go
// generateSign 生成签名
func (tm *TokenManager) generateSign(timestamp string) string {
	// 签名算法：MD5(appkey + timestamp + master_secret)
	signStr := tm.config.AppKey + timestamp + tm.config.MasterSecret
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}
```

**修改后** (`token.go`):
```go
// generateSign 生成签名
func (tm *TokenManager) generateSign(timestamp string) string {
	// 签名算法：SHA256(appkey + timestamp + master_secret)
	signStr := tm.config.AppKey + timestamp + tm.config.MasterSecret
	return fmt.Sprintf("%x", sha256.Sum256([]byte(signStr)))
}
```

### 2. 更新导入包

**修改前**:
```go
import (
	"crypto/md5"
	// ...
)
```

**修改后**:
```go
import (
	"crypto/sha256"
	// ...
)
```

### 3. 更新测试用例

更新了相关的测试用例以适应SHA256签名：

- `TestTokenManager_SignGeneration`: 验证签名长度为64位
- `TestTokenManager_RequestFormat`: 验证请求格式中的签名长度
- `TestTokenManager_ConcurrentAccess`: 调整并发测试逻辑

## 验证结果

### 修复前
```
算法1: MD5(appkey + timestamp + master_secret)
   签名: 8867e9ec87c536a01d9591f4a3cd16dc
   结果: ❌ 失败: code=20001, message=sign is invalid
```

### 修复后
```
SHA256算法1: SHA256(appkey + timestamp + master_secret)
   签名: a2c51ae125ed812de09c9db84f106da67302978b4d0811ed2758ca2a2d896dee
   结果: ✅ 成功! Token: {"expire_t****6242f69d"}
```

## 测试验证

运行所有测试确认修复成功：

```bash
go test -v -run TestTokenManager
```

**结果**: 所有TokenManager测试通过 ✅

```bash
go test -v
```

**结果**: 所有测试通过 ✅

## 关键发现

1. **个推官方使用SHA256而不是MD5**进行签名
2. **签名参数顺序**: `appkey + timestamp + master_secret`
3. **时间戳格式**: 毫秒时间戳（13位数字）
4. **签名长度**: SHA256产生64位十六进制字符串

## 影响范围

- ✅ Token获取功能正常工作
- ✅ 推送功能正常工作（依赖Token）
- ✅ 所有测试用例通过
- ✅ 向后兼容性保持

## 建议

1. **参考官方文档**: 在实现第三方API时，务必参考官方文档
2. **全面测试**: 实现多种签名算法进行对比测试
3. **错误分析**: 仔细分析API错误信息，区分配置错误和实现错误
4. **版本控制**: 保持代码版本控制，便于回滚和对比

## 总结

通过系统性的调试和测试，成功识别并解决了签名算法错误。现在个推Go SDK可以正常获取Token并执行推送操作。 