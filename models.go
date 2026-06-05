package lansenger

type SendMessageResult struct {
	Success     bool
	MessageID   string
	Error       string
	Platform    string
	MsgType     string
	Operation   string
	RawResponse map[string]interface{}
	Retryable   bool
}

type StaffBasicInfoResult struct {
	Success     bool
	OrgID       string
	OrgName     string
	Name        string
	Gender      string
	Signature   string
	AvatarURL   string
	AvatarID    string
	Status      string
	Departments []map[string]interface{}
	Error       string
	RawResponse map[string]interface{}
}

type StaffDetailResult struct {
	Success              bool
	Name                 string
	Signature            string
	AvatarID             string
	AvatarURL            string
	Status               string
	Departments          []map[string]interface{}
	Gender               string
	OrgID                string
	OrgName              string
	LoginName            string
	EmployeeNumber       string
	Email                string
	ExternalID           string
	Nationality          string
	Birthdate            string
	IDNumber             string
	NativePlace          string
	Duties               []map[string]interface{}
	Parties              []map[string]interface{}
	Address              string
	MobilePhone          string
	MobilePhoneCountryCode string
	ExtraPhones          []string
	Introduction         string
	Education            []map[string]interface{}
	Career               []map[string]interface{}
	LoginWays            []map[string]interface{}
	Tags                 []map[string]interface{}
	ExtraFieldSet        []map[string]interface{}
	Leaders              []map[string]interface{}
	JoinDate             string
	Error                string
	RawResponse          map[string]interface{}
}

type DepartmentAncestorsResult struct {
	Success        bool
	AncestorGroups [][]map[string]string
	Error          string
	RawResponse    map[string]interface{}
}

type StaffIdMappingResult struct {
	Success     bool
	StaffID     string
	Error       string
	RawResponse map[string]interface{}
}

type OrgInfoResult struct {
	Success           bool
	OrgID             string
	OrgName           string
	IconURL           string
	OrgMaxMemberLimit int
	OrgOrderType      string
	OrgDaysLimit      int
	OrgBillingDate    string
	Error             string
	RawResponse       map[string]interface{}
}

type ExtraFieldIdsResult struct {
	Success       bool
	HasMore       bool
	Total         int
	ExtraFieldIDs []string
	Error         string
	RawResponse   map[string]interface{}
}

type StaffSearchResult struct {
	Success     bool
	HasMore     bool
	Total       int
	StaffInfo   []map[string]interface{}
	Error       string
	RawResponse map[string]interface{}
}

type QueryGroupsResult struct {
	Success       bool
	TotalGroupIDs int
	GroupIDs      []string
	Error         string
	Platform      string
	Operation     string
	RawResponse   map[string]interface{}
}

type UploadMediaResult struct {
	Success     bool
	MediaID     string
	CreatedTime string
	Error       string
}

type UploadAppMediaResult struct {
	Success     bool
	MediaID     string
	Error       string
}

type DownloadMediaResult struct {
	Success bool
	Data    []byte
	Error   string
}

type MediaPathResult struct {
	Success     bool
	MediaPath   string
	Name        string
	Type        string
	Size        string
	Error       string
	RawResponse map[string]interface{}
}

type AppCardParams struct {
	BodyTitle      string
	ChatID         string
	HeadTitle      string
	BodySubTitle   string
	BodyContent    string
	Signature      string
	Fields         []map[string]interface{}
	Links          []map[string]interface{}
	CardLink       string
	PcCardLink     string
	PadCardLink    string
	IsDynamic      bool
	HeadStatusInfo map[string]interface{}
	StaffID        string
	HeadIconURL    string
	IsGroup        bool
	UserToken      string
	SenderID       string
}

type LinkCardParams struct {
	ChatID       string
	Title        string
	Link         string
	Description  string
	IconLink     string
	PcLink       string
	PadLink      string
	FromName     string
	FromIconLink string
	IsGroup      bool
	UserToken    string
	SenderID     string
}

type OaCardParams struct {
	ChatID     string
	Head       string
	Title      string
	SubTitle   string
	StaffID    string
	Fields     []map[string]interface{}
	Link       string
	PcLink     string
	PadLink    string
	CardAction map[string]interface{}
	IsGroup    bool
	UserToken  string
	SenderID   string
}

type DynamicCardUpdateParams struct {
	MsgID          string
	UserId         string
	HeadStatusInfo map[string]interface{}
	Links          []map[string]interface{}
	IsLastUpdate   bool
}

type SysMsgParams struct {
	Content string
	MediaID string
}

type UserTokenResult struct {
	Success          bool
	UserToken        string
	ExpiresIn        int
	RefreshToken     string
	RefreshExpiresIn int
	StaffID          string
	Scope            string
	State            string
	Error            string
	RawResponse      map[string]interface{}
}

type UserInfoResult struct {
	Success        bool
	StaffID        string
	Name           string
	OrgID          string
	OrgName        string
	AvatarID       string
	AvatarURL      string
	Email          string
	EmployeeNumber string
	LoginName      string
	ExternalID     string
	Departments    []map[string]interface{}
	Error          string
	RawResponse    map[string]interface{}
}

type AccountMessageResult struct {
	Success           bool
	MessageID         string
	InvalidStaff      []string
	InvalidDepartment []string
	Error             string
	RawResponse       map[string]interface{}
}

type UserMessageResult struct {
	Success     bool
	MessageID   string
	Error       string
	RawResponse map[string]interface{}
}

type BotMessageResult struct {
	Success           bool
	MessageID         string
	InvalidStaff      []string
	InvalidDepartment []string
	Error             string
	RawResponse       map[string]interface{}
}

type StreamMessageResult struct {
	Success     bool
	MessageID   string
	Error       string
	RawResponse map[string]interface{}
}

type GroupCreateInfo struct {
	Name                 string
	OrgID                int
	OwnerID              string
	Description          string
	AvatarID             string
	StaffIDList          []string
	DepartmentIDList     []string
	ApplyRequestID       string
	ApplyNotes           string
	ApplyGlobalUniqueID  string
	ApplySessionUniqueID string
	I18nApplyNotes       map[string]interface{}
}

type CreateGroupResult struct {
	Success           bool
	GroupID           string
	TotalMembers      int
	InvalidStaff      []string
	InvalidDepartment []string
	Error             string
	RawResponse       map[string]interface{}
}

type GroupInfoResult struct {
	Success            bool
	Name               string
	Description        string
	AvatarID           string
	AvatarURL          string
	OwnerStaffID       string
	OwnerName          string
	CreatorStaffID     string
	CreatorName        string
	State              string
	ManageMode         string
	LocationShare      bool
	NeedsConfirm       bool
	IsPublic           bool
	MaxMembers         int
	MaxHistoryMsgCount int
	TotalMembers       int
	RemindAll          bool
	SendMsgStatus      string
	Error              string
	RawResponse        map[string]interface{}
}

type GroupMemberResult struct {
	Success      bool
	TotalMembers int
	Members      []map[string]interface{}
	Error        string
	RawResponse  map[string]interface{}
}

type UpdateGroupResult struct {
	Success     bool
	Error       string
	RawResponse map[string]interface{}
}

type UpdateGroupMembersResult struct {
	Success           bool
	TotalMembers      int
	AddedStaffCount   int
	DeletedStaffCount int
	InvalidStaff      []string
	InvalidDepartment []string
	Error             string
	RawResponse       map[string]interface{}
}

type GroupListResult struct {
	Success       bool
	TotalGroupIDs int
	GroupIDs      []string
	Error         string
	RawResponse   map[string]interface{}
}

type IsInGroupResult struct {
	Success     bool
	IsInGroup   bool
	Error       string
	RawResponse map[string]interface{}
}

type DepartmentDetailResult struct {
	Success               bool
	ID                    string
	Name                  string
	ExternalID            string
	ParentID              string
	Order                 float64
	HasChildren           bool
	NormalMembers         int
	InactiveMembers       int
	FrozenMembers         int
	DeletedMembers        int
	NormalMembersUnique   int
	InactiveMembersUnique int
	FrozenMembersUnique   int
	DeletedMembersUnique  int
	Tags                  []map[string]interface{}
	AncestorDepartments   []map[string]interface{}
	Leaders               []map[string]interface{}
	Emails                []string
	Phones                []string
	Addresses             []string
	Introductions        []string
	DeptType              int
	Error                 string
	RawResponse           map[string]interface{}
}

type DepartmentChildrenResult struct {
	Success     bool
	Departments []map[string]interface{}
	Error       string
	RawResponse map[string]interface{}
}

type DepartmentStaffsResult struct {
	Success     bool
	HasMore     bool
	Total       int
	Staffs      []map[string]interface{}
	Error       string
	RawResponse map[string]interface{}
}

type TodoTaskCreateResult struct {
	Success     bool
	TodotaskID  string
	Error       string
	RawResponse map[string]interface{}
}

type TodoTaskInfoResult struct {
	Success     bool
	TodotaskID  string
	SourceID    string
	Title       string
	Desc        string
	Status      string
	Type        int
	Link        string
	PcLink      string
	SenderID    string
	ExecutorIDs []string
	CreateTime  string
	AppID       string
	Error       string
	RawResponse map[string]interface{}
}

type TodoTaskListResult struct {
	Success      bool
	Total        int
	TodotaskList []map[string]interface{}
	Error        string
	RawResponse  map[string]interface{}
}

type TodoTaskStatusCountResult struct {
	Success      bool
	StatusCounts []map[string]interface{}
	Error        string
	RawResponse  map[string]interface{}
}

type TodoTaskExecutorListResult struct {
	Success      bool
	Total        int
	ExecutorList []map[string]interface{}
	Error        string
	RawResponse  map[string]interface{}
}

type CalendarPrimaryResult struct {
	Success     bool
	CalendarID  string
	Summary     string
	Description string
	Permissions string
	Color       string
	Type        string
	Role        string
	Error       string
	RawResponse map[string]interface{}
}

type ScheduleCreateResult struct {
	Success     bool
	ScheduleID  string
	Error       string
	RawResponse map[string]interface{}
}

type ScheduleDeleteResult struct {
	Success      bool
	ScheduleIDs  []string
	Error        string
	RawResponse  map[string]interface{}
}

type ScheduleAttendeesDeleteResult struct {
	Success      bool
	ScheduleIDs  []string
	Error        string
	RawResponse  map[string]interface{}
}

type ScheduleInfoResult struct {
	Success            bool
	ScheduleID         string
	Summary            string
	Description        string
	RepeatType         string
	AllDay             string
	StartTime          string
	EndTime            string
	Creator            string
	RsvpStatus         string
	PrimaryScheduleID  string
	ExpireDateType     string
	AttendeePermissions string
	Color              string
	Error              string
	RawResponse        map[string]interface{}
}

type ScheduleListResult struct {
	Success      bool
	ScheduleList []map[string]interface{}
	Error        string
	RawResponse  map[string]interface{}
}

type ScheduleAttendeesResult struct {
	Success     bool
	Total       int
	Attendees   []map[string]interface{}
	Error       string
	RawResponse map[string]interface{}
}

type ChatStaffInfo struct {
	StaffID     string
	StaffName   string
	SectorNames []string
}

type ChatGroupInfo struct {
	GroupID   string
	GroupName string
}

type ChatListResult struct {
	Success     bool
	StaffInfos  []ChatStaffInfo
	GroupInfos  []ChatGroupInfo
	Error       string
	RawResponse map[string]interface{}
}

type ChatMessageInfo struct {
	SendTime    string
	Sender      string
	MessageType string
	Content     map[string]interface{}
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
	Success     bool
	HasMore     bool
	Total       int
	LastVersion string
	Name        string
	ChatType    string
	Messages    []ChatMessageInfo
	Error       string
	RawResponse map[string]interface{}
}
