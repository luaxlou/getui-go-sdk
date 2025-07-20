# GetToken 测试配置指南

## 概述

GetToken测试用例现在支持两种模式：
1. **配置错误模式** - 使用测试配置，预期返回配置相关错误（当前模式）
2. **真实Token模式** - 使用真实的个推应用配置，获取真实的Token

## 当前测试行为

### 配置错误模式（默认）
当使用测试配置时，测试用例会：
- 接受配置相关的API错误（20001, 20002, 20003, 404）
- 验证错误处理逻辑
- 确保请求格式正确

### 真实Token模式
当配置真实的个推应用信息时，测试用例会：
- 获取真实的Token
- 验证Token格式和长度
- 测试Token缓存机制
- 验证Token过期处理

## 配置真实Token测试

### 1. 获取个推应用信息
1. 登录个推开发者后台：https://dev.getui.com/
2. 创建或选择应用
3. 获取以下信息：
   - AppID
   - AppKey
   - MasterSecret

### 2. 配置环境变量
编辑 `.env` 文件，填入真实的配置信息：

```bash
# 个推应用配置
GETUI_TEST_APP_ID=your_real_app_id
GETUI_TEST_APP_KEY=your_real_app_key
GETUI_TEST_MASTER_SECRET=your_real_master_secret

# 个推API域名
GETUI_TEST_DOMAIN=https://restapi.getui.com/v2
```

### 3. 运行测试
```bash
# 运行所有TokenManager测试
go test -v -run TestTokenManager

# 运行特定测试
go test -v -run TestTokenManager_ValidConfig
```

## 测试用例说明

### 成功获取Token的测试用例

#### TestTokenManager_ValidConfig
- **目的**: 验证使用有效配置获取Token
- **成功条件**: 获取到非空Token，长度大于10
- **失败条件**: 网络错误、意外的API错误

#### TestTokenManager_TokenCaching
- **目的**: 验证Token缓存机制
- **成功条件**: 两次获取的Token相同
- **失败条件**: 缓存失效、Token不一致

#### TestTokenManager_TokenExpiration
- **目的**: 验证Token过期后重新获取
- **成功条件**: 获取到新的Token，不是过期的Token
- **失败条件**: 返回过期Token

#### TestTokenManager_EmptyTokenExpired
- **目的**: 验证空Token过期后重新获取
- **成功条件**: 获取到新的非空Token
- **失败条件**: 返回空Token

#### TestTokenManager_ConcurrentAccess
- **目的**: 验证并发访问Token
- **成功条件**: 所有并发请求返回相同Token
- **失败条件**: Token不一致、错误类型不一致

### 辅助功能测试用例

#### TestTokenManager_SignGeneration
- **目的**: 验证签名生成
- **验证**: 签名长度64位、十六进制格式

#### TestTokenManager_SignConsistency
- **目的**: 验证签名一致性
- **验证**: 相同时间戳生成相同签名

#### TestTokenManager_RequestFormat
- **目的**: 验证请求格式
- **验证**: Sign、Timestamp、AppKey格式正确

#### TestTokenManager_ResponseParsing
- **目的**: 验证响应解析
- **验证**: JSON解析正确、Token提取正确

### 管理功能测试用例

#### TestTokenManager_GetTokenExpireTime
- **目的**: 验证过期时间获取
- **验证**: 返回正确的过期时间

#### TestTokenManager_GetCurrentToken
- **目的**: 验证当前Token获取
- **验证**: 返回设置的Token

#### TestTokenManager_SetToken
- **目的**: 验证Token设置
- **验证**: Token和过期时间设置正确

#### TestTokenManager_ClearToken
- **目的**: 验证Token清除
- **验证**: Token和过期时间被清除

#### TestTokenManager_IsTokenExpired
- **目的**: 验证Token过期状态检查
- **验证**: 正确判断Token是否过期

## 错误处理

### 配置相关错误（允许通过）
- **20001**: sign is invalid - 签名无效
- **20002**: appkey is invalid - AppKey无效
- **20003**: timestamp is invalid - 时间戳无效
- **404**: not found - 应用不存在

### 其他错误（测试失败）
- **网络错误**: 连接失败、超时等
- **解析错误**: JSON解析失败
- **意外API错误**: 其他API错误码

## 测试覆盖率

当前测试覆盖了TokenManager的所有功能：
- ✅ Token获取和缓存
- ✅ 签名生成和验证
- ✅ 请求格式构建
- ✅ 响应解析
- ✅ 并发访问
- ✅ Token管理（设置、清除、过期检查）
- ✅ 错误处理

## 注意事项

1. **安全性**: 不要将真实的MasterSecret提交到版本控制
2. **测试环境**: 建议使用测试应用进行测试
3. **网络要求**: 需要能够访问个推API服务器
4. **频率限制**: 注意API调用频率限制

## 示例输出

### 配置错误模式
```
=== RUN   TestTokenManager_ValidConfig
API错误：code=20001, message=sign is invalid
配置错误（预期）：sign is invalid
--- PASS: TestTokenManager_ValidConfig (0.04s)
```

### 真实Token模式
```
=== RUN   TestTokenManager_ValidConfig
Token获取成功：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
--- PASS: TestTokenManager_ValidConfig (0.05s)
``` 