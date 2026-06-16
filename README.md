[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# lansenger-sdk-go

Go SDK for the Lansenger (蓝信) platform — supports Lansenger apps, organization bots, and personal bots.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version: 0.9.18](https://img.shields.io/badge/Version-0.9.18-blue)](https://github.com/lansenger-pm/lansenger-sdk-go)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-blue)](https://go.dev/)
[![Tests: 148](https://img.shields.io/badge/Tests-148-green)](https://github.com/lansenger-pm/lansenger-sdk-go)

> SDK library has zero external dependencies (stdlib only). The CLI (`cmd/lansenger`) uses [spf13/cobra](https://github.com/spf13/cobra).

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
- **Media upload/download** — files, images, videos with auto type detection, fetch download path, app/bot media upload
- **Message management** — revoke, dynamic card update, urgent reminder
- **Groups V2** — create, info, members, list, membership check, update settings & members, dissolve
- **Calendar & schedule** — primary calendar, schedule CRUD, attendee management, attendee metadata update
- **Unified todo** — create, update, delete, query, executor management, status counts
- **Callback events** — 24 event types, structured data parsing, signature verification

## Quick Install

**SDK (library)**:
```bash
go get github.com/lansenger-pm/lansenger-sdk-go
```

**CLI (for AI agents & debugging)**:
```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
lansenger version
```

The CLI shares credentials with the SDK via `~/.lansenger/sdk_state.json`. After installing, configure credentials:
```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
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
| **Personal Bot** | Lansenger desktop → Contacts → Bots → Personal Bots → click ℹ️ icon (mobile client does NOT show credentials) |
| **Lansenger App** | Create at Lansenger Developer Center — may require organization admin approval |
| **Organization Bot** | Create at Lansenger Developer Center — may require organization admin approval |

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
result, err := client.SendText(ctx, "staff123", "Hello!", "", 0, "", false, nil, false, "", "")
result, err := client.SendMarkdown(ctx, "staff123", "**Bold**", false, nil, false, "", "")
result, err := client.SendFile(ctx, "staff123", "/path/to/report.pdf", "", 0, "", false, "", "")
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
result, err := client.SendText(ctx, "group123", "Notice", "", 0, "", false, nil, true, "", "")

// Human → group (with userToken)
result, err := client.SendGroupMessage(ctx, "group123", "text",
    map[string]interface{}{"content": "I'll handle it"}, "ut", "", false, nil, "", "", "")

// @mention in group
result, err := client.SendText(ctx, "group123", "Important!", "", 0, "", true, nil, true, "", "")
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
// Upload (core service — numeric type)
upload, err := client.UploadAppMedia(ctx, "/path/to/file.pdf", lansenger.AppMediaTypeFile, 0, 0, 0)

// Upload (app/bot — string type, supports width/height/duration)
upload, err := client.UploadAppMedia(ctx, "/path/to/video.mp4",
    lansenger.AppMediaTypeVideo, 680, 480, 300)

// Download
download, err := client.DownloadMedia(ctx, "media123")

// Download and save to file
path, err := client.DownloadMediaToFile(ctx, "media123", "/path/to/save.pdf")

// Fetch download path info
pathInfo, err := client.FetchMediaPath(ctx, "media123", "ut")

// Revoke messages
result, err := client.RevokeMessage(ctx, []string{"msg1", "msg2"}, "bot", "")

// Send urgent reminder
result, err := client.SendReminder(ctx, "msg123", []int{1, 2}, []string{"staff1", "staff2"})
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

// Dissolve group
result, err := client.DissolveGroup(ctx, "groupId", "ut")
```

## 6. Calendar & Schedule

```go
// Get primary calendar (requires userToken or userID)
cal, err := client.FetchPrimaryCalendar(ctx, "ut", "uid1")

// Create schedule (startTime/endTime are map objects, allDay is "yes"/"no")
schedule, err := client.CreateSchedule(ctx, cal.CalendarID, "Team Meeting",
    map[string]interface{}{"time": "2024-01-15T09:00"},
    map[string]interface{}{"time": "2024-01-15T10:00"},
    nil, "", "no", "", nil, "", "", "", "ut", "")

// Fetch/delete schedule
info, err := client.FetchSchedule(ctx, "cal1", "sch1", "ut", "")
result, err := client.DeleteSchedule(ctx, "cal1", "sch1", "", "", "", "ut", "")

// Update schedule
result, err := client.UpdateSchedule(ctx, "cal1", "sch1",
    map[string]interface{}{"summary": "Updated Meeting"}, "ut", "")

// Schedule list in time range
schedules, err := client.FetchScheduleList(ctx, "cal1",
    map[string]interface{}{"time": "2024-01-15T00:00"},
    map[string]interface{}{"time": "2024-01-15T23:59"}, "ut", "")

// Attendee management (attendees are []string)
attendees, err := client.FetchScheduleAttendees(ctx, "cal1", "sch1", 1, 10, "ut", "")
result, err := client.AddScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")
result, err := client.DeleteScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")

// Update attendee metadata
result, err := client.UpdateScheduleAttendeeMeta(ctx, "cal1", "sch1",
    map[string]interface{}{"rsvpStatus": "accepted"}, "ut", "")
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
// Parse plain (unencrypted) webhook payload — query string or JSON
events, err := lansenger.ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")

// Parse plain JSON callback
events, err = lansenger.ParseCallbackPayload(`{"events":[{"eventType":"staff_modify","data":{"staffId":"s001"}}],"orgId":"org1","appId":"app1"}`)

// Decrypt encrypted callback payload (AES-256-CBC)
result, err := lansenger.DecryptCallbackPayload(encryptedData, encodingKey, knownAppID)
fmt.Println(result.OrgID, result.AppID, result.Events)

// Verify signature (SHA1-based, matching Lansenger protocol)
valid := lansenger.VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, callbackToken)

// Available event types (24 types, structured field mapping)
types := lansenger.GetCallbackEventTypes()
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

### Credential Overview

All credentials are persisted per profile in `~/.lansenger/sdk_state.json` (0600 permissions):

| Credential | Required | CLI key | Description |
|-------------|----------|---------|-------------|
| App ID | ✓ | `app_id` | Lansenger app/bot ID |
| App Secret | ✓ | `app_secret` | Lansenger app/bot secret |
| API Gateway URL | ✓ | `api_gateway_url` | API gateway endpoint (default: `https://open.e.lanxin.cn/open/apigw`) |
| Passport URL | OAuth2 only | `passport_url` | OAuth2 authorize page URL |
| Redirect URI | OAuth2 only | `redirect_uri` | OAuth2 callback redirect URI (default: `http://localhost:8765`) |
| Encoding Key | Callbacks only | `encoding_key` | AES-256-CBC key for callback decryption |
| Callback Token | Callbacks only | `callback_token` | Token for callback signature verification |

### CLI Configuration

```bash
# Step 1: Set required credentials
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw

# Step 2 (optional): Set OAuth2 passport URL (needed for userToken)
lansenger config set passport_url https://passport.lx.qianxin.com
lansenger config set redirect_uri http://localhost:8765   # OAuth2 redirect URI (default)

# Step 3 (optional): Set callback credentials (needed for webhook decryption)
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN

# Verify configuration
lansenger config show

# Multi-profile support (e.g. separate orgs/apps)
lansenger config set app_id APP2_ID --profile org2
lansenger config set app_secret APP2_SECRET --profile org2
lansenger --profile org2 staff basic-info STAFF_ID
```

### SDK Configuration

**From code** (direct):
```go
client := lansenger.NewClient("app_id", "app_secret")
// Override gateway URL if needed
cfg := lansenger.NewConfig("app_id", "app_secret")
cfg.APIGatewayURL = "https://custom-gateway.example.com"
cfg.PassportURL = "https://passport.example.com"
cfg.RedirectURI = "https://myapp.com/callback"
cfg.EncodingKey = "your_encoding_key"
cfg.CallbackToken = "your_callback_token"
client := lansenger.NewClientWithConfig(cfg)
```

**From environment** (auto-detect):

| Variable | Required | Description | Default |
|----------|----------|-------------|---------|
| `LANSENGER_APP_ID` | ✓ | App/Bot ID | — |
| `LANSENGER_APP_SECRET` | ✓ | App/Bot Secret | — |
| `LANSENGER_API_GATEWAY_URL` | ✗ | API Gateway URL | `https://open.e.lanxin.cn/open/apigw` |
| `LANSENGER_PASSPORT_URL` | ✗ | Passport URL (for OAuth2) | — |
| `LANSENGER_REDIRECT_URI` | ✗ | OAuth2 redirect URI | `http://localhost:8765` |
| `LANSENGER_ENCODING_KEY` | ✗ | Encoding key for callback decryption | — |
| `LANSENGER_CALLBACK_TOKEN` | ✗ | Callback token (defaults to encoding_key) | — |
| `LANSENGER_HTTP_TIMEOUT` | ✗ | HTTP timeout (seconds) | `30` |

```go
client, err := lansenger.NewClientFromEnv()
```

### Credential & Token Persistence

By default, credentials and tokens stay in memory only (lost on process exit). Enable file persistence with `CredentialStore`:

```go
// Auto-persist to ~/.lansenger/sdk_state.json (0600 permissions)
store := lansenger.NewCredentialStore("", "default")
store.SaveCredentials("app_id", "app_secret", "https://apigw.lx.qianxin.com", "https://passport.lx.qianxin.com")
store.SaveCallbackConfig("encoding_key", "callback_token")

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

## Identity & Permissions

### Identity Capability Matrix

The Lansenger platform has three identity types with different API access:

| Command Domain | Personal Bot | Org App (Self-built) | Org App + Bot | Notes |
|--------|:---:|:---:|:---:|------|
| `message send-text/markdown/file/...` (bot DM) | **Y** | N | **Y** | Only bots can send bot DMs |
| `message send-text --group` (group chat) | N* | N | **Y** | Personal bot API supports it but no join-group feature yet |
| `message send-group-message` | N* | N | **Y** | Same as above |
| `message send-account-message` (public account) | N | **Y** | **Y** | Requires public account capability |
| `message send-user-message` (user-to-user) | N | **Y** | **Y** | Requires userToken + OAuth2 |
| `message revoke` | **Y** | **Y** | **Y** | Revoke own messages |
| `staff *` (contacts read-only) | N | **Y** | **Y** | `search` additionally requires userToken |
| `department *` | N | **Y** | **Y** | Org-level apps only |
| `calendar *` | N | **Y** | **Y** | With userToken = user identity; without = bot identity |
| `todo *` | N | **Y** | **Y** | Org-level apps only |
| `chat list/messages` | N | **Y** | **Y** | Org-level apps only |
| `group *` (group management V2) | N | N | **Y** | Requires bot to be in group |
| `media upload` | **Y** | **Y** | **Y** | General upload |
| `media upload-app` | **Y** | **Y** | **Y** | Self-built apps only (not ISV) |
| `media download/path` | **Y** | **Y** | **Y** | General download |
| `oauth *` | N | **Y** | **Y** | Org-level apps only |
| `streaming *` | N | **Y** | **Y** | Org-level apps only |
| `callback *` (event parsing) | N/A | N/A | N/A | Pure data operation, no identity required |

> \* **N\*** = API capability exists, but join-group feature not yet available.

> **Personal Bot** can only send/receive messages and upload/download files. Cannot access contacts, groups, calendars, or OAuth2.
>
> **Org App vs Org App + Bot**: Same appID/appSecret. The only difference is messaging channels — only bots can send bot DMs and group messages (because only bots can join groups). All other APIs (contacts, calendar, todo, chat, OAuth2, streaming) work identically for both. Currently only self-built apps support bot capability.

### Developer Center Permissions

Beyond identity type, specific API calls also depend on permission toggles in the Lansenger Developer Center. The organization may restrict developer access, requiring admin assistance.

**Basic Permissions (enabled by default):**

| Permission | Description |
|------|------|
| Get basic user info | Get personnel basic info for system/app login |
| Send notification messages | Get org message channels to send messages to people/groups |

**Advanced Permissions (disabled by default, must be manually enabled):**

| Permission | Description |
|------|------|
| Contacts read-only | Read access to contacts |
| Contacts edit | Edit access to contacts (create/update/delete staff) |
| Sensitive info - Phone | Access user phone numbers |
| Sensitive info - Email | Access user emails |
| Sensitive info - ID number | Access user ID numbers |
| Sensitive info - Employee ID | Access user employee IDs |
| Map unique attribute to staff ID | Map phone/email/employee ID to staff ID |
| App edit | Create and update apps |
| Groups read-only | Read access to groups |
| Groups edit | Edit access to groups |
| Calendar read-only | Read access to calendar & schedules |
| Calendar edit | Edit access to calendar & schedules |
| Upload media | Upload media file permission |
| Workbench template read | Read access to workbench templates |
| Workbench template write | Write access to workbench templates |

When encountering permission errors, first verify the identity type supports the operation, then prompt the user to enable the corresponding advanced permission in the Developer Center (contact org admin if unable to access).

## Project Structure

```
lansenger-sdk-go/
├── client.go            # LansengerClient — core client with HTTP helpers
├── config.go            # Config — configuration + env vars
├── constants.go         # API endpoints, media types, callback event types
├── errors.go            # LansengerError hierarchy (Auth/Config/API/Network/File)
├── models.go            # 50+ result/params struct types
├── auth.go              # TokenManager — appToken lifecycle with auto-refresh
├── user_token_manager.go # UserTokenManager — userToken lifecycle with auto-refresh
├── url_helpers.go       # BuildAPIURL — Options pattern for URL construction
├── oauth.go             # OAuth2 authorize URL, code exchange, token refresh
├── contacts.go          # Staff & org info APIs
├── users.go             # User profile APIs
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
├── calendars.go         # Calendar & schedule (4.23) — 10 endpoints
├── callbacks.go         # Callback event parsing + AES-256-CBC decryption + SHA1 signature verification
├── persistence.go       # CredentialStore — JSON file persistence
├── version.go           # SDK version constant
├── *_test.go            # 136 unit tests + 10 integration tests
├── cmd/lansenger/       # CLI tool (config, oauth, messaging, staff, etc.)
├── go.mod
└── README.md
```

## Development

```bash
go test . -v                        # all tests (136 unit + 10 integration)
go test . -run TestIntegration      # integration tests only (requires ~/.lansenger/sdk_state.json + network)
```

## License

MIT — see [LICENSE](LICENSE).