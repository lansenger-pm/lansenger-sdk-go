package main

import (
	"context"

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

var calendarListSchedulesCmd = &cobra.Command{
	Use:   "list-schedules CALENDAR_ID START_TIME END_TIME",
	Short: "List schedules in a time range",
	Args:  cobra.ExactArgs(3),
	Run:   runCalendarListSchedules,
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

var (
	calPrimaryUserToken string
	calPrimaryUserID    string

	calCreateDesc              string
	calCreateAllDay            bool
	calCreateRepeatType        string
	calCreateReminderType      int
	calCreateAttendeePerms     int
	calCreateExpireDateType    int
	calCreateRule               string
	calCreateUserToken          string

	calFetchUserToken string
	calFetchUserID    string

	calDeleteUserToken string

	calListUserToken string

	calAttendeesUserToken string
	calAttendeesPage      int
	calAttendeesSize      int

	calAddAttendeesReminderType int
	calAddAttendeesUserToken    string

	calDelAttendeesReminderType int
	calDelAttendeesUserToken    string
)

func init() {
	calendarPrimaryCmd.Flags().StringVar(&calPrimaryUserToken, "user-token", "", "User token")
	calendarPrimaryCmd.Flags().StringVar(&calPrimaryUserID, "user-id", "", "User ID")

	calendarCreateScheduleCmd.Flags().StringVarP(&calCreateDesc, "desc", "d", "", "Schedule description")
	calendarCreateScheduleCmd.Flags().BoolVar(&calCreateAllDay, "all-day", false, "All-day event")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateRepeatType, "repeat", "", "Repeat type: no/daily/weekly/monthly/yearly/work_day/custom")
	calendarCreateScheduleCmd.Flags().IntVar(&calCreateReminderType, "reminder", 0, "Reminder type")
	calendarCreateScheduleCmd.Flags().IntVar(&calCreateAttendeePerms, "attendee-perms", 0, "Attendee permissions")
	calendarCreateScheduleCmd.Flags().IntVar(&calCreateExpireDateType, "expire-date-type", 0, "Expire date type")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateRule, "rule", "", "Repeat rule (JSON)")
	calendarCreateScheduleCmd.Flags().StringVar(&calCreateUserToken, "user-token", "", "User token")

	calendarFetchScheduleCmd.Flags().StringVar(&calFetchUserToken, "user-token", "", "User token")
	calendarFetchScheduleCmd.Flags().StringVar(&calFetchUserID, "user-id", "", "User ID")

	calendarDeleteScheduleCmd.Flags().StringVar(&calDeleteUserToken, "user-token", "", "User token")

	calendarListSchedulesCmd.Flags().StringVar(&calListUserToken, "user-token", "", "User token")

	calendarAttendeesCmd.Flags().StringVar(&calAttendeesUserToken, "user-token", "", "User token")
	calendarAttendeesCmd.Flags().IntVarP(&calAttendeesPage, "page", "p", 1, "Page number")
	calendarAttendeesCmd.Flags().IntVarP(&calAttendeesSize, "size", "s", 20, "Page size")

	calendarAddAttendeesCmd.Flags().IntVar(&calAddAttendeesReminderType, "reminder", 0, "Reminder type")
	calendarAddAttendeesCmd.Flags().StringVar(&calAddAttendeesUserToken, "user-token", "", "User token")

	calendarDeleteAttendeesCmd.Flags().IntVar(&calDelAttendeesReminderType, "reminder", 0, "Reminder type")
	calendarDeleteAttendeesCmd.Flags().StringVar(&calDelAttendeesUserToken, "user-token", "", "User token")

	calendarCmd.AddCommand(calendarPrimaryCmd)
	calendarCmd.AddCommand(calendarCreateScheduleCmd)
	calendarCmd.AddCommand(calendarFetchScheduleCmd)
	calendarCmd.AddCommand(calendarDeleteScheduleCmd)
	calendarCmd.AddCommand(calendarListSchedulesCmd)
	calendarCmd.AddCommand(calendarAttendeesCmd)
	calendarCmd.AddCommand(calendarAddAttendeesCmd)
	calendarCmd.AddCommand(calendarDeleteAttendeesCmd)
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
	startTime := args[2]
	endTime := args[3]

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

	result, err := client.CreateSchedule(ctx, calendarID, summary, startTime, endTime, attendees, calCreateDesc, calCreateAllDay, calCreateRepeatType, rule, calCreateExpireDateType, calCreateReminderType, calCreateAttendeePerms, calCreateUserToken)
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

	result, err := client.DeleteSchedule(ctx, args[0], args[1], 0, "", "", calDeleteUserToken)
	checkError(err)
	outputResult(result)
}

func runCalendarListSchedules(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchScheduleList(ctx, args[0], args[1], args[2], calListUserToken)
	checkError(err)
	outputResultFields(result, []string{"schedule_list"})
}

func runCalendarAttendees(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchScheduleAttendees(ctx, args[0], args[1], calAttendeesPage, calAttendeesSize)
	checkError(err)
	outputResultFields(result, []string{"total"})
}

func runCalendarAddAttendees(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	attendeeMaps, err := parseJSONArray(args[2])
	checkError(err)

	var attendees []map[string]interface{}
	for _, m := range attendeeMaps {
		converted := make(map[string]interface{}, len(m))
		for k, v := range m {
			converted[k] = v
		}
		attendees = append(attendees, converted)
	}

	result, err := client.AddScheduleAttendees(ctx, args[0], args[1], attendees, calAddAttendeesReminderType)
	checkError(err)
	outputResult(result)
}

func runCalendarDeleteAttendees(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	attendeeMaps, err := parseJSONArray(args[2])
	checkError(err)

	var attendees []map[string]interface{}
	for _, m := range attendeeMaps {
		converted := make(map[string]interface{}, len(m))
		for k, v := range m {
			converted[k] = v
		}
		attendees = append(attendees, converted)
	}

	result, err := client.DeleteScheduleAttendees(ctx, args[0], args[1], attendees, calDelAttendeesReminderType)
	checkError(err)
	outputResult(result)
}