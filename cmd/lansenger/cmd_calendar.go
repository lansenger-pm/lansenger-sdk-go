package main

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var calendarCmd = &cobra.Command{
	Use:   "calendar",
	Short: "Calendar and schedule commands",
}

var calendarPrimaryCmd = &cobra.Command{
	Use:   "primary",
	Short: "Fetch primary calendar",
	Args:  cobra.NoArgs,
	Run:   runCalendarPrimary,
}

var calendarCreateScheduleCmd = &cobra.Command{
	Use:   "create-schedule CALENDAR_ID SUMMARY START_TIME END_TIME ATTENDEES",
	Short: "Create a schedule",
	Args:  cobra.ExactArgs(5),
	Run:   runCalendarCreateSchedule,
}

var calendarListSchedulesCmd = &cobra.Command{
	Use:   "list-schedules CALENDAR_ID START_TIME END_TIME",
	Short: "List schedules in a time range",
	Args:  cobra.ExactArgs(3),
	Run:   runCalendarListSchedules,
}

var calendarFetchScheduleCmd = &cobra.Command{
	Use:   "fetch-schedule CALENDAR_ID SCHEDULE_ID",
	Short: "Fetch schedule info",
	Args:  cobra.ExactArgs(2),
	Run:   runCalendarFetchSchedule,
}

var calendarDeleteScheduleCmd = &cobra.Command{
	Use:   "delete-schedule CALENDAR_ID SCHEDULE_ID",
	Short: "Delete a schedule",
	Args:  cobra.ExactArgs(2),
	Run:   runCalendarDeleteSchedule,
}

var calendarAttendeesCmd = &cobra.Command{
	Use:   "attendees CALENDAR_ID SCHEDULE_ID",
	Short: "Fetch schedule attendees",
	Args:  cobra.ExactArgs(2),
	Run:   runCalendarAttendees,
}

var calendarAddAttendeesCmd = &cobra.Command{
	Use:   "add-attendees CALENDAR_ID SCHEDULE_ID ATTENDEES",
	Short: "Add attendees to a schedule",
	Args:  cobra.ExactArgs(3),
	Run:   runCalendarAddAttendees,
}

var calendarDeleteAttendeesCmd = &cobra.Command{
	Use:   "delete-attendees CALENDAR_ID SCHEDULE_ID ATTENDEES",
	Short: "Delete attendees from a schedule",
	Args:  cobra.ExactArgs(3),
	Run:   runCalendarDeleteAttendees,
}

var calendarUpdateScheduleCmd = &cobra.Command{
	Use:   "update-schedule CALENDAR_ID SCHEDULE_ID",
	Short: "Update a schedule",
	Args:  cobra.ExactArgs(2),
	Run:   runCalendarUpdateSchedule,
}

var calendarAttendeeMetaCmd = &cobra.Command{
	Use:   "attendee-meta CALENDAR_ID SCHEDULE_ID",
	Short: "Update attendee metadata",
	Args:  cobra.ExactArgs(2),
	Run:   runCalendarAttendeeMeta,
}

var (
	calPrimaryUserToken string
	calPrimaryUserID    string

	calCreateDesc              string
	calCreateAllDay            string
	calCreateRepeatType        string
	calCreateReminderType      string
	calCreateAttendeePerms     string
	calCreateExpireDateType    string
	calCreateRule              string
	calCreateTz                string
	calCreateDate              string
	calCreateUserToken         string
	calCreateUserID            string

	calFetchUserToken string
	calFetchUserID    string

	calDeleteUserToken string
	calDeleteUserID    string

	calListUserToken string
	calListUserID    string

	calAttendeesUserToken string
	calAttendeesUserID    string
	calAttendeesPage      int
	calAttendeesSize      int

	calAddAttendeesReminderType string
	calAddAttendeesUserToken    string
	calAddAttendeesUserID       string

	calDelAttendeesReminderType string
	calDelAttendeesUserToken    string
	calDelAttendeesUserID       string

	calUpdateSummary       string
	calUpdateDesc          string
	calUpdateOp            string
	calUpdateCurrentTime   int
	calUpdateReminder      string
	calUpdateRepeat        string
	calUpdateRule          string
	calUpdateExpire        string
	calUpdateAllDay        string
	calUpdatePermissions   string
	calUpdateStartTime     string
	calUpdateEndTime       string
	calUpdateScheduleUserToken string
	calUpdateScheduleUserID    string

	calAttendeeMetaRsvp        string
	calAttendeeMetaColor       string
	calAttendeeMetaPermissions string
	calAttendeeMetaBusyFree    string
	calAttendeeMetaRemindTimes string
	calAttendeeMetaUserToken   string
	calAttendeeMetaUserID      string
)

func init() {
	calendarPrimaryCmd.Flags().StringVar(&calPrimaryUserToken, "user-token", "", "User token")
	calendarPrimaryCmd.Flags().StringVar(&calPrimaryUserID, "user-id", "", "User ID")

	calendarCreateScheduleCmd.Flags().StringVarP(&calCreateDesc, "desc", "d", "", "Schedule description")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateAllDay, "all-day", "no", "All-day event (yes/no)")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateRepeatType, "repeat", "no", "Repeat type: no/daily/weekly/monthly/yearly/work_day/custom")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateReminderType, "reminder", "yes", "Reminder type (yes/no)")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateAttendeePerms, "attendee-perms", "", "Attendee permissions")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateExpireDateType, "expire", "", "Expire date type: yes or no")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateRule, "rule", "", "Repeat rule (JSON)")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateTz, "tz", "Asia/Shanghai", "Time zone")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateDate, "date", "", "Date (for all-day events)")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateUserToken, "user-token", "", "User token")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateUserID, "user-id", "", "User ID")

	calendarFetchScheduleCmd.Flags().StringVar(&calFetchUserToken, "user-token", "", "User token")
	calendarFetchScheduleCmd.Flags().StringVar(&calFetchUserID, "user-id", "", "User ID")

	calendarDeleteScheduleCmd.Flags().StringVar(&calDeleteUserToken, "user-token", "", "User token")
	calendarDeleteScheduleCmd.Flags().StringVar(&calDeleteUserID, "user-id", "", "User ID")

	calendarListSchedulesCmd.Flags().StringVar(&calListUserToken, "user-token", "", "User token")
	calendarListSchedulesCmd.Flags().StringVar(&calListUserID, "user-id", "", "User ID")

	calendarAttendeesCmd.Flags().StringVar(&calAttendeesUserToken, "user-token", "", "User token")
	calendarAttendeesCmd.Flags().StringVar(&calAttendeesUserID, "user-id", "", "User ID")
	calendarAttendeesCmd.Flags().IntVarP(&calAttendeesPage, "page", "p", 1, "Page number")
	calendarAttendeesCmd.Flags().IntVarP(&calAttendeesSize, "size", "s", 500, "Page size")

	calendarAddAttendeesCmd.Flags().StringVar(&calAddAttendeesReminderType, "reminder", "", "Reminder type (yes/no)")
	calendarAddAttendeesCmd.Flags().StringVar(&calAddAttendeesUserToken, "user-token", "", "User token")
	calendarAddAttendeesCmd.Flags().StringVar(&calAddAttendeesUserID, "user-id", "", "User ID")

	calendarDeleteAttendeesCmd.Flags().StringVar(&calDelAttendeesReminderType, "reminder", "", "Reminder type (yes/no)")
	calendarDeleteAttendeesCmd.Flags().StringVar(&calDelAttendeesUserToken, "user-token", "", "User token")
	calendarDeleteAttendeesCmd.Flags().StringVar(&calDelAttendeesUserID, "user-id", "", "User ID")

	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateSummary, "summary", "", "New schedule summary")
	calendarUpdateScheduleCmd.Flags().StringVarP(&calUpdateDesc, "desc", "d", "", "New description")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateOp, "op", "modify_all", "Operation type: modify_all, modify_current, modify_current_after")
	calendarUpdateScheduleCmd.Flags().IntVar(&calUpdateCurrentTime, "current-time", 0, "Required when op != modify_all")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateReminder, "reminder", "", "Reminder type: yes or no")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateRepeat, "repeat", "", "Repeat type: no/daily/weekly/monthly/yearly/work_day/custom")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateRule, "rule", "", "RFC 5545 repeat rule")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateExpire, "expire", "", "Expire date type: yes or no")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateAllDay, "all-day", "", "All day: yes or no")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdatePermissions, "permissions", "", "Attendee permissions")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateStartTime, "start-time", "", "Start time as JSON dict")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateEndTime, "end-time", "", "End time as JSON dict")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateScheduleUserToken, "user-token", "", "User token")
	calendarUpdateScheduleCmd.Flags().StringVar(&calUpdateScheduleUserID, "user-id", "", "User ID")

	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaRsvp, "rsvp", "", "RSVP status: accept, tentative, decline")
	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaColor, "color", "", "Hex color (e.g. #FF347AFC)")
	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaPermissions, "permissions", "", "Visibility: private, public, default")
	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaBusyFree, "busy-free", "", "Busy/free state: busy, free")
	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaRemindTimes, "remind-times", "", "Reminder offsets in minutes as JSON list")
	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaUserToken, "user-token", "", "User token")
	calendarAttendeeMetaCmd.Flags().StringVar(&calAttendeeMetaUserID, "user-id", "", "User ID")

	calendarCmd.AddCommand(calendarPrimaryCmd)
	calendarCmd.AddCommand(calendarCreateScheduleCmd)
	calendarCmd.AddCommand(calendarFetchScheduleCmd)
	calendarCmd.AddCommand(calendarDeleteScheduleCmd)
	calendarCmd.AddCommand(calendarListSchedulesCmd)
	calendarCmd.AddCommand(calendarAttendeesCmd)
	calendarCmd.AddCommand(calendarAddAttendeesCmd)
	calendarCmd.AddCommand(calendarDeleteAttendeesCmd)
	calendarCmd.AddCommand(calendarUpdateScheduleCmd)
	calendarCmd.AddCommand(calendarAttendeeMetaCmd)
	rootCmd.AddCommand(calendarCmd)
}

func runCalendarPrimary(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchPrimaryCalendar(ctx, calPrimaryUserToken, calPrimaryUserID)
	checkError(err)
	outputResultFields(result, []string{"calendar_id", "summary", "description", "permissions", "color", "type", "role"})
}

func runCalendarCreateSchedule(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	calendarID := args[0]
	summary := args[1]

	var startTimeInt, endTimeInt int64
	fmt.Sscanf(args[2], "%d", &startTimeInt)
	fmt.Sscanf(args[3], "%d", &endTimeInt)

	startTime := map[string]interface{}{"time": startTimeInt, "timeZone": calCreateTz}
	if calCreateDate != "" {
		startTime["date"] = calCreateDate
	}

	endTime := map[string]interface{}{"time": endTimeInt, "timeZone": calCreateTz}

	attendeeMaps, err := parseJSONArray(args[4])
	checkError(err)

	var attendees []map[string]interface{}
	for _, m := range attendeeMaps {
		converted := make(map[string]interface{}, len(m))
		for k, v := range m {
			converted[k] = v
		}
		attendees = append(attendees, converted)
	}

	var rule map[string]interface{}
	if calCreateRule != "" {
		raw, rErr := parseJSONRaw(calCreateRule)
		checkError(rErr)
		if m, ok := raw.(map[string]interface{}); ok {
			rule = m
		}
	}

	result, err := client.CreateSchedule(ctx, calendarID, summary, startTime, endTime, attendees, calCreateDesc, calCreateAllDay, calCreateRepeatType, rule, calCreateExpireDateType, calCreateReminderType, calCreateAttendeePerms, calCreateUserToken, calCreateUserID)
	checkError(err)
	outputResultFields(result, []string{"schedule_id"})
}

func runCalendarFetchSchedule(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchSchedule(ctx, args[0], args[1], calFetchUserToken, calFetchUserID)
	checkError(err)
	outputResultFields(result, []string{"schedule_id", "summary", "description", "repeat_type", "all_day", "start_time", "end_time", "creator", "rsvp_status"})
}

func runCalendarDeleteSchedule(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.DeleteSchedule(ctx, args[0], args[1], "", "", "", calDeleteUserToken, calDeleteUserID)
	checkError(err)
	outputResult(result)
}

func runCalendarListSchedules(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	startTime, _ := strconv.ParseInt(args[1], 10, 64)
	endTime, _ := strconv.ParseInt(args[2], 10, 64)
	result, err := client.FetchScheduleList(ctx, args[0], startTime, endTime, calListUserToken, calListUserID)
	checkError(err)
	outputResultFields(result, []string{"schedule_list"})
}

func runCalendarAttendees(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchScheduleAttendees(ctx, args[0], args[1], calAttendeesPage, calAttendeesSize, calAttendeesUserToken, calAttendeesUserID)
	checkError(err)
	outputResultFields(result, []string{"total"})
}

func runCalendarAddAttendees(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	attendees := parseStringList(args[2])

	result, err := client.AddScheduleAttendees(ctx, args[0], args[1], attendees, calAddAttendeesReminderType, "", "", calAddAttendeesUserToken, calAddAttendeesUserID)
	checkError(err)
	outputResult(result)
}

func runCalendarDeleteAttendees(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	attendees := parseStringList(args[2])

	result, err := client.DeleteScheduleAttendees(ctx, args[0], args[1], attendees, calDelAttendeesReminderType, "", "", calDelAttendeesUserToken, calDelAttendeesUserID)
	checkError(err)
	outputResult(result)
}

func runCalendarUpdateSchedule(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	params := map[string]interface{}{}

	if calUpdateSummary != "" {
		params["summary"] = calUpdateSummary
	}
	if calUpdateDesc != "" {
		params["description"] = calUpdateDesc
	}
	if calUpdateOp != "" {
		params["operationType"] = calUpdateOp
	}
	if calUpdateCurrentTime != 0 {
		params["currentTime"] = calUpdateCurrentTime
	}
	if calUpdateReminder != "" {
		params["reminderType"] = calUpdateReminder
	}
	if calUpdateRepeat != "" {
		params["repeatType"] = calUpdateRepeat
	}
	if calUpdateRule != "" {
		rule, rErr := parseJSONRaw(calUpdateRule)
		checkError(rErr)
		if m, ok := rule.(map[string]interface{}); ok {
			params["rule"] = m
		}
	}
	if calUpdateExpire != "" {
		params["expireDateType"] = calUpdateExpire
	}
	if calUpdateAllDay != "" {
		params["allDay"] = calUpdateAllDay
	}
	if calUpdatePermissions != "" {
		params["attendeePermissions"] = calUpdatePermissions
	}
	if calUpdateStartTime != "" {
		st, sErr := parseJSONRaw(calUpdateStartTime)
		checkError(sErr)
		if m, ok := st.(map[string]interface{}); ok {
			params["startTime"] = m
		}
	}
	if calUpdateEndTime != "" {
		et, eErr := parseJSONRaw(calUpdateEndTime)
		checkError(eErr)
		if m, ok := et.(map[string]interface{}); ok {
			params["endTime"] = m
		}
	}

	result, err := client.UpdateSchedule(ctx, args[0], args[1], params, calUpdateScheduleUserToken, calUpdateScheduleUserID)
	checkError(err)
	outputResult(result)
}

func runCalendarAttendeeMeta(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	params := map[string]interface{}{}

	if calAttendeeMetaRsvp != "" {
		params["rsvpStatus"] = calAttendeeMetaRsvp
	}
	if calAttendeeMetaColor != "" {
		params["color"] = calAttendeeMetaColor
	}
	if calAttendeeMetaPermissions != "" {
		params["permissions"] = calAttendeeMetaPermissions
	}
	if calAttendeeMetaBusyFree != "" {
		params["busyFreeState"] = calAttendeeMetaBusyFree
	}
	if calAttendeeMetaRemindTimes != "" {
		rt, rtErr := parseJSONRaw(calAttendeeMetaRemindTimes)
		checkError(rtErr)
		params["remindTimes"] = rt
	}

	result, err := client.UpdateScheduleAttendeeMeta(ctx, args[0], args[1], params, calAttendeeMetaUserToken, calAttendeeMetaUserID)
	checkError(err)
	outputResult(result)
}