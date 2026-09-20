package main

import (
	"context"
	"strconv"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var boardroomCmd = &cobra.Command{
	Use:   "boardroom",
	Short: "Meeting-room reservation (会议室预定 V2)",
}

var boardroomRoomsCmd = &cobra.Command{
	Use:   "rooms",
	Short: "Filter meeting rooms (paged)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchBoardroomList(context.Background(),
			brGradingID, brAreaOfficeID, splitCommaList(brFloorIDs), splitCommaList(brEquipment),
			brTimeStart, brTimeEnd, brQueryDate, brPage, brLimit, "", "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"count", "items"})
	},
}

var boardroomRoomDetailCmd = &cobra.Command{
	Use:   "room-detail ROOM_ID",
	Short: "Fetch meeting-room detail",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchBoardroomDetail(context.Background(), args[0], "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"room_id", "name", "status", "people_num", "can_reserve_flag", "area_name", "address"})
	},
}

var boardroomScheduleCmd = &cobra.Command{
	Use:   "schedule ROOM_ID QUERY_DATE",
	Short: "Fetch a room's bookings + deactivations for a date",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchBoardroomSchedule(context.Background(), args[0], args[1], brScheduleGradingID, "", "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"room_id", "name", "people_num", "can_reserve_flag", "reserves", "deactivations"})
	},
}

var boardroomReserveDetailCmd = &cobra.Command{
	Use:   "reserve-detail RESERVE_ROOM_ID",
	Short: "Fetch reservation detail (attendees, approval flow)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchBoardroomReserveDetail(context.Background(), args[0], brReserveDetailGradingID, "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"reserve_id", "boardroom_name", "meeting_name", "status", "reserve_time_start", "reserve_time_end", "reserve_user_name", "people_number"})
	},
}

var boardroomReserveCmd = &cobra.Command{
	Use:   "reserve BOARD_ROOM_ID NAME",
	Short: "Reserve a meeting room (single or repeating)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.ReserveBoardroom(context.Background(), &lansenger.BoardroomReserveParams{
			BoardRoomID:        args[0],
			Name:               args[1],
			GradingID:          brReserveGradingID,
			ReserveTimeStart:   brReserveStart,
			ReserveTimeEnd:     brReserveEnd,
			NoticeTime:         brReserveNoticeTime,
			PeopleNumber:       brReservePeople,
			Toastmaster:        brReserveToastmaster,
			Leader:             brReserveLeader,
			LeaderAttend:       brReserveLeaderAttend,
			IsVideo:            brReserveIsVideo,
			VideoName:          brReserveVideoName,
			OtherDemand:        brReserveOtherDemand,
			TableCards:         brReserveTableCards,
			InvitationUserList: splitCommaList(brReserveInvite),
			UserList:           splitCommaList(brReserveApprovers),
			ReserveType:        brReserveType,
			RepeatType:         brReserveRepeatType,
			RepeatDays:         brReserveRepeatDaysParsed(),
			Skip:               brReserveSkip,
			RepeatEndDateStr:   brReserveRepeatEnd,
			UserToken:          brReserveUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"reserve_id", "reserve_code", "boardroom_name", "meeting_name", "status", "reserve_time"})
	},
}

var boardroomEditReserveCmd = &cobra.Command{
	Use:   "edit-reserve RESERVE_ID BOARD_ROOM_ID NAME",
	Short: "Edit a reservation (only non-approval-flow bookings)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.EditBoardroomReserve(context.Background(), &lansenger.BoardroomEditParams{
			ReserveID:        args[0],
			BoardRoomID:      args[1],
			Name:             args[2],
			GradingID:        brEditGradingID,
			ReserveTimeStart: brEditStart,
			ReserveTimeEnd:   brEditEnd,
			NoticeTime:       brEditNoticeTime,
			EditType:         brEditType,
			PeopleNumber:     brEditPeople,
			UserToken:        brEditUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"reserve_id", "reserve_code", "meeting_name", "status", "reserve_time"})
	},
}

var boardroomCancelCmd = &cobra.Command{
	Use:   "cancel RESERVE_ID",
	Short: "Cancel a reservation (status 0/1/5 only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		confirmHighRisk("cancel", "reservation "+args[0], brCancelYes, brCancelDryRun)
		if brCancelDryRun {
			return
		}
		client := getClient()
		var isSend *bool
		if brCancelNotify {
			v := true
			isSend = &v
		}
		result, err := client.CancelBoardroomReserve(context.Background(), args[0], "", "", brCancelReason, isSend, nil, "", brCancelType, brCancelUserToken)
		checkError(err)
		outputResultFields(result, []string{"done"})
	},
}

var boardroomConfirmSignCmd = &cobra.Command{
	Use:   "confirm-sign RESERVE_ID",
	Short: "Scan-code confirmation (status 1 only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.ConfirmBoardroomSign(context.Background(), args[0], "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"done"})
	},
}

var boardroomMyReservesCmd = &cobra.Command{
	Use:   "my-reserves",
	Short: "Page my reservations with filters",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchMyBoardroomReserves(context.Background(), brMyGradingID, brMyKeys, brMyStartTime, brMyEndTime, brMyRoomID, nil, brPage, brLimit, "", "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"count", "items"})
	},
}

var boardroomGradingsCmd = &cobra.Command{
	Use:   "gradings",
	Short: "Fetch gradings visible to the user (gradingId source)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchBoardroomGradings(context.Background(), "", "", brUserToken)
		checkError(err)
		outputResultFields(result, []string{"total", "gradings"})
	},
}

var boardroomAreaOfficesCmd = &cobra.Command{
	Use:   "area-offices GRADING_ID",
	Short: "Fetch office areas under a grading",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchBoardroomAreaOffices(context.Background(), args[0], brUserToken)
		checkError(err)
		outputResultFields(result, []string{"total", "areas"})
	},
}

var (
	brGradingID, brAreaOfficeID, brFloorIDs, brEquipment                  string
	brTimeStart, brTimeEnd, brQueryDate, brUserToken                      string
	brPage, brLimit                                                       int
	brScheduleGradingID                                                   string
	brReserveDetailGradingID                                              string
	brReserveGradingID, brReserveStart, brReserveEnd, brReserveNoticeTime string
	brReservePeople, brReserveToastmaster, brReserveLeader                string
	brReserveLeaderAttend, brReserveOtherDemand, brReserveTableCards      string
	brReserveIsVideo, brReserveVideoName, brReserveInvite                 string
	brReserveApprovers, brReserveType, brReserveRepeatType                string
	brReserveRepeatDays, brReserveSkip, brReserveRepeatEnd                string
	brReserveUserToken                                                    string
	brEditGradingID, brEditStart, brEditEnd, brEditNoticeTime             string
	brEditType, brEditPeople, brEditUserToken                             string
	brCancelReason, brCancelType, brCancelUserToken                       string
	brCancelNotify, brCancelYes, brCancelDryRun                           bool
	brMyGradingID, brMyKeys, brMyStartTime, brMyEndTime, brMyRoomID       string
)

func init() {
	boardroomRoomsCmd.Flags().StringVar(&brGradingID, "grading-id", "", "Grading (分区) ID")
	boardroomRoomsCmd.Flags().StringVar(&brAreaOfficeID, "area-office-id", "", "Office area ID")
	boardroomRoomsCmd.Flags().StringVar(&brFloorIDs, "floor-ids", "", "Comma-separated floor IDs")
	boardroomRoomsCmd.Flags().StringVar(&brEquipment, "equipment", "", "Comma-separated equipment")
	boardroomRoomsCmd.Flags().StringVar(&brTimeStart, "time-start", "", "Filter start (yyyy-MM-dd HH:mm)")
	boardroomRoomsCmd.Flags().StringVar(&brTimeEnd, "time-end", "", "Filter end (yyyy-MM-dd HH:mm)")
	boardroomRoomsCmd.Flags().StringVar(&brQueryDate, "date", "", "Query date (yyyy-MM-dd)")
	boardroomRoomsCmd.Flags().IntVar(&brPage, "page", 1, "Page number")
	boardroomRoomsCmd.Flags().IntVar(&brLimit, "size", 10, "Page size")
	boardroomRoomsCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")

	boardroomScheduleCmd.Flags().StringVar(&brScheduleGradingID, "grading-id", "", "Grading (分区) ID (required)")
	boardroomScheduleCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")

	boardroomReserveDetailCmd.Flags().StringVar(&brReserveDetailGradingID, "grading-id", "", "Grading (分区) ID")
	boardroomReserveDetailCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")

	boardroomReserveCmd.Flags().StringVar(&brReserveGradingID, "grading-id", "", "Grading (分区) ID (required)")
	boardroomReserveCmd.Flags().StringVar(&brReserveStart, "start", "", "Reserve start (yyyy-MM-dd HH:mm:ss) (required)")
	boardroomReserveCmd.Flags().StringVar(&brReserveEnd, "end", "", "Reserve end (yyyy-MM-dd HH:mm:ss) (required)")
	boardroomReserveCmd.Flags().StringVar(&brReserveNoticeTime, "notice-time", "", "Meeting reminder (required)")
	boardroomReserveCmd.Flags().StringVar(&brReservePeople, "people", "", "Attendee count")
	boardroomReserveCmd.Flags().StringVar(&brReserveToastmaster, "toastmaster", "", "Host (max 10 chars)")
	boardroomReserveCmd.Flags().StringVar(&brReserveLeader, "leader", "", "Attending leader (max 200 chars)")
	boardroomReserveCmd.Flags().StringVar(&brReserveLeaderAttend, "leader-attend", "", "Leader attendance: 0=attend, 1=not attend")
	boardroomReserveCmd.Flags().StringVar(&brReserveIsVideo, "is-video", "", "Video meeting: 0=on, 1=off")
	boardroomReserveCmd.Flags().StringVar(&brReserveVideoName, "video-name", "", "Video meeting name")
	boardroomReserveCmd.Flags().StringVar(&brReserveOtherDemand, "other-demand", "", "Other meeting requirements")
	boardroomReserveCmd.Flags().StringVar(&brReserveTableCards, "table-cards", "", "Table cards: 0=on, 1=off")
	boardroomReserveCmd.Flags().StringVar(&brReserveInvite, "invite", "", "Comma-separated attendee staff IDs")
	boardroomReserveCmd.Flags().StringVar(&brReserveApprovers, "approvers", "", "Comma-separated approver staff IDs")
	boardroomReserveCmd.Flags().StringVar(&brReserveType, "reserve-type", "0", "0=single, 1=repeat")
	boardroomReserveCmd.Flags().StringVar(&brReserveRepeatType, "repeat-type", "", "day/week/month")
	boardroomReserveCmd.Flags().StringVar(&brReserveRepeatDays, "repeat-days", "", "Comma-separated repeat day numbers")
	boardroomReserveCmd.Flags().StringVar(&brReserveSkip, "skip", "", "Skip weekends/holidays: 0=skip, 1=don't")
	boardroomReserveCmd.Flags().StringVar(&brReserveRepeatEnd, "repeat-end", "", "Repeat end (yyyy-MM-dd HH:mm:ss)")
	boardroomReserveCmd.Flags().StringVar(&brReserveUserToken, "user-token", "", "User token")

	boardroomEditReserveCmd.Flags().StringVar(&brEditGradingID, "grading-id", "", "Grading (分区) ID (required)")
	boardroomEditReserveCmd.Flags().StringVar(&brEditStart, "start", "", "Reserve start (required)")
	boardroomEditReserveCmd.Flags().StringVar(&brEditEnd, "end", "", "Reserve end (required)")
	boardroomEditReserveCmd.Flags().StringVar(&brEditNoticeTime, "notice-time", "", "Meeting reminder (required)")
	boardroomEditReserveCmd.Flags().StringVar(&brEditType, "edit-type", "1", "1=this booking, 2=this and following")
	boardroomEditReserveCmd.Flags().StringVar(&brEditPeople, "people", "", "Attendee count")
	boardroomEditReserveCmd.Flags().StringVar(&brEditUserToken, "user-token", "", "User token")

	boardroomCancelCmd.Flags().StringVar(&brCancelReason, "reason", "", "Cancel reason (max 200 chars)")
	boardroomCancelCmd.Flags().StringVar(&brCancelType, "cancel-type", "", "1=this, 2=this and following, 3=all unfinished")
	boardroomCancelCmd.Flags().BoolVar(&brCancelNotify, "notify", false, "Notify attendees")
	boardroomCancelCmd.Flags().StringVar(&brCancelUserToken, "user-token", "", "User token")
	boardroomCancelCmd.Flags().BoolVarP(&brCancelYes, "yes", "y", false, "Confirm cancellation before executing")
	boardroomCancelCmd.Flags().BoolVar(&brCancelDryRun, "dry-run", false, "Validate inputs without cancelling")

	boardroomConfirmSignCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")

	boardroomMyReservesCmd.Flags().StringVar(&brMyGradingID, "grading-id", "", "Grading (分区) ID (required)")
	boardroomMyReservesCmd.Flags().StringVar(&brMyKeys, "keys", "", "Keyword search")
	boardroomMyReservesCmd.Flags().StringVar(&brMyStartTime, "start-time", "", "Reserve start filter")
	boardroomMyReservesCmd.Flags().StringVar(&brMyEndTime, "end-time", "", "Reserve end filter")
	boardroomMyReservesCmd.Flags().StringVar(&brMyRoomID, "room-id", "", "Meeting-room ID")
	boardroomMyReservesCmd.Flags().IntVar(&brPage, "page", 1, "Page number")
	boardroomMyReservesCmd.Flags().IntVar(&brLimit, "size", 10, "Page size")
	boardroomMyReservesCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")

	boardroomGradingsCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")
	boardroomAreaOfficesCmd.Flags().StringVar(&brUserToken, "user-token", "", "User token")

	boardroomCmd.AddCommand(boardroomRoomsCmd)
	boardroomCmd.AddCommand(boardroomRoomDetailCmd)
	boardroomCmd.AddCommand(boardroomScheduleCmd)
	boardroomCmd.AddCommand(boardroomReserveDetailCmd)
	boardroomCmd.AddCommand(boardroomReserveCmd)
	boardroomCmd.AddCommand(boardroomEditReserveCmd)
	boardroomCmd.AddCommand(boardroomCancelCmd)
	boardroomCmd.AddCommand(boardroomConfirmSignCmd)
	boardroomCmd.AddCommand(boardroomMyReservesCmd)
	boardroomCmd.AddCommand(boardroomGradingsCmd)
	boardroomCmd.AddCommand(boardroomAreaOfficesCmd)
	rootCmd.AddCommand(boardroomCmd)
}

// brReserveRepeatDaysParsed converts the comma string into ints for the SDK call.
func brReserveRepeatDaysParsed() []int {
	if brReserveRepeatDays == "" {
		return nil
	}
	out := make([]int, 0)
	for _, s := range splitCommaList(brReserveRepeatDays) {
		if v, err := strconv.Atoi(s); err == nil {
			out = append(out, v)
		}
	}
	return out
}
