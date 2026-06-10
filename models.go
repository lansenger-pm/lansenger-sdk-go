package lansenger

type SendMessageResult struct {
	Success     bool                   `json:"success"`
	MessageID   string                 `json:"message_id"`
	Error       string                 `json:"error"`
	Platform    string                 `json:"platform"`
	MsgType     string                 `json:"msg_type"`
	Operation   string                 `json:"operation"`
	RawResponse map[string]interface{} `json:"raw_response"`
	Retryable   bool                   `json:"retryable"`
}

type StaffBasicInfoResult struct {
	Success     bool                   `json:"success"`
	OrgID       string                 `json:"org_id"`
	OrgName     string                 `json:"org_name"`
	Name        string                 `json:"name"`
	Gender      string                 `json:"gender"`
	Signature   string                 `json:"signature"`
	AvatarURL   string                 `json:"avatar_url"`
	AvatarID    string                 `json:"avatar_id"`
	Status      string                 `json:"status"`
	Departments []map[string]interface{} `json:"departments"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type StaffDetailResult struct {
	Success                bool                   `json:"success"`
	Name                   string                 `json:"name"`
	Signature              string                 `json:"signature"`
	AvatarID               string                 `json:"avatar_id"`
	AvatarURL              string                 `json:"avatar_url"`
	Status                 string                 `json:"status"`
	Departments            []map[string]interface{} `json:"departments"`
	Gender                 string                 `json:"gender"`
	OrgID                  string                 `json:"org_id"`
	OrgName                string                 `json:"org_name"`
	LoginName              string                 `json:"login_name"`
	EmployeeNumber         string                 `json:"employee_number"`
	Email                  string                 `json:"email"`
	ExternalID             string                 `json:"external_id"`
	Nationality            string                 `json:"nationality"`
	Birthdate              string                 `json:"birthdate"`
	IDNumber               string                 `json:"id_number"`
	NativePlace            string                 `json:"native_place"`
	Duties                 []map[string]interface{} `json:"duties"`
	Parties                []map[string]interface{} `json:"parties"`
	Address                string                 `json:"address"`
	MobilePhone            string                 `json:"mobile_phone"`
	MobilePhoneCountryCode string                 `json:"mobile_phone_country_code"`
	ExtraPhones            []string               `json:"extra_phones"`
	Introduction           string                 `json:"introduction"`
	Education              []map[string]interface{} `json:"education"`
	Career                 []map[string]interface{} `json:"career"`
	LoginWays              []map[string]interface{} `json:"login_ways"`
	Tags                   []map[string]interface{} `json:"tags"`
	ExtraFieldSet          []map[string]interface{} `json:"extra_field_set"`
	Leaders                []map[string]interface{} `json:"leaders"`
	JoinDate               string                 `json:"join_date"`
	Error                  string                 `json:"error"`
	RawResponse            map[string]interface{} `json:"raw_response"`
}

type DepartmentAncestorsResult struct {
	Success        bool                   `json:"success"`
	AncestorGroups [][]map[string]string  `json:"ancestor_groups"`
	Error          string                 `json:"error"`
	RawResponse    map[string]interface{} `json:"raw_response"`
}

type StaffIdMappingResult struct {
	Success     bool                   `json:"success"`
	StaffID     string                 `json:"staff_id"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type OrgInfoResult struct {
	Success           bool                   `json:"success"`
	OrgID             string                 `json:"org_id"`
	OrgName           string                 `json:"org_name"`
	IconURL           string                 `json:"icon_url"`
	OrgMaxMemberLimit int                    `json:"org_max_member_limit"`
	OrgOrderType      string                 `json:"org_order_type"`
	OrgDaysLimit      int                    `json:"org_days_limit"`
	OrgBillingDate    string                 `json:"org_billing_date"`
	Error             string                 `json:"error"`
	RawResponse       map[string]interface{} `json:"raw_response"`
}

type ExtraFieldIdsResult struct {
	Success       bool                   `json:"success"`
	HasMore       bool                   `json:"has_more"`
	Total         int                    `json:"total"`
	ExtraFieldIDs []string               `json:"extra_field_ids"`
	Error         string                 `json:"error"`
	RawResponse   map[string]interface{} `json:"raw_response"`
}

type StaffSearchResult struct {
	Success     bool                   `json:"success"`
	HasMore     bool                   `json:"has_more"`
	Total       int                    `json:"total"`
	StaffInfo   []map[string]interface{} `json:"staff_info"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type QueryGroupsResult struct {
	Success       bool                   `json:"success"`
	TotalGroupIDs int                    `json:"total_group_ids"`
	GroupIDs      []string               `json:"group_ids"`
	Error         string                 `json:"error"`
	Platform      string                 `json:"platform"`
	Operation     string                 `json:"operation"`
	RawResponse   map[string]interface{} `json:"raw_response"`
}

type UploadMediaResult struct {
	Success     bool   `json:"success"`
	MediaID     string `json:"media_id"`
	CreatedTime string `json:"created_time"`
	Error       string `json:"error"`
}

type UploadAppMediaResult struct {
	Success bool   `json:"success"`
	MediaID string `json:"media_id"`
	Error   string `json:"error"`
}

type DownloadMediaResult struct {
	Success bool   `json:"success"`
	Data    []byte `json:"data"`
	Error   string `json:"error"`
}

type MediaPathResult struct {
	Success     bool                   `json:"success"`
	MediaPath   string                 `json:"media_path"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Size        string                 `json:"size"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type AppCardParams struct {
	BodyTitle      string                 `json:"body_title"`
	ChatID         string                 `json:"chat_id"`
	HeadTitle      string                 `json:"head_title"`
	BodySubTitle   string                 `json:"body_sub_title"`
	BodyContent    string                 `json:"body_content"`
	Signature      string                 `json:"signature"`
	Fields         []map[string]interface{} `json:"fields"`
	Links          []map[string]interface{} `json:"links"`
	CardLink       string                 `json:"card_link"`
	PcCardLink     string                 `json:"pc_card_link"`
	PadCardLink    string                 `json:"pad_card_link"`
	IsDynamic      bool                   `json:"is_dynamic"`
	HeadStatusInfo map[string]interface{} `json:"head_status_info"`
	StaffID        string                 `json:"staff_id"`
	HeadIconURL    string                 `json:"head_icon_url"`
	IsGroup        bool                   `json:"is_group"`
	UserToken      string                 `json:"user_token"`
	SenderID       string                 `json:"sender_id"`
}

type LinkCardParams struct {
	ChatID       string `json:"chat_id"`
	Title        string `json:"title"`
	Link         string `json:"link"`
	Description  string `json:"description"`
	IconLink     string `json:"icon_link"`
	PcLink       string `json:"pc_link"`
	PadLink      string `json:"pad_link"`
	FromName     string `json:"from_name"`
	FromIconLink string `json:"from_icon_link"`
	IsGroup      bool   `json:"is_group"`
	UserToken    string `json:"user_token"`
	SenderID     string `json:"sender_id"`
}

type OaCardParams struct {
	ChatID     string                 `json:"chat_id"`
	Head       string                 `json:"head"`
	Title      string                 `json:"title"`
	SubTitle   string                 `json:"sub_title"`
	StaffID    string                 `json:"staff_id"`
	Fields     []map[string]interface{} `json:"fields"`
	Link       string                 `json:"link"`
	PcLink     string                 `json:"pc_link"`
	PadLink    string                 `json:"pad_link"`
	CardAction map[string]interface{} `json:"card_action"`
	IsGroup    bool                   `json:"is_group"`
	UserToken  string                 `json:"user_token"`
	SenderID   string                 `json:"sender_id"`
}

type DynamicCardUpdateParams struct {
	MsgID          string                 `json:"msg_id"`
	UserId         string                 `json:"user_id"`
	HeadStatusInfo map[string]interface{} `json:"head_status_info"`
	Links          []map[string]interface{} `json:"links"`
	IsLastUpdate   bool                   `json:"is_last_update"`
}

type SysMsgParams struct {
	Content string `json:"content"`
	MediaID string `json:"media_id"`
}

type UserTokenResult struct {
	Success          bool                   `json:"success"`
	UserToken        string                 `json:"user_token"`
	ExpiresIn        int                    `json:"expires_in"`
	RefreshToken     string                 `json:"refresh_token"`
	RefreshExpiresIn int                    `json:"refresh_expires_in"`
	StaffID          string                 `json:"staff_id"`
	Scope            string                 `json:"scope"`
	State            string                 `json:"state"`
	Error            string                 `json:"error"`
	RawResponse      map[string]interface{} `json:"raw_response"`
}

type UserInfoResult struct {
	Success        bool                   `json:"success"`
	StaffID        string                 `json:"staff_id"`
	Name           string                 `json:"name"`
	OrgID          string                 `json:"org_id"`
	OrgName        string                 `json:"org_name"`
	AvatarID       string                 `json:"avatar_id"`
	AvatarURL      string                 `json:"avatar_url"`
	Email          string                 `json:"email"`
	EmployeeNumber string                 `json:"employee_number"`
	LoginName      string                 `json:"login_name"`
	ExternalID     string                 `json:"external_id"`
	Departments    []map[string]interface{} `json:"department"`
	Error          string                 `json:"error"`
	RawResponse    map[string]interface{} `json:"raw_response"`
}

type AccountMessageResult struct {
	Success           bool                   `json:"success"`
	MessageID         string                 `json:"message_id"`
	InvalidStaff      []string               `json:"invalid_staff"`
	InvalidDepartment []string               `json:"invalid_department"`
	Error             string                 `json:"error"`
	RawResponse       map[string]interface{} `json:"raw_response"`
}

type UserMessageResult struct {
	Success     bool                   `json:"success"`
	MessageID   string                 `json:"message_id"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type BotMessageResult struct {
	Success           bool                   `json:"success"`
	MessageID         string                 `json:"message_id"`
	InvalidStaff      []string               `json:"invalid_staff"`
	InvalidDepartment []string               `json:"invalid_department"`
	Error             string                 `json:"error"`
	RawResponse       map[string]interface{} `json:"raw_response"`
}

type StreamMessageResult struct {
	Success     bool                   `json:"success"`
	MessageID   string                 `json:"message_id"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type GroupCreateInfo struct {
	Name                 string                 `json:"name"`
	OrgID                int                    `json:"org_id"`
	OwnerID              string                 `json:"owner_id"`
	Description          string                 `json:"description"`
	AvatarID             string                 `json:"avatar_id"`
	StaffIDList          []string               `json:"staff_id_list"`
	DepartmentIDList     []string               `json:"department_id_list"`
	ApplyRequestID       string                 `json:"apply_request_id"`
	ApplyNotes           string                 `json:"apply_notes"`
	ApplyGlobalUniqueID  string                 `json:"apply_global_unique_id"`
	ApplySessionUniqueID string                 `json:"apply_session_unique_id"`
	I18nApplyNotes       map[string]interface{} `json:"i18n_apply_notes"`
}

type CreateGroupResult struct {
	Success           bool                   `json:"success"`
	GroupID           string                 `json:"group_id"`
	TotalMembers      int                    `json:"total_members"`
	InvalidStaff      []string               `json:"invalid_staff"`
	InvalidDepartment []string               `json:"invalid_department"`
	Error             string                 `json:"error"`
	RawResponse       map[string]interface{} `json:"raw_response"`
}

type GroupInfoResult struct {
	Success            bool                   `json:"success"`
	Name               string                 `json:"name"`
	Description        string                 `json:"description"`
	AvatarID           string                 `json:"avatar_id"`
	AvatarURL          string                 `json:"avatar_url"`
	OwnerStaffID       string                 `json:"owner_staff_id"`
	OwnerName          string                 `json:"owner_name"`
	CreatorStaffID     string                 `json:"creator_staff_id"`
	CreatorName        string                 `json:"creator_name"`
	State              string                 `json:"state"`
	ManageMode         string                 `json:"manage_mode"`
	LocationShare      bool                   `json:"location_share"`
	NeedsConfirm       bool                   `json:"needs_confirm"`
	IsPublic           bool                   `json:"is_public"`
	MaxMembers         int                    `json:"max_members"`
	MaxHistoryMsgCount int                    `json:"max_history_msg_count"`
	TotalMembers       int                    `json:"total_members"`
	RemindAll          bool                   `json:"remind_all"`
	SendMsgStatus      bool                   `json:"send_msg_status"`
	Error              string                 `json:"error"`
	RawResponse        map[string]interface{} `json:"raw_response"`
}

type GroupMemberResult struct {
	Success      bool                   `json:"success"`
	TotalMembers int                    `json:"total_members"`
	Members      []map[string]interface{} `json:"members"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type UpdateGroupResult struct {
	Success     bool                   `json:"success"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type UpdateGroupMembersResult struct {
	Success           bool                   `json:"success"`
	TotalMembers      int                    `json:"total_members"`
	AddedStaffCount   int                    `json:"added_staff_count"`
	DeletedStaffCount int                    `json:"deleted_staff_count"`
	InvalidStaff      []string               `json:"invalid_staff"`
	InvalidDepartment []string               `json:"invalid_department"`
	Error             string                 `json:"error"`
	RawResponse       map[string]interface{} `json:"raw_response"`
}

type GroupListResult struct {
	Success       bool                   `json:"success"`
	TotalGroupIDs int                    `json:"total_group_ids"`
	GroupIDs      []string               `json:"group_ids"`
	Error         string                 `json:"error"`
	RawResponse   map[string]interface{} `json:"raw_response"`
}

type IsInGroupResult struct {
	Success     bool                   `json:"success"`
	IsInGroup   bool                   `json:"is_in_group"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type DepartmentDetailResult struct {
	Success               bool                   `json:"success"`
	ID                    string                 `json:"id"`
	Name                  string                 `json:"name"`
	ExternalID            string                 `json:"external_id"`
	ParentID              string                 `json:"parent_id"`
	Order                 float64                `json:"order"`
	HasChildren           bool                   `json:"has_children"`
	NormalMembers         int                    `json:"normal_members"`
	InactiveMembers       int                    `json:"inactive_members"`
	FrozenMembers         int                    `json:"frozen_members"`
	DeletedMembers        int                    `json:"deleted_members"`
	NormalMembersUnique   int                    `json:"normal_members_unique"`
	InactiveMembersUnique int                    `json:"inactive_members_unique"`
	FrozenMembersUnique   int                    `json:"frozen_members_unique"`
	DeletedMembersUnique  int                    `json:"deleted_members_unique"`
	Tags                  []map[string]interface{} `json:"tags"`
	AncestorDepartments   []map[string]interface{} `json:"ancestor_departments"`
	Leaders               []map[string]interface{} `json:"leaders"`
	Emails                []string               `json:"emails"`
	Phones                []string               `json:"phones"`
	Addresses             []string               `json:"addresses"`
	Introductions        []string               `json:"introductions"`
	DeptType              int                    `json:"dept_type"`
	Error                 string                 `json:"error"`
	RawResponse           map[string]interface{} `json:"raw_response"`
}

type DepartmentChildrenResult struct {
	Success     bool                   `json:"success"`
	Departments []map[string]interface{} `json:"departments"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type DepartmentStaffsResult struct {
	Success     bool                   `json:"success"`
	HasMore     bool                   `json:"has_more"`
	Total       int                    `json:"total"`
	Staffs      []map[string]interface{} `json:"staffs"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type TodoTaskCreateResult struct {
	Success     bool                   `json:"success"`
	TodotaskID  string                 `json:"todotask_id"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type TodoTaskInfoResult struct {
	Success     bool                   `json:"success"`
	TodotaskID  string                 `json:"todotask_id"`
	SourceID    string                 `json:"source_id"`
	Title       string                 `json:"title"`
	Desc        string                 `json:"desc"`
	Status      string                 `json:"status"`
	Type        int                    `json:"type"`
	Link        string                 `json:"link"`
	PcLink      string                 `json:"pc_link"`
	SenderID    string                 `json:"sender_id"`
	ExecutorIDs []string               `json:"executor_ids"`
	CreateTime  string                 `json:"create_time"`
	AppID       string                 `json:"app_id"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type TodoTaskListResult struct {
	Success      bool                   `json:"success"`
	Total        int                    `json:"total"`
	TodotaskList []map[string]interface{} `json:"todotask_list"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type TodoTaskStatusCountResult struct {
	Success      bool                   `json:"success"`
	StatusCounts []map[string]interface{} `json:"status_counts"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type TodoTaskExecutorListResult struct {
	Success      bool                   `json:"success"`
	Total        int                    `json:"total"`
	ExecutorList []map[string]interface{} `json:"executor_list"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type CalendarPrimaryResult struct {
	Success     bool                   `json:"success"`
	CalendarID  string                 `json:"calendar_id"`
	Summary     string                 `json:"summary"`
	Description string                 `json:"description"`
	Permissions string                 `json:"permissions"`
	Color       string                 `json:"color"`
	Type        string                 `json:"type"`
	Role        string                 `json:"role"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type ScheduleCreateResult struct {
	Success     bool                   `json:"success"`
	ScheduleID  string                 `json:"schedule_id"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type ScheduleDeleteResult struct {
	Success      bool                   `json:"success"`
	ScheduleIDs  []string               `json:"schedule_ids"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type ScheduleAttendeesDeleteResult struct {
	Success      bool                   `json:"success"`
	ScheduleIDs  []string               `json:"schedule_ids"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type ScheduleInfoResult struct {
	Success             bool                   `json:"success"`
	ScheduleID          string                 `json:"schedule_id"`
	Summary             string                 `json:"summary"`
	Description         string                 `json:"description"`
	RepeatType          string                 `json:"repeat_type"`
	AllDay              string                 `json:"all_day"`
	StartTime           map[string]interface{} `json:"start_time"`
	EndTime             map[string]interface{} `json:"end_time"`
	Creator             map[string]interface{} `json:"creator"`
	RsvpStatus          string                 `json:"rsvp_status"`
	PrimaryScheduleID   string                 `json:"primary_schedule_id"`
	ExpireDateType      string                 `json:"expire_date_type"`
	AttendeePermissions string                 `json:"attendee_permissions"`
	Color               string                 `json:"color"`
	Error               string                 `json:"error"`
	RawResponse         map[string]interface{} `json:"raw_response"`
}

type ScheduleListResult struct {
	Success      bool                   `json:"success"`
	ScheduleList []map[string]interface{} `json:"schedule_list"`
	Error        string                 `json:"error"`
	RawResponse  map[string]interface{} `json:"raw_response"`
}

type ScheduleAttendeesResult struct {
	Success     bool                   `json:"success"`
	Total       int                    `json:"total"`
	Attendees   []map[string]interface{} `json:"attendees"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type ChatStaffInfo struct {
	StaffID     string   `json:"staff_id"`
	StaffName   string   `json:"staff_name"`
	SectorNames []string `json:"sector_names"`
}

type ChatGroupInfo struct {
	GroupID   string `json:"group_id"`
	GroupName string `json:"group_name"`
}

type ChatListResult struct {
	Success     bool                   `json:"success"`
	StaffInfos  []ChatStaffInfo        `json:"staff_infos"`
	GroupInfos  []ChatGroupInfo        `json:"group_infos"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}

type ChatMessageInfo struct {
	SendTime    string                 `json:"send_time"`
	Sender      string                 `json:"sender"`
	MessageType string                 `json:"message_type"`
	Content     map[string]interface{} `json:"content"`
}

func (m *ChatMessageInfo) PlainText() string {
	if m.Content == nil {
		return ""
	}
	if text, ok := m.Content["text"].(string); ok && text != "" {
		return text
	}
	if ft, ok := m.Content["formatText"].(map[string]interface{}); ok {
		if content, ok := ft["content"].(string); ok {
			return content
		}
	}
	return ""
}

type ChatMessagesResult struct {
	Success     bool                   `json:"success"`
	HasMore     bool                   `json:"has_more"`
	Total       int                    `json:"total"`
	LastVersion string                 `json:"last_version"`
	Name        string                 `json:"name"`
	ChatType    string                 `json:"chat_type"`
	Messages    []ChatMessageInfo      `json:"messages"`
	Error       string                 `json:"error"`
	RawResponse map[string]interface{} `json:"raw_response"`
}