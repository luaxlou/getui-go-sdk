# 个推Go SDK 测试文档

## 测试概述

本SDK包含完整的单元测试和集成测试，确保代码质量和功能正确性。

## 测试覆盖率

当前测试覆盖率为 **32.6%**，主要覆盖了核心功能和错误处理逻辑。

## GetToken 测试用例

### 测试用例列表

#### 1. TestGetToken_ValidConfig
**目的**: 测试有效配置下的Token获取
**验证点**:
- 请求格式正确
- 能够发送到个推服务器
- 正确处理API响应

#### 2. TestGetToken_TokenCaching
**目的**: 测试Token缓存机制
**验证点**:
- 第一次获取Token后缓存
- 第二次获取使用缓存的Token
- 缓存Token一致性

#### 3. TestGetToken_InvalidConfig
**目的**: 测试无效配置处理
**验证点**:
- 空AppID导致panic
- 配置验证机制正确
- 错误处理及时

#### 4. TestGetToken_NetworkError
**目的**: 测试网络错误处理
**验证点**:
- 无效域名导致网络错误
- 网络错误类型正确
- 错误信息准确

#### 5. TestGetToken_SignGeneration
**目的**: 测试签名生成功能
**验证点**:
- 签名不为空
- MD5签名长度32位
- 签名格式为十六进制

#### 6. TestGetToken_SignConsistency
**目的**: 测试签名一致性
**验证点**:
- 相同时间戳生成相同签名
- 签名算法稳定性
- 结果可重现

#### 7. TestGetToken_TokenExpiration
**目的**: 测试Token过期处理
**验证点**:
- 过期Token重新获取
- 过期时间判断正确
- 新Token替换旧Token

#### 8. TestGetToken_ValidTokenNotExpired
**目的**: 测试有效Token缓存使用
**验证点**:
- 未过期Token直接返回
- 缓存机制正常工作
- 不重复请求服务器

#### 9. TestGetToken_EmptyTokenExpired
**目的**: 测试空Token处理
**验证点**:
- 空Token重新获取
- 空Token判断逻辑正确
- 新Token获取成功

#### 10. TestGetToken_RequestFormat
**目的**: 测试请求格式验证
**验证点**:
- 认证请求格式正确
- 请求头设置正确
- 请求体格式正确

#### 11. TestGetToken_ResponseParsing
**目的**: 测试响应解析
**验证点**:
- Token响应解析正确
- 错误响应处理正确
- 数据格式验证

#### 12. TestGetToken_ConcurrentAccess
**目的**: 测试并发访问
**验证点**:
- 多协程并发获取Token
- 线程安全
- 并发错误处理

## PushToSingleByCID 测试用例

### 测试用例列表

#### 1. TestPushToSingleByCID_ValidRequest
**目的**: 测试有效的推送请求
**验证点**:
- 请求格式正确
- 能够发送到个推服务器
- 正确处理API响应

#### 2. TestPushToSingleByCID_EmptyRequestID
**目的**: 测试自动生成RequestID功能
**验证点**:
- 当RequestID为空时，自动生成
- 生成的RequestID不为空
- 请求能够正常发送

#### 3. TestPushToSingleByCID_InvalidRequestID
**目的**: 测试RequestID验证
**验证点**:
- 过短的RequestID被拒绝
- 返回正确的错误类型
- 错误信息准确

#### 4. TestPushToSingleByCID_NilPushDTO
**目的**: 测试空指针处理
**验证点**:
- 传入nil时返回错误
- 错误信息明确
- 不会导致panic

#### 5. TestPushToSingleByCID_EmptyAudience
**目的**: 测试受众验证
**验证点**:
- Audience为空时返回错误
- 返回ErrEmptyAudience错误类型
- 错误处理正确

#### 6. TestPushToSingleByCID_EmptyPushMessage
**目的**: 测试推送消息验证
**验证点**:
- PushMessage为空时返回错误
- 返回ErrEmptyPushMessage错误类型
- 验证逻辑正确

#### 7. TestPushToSingleByCID_TransmissionMessage
**目的**: 测试透传消息
**验证点**:
- 透传消息格式正确
- 能够发送到服务器
- 响应处理正确

#### 8. TestPushToSingleByCID_WithSettings
**目的**: 测试带设置的推送
**验证点**:
- TTL设置正确
- 推送策略配置正确
- 设置参数传递正确

#### 9. TestPushToSingleByCID_WithPushChannel
**目的**: 测试厂商通道推送
**验证点**:
- Android厂商通道配置
- iOS APNS配置
- 通道参数正确传递

#### 10. TestPushToSingleByCID_MultipleCIDs
**目的**: 测试多CID推送
**验证点**:
- 多个CID正确处理
- 批量推送逻辑正确
- 响应处理正确

#### 11. TestPushToSingleByCID_WithTaskName
**目的**: 测试任务名称设置
**验证点**:
- 任务名称正确设置
- 参数传递正确
- 不影响其他功能

## 运行测试

### 基本测试
```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestPushToSingleByCID

# 运行测试并显示覆盖率
go test -v -cover
```

### 环境变量配置
```bash
# 设置测试环境变量
export GETUI_TEST_APP_ID="your_test_app_id"
export GETUI_TEST_APP_KEY="your_test_app_key"
export GETUI_TEST_MASTER_SECRET="your_test_master_secret"

# 运行测试
go test -v
```

### 生成覆盖率报告
```bash
# 生成覆盖率文件
go test -coverprofile=coverage.out

# 生成HTML报告
go tool cover -html=coverage.out -o coverage.html

# 查看覆盖率详情
go tool cover -func=coverage.out
```

## 测试策略

### 单元测试
- 验证输入参数验证逻辑
- 测试错误处理机制
- 确保边界条件处理正确

### 集成测试
- 测试与个推API的实际交互
- 验证认证机制
- 测试网络错误处理

### 错误处理测试
- 测试各种错误场景
- 验证错误信息准确性
- 确保错误类型正确

## 测试数据

### 测试配置
- 使用环境变量获取真实配置
- 提供默认测试配置作为回退
- 支持不同环境的配置切换

### 测试消息
- 标准通知消息
- 透传消息
- 厂商通道消息
- 各种点击类型

### 测试受众
- 单个CID
- 多个CID
- 不同受众类型

## 持续集成

### 自动化测试
- 每次提交自动运行测试
- 覆盖率报告生成
- 测试结果通知

### 质量门禁
- 测试覆盖率不低于20%
- 所有测试必须通过
- 无严重错误或警告

## 故障排除

### 常见问题

1. **API错误 20001: appid is invalid**
   - 检查环境变量配置
   - 确认AppID正确
   - 验证应用状态

2. **网络连接错误**
   - 检查网络连接
   - 验证域名配置
   - 确认防火墙设置

3. **认证失败**
   - 检查AppKey和MasterSecret
   - 验证签名算法
   - 确认时间戳正确

### 调试技巧

1. **启用详细日志**
   ```bash
   go test -v -run TestPushToSingleByCID
   ```

2. **查看覆盖率详情**
   ```bash
   go tool cover -func=coverage.out
   ```

3. **分析测试失败**
   - 检查错误信息
   - 验证测试数据
   - 确认环境配置

## 扩展测试

### 添加新测试用例
1. 在相应的测试文件中添加测试函数
2. 遵循命名规范：`TestFunctionName_Scenario`
3. 包含清晰的测试目的和验证点
4. 添加适当的错误处理

### 测试最佳实践
1. 每个测试用例只测试一个功能点
2. 使用描述性的测试名称
3. 包含正面和负面测试场景
4. 验证错误处理和边界条件
5. 保持测试代码的可读性和可维护性 