[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# Lansenger CLI（Go）

Lansenger 命令行工具 — 在终端直接调用蓝信开放平台 API，发送消息、管理群组、查询人员/部门、操作日程与待办等。

命令语法与 Python 版、TypeScript 版完全一致，安装任一版本即可使用。

## 安装

```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
```

或从源码安装：

```bash
git clone https://github.com/lansenger-pm/lansenger-sdk-go.git
cd lansenger-sdk-go/cmd/lansenger
go build -o lansenger .
```

需要 Go 1.26+。

## 快速开始

### 1. 配置凭证

通过 `config set` 命令保存凭证（按 profile 隔离存储在 `~/.lansenger/sdk_state.json`，密钥脱敏显示，文件权限 0600）：

**基本凭证（所有用户必填）**：

```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw
```

**OAuth2 用户认证（需要获取 userToken 时填写）**：

```bash
lansenger config set passport_url https://passport.lx.qianxin.com
lansenger config set redirect_uri http://localhost:8765   # OAuth2 回调地址（默认值）
```

**回调接收（需要解析/验签回调 Webhook 时填写）**：

```bash
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN
```

也可以通过环境变量配置（适合 CI/CD 或临时使用）：

```bash
export LANSENGER_APP_ID=YOUR_APP_ID
export LANSENGER_APP_SECRET=YOUR_APP_SECRET
export LANSENGER_REDIRECT_URI=http://localhost:8765
```

### 2. 查看配置

```bash
lansenger config show
```

### 3. 发送第一条消息

```bash
lansenger message send-text staff001 "Hello from CLI!"
```

## 命令总览

| 命令组 | 说明 | 子命令 |
|--------|------|--------|
| `config` | 管理凭证配置 | `set`, `show`, `clear`, `list-profiles`, `delete-profile`, `list-users` |
| `message` | 发送与管理消息 | `send-text`, `send-markdown`, `send-file`, `send-image-url`, `send-link-card`, `send-app-articles`, `send-app-card`, `send-oacard`, `send-bot-message`, `send-group-message`, `send-account-message`, `send-user-message`, `update-dynamic-card`, `revoke`, `query-groups`, `send-reminder` |
| `group` | 管理群组 | `create`, `info`, `members`, `list`, `check`, `update`, `update-members`, `dismiss` |
| `staff` | 查询人员信息 | `basic-info`, `detail`, `ancestors`, `id-mapping`, `org-extra-fields`, `search`, `org-info` |
| `department` | 查询部门信息 | `detail`, `children`, `staffs` |
| `calendar` | 日程操作 | `primary`, `create-schedule`, `fetch-schedule`, `delete-schedule`, `list-schedules`, `attendees`, `add-attendees`, `delete-attendees`, `update-schedule`, `attendee-meta` |
| `todo` | 待办任务管理 | `create`, `update`, `update-status`, `delete`, `list`, `fetch-by-id`, `fetch-by-source`, `status-counts`, `executor-status`, `add-executors`, `delete-executors`, `executor-list` |
| `oauth` | OAuth2 用户认证 | `authorize-url`, `exchange-code`, `refresh-token`, `user-info`, `parse-callback`, `validate-state` |
| `callback` | 回调事件解析 | `parse-payload`, `decrypt-payload`, `verify-signature`, `event-types` |
| `media` | 媒体文件操作 | `upload`, `upload-app`, `download`, `download-to-file`, `path` |
| `streaming` | 流式消息（AI 场景） | `create`, `fetch` |
| `chat` | 会话与消息记录 | `list`, `messages` |
| `health` | 连接健康检查 | `check` |

## 常用示例

### 消息发送

```bash
# 发送纯文本消息
lansenger message send-text chat123 "你好！"

# 发送 Markdown 消息
lansenger message send-markdown chat123 "**粗体** *斜体*"

# 发送文件
lansenger message send-file chat123 /path/to/report.pdf

# 发送网络图片
lansenger message send-image-url chat123 https://example.com/photo.jpg

# 发送链接卡片
lansenger message send-link-card chat123 "文档" "点击查看" https://docs.example.com

# 发送应用卡片
lansenger message send-app-card chat123 "卡片标题" --content "正文内容" --card-link https://example.com

# 发送多条图文（appArticles）
lansenger message send-app-articles chat123 '{"title":"文章1","url":"https://a.com"}' '{"title":"文章2","url":"https://b.com"}'

# 发送 OA 审批卡片
lansenger message send-oacard chat123 "审批标题" --head "审批通知" --field '{"key":"申请人","value":"张三"}'

# 群内发送并 @all（user_token 可选，无则显示为 bot）
lansenger message send-text group123 "全员通知" --group --mention-all

# 群内 @指定人
lansenger message send-text group123 "请查看" --group --mention staff001

# 机器人通道发送消息
lansenger message send-bot-message text '{"content":"通知内容"}' --chat-id user001 --chat-id user002
```

### 群组管理

```bash
# 创建群组
lansenger group create "项目群" org001 --staff staff001 --staff staff002

# 查看群信息
lansenger group info group123

# 查看群成员
lansenger group members group123

# 查看群列表（bot 可查看所在的群，传 user_token 可查看用户所在的群）
lansenger group list

# 查看用户所在的群列表（需要 user_token）
lansenger group list --user-token YOUR_USER_TOKEN

# 检查用户是否在群内
lansenger group check group123 --staff-id staff001

# 更新群信息
lansenger group update group123 --name "新名称" --desc "新描述"

# 添加/移除成员
lansenger group update-members group123 --add staff003 --remove staff001
```

### 人员查询

```bash
# 查看人员基本信息
lansenger staff basic-info staff001

# 查看人员详细信息
lansenger staff detail staff001

# 搜索人员
lansenger staff search "张三" --user-token YOUR_USER_TOKEN

# ID 映射（手机号 → staffId）
lansenger staff id-mapping org001 mobile 13800138000

# 查看部门祖先链
lansenger staff ancestors staff001
```

### 媒体文件

```bash
# 上传核心平台文件
lansenger media upload /path/to/file.pdf --media-type 3

# 上传应用/机器人媒体文件（用于 send-text / send-file 等）
lansenger media upload-app /path/to/file.pdf --media-type file

# 下载媒体文件到本地
lansenger media download-to-file MEDIA_ID --output /path/to/save.pdf
```

## 全局选项

| 选项 | 说明 |
|------|------|
| `--json` / `-j` | 输出原始 JSON 格式而非表格 |
| `--profile` / `-P` | 使用指定的凭证 profile（默认：`default`） |
| `--as <staff_id>` | 从凭证存储中自动加载并自动刷新指定 staff_id 的 user token |

## 多应用/多机器人配置（Profile）

CLI 支持多 profile，每个 profile 对应一个 appID（一个应用或一个机器人），凭证互相隔离：

```bash
# 配置第一个应用（个人机器人）
lansenger config set app_id xxx1 --profile my-bot
lansenger config set app_secret xxx1 --profile my-bot

# 配置第二个应用（组织机器人）
lansenger config set app_id xxx2 --profile org-bot
lansenger config set app_secret xxx2 --profile org-bot

# 删除指定 profile（如为当前 active 则自动切换到 default）
lansenger config delete-profile my-bot

# 使用指定 profile
lansenger --profile org-bot staff basic-info STAFF_ID
```

## 安全性

- 凭证存储在 `~/.lansenger/sdk_state.json`，文件权限 `0600`
- `config show` 对所有密钥类字段脱敏显示（`***`），仅 `api_gateway_url` 和 `passport_url` 明文展示
- 支持环境变量 `LANSENGER_APP_ID` / `LANSENGER_APP_SECRET` / `LANSENGER_ENCODING_KEY` / `LANSENGER_CALLBACK_TOKEN`，适合 CI/CD 场景

## 身份与权限

### 身份能力矩阵

蓝信平台有三种身份类型，各自拥有不同的 API 访问权限：

| 命令域 | 个人机器人 | 组织应用（自建） | 组织应用 + 机器人 | 备注 |
|--------|:---:|:---:|:---:|------|
| `message send-text/markdown/file/...` (bot DM) | **Y** | N | **Y** | 仅机器人可发送机器人私聊消息 |
| `message send-text --group` (群聊) | N* | N | **Y** | 个人机器人 API 支持，但暂不支持入群 |
| `message send-group-message` | N* | N | **Y** | 同上 |
| `message send-account-message` (公共号) | N | **Y** | **Y** | 需要公共号能力 |
| `message send-user-message` (用户代发) | N | **Y** | **Y** | 需要 userToken + OAuth2 |
| `message revoke` | **Y** | **Y** | **Y** | 撤回自己的消息 |
| `staff *` (通讯录只读) | N | **Y** | **Y** | `search` 额外需要 userToken |
| `department *` | N | **Y** | **Y** | 仅组织级应用 |
| `calendar *` | N | **Y** | **Y** | 带 userToken = 用户身份；不带 = 机器人身份 |
| `todo *` | N | **Y** | **Y** | 仅组织级应用 |
| `chat list/messages` | N | **Y** | **Y** | 仅组织级应用 |
| `group *` (群组管理 V2) | N | N | **Y** | 需要机器人已在群内 |
| `media upload` | **Y** | **Y** | **Y** | 通用上传 |
| `media upload-app` | N | **Y** | **Y** | 仅自建应用（非 ISV） |
| `media download/path` | **Y** | **Y** | **Y** | 通用下载 |
| `oauth *` | N | **Y** | **Y** | 仅组织级应用 |
| `streaming *` | N | **Y** | **Y** | 仅组织级应用 |
| `callback *` (事件解析) | N/A | N/A | N/A | 纯数据操作，与身份无关 |

> \* **N\*** = API 能力存在，但入群功能暂未开放。

> **个人机器人** 只能收发消息和上传/下载文件，无法访问通讯录、群组、日程或 OAuth2。
>
> **组织应用 vs 组织应用 + 机器人**：使用相同的 appID/appSecret。唯一区别在于消息通道 —— 仅机器人可以发送机器人私聊消息和群聊消息（因为只有机器人能加入群聊）。其他所有 API（通讯录、日程、待办、聊天记录、OAuth2、流式消息）两者完全一致。目前仅自建应用支持机器人能力。

### 开发者中心权限

除了身份类型，特定 API 调用还取决于蓝信开发者中心中的权限开关。组织可能限制开发者访问权限，需要管理员协助。

**基础权限（默认开启）：**

| 权限 | 说明 |
|------|------|
| Get basic user info | 获取人员基本信息，用于系统/应用登录 |
| Send notification messages | 获取组织消息通道，向人员/群组发送消息 |

**高级权限（默认关闭，需手动开启）：**

| 权限 | 说明 |
|------|------|
| Contacts read-only | 通讯录只读访问 |
| Contacts edit | 通讯录编辑访问（创建/更新/删除人员） |
| Sensitive info - Phone | 访问用户手机号 |
| Sensitive info - Email | 访问用户邮箱 |
| Sensitive info - ID number | 访问用户身份证号 |
| Sensitive info - Employee ID | 访问用户工号 |
| Map unique attribute to staff ID | 将手机/邮箱/工号映射为人员 ID |
| App edit | 创建及更新应用 |
| Groups read-only | 群组只读访问 |
| Groups edit | 群组编辑访问 |
| Calendar read-only | 日程及日历只读访问 |
| Calendar edit | 日程及日历编辑访问 |
| Upload media | 媒体文件上传权限 |
| Workbench template read | 工作台模板读取 |
| Workbench template write | 工作台模板写入 |

当遇到权限错误时，请先确认身份类型是否支持该操作，然后提示用户在开发者中心开启对应的高级权限（如无法访问请联系组织管理员）。

## CLI 兼容性

本 Go CLI 与 Python 版、TypeScript 版命令语法完全一致：

```bash
# Python CLI
pip install lansenger-cli

# Go CLI
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest

# TypeScript CLI
npm install -g lansenger-cli
```

## 许可证

MIT License
