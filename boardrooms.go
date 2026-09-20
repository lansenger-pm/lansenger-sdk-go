package lansenger

import (
	"context"
)

// BoardroomReserveParams carries the fields for ReserveBoardroom (会议室预定 V2 /v2/reserveRoom).
type BoardroomReserveParams struct {
	BoardRoomID        string
	Name               string
	GradingID          string
	ReserveTimeStart   string // yyyy-MM-dd HH:mm:ss
	ReserveTimeEnd     string
	NoticeTime         string
	ReserveUser        string
	OrgID              string
	Toastmaster        string
	Leader             string
	LeaderAttend       string // 0出席 1不出席
	PeopleNumber       string
	OtherDemand        string
	IsVideo            string // 0开启 1不开启
	VideoName          string
	UserList           []string
	InvitationUserList []string
	TableCards         string // 0开启 1关闭
	ReserveType        string // 0单次 1重复
	RepeatType         string // day/week/month
	RepeatDays         []int
	Skip               string // 0跳过 1不跳过
	RepeatEndDateStr   string
	UserToken          string
}

// BoardroomEditParams carries the fields for EditBoardroomReserve (会议室预定 V2 /v2/editReserve).
type BoardroomEditParams struct {
	ReserveID        string
	BoardRoomID      string
	Name             string
	GradingID        string
	ReserveTimeStart string
	ReserveTimeEnd   string
	NoticeTime       string
	EditType         string // 1当前 2当前及以后
	ReserveUser      string
	OrgID            string
	PeopleNumber     string
	UserToken        string
}

func fillBoardroomPage(res *BoardroomListResult, result map[string]interface{}) {
	data := extractData(result)
	if data == nil {
		return
	}
	res.Count = intFromMap(data, "count")
	if list, ok := data["data"].([]interface{}); ok {
		res.Items = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Items = append(res.Items, m)
			}
		}
	}
}

func fillBoardroomDetail(res *BoardroomDetailResult, data map[string]interface{}) {
	res.RoomID = strFromMap(data, "id")
	res.Name = strFromMap(data, "name")
	res.Status = strFromMap(data, "status")
	res.PeopleNum = intFromMap(data, "peopleNum")
	res.CanReserveFlag = strFromMap(data, "canReserveFlag")
	res.Address = strFromMap(data, "address")
	res.AreaName = strFromMap(data, "areaName")
	res.GradingID = strFromMap(data, "gradingId")
}

func fillBoardroomSchedule(res *BoardroomScheduleResult, data map[string]interface{}) {
	res.RoomID = strFromMap(data, "id")
	res.Name = strFromMap(data, "name")
	res.PeopleNum = intFromMap(data, "peopleNum")
	res.CanReserveFlag = strFromMap(data, "canReserveFlag")
	if list, ok := data["reserveDtoList"].([]interface{}); ok {
		res.Reserves = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Reserves = append(res.Reserves, m)
			}
		}
	}
	if list, ok := data["deactivatedInfoList"].([]interface{}); ok {
		res.Deactivations = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Deactivations = append(res.Deactivations, m)
			}
		}
	}
}

func fillBoardroomReserveDetail(res *BoardroomReserveDetailResult, data map[string]interface{}) {
	res.ReserveID = strFromMap(data, "id")
	res.BoardroomName = strFromMap(data, "boardRoomName")
	res.MeetingName = strFromMap(data, "name")
	res.Status = strFromMap(data, "status")
	res.ReserveTimeStart = strFromMap(data, "reserveTimeStart")
	res.ReserveTimeEnd = strFromMap(data, "reserveTimeEnd")
	res.ReserveTime = strFromMap(data, "reserveTime")
	res.ReserveUserName = strFromMap(data, "reserveUserName")
	res.PeopleNumber = strFromMap(data, "peopleNumber")
}

func fillBoardroomReserve(res *BoardroomReserveResult, data map[string]interface{}) {
	res.ReserveID = strFromMap(data, "id")
	res.ReserveCode = strFromMap(data, "reserveCode")
	res.BoardroomName = strFromMap(data, "boardRoomName")
	res.MeetingName = strFromMap(data, "name")
	res.Status = strFromMap(data, "status")
	res.ReserveTimeStart = strFromMap(data, "reserveTimeStart")
	res.ReserveTimeEnd = strFromMap(data, "reserveTimeEnd")
	res.ReserveTime = strFromMap(data, "reserveTime")
}

// FetchBoardroomList filters meeting rooms by area/floor/equipment/time/capacity.
func (c *LansengerClient) FetchBoardroomList(ctx context.Context, gradingID, areaOfficeID string, floorIDs, equipment []string, reserveTimeStart, reserveTimeEnd, queryDate string, page, limit int, lxUserID, orgID, userToken string) (*BoardroomListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "room_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"page": page, "limit": limit}
	if gradingID != "" {
		body["gradingId"] = gradingID
	}
	if areaOfficeID != "" {
		body["areaOfficeId"] = areaOfficeID
	}
	if len(floorIDs) > 0 {
		body["areaOfficeFoolerIds"] = floorIDs
	}
	if len(equipment) > 0 {
		body["equipment"] = equipment
	}
	if reserveTimeStart != "" {
		body["reserveTimeStartStr"] = reserveTimeStart
	}
	if reserveTimeEnd != "" {
		body["reserveTimeEndStr"] = reserveTimeEnd
	}
	if queryDate != "" {
		body["queryDate"] = queryDate
	}
	if lxUserID != "" {
		body["lxUserId"] = lxUserID
	}
	if orgID != "" {
		body["orgId"] = orgID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomListResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomListResult{Success: true, RawResponse: result}
	fillBoardroomPage(res, result)
	return res, nil
}

// FetchBoardroomDetail fetches meeting-room detail.
func (c *LansengerClient) FetchBoardroomDetail(ctx context.Context, roomID, orgID, userToken string) (*BoardroomDetailResult, error) {
	if roomID == "" {
		return &BoardroomDetailResult{Success: false, Error: "room_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "room_detail", token, WithUserToken(userToken))
	body := map[string]interface{}{"id": roomID}
	if orgID != "" {
		body["orgId"] = orgID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillBoardroomDetail(res, data)
	}
	return res, nil
}

// FetchBoardroomSchedule fetches a room's bookings + deactivations for a date.
func (c *LansengerClient) FetchBoardroomSchedule(ctx context.Context, roomID, queryDate, gradingID, reserveUserID, orgID, userToken string) (*BoardroomScheduleResult, error) {
	if roomID == "" {
		return &BoardroomScheduleResult{Success: false, Error: "room_id is required"}, nil
	}
	if queryDate == "" {
		return &BoardroomScheduleResult{Success: false, Error: "query_date is required"}, nil
	}
	if gradingID == "" {
		return &BoardroomScheduleResult{Success: false, Error: "grading_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "room_schedule", token, WithUserToken(userToken))
	body := map[string]interface{}{"roomId": roomID, "queryDate": queryDate, "gradingId": gradingID}
	if reserveUserID != "" {
		body["reserveUserId"] = reserveUserID
	}
	if orgID != "" {
		body["orgId"] = orgID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomScheduleResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomScheduleResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillBoardroomSchedule(res, data)
	}
	return res, nil
}

// FetchBoardroomReserveDetail fetches reservation detail (attendees, approval flow).
func (c *LansengerClient) FetchBoardroomReserveDetail(ctx context.Context, reserveRoomID, gradingID, orgID, userToken string) (*BoardroomReserveDetailResult, error) {
	if reserveRoomID == "" {
		return &BoardroomReserveDetailResult{Success: false, Error: "reserve_room_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "reserve_detail", token, WithUserToken(userToken))
	body := map[string]interface{}{"reserveRoomId": reserveRoomID}
	if gradingID != "" {
		body["gradingId"] = gradingID
	}
	if orgID != "" {
		body["orgId"] = orgID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomReserveDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomReserveDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillBoardroomReserveDetail(res, data)
	}
	return res, nil
}

// ReserveBoardroom reserves a meeting room (single or repeating).
func (c *LansengerClient) ReserveBoardroom(ctx context.Context, p *BoardroomReserveParams) (*BoardroomReserveResult, error) {
	if p == nil {
		return &BoardroomReserveResult{Success: false, Error: "params is required"}, nil
	}
	if p.BoardRoomID == "" {
		return &BoardroomReserveResult{Success: false, Error: "boardroom_id is required"}, nil
	}
	if p.Name == "" {
		return &BoardroomReserveResult{Success: false, Error: "name is required"}, nil
	}
	if p.GradingID == "" {
		return &BoardroomReserveResult{Success: false, Error: "grading_id is required"}, nil
	}
	if p.ReserveTimeStart == "" {
		return &BoardroomReserveResult{Success: false, Error: "reserve_time_start is required"}, nil
	}
	if p.ReserveTimeEnd == "" {
		return &BoardroomReserveResult{Success: false, Error: "reserve_time_end is required"}, nil
	}
	if p.NoticeTime == "" {
		return &BoardroomReserveResult{Success: false, Error: "notice_time is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "reserve_room", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"boardRoomId":         p.BoardRoomID,
		"name":                p.Name,
		"gradingId":           p.GradingID,
		"reserveTimeStartStr": p.ReserveTimeStart,
		"reserveTimeEndStr":   p.ReserveTimeEnd,
		"noticeTime":          p.NoticeTime,
	}
	if p.ReserveUser != "" {
		body["reserveUser"] = p.ReserveUser
	}
	if p.OrgID != "" {
		body["orgId"] = p.OrgID
	}
	if p.Toastmaster != "" {
		body["toastmaster"] = p.Toastmaster
	}
	if p.Leader != "" {
		body["leader"] = p.Leader
	}
	if p.LeaderAttend != "" {
		body["leaderAttend"] = p.LeaderAttend
	}
	if p.PeopleNumber != "" {
		body["peopleNumber"] = p.PeopleNumber
	}
	if p.OtherDemand != "" {
		body["otherDemand"] = p.OtherDemand
	}
	if p.IsVideo != "" {
		body["isVideo"] = p.IsVideo
	}
	if p.VideoName != "" {
		body["videoName"] = p.VideoName
	}
	if len(p.UserList) > 0 {
		body["userList"] = p.UserList
	}
	if len(p.InvitationUserList) > 0 {
		body["invitationUserList"] = p.InvitationUserList
	}
	if p.TableCards != "" {
		body["tableCards"] = p.TableCards
	}
	if p.ReserveType != "" {
		body["reserveType"] = p.ReserveType
	}
	if p.RepeatType != "" {
		body["repeatType"] = p.RepeatType
	}
	if len(p.RepeatDays) > 0 {
		body["repeatDays"] = p.RepeatDays
	}
	if p.Skip != "" {
		body["skip"] = p.Skip
	}
	if p.RepeatEndDateStr != "" {
		body["repeatEndDateStr"] = p.RepeatEndDateStr
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomReserveResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomReserveResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillBoardroomReserve(res, data)
	}
	return res, nil
}

// EditBoardroomReserve edits a reservation (only non-approval-flow bookings).
func (c *LansengerClient) EditBoardroomReserve(ctx context.Context, p *BoardroomEditParams) (*BoardroomReserveResult, error) {
	if p == nil {
		return &BoardroomReserveResult{Success: false, Error: "params is required"}, nil
	}
	if p.ReserveID == "" {
		return &BoardroomReserveResult{Success: false, Error: "reserve_id is required"}, nil
	}
	if p.BoardRoomID == "" {
		return &BoardroomReserveResult{Success: false, Error: "boardroom_id is required"}, nil
	}
	if p.Name == "" {
		return &BoardroomReserveResult{Success: false, Error: "name is required"}, nil
	}
	if p.GradingID == "" {
		return &BoardroomReserveResult{Success: false, Error: "grading_id is required"}, nil
	}
	if p.ReserveTimeStart == "" {
		return &BoardroomReserveResult{Success: false, Error: "reserve_time_start is required"}, nil
	}
	if p.ReserveTimeEnd == "" {
		return &BoardroomReserveResult{Success: false, Error: "reserve_time_end is required"}, nil
	}
	if p.NoticeTime == "" {
		return &BoardroomReserveResult{Success: false, Error: "notice_time is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "edit_reserve", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"id":                  p.ReserveID,
		"boardRoomId":         p.BoardRoomID,
		"name":                p.Name,
		"gradingId":           p.GradingID,
		"reserveTimeStartStr": p.ReserveTimeStart,
		"reserveTimeEndStr":   p.ReserveTimeEnd,
		"noticeTime":          p.NoticeTime,
	}
	if p.EditType != "" {
		body["editType"] = p.EditType
	}
	if p.ReserveUser != "" {
		body["reserveUser"] = p.ReserveUser
	}
	if p.OrgID != "" {
		body["orgId"] = p.OrgID
	}
	if p.PeopleNumber != "" {
		body["peopleNumber"] = p.PeopleNumber
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomReserveResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomReserveResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillBoardroomReserve(res, data)
	}
	return res, nil
}

// CancelBoardroomReserve cancels a reservation (status 0/1/5 only).
func (c *LansengerClient) CancelBoardroomReserve(ctx context.Context, reserveID, cancelUserID, orgID, cancelReason string, isSend *bool, notifyUserList []string, cancelVideo, cancelType, userToken string) (*BoardroomOpResult, error) {
	if reserveID == "" {
		return &BoardroomOpResult{Success: false, Error: "reserve_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "reserve_cancel", token, WithUserToken(userToken))
	body := map[string]interface{}{"id": reserveID}
	if cancelUserID != "" {
		body["cancelUserId"] = cancelUserID
	}
	if orgID != "" {
		body["orgId"] = orgID
	}
	if cancelReason != "" {
		body["cancelReason"] = cancelReason
	}
	if isSend != nil {
		body["isSend"] = *isSend
	}
	if len(notifyUserList) > 0 {
		body["userList"] = notifyUserList
	}
	if cancelVideo != "" {
		body["cancelVideo"] = cancelVideo
	}
	if cancelType != "" {
		body["cancelType"] = cancelType
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomOpResult{Success: true, RawResponse: result}
	res.Done = boolDataValue(result)
	return res, nil
}

// ConfirmBoardroomSign confirms a reservation awaiting scan-code (status 1 only).
func (c *LansengerClient) ConfirmBoardroomSign(ctx context.Context, reserveID, orgID, userToken string) (*BoardroomOpResult, error) {
	if reserveID == "" {
		return &BoardroomOpResult{Success: false, Error: "reserve_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "confirm_sign", token, WithUserToken(userToken))
	body := map[string]interface{}{"id": reserveID}
	if orgID != "" {
		body["orgId"] = orgID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomOpResult{Success: true, RawResponse: result}
	res.Done = boolDataValue(result)
	return res, nil
}

// FetchMyBoardroomReserves pages the user's reservations with filters.
func (c *LansengerClient) FetchMyBoardroomReserves(ctx context.Context, gradingID, keys, startTime, endTime, boardRoomID string, floorIDs []string, page, limit int, lxUserID, orgID, userToken string) (*BoardroomListResult, error) {
	if gradingID == "" {
		return &BoardroomListResult{Success: false, Error: "grading_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "my_reserve_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"gradingId": gradingID, "page": page, "limit": limit}
	if keys != "" {
		body["keys"] = keys
	}
	if startTime != "" {
		body["startTime"] = startTime
	}
	if endTime != "" {
		body["endTime"] = endTime
	}
	if boardRoomID != "" {
		body["boardRoomId"] = boardRoomID
	}
	if len(floorIDs) > 0 {
		body["areaOfficeFoolerIds"] = floorIDs
	}
	if lxUserID != "" {
		body["lxUserId"] = lxUserID
	}
	if orgID != "" {
		body["orgId"] = orgID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomListResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomListResult{Success: true, RawResponse: result}
	fillBoardroomPage(res, result)
	return res, nil
}

// FetchBoardroomGradings fetches gradings visible to the user.
func (c *LansengerClient) FetchBoardroomGradings(ctx context.Context, lxUserID, orgID, userToken string) (*BoardroomGradingListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "grading_list", token, WithUserToken(userToken))
	body := map[string]interface{}{}
	if lxUserID != "" {
		body["lxUserId"] = lxUserID
	}
	if orgID != "" {
		body["orgId"] = orgID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BoardroomGradingListResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomGradingListResult{Success: true, RawResponse: result}
	if list := extractDataArray(result); list != nil {
		res.Gradings = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Gradings = append(res.Gradings, m)
			}
		}
		res.Total = len(res.Gradings)
	}
	return res, nil
}

// FetchBoardroomAreaOffices fetches office areas under a grading.
func (c *LansengerClient) FetchBoardroomAreaOffices(ctx context.Context, gradingID, userToken string) (*BoardroomAreaListResult, error) {
	if gradingID == "" {
		return &BoardroomAreaListResult{Success: false, Error: "grading_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "boardrooms", "area_office_list", token, WithUserToken(userToken))
	result, err := c.doPost(ctx, url, map[string]interface{}{"gradingId": gradingID})
	if err != nil {
		return &BoardroomAreaListResult{Success: false, Error: err.Error()}, nil
	}
	res := &BoardroomAreaListResult{Success: true, RawResponse: result}
	if list := extractDataArray(result); list != nil {
		res.Areas = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Areas = append(res.Areas, m)
			}
		}
		res.Total = len(res.Areas)
	}
	return res, nil
}
