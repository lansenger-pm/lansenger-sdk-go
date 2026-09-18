package main

import (
	"context"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var noticeCmd = &cobra.Command{
	Use:   "notice",
	Short: "Notice operations via official accounts (通知系统)",
}

var noticeSendCmd = &cobra.Command{
	Use:   "send TITLE ACCOUNT_CODE",
	Short: "Send a notice via an official account",
	Args:  cobra.ExactArgs(2),
	Run:   runNoticeSend,
}

var noticeAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List official accounts (fetch accountCode for sending)",
	Args:  cobra.NoArgs,
	Run:   runNoticeAccounts,
}

var (
	noticeSendContent      string
	noticeSendLink         string
	noticeSendContentType  int
	noticeSendUserType     int
	noticeSendPhones       string
	noticeSendCCPhones     string
	noticeSendRange        string
	noticeSendCCStaffIDs   string
	noticeSendCreateMobile string
	noticeSendCreateUserID string
	noticeSendLocation     string
	noticeSendResources    string
	noticeSendExtendID     string

	noticeSendConfirmFlag   int
	noticeSendForwardFlag   int
	noticeSendReplyFlag     int
	noticeSendAnonymousFlag int

	noticeSendRemindStatus     int
	noticeSendRemindMsgType    string
	noticeSendAtOnceFlag       int
	noticeSendRemindAfterType  string
	noticeSendRemindMaxCount   int
	noticeSendRemindInterval   int
	noticeSendRemindUnit       string
	noticeSendRemindRange      string
	noticeSendRemindExclude    string
	noticeSendUserToken        string

	noticeAccountsOrgID    string
	noticeAccountsUserToken string
)

func init() {
	noticeSendCmd.Flags().StringVar(&noticeSendContent, "content", "", "Text content (required when --content-type is 1)")
	noticeSendCmd.Flags().StringVar(&noticeSendLink, "link", "", "Link URL (required when --content-type is 2)")
	noticeSendCmd.Flags().IntVar(&noticeSendContentType, "content-type", 1, "Content type: 1=text, 2=link")
	noticeSendCmd.Flags().IntVar(&noticeSendUserType, "user-type", 1, "Targeting type: 1=phone, 2=staff/department")
	noticeSendCmd.Flags().StringVar(&noticeSendPhones, "release-phones", "", "Comma-separated receiver phones (user-type=1, max 10)")
	noticeSendCmd.Flags().StringVar(&noticeSendCCPhones, "cc-phones", "", "Comma-separated cc phones (user-type=1, max 10)")
	noticeSendCmd.Flags().StringVar(&noticeSendRange, "release-range", "", `JSON list (user-type=2, max 200): [{"objId":"dept-1","objName":"RD","objType":2}]`)
	noticeSendCmd.Flags().StringVar(&noticeSendCCStaffIDs, "cc-staff-ids", "", "Comma-separated cc staff IDs (user-type=2, max 200)")
	noticeSendCmd.Flags().StringVar(&noticeSendCreateMobile, "create-mobile", "", "Operator mobile (user-type=1; omit when --as/--user-token is set)")
	noticeSendCmd.Flags().StringVar(&noticeSendCreateUserID, "create-user-id", "", "Creator staff ID (user-type=2; omit when --as/--user-token is set)")
	noticeSendCmd.Flags().StringVar(&noticeSendLocation, "location", "", "Notice address")
	noticeSendCmd.Flags().StringVar(&noticeSendResources, "resources", "", `JSON list of attachments: [{"fileName":"a.pdf","resourceId":"res-1","fileType":"application/pdf","fileSize":1024}]`)
	noticeSendCmd.Flags().StringVar(&noticeSendExtendID, "extend-id", "", "Caller-side correlation ID")

	noticeSendCmd.Flags().IntVar(&noticeSendConfirmFlag, "confirm-flag", -1, "Require read confirmation: 1=yes, 0=no (-1=omit)")
	noticeSendCmd.Flags().IntVar(&noticeSendForwardFlag, "forward-flag", -1, "Allow forwarding: 1=yes, 0=no (-1=omit)")
	noticeSendCmd.Flags().IntVar(&noticeSendReplyFlag, "reply-flag", -1, "Allow reply: 1=yes, 0=no (-1=omit)")
	noticeSendCmd.Flags().IntVar(&noticeSendAnonymousFlag, "anonymous-flag", -1, "Allow anonymous reply: 1=yes, 0=no (-1=omit)")

	noticeSendCmd.Flags().IntVar(&noticeSendRemindStatus, "remind-status", -1, "Enable reminding: 1=yes, 0=no (-1=omit)")
	noticeSendCmd.Flags().StringVar(&noticeSendRemindMsgType, "remind-msg-type", "", "Remind channels, comma-separated: mobile,sms,app")
	noticeSendCmd.Flags().IntVar(&noticeSendAtOnceFlag, "at-once", -1, "Remind immediately: 1=yes, 0=no (-1=omit)")
	noticeSendCmd.Flags().StringVar(&noticeSendRemindAfterType, "remind-after", "", "Follow-up remind type: never, unOperate, count")
	noticeSendCmd.Flags().IntVar(&noticeSendRemindMaxCount, "remind-max-count", -1, "Max remind count (remind-after=count; -1=omit)")
	noticeSendCmd.Flags().IntVar(&noticeSendRemindInterval, "remind-interval", -1, "Remind interval value (-1=omit)")
	noticeSendCmd.Flags().StringVar(&noticeSendRemindUnit, "remind-interval-unit", "", "Interval unit: minutes, hour, day")
	noticeSendCmd.Flags().StringVar(&noticeSendRemindRange, "remind-range", "", "Remind range type: all, receiver, partialRemind, notReminder")
	noticeSendCmd.Flags().StringVar(&noticeSendRemindExclude, "remind-exclude", "", "Comma-separated staff IDs excluded from reminding")
	noticeSendCmd.Flags().StringVar(&noticeSendUserToken, "user-token", "", "User token")

	noticeAccountsCmd.Flags().StringVar(&noticeAccountsOrgID, "org-id", "", "Organization ID")
	noticeAccountsCmd.Flags().StringVar(&noticeAccountsUserToken, "user-token", "", "User token")

	noticeCmd.AddCommand(noticeSendCmd)
	noticeCmd.AddCommand(noticeAccountsCmd)
	rootCmd.AddCommand(noticeCmd)
}

// noticeFlagPtr converts a tri-state int flag (-1 = omit) into the *int
// the SDK expects (nil = omit, 0/1 = explicit value).
func noticeFlagPtr(v int) *int {
	if v < 0 {
		return nil
	}
	val := v
	return &val
}

func noticeJSONMapList(raw string) []map[string]interface{} {
	if raw == "" {
		return nil
	}
	decoded, err := parseJSONRaw(raw)
	checkError(err)
	list, ok := decoded.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func runNoticeSend(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	p := &lansenger.NoticeSendParams{
		Title:                      args[0],
		AccountCode:                args[1],
		ContentType:                noticeSendContentType,
		UserType:                   noticeSendUserType,
		Content:                    noticeSendContent,
		NoticeLink:                 noticeSendLink,
		NoticeLocation:             noticeSendLocation,
		ReleasePhones:              splitCommaList(noticeSendPhones),
		CCPhones:                   splitCommaList(noticeSendCCPhones),
		ReleaseRange:               noticeJSONMapList(noticeSendRange),
		CCStaffIDs:                 splitCommaList(noticeSendCCStaffIDs),
		CreateMobile:               noticeSendCreateMobile,
		CreateUserID:               noticeSendCreateUserID,
		ResourceList:               noticeJSONMapList(noticeSendResources),
		ExtendID:                   noticeSendExtendID,
		ConfirmFlag:                noticeFlagPtr(noticeSendConfirmFlag),
		ForwardFlag:                noticeFlagPtr(noticeSendForwardFlag),
		ReplyFlag:                  noticeFlagPtr(noticeSendReplyFlag),
		AnonymousFlag:              noticeFlagPtr(noticeSendAnonymousFlag),
		RemindStatus:               noticeFlagPtr(noticeSendRemindStatus),
		RemindMsgType:              noticeSendRemindMsgType,
		AtOnceFlag:                 noticeFlagPtr(noticeSendAtOnceFlag),
		RemindAfterType:            noticeSendRemindAfterType,
		RemindMaxCount:             noticeFlagPtr(noticeSendRemindMaxCount),
		RemindIntervalTime:         noticeFlagPtr(noticeSendRemindInterval),
		RemindIntervalTimeDuration: noticeSendRemindUnit,
		RemindRangeType:            noticeSendRemindRange,
		RemindRangeStaffIDs:        splitCommaList(noticeSendRemindExclude),
		UserToken:                  noticeSendUserToken,
	}

	result, err := client.SendNotice(ctx, p)
	checkError(err)
	outputResultFields(result, []string{"notice_code", "title", "notice_status", "confirm_status", "publish_user_name"})
}

func runNoticeAccounts(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchNoticeAccounts(ctx, noticeAccountsOrgID, noticeAccountsUserToken)
	checkError(err)
	outputResultFields(result, []string{"total", "accounts"})
}
