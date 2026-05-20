package lansenger

import (
	"context"
	"fmt"
)

func (c *LansengerClient) FetchPrimaryCalendar(ctx context.Context, userToken, userID string) (*CalendarPrimaryResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_primary_calendar")
	}

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

func (c *LansengerClient) CreateSchedule(ctx context.Context, calendarID, summary, startTime, endTime string, attendees []map[string]interface{}, description string, allDay bool, repeatType string, rule map[string]interface{}, expireDateType, reminderType, attendeePermissions int, userToken string) (*ScheduleCreateResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for create_schedule")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_create", token,
		WithUserToken(userToken),
		WithPathVar("calendar_id", calendarID),
	)

	body := map[string]interface{}{
		"summary":   summary,
		"startTime": startTime,
		"endTime":   endTime,
	}
	if len(attendees) > 0 {
		body["attendees"] = attendees
	}
	if description != "" {
		body["description"] = description
	}
	if allDay {
		body["allDay"] = true
	}
	if repeatType != "" {
		body["repeatType"] = repeatType
	}
	if rule != nil {
		body["rule"] = rule
	}
	if expireDateType != 0 {
		body["expireDateType"] = expireDateType
	}
	if reminderType != 0 {
		body["reminderType"] = reminderType
	}
	if attendeePermissions != 0 {
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
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_schedule")
	}

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
		Success:     true,
		ScheduleID:  strFromMap(data, "scheduleId"),
		Summary:     strFromMap(data, "summary"),
		Description: strFromMap(data, "description"),
		RepeatType:  strFromMap(data, "repeatType"),
		AllDay:      boolFromMap(data, "allDay"),
		StartTime:   strFromMap(data, "startTime"),
		EndTime:     strFromMap(data, "endTime"),
		Creator:     strFromMap(data, "creator"),
		RsvpStatus:  strFromMap(data, "rsvpStatus"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) DeleteSchedule(ctx context.Context, calendarID, scheduleID string, reminderType int, operationType string, currentTime string, userToken string) (*ScheduleCreateResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for delete_schedule")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_delete", token,
		WithUserToken(userToken),
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{
		"reminder_type": reminderType,
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

func (c *LansengerClient) FetchScheduleList(ctx context.Context, calendarID, startTime, endTime string, userToken string) (*ScheduleListResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_schedule_list")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_list_fetch", token,
		WithUserToken(userToken),
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

func (c *LansengerClient) FetchScheduleAttendees(ctx context.Context, calendarID, scheduleID string, page, pageSize int) (*ScheduleAttendeesResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_fetch", token,
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
	}
	return res, nil
}

func (c *LansengerClient) AddScheduleAttendees(ctx context.Context, calendarID, scheduleID string, attendees []map[string]interface{}, reminderType int) (*ScheduleCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_create", token,
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{
		"attendees": attendees,
	}
	if reminderType != 0 {
		body["reminderType"] = reminderType
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

func (c *LansengerClient) DeleteScheduleAttendees(ctx context.Context, calendarID, scheduleID string, attendees []map[string]interface{}, reminderType int) (*ScheduleCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "calendars", "schedules_members_delete", token,
		WithPathVar("calendar_id", calendarID),
		WithPathVar("schedule_id", scheduleID),
	)

	body := map[string]interface{}{
		"attendees": attendees,
	}
	if reminderType != 0 {
		body["reminderType"] = reminderType
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
