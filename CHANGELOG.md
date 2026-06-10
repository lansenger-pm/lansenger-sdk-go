# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

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
