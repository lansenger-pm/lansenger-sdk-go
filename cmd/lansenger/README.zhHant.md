[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# Lansenger CLI（Go）

Lansenger 命令行工具 — 在終端直接調用藍信開放平台 API，傳送訊息、管理群組、查詢人員/部門、操作行事曆與待辦等。

命令語法與 Python 版、TypeScript 版完全一致，安裝任一版本即可使用。

## 安裝

```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
```

或從原始碼安裝：

```bash
git clone https://github.com/lansenger-pm/lansenger-sdk-go.git
cd lansenger-sdk-go/cmd/lansenger
go build -o lansenger .
```

需要 Go 1.26+。

## 快速開始

### 1. 設定憑證

透過 `config set` 命令儲存憑證（按 profile 隔離儲存在 `~/.lansenger/sdk_state.json`，密鑰脫敏顯示，檔案權限 0600）：

**基本憑證（所有使用者必填）**：

```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw
```

**OAuth2 使用者認證（需要取得 userToken 時填寫）**：

```bash
lansenger config set passport_url https://passport.lx.qianxin.com
lansenger config set redirect_uri http://localhost:8765   # OAuth2 回呼地址（預設值）
```

**回呼接收（需要解析/驗簽回呼 Webhook 時填寫）**：

```bash
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN
```

也可以透過環境變數設定（適合 CI/CD 或臨時使用）：

```bash
export LANSENGER_APP_ID=YOUR_APP_ID
export LANSENGER_APP_SECRET=YOUR_APP_SECRET
export LANSENGER_REDIRECT_URI=http://localhost:8765
```

### 2. 檢視設定

```bash
lansenger config show
```

### 3. 傳送第一條訊息

```bash
lansenger message send-text staff001 "Hello from CLI!"
```

## 命令總覽

| 命令組 | 說明 | 子命令 |
|--------|------|--------|
| `config` | 管理憑證設定 | `set`, `show`, `clear`, `list-profiles` |
| `message` | 傳送與管理訊息 | `send-text`, `send-markdown`, `send-file`, `send-image-url`, `send-link-card`, `send-app-articles`, `send-app-card`, `send-oacard`, `send-bot-message`, `send-group-message`, `send-account-message`, `send-user-message`, `update-dynamic-card`, `revoke`, `query-groups`, `send-reminder` |
| `group` | 管理群組 | `create`, `info`, `members`, `list`, `check`, `update`, `update-members`, `dismiss` |
| `staff` | 查詢人員資訊 | `basic-info`, `detail`, `ancestors`, `id-mapping`, `org-extra-fields`, `search`, `org-info` |
| `department` | 查詢部門資訊 | `detail`, `children`, `staffs` |
| `calendar` | 行事曆操作 | `primary`, `create-schedule`, `fetch-schedule`, `delete-schedule`, `list-schedules`, `attendees`, `add-attendees`, `delete-attendees`, `update-schedule`, `attendee-meta` |
| `todo` | 待辦任務管理 | `create`, `update`, `update-status`, `delete`, `list`, `fetch-by-id`, `fetch-by-source`, `status-counts`, `executor-status`, `add-executors`, `delete-executors`, `executor-list` |
| `oauth` | OAuth2 使用者認證 | `authorize-url`, `exchange-code`, `refresh-token`, `user-info`, `parse-callback`, `validate-state` |
| `callback` | 回呼事件解析 | `parse-payload`, `decrypt-payload`, `verify-signature`, `event-types` |
| `media` | 媒體檔案操作 | `upload`, `upload-app`, `download`, `download-to-file`, `path` |
| `streaming` | 串流訊息（AI 場景） | `create`, `fetch` |
| `chat` | 會話與訊息記錄 | `list`, `messages` |
| `health` | 連線健康檢查 | `check` |

## 常用範例

### 訊息傳送

```bash
# 傳送純文字訊息
lansenger message send-text chat123 "你好！"

# 傳送 Markdown 訊息
lansenger message send-markdown chat123 "**粗體** *斜體*"

# 傳送檔案
lansenger message send-file chat123 /path/to/report.pdf

# 傳送網路圖片
lansenger message send-image-url chat123 https://example.com/photo.jpg

# 傳送連結卡片
lansenger message send-link-card chat123 "文件" "點擊檢視" https://docs.example.com

# 傳送應用卡片
lansenger message send-app-card chat123 "卡片標題" --content "正文內容" --card-link https://example.com

# 傳送多條圖文（appArticles）
lansenger message send-app-articles chat123 '{"title":"文章1","url":"https://a.com"}' '{"title":"文章2","url":"https://b.com"}'

# 群內傳送並 @all（user_token 可選，無則顯示為 bot）
lansenger message send-text group123 "全員通知" --group --mention-all

# 群內 @指定人
lansenger message send-text group123 "請檢視" --group --mention staff001

# 機器人通道傳送訊息
lansenger message send-bot-message text '{"content":"通知內容"}' --chat-id user001 --chat-id user002
```

### 群組管理

```bash
# 建立群組
lansenger group create "專案群" org001 --staff staff001 --staff staff002

# 檢視群組資訊
lansenger group info group123

# 檢視群組成員
lansenger group members group123

# 檢視群組列表（bot 可檢視所在的群組，傳 user_token 可檢視使用者所在的群組）
lansenger group list

# 檢視使用者所在的群組列表（需要 user_token）
lansenger group list --user-token YOUR_USER_TOKEN

# 檢查使用者是否在群組內
lansenger group check group123 --staff-id staff001

# 更新群組資訊
lansenger group update group123 --name "新名稱" --desc "新描述"

# 新增/移除成員
lansenger group update-members group123 --add staff003 --remove staff001
```

### 人員查詢

```bash
# 檢視人員基本資訊
lansenger staff basic-info staff001

# 檢視人員詳細資訊
lansenger staff detail staff001

# 搜尋人員
lansenger staff search "張三" --user-token YOUR_USER_TOKEN

# ID 對映（手機號 → staffId）
lansenger staff id-mapping org001 mobile 13800138000

# 檢視部門祖先鏈
lansenger staff ancestors staff001
```

### 媒體檔案

```bash
# 上傳核心平台檔案
lansenger media upload /path/to/file.pdf --media-type 3

# 上傳應用/機器人媒體檔案（用於 send-text / send-file 等）
lansenger media upload-app /path/to/file.pdf --media-type file

# 下載媒體檔案到本地
lansenger media download-to-file MEDIA_ID --output /path/to/save.pdf
```

## 全域選項

| 選項 | 說明 |
|------|------|
| `--json` / `-j` | 輸出原始 JSON 格式而非表格 |
| `--profile` / `-P` | 使用指定的憑證 profile（預設：`default`） |

## 多應用/多機器人設定（Profile）

CLI 支援多 profile，每個 profile 對應一個 appID（一個應用或一個機器人），憑證互相隔離：

```bash
# 設定第一個應用（個人機器人）
lansenger config set app_id xxx1 --profile my-bot
lansenger config set app_secret xxx1 --profile my-bot

# 設定第二個應用（組織機器人）
lansenger config set app_id xxx2 --profile org-bot
lansenger config set app_secret xxx2 --profile org-bot

# 使用指定 profile
lansenger --profile org-bot staff basic-info STAFF_ID
```

## 安全性

- 憑證儲存在 `~/.lansenger/sdk_state.json`，檔案權限 `0600`
- `config show` 對所有密鑰類欄位脫敏顯示（`***`），僅 `api_gateway_url` 和 `passport_url` 明文展示
- 支援環境變數 `LANSENGER_APP_ID` / `LANSENGER_APP_SECRET` / `LANSENGER_ENCODING_KEY` / `LANSENGER_CALLBACK_TOKEN`，適合 CI/CD 場景

## CLI 相容性

本 Go CLI 與 Python 版、TypeScript 版命令語法完全一致：

```bash
# Python CLI
pip install lansenger-cli

# Go CLI
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest

# TypeScript CLI
npm install -g lansenger-cli
```

## 授權條款

MIT License
