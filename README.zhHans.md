[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# lansenger-sdk-go

Lansenger（蓝信）平台的 Go SDK — 支持蓝信应用、组织机器人和个人机器人。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version: 0.9.0](https://img.shields.io/badge/Version-0.9.0-blue)](https://github.com/lansenger-pm/lansenger-sdk-go)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-blue)](https://go.dev/)
[![Tests: 146](https://img.shields.io/badge/Tests-146-green)](https://github.com/lansenger-pm/lansenger-sdk-go)

> SDK 库零外部依赖（仅标准库）。CLI（`cmd/lansenger`）使用 [spf13/cobra](https://github.com/spf13/cobra)。

## 支持的机器人类型

| 机器人类型 | 认证 | WebSocket 入站 | 全部 API |
|----------|------|-----------------|----------|
| **蓝信应用** | appToken + userToken | ✗（使用 webhook） | ✓ |
| **组织机器人** | appToken + userToken | ✗（使用 webhook） | ✓ |
| **个人机器人** | appToken | ✓（WebSocket） | ✓（非机器人 API 受限） |

三种机器人类型使用相同的认证机制：每次 API 调用都需要 `appToken`；仅在特定用户级别操作（用户信息、员工搜索、日历等）时才需要 `userToken`。

## 功能特性

- **单一客户端** — `LansengerClient` 配合 `context.Context` 用于所有 API 调用
- **凭据与令牌持久化** — `CredentialStore` 将凭据和令牌保存到 JSON 文件（重启后仍有效）
- **OAuth2 用户认证** — 授权 URL、代码交换、令牌刷新
- **组织与部门** — 组织信息、部门详情/子部门/员工
- **员工与联系人** — 基本/详细信息、ID 映射、部门祖先链、搜索
- **消息发送** — 3 种私聊通道（机器人、官方账号、用户代发）+ 群聊，支持所有消息类型、@提及、人类/机器人发送者身份
- **富卡片** — appCard（支持动态状态更新）、oacard、linkCard、appArticles
- **流式消息** — 基于 SSE 的实时投递，适用于 AI 代理
- **媒体上传/下载** — 文件、图片、视频，自动类型检测，获取下载路径，应用/机器人媒体上传
- **消息管理** — 撤回、动态卡片更新、加急提醒
- **群组 V2** — 创建、信息、成员、列表、成员检查、更新设置与成员、解散
- **日历与日程** — 主日历、日程 CRUD、参会人管理、参会人元数据更新
- **统一待办** — 创建、更新、删除、查询、执行人管理、状态计数
- **回调事件** — 24 种事件类型、结构化数据解析、签名验证

## 快速安装

**SDK（库）**:
```bash
go get github.com/lansenger-pm/lansenger-sdk-go
```

**CLI（用于 AI 代理和调试）**:
```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
lansenger version
```

CLI 通过 `~/.lansenger/sdk_state.json` 与 SDK 共享凭据。安装后配置凭据：
```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
```

## 1. 认证

### appToken — 所有 API 调用必需

每个 SDK 方法都需要 `appToken`。客户端使用你的 `appID` + `appSecret` 自动获取并刷新令牌。你无需手动管理 appToken — `TokenManager` 处理整个生命周期：

1. **首次调用** → `GET /v1/apptoken/create` 使用 appID + appSecret → 返回 `appToken`（有效期 2 小时）
2. **后续调用** → 复用缓存的 appToken 直到过期
3. **令牌过期** → 通过相同端点自动刷新

```go
// appToken 自动管理 — 只需配置 appID + appSecret
client := lansenger.NewClient("your-appid", "your-secret")

// 也可以手动获取/失效令牌
token, err := client.GetToken(ctx)
client.InvalidateToken() // 强制下次调用时刷新
```

### userToken — 仅特定端点需要

`userToken` 代表特定蓝信用户的授权（通过 OAuth2 获取）。仅在以下场景需要：
- 用户级别信息（FetchUserInfo、FetchStaffDetail、SearchStaff）
- 日历与日程操作（FetchPrimaryCalendar、CreateSchedule 等）
- 以人类发送者身份进行群组操作

### 获取凭据

| 机器人类型 | 如何获取 appID + appSecret |
|----------|----------------------------|
| **个人机器人** | 蓝信桌面端 → 联系人 → 智能机器人 → 个人机器人 → 点击 ℹ️ 图标（移动端不显示凭据） |
| **蓝信应用** | 在 [蓝信开发者中心](https://dev.lanxin.cn) 创建 — 可能需要组织管理员审批 |
| **组织机器人** | 在 [蓝信开发者中心](https://dev.lanxin.cn) 创建 — 可能需要组织管理员审批 |

### OAuth2 用户级别认证

```go
// 构建授权 URL — 将用户重定向到蓝信通行证
url := client.BuildAuthorizeURL("https://myapp.com/callback", "", "state123")

// 用户授权后，用 code 交换 userToken + refreshToken
tokenResult, err := client.ExchangeCode(ctx, "auth_code_from_callback", "https://myapp.com/callback")

// 刷新过期的 userToken
newToken, err := client.RefreshUserToken(ctx, tokenResult.RefreshToken, "")

// 获取用户信息
userInfo, err := client.FetchUserInfo(ctx, tokenResult.UserToken)
```

## 2. 组织与部门

```go
// 组织信息
org, err := client.FetchOrgInfo(ctx, "orgId", "")

// 部门层级
detail, err := client.FetchDepartmentDetail(ctx, "deptId", "", "")
children, err := client.FetchDepartmentChildren(ctx, "deptId", "")
staffs, err := client.FetchDepartmentStaffs(ctx, "deptId", "", 1, 100)
```

## 3. 员工与联系人

```go
// 基本员工信息
staff, err := client.FetchStaffBasicInfo(ctx, "staffOpenId", "")

// 详细资料（推荐使用 userToken）
detail, err := client.FetchStaffDetail(ctx, "staffOpenId", "ut")

// 手机号 → staffId 映射
mapping, err := client.FetchStaffIdMapping(ctx, "orgId", "mobile", "13800138000", "")

// 员工的部门祖先链
ancestors, err := client.FetchDepartmentAncestors(ctx, "staffOpenId", "")

// 搜索员工（需要 userToken 或 userID）
results, err := client.SearchStaff(ctx, "张三", "ut", "", true, nil, 1, 10)

// 组织扩展字段 ID
fields, err := client.FetchOrgExtraFieldIDs(ctx, "orgId", "", 1, 1000)
```

## 4. 消息与媒体

#### 机器人私聊 — 最常用

```go
result, err := client.SendText(ctx, "staff123", "你好！", "", 0, "", false, nil, false, "", "")
result, err := client.SendMarkdown(ctx, "staff123", "**加粗**", false, nil, false, "", "")
result, err := client.SendFile(ctx, "staff123", "/path/to/report.pdf", "", 0, "", false, "", "")
```

#### 公共账号通道

```go
result, err := client.SendAccountMessage(ctx, "text",
    map[string]interface{}{"content": "系统通知"},
    []string{"staff1", "staff2"}, nil, "524288-xxxx", "", "", "")
```

#### 用户代发通道（需要 userToken）

```go
result, err := client.SendUserMessage(ctx, "staff456", "text",
    map[string]interface{}{"content": "你好"}, "ut", "")
```

#### 群聊

```go
// 机器人 → 群组
result, err := client.SendText(ctx, "group123", "通知", "", 0, "", false, nil, true, "", "")

// 人类 → 群组（使用 userToken）
result, err := client.SendGroupMessage(ctx, "group123", "text",
    map[string]interface{}{"content": "我来处理"}, "ut", "", false, nil, "", "", "")

// 群聊 @提及
result, err := client.SendText(ctx, "group123", "重要！", "", 0, "", true, nil, true, "", "")
```

#### 富卡片

```go
// appCard
params := &lansenger.AppCardParams{
    ChatID: "staff123", BodyTitle: "审批", IsDynamic: true,
}
result, err := client.SendAppCardWithParams(ctx, params)

// linkCard
params := &lansenger.LinkCardParams{
    ChatID: "staff123", Title: "文章", Link: "https://...",
}
result, err := client.SendLinkCardWithParams(ctx, params)

// 更新动态卡片状态
updateParams := &lansenger.DynamicCardUpdateParams{
    MsgID: "msg123", IsLastUpdate: true,
}
result, err := client.UpdateDynamicCard(ctx, updateParams)
```

#### 流式消息（用于 AI 代理）

```go
result, err := client.CreateStreamMessage(ctx, "staff1", "staff", "stream1")
result, err := client.FetchStreamMessage(ctx, "msg123")
```

#### 媒体

```go
// 上传（核心服务 — 数字类型）
upload, err := client.UploadMedia(ctx, "/path/to/file.pdf", lansenger.MediaTypeFile)

// 上传（应用/机器人 — 字串类型，支持 width/height/duration）
upload, err := client.UploadAppMedia(ctx, "/path/to/video.mp4",
    lansenger.AppMediaTypeVideo, 680, 480, 300)

// 下载
download, err := client.DownloadMedia(ctx, "media123")

// 下载并保存到文件
path, err := client.DownloadMediaToFile(ctx, "media123", "/path/to/save.pdf")

// 获取下载路径信息
pathInfo, err := client.FetchMediaPath(ctx, "media123", "ut")

// 撤回消息
result, err := client.RevokeMessage(ctx, []string{"msg1", "msg2"}, "bot", "")

// 发送加急提醒
result, err := client.SendReminder(ctx, "msg123", []int{1, 2}, []string{"staff1", "staff2"})
```

## 5. 群组

```go
// 创建群组
info := &lansenger.GroupCreateInfo{
    Name: "项目聊天", OrgID: 1, StaffIDList: []string{"s1", "s2", "s3"},
}
group, err := client.CreateGroup(ctx, info, "")

// 获取信息与成员
info, err := client.FetchGroupInfo(ctx, "groupOpenId", "")
members, err := client.FetchGroupMembers(ctx, "groupOpenId", "", 0, 100)
groups, err := client.FetchGroupList(ctx, "", 0, 100)

// 检查成员关系
result, err := client.CheckIsInGroup(ctx, "groupOpenId", "", "staff1")

// 更新设置
result, err := client.UpdateGroupInfo(ctx, "groupId", map[string]interface{}{"name": "新名称"}, "")

// 添加/移除成员
result, err := client.UpdateGroupMembers(ctx, "groupId",
    []string{"staff4"}, []string{"staff3"}, nil, "")

// 解散群组
result, err := client.DissolveGroup(ctx, "groupId", "ut")
```

## 6. 日历与日程

```go
// 获取主日历（需要 userToken 或 userID）
cal, err := client.FetchPrimaryCalendar(ctx, "ut", "uid1")

// 创建日程（startTime/endTime 为 map 对象，allDay 为 "yes"/"no")
schedule, err := client.CreateSchedule(ctx, cal.CalendarID, "团队会议",
    map[string]interface{}{"time": "2024-01-15T09:00"},
    map[string]interface{}{"time": "2024-01-15T10:00"},
    nil, "", "no", "", nil, "", "", "", "ut", "")

// 获取/删除日程
info, err := client.FetchSchedule(ctx, "cal1", "sch1", "ut", "")
result, err := client.DeleteSchedule(ctx, "cal1", "sch1", "", "", "", "ut", "")

// 更新日程
result, err := client.UpdateSchedule(ctx, "cal1", "sch1",
    map[string]interface{}{"summary": "更新后的会议"}, "ut", "")

// 时间范围内的日程列表
schedules, err := client.FetchScheduleList(ctx, "cal1",
    map[string]interface{}{"time": "2024-01-15T00:00"},
    map[string]interface{}{"time": "2024-01-15T23:59"}, "ut", "")

// 参会人管理（参会人为 []string）
attendees, err := client.FetchScheduleAttendees(ctx, "cal1", "sch1", 1, 10, "ut", "")
result, err := client.AddScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")
result, err := client.DeleteScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")

// 更新参会人元数据
result, err := client.UpdateScheduleAttendeeMeta(ctx, "cal1", "sch1",
    map[string]interface{}{"rsvpStatus": "accepted"}, "ut", "")
```

## 7. 统一待办

```go
// 创建待办任务
todo, err := client.CreateTodoTask(ctx, "审批请求", lansenger.TodoTypeApproval,
    "https://app.com/a/1", "https://pc.app.com/a/1", []string{"staff1"}, "org1", "", "", "", "")

// 更新状态（11=待读, 12=已读, 21=待办, 22=已完成）
result, err := client.UpdateTodoTaskStatus(ctx, "taskId", lansenger.TodoStatusDone, "org1", "", "")

// 更新内容
result, err := client.UpdateTodoTask(ctx, "taskId", "已更新", "l", "p", "org1", "", "")

// 删除（仅发送者）
result, err := client.DeleteTodoTask(ctx, "taskId", "org1", "", "")

// 查询
list, err := client.FetchTodoTaskList(ctx, "org1", nil, "", nil, "")
task, err := client.FetchTodoTaskByID(ctx, "taskId", "org1", "", "")
task, err := client.FetchTodoTaskBySourceID(ctx, "src1", "org1", "", "")
counts, err := client.FetchTodoTaskStatusCounts(ctx, "staff1", "org1", "", "", "")

// 执行人管理
result, err := client.AddExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
result, err := client.DeleteExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
executors, err := client.FetchExecutorList(ctx, "taskId", "org1", "", nil, "")
```

## 8. 回调事件

```go
// 解析明文（未加密）webhook 数据 — URL 查询字符串或 JSON
events, err := lansenger.ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")

// 解析明文 JSON 回调
events, err = lansenger.ParseCallbackPayload(`{"events":[{"eventType":"staff_modify","data":{"staffId":"s001"}}],"orgId":"org1","appId":"app1"}`)

// 解密加密回调数据（AES-256-CBC）
result, err := lansenger.DecryptCallbackPayload(encryptedData, encodingKey, knownAppID)
fmt.Println(result.OrgID, result.AppID, result.Events)

// 验证签名（基于 SHA1，匹配蓝信协议）
valid := lansenger.VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, callbackToken)

// 可用事件类型（24 种，结构化字段映射）
types := lansenger.GetCallbackEventTypes()
```

## 9. 聊天阅读

```go
// 获取用户聊天列表（私聊 + 群聊）
chats, err := client.FetchChatList(ctx, "ut", "private", "", "", "")

// 获取与特定人的私聊消息
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "s001", "", "", "", "")

// 获取群聊消息
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "", "g001", "", "", "")
```

## 消息类型能力矩阵

| msgType | Markdown | @提及 | 附件 | 私聊通道 | 群聊 | 备注 |
|---------|----------|--------|------|----------|------|------|
| `text` | ✗ | ✓（群聊） | ✓ | 机器人、官方账号、用户代发 | ✓ | 最多 6000 字节 |
| `formatText` | ✓ | ✗ | ✗ | 仅用户代发 | ✓ | 通过 formatType=1 实现 Markdown |
| `oacard` | ✗ | ✗ | ✗ | 机器人、官方账号、用户代发 | ✓ | 带字段的简单卡片 |
| `appCard` | ✓（div 标签） | ✗ | ✗ | 机器人、官方账号、用户代发 | ✓ | 富卡片，支持动态更新 |
| `linkCard` | ✗ | ✗ | ✗ | 机器人、官方账号 | ✓ | 链接预览卡片 |
| `appArticles` | ✗ | ✗ | ✗ | 仅机器人私聊 | ✓ | 文章列表（1+ 篇文章） |

**群聊**支持所有消息类型。只有群聊支持 @提及。

## 配置

### 凭证概览

所有凭证按 profile 持久化存储在 `~/.lansenger/sdk_state.json`（0600 权限）：

| 凭证 | 必填 | CLI 键名 | 说明 |
|------|------|----------|------|
| App ID | ✓ | `app_id` | 蓝信应用/机器人 ID |
| App Secret | ✓ | `app_secret` | 蓝信应用/机器人密钥 |
| API Gateway URL | ✓ | `api_gateway_url` | API 网关地址（默认：`https://open.e.lanxin.cn/open/apigw`） |
| Passport URL | 仅 OAuth2 | `passport_url` | OAuth2 授权页地址 |
| Encoding Key | 仅回调 | `encoding_key` | AES-256-CBC 解密密钥 |
| Callback Token | 仅回调 | `callback_token` | 回调签名验证令牌 |

### CLI 配置

```bash
# 第1步：设置必填凭证
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw

# 第2步（可选）：设置 OAuth2 授权页地址（获取 userToken 需要）
lansenger config set passport_url https://passport.lx.qianxin.com

# 第3步（可选）：设置回调凭证（接收 Webhook 回调需要）
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN

# 验证配置
lansenger config show

# 多 profile 支持（如不同组织/应用）
lansenger config set app_id APP2_ID --profile org2
lansenger config set app_secret APP2_SECRET --profile org2
lansenger --profile org2 staff basic-info STAFF_ID
```

### SDK 配置

**代码方式**（直接传入）：
```go
client := lansenger.NewClient("app_id", "app_secret")
// 如需自定义网关地址
cfg := lansenger.NewConfig("app_id", "app_secret")
cfg.APIGatewayURL = "https://custom-gateway.example.com"
cfg.PassportURL = "https://passport.example.com"
cfg.EncodingKey = "your_encoding_key"
cfg.CallbackToken = "your_callback_token"
client := lansenger.NewClientWithConfig(cfg)
```

**环境变量方式**（自动检测）：

| 变量 | 必填 | 说明 | 默认值 |
|------|------|------|--------|
| `LANSENGER_APP_ID` | ✓ | 应用/机器人 ID | — |
| `LANSENGER_APP_SECRET` | ✓ | 应用/机器人密钥 | — |
| `LANSENGER_API_GATEWAY_URL` | ✗ | API 网关地址 | `https://open.e.lanxin.cn/open/apigw` |
| `LANSENGER_PASSPORT_URL` | ✗ | 授权页地址（OAuth2） | — |
| `LANSENGER_ENCODING_KEY` | ✗ | 回调解密密钥 | — |
| `LANSENGER_CALLBACK_TOKEN` | ✗ | 回调验证令牌（默认同 encoding_key） | — |
| `LANSENGER_HTTP_TIMEOUT` | ✗ | HTTP 超时（秒） | `30` |

```go
client, err := lansenger.NewClientFromEnv()
```

### 凭证与令牌持久化

默认情况下，凭证和令牌仅在内存中保留（进程退出后丢失）。使用 `CredentialStore` 启用文件持久化：

```go
// 自动持久化到 ~/.lansenger/sdk_state.json（0600 权限）
store := lansenger.NewCredentialStore("", "default")
store.SaveCredentials("app_id", "app_secret", "https://apigw.lx.qianxin.com", "https://passport.lx.qianxin.com")
store.SaveCallbackConfig("encoding_key", "callback_token")

// 保存令牌
store.SaveAppToken("token123", 7200)
store.SaveUserToken("ut123", "rt123", 7200)

// 加载令牌（过期时返回空字符串）
token, err := store.LoadAppToken()

// 凭证与 Python SDK 共享（相同 ~/.lansenger/sdk_state.json 格式）
```

启用持久化后：
- **appToken** 可在重启后保存与恢复（跳过冗余 API 请求）
- **userToken + refreshToken** 可在 OAuth2 交换后保存
- **凭证 + URL** 一并保存，完整恢复配置

## 项目结构

```
lansenger-sdk-go/
├── client.go            # LansengerClient — 核心客户端与 HTTP 辅助方法
├── config.go            # Config — 配置 + 环境变量
├── constants.go         # API 端点、媒体类型、回调事件类型
├── errors.go            # LansengerError 层级（Auth/Config/API/Network/File）
├── models.go            # 50+ 结果/参数结构体类型
├── auth.go              # TokenManager — appToken 生命周期与自动刷新
├── user_token_manager.go # UserTokenManager — userToken 生命周期与自动刷新
├── url_helpers.go       # BuildAPIURL — Options 模式构建 URL
├── oauth.go             # OAuth2 授权 URL、代码交换、令牌刷新
├── contacts.go          # 员工与组织信息 API
├── users.go             # 用户资料 API
├── departments.go       # 部门 API
├── groups.go            # 群组 V2 API
├── chats.go             # 聊天列表与消息 API
├── account_messages.go  # 公共账号通道（4.6.1）
├── user_messages.go     # 用户代发通道（4.6.3）
├── group_messages.go    # 群聊通道（4.6.2）
├── bot_messages.go      # 机器人通道（4.6.12）
├── messaging.go         # 便捷方法 + 撤回 + 动态更新
├── streaming.go         # SSE 流式消息
├── media.go             # 上传/下载文件与图片
├── todos.go             # 统一待办（4.33）— 12 个端点
├── calendars.go         # 日历与日程（4.23）— 10 个端点
├── callbacks.go         # 回调事件解析 + AES-256-CBC 解密 + SHA1 签名验证
├── persistence.go       # CredentialStore — JSON 文件持久化
├── version.go           # SDK 版本常量
├── *_test.go            # 136 单元测试 + 10 集成测试
├── cmd/lansenger/       # CLI 工具（config、oauth、消息、员工等）
├── go.mod
└── README.md
```

## 开发

```bash
go test . -v                        # 全部测试（136 单元 + 10 集成）
go test . -run TestIntegration      # 仅集成测试（需要 ~/.lansenger/sdk_state.json + 网络）
```

## 许可证

MIT — 详见 [LICENSE](LICENSE)。