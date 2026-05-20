package lansenger

const (
	DefaultAPIGatewayURL = "https://open.e.lanxin.cn/open/apigw"
	DefaultStateDir      = "~/.lansenger"
	DefaultStateFile     = "sdk_state.json"
	DefaultProfile       = "default"
	TokenRefreshMargin   = 300

	MediaTypeVideo = 1
	MediaTypeImage = 2
	MediaTypeFile  = 3

	MaxMessageLength = 4000

	OAuth2ScopeBasicUserInfo = "basic_userinfor"

	TodoStatusPendingRead = "11"
	TodoStatusRead        = "12"
	TodoStatusPendingDo   = "21"
	TodoStatusDone        = "22"

	TodoTypeNotification = 1
	TodoTypeApproval     = 2
)

var ImageExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".webp": true,
	".gif":  true,
}

var VideoExtensions = map[string]bool{
	".mp4":  true,
	".mov":  true,
	".avi":  true,
	".mkv":  true,
	".webm": true,
	".3gp":  true,
}

var APIEndpoints = map[string]map[string]string{
	"auth": {
		"app_token_create":    "/v1/apptoken/create",
		"tenant_access_token": "/auth/v3/tenant_access_token/internal",
	},
	"oauth": {
		"user_token_create":    "/v2/user_token/create",
		"refresh_token_create": "/v1/refresh_token/create",
	},
	"users": {
		"fetch": "/v1/users/fetch",
	},
	"staffs": {
		"basic_info_fetch":           "/v1/staffs/{staff_id}/fetch",
		"detail_fetch":               "/v1/staffs/{staff_id}/infor/fetch",
		"department_ancestors_fetch": "/v1/staffs/{staff_id}/departmentancestors/fetch",
		"id_mapping_fetch":           "/v2/staffs/id_mapping/fetch",
		"search":                     "/v2/staffs/search",
	},
	"org": {
		"extra_field_ids_fetch": "/v1/org/{org_id}/extrafieldids/fetch",
		"info_fetch":            "/v1/org/{org_id}/fetch",
	},
	"departments": {
		"detail_fetch":   "/v1/departments/{department_id}/fetch",
		"children_fetch": "/v1/departments/{department_id}/children/fetch",
		"staffs_fetch":   "/v1/departments/{department_id}/staffs/fetch",
	},
	"groups_v2": {
		"create":         "/v2/groups/create",
		"info_fetch":     "/v2/groups/{group_id}/info/fetch",
		"members_fetch":  "/v2/groups/{group_id}/members/fetch",
		"list_fetch":     "/v2/groups/fetch",
		"is_in_group":    "/v2/groups/{group_id}/members/is_in_group",
		"info_update":    "/v2/groups/{group_id}/info/update",
		"members_update": "/v2/groups/{group_id}/members/update",
	},
	"messages": {
		"create":         "/v1/messages/create",
		"chat_create":    "/v1/messages/chat/create",
		"group_create":   "/v1/messages/group/create",
		"revoke":         "/v1/messages/revoke",
		"dynamic_update": "/v1/messages/dynamic/update",
		"fetch":          "/v1/messages/fetch",
	},
	"bot": {
		"messages_create": "/v1/bot/messages/create",
	},
	"sse": {
		"msg_create": "/v1/sse/msg/create",
		"msg_fetch":  "/v1/sse/msg/fetch",
	},
	"medias": {
		"create": "/v1/medias/create",
		"fetch":  "/v1/medias/{media_id}/fetch",
	},
	"chats": {
		"fetch": "/v1/chats/fetch",
	},
	"calendars": {
		"primary_fetch":            "/v1/calendars/primary",
		"schedules_create":         "/v1/calendars/{calendar_id}/schedules/create",
		"schedules_fetch":          "/v1/calendars/{calendar_id}/schedules/{schedule_id}/fetch",
		"schedules_delete":         "/v1/calendars/{calendar_id}/schedules/{schedule_id}/delete",
		"schedules_list_fetch":     "/v1/calendars/{calendar_id}/schedules/fetch",
		"schedules_members_fetch":  "/v1/calendars/{calendar_id}/schedules/{schedule_id}/members/fetch",
		"schedules_members_create": "/v1/calendars/{calendar_id}/schedules/{schedule_id}/members/create",
		"schedules_members_delete": "/v1/calendars/{calendar_id}/schedules/{schedule_id}/members/delete",
	},
	"todo": {
		"create":                  "/xtra/task/unified/v1/todotask/create",
		"info_update":             "/xtra/task/unified/v1/todotask/info/update",
		"status_update":           "/xtra/task/unified/v1/todotask/status/update",
		"sender_delete":           "/xtra/task/unified/v1/sender/todotask/delete",
		"list_fetch":              "/xtra/task/unified/v1/todotask/list/fetch",
		"info_fetch_by_source_id": "/xtra/task/unified/v1/todotask/info/fetchbysourceid",
		"info_fetch":              "/xtra/task/unified/v1/todotask/info/fetch",
		"status_count_list_fetch": "/xtra/task/unified/v1/todotask/status/countList/fetch",
		"executor_status_update":  "/xtra/task/unified/v1/todotask/executor/status/update",
		"executor_create":         "/xtra/task/unified/v1/todotask/executor/create",
		"executor_delete":         "/xtra/task/unified/v1/todotask/executor/delete",
		"executor_list_fetch":     "/xtra/task/unified/v1/todotask/executor/list/fetch",
	},
	"websocket": {
		"endpoint": "/v1/ws/endpoint/create",
	},
}

var CallbackEventTypes = map[string]string{
	"account_subscribe":    "public_account",
	"account_unsubscribe":  "public_account",
	"staff_info":           "staff",
	"staff_modify":         "staff",
	"staff_create":         "staff",
	"staff_delete":         "staff",
	"dept_modify":          "department",
	"dept_create":          "department",
	"dept_delete":          "department",
	"tag_member":           "tag",
	"app_install_org":      "app",
	"app_uninstall_org":    "app",
	"bot_private_message":  "bot",
	"bot_group_message":    "bot",
	"group_create_approve": "group",
	"telephone_track":      "notification",
	"ua_cert_create":       "certificate",
	"ua_cert_delete":       "certificate",
	"report_location":      "location",
	"user_logout":          "auth",
	"data_scope":           "data_scope",
	"wb_visible_config":    "workbench",
	"schedule_modify":      "calendar",
	"schedule_delete":      "calendar",
}

func GuessMediaType(filePath string) int {
	ext := ""
	for i := len(filePath) - 1; i >= 0; i-- {
		if filePath[i] == '.' {
			ext = filePath[i:]
			break
		}
	}
	if ImageExtensions[ext] {
		return MediaTypeImage
	}
	if VideoExtensions[ext] {
		return MediaTypeVideo
	}
	return MediaTypeFile
}
