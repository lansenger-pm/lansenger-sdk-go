# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [0.14.1] - 2026-09-23

### Added

- **personal_todos**: 挂附件条目构造助手。上传接口返回 `mimeType`/`size`，而挂附件写体的条目必须叫 `fileType`/`fileSize`，直接把上传响应塞进 `resources` 会被后端以 errCode 500 打回——两个入口统一做这层映射（与 Python/TS SDK 同名助手等价）：
  - `BuildPersonalTodoResourceEntry(resourceID, fileName, fileType, fileSize, opt)`：显式构造（条目结构的唯一定义处）；
  - `PersonalTodoResourceEntryFromUpload(upload, opt)`：从上传结果构造，`upload` 接受 `*PersonalTodoResourceResult` **或原始响应 map**（含内层 `data`）两种形态；无法解析时返回 error，不会静默产出缺 `resourceId` 的条目。

### Fixed

- **videoconferences**: `ModifyMeeting` 补齐 `UserStopTime`（自动结束时间）透传，与 `CreateMeeting` 对齐。此前该参数只在创建时可传——修改时不生效，且会把会议结束时间重置为开始时间 +24 小时。
- **videoconferences**: `ControlMember` 移除 `opCode` 客户端校验，改为原样透传（服务端才是权威），`VCOpCodes` 保留为已知值参考表。客户端硬校验既拦掉了服务端实际接受的取值（如 `mute`），又放行了服务端不认的值（如 `applyAudio`）。
- **cli**: `videoconference modify` 新增 `--user-stop-time`；`member-control` 移除 `OP_CODE` 白名单校验。

---

## [0.14.0] - 2026-09-23

### Fixed

- **calendars**: default `attendeeFlag` changed from `"required"` to `"yes"` — not in the server enum (`yes`/`option`/`no`), made the auto-fill path fail with 40060 (LXBUGS-128487). Added `ATTENDEE_FLAGS` local validation before any HTTP call; docstrings corrected to the unix-seconds time structure.
- **contacts**: staff search applies `page_size` alone by defaulting `page=1` (LXBUGS-128510); documented unreliable `hasMore`.
- **models**: chat message plain-text extraction now parses appCard/i18nAppCard, linkCard, appArticles, formatText (`text` + legacy `content` keys) and object text shapes, with a depth-limited fallback scan (LXBUGS-128493).
- **streaming**: stream message creation downgrades success-with-empty-msgId to a failure (LXBUGS-128497).
- **groups**: is-in-group error explains the errCode=10000 ambiguity (LXBUGS-128498).
- **videoconferences**: VC ops annotated — `muteall`/`unmuteall` rejected live with 105601 (LXBUGS-128490).
- **media/client**: documented download/health-check return shapes (LXBUGS-128494).
- READMEs: attendee example fixed.

---

## [0.13.0] - 2026-09-20

### Added

- **videoconferences**: 视频会议开放能力域（`/xtra/videoconference/openapi/v1/*`，20 端点）——会议创建/修改/取消/结束、详情/列表/操作记录/进出记录、固定会议室、批量状态、事件订阅、会议参数、历史与进行中会议、主持人会控、成员邀请与列表、录像列表与下载链接、组织配置。前提：组织安装视频会议应用并在开发者中心开通开放能力。
- **cli**: `videoconference` 命令组，覆盖上述 20 个端点；`--member`/`--vods` 以 JSON 入参，`cancel`/`stop` 接入高风险门禁（`--yes`/`--dry-run`）。

### Notes

- 创建/修改要求成员列表恰好一名 `role="admin"` 主持人；录像下载链接每次最多 3 个；`fetchRange='person'` 必须传 `staffId`；`mids` 不得为空。

---

## [0.15.0] - 2026-09-24

### Fixed

- **dynamic-card**: `UpdateDynamicCard` / `DynamicCardUpdateParams` now accept `UserToken` (query) and `UserId` (body) — the update must carry the identity that originally SENT the card. Bot-sent cards update with the app identity alone; human-sent cards (e.g. group messages) require the sender identity; a mismatch fails with 10005 无权限 (LXBUGS-128492, verified live by 李川 lxtest).
- **videoconferences**: `VCOpCodes` removed `muteall`/`unmuteall` — `/meeting/member/control` is per-member only (backend confirmation 邹治会 2026-09-24); `mute` (per-member mute) verified live errCode=0 (LXBUGS-128490).

---

## [0.12.1] - 2026-09-20

### Added

- **personal_todos**: 个人待办 `/xtra/tdtask/server/openapi/...` 6 个端点 — 创建、按字段编辑、用户待办分页，以及资源上传、下载 URL、预签名上传 URL。
- **cli**: `personal-todo` 命令组 6 个子命令。
- **models**: `PersonalTodoSaveResult` / `PersonalTodoListResult` / `PersonalTodoResourceResult` / `PersonalTodoURLResult`。
- **notices**: 新增通知发送、官方账号查询及 CLI 命令。
- **questionnaires**: 新增问卷系统 22 个接口和 CLI 命令组。
- **boardrooms**: 新增会议室预定 V2 全部 11 个接口和 CLI 命令组。

### Changed

- **http**: `doPost` 兼容成功码 `0` 和旧环境写接口成功码 `200`。
- **notices**: 默认使用注入身份，校验创建人字段，并允许显式发送经纬度 `0`。
- **questionnaires**: 修复字符串、数字、布尔 `data` 响应解析，并注入 CLI 身份字段。
- **boardrooms**: 修复取消/扫码确认的布尔响应解析，补齐 CLI 身份与可选参数。

### Notes

- 个人待办与应用身份统一待办完全分离；`orgId` 必须显式传入，编辑接口的 `orgId` 位于请求体顶层。
- 服务端当前不提供个人待办完成/删除能力。
- 创建接口 `finishTime` 默认发送 `0`，与 stage 实测可调用请求一致，不使用 `null`。

---

## [0.13.0] - 2026-09-21

### Added

- **videoconferences**: New videoconference domain (视频会议开放能力) — 20 endpoints on `/xtra/videoconference/openapi/v1`: create/modify/cancel/stop meeting, detail, meeting list, record list, member simplerecord, fixroom list, batch status, events subscribe, meeting params, history/active fetch, member control/invite/member list, vod list, vod download URLs, org conf fetch.
- **videoconferences**: Client-side validation — exactly one `role='admin'` member on create/modify, opCode whitelist for member control, 1..3 vods per download call, `fetchRange='person'` requires staffId, mids non-empty.
- **videoconferences**: Typed `VideoconferenceMember`/`VideoconferenceVod` structs, role/fetchRange/createSource/opCode constants, unit tests + domain contract cases and input guards.
- **cli**: `videoconference` command group covering all 20 endpoints — create/modify/cancel/stop/detail/list/record-list/simplerecord/fixroom-list/status/subscribe/params/history/active/member-control/invite/member-list/vod-list/vod-download/conf — with `--member`/`--vods` JSON flags, the high-risk gate (`--yes`/`--dry-run`) on cancel and stop, and local guards for admin count, opCode whitelist and vod count.

---

## [0.12.0] - 2026-08-28

### Added

- **persistence**: CredentialStore 支持 `identity_type` 身份类型持久化（`ValidIdentityTypes` / `LoadIdentityType` / `SaveIdentityType`）；CLI `config set identity_type` 命令，`config show` / `config list-profiles` 显示身份类型——用于区分个人机器人与组织应用凭证。


## [0.10.0] - 2026-07-29

### Added

- **media**: `UploadAppMediaV2()` — 4.5.5 V2 app media upload with required `userToken`, coexists with V1 (`/v2/app/medias/create`).
- **media**: `DownloadMediaByShareID()` — 4.5.6 download media by share ID (`/v1/media/share/{share_id}/fetch`).
- **media**: Unit tests for both new APIs (6 tests).
- **cli**: `upload-app-v2` and `download-share` media subcommands.

## [0.9.29] - 2026-07-29

### Added

- **cli**: `parseFieldOrJSON()` helper — parses `--fields` values as JSON with automatic `key=value` fallback, fixing PowerShell quoting issues.

### Fixed

- **cli**: `send-approve-card` now accepts `key=value` format for `--fields` parameter in addition to JSON array.

## [0.9.28] - 2026-07-20

### Added

- **auth**: `UserTokenManager` now supports **external mode** — when `config.UserToken` is set, `GetToken()` returns the provided token directly without expiry checks or auto-refresh.
- **client**: `NewClientWithConfig` now auto-initializes `UserTokenManager` when `cfg.UserToken` is non-empty, enabling external user token support.
- **cli**: External mode now auto-injects `--user-token` via `SetDefaultUserToken()` for transparent command usage.
- **docs**: All READMEs updated with external mode documentation (5 languages).

## [0.9.27] - 2026-07-17

### Changed

- **messaging**: `QueryGroups()` now delegates to `FetchGroupList()` internally (deprecated).

### Fixed

- **docs**: README version badge updated from `0.9.23` to `0.9.26`.

## [0.9.26] - 2026-07-16

### Added

- **linting**: `.golangci.yml` config with errcheck, govet, staticcheck, unused, misspell.

### Fixed

- **config_test**: Updated tests that expected hardcoded default gateway/passport URLs.

## [0.9.25] - 2026-07-16

### Added

- **calendars**: `CreateSchedule` now auto-fills attendees when empty and userID is provided.

## [0.9.24] - 2026-07-16

### Added

- **logging**: `DebugLogger` package variable with nil-safe guards. CLI `--verbose` flag wired via `SetDefaultUserToken`/`SetDefaultUserID`.

### Changed

- **config**: Removed hardcoded `DefaultAPIGatewayURL`/`DefaultPassportURL` fallbacks. Both must now be explicitly provided.

### Fixed

- **auth**: Data race in `TokenManager.GetToken` — `expiresAt` now captured before unlock.

## [0.9.23] - 2026-07-09

### Fixed

- **callbacks**: `DecryptCallbackPayload` now supports JSON-format decrypted data (in addition to the documented binary format). Some platforms return `{"random":"...","orgId":"...","appId":"...","events":[...]}` instead of the binary `random(16B)+eventsLen(4B)+orgId+appId+events` structure.

## [0.9.22] - 2026-07-09

### Added

- **oauth**: `authorize-url` command now auto-copies the URL to clipboard.

### Fixed

- **groups**: corrected API routing from `groups_v2` to `groups` in groups.go.
- **callbacks**: `bot_group_message` events now correctly extract `isAtMe`, `isAtAll`, `bots`, and `staffs` from the nested `reminder` object (was incorrectly reading from top-level, per OpenAPI 4.10.1.3 update).
- **callbacks**: `bot_group_message` parsing now maps the `magic` field.
- **callbacks**: `staff_info` field map now includes `employeeId` alias (alongside `employId`).
- **streaming**: `FetchStreamMessage` request body field corrected from `model` to `msgId` (per OpenAPI 4.6.16).

## [0.9.21] - 2026-07-01

### Added

- **config**: `Config.AppToken` and `Config.UserToken` fields for **external token mode**.
- **auth**: `TokenManager` now supports **external mode** — `GetToken()` returns externally-provided token directly.
- **messaging**: `SendApproveCard` and `UpdateApproveCard` for approveCard (审批卡片) messages (4.6.4.12/13).
- **models**: `ApproveCardParams` and `ApproveCardUpdateParams` structs.
- **calendars**: `UpdateScheduleAttendees` for batch add/delete of schedule attendees (4.23.19).
- **bot_commands**: New module with `CreateBotCommands`, `FetchBotCommands`, `DeleteBotCommands` (4.37).
- **personal_apps**: New module with CRUD + List methods for personal apps/bots (4.38).
- **models**: `ScheduleAttendeesUpdateResult`, `BotCommandResult`, `BotCommandQueryResult`, `PersonalAppCreateResult`, `PersonalAppInfoResult`, `PersonalAppListResult` structs.
- **cli**: `bot-command` (create/query/delete) and `personal-app` (create/update/info/delete/list) command groups.
- **cli**: `calendar update-attendees` command for batch add/delete schedule attendees.
- **cli**: `--app-token` and `--user-token` global flags for external token mode.
- **tests**: Test suites for bot_commands, personal_apps, and update_schedule_attendees.

### Changed

- **docs**: READMEs updated to reflect personal bots support group chat.
- **cli**: `--app-token` now triggers external mode — credential file is skipped entirely when provided.

### Fixed

- **auth**: `TokenManager` external mode prevents accidental auto-refresh of externally-provided tokens.

## [0.9.20] - 2026-06-17

### Added

- **callbacks**: BotPrivateMessageData now includes `MsgId` and `ReferenceMsg` fields; BotGroupMessageData now includes `ReferenceMsg` field.
- **messaging**: SendText, SendMarkdown, SendGroupMessage now support `reminderBotIDs` parameter for @mentioning bots.
- **messaging**: SendText, SendMarkdown, SendBotMessage, SendGroupMessage now support `refMsgID` parameter for replying to messages (prs5.9.0).
- **cli**: send-text, send-markdown, send-group-message, send-bot-message commands now support `--mention-bot` and `--ref-msg-id` flags.

### Fixed

- **cli**: query-groups command default page offset changed from 1 to 0 to match V2 API specification.

## [0.9.19] - 2026-06-17

### Fixed

- **cli**: `config list-users --show-tokens` showed 0 for both expiry times due to wrong key names (`expires_in` → `user_token_expiry`, `refresh_expires_in` → `refresh_token_expiry`).

## [0.9.18] - 2026-06-16

### Added

- **persistence**: `CredentialStore.ListUserTokens()` method to list all staffIDs with stored user tokens in the current profile.
- **persistence**: `CredentialStore.HasFullConfig()` method to check if `app_id`, `app_secret`, and `api_gateway_url` are all present for the current profile. Matches Python SDK `has_full_config` and TS SDK `hasFullConfig`.
- **client**: `GetUserTokenWithStaffID(ctx, staffID)` method to retrieve token for a specific user from CredentialStore with auto-refresh support. `GetUserToken(ctx)` maintains backward compatibility with single-user mode.
- **client**: `SetDefaultUserToken()` and `getDefaultUserToken()` functions for thread-safe fallback user token injection when no explicit `userToken` is provided to API methods.
- **cli**: `config list-users` command to list all users with stored tokens in the current profile.
- **cli**: `config list-users --show-tokens` flag to display complete token information (user_token, refresh_token, expires_in, refresh_expires_in) for each stored user.
- **cli**: `--as <staff_id>` global flag (short for "act as") that sets the default user token on the client via `SetDefaultUserToken()`, auto-loading and auto-refreshing user tokens from the CredentialStore. Works transparently — no command files were modified.
- **tests**: Test suite for `ListUserTokens` (empty, single user, multiple users, profile isolation).
- **tests**: Additional boundary tests for multi-user token isolation: auto-migration on save, no-staff_id fallback, and non-existent staff_id graceful degradation.

### Changed

- **persistence**: `CredentialStore.DeleteProfileByName(name)` now returns `bool` (`true` on success, `false` if not found) matching Python/TS SDK behavior, instead of `error`.

## [0.9.17] - 2026-06-16

### Fixed

- **cli**: `calendar list-schedules` now renders a per-row `scheduleId | summary` table instead of dumping raw JSON.

## [0.9.16] - 2026-06-16

### Fixed

- **persistence**: `ensureMigrated` Phase 2 now always merges flat fields into nested entries, even when a nested entry already exists. Fixes stale flat fields left by old SDK after migration.

### Added

- **tests**: Test for stale flat field migration cleanup.

## [0.9.15] - 2026-06-16

### Added

- **tests**: Multi-user userToken isolation test suite covering: two users not overwriting each other, cross-staff independence, legacy flat format auto-migration, and backward-compatible no-staff_id fallback.

## [0.9.14] - 2026-06-16

### Fixed

- **persistence**: Fix multi-user userToken overwrite bug in `CredentialStore`. Previously `SaveUserToken()` wrote tokens as flat fields in the profile, so each new OAuth2 authorization for the same app overwrote the previous user's tokens. Tokens are now stored per-staffID in `user_tokens[staffID]` so multiple users can coexist in the same profile. Legacy flat-format stores are auto-migrated on load.

## [0.9.13] - 2026-06-15

### Added

- **persistence**: `CredentialStore.DeleteProfileByName(name)` method to delete a specific profile by name. Automatically falls back to `"default"` if the deleted profile was the active one. Returns `true` on success, `false` if the profile does not exist.
- **cli**: `config delete-profile` command to permanently remove a credential profile and all its data.

## [0.9.12] - 2026-06-12

### Fixed

- **user_token_manager**: `refreshExpiresAt` now conditionally updated only when `RefreshExpiresIn > 0`, preventing the in-memory expiry from being reset to `time.Now()` when the API returns 0. Also added 300-second margin to refreshToken expiry check to avoid race at exact boundary. Fixed memory-disk inconsistency where persisted expiry kept the old value while in-memory was overwritten.

## [0.9.11] - 2026-06-12

### Changed

- **constants**: Merged `groups_v2` endpoint category into `groups`, removing the redundant separate category. All group APIs now use the unified `"groups"` key in `APIEndpoints`.

## [0.9.10] - 2026-06-12

### Changed

- **help docs**: Improved module-level and search command descriptions for clarity and consistency with Python CLI.
- **staff search**: `--user-token` and `--user-id` flags now clearly state "one of the two is required".

---

## [0.9.9] - 2026-06-10

### Fixed

- **SendText / SendFile / SendImageURL**: Message body `mediaType` now correctly sent as `int` (1/2/3) per OpenAPI spec, while upload still uses `UploadAppMedia` (4.5.4) with string type. Added `AppToMsgMediaType` mapping.

## [0.9.8] - 2026-06-10

### Changed

- **SendText / SendFile**: File uploads now use app/bot upload endpoint (4.5.4) instead of core service endpoint (4.5.1). The `mediaType` parameter type changed from `int` to `string` (values: `"file"`, `"video"`, `"image"`, `"audio"`).
- **SendImageURL**: Uses `UploadAppMedia` (4.5.4) with `AppMediaTypeImage` instead of `UploadMedia` (4.5.1).

### Fixed

- **GuessMediaType()**: Now returns `0` for unknown file extensions instead of `MediaTypeImage`, allowing callers to fall back to their own default.

### Added

- **GuessAppMediaType()**: New function for app/bot upload (4.5.4) — returns string type (`"file"`, `"video"`, `"image"`, `"audio"`), defaulting to `"file"` for unknown extensions.

### CLI Changes

- **send-text / send-file**: `--media-type` changed from `int` to `string` (`file`/`video`/`image`/`audio`).

## [0.9.7] - 2026-06-10

### Fixed
- `NewClientFromStore`: Now initializes `UserTokenManager` with store to auto-load user token from cache, matching Python SDK behavior.

## [0.9.6] — v0.7.2 tag — 2026-06-10

### Added
- `Config.RedirectURI` field + `LANSENGER_REDIRECT_URI` env var support
- `CredentialStore` now persists `redirect_uri` in `SaveCredentials` / `LoadCredentials`
- `BuildAuthorizeURL` uses `config.RedirectURI` as default fallback when `redirectURI` arg is empty

## [0.9.6] - 2026-06-10

### Fixed
- `SendBotMessage`: Removed non-existent `groupIdList` parameter (bot API only supports `userIdList`/`departmentIdList` per OpenAPI 4.6.12)
- `FetchStreamMessage`: Fixed body field from `msgId` to `model` (per OpenAPI 4.6.16)
- `UpdateDynamicCard`: Removed extra `appCardUpdateMsg` wrapper (per OpenAPI 4.6.5)
- `CreateGroupShareID`: Added required `creator` and `expiresIn` parameters (per OpenAPI 4.28.8)
- `GroupInfoResult.SendMsgStatus`: Changed type from `string` to `bool` (per OpenAPI 4.28.3)
- `ScheduleInfoResult`: Changed `StartTime`, `EndTime`, `Creator` types from `string` to `map[string]interface{}` (per OpenAPI 4.23.11)
- `FetchUserInfo`: Added parsing of `departments` array from API response
- `BuildAuthorizeURL`: Added auto-generation of random `state` when empty
- `MediaType` constants: Removed incorrect `MediaTypeFile=3`, fixed `MediaTypeAudio=3` (per OpenAPI 4.5.1)
- `GuessMediaType`: Default unknown files to `MediaTypeImage`

### CLI Changes
- `send-bot-message --group`: Route to `SendGroupMessage` (OpenAPI 4.6.2) instead of removing flag
- `group update`: Added 11 missing flags (`--assistant`, `--demote-assistant`, `--manage-mode`, etc.)
- `send-group-message`: Inject mention flags into `msgData`
- `calendar create-schedule`: Parse timestamps as `int64` instead of `string`
- `media upload`: Fixed help text (1=video, 2=image, 3=audio), added `created_time` output

## [0.9.5] - 2026-06-08

### Documentation
- Updated README version badges to 0.9.5

## [0.9.4] - 2026-06-07

### Miscellaneous
- Bump version to 0.9.4

## [0.9.3] - 2026-06-06

### Fixed
- Fixed indentation errors in `groups.go` and `messaging.go`

## [0.9.2] - 2026-06-05

### Miscellaneous
- Bump version to 0.9.2

## [0.9.1] - 2026-06-04

### Fixed
- Fixed parsing of chat messages API response to support both `messageInfo` and `messageInfos` fields

## [0.9.0] - 2026-06-03

### Fixed
- Fixed chat list/messages API bugs
- Fixed endpoint key consistency
- Fixed models JSON tags

## [0.8.0] - 2026-05-28

### Miscellaneous
- Bump version to 0.8.0

## [0.7.1] - 2026-05-25

### Added
- OAuth local-callback: Added `--redirect-uri` option

## [0.7.0] - 2026-05-24

### Fixed
- Token management bugs:
  - Preserve `refreshToken` on refresh
  - Subtract margin from expiry
  - Persist `refreshExpiresIn`
  - Restore expiry from store on load

## [0.6.0] - 2026-05-20

### Added
- `UserTokenManager`: Added auto-refresh functionality
- OAuth: Added local-callback support

## [0.5.1] - 2026-05-18

### Changed
- Cross-SDK spec compliance:
  - `PlainText` formatting
  - `SaveUserToken` with `refreshExpiresIn`
  - `chat --split-month`/`--progress` flags
  - OAuth `--json` flag
  - Exchange-code persistence

## [0.5.0] - 2026-05-15

### Fixed
- Aligned CLI parameters with Python SDK

[0.9.6]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.9.5...v0.9.6
[0.9.5]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.9.4...v0.9.5
[0.9.4]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.9.3...v0.9.4
[0.9.3]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.9.2...v0.9.3
[0.9.2]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.9.1...v0.9.2
[0.9.1]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.9.0...v0.9.1
[0.9.0]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.8.0...v0.9.0
[0.8.0]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.7.1...v0.8.0
[0.7.1]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.7.0...v0.7.1
[0.7.0]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.6.0...v0.7.0
[0.6.0]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.5.1...v0.6.0
[0.5.1]: https://github.com/lansenger-pm/lansenger-sdk-go/compare/v0.5.0...v0.5.1
[0.5.0]: https://github.com/lansenger-pm/lansenger-sdk-go/releases/tag/v0.5.0
