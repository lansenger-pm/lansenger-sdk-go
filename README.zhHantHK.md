[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# lansenger-sdk-go

Lansenger（藍信）平台的 Go SDK — 支援藍信應用、組織機械人和個人機械人。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version: 0.9.23](https://img.shields.io/badge/Version-0.9.23-blue)](https://github.com/lansenger-pm/lansenger-sdk-go)
[![Go 1.26+](https://img.shields.io/badge/Go-1.26%2B-blue)](https://go.dev/)
[![Tests: 148](https://img.shields.io/badge/Tests-148-green)](https://github.com/lansenger-pm/lansenger-sdk-go)

> SDK 函式庫零外部依賴（僅標準函式庫）。CLI（`cmd/lansenger`）使用 [spf13/cobra](https://github.com/spf13/cobra)。

## 支援的機械人類型

| 機械人類型 | 認證 | WebSocket 入站 | 全部 API |
|----------|------|-----------------|----------|
| **藍信應用** | appToken + userToken | ✗（使用 webhook） | ✓ |
| **組織機械人** | appToken + userToken | ✗（使用 webhook） | ✓ |
| **個人機械人** | appToken | ✓（WebSocket） | ✓（非機械人 API 受限） |

三種機械人類型使用相同的認證機制：每次 API 調用都需要 `appToken`；僅在特定使用者級別操作（使用者資訊、員工搜尋、行事曆等）時才需要 `userToken`。

## 功能特性

- **單一客戶端** — `LansengerClient` 配合 `context.Context` 用於所有 API 調用
- **憑證與令牌持久化** — `CredentialStore` 將憑證和令牌儲存到 JSON 檔案（重啟後仍有效）
- **OAuth2 使用者認證** — 授權 URL、程式碼交換、令牌刷新
- **組織與部門** — 組織資訊、部門詳情/子部門/員工
- **員工與聯絡人** — 基本/詳細資訊、ID 對映、部門祖先鏈、搜尋
- **訊息發送** — 3 種私聊通道（機械人、官方帳號、使用者代發）+ 群聊，支援所有訊息類型、@提及、人類/機械人發送者身份
- **富卡片** — appCard（支援動態狀態更新）、oacard、linkCard、appArticles
- **串流訊息** — 基於 SSE 的即時投遞，適用於 AI 代理
- **媒體上載/下載** — 檔案、圖片、影片，自動類型偵測，取得下載路徑，應用/機械人媒體上載
- **訊息管理** — 撤回、動態卡片更新、加急提醒
- **群組 V2** — 建立、資訊、成員、列表、成員檢查、更新設定與成員、解散
- **行事曆與排程** — 主行事曆、排程 CRUD、出席人管理、出席人元資料更新、update_schedule_attendees()
- **統一待辦** — 建立、更新、刪除、查詢、執行人管理、狀態計數
- **機械人命令** — 創建/查詢/刪除機械人快捷命令
- **個人應用** — 創建/修改/查詢/刪除/列表個人機械人應用
- **回呼事件** — 24 種事件類型、結構化資料解析、簽章驗證

## 快速安裝

**SDK（函式庫）**:
```bash
go get github.com/lansenger-pm/lansenger-sdk-go
```

**CLI（用於 AI 代理與除錯）**:
```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
lansenger version
```

CLI 透過 `~/.lansenger/sdk_state.json` 與 SDK 共用憑證。安裝後設定憑證：
```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
```

## 1. 認證

### appToken — 所有 API 調用必需

每個 SDK 方法都需要 `appToken`。客戶端使用你的 `appID` + `appSecret` 自動取得並刷新令牌。你無需手動管理 appToken — `TokenManager` 處理整個生命週期：

1. **首次調用** → `GET /v1/apptoken/create` 使用 appID + appSecret → 回傳 `appToken`（有效期 2 小時）
2. **後續調用** → 複用已緩存的 appToken 直到過期
3. **令牌過期** → 透過相同端點自動刷新

```go
// appToken 自動管理 — 只需設定 appID + appSecret
client := lansenger.NewClient("your-appid", "your-secret")

// 也可以手動取得/失效令牌
token, err := client.GetToken(ctx)
client.InvalidateToken() // 強制下次調用時刷新
```

### userToken — 僅特定端點需要

`userToken` 代表特定藍信使用者的授權（透過 OAuth2 取得）。僅在以下場景需要：
- 使用者級別資訊（FetchUserInfo、FetchStaffDetail、SearchStaff）
- 行事曆與排程操作（FetchPrimaryCalendar、CreateSchedule 等）
- 以人類發送者身份進行群組操作

### 取得憑證

| 機械人類型 | 如何取得 appID + appSecret |
|----------|----------------------------|
| **個人機械人** | 藍信桌面端 → 聯絡人 → 智能機械人 → 個人機械人 → 點擊 ℹ️ 圖示（流動端不顯示憑證） |
| **藍信應用** | 在 藍信開發者中心 建立 — 可能需要組織管理員審批 |
| **組織機械人** | 在 藍信開發者中心 建立 — 可能需要組織管理員審批 |

### OAuth2 使用者級別認證

```go
// 建立授權 URL — 將使用者重定向到藍信通行證
url := client.BuildAuthorizeURL("https://myapp.com/callback", "", "state123")

// 使用者授權後，用程式碼交換 userToken + refreshToken
tokenResult, err := client.ExchangeCode(ctx, "auth_code_from_callback", "https://myapp.com/callback")

// 刷新過期的 userToken
newToken, err := client.RefreshUserToken(ctx, tokenResult.RefreshToken, "")

// 取得使用者資訊
userInfo, err := client.FetchUserInfo(ctx, tokenResult.UserToken)
```

## 2. 組織與部門

```go
// 組織資訊
org, err := client.FetchOrgInfo(ctx, "orgId", "")

// 部門層級
detail, err := client.FetchDepartmentDetail(ctx, "deptId", "", "")
children, err := client.FetchDepartmentChildren(ctx, "deptId", "")
staffs, err := client.FetchDepartmentStaffs(ctx, "deptId", "", 1, 100)
```

## 3. 員工與聯絡人

```go
// 基本員工資訊
staff, err := client.FetchStaffBasicInfo(ctx, "staffOpenId", "")

// 詳細資料（建議使用 userToken）
detail, err := client.FetchStaffDetail(ctx, "staffOpenId", "ut")

// 流動電話 → staffId 對映
mapping, err := client.FetchStaffIdMapping(ctx, "orgId", "mobile", "13800138000", "")

// 員工的部門祖先鏈
ancestors, err := client.FetchDepartmentAncestors(ctx, "staffOpenId", "")

// 搜尋員工（需要 userToken 或 userID）
results, err := client.SearchStaff(ctx, "張三", "ut", "", true, nil, 1, 10)

// 組織擴充欄位 ID
fields, err := client.FetchOrgExtraFieldIDs(ctx, "orgId", "", 1, 1000)
```

## 4. 訊息與媒體

#### 機械人私聊 — 最常用

```go
result, err := client.SendText(ctx, "staff123", "你好！", "", 0, "", false, nil, false, "", "")
result, err := client.SendMarkdown(ctx, "staff123", "**加粗**", false, nil, false, "", "")
result, err := client.SendFile(ctx, "staff123", "/path/to/report.pdf", "", 0, "", false, "", "")
```

#### 公共帳號通道

```go
result, err := client.SendAccountMessage(ctx, "text",
    map[string]interface{}{"content": "系統通知"},
    []string{"staff1", "staff2"}, nil, "524288-xxxx", "", "", "")
```

#### 使用者代發通道（需要 userToken）

```go
result, err := client.SendUserMessage(ctx, "staff456", "text",
    map[string]interface{}{"content": "你好"}, "ut", "")
```

#### 群聊

```go
// 機械人 → 群組
result, err := client.SendText(ctx, "group123", "通知", "", 0, "", false, nil, true, "", "")

// 人類 → 群組（使用 userToken）
result, err := client.SendGroupMessage(ctx, "group123", "text",
    map[string]interface{}{"content": "我來處理"}, "ut", "", false, nil, "", "", "")

// 群聊 @提及
result, err := client.SendText(ctx, "group123", "重要！", "", 0, "", true, nil, true, "", "")
```

#### 富卡片

```go
// appCard
params := &lansenger.AppCardParams{
    ChatID: "staff123", BodyTitle: "審批", IsDynamic: true,
}
result, err := client.SendAppCardWithParams(ctx, params)

// linkCard
params := &lansenger.LinkCardParams{
    ChatID: "staff123", Title: "文章", Link: "https://...",
}
result, err := client.SendLinkCardWithParams(ctx, params)

// 更新動態卡片狀態
updateParams := &lansenger.DynamicCardUpdateParams{
    MsgID: "msg123", IsLastUpdate: true,
}
result, err := client.UpdateDynamicCard(ctx, updateParams)
```

#### 串流訊息（用於 AI 代理）

```go
result, err := client.CreateStreamMessage(ctx, "staff1", "staff", "stream1")
result, err := client.FetchStreamMessage(ctx, "msg123")
```

#### 媒體

```go
// 上載（核心服務 — 數字類型）
upload, err := client.UploadAppMedia(ctx, "/path/to/file.pdf", lansenger.AppMediaTypeFile, 0, 0, 0)

// 上載（應用/機械人 — 字串類型，支援 width/height/duration）
upload, err := client.UploadAppMedia(ctx, "/path/to/video.mp4",
    lansenger.AppMediaTypeVideo, 680, 480, 300)

// 下載
download, err := client.DownloadMedia(ctx, "media123")

// 下載並儲存到檔案
path, err := client.DownloadMediaToFile(ctx, "media123", "/path/to/save.pdf")

// 取得下載路徑資訊
pathInfo, err := client.FetchMediaPath(ctx, "media123", "ut")

// 撤回訊息
result, err := client.RevokeMessage(ctx, []string{"msg1", "msg2"}, "bot", "")

// 發送加急提醒
result, err := client.SendReminder(ctx, "msg123", []int{1, 2}, []string{"staff1", "staff2"})
```

## 5. 群組

```go
// 建立群組
info := &lansenger.GroupCreateInfo{
    Name: "專案聊天", OrgID: 1, StaffIDList: []string{"s1", "s2", "s3"},
}
group, err := client.CreateGroup(ctx, info, "")

// 取得資訊與成員
info, err := client.FetchGroupInfo(ctx, "groupOpenId", "")
members, err := client.FetchGroupMembers(ctx, "groupOpenId", "", 0, 100)
groups, err := client.FetchGroupList(ctx, "", 0, 100)

// 檢查成員關係
result, err := client.CheckIsInGroup(ctx, "groupOpenId", "", "staff1")

// 更新設定
result, err := client.UpdateGroupInfo(ctx, "groupId", map[string]interface{}{"name": "新名稱"}, "")

// 新增/移除成員
result, err := client.UpdateGroupMembers(ctx, "groupId",
    []string{"staff4"}, []string{"staff3"}, nil, "")

// 解散群組
result, err := client.DissolveGroup(ctx, "groupId", "ut")
```

## 6. 行事曆與排程

```go
// 取得主行事曆（需要 userToken 或 userID）
cal, err := client.FetchPrimaryCalendar(ctx, "ut", "uid1")

// 建立排程（startTime/endTime 為 map 物件，allDay 為 "yes"/"no")
schedule, err := client.CreateSchedule(ctx, cal.CalendarID, "團隊會議",
    map[string]interface{}{"time": "2024-01-15T09:00"},
    map[string]interface{}{"time": "2024-01-15T10:00"},
    nil, "", "no", "", nil, "", "", "", "ut", "")

// 取得/刪除排程
info, err := client.FetchSchedule(ctx, "cal1", "sch1", "ut", "")
result, err := client.DeleteSchedule(ctx, "cal1", "sch1", "", "", "", "ut", "")

// 更新排程
result, err := client.UpdateSchedule(ctx, "cal1", "sch1",
    map[string]interface{}{"summary": "更新後的會議"}, "ut", "")

// 時間範圍內的排程列表
schedules, err := client.FetchScheduleList(ctx, "cal1",
    map[string]interface{}{"time": "2024-01-15T00:00"},
    map[string]interface{}{"time": "2024-01-15T23:59"}, "ut", "")

// 出席人管理（出席人為 []string）
attendees, err := client.FetchScheduleAttendees(ctx, "cal1", "sch1", 1, 10, "ut", "")
result, err := client.AddScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")
result, err := client.DeleteScheduleAttendees(ctx, "cal1", "sch1",
    []string{"staff2"}, "", "", "", "ut", "")

// 更新出席人元資料
result, err := client.UpdateScheduleAttendeeMeta(ctx, "cal1", "sch1",
    map[string]interface{}{"rsvpStatus": "accepted"}, "ut", "")
```

## 7. 統一待辦

```go
// 建立待辦任務
todo, err := client.CreateTodoTask(ctx, "審批請求", lansenger.TodoTypeApproval,
    "https://app.com/a/1", "https://pc.app.com/a/1", []string{"staff1"}, "org1", "", "", "", "")

// 更新狀態（11=待讀, 12=已讀, 21=待辦, 22=已完成）
result, err := client.UpdateTodoTaskStatus(ctx, "taskId", lansenger.TodoStatusDone, "org1", "", "")

// 更新內容
result, err := client.UpdateTodoTask(ctx, "taskId", "已更新", "l", "p", "org1", "", "")

// 刪除（僅發送者）
result, err := client.DeleteTodoTask(ctx, "taskId", "org1", "", "")

// 查詢
list, err := client.FetchTodoTaskList(ctx, "org1", nil, "", nil, "")
task, err := client.FetchTodoTaskByID(ctx, "taskId", "org1", "", "")
task, err := client.FetchTodoTaskBySourceID(ctx, "src1", "org1", "", "")
counts, err := client.FetchTodoTaskStatusCounts(ctx, "staff1", "org1", "", "", "")

// 執行人管理
result, err := client.AddExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
result, err := client.DeleteExecutors(ctx, []string{"staff2"}, "org1", "taskId", "")
executors, err := client.FetchExecutorList(ctx, "taskId", "org1", "", nil, "")
```

## 8. 回呼事件

```go
// 解析明文（未加密）webhook 資料 — URL 查詢字串或 JSON
events, err := lansenger.ParseCallbackPayload("eventType=staff_modify&staffId=s001&orgId=org1")

// 解析明文 JSON 回呼
events, err = lansenger.ParseCallbackPayload(`{"events":[{"eventType":"staff_modify","data":{"staffId":"s001"}}],"orgId":"org1","appId":"app1"}`)

// 解密加密回呼資料（AES-256-CBC）
result, err := lansenger.DecryptCallbackPayload(encryptedData, encodingKey, knownAppID)
fmt.Println(result.OrgID, result.AppID, result.Events)

// 驗證簽章（基於 SHA1，匹配藍信協議）
valid := lansenger.VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, callbackToken)

// 可用事件類型（24 種，結構化欄位對映）
types := lansenger.GetCallbackEventTypes()
```

## 9. 聊天閱讀

```go
// 取得使用者聊天列表（私聊 + 群聊）
chats, err := client.FetchChatList(ctx, "ut", "private", "", "", "")

// 取得與特定人的私聊訊息
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "s001", "", "", "", "")

// 取得群聊訊息
msgs, err := client.FetchChatMessages(ctx, "ut", 10, "", "", "g001", "", "", "")
```

## 訊息類型能力矩陣

| msgType | Markdown | @提及 | 附件 | 私聊通道 | 群聊 | 備註 |
|---------|----------|--------|------|----------|------|------|
| `text` | ✗ | ✓（群聊） | ✓ | 機械人、官方帳號、使用者代發 | ✓ | 最多 6000 字節 |
| `formatText` | ✓ | ✗ | ✗ | 僅使用者代發 | ✓ | 透過 formatType=1 實現 Markdown |
| `oacard` | ✗ | ✗ | ✗ | 機械人、官方帳號、使用者代發 | ✓ | 帶欄位的簡單卡片 |
| `appCard` | ✓（div 標籤） | ✗ | ✗ | 機械人、官方帳號、使用者代發 | ✓ | 富卡片，支援動態更新 |
| `linkCard` | ✗ | ✗ | ✗ | 機械人、官方帳號 | ✓ | 連結預覽卡片 |
| `appArticles` | ✗ | ✗ | ✗ | 唯機械人私聊 | ✓ | 文章列表（1+ 篇文章） |

**群聊**支援所有訊息類型。只有群聊支援 @提及。

## 設定

### 憑證概覽

所有憑證按 profile 持久化儲存於 `~/.lansenger/sdk_state.json`（0600 權限）：

| 憑證 | 必填 | CLI 鍵名 | 說明 |
|------|------|----------|------|
| App ID | ✓ | `app_id` | 藍信應用/機械人 ID |
| App Secret | ✓ | `app_secret` | 藍信應用/機械人密鑰 |
| API Gateway URL | ✓ | `api_gateway_url` | API 閘道地址 |
| Passport URL | 僅 OAuth2 | `passport_url` | OAuth2 授權頁地址 |
| Redirect URI | 僅 OAuth2 | `redirect_uri` | OAuth2 回呼地址（預設：`http://localhost:8765`） |
| Encoding Key | 僅回呼 | `encoding_key` | AES-256-CBC 解密密鑰 |
| Callback Token | 僅回呼 | `callback_token` | 回呼簽章驗證令牌 |

### CLI 設定

```bash
# 第1步：設定必填憑證
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://your-gateway.example.com

# 第2步（可選）：設定 OAuth2 授權頁地址（取得 userToken 需要）
lansenger config set passport_url https://your-passport.example.com
lansenger config set redirect_uri http://localhost:8765   # OAuth2 回呼地址（預設值）

# 第3步（可選）：設定回呼憑證（接收 Webhook 回呼需要）
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN

# 驗證設定
lansenger config show

# 多 profile 支援（如不同組織/應用）
lansenger config set app_id APP2_ID --profile org2
lansenger config set app_secret APP2_SECRET --profile org2
lansenger --profile org2 staff basic-info STAFF_ID
```

### SDK 設定

**程式碼方式**（直接傳入）：
```go
client := lansenger.NewClient("app_id", "app_secret")
// 如需自訂網關地址
cfg := lansenger.NewConfig("app_id", "app_secret")
cfg.APIGatewayURL = "https://custom-gateway.example.com"
cfg.PassportURL = "https://passport.example.com"
cfg.RedirectURI = "https://myapp.com/callback"
cfg.EncodingKey = "your_encoding_key"
cfg.CallbackToken = "your_callback_token"
client := lansenger.NewClientWithConfig(cfg)
```

**環境變數方式**（自動偵測）：

| 變數 | 必填 | 說明 | 預設值 |
|------|------|------|--------|
| `LANSENGER_APP_ID` | ✓ | 應用/機械人 ID | — |
| `LANSENGER_APP_SECRET` | ✓ | 應用/機械人密鑰 | — |
| `LANSENGER_API_GATEWAY_URL` | ✓ | API 閘道地址 | — |
| `LANSENGER_PASSPORT_URL` | ✗ | 授權頁地址（OAuth2） | — |
| `LANSENGER_REDIRECT_URI` | ✗ | OAuth2 回呼地址 | `http://localhost:8765` |
| `LANSENGER_ENCODING_KEY` | ✗ | 回呼解密密鑰 | — |
| `LANSENGER_CALLBACK_TOKEN` | ✗ | 回呼驗證令牌（預設同 encoding_key） | — |
| `LANSENGER_HTTP_TIMEOUT` | ✗ | HTTP 逾時（秒） | `30` |

```go
client, err := lansenger.NewClientFromEnv()
```

### 憑證與令牌持久化

預設情況下，憑證和令牌僅在記憶體中保留（程序結束後消失）。使用 `CredentialStore` 啟用檔案持久化：

```go
// 自動持久化到 ~/.lansenger/sdk_state.json（0600 權限）
store := lansenger.NewCredentialStore("", "default")
store.SaveCredentials("app_id", "app_secret", "https://your-gateway.example.com", "https://your-passport.example.com")
store.SaveCallbackConfig("encoding_key", "callback_token")

// 儲存令牌
store.SaveAppToken("token123", 7200)
store.SaveUserToken("ut123", "rt123", 7200)

// 載入令牌（過期時回傳空字串）
token, err := store.LoadAppToken()

// 憑證與 Python SDK 共用（相同 ~/.lansenger/sdk_state.json 格式）
```

啟用持久化後：
- **appToken** 可在重啟後儲存與恢復（跳過冗餘 API 請求）
- **userToken + refreshToken** 可在 OAuth2 交換後儲存
- **憑證 + URL** 一併儲存，完整恢復設定

## 身份與權限

### 身份能力矩陣

藍信平台有三種身份類型，各自擁有不同的 API 存取權限：

| 命令域 | 個人機器人 | 組織應用（自建） | 組織應用 + 機器人 | 備註 |
|--------|:---:|:---:|:---:|------|
| `message send-text/markdown/file/...` (bot DM) | **Y** | N | **Y** | 僅機器人可傳送機器人私聊訊息 |
| `message send-text --group` (群聊) | **Y** | N | **Y** | 個人機械人現已支援群聊 |
| `message send-group-message` | **Y** | N | **Y** | 同上 |
| `message send-account-message` (公共號) | N | **Y** | **Y** | 需要公共號能力 |
| `message send-user-message` (使用者代發) | N | **Y** | **Y** | 需要 userToken + OAuth2 |
| `message revoke` | **Y** | **Y** | **Y** | 撤回自己的訊息 |
| `staff *` (通訊錄唯讀) | N | **Y** | **Y** | `search` 額外需要 userToken |
| `department *` | N | **Y** | **Y** | 僅組織級應用 |
| `calendar *` | N | **Y** | **Y** | 帶 userToken = 使用者身份；不帶 = 機器人身份 |
| `todo *` | N | **Y** | **Y** | 僅組織級應用 |
| `chat list/messages` | N | **Y** | **Y** | 僅組織級應用 |
| `group *` (群組管理 V2) | N | N | **Y** | 需要機器人已在群內 |
| `media upload` | **Y** | **Y** | **Y** | 通用上傳 |
| `media upload-app` | **Y** | **Y** | **Y** | 僅自建應用（非 ISV） |
| `media download/path` | **Y** | **Y** | **Y** | 通用下載 |
| `oauth *` | N | **Y** | **Y** | 僅組織級應用 |
| `streaming *` | N | **Y** | **Y** | 僅組織級應用 |
| `callback *` (事件解析) | N/A | N/A | N/A | 純資料操作，與身份無關 |


> **個人機器人** 只能收發訊息和上傳/下載檔案，無法存取通訊錄、群組、行事曆或 OAuth2。
>
> **組織應用 vs 組織應用 + 機器人**：使用相同的 appID/appSecret。唯一區別在於訊息通道 —— 僅機器人可以傳送機器人私聊訊息和群聊訊息（因為只有機器人能加入群聊）。其他所有 API（通訊錄、行事曆、待辦、聊天記錄、OAuth2、串流訊息）兩者完全一致。目前僅自建應用支援機器人能力。

### 開發者中心權限

除了身份類型，特定 API 調用還取決於藍信開發者中心中的權限開關。組織可能限制開發者存取權限，需要管理員協助。

**基礎權限（預設開啟）：**

| 權限 | 說明 |
|------|------|
| Get basic user info | 取得人員基本資訊，用於系統/應用登入 |
| Send notification messages | 取得組織訊息通道，向人員/群組傳送訊息 |

**進階權限（預設關閉，需手動開啟）：**

| 權限 | 說明 |
|------|------|
| Contacts read-only | 通訊錄唯讀存取 |
| Contacts edit | 通訊錄編輯存取（建立/更新/刪除人員） |
| Sensitive info - Phone | 存取使用者手機號碼 |
| Sensitive info - Email | 存取使用者電子郵件 |
| Sensitive info - ID number | 存取使用者身份證號 |
| Sensitive info - Employee ID | 存取使用者工號 |
| Map unique attribute to staff ID | 將手機/電子郵件/工號對映為人員 ID |
| App edit | 建立及更新應用 |
| Groups read-only | 群組唯讀存取 |
| Groups edit | 群組編輯存取 |
| Calendar read-only | 行事曆及排程唯讀存取 |
| Calendar edit | 行事曆及排程編輯存取 |
| Upload media | 媒體檔案上傳權限 |
| Workbench template read | 工作台範本讀取 |
| Workbench template write | 工作台範本寫入 |

遇到權限錯誤時，請先確認身份類型是否支援該操作，然後提示使用者在開發者中心開啟對應的進階權限（如無法存取請聯絡組織管理員）。

## 專案結構

```
lansenger-sdk-go/
├── client.go            # LansengerClient — 核心客戶端與 HTTP 輔助方法
├── config.go            # Config — 配置 + 環境變數
├── constants.go         # API 端點、媒體類型、回呼事件類型
├── errors.go            # LansengerError 層級（Auth/Config/API/Network/File）
├── models.go            # 50+ 結果/參數結構體類型
├── auth.go              # TokenManager — appToken 生命週期與自動刷新
├── user_token_manager.go # UserTokenManager — userToken 生命週期與自動刷新
├── url_helpers.go       # BuildAPIURL — Options 模式建構 URL
├── oauth.go             # OAuth2 授權 URL、程式碼交換、令牌刷新
├── contacts.go          # 員工與組織資訊 API
├── users.go             # 使用者資料 API
├── departments.go       # 部門 API
├── groups.go            # 群組 V2 API
├── chats.go             # 聊天列表與訊息 API
├── account_messages.go  # 公共帳號通道（4.6.1）
├── user_messages.go     # 使用者代發通道（4.6.3）
├── group_messages.go    # 群聊通道（4.6.2）
├── bot_messages.go      # 機械人通道（4.6.12）
├── messaging.go         # 便捷方法 + 撤回 + 動態更新
├── streaming.go         # SSE 串流訊息
├── media.go             # 上載/下載檔案與圖片
├── todos.go             # 結合待辦（4.33）— 12 個端點
├── calendars.go         # 行事曆與排程（4.23）— 10 個端點
├── callbacks.go         # 回呼事件解析 + AES-256-CBC 解密 + SHA1 簽章驗證
├── persistence.go       # CredentialStore — JSON 檔案持久化
├── version.go           # SDK 版本常數
├── *_test.go            # 136 單元測試 + 10 整合測試
├── cmd/lansenger/       # CLI 工具（config、oauth、訊息、員工等）
├── go.mod
└── README.md
```

## 開發

```bash
go test . -v                        # 全部測試（136 單元 + 10 整合）
go test . -run TestIntegration      # 僅整合測試（需要 ~/.lansenger/sdk_state.json + 網路）
```

## 授權

MIT — 詳見 [LICENSE](LICENSE)。