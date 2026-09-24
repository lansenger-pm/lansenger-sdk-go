package main

import (
	"context"
	"encoding/json"
	"fmt"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var videoconferenceCmd = &cobra.Command{
	Use:   "videoconference",
	Short: "Video-meeting open APIs (视频会议开放能力)",
}

var (
	vcStartTime       int64
	vcEndTime         int64
	vcMembers         string
	vcEvents          string
	vcVods            string
	vcSubject         string
	vcOrgID           string
	vcOperator        string
	vcFetchRange      string
	vcStaffID         string
	vcLimit           int
	vcOffset          int
	vcAutoRecord      int
	vcMeetingType     int
	vcGroupNew        int
	vcConfPassword    string
	vcControlPassword string
	vcMaskType        int
	vcExtAttr         string
	vcJoinMute        int
	vcOpenMute        int
	vcEnablePreJoin   int
	vcUserStopTime    int64
	vcInviteAdmin     int
	vcCallBackInfo    string
	vcAdmin           string
	vcCreateSource    int
	vcMeetingNumber   string
	vcYes             bool
	vcDryRun          bool
	vcUserToken       string
)

// parseVCMembers decodes a --member JSON list into VideoconferenceMember
// entries and enforces the exactly-one-host rule locally for a clear CLI
// error before any network call.
func parseVCMembers(raw string, flagName string) ([]lansenger.VideoconferenceMember, error) {
	if raw == "" {
		return nil, fmt.Errorf("%s is required (JSON list with exactly one role='admin')", flagName)
	}
	var members []lansenger.VideoconferenceMember
	if err := json.Unmarshal([]byte(raw), &members); err != nil {
		return nil, fmt.Errorf("%s must be a JSON list of members: %v", flagName, err)
	}
	hosts := 0
	for _, m := range members {
		if m.Role == lansenger.VCMemberRoleHost {
			hosts++
		}
	}
	if hosts != 1 {
		return nil, fmt.Errorf("%s must contain exactly one member with role='admin'", flagName)
	}
	return members, nil
}

// parseVCVods decodes --vods as a JSON list of vodId strings or objects and
// enforces the 1..3 limit locally.
func parseVCVods(raw string) ([]lansenger.VideoconferenceVod, error) {
	if raw == "" {
		return nil, fmt.Errorf("--vods is required (1..3 entries)")
	}
	var rawList []json.RawMessage
	if err := json.Unmarshal([]byte(raw), &rawList); err != nil {
		return nil, fmt.Errorf("--vods must be a JSON list: %v", err)
	}
	if len(rawList) == 0 || len(rawList) > 3 {
		return nil, fmt.Errorf("--vods must contain 1..3 entries")
	}
	vods := make([]lansenger.VideoconferenceVod, 0, len(rawList))
	for _, item := range rawList {
		var s string
		if err := json.Unmarshal(item, &s); err == nil {
			vods = append(vods, lansenger.VideoconferenceVod{VodID: s})
			continue
		}
		var obj lansenger.VideoconferenceVod
		if err := json.Unmarshal(item, &obj); err != nil {
			return nil, fmt.Errorf("--vods entries must be vodId strings or {\"vodId\": ...} objects")
		}
		vods = append(vods, obj)
	}
	return vods, nil
}

var vcCreateCmd = &cobra.Command{
	Use:   "create SUBJECT",
	Short: "Create a meeting (instant or reserved)",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if vcStartTime <= 0 {
			return fmt.Errorf("--start-time is required (epoch milliseconds, future for reserved meetings)")
		}
		if vcOrgID == "" {
			return fmt.Errorf("--org-id is required")
		}
		_, err := parseVCMembers(vcMembers, "--member")
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		members, err := parseVCMembers(vcMembers, "--member")
		checkError(err)
		result, err := getClient().CreateMeeting(context.Background(), &lansenger.VideoconferenceCreateParams{
			Subject:         args[0],
			StartTime:       vcStartTime,
			Members:         members,
			OrgID:           vcOrgID,
			AutoRecord:      vcAutoRecord,
			Type:            vcMeetingType,
			GroupNew:        vcGroupNew,
			ConfPassword:    vcConfPassword,
			ControlPassword: vcControlPassword,
			MaskType:        vcMaskType,
			ExtAttr:         vcExtAttr,
			JoinMute:        intPtrOrNil(vcJoinMute),
			OpenMute:        intPtrOrNil(vcOpenMute),
			EnablePreJoin:   intPtrOrNil(vcEnablePreJoin),
			UserStopTime:    int64PtrOrNil(vcUserStopTime),
			InviteAdmin:     intPtrOrNil(vcInviteAdmin),
			UserToken:       vcUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"mid", "subject", "meeting_number", "start_time", "type", "status"})
	},
}

var vcModifyCmd = &cobra.Command{
	Use:   "modify MID",
	Short: "Modify a meeting that has not started",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if vcOrgID == "" {
			return fmt.Errorf("--org-id is required")
		}
		_, err := parseVCMembers(vcMembers, "--member")
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		members, err := parseVCMembers(vcMembers, "--member")
		checkError(err)
		result, err := getClient().ModifyMeeting(context.Background(), &lansenger.VideoconferenceModifyParams{
			Mid:             args[0],
			Subject:         vcSubject,
			StartTime:       vcStartTime,
			Members:         members,
			OrgID:           vcOrgID,
			Operator:        vcOperator,
			AutoRecord:      vcAutoRecord,
			Type:            vcMeetingType,
			GroupNew:        vcGroupNew,
			ConfPassword:    vcConfPassword,
			ControlPassword: vcControlPassword,
			UserStopTime:    int64PtrOrNil(vcUserStopTime),
			UserToken:       vcUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"done", "message"})
	},
}

var vcCancelCmd = &cobra.Command{
	Use:   "cancel MID",
	Short: "Cancel a meeting that has not started",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		confirmHighRisk("cancel", "meeting "+args[0], vcYes, vcDryRun)
		if vcDryRun {
			return
		}
		result, err := getClient().CancelMeeting(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"done", "message"})
	},
}

var vcStopCmd = &cobra.Command{
	Use:   "stop MID",
	Short: "End a running meeting",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		confirmHighRisk("stop", "meeting "+args[0], vcYes, vcDryRun)
		if vcDryRun {
			return
		}
		result, err := getClient().StopMeeting(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"done", "message"})
	},
}

var vcDetailCmd = &cobra.Command{
	Use:   "detail MID",
	Short: "Fetch meeting detail by mid",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMeetingDetail(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"mid", "subject", "meeting_number", "start_time", "stop_time", "type", "status", "admin"})
	},
}

var vcListCmd = &cobra.Command{
	Use:   "list ORG_ID START_TIME END_TIME",
	Short: "Meeting list by time range (epoch milliseconds)",
	Args:  cobra.ExactArgs(3),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		start, err := parseEpochArg(args[1], "START_TIME")
		if err != nil {
			return err
		}
		end, err := parseEpochArg(args[2], "END_TIME")
		if err != nil {
			return err
		}
		vcStartTime, vcEndTime = start, end
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMeetingList(context.Background(), args[0], vcStartTime, vcEndTime, vcFetchRange, vcStaffID, vcLimit, vcOffset, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcRecordListCmd = &cobra.Command{
	Use:   "record-list ORG_ID START_TIME END_TIME",
	Short: "Meeting operation record list",
	Args:  cobra.ExactArgs(3),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		start, err := parseEpochArg(args[1], "START_TIME")
		if err != nil {
			return err
		}
		end, err := parseEpochArg(args[2], "END_TIME")
		if err != nil {
			return err
		}
		vcStartTime, vcEndTime = start, end
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMeetingRecordList(context.Background(), args[0], vcStartTime, vcEndTime, vcAdmin, vcCreateSource, vcLimit, vcOffset, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcSimpleRecordCmd = &cobra.Command{
	Use:   "simplerecord MID",
	Short: "Member join/leave records of a meeting",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMemberSimpleRecord(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken, vcLimit, vcOffset)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcFixroomListCmd = &cobra.Command{
	Use:   "fixroom-list",
	Short: "Fixed (cloud) meeting-room list",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchFixroomList(context.Background(), vcOrgID, vcOperator, vcUserToken, vcLimit, vcOffset)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcStatusCmd = &cobra.Command{
	Use:   "status MIDS ORG_ID",
	Short: "Batch meeting status by comma-separated mids",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMeetingStatus(context.Background(), splitCommaList(args[0]), args[1], vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"statuses"})
	},
}

var vcSubscribeCmd = &cobra.Command{
	Use:   "subscribe MID",
	Short: "Subscribe meeting status-change events",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		events := parseJSONListFlag(vcEvents, "--events")
		result, err := getClient().SubscribeMeetingEvents(context.Background(), args[0], vcOrgID, events, vcCallBackInfo, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"done", "message"})
	},
}

var vcParamsCmd = &cobra.Command{
	Use:   "params MEETING_NUMBER",
	Short: "Fetch meeting params by meetingNumber (PRS >=3.8)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMeetingParams(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"meeting_info"})
	},
}

var vcHistoryCmd = &cobra.Command{
	Use:   "history ORG_ID",
	Short: "A person's past meetings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchHistoryMeetings(context.Background(), args[0], vcOperator, vcUserToken, vcLimit, vcOffset)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcActiveCmd = &cobra.Command{
	Use:   "active ORG_ID",
	Short: "A person's running + reserved meetings",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchActiveMeetings(context.Background(), args[0], vcOperator, vcUserToken, vcLimit, vcOffset)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcMemberControlCmd = &cobra.Command{
	Use:   "member-control MID STAFF_ID OP_CODE",
	Short: "Host control on a member (kick/mute/setHost/...; per-member only)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().ControlMember(context.Background(), args[0], args[1], args[2], vcOperator, vcOrgID, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"done", "message"})
	},
}

var vcInviteCmd = &cobra.Command{
	Use:   "invite MEETING_NUMBER",
	Short: "Invite members to a running meeting",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		members, err := parseVCMembers(vcMembers, "--member")
		checkError(err)
		result, err := getClient().InviteMembers(context.Background(), args[0], members, vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"done", "message"})
	},
}

var vcMemberListCmd = &cobra.Command{
	Use:   "member-list MID",
	Short: "Paged member list of a meeting",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchMemberList(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken, vcLimit, vcOffset)
		checkError(err)
		outputResultFields(result, []string{"offset", "total", "items"})
	},
}

var vcVodListCmd = &cobra.Command{
	Use:   "vod-list MID",
	Short: "Recording list of a meeting",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchVodList(context.Background(), args[0], vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"items"})
	},
}

var vcVodDownloadCmd = &cobra.Command{
	Use:   "vod-download",
	Short: "Recording download URLs (max 3 vods per call)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		vods, err := parseVCVods(vcVods)
		checkError(err)
		result, err := getClient().FetchVodDownloadURLs(context.Background(), vods, vcOrgID, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"data"})
	},
}

var vcConfCmd = &cobra.Command{
	Use:   "conf ORG_ID",
	Short: "Org videoconference config (PRS >=3.8)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchOrgConf(context.Background(), args[0], vcMeetingNumber, vcOperator, vcUserToken)
		checkError(err)
		outputResultFields(result, []string{"max_person", "default_max_person", "allowed_record_flag", "force_passwd_flag", "space_size"})
	},
}

func parseEpochArg(arg, name string) (int64, error) {
	var v int64
	if _, err := fmt.Sscanf(arg, "%d", &v); err != nil {
		return 0, fmt.Errorf("%s must be epoch milliseconds", name)
	}
	return v, nil
}

func init() {
	vcCreateCmd.Flags().Int64Var(&vcStartTime, "start-time", 0, "Start time in epoch milliseconds (required; future for reserved meetings)")
	_ = vcCreateCmd.MarkFlagRequired("start-time")
	vcCreateCmd.Flags().StringVar(&vcMembers, "member", "", "JSON list of members [{staffId, employeeName, role}] with exactly one role='admin'")
	_ = vcCreateCmd.MarkFlagRequired("member")
	vcCreateCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	_ = vcCreateCmd.MarkFlagRequired("org-id")
	vcCreateCmd.Flags().IntVar(&vcAutoRecord, "auto-record", 0, "Auto record: 0=no, 1=yes")
	vcCreateCmd.Flags().IntVar(&vcMeetingType, "type", 1, "0=instant meeting, 1=reserved")
	vcCreateCmd.Flags().IntVar(&vcGroupNew, "group-new", 0, "Create group chat: 0=no, 1=yes")
	vcCreateCmd.Flags().StringVar(&vcConfPassword, "conf-password", "", "Meeting password")
	vcCreateCmd.Flags().StringVar(&vcControlPassword, "control-password", "", "Host control password")
	vcCreateCmd.Flags().IntVar(&vcMaskType, "mask-type", 0, "Mask type")
	vcCreateCmd.Flags().StringVar(&vcExtAttr, "ext-attr", "", "Extension attributes")
	vcCreateCmd.Flags().IntVar(&vcJoinMute, "join-mute", -1, "Join muted: 0=no, 1=yes (-1=omit)")
	vcCreateCmd.Flags().IntVar(&vcOpenMute, "open-mute", -1, "Open muted: 0=no, 1=yes (-1=omit)")
	vcCreateCmd.Flags().IntVar(&vcEnablePreJoin, "enable-pre-join", -1, "Pre-join screen: 0=no, 1=yes (-1=omit)")
	vcCreateCmd.Flags().Int64Var(&vcUserStopTime, "user-stop-time", -1, "Auto-stop time in epoch milliseconds (-1=omit)")
	vcCreateCmd.Flags().IntVar(&vcInviteAdmin, "invite-admin", -1, "Invite admin: 0=no, 1=yes (-1=omit)")
	vcCreateCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcModifyCmd.Flags().Int64Var(&vcStartTime, "start-time", 0, "New start time in epoch milliseconds")
	vcModifyCmd.Flags().StringVar(&vcSubject, "subject", "", "New subject")
	vcModifyCmd.Flags().StringVar(&vcMembers, "member", "", "JSON list of members with exactly one role='admin'")
	_ = vcModifyCmd.MarkFlagRequired("member")
	vcModifyCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	_ = vcModifyCmd.MarkFlagRequired("org-id")
	vcModifyCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcModifyCmd.Flags().IntVar(&vcAutoRecord, "auto-record", 0, "Auto record: 0=no, 1=yes")
	vcModifyCmd.Flags().IntVar(&vcMeetingType, "type", 1, "0=instant meeting, 1=reserved")
	vcModifyCmd.Flags().IntVar(&vcGroupNew, "group-new", 0, "Create group chat: 0=no, 1=yes")
	vcModifyCmd.Flags().StringVar(&vcConfPassword, "conf-password", "", "Meeting password")
	vcModifyCmd.Flags().StringVar(&vcControlPassword, "control-password", "", "Host control password")
	vcModifyCmd.Flags().Int64Var(&vcUserStopTime, "user-stop-time", -1, "Auto-stop time in epoch milliseconds (-1=omit)")
	vcModifyCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	for _, c := range []*cobra.Command{vcCancelCmd, vcStopCmd} {
		c.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
		c.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
		c.Flags().StringVar(&vcUserToken, "user-token", "", "User token")
		c.Flags().BoolVar(&vcYes, "yes", false, "Confirm the high-risk operation")
		c.Flags().BoolVar(&vcDryRun, "dry-run", false, "Print the planned action without performing it")
	}

	vcDetailCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcDetailCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcDetailCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcListCmd.Flags().StringVar(&vcFetchRange, "fetch-range", lansenger.VCFetchRangeAll, "my | all | person (person requires --staff-id)")
	vcListCmd.Flags().StringVar(&vcStaffID, "staff-id", "", "Staff ID (required when --fetch-range=person)")
	vcListCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcListCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcListCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcRecordListCmd.Flags().StringVar(&vcAdmin, "admin", "", "Filter by host staff ID")
	vcRecordListCmd.Flags().IntVar(&vcCreateSource, "create-source", 0, "0=platform client, 1=third-party app")
	vcRecordListCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcRecordListCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcRecordListCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcSimpleRecordCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcSimpleRecordCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcSimpleRecordCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcSimpleRecordCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcSimpleRecordCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcFixroomListCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcFixroomListCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcFixroomListCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcFixroomListCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcFixroomListCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcSubscribeCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcSubscribeCmd.Flags().StringVar(&vcEvents, "events", "", "JSON list of events to subscribe")
	_ = vcSubscribeCmd.MarkFlagRequired("events")
	vcSubscribeCmd.Flags().StringVar(&vcCallBackInfo, "callback-info", "", "Callback info")
	vcSubscribeCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcParamsCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcParamsCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcParamsCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcHistoryCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcHistoryCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcHistoryCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcHistoryCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcActiveCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcActiveCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcActiveCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcActiveCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcMemberControlCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcMemberControlCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcMemberControlCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcInviteCmd.Flags().StringVar(&vcMembers, "member", "", "JSON list of members [{staffId, employeeName, type, video, audio, typeMask}]")
	_ = vcInviteCmd.MarkFlagRequired("member")
	vcInviteCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcInviteCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcInviteCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcMemberListCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcMemberListCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcMemberListCmd.Flags().IntVar(&vcLimit, "limit", 10, "Page size")
	vcMemberListCmd.Flags().IntVar(&vcOffset, "offset", 0, "Page offset")
	vcMemberListCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcVodListCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcVodListCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcVodListCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcVodDownloadCmd.Flags().StringVar(&vcVods, "vods", "", "JSON list of vodIds (1..3 entries)")
	_ = vcVodDownloadCmd.MarkFlagRequired("vods")
	vcVodDownloadCmd.Flags().StringVar(&vcOrgID, "org-id", "", "Organization ID")
	vcVodDownloadCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcVodDownloadCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	vcConfCmd.Flags().StringVar(&vcMeetingNumber, "meeting-number", "", "Fetch the config of the org owning this meeting")
	vcConfCmd.Flags().StringVar(&vcOperator, "operator", "", "Operator staff ID")
	vcConfCmd.Flags().StringVar(&vcUserToken, "user-token", "", "User token")

	videoconferenceCmd.AddCommand(
		vcCreateCmd,
		vcModifyCmd,
		vcCancelCmd,
		vcStopCmd,
		vcDetailCmd,
		vcListCmd,
		vcRecordListCmd,
		vcSimpleRecordCmd,
		vcFixroomListCmd,
		vcStatusCmd,
		vcSubscribeCmd,
		vcParamsCmd,
		vcHistoryCmd,
		vcActiveCmd,
		vcMemberControlCmd,
		vcInviteCmd,
		vcMemberListCmd,
		vcVodListCmd,
		vcVodDownloadCmd,
		vcConfCmd,
	)
	rootCmd.AddCommand(videoconferenceCmd)
}
