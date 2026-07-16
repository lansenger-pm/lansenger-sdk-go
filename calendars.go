package lansenger

import (
	"context"
)

func (c *LansengerClient) FetchPrimaryCalendar(ctx context.Context, userToken, userID string) (*CalendarPrimaryResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "primary_fetch", token,
		WithUserToken(userToken),
		WithUserID(userID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &CalendarPrimaryResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &CalendarPrimaryResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &CalendarPrimaryResult{
		Success:     true,
		CalendarID:  strFromMap(data, "calendarId"),
		Summary:     strFromMap(data, "summary"),
		Description: strFromMap(data, "description"),
		Permissions: strFromMap(data, "permissions"),
		Color:       strFromMap(data, "color"),
		Type:        strFromMap(data, "type"),
		Role:        strFromMap(data, "role"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) CreateSchedule(ctx context.Context, calendarID, summary string, startTime, endTime map[string]interface{}, attendees []map[string]interface{}, description string, allDay string, repeatType string, rule map[string]interface{}, expireDateType, reminderType, attendeePermissions string, userToken, userID string) (*ScheduleCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_create", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
	)

	body := map[string]interface{}{
		"summary":   summary,
		"startTime": startTime,
		"endTime":   endTime,
	}
	if len(attendees) == 0 {
		if userID == "" {
			return &ScheduleCreateResult{Success: false, Error: "attendees is required (or provide user_id to auto-fill creator)"}, nil
		}
		attendees = []map[string]interface{}{{"staffId": userID, "attendeeFlag": "required"}}
	}
	body["attendees"] = attendees
	if description != "" {
		body["description"] = description
	}
	if allDay != "" {
		body["allDay"] = allDay
	}
	if repeatType != "" {
		body["repeatType"] = repeatType
	}
	if rule != nil {
		body["rule"] = rule
	}
	if expireDateType != "" {
		body["expireDateType"] = expireDateType
	}
	if reminderType != "" {
		body["reminderType"] = reminderType
	}
	if attendeePermissions != "" {
		body["attendeePermissions"] = attendeePermissions
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &ScheduleCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &ScheduleCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.ScheduleID = strFromMap(data, "scheduleId")
	}
	return res, nil
}

func (c *LansengerClient) FetchSchedule(ctx context.Context, calendarID, scheduleID, userToken, userID string) (*ScheduleInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_fetch", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &ScheduleInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &ScheduleInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &ScheduleInfoResult{
		Success:          true,
		ScheduleID:       strFromMap(data, "scheduleId"),
		Summary:          strFromMap(data, "summary"),
		Description:      strFromMap(data, "description"),
		RepeatType:       strFromMap(data, "repeatType"),
		AllDay:           strFromMap(data, "allDay"),
		StartTime:        mapFromMap(data, "startTime"),
		EndTime:          mapFromMap(data, "endTime"),
		Creator:          mapFromMap(data, "creator"),
		RsvpStatus:       strFromMap(data, "rsvpStatus"),
		PrimaryScheduleID: strFromMap(data, "primaryScheduleId"),
		ExpireDateType:   strFromMap(data, "expireDateType"),
		AttendeePermissions: strFromMap(data, "attendeePermissions"),
		Color:            strFromMap(data, "color"),
		RawResponse:      result,
	}, nil
}

func (c *LansengerClient) DeleteSchedule(ctx context.Context, calendarID, scheduleID string, reminderType string, operationType string, currentTime string, userToken, userID string) (*ScheduleDeleteResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_delete", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{
		"reminderType": reminderType,
	}
	if operationType != "" {
		body["operationType"] = operationType
	}
	if currentTime != "" {
		body["currentTime"] = currentTime
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &ScheduleDeleteResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &ScheduleDeleteResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.ScheduleIDs = stringArrayFromMap(data, "scheduleIds")
	}
	return res, nil
}

func (c *LansengerClient) FetchScheduleList(ctx context.Context, calendarID string, startTime, endTime int64, userToken, userID string) (*ScheduleListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_list_fetch", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
	)

	body := map[string]interface{}{
		"startTime": startTime,
		"endTime":   endTime,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &ScheduleListResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &ScheduleListResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		if list, ok := data["scheduleList"].([]interface{}); ok {
			res.ScheduleList = make([]map[string]interface{}, 0, len(list))
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					res.ScheduleList = append(res.ScheduleList, m)
				}
			}
		}
	}
	return res, nil
}

func (c *LansengerClient) FetchScheduleAttendees(ctx context.Context, calendarID, scheduleID string, page, pageSize int, userToken, userID string) (*ScheduleAttendeesResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_fetch", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
		WithPage(page),
		WithPageSize(pageSize),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &ScheduleAttendeesResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &ScheduleAttendeesResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.Total = intFromMap(data, "total")
		if attendees, ok := data["attendees"].([]interface{}); ok {
			res.Attendees = make([]map[string]interface{}, 0, len(attendees))
			for _, item := range attendees {
				if m, ok := item.(map[string]interface{}); ok {
					res.Attendees = append(res.Attendees, m)
				}
			}
		}
	}
	return res, nil
}

func (c *LansengerClient) AddScheduleAttendees(ctx context.Context, calendarID, scheduleID string, attendees []string, reminderType string, operationType, currentTime string, userToken, userID string) (*ScheduleCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_create", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{
		"attendees": attendees,
	}
	if reminderType != "" {
		body["reminderType"] = reminderType
	}
	if operationType != "" {
		body["operationType"] = operationType
	}
	if currentTime != "" {
		body["currentTime"] = currentTime
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &ScheduleCreateResult{Success: false, Error: err.Error()}, nil
	}

	return &ScheduleCreateResult{
		Success:     true,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) DeleteScheduleAttendees(ctx context.Context, calendarID, scheduleID string, attendees []string, reminderType string, operationType, currentTime string, userToken, userID string) (*ScheduleAttendeesDeleteResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_delete", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{
		"attendees": attendees,
	}
	if reminderType != "" {
		body["reminderType"] = reminderType
	}
	if operationType != "" {
		body["operationType"] = operationType
	}
	if currentTime != "" {
		body["currentTime"] = currentTime
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &ScheduleAttendeesDeleteResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &ScheduleAttendeesDeleteResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.ScheduleIDs = stringArrayFromMap(data, "scheduleIds")
	}
	return res, nil
}

func (c *LansengerClient) UpdateSchedule(ctx context.Context, calendarID, scheduleID string, params map[string]interface{}, userToken, userID string) (*ScheduleDeleteResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_update", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	result, err := c.doPost(ctx, url, params)
	if err != nil {
		return &ScheduleDeleteResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &ScheduleDeleteResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.ScheduleIDs = stringArrayFromMap(data, "scheduleIds")
	}
	return res, nil
}

func (c *LansengerClient) UpdateScheduleAttendeeMeta(ctx context.Context, calendarID, scheduleID string, params map[string]interface{}, userToken, userID string) (*SendMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_meta_update", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	result, err := c.doPost(ctx, url, params)
	if err != nil {
		return &SendMessageResult{Success: false, Error: err.Error(), Platform: "lansenger"}, nil
	}

	return &SendMessageResult{
		Success:     true,
		Platform:    "lansenger",
		Operation:   "update_schedule_attendee_meta",
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) UpdateScheduleAttendees(ctx context.Context, calendarID, scheduleID string, addAttendees, deleteAttendees []string, reminderType, operationType string, currentTime int, userToken, userID string) (*ScheduleAttendeesUpdateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_update", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{}
	if len(addAttendees) > 0 {
		body["addAttendees"] = addAttendees
	}
	if len(deleteAttendees) > 0 {
		body["deleteAttendees"] = deleteAttendees
	}
	if reminderType != "" {
		body["reminderType"] = reminderType
	}
	if operationType != "" {
		body["operationType"] = operationType
	}
	if currentTime != 0 {
		body["currentTime"] = currentTime
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &ScheduleAttendeesUpdateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	res := &ScheduleAttendeesUpdateResult{
		Success:     true,
		RawResponse: result,
	}

	if scheduleIDs, ok := data["scheduleIds"].([]interface{}); ok {
		for _, id := range scheduleIDs {
			if s, ok := id.(string); ok {
				res.ScheduleIDs = append(res.ScheduleIDs, s)
			}
		}
	}
	if attendees, ok := data["attendees"].([]interface{}); ok {
		for _, a := range attendees {
			if s, ok := a.(string); ok {
				res.FailedAttendees = append(res.FailedAttendees, s)
			}
		}
	}

	return res, nil
}