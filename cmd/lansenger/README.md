[English](README.md) | [简体中文](README.zhHans.md) | [繁體中文](README.zhHant.md) | [繁體中文（香港）](README.zhHantHK.md) | [Français](README.fr.md)

# Lansenger CLI (Go)

Lansenger command-line tool — interact with Lansenger APIs directly from the terminal: send messages, manage groups, query staff/departments, operate calendars and todos, and more.

Command syntax is consistent with the Python and TypeScript versions. Install any one.

## Install

```bash
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest
```

Or build from source:

```bash
git clone https://github.com/lansenger-pm/lansenger-sdk-go.git
cd lansenger-sdk-go/cmd/lansenger
go build -o lansenger .
```

Requires Go 1.26+.

## Quick Start

### 1. Configure Credentials

Save credentials via `config set` (stored per profile in `~/.lansenger/sdk_state.json`, keys masked, file permissions 0600):

**Required credentials**:

```bash
lansenger config set app_id YOUR_APP_ID
lansenger config set app_secret YOUR_APP_SECRET
lansenger config set api_gateway_url https://open.e.lanxin.cn/open/apigw
```

**OAuth2 user auth (fill in when you need userToken)**:

```bash
lansenger config set passport_url https://passport.lx.qianxin.com
lansenger config set redirect_uri http://localhost:8765   # OAuth2 redirect URI (default)
```

**Callback receiver (fill in when you need to parse/verify webhook callbacks)**:

```bash
lansenger config set encoding_key YOUR_ENCODING_KEY
lansenger config set callback_token YOUR_CALLBACK_TOKEN
```

You can also configure via environment variables (CI/CD friendly):

```bash
export LANSENGER_APP_ID=YOUR_APP_ID
export LANSENGER_APP_SECRET=YOUR_APP_SECRET
export LANSENGER_REDIRECT_URI=http://localhost:8765
```

### 2. View Configuration

```bash
lansenger config show
```

### 3. Send Your First Message

```bash
lansenger message send-text staff001 "Hello from CLI!"
```

## Command Overview

| Group | Description | Subcommands |
|--------|------|--------|
| `config` | Manage credentials | `set`, `show`, `clear`, `list-profiles`, `delete-profile`, `list-users` |
| `message` | Send & manage messages | `send-text`, `send-markdown`, `send-file`, `send-image-url`, `send-link-card`, `send-app-articles`, `send-app-card`, `send-oacard`, `send-bot-message`, `send-group-message`, `send-account-message`, `send-user-message`, `update-dynamic-card`, `revoke`, `query-groups`, `send-reminder` |
| `group` | Manage groups | `create`, `info`, `members`, `list`, `check`, `update`, `update-members`, `dismiss` |
| `staff` | Query staff info | `basic-info`, `detail`, `ancestors`, `id-mapping`, `org-extra-fields`, `search`, `org-info` |
| `department` | Query department info | `detail`, `children`, `staffs` |
| `calendar` | Calendar & schedule | `primary`, `create-schedule`, `fetch-schedule`, `delete-schedule`, `list-schedules`, `attendees`, `add-attendees`, `delete-attendees`, `update-schedule`, `attendee-meta` |
| `todo` | Todo task management | `create`, `update`, `update-status`, `delete`, `list`, `fetch-by-id`, `fetch-by-source`, `status-counts`, `executor-status`, `add-executors`, `delete-executors`, `executor-list` |
| `oauth` | OAuth2 user auth | `authorize-url`, `exchange-code`, `refresh-token`, `user-info`, `parse-callback`, `validate-state` |
| `callback` | Callback event parsing | `parse-payload`, `decrypt-payload`, `verify-signature`, `event-types` |
| `media` | Media file operations | `upload`, `upload-app`, `download`, `download-to-file`, `path` |
| `streaming` | Streaming messages (AI) | `create`, `fetch` |
| `chat` | Conversations & messages | `list`, `messages` |
| `health` | Connection health check | `check` |

## Common Examples

### Messaging

```bash
# Send plain text
lansenger message send-text chat123 "Hello!"

# Send markdown
lansenger message send-markdown chat123 "**Bold** *italic*"

# Send file
lansenger message send-file chat123 /path/to/report.pdf

# Send image from URL
lansenger message send-image-url chat123 https://example.com/photo.jpg

# Send link card
lansenger message send-link-card chat123 "Documentation" "Read this" https://docs.example.com

# Send app card
lansenger message send-app-card chat123 "Card Title" --content "Body text" --card-link https://example.com

# Send app articles
lansenger message send-app-articles chat123 '{"title":"Article 1","url":"https://a.com"}' '{"title":"Article 2","url":"https://b.com"}'

# Send OA approval card
lansenger message send-oacard chat123 "Approval Title" --head "Notification" --field '{"key":"Applicant","value":"John"}'

# Send in group with @all (user_token optional, shows as bot without it)
lansenger message send-text group123 "Announcement" --group --mention-all

# @mention specific people in group
lansenger message send-text group123 "Please check" --group --mention staff001

# @mention specific bots in group
lansenger message send-text group123 "Bot check" --group --mention-bot bot001 --mention-bot bot002

# Reply to a message (message reference)
lansenger message send-text group123 "Got it" --group --ref-msg-id 524288-xxx

# Bot channel broadcast
lansenger message send-bot-message text '{"content":"Notice"}' --chat-id user001 --chat-id user002

# Bot channel reply (message reference)
lansenger message send-bot-message text '{"content":"Reply"}' --chat-id user001 --ref-msg-id 524288-xxx

# Query group ID list
lansenger message query-groups --page 0 --size 100
```

### Group Management

```bash
# Create group
lansenger group create "Project Group" org001 --staff staff001 --staff staff002

# View group info
lansenger group info group123

# View group members
lansenger group members group123

# View group list (bot can list groups it belongs to)
lansenger group list

# View group list as user (requires user_token)
lansenger group list --user-token YOUR_USER_TOKEN

# Check if user is in group
lansenger group check group123 --staff-id staff001

# Update group info
lansenger group update group123 --name "New Name" --desc "New Description"

# Add/remove members
lansenger group update-members group123 --add staff003 --remove staff001
```

### Staff Query

```bash
# Basic staff info
lansenger staff basic-info staff001

# Detailed staff info
lansenger staff detail staff001

# Search staff
lansenger staff search "张三" --user-token YOUR_USER_TOKEN

# ID mapping (phone → staffId)
lansenger staff id-mapping org001 mobile 13800138000

# Department ancestors
lansenger staff ancestors staff001
```

### Media Files

```bash
# Upload core platform file
lansenger media upload /path/to/file.pdf --media-type 3

# Upload app/bot media file (used for send-text / send-file etc.)
lansenger media upload-app /path/to/file.pdf --media-type file

# Download media to local file
lansenger media download-to-file MEDIA_ID --output /path/to/save.pdf
```

## Global Options

| Option | Description |
|------|------|
| `--json` / `-j` | Output raw JSON instead of formatted tables |
| `--profile` / `-P` | Use a specific credential profile (default: `default`) |
| `--as <staff_id>` | Auto-load & auto-refresh user token for the given staff_id from credential store |

## Multi-app / Multi-bot Profiles

CLI supports multiple profiles, each corresponding to one appID (one app or bot), with isolated credentials:

```bash
# Configure first app (personal bot)
lansenger config set app_id xxx1 --profile my-bot
lansenger config set app_secret xxx1 --profile my-bot

# Configure second app (organization bot)
lansenger config set app_id xxx2 --profile org-bot
lansenger config set app_secret xxx2 --profile org-bot

# Delete a profile (auto-switches to default if active)
lansenger config delete-profile my-bot

# Use a specific profile
lansenger --profile org-bot staff basic-info STAFF_ID
```

## Security

- Credentials stored in `~/.lansenger/sdk_state.json` with `0600` permissions
- `config show` masks all secret fields (`***`), only `api_gateway_url` and `passport_url` shown in plaintext
- Environment variables `LANSENGER_APP_ID` / `LANSENGER_APP_SECRET` / `LANSENGER_ENCODING_KEY` / `LANSENGER_CALLBACK_TOKEN` supported for CI/CD

## Identity & Permissions

### Identity Capability Matrix

The Lansenger platform has three identity types with different API access:

| Command Domain | Personal Bot | Org App (Self-built) | Org App + Bot | Notes |
|--------|:---:|:---:|:---:|------|
| `message send-text/markdown/file/...` (bot DM) | **Y** | N | **Y** | Only bots can send bot DMs |
| `message send-text --group` (group chat) | **Y** | N | **Y** | Personal bot API supports it |
| `message send-group-message` | **Y** | N | **Y** | Same as above |
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


> **Personal Bot** can only send/receive messages and upload/download files. Cannot access contacts, calendars, or OAuth2.
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

## CLI Compatibility

This Go CLI shares the same command syntax as the Python and TypeScript versions:

```bash
# Python CLI
pip install lansenger-cli

# Go CLI
go install github.com/lansenger-pm/lansenger-sdk-go/cmd/lansenger@latest

# TypeScript CLI
npm install -g lansenger-cli
```

## License

MIT License
