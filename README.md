[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# lansenger-sdk-go

Go SDK for the Lansenger (蓝信) platform — supports Lansenger apps, organization bots, and personal bots.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go 1.21+](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![Tests: 115](https://img.shields.io/badge/Tests-115-green)](https://github.com/lansenger-pm/lansenger-sdk-go)

> Zero external dependencies — only Go standard library. Works with any Go project.

## Supported Bot Types

| Bot Type | Auth | WebSocket Inbound | All APIs |
|----------|------|-------------------|----------|
| **Lansenger App** | appToken + userToken | ✗ (uses webhook) | ✓ |
| **Organization Bot** | appToken + userToken | ✗ (uses webhook) | ✓ |
| **Personal Bot** | appToken | ✓ (WebSocket) | ✓ (limited for non-bot APIs) |

All three bot types use the same auth mechanism: `appToken` is required for every API call; `userToken` is only needed for specific user-level operations (user info, staff search, calendar, etc.).

## Features

- **Single client** — `LansengerClient` with `context.Context` for all API calls
- **Credential & token persistence** — `CredentialStore` saves credentials and tokens to JSON file (survives restarts)
- **OAuth2 user authentication** — authorize URL, code exchange, token refresh
- **Organization & departments** — org info, department detail/children/staff
- **Staff & contacts** — basic/detailed info, ID mapping, department ancestors, search
- **Messaging** — 3 private chat channels (bot, official account, user impersonate) + group chat, all message types, @mention, human/bot sender identity
- **Rich cards** — appCard (with dynamic status updates), oacard, linkCard, appArticles
- **Streaming messages** — SSE-based real-time delivery for AI agents
- **Media upload/download** — files, images, videos with auto type detection
- **Message management** — revoke, dynamic card update
- **Groups V2** — create, info, members, list, membership check, update settings & members
- **Calendar & schedule** — primary calendar, schedule CRUD, attendee management
- **Unified todo** — create, update, delete, query, executor management, status counts
- **Callback events** — 24 event types, structured data parsing, signature verification

## Quick Install

```bash
go get lansenger-sdk-go
```

## 1. Authentication

### appToken — Required for all API calls

Every SDK method requires `appToken`. The client automatically obtains and refreshes it using your `appID` + `appSecret`. You never need to manage appToken manually — the `TokenManager` handles the lifecycle:

1. **First call** → `GET /v1/apptoken/create` with appID + appSecret → returns `appToken` (valid 2 hours)
2. **Subsequent calls** → reuse cached appToken until expiry
3. **Token expired** → automatically refresh via the same endpoint

```go
// appToken is managed automatically — just configure appID + appSecret
client := lansenger.NewClient("your-appid", "your-secret")

// You can also get/invalidate token manually
token, err := client.GetToken(ctx)
client.InvalidateToken() // force refresh on next call
```

### userToken — Only needed for specific endpoints

`userToken` represents a specific Lansenger user's authorization (obtained via OAuth2). It's only required for:
- User-level information (FetchUserInfo, FetchStaffDetail, SearchStaff)
- Calendar & schedule operations (FetchPrimaryCalendar, CreateSchedule, etc.)
- Group operations as a human sender

### Getting Credentials

| Bot Type | How to get appID + appSecret |
|----------|------------------------------|
| **Personal Bot** | Lansenger desktop → Contacts → Smart Bots → Personal Bots → click ℹ️ icon (mobile client does NOT show credentials) |
| **Lansenger App** | Create at [Lansenger Developer Center](https://dev.lanxin.cn) — may require organization admin approval |
| **Organization Bot** | Create at [Lansenger Developer Center](https://dev.lanxin.cn) — may require organization admin approval |

### OAuth2 user-level auth

```go
// Build authorize URL — redirect user to Lansenger passport
url := client.BuildAuthorizeURL("https://myapp.com/callback", "", "state123")

// After user authorizes, exchange code for userToken + refreshToken
tokenResult, err := client.ExchangeCode(ctx, "auth_code_from_callback", "https://myapp.com/callback")

// Refresh expired userToken
newToken, err := client.RefreshUserToken(ctx, tokenResult.RefreshToken, "")

// Fetch user profile
userInfo, err := client.FetchUserInfo(ctx, tokenResult.UserToken)
```

## 2. Organization & Departments

```go
// Organization info
org, err := client.FetchOrgInfo(ctx, "orgId", "")

// Department hierarchy
detail, err := client.FetchDepartmentDetail(ctx, "deptId", "", "")
children, err := client.FetchDepartmentChildren(ctx, "deptId", "")
staffs, err := client.FetchDepartmentStaffs(ctx, "deptId", "", 1, 100)
```

## 3. Staff & Contacts

```go
// Basic staff info
staff, err := client.FetchStaffBasicInfo(ctx, "staffOpenId", "")

// Detailed profile (userToken recommended)
detail, err := client.FetchStaffDetail(ctx, "staffOpenId", "ut")

// Map phone → staffId
mapping, err := client.FetchStaffIdMapping(ctx, "orgId", "mobile", "13800138000", "")

// Department ancestors for a staff member
ancestors, err := client.FetchDepartmentAncestors(ctx, "staffOpenId", "")

// Search staff (requires userToken or userID)
results, err := client.SearchStaff(ctx, "Zhang San", "ut", "", true, nil, 1, 10)

// Org extra field IDs
fields, err := client.FetchOrgExtraFieldIDs(ctx, "orgId", "", 1, 1000)
```

## 4. Messaging & Media

#### Bot private chat — most common

```go
result, err := client.SendText(ctx, "staff123", "Hello!", "", 0, false, nil, false, "", "")
result, err := client.SendMarkdown(ctx, "staff123", "**Bold**", false, nil, false, "", "")
result, err := client.SendFile(ctx, "staff123", "/path/to/report.pdf", "", 0, false, "", "")
```

#### Public account channel

```go
result, err := client.SendAccountMessage(ctx, "text",
    map[string]interface{}{"content": "System notice"},
    []string{"staff1", "staff2"}, nil, "524288-xxxx", "", "", "")
```

#### User impersonate channel (requires userToken)

```go
result, err := client.SendUserMessage(ctx, "staff456", "text",
    map[string]interface{}{"content": "Hello"}, "ut", "")
```

#### Group chat

```go
// Bot → group
result, err := client.SendText(ctx, "group123", "Notice", "", 0, false, nil, true, "", "")

// Human → group (with userToken)
result, err := client.SendGroupMessage(ctx, "group123", "text",
    map[string]interface{}{"content": "I'll handle it"}, "ut", "", false, nil, "", "", "")

// @mention in group
result, err := client.SendText(ctx, "group123", "Important!", "", 0, true, nil, true, "", "")
```

#### Rich cards

```go
// appCard
params := &lansenger.AppCardParams{
    ChatID: "staff123", BodyTitle: "Approval", IsDynamic: true,
}
result, err := client.SendAppCardWithParams(ctx, params)

// linkCard
params := &lansenger.LinkCardParams{
    ChatID: "staff123", Title: "Article", Link: "https://...",
}
result, err := client.SendLinkCardWithParams(ctx, params)

// Update dynamic card status
updateParams := &lansenger.DynamicCardUpdateParams{
    MsgID: "msg123", IsLastUpdate: true,
}
result, err := client.UpdateDynamicCard(ctx, updateParams)
```

#### Streaming messages (for AI agents)

```go
result, err := client.CreateStreamMessage(ctx, "staff1", "staff", "stream1")
result, err := client.FetchStreamMessage(ctx, "msg123")
```

#### Media

```go
// Upload
upload, err := client.UploadMedia(ctx, "/path/to/file.pdf", lansenger.MediaTypeFile)

// Download
download, err := client.DownloadMedia(ctx, "media123")

// Download and save to file
path, err := client.DownloadMediaToFile(ctx, "media123", "/path/to/save.pdf")

// Revoke messages
result, err := client.RevokeMessage(ctx, []string{"msg1", "msg2"}, "bot", "")
```

## 5. Groups

```go
// Create group
info := &lansenger.GroupCreateInfo{
    Name: "Project Chat", OrgID: 1, StaffIDList: []string{"s1", "s2", "s3"},
}
group, err := client.CreateGroup(ctx, info, "")

// Fetch info & members
info, err := client.FetchGroupInfo(ctx, "groupOpenId", "")
members, err := client.FetchGroupMembers(ctx, "groupOpenId", "", 0, 100)
groups, err := client.FetchGroupList(ctx, "", 0, 100)

// Check membership
result, err := client.CheckIsInGroup(ctx, "groupOpenId", "", "staff1")

// Update settings
result, err := client.UpdateGroupInfo(ctx, "groupId", map[string]interface{}{"name": "New Name"}, "")

// Add/remove members
result, err := client.UpdateGroupMembers(ctx, "groupId",
    []string{"staff4"}, []string{"staff3"}, nil, "")
```

## 6. Calendar & Schedule

```go
// Get primary calendar (requires userToken or userID)
cal, err := client.FetchPrimaryCalendar(ctx, "ut", "uid1")

// Create schedule
schedule, err := client.CreateSchedule(ctx, cal.CalendarID, "Team Meeting",
    "2024-01-15T09:00", "2024-01-15T10:00", nil, "", false, "", nil, 0, 0, 0, "ut")

// Fetch/delete schedule
info, err := client.FetchSchedule(ctx, "cal1", "sch1", "ut", "")
result, err := client.DeleteSchedule(ctx, "cal1", "sch1", 1, "", "", "ut")

// Schedule list in time range
schedules, err := client.FetchScheduleList(ctx, "cal1", "2024-01-15T00:00", "2024-01-15T23:59", "ut")

// Attendee management
attendees, err := client.FetchScheduleAttendees(ctx, "cal1", "sch1", 1, 10)
result, err := client.AddScheduleAttendees(ctx, "cal1", "sch1",
    []map[string]interface{}{{"staffId": "staff2"}}, 0)
result, err := client.DeleteScheduleAttendees(ctx, "cal1", "sch1",
    []map[string]interface{}{{"staffId": "staff2"}}, 0)
```

## 7. Unified Todo

```go
// Create todo task
todo, err := client.CreateTodoTask(ctx, "Approval Request", lansenger.TodoTypeApproval,
    "https://app.com/a/1", "https://pc.app.com/a/1", []string{"staff1"}, "org1", "", "", "", "")

// Update status (11=pending-read, 12=read, 21=pending-do, 22=done)
result, err := client.UpdateTodoTaskStatus(ctx, "taskId", lansenger.TodoStatusDone, "org1", "", "")

// Update content
result, err := client.UpdateTodoTask(ctx, "taskId", "Updated", "l", "p", "org1", "", "")

// Delete (sender only)
result, err := client.DeleteTodoTask(ctx, "taskId", "org1", "", "")

// Query
list, err := client.FetchTodoTaskList(ctx, "org1", nil, "", nil, "")
task, err := client.FetchTodoTaskByID(ctx, "taskId", "org1", "", "")
task, err := client.FetchTodoTaskBySourceID(ctx, "src1", "org1", "", "")
counts, err := client.FetchTodoTaskStatusCounts(ctx, "staff1", "org1", "", "", "")

// Executor management
result, err := client.AddExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
result, err := client.DeleteExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
executors, err := client.FetchExecutorList(ctx, "taskId", "org1", "", nil, "")
```

## 8. Callback Events

```go
// Parse webhook payload
events, err := lansenger.ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")

// Verify signature
isValid := lansenger.VerifyCallbackSignature(queryString, "app_secret")

// Available event types
types := lansenger.GetCallbackEventTypes() // 24 event types across 14 categories
```

## 9. Chat Reading

```go
// Fetch user's chat list (private + group)
chats, err := client.FetchChatList(ctx, "ut", "private", "", "", "")

// Fetch private chat messages with a specific person
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "s001", "", "", "", "")

// Fetch group chat messages
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "", "g001", "", "", "")
```

## Message Type Capability Matrix

| msgType | Markdown | @mention | Attachments | Private Channels | Group Chat | Notes |
|---------|----------|----------|-------------|------------------|------------|-------|
| `text` | ✗ | ✓ (group) | ✓ | Bot, Official Account, User Impersonate | ✓ | Up to 6000 bytes |
| `formatText` | ✓ | ✗ | ✗ | User Impersonate only | ✓ | Markdown via formatType=1 |
| `oacard` | ✗ | ✗ | ✗ | Bot, Official Account, User Impersonate | ✓ | Simple card with fields |
| `appCard` | ✓ (div tags) | ✗ | ✗ | Bot, Official Account, User Impersonate | ✓ | Rich card, dynamic updates |
| `linkCard` | ✗ | ✗ | ✗ | Bot, Official Account | ✓ | Link preview card |
| `appArticles` | ✗ | ✗ | ✗ | Bot private only | ✓ | Article list (1+ articles) |

**Group chat** supports all message types. Only group chat supports @mention.

## Configuration

### Environment Variables

| Variable | Required | Description | Default |
|----------|----------|-------------|---------|
| `LANSENGER_APP_ID` | ✓ | App/Bot ID | — |
| `LANSENGER_APP_SECRET` | ✓ | App/Bot Secret | — |
| `LANSENGER_API_GATEWAY_URL` | ✗ | API Gateway URL | `https://open.e.lanxin.cn/open/apigw` |
| `LANSENGER_PASSPORT_URL` | ✗ | Passport URL (for OAuth2) | — |
| `LANSENGER_HTTP_TIMEOUT` | ✗ | HTTP timeout (seconds) | `30` |

### From Environment

```go
client, err := lansenger.NewClientFromEnv()
```

### Credential & Token Persistence

By default, credentials and tokens stay in memory only (lost on process exit). Enable file persistence with `CredentialStore`:

```go
// Auto-persist to ~/.lansenger/sdk_state.json (0600 permissions)
store := lansenger.NewCredentialStore("", "default")
store.SaveCredentials("app_id", "app_secret", "https://apigw.lx.qianxin.com", "https://passport.lx.qianxin.com")

// Save tokens
store.SaveAppToken("token123", 7200)
store.SaveUserToken("ut123", "rt123", 7200)

// Load tokens (returns empty string if expired)
token, err := store.LoadAppToken()

// Credentials are shared with Python SDK (same ~/.lansenger/sdk_state.json format)
```

When persistence is enabled:
- **appToken** can be saved and restored on restart (skips redundant API calls)
- **userToken + refreshToken** can be saved after OAuth2 exchange
- **Credentials + URLs** are saved together for full config recovery

## Project Structure

```
lansenger-sdk-go/
├── client.go            # LansengerClient — core client with HTTP helpers
├── config.go            # Config — configuration + env vars
├── constants.go         # API endpoints, media types, callback event types
├── errors.go            # LansengerError hierarchy (Auth/Config/API/Network/File)
├── models.go            # 35+ result/params struct types
├── auth.go              # TokenManager — appToken lifecycle with auto-refresh
├── url_helpers.go       # BuildAPIURL — Options pattern for URL construction
├── oauth.go             # OAuth2 authorize URL, code exchange, token refresh
├── contacts.go          # Staff & org info APIs
├── departments.go       # Department APIs
├── groups.go            # Groups V2 APIs
├── chats.go             # Chat list & messages APIs
├── account_messages.go  # Public account channel (4.6.1)
├── user_messages.go     # User impersonate channel (4.6.3)
├── group_messages.go    # Group chat channel (4.6.2)
├── bot_messages.go      # Bot channel (4.6.12)
├── messaging.go         # Convenience methods + revoke + dynamic update
├── streaming.go         # SSE streaming messages
├── media.go             # Upload/download files & images
├── todos.go             # Unified todo (4.33) — 12 endpoints
├── calendars.go         # Calendar & schedule (4.23) — 8 endpoints
├── callbacks.go         # Callback event parsing + signature verification
├── persistence.go       # CredentialStore — JSON file persistence
├── *_test.go            # 115 unit tests + 10 integration tests
├── go.mod
└── README.md
```

## Development

```bash
go test ./... -v                    # unit tests (115 tests)
go test ./... -run TestIntegration  # integration tests (10 tests, requires ~/.lansenger/sdk_state.json)
```

## License

MIT — see [LICENSE](LICENSE).