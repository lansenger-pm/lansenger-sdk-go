package lansenger

import (
	"context"
	"strconv"
)

const (
	VCMemberRoleHost        = "admin"
	VCMemberRoleJoinHost    = "joinHost"
	VCMemberRoleParticipant = "participant"

	VCFetchRangeMy     = "my"
	VCFetchRangeAll    = "all"
	VCFetchRangePerson = "person"

	VCCreateSourceClient     = 0
	VCCreateSourceThirdParty = 1
)

// VCOpCodes lists known opCode values for ControlMember (member/control). Reference
// only — ControlMember forwards the caller's opCode verbatim, because the server is the
// authority: a client-side check both rejected values the server accepts and accepted
// values it does not recognise. The list may be incomplete.
// Live-verified 2026-09-23: the server accepts "mute" (per-member mute, errCode 0) and
// rejects "applyAudio" (errCode 105601, opCode does not exist) — hence the swap.
// Verified live (LXBUGS-128490): "kick" works; "muteall"/"unmuteall" are rejected by
// the meeting server (errCode=105601 opCode不存在) in the tested environment. They are
// kept here pending server-side confirmation — treat 105601 as "value unsupported in
// this environment".
var VCOpCodes = []string{
	"kick", "quit", "join", "handup", "openScreenShare", "closeScreenShare",
	"openVideo", "closeVideo", "mute", "applyVideo", "shareVideo",
	"cancelShareVideo", "muteall", "unmuteall", "remove", "call",
	"enforceOpenVideo", "setJoinHost", "cancelJoinHost", "inviteOpenAudio",
	"setHost", "grabHost",
}

// VideoconferenceMember is one attendee entry ({staffId, employeeName, role}
// on create/modify; {staffId, employeeName, type, video, audio, typeMask} on
// invite). Role: admin (host) / joinHost / participant; Type: 0=platform
// member, 1=Xiaoyu device; TypeMask: 1=PSTN, 2=Xiaoyu, 3=H323.
type VideoconferenceMember struct {
	StaffID      string `json:"staffId"`
	EmployeeName string `json:"employeeName,omitempty"`
	Role         string `json:"role,omitempty"`
	Type         *int   `json:"type,omitempty"`
	Video        *int   `json:"video,omitempty"`
	Audio        *int   `json:"audio,omitempty"`
	TypeMask     *int   `json:"typeMask,omitempty"`
}

// VideoconferenceVod is one entry ({vodId}) of the vod download request.
type VideoconferenceVod struct {
	VodID interface{} `json:"vodId"`
}

// VideoconferenceCreateParams carries the fields for CreateMeeting
// (video conference /meeting/create). StartTime is epoch milliseconds; Type
// 0=instant meeting, 1=reserved (reserved requires a future StartTime).
type VideoconferenceCreateParams struct {
	Subject         string
	StartTime       int64
	Members         []VideoconferenceMember
	OrgID           string
	AutoRecord      int
	Type            int
	GroupNew        int
	ConfPassword    string
	ControlPassword string
	MaskType        int
	ExtAttr         string
	JoinMute        *int
	OpenMute        *int
	EnablePreJoin   *int
	UserStopTime    *int64
	InviteAdmin     *int
	UserToken       string
}

// VideoconferenceModifyParams carries the fields for ModifyMeeting
// (video conference /meeting/modify).
type VideoconferenceModifyParams struct {
	Mid             string
	Subject         string
	StartTime       int64
	Members         []VideoconferenceMember
	OrgID           string
	Operator        string
	AutoRecord      int
	Type            int
	GroupNew        int
	ConfPassword    string
	ControlPassword string
	UserStopTime    *int64
	UserToken       string
}

// vcOrgIDValue mirrors the Python SDK: numeric org IDs are sent as JSON
// numbers, anything else as a string.
func vcOrgIDValue(orgID string) interface{} {
	if n, err := strconv.Atoi(orgID); err == nil {
		return n
	}
	return orgID
}

// vcMidValue converts a mid to the numeric JSON value the server expects.
func vcMidValue(mid string) (int64, error) {
	if mid == "" {
		return 0, &vcClientError{msg: "mid is required"}
	}
	n, err := strconv.ParseInt(mid, 10, 64)
	if err != nil {
		return 0, &vcClientError{msg: "mid must be numeric"}
	}
	return n, nil
}

type vcClientError struct{ msg string }

func (e *vcClientError) Error() string { return e.msg }

func vcHostErr(members []VideoconferenceMember) error {
	if len(members) == 0 {
		return &vcClientError{msg: "member is required (with exactly one host)"}
	}
	hosts := 0
	for _, m := range members {
		if m.Role == VCMemberRoleHost {
			hosts++
		}
	}
	if hosts != 1 {
		return &vcClientError{msg: "exactly one member must have role='admin'"}
	}
	return nil
}

// fillVCOp fills an op result. The outer errCode has already been validated by the
// response layer, so Done means "the request completed":
//   - endpoints returning an inner {code, message} (cancel / stop / member control)
//     are judged by code == 0;
//   - endpoints returning a business object (modify returns the meeting object, which
//     carries no inner code) count as completed. Judging those by the inner code made
//     Done permanently false for ModifyMeeting, unlike CreateMeeting (detail result).
//   - a success response with no payload at all also counts as completed, so that all
//     three SDKs answer identically (the Python/TypeScript _op do the same).
func fillVCOp(res *VideoconferenceOpResult, result map[string]interface{}) {
	data := extractData(result)
	if data == nil {
		res.Done = true
		return
	}
	inner, ok := data["data"].(map[string]interface{})
	if !ok {
		inner = data
	}
	if code, hasCode := inner["code"]; hasCode {
		res.Done = vcCodeZero(code)
	} else {
		res.Done = true
	}
	res.Message = strFromMap(inner, "message")
}

func vcCodeZero(v interface{}) bool {
	switch n := v.(type) {
	case float64:
		return n == 0
	case int:
		return n == 0
	case string:
		return n == "0"
	}
	return false
}

func fillVCPage(res *VideoconferenceListResult, result map[string]interface{}) {
	data := extractData(result)
	if data == nil {
		return
	}
	res.Offset = intFromMap(data, "offset")
	res.Total = intFromMap(data, "total")
	list, ok := data["items"].([]interface{})
	if !ok {
		list, _ = data["mids"].([]interface{})
	}
	if list != nil {
		res.Items = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Items = append(res.Items, m)
			}
		}
		if len(res.Items) == 0 && len(list) > 0 {
			res.Scalars = list
		}
	}
}

func fillVCDetail(res *VideoconferenceDetailResult, data map[string]interface{}) {
	res.Mid = strFromMap(data, "id")
	if res.Mid == "" {
		res.Mid = strFromMap(data, "mid")
	}
	res.Subject = strFromMap(data, "subject")
	res.MeetingNumber = strFromMap(data, "meetingNumber")
	res.StartTime = int64FromMap(data, "startTime")
	res.StopTime = int64FromMap(data, "stopTime")
	res.Type = intFromMap(data, "type")
	res.Status = intFromMap(data, "status")
	res.Admin = strFromMap(data, "admin")
}

// CreateMeeting creates a meeting (instant or reserved) (/meeting/create).
// Members require exactly one role='admin'.
func (c *LansengerClient) CreateMeeting(ctx context.Context, p *VideoconferenceCreateParams) (*VideoconferenceDetailResult, error) {
	if p == nil {
		return &VideoconferenceDetailResult{Success: false, Error: "params is required"}, nil
	}
	if p.Subject == "" {
		return &VideoconferenceDetailResult{Success: false, Error: "subject is required"}, nil
	}
	if err := vcHostErr(p.Members); err != nil {
		return &VideoconferenceDetailResult{Success: false, Error: err.Error()}, nil
	}
	if p.StartTime == 0 {
		return &VideoconferenceDetailResult{Success: false, Error: "start_time is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_create", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"subject":         p.Subject,
		"startTime":       p.StartTime,
		"autoRecord":      p.AutoRecord,
		"type":            p.Type,
		"groupNew":        p.GroupNew,
		"confPassword":    p.ConfPassword,
		"controlPassword": p.ControlPassword,
		"member":          p.Members,
		"extAttr":         p.ExtAttr,
		"maskType":        p.MaskType,
		"orgId":           vcOrgIDValue(p.OrgID),
	}
	if p.JoinMute != nil {
		body["joinMute"] = *p.JoinMute
	}
	if p.OpenMute != nil {
		body["openMute"] = *p.OpenMute
	}
	if p.EnablePreJoin != nil {
		body["enablePreJoin"] = *p.EnablePreJoin
	}
	if p.UserStopTime != nil {
		body["userStopTime"] = *p.UserStopTime
	}
	if p.InviteAdmin != nil {
		body["inviteAdmin"] = *p.InviteAdmin
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillVCDetail(res, data)
	}
	return res, nil
}

// ModifyMeeting modifies a meeting that has not started (/meeting/modify).
// Members require exactly one role='admin'.
func (c *LansengerClient) ModifyMeeting(ctx context.Context, p *VideoconferenceModifyParams) (*VideoconferenceOpResult, error) {
	if p == nil {
		return &VideoconferenceOpResult{Success: false, Error: "params is required"}, nil
	}
	if err := vcHostErr(p.Members); err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	mid, err := vcMidValue(p.Mid)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_modify", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"orgId":           vcOrgIDValue(p.OrgID),
		"operator":        p.Operator,
		"mid":             mid,
		"subject":         p.Subject,
		"startTime":       p.StartTime,
		"autoRecord":      p.AutoRecord,
		"type":            p.Type,
		"groupNew":        p.GroupNew,
		"confPassword":    p.ConfPassword,
		"controlPassword": p.ControlPassword,
		"member":          p.Members,
	}
	// Same contract as CreateMeeting: when the field is omitted the server resets the
	// meeting to startTime + 24h, so only send it when the caller gave a value.
	if p.UserStopTime != nil {
		body["userStopTime"] = *p.UserStopTime
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceOpResult{Success: true, RawResponse: result}
	fillVCOp(res, result)
	return res, nil
}

// CancelMeeting cancels a meeting that has not started (/meeting/cancle —
// the endpoint keeps the server's historical "cancle" spelling).
func (c *LansengerClient) CancelMeeting(ctx context.Context, mid, orgID, operator, userToken string) (*VideoconferenceOpResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_cancel", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"mid":      midVal,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceOpResult{Success: true, RawResponse: result}
	fillVCOp(res, result)
	return res, nil
}

// StopMeeting ends a running meeting (/meeting/stop).
func (c *LansengerClient) StopMeeting(ctx context.Context, mid, orgID, operator, userToken string) (*VideoconferenceOpResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_stop", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"mid":      midVal,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceOpResult{Success: true, RawResponse: result}
	fillVCOp(res, result)
	return res, nil
}

// FetchMeetingDetail fetches meeting detail by mid (/meeting/detail).
func (c *LansengerClient) FetchMeetingDetail(ctx context.Context, mid, orgID, operator, userToken string) (*VideoconferenceDetailResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceDetailResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_detail", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"mid":      midVal,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillVCDetail(res, data)
	}
	return res, nil
}

// FetchMeetingList fetches the meeting list by time range (/meeting/list).
// fetchRange: my | all | person (person requires staffID).
func (c *LansengerClient) FetchMeetingList(ctx context.Context, orgID string, startTime, endTime int64, fetchRange, staffID string, limit, offset int, userToken string) (*VideoconferenceListResult, error) {
	if fetchRange == VCFetchRangePerson && staffID == "" {
		return &VideoconferenceListResult{Success: false, Error: "staff_id is required when fetch_range='person'"}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_list", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":      vcOrgIDValue(orgID),
		"limit":      limit,
		"offset":     offset,
		"startTime":  startTime,
		"endTime":    endTime,
		"fetchRange": fetchRange,
		"staffId":    staffID,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// FetchMeetingRecordList fetches the meeting operation record list
// (/meeting/record/list). createSource: 0=platform client, 1=third-party app.
func (c *LansengerClient) FetchMeetingRecordList(ctx context.Context, orgID string, startTime, endTime int64, admin string, createSource, limit, offset int, userToken string) (*VideoconferenceListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "meeting_record_list", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":        vcOrgIDValue(orgID),
		"limit":        limit,
		"offset":       offset,
		"startTime":    startTime,
		"endTime":      endTime,
		"admin":        admin,
		"createSource": createSource,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// FetchMemberSimpleRecord fetches member join/leave records of a meeting
// (/meeting/member/simplerecord).
func (c *LansengerClient) FetchMemberSimpleRecord(ctx context.Context, mid, orgID, operator, userToken string, limit, offset int) (*VideoconferenceListResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "member_simplerecord", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"mid":      midVal,
		"operator": operator,
		"limit":    limit,
		"offset":   offset,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// FetchFixroomList fetches the fixed (cloud) meeting-room list
// (/meeting/fixroom/list).
func (c *LansengerClient) FetchFixroomList(ctx context.Context, orgID, operator, userToken string, limit, offset int) (*VideoconferenceListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "fixroom_list", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"limit":    limit,
		"offset":   offset,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// FetchMeetingStatus batch-fetches meeting status by mids
// (/meeting/status/fetchmore).
func (c *LansengerClient) FetchMeetingStatus(ctx context.Context, mids []string, orgID, userToken string) (*VideoconferenceStatusListResult, error) {
	if len(mids) == 0 {
		return &VideoconferenceStatusListResult{Success: false, Error: "mids is required"}, nil
	}
	midVals := make([]int64, 0, len(mids))
	for _, mid := range mids {
		n, err := vcMidValue(mid)
		if err != nil {
			return &VideoconferenceStatusListResult{Success: false, Error: err.Error()}, nil
		}
		midVals = append(midVals, n)
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "status_fetchmore", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId": vcOrgIDValue(orgID),
		"mids":  midVals,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceStatusListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceStatusListResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		if list, ok := data["mids"].([]interface{}); ok {
			res.Statuses = list
		}
	}
	return res, nil
}

// SubscribeMeetingEvents subscribes meeting status-change events
// (/meeting/events/subscribe).
func (c *LansengerClient) SubscribeMeetingEvents(ctx context.Context, mid, orgID string, events []map[string]interface{}, callBackInfo, userToken string) (*VideoconferenceOpResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "events_subscribe", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":  vcOrgIDValue(orgID),
		"mid":    midVal,
		"events": events,
	}
	if callBackInfo != "" {
		body["callBackInfo"] = callBackInfo
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceOpResult{Success: true, RawResponse: result}
	fillVCOp(res, result)
	return res, nil
}

// FetchMeetingParams fetches meeting params by meetingNumber; PRS >=3.8
// (/meeting/param/fetch).
func (c *LansengerClient) FetchMeetingParams(ctx context.Context, meetingNumber, orgID, operator, userToken string) (*VideoconferenceParamResult, error) {
	if meetingNumber == "" {
		return &VideoconferenceParamResult{Success: false, Error: "meeting_number is required"}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "param_fetch", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":         vcOrgIDValue(orgID),
		"meetingNumber": meetingNumber,
		"operator":      operator,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceParamResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceParamResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		if info, ok := data["meetingInfo"].(map[string]interface{}); ok {
			res.MeetingInfo = info
		}
	}
	return res, nil
}

// FetchHistoryMeetings fetches a person's past meetings
// (/meeting/history/fetch).
func (c *LansengerClient) FetchHistoryMeetings(ctx context.Context, orgID, operator, userToken string, limit, offset int) (*VideoconferenceListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "history_fetch", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"limit":    limit,
		"offset":   offset,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// FetchActiveMeetings fetches a person's running + reserved meetings
// (/meeting/active/fetch).
func (c *LansengerClient) FetchActiveMeetings(ctx context.Context, orgID, operator, userToken string, limit, offset int) (*VideoconferenceListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "active_fetch", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"limit":    limit,
		"offset":   offset,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// ControlMember performs a host control operation on a member
// (/meeting/member/control). opCode is forwarded verbatim — the server is the
// authority on accepted values; VCOpCodes is a reference list, not a validator.
// Some values (e.g. muteall/unmuteall) may be rejected with errCode=105601
// depending on the meeting server build.
func (c *LansengerClient) ControlMember(ctx context.Context, mid, staffID, opCode, operator, orgID, userToken string) (*VideoconferenceOpResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "member_control", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"staffId":  staffID,
		"mid":      midVal,
		"opCode":   opCode,
		"operator": operator,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceOpResult{Success: true, RawResponse: result}
	fillVCOp(res, result)
	return res, nil
}

// InviteMembers invites members to a running meeting (/meeting/member/invite).
func (c *LansengerClient) InviteMembers(ctx context.Context, meetingNumber string, members []VideoconferenceMember, orgID, operator, userToken string) (*VideoconferenceOpResult, error) {
	if len(members) == 0 {
		return &VideoconferenceOpResult{Success: false, Error: "member is required"}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "member_invite", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":         vcOrgIDValue(orgID),
		"operator":      operator,
		"meetingNumber": meetingNumber,
		"member":        members,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceOpResult{Success: true, RawResponse: result}
	fillVCOp(res, result)
	return res, nil
}

// FetchMemberList fetches the paged member list of a meeting
// (/meeting/member/list). The endpoint is documented as GET with a JSON
// body; this gateway's convention is POST, so we send POST like the rest
// of the module.
func (c *LansengerClient) FetchMemberList(ctx context.Context, mid, orgID, operator, userToken string, limit, offset int) (*VideoconferenceListResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "member_list", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"mid":      midVal,
		"operator": operator,
		"limit":    limit,
		"offset":   offset,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceListResult{Success: true, RawResponse: result}
	fillVCPage(res, result)
	return res, nil
}

// FetchVodList fetches the recording list of a meeting (/meeting/vod/list).
func (c *LansengerClient) FetchVodList(ctx context.Context, mid, orgID, operator, userToken string) (*VideoconferenceVodListResult, error) {
	midVal, err := vcMidValue(mid)
	if err != nil {
		return &VideoconferenceVodListResult{Success: false, Error: err.Error()}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "vod_list", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"mid":      midVal,
		"operator": operator,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceVodListResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceVodListResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		if list, ok := data["items"].([]interface{}); ok {
			res.Items = list
		}
	}
	return res, nil
}

// FetchVodDownloadURLs fetches recording download URLs
// (/vod/url/download/fetch); max 3 vods per call.
func (c *LansengerClient) FetchVodDownloadURLs(ctx context.Context, vods []VideoconferenceVod, orgID, operator, userToken string) (*VideoconferenceVodUrlResult, error) {
	if len(vods) == 0 || len(vods) > 3 {
		return &VideoconferenceVodUrlResult{Success: false, Error: "vods must contain 1..3 entries"}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "vod_download_url", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    vcOrgIDValue(orgID),
		"operator": operator,
		"vods":     vods,
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceVodUrlResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceVodUrlResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.Data = data
	}
	return res, nil
}

// FetchOrgConf fetches the org videoconference config; PRS >=3.8
// (/conf/fetch). Pass meetingNumber to get the config of the org owning
// that meeting.
func (c *LansengerClient) FetchOrgConf(ctx context.Context, orgID, meetingNumber, operator, userToken string) (*VideoconferenceConfResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "videoconferences", "conf_fetch", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId": vcOrgIDValue(orgID),
	}
	if meetingNumber != "" {
		body["meetingNumber"] = meetingNumber
	}
	if operator != "" {
		body["operator"] = operator
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &VideoconferenceConfResult{Success: false, Error: err.Error()}, nil
	}
	res := &VideoconferenceConfResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.MaxPerson = intFromMap(data, "maxPerson")
		res.DefaultMaxPerson = intFromMap(data, "defaultMaxPerson")
		res.AllowedRecordFlag = intFromMap(data, "allowedRecordFlag")
		res.ForcePasswdFlag = intFromMap(data, "forcePasswdFlag")
		res.SpaceSize = int64FromMap(data, "spaceSize")
	}
	return res, nil
}
