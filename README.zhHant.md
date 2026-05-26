[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# lansenger-sdk-go

Lansenger（藍信）平台的 Go SDK — 支援藍信應用、組織機器人和個人機器人。

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Version: 0.4.0](https://img.shields.io/badge/Version-0.4.0-blue)](https://github.com/lansenger-pm/lansenger-sdk-go)
[![Go 1.21+](https://img.shields.io/badge/Go-1.21%2B-blue)](https://go.dev/)
[![Tests: 146](https://img.shields.io/badge/Tests-146-green)](https://github.com/lansenger-pm/lansenger-sdk-go)

> 零外部依賴—僅使用 Go 標準庫。適用於任何 Go 專案。

## 支援的機器人類型

| 機器人類型 | 認證 | WebSocket 入站 | 全部 API |
|----------|------|-----------------|----------|
| **藍信應用** | appToken + userToken | ✗（使用 webhook） | ✓ |
| **組織機器人** | appToken + userToken | ✗（使用 webhook） | ✓ |
| **個人機器人** | appToken | ✓（WebSocket） | ✓（非機器人 API 受限） |

三種機器人類型使用相同的認證機制：每次 API 調用都需要 `appToken`；僅在特定使用者級別操作（使用者資訊、員工搜尋、行事曆等）時才需要 `userToken`。

## 功能特性

- **單一客戶端** — `LansengerClient` 配合 `context.Context` 用於所有 API 調用
- **憑證與令牌持久化** — `CredentialStore` 將憑證和令牌儲存到 JSON 檔案（重啟後仍有效）
- **OAuth2 使用者認證** — 授權 URL、程式碼交換、令牌刷新
- **組織與部門** — 組織資訊、部門詳情/子部門/員工
- **員工與聯絡人** — 基本/詳細資訊、ID 對映、部門祖先鏈、搜尋
- **訊息傳送** — 3 種私聊通道（機器人、官方帳號、使用者代發）+ 群聊，支援所有訊息類型、@提及、人類/機器人發送者身份
- **富卡片** — appCard（支援動態狀態更新）、oacard、linkCard、appArticles
- **串流訊息** — 基於 SSE 的即時投遞，適用於 AI 代理
- **媒體上傳/下載** — 檔案、圖片、影片，自動類型偵測，取得下載路徑，應用/機械人媒體上傳
- **訊息管理** — 撤回、動態卡片更新、加急提醒
- **群組 V2** — 建立、資訊、成員、列表、成員檢查、更新設定與成員、解散
- **行事曆與排程** — 主行事曆、排程 CRUD、出席人管理、出席人元資料更新
- **統一待辦** — 建立、更新、刪除、查詢、執行人管理、狀態計數
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
lansenger config set --app-id YOUR_APP_ID --app-secret YOUR_APP_SECRET
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

### userToken —僅特定端點需要

`userToken` 代表特定藍信使用者的授權（透過 OAuth2 取得）。僅在以下場景需要：
- 使用者級別資訊（FetchUserInfo、FetchStaffDetail、SearchStaff）
- 行事曆與排程操作（FetchPrimaryCalendar、CreateSchedule 等）
- 以人類發送者身份進行群組操作

### 取得憑證

| 機器人類型 | 如何取得 appID + appSecret |
|----------|----------------------------|
| **個人機器人** | 藍信桌面端 → 聯絡人 → 智慧機器人 → 個人機器人 → 點擊 ℹ️ 圖示（行動端不顯示憑證） |
| **藍信應用** | 在 [藍信開發者中心](https://dev.lanxin.cn) 建立 — 可能需要組織管理員審批 |
| **組織機器人** | 在 [藍信開發者中心](https://dev.lanxin.cn) 建立 — 可能需要組織管理員審批 |

### OAuth2 使用者級別認證

```go
// 建立授權 URL —將使用者重定向到藍信通行證
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

// 手機號 → staffId 對映
mapping, err := client.FetchStaffIdMapping(ctx, "orgId", "mobile", "13800138000", "")

// 員工的部門祖先鏈
ancestors, err := client.FetchDepartmentAncestors(ctx, "staffOpenId", "")

// 搜尋員工（需要 userToken 或 userID）
results, err := client.SearchStaff(ctx, "張三", "ut", "", true, nil, 1, 10)

// 組織擴充欄位 ID
fields, err := client.FetchOrgExtraFieldIDs(ctx, "orgId", "", 1, 1000)
```

## 4. 訊息與媒體

#### 機器人私聊 —最常用

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
// 機器人 → 群組
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
// 上傳（核心服務 — 數字類型）
upload, err := client.UploadMedia(ctx, "/path/to/file.pdf", lansenger.MediaTypeFile)

// 上傳（應用/機械人 — 字串類型，支援 width/height/duration）
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
| `text` | ✗ | ✓（群聊） | ✓ | 機器人、官方帳號、使用者代發 | ✓ | 最多 6000 字節 |
| `formatText` | ✓ | ✗ | ✗ |僅使用者代發 | ✓ | 透過 formatType=1 實現 Markdown |
| `oacard` | ✗ | ✗ | ✗ | 機器人、官方帳號、使用者代發 | ✓ | 帶欄位的簡單卡片 |
| `appCard` | ✓（div 標籤） | ✗ | ✗ | 機器人、官方帳號、使用者代發 | ✓ | 富卡片，支援動態更新 |
| `linkCard` | ✗ | ✗ | ✗ | 機器人、官方帳號 | ✓ | 連結預覽卡片 |
| `appArticles` | ✗ | ✗ | ✗ |僅機器人私聊 | ✓ | 文章列表（1+ 篇文章） |

**群聊**支援所有訊息類型。只有群聊支援 @提及。

## 配置

### 環境變數

| 變數 | 必填 | 描述 | 預設值 |
|------|------|------|--------|
| `LANSENGER_APP_ID` | ✓ | 應用/機器人 ID | — |
| `LANSENGER_APP_SECRET` | ✓ | 應用/機器人密鑰 | — |
| `LANSENGER_API_GATEWAY_URL` | ✗ | API 網關 URL | `https://open.e.lanxin.cn/open/apigw` |
| `LANSENGER_PASSPORT_URL` | ✗ | 通行證 URL（用於 OAuth2） | — |
| `LANSENGER_ENCODING_KEY` | ✗ | 回呼解密編碼密鑰 | — |
| `LANSENGER_CALLBACK_TOKEN` | ✗ | 回呼令牌（預設使用 encoding_key） | — |
| `LANSENGER_HTTP_TIMEOUT` | ✗ | HTTP 超時（秒） | `30` |

### 從環境變數建立

```go
client, err := lansenger.NewClientFromEnv()
```

### 憑證與令牌持久化

預設情況下，憑證和令牌僅儲存在記憶體中（程式退出後遺失）。使用 `CredentialStore` 啟用檔案持久化：

```go
// 自動持久化到 ~/.lansenger/sdk_state.json（0600 權限）
store := lansenger.NewCredentialStore("", "default")
store.SaveCredentials("app_id", "app_secret", "https://apigw.lx.qianxin.com", "https://passport.lx.qianxin.com")

// 儲存令牌
store.SaveAppToken("token123", 7200)
store.SaveUserToken("ut123", "rt123", 7200)

// 載入令牌（過期則回傳空字串）
token, err := store.LoadAppToken()

// 憑證與 Python SDK 共用（相同的 ~/.lansenger/sdk_state.json 格式）
```

啟用持久化後：
- **appToken** 可以在重啟後儲存和恢復（跳過冗餘 API 調用）
- **userToken + refreshToken** 可以在 OAuth2 交換後儲存
- **憑證 + URL** 一同儲存，實現完整配置恢復

## 專案結構

```
lansenger-sdk-go/
├── client.go            # LansengerClient — 核心客戶端與 HTTP 輔助方法
├── config.go            # Config — 配置 + 環境變數
├── constants.go         # API 端點、媒體類型、回呼事件類型
├── errors.go            # LansengerError 層級（Auth/Config/API/Network/File）
├── models.go            # 35+ 結果/參數結構體類型
├── auth.go              # TokenManager — appToken 生命週期與自動刷新
├── url_helpers.go       # BuildAPIURL — Options 模式建構 URL
├── oauth.go             # OAuth2 授權 URL、程式碼交換、令牌刷新
├── contacts.go          # 員工與組織資訊 API
├── departments.go       # 部門 API
├── groups.go            # 群組 V2 API
├── chats.go             # 聊天列表與訊息 API
├── account_messages.go  # 公共帳號通道（4.6.1）
├── user_messages.go     # 使用者代發通道（4.6.3）
├── group_messages.go    # 群聊通道（4.6.2）
├── bot_messages.go      # 機器人通道（4.6.12）
├── messaging.go         # 便捷方法 + 撤回 + 動態更新
├── streaming.go         # SSE 串流訊息
├── media.go             # 上傳/下載檔案與圖片
├── todos.go             # 結合待辦（4.33）— 12 個端點
├── calendars.go         # 行事曆與排程（4.23）— 8 個端點
├── callbacks.go         # 回呼事件解析 + AES-256-CBC 解密 + SHA1 簽章驗證
├── persistence.go       # CredentialStore — JSON 檔案持久化
├── *_test.go            # 115 單元測試 + 10 整合測試
├── go.mod
└── README.md
```

## 開發

```bash
go test ./... -v                    # 單元測試（115 個測試）
go test ./... -run TestIntegration  # 整合測試（10 個測試，需要 ~/.lansenger/sdk_state.json）
```

## 授權

MIT —詳見 [LICENSE](LICENSE)。