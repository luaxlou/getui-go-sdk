# 个推Go SDK

该SDK以个推官方Java SDK为蓝本，使用Go语言实现，提供了完整的个推推送服务API封装。

个推官方Java SDK：https://github.com/GetuiLaboratory/getui-pushapi-java-client-v2

## 功能特性

- 支持个推推送服务V2版本API
- 支持单推、批量推送、群推等推送方式
- 支持Android、iOS、鸿蒙等多平台推送
- 支持透传消息和通知消息
- 支持定时推送和任务管理
- 支持用户管理和统计分析
- 完整的错误处理和重试机制

## 环境要求

- Go 1.21 或更高版本
- 个推开发者账号和应用配置

## 安装

```bash
go get github.com/luaxlou/getui-go-sdk
```

## 快速开始

### 1. 创建客户端

```go
package main

import (
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
    
    // 使用客户端进行推送
    // ...
}
```

### 2. 单推示例

```go
// 创建推送消息
pushMessage := &getui.PushMessage{
    Notification: &getui.Notification{
        Title:    "推送标题",
        Body:     "推送内容",
        ClickType: "url",
        URL:      "https://www.getui.com",
    },
}

// 创建推送请求
pushDTO := &getui.PushDTO{
    RequestID:   "unique_request_id",
    PushMessage: pushMessage,
    Audience: &getui.Audience{
        CIDs: []string{"target_cid"},
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
```

### 3. 批量推送示例

```go
// 创建批量推送请求
batchDTO := &getui.PushBatchDTO{
    RequestID:   "batch_request_id",
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
```

### 4. 群推示例

```go
// 创建群推请求
pushDTO := &getui.PushDTO{
    RequestID:   "all_push_request_id",
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
    log.Printf("群推成功，任务ID: %s", result.Data.TaskID)
} else {
    log.Printf("群推失败: code=%d, msg=%s", result.Code, result.Msg)
}
```

## API 接口

### PushAPI - 推送相关接口

- `PushToSingleByCID(pushDTO *PushDTO) (*ApiResult, error)` - 根据CID单推
- `PushToSingleByAlias(pushDTO *PushDTO) (*ApiResult, error)` - 根据别名单推
- `PushBatchByCID(batchDTO *PushBatchDTO) (*ApiResult, error)` - 根据CID批量推送
- `PushBatchByAlias(batchDTO *PushBatchDTO) (*ApiResult, error)` - 根据别名批量推送
- `PushAll(pushDTO *PushDTO) (*ApiResult, error)` - 群推
- `PushByTag(pushDTO *PushDTO) (*ApiResult, error)` - 根据标签推送
- `StopPush(taskID string) (*ApiResult, error)` - 停止推送任务
- `QueryScheduleTask(taskID string) (*ApiResult, error)` - 查询定时任务

### UserAPI - 用户管理接口

- `QueryUserStatus(cids []string) (*ApiResult, error)` - 查询用户状态
- `QueryAliasByCID(cid string) (*ApiResult, error)` - 根据CID查询别名
- `QueryCIDByAlias(alias string) (*ApiResult, error)` - 根据别名查询CID
- `BindAlias(alias string, cid string) (*ApiResult, error)` - 绑定别名
- `UnbindAlias(alias string, cid string) (*ApiResult, error)` - 解绑别名

### StatisticAPI - 统计分析接口

- `QueryPushResultByTaskIDs(taskIDs []string) (*ApiResult, error)` - 根据任务ID查询推送结果
- `QueryPushResultByDate(date string) (*ApiResult, error)` - 根据日期查询推送结果

## 配置选项

```go
config := &getui.Config{
    AppID:        "your_app_id",
    AppKey:       "your_app_key",
    MasterSecret: "your_master_secret",
    Domain:       "https://restapi.getui.com/v2",
    
    // 可选配置
    SocketTimeout:           30000,  // HTTP读取超时时间(ms)
    ConnectTimeout:          10000,  // HTTP连接超时时间(ms)
    MaxHTTPTryTime:          1,      // HTTP重试次数
    TrustSSL:               false,  // 是否信任SSL证书
    OpenAnalyseStableDomain: true,   // 是否开启稳定域名检测
}
```

## 错误处理

SDK提供了完整的错误处理机制：

```go
result, err := client.PushAPI.PushToSingleByCID(pushDTO)
if err != nil {
    // 网络错误或其他异常
    log.Printf("请求异常: %v", err)
    return
}

if !result.IsSuccess() {
    // API返回错误
    log.Printf("API错误: code=%d, msg=%s", result.Code, result.Msg)
    return
}

// 成功处理
log.Printf("推送成功: %+v", result.Data)
```

## 测试

### 环境变量配置

为了运行测试，您需要设置以下环境变量。我们提供了便捷的脚本来自动设置：

#### 使用脚本设置（推荐）

**Linux/macOS:**
```bash
# 首次设置
source scripts/setup_env.sh

# 后续加载
source scripts/load_env.sh
```

**Windows:**
```cmd
# 首次设置
scripts\setup_env.bat

# 后续加载
scripts\load_env.bat
```

#### 使用.env文件（推荐）

创建 `.env` 文件（参考 `env.example`）：

```bash
cp env.example .env
# 编辑 .env 文件，填入实际的测试配置
```

#### 手动设置环境变量

```bash
# 设置测试环境变量
export GETUI_TEST_APP_ID="your_test_app_id"
export GETUI_TEST_APP_KEY="your_test_app_key"
export GETUI_TEST_MASTER_SECRET="your_test_master_secret"
```

#### 脚本功能特性

- 🔧 **交互式配置**: 引导式输入个推配置信息
- ✅ **自动验证**: 确保输入不为空
- 💾 **配置保存**: 可选择保存到.env文件
- 🔄 **重复使用**: 支持重新加载已保存的配置
- 🛡️ **安全保护**: .env文件已添加到.gitignore中

#### .env文件特性

- 📁 **自动加载**: 测试文件自动从.env文件读取配置
- 🔒 **安全存储**: 敏感配置信息本地化存储
- 📝 **格式支持**: 支持注释、空行、引号等格式
- 🛡️ **容错处理**: 文件不存在时自动使用默认配置

### 运行测试

```bash
# 运行所有测试
go test -v

# 运行特定测试
go test -v -run TestNewClient

# 运行测试并显示覆盖率
go test -v -cover
```

## 许可证

MIT License