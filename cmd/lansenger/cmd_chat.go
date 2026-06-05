package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Chat list and message history (4.24 MCP)",
}

var chatListCmd = &cobra.Command{
	Use:   "list",
	Short: "Fetch personal chat list",
	Args:  cobra.NoArgs,
	Run:   runChatList,
}

var chatMessagesCmd = &cobra.Command{
	Use:   "messages",
	Short: "Fetch messages from a specific conversation",
	Args:  cobra.NoArgs,
	Run:   runChatMessages,
}

var (
	chatListType      int
	chatListKeyword   string
	chatListStartTime int64
	chatListEndTime   int64
	chatListUserToken string

	chatMsgStaffID    string
	chatMsgGroupID    string
	chatMsgPageSize   int
	chatMsgVersion    string
	chatMsgStartTime  int64
	chatMsgEndTime    int64
	chatMsgSenderID   string
	chatMsgUserToken  string
	chatMsgSplitMonth bool
	chatMsgProgress   bool
)

func init() {
	chatListCmd.Flags().IntVarP(&chatListType, "type", "t", 0, "0=all, 1=private, 2=group")
	chatListCmd.Flags().StringVarP(&chatListKeyword, "keyword", "k", "", "Search keyword (only for type 1 or 2)")
	chatListCmd.Flags().Int64Var(&chatListStartTime, "start", 0, "Start time in microseconds")
	chatListCmd.Flags().Int64Var(&chatListEndTime, "end", 0, "End time in microseconds")
	chatListCmd.Flags().StringVar(&chatListUserToken, "user-token", "", "User token")

	chatMessagesCmd.Flags().StringVar(&chatMsgStaffID, "staff-id", "", "Private chat partner staffId")
	chatMessagesCmd.Flags().StringVar(&chatMsgGroupID, "group-id", "", "Group openId")
	chatMessagesCmd.Flags().IntVarP(&chatMsgPageSize, "size", "s", 100, "Per-page count (max 100)")
	chatMessagesCmd.Flags().StringVar(&chatMsgVersion, "version", "0", "Deep pagination cursor (first call: 0)")
	chatMessagesCmd.Flags().Int64Var(&chatMsgStartTime, "start", 0, "Start time in microseconds")
	chatMessagesCmd.Flags().Int64Var(&chatMsgEndTime, "end", 0, "End time in microseconds")
	chatMessagesCmd.Flags().StringVar(&chatMsgSenderID, "sender-id", "", "Filter by sender staffId")
	chatMessagesCmd.Flags().StringVar(&chatMsgUserToken, "user-token", "", "User token")
	chatMessagesCmd.Flags().BoolVar(&chatMsgSplitMonth, "split-month", false, "Auto-split query by month when time range exceeds 1 month")
	chatMessagesCmd.Flags().BoolVar(&chatMsgProgress, "progress", false, "Show pagination progress")

	chatCmd.AddCommand(chatListCmd)
	chatCmd.AddCommand(chatMessagesCmd)
	rootCmd.AddCommand(chatCmd)
}

func runChatList(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	chatTypeStr := strconv.Itoa(chatListType)
	startTimeStr := ""
	if chatListStartTime != 0 {
		startTimeStr = strconv.FormatInt(chatListStartTime, 10)
	}
	endTimeStr := ""
	if chatListEndTime != 0 {
		endTimeStr = strconv.FormatInt(chatListEndTime, 10)
	}

	result, err := client.FetchChatList(ctx, chatListUserToken, chatTypeStr, chatListKeyword, startTimeStr, endTimeStr)
	checkError(err)

	if jsonOutput {
		outputJSON(result)
		return
	}

	fmt.Printf("%-20s %s\n", "Field", "Value")
	fmt.Printf("%-20s %s\n", strings.Repeat("━", 20), strings.Repeat("━", 60))
	fmt.Printf("%-20s %s\n", "success", fmtVal(result.Success))

	if result.StaffInfos != nil {
		fmt.Println()
		fmt.Println("Staff Infos (Private Chats):")
		fmt.Printf("  %-20s %-30s %s\n", "Staff ID", "Name", "Sectors")
		fmt.Printf("  %-20s %-30s %s\n", strings.Repeat("─", 20), strings.Repeat("─", 30), strings.Repeat("─", 40))
		for _, s := range result.StaffInfos {
			sectors := strings.Join(s.SectorNames, ", ")
			fmt.Printf("  %-20s %-30s %s\n", s.StaffID, s.StaffName, sectors)
		}
	}

	if result.GroupInfos != nil {
		fmt.Println()
		fmt.Println("Group Infos (Group Chats):")
		fmt.Printf("  %-20s %s\n", "Group ID", "Name")
		fmt.Printf("  %-20s %s\n", strings.Repeat("─", 20), strings.Repeat("─", 40))
		for _, g := range result.GroupInfos {
			fmt.Printf("  %-20s %s\n", g.GroupID, g.GroupName)
		}
	}

	if !result.Success {
		fmt.Printf("%-20s %s\n", "error", result.Error)
	}
}

func runChatMessages(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	if chatMsgSplitMonth {
		runChatMessagesSplitMonth(client, ctx)
		return
	}

	startTimeStr := ""
	if chatMsgStartTime != 0 {
		startTimeStr = strconv.FormatInt(chatMsgStartTime, 10)
	}
	endTimeStr := ""
	if chatMsgEndTime != 0 {
		endTimeStr = strconv.FormatInt(chatMsgEndTime, 10)
	}

	result, err := client.FetchChatMessages(ctx, chatMsgUserToken, chatMsgPageSize, chatMsgVersion, chatMsgStaffID, chatMsgGroupID, startTimeStr, endTimeStr, chatMsgSenderID)
	checkError(err)

	if jsonOutput {
		outputJSON(result)
		return
	}

	outputResultFields(result, []string{"has_more", "total", "last_version", "name", "chat_type"})

	if result.Messages != nil {
		fmt.Println()
		fmt.Println("Messages:")
		fmt.Printf("  %-20s %-30s %s\n", "Time", "Sender", "Type")
		fmt.Printf("  %-20s %-30s %s\n", strings.Repeat("─", 20), strings.Repeat("─", 30), strings.Repeat("─", 20))
		for _, m := range result.Messages {
			fmt.Printf("  %-20s %-30s %s\n", m.SendTime, m.Sender, m.MessageType)
		}
	}
}

func splitMonths(startUs, endUs int64) [][2]int64 {
	startTime := time.Unix(startUs/1_000_000, (startUs%1_000_000)*1000)
	if startUs == 0 {
		startTime = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	}

	var intervals [][2]int64
	y := startTime.Year()
	m := startTime.Month()
	for {
		monthStart := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
		nextMonth := m + 1
		nextY := y
		if nextMonth > 12 {
			nextMonth = 1
			nextY++
		}
		monthEndLastDay := time.Date(nextY, nextMonth, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)

		msUs := monthStart.Unix() * 1_000_000
		meUs := monthEndLastDay.Unix() * 1_000_000

		if msUs > endUs {
			break
		}
		if meUs > endUs {
			meUs = endUs
		}
		if msUs < startUs {
			msUs = startUs
		}

		intervals = append(intervals, [2]int64{msUs, meUs})

		y = nextY
		m = nextMonth
	}
	return intervals
}

func runChatMessagesSplitMonth(client *lansenger.LansengerClient, ctx context.Context) {
	if chatMsgEndTime == 0 {
		fmt.Fprintf(os.Stderr, "Error: --end is required when using --split-month\n")
		os.Exit(1)
	}

	intervals := splitMonths(chatMsgStartTime, chatMsgEndTime)

	var allMessages []lansenger.ChatMessageInfo
	totalPages := 0

	for i, interval := range intervals {
		monthNum := i + 1
		startStr := strconv.FormatInt(interval[0], 10)
		endStr := strconv.FormatInt(interval[1], 10)
		cursor := "0"
		monthMsgCount := 0
		pages := 0

		for {
			result, err := client.FetchChatMessages(ctx, chatMsgUserToken, chatMsgPageSize, cursor, chatMsgStaffID, chatMsgGroupID, startStr, endStr, chatMsgSenderID)
			checkError(err)

			allMessages = append(allMessages, result.Messages...)
			monthMsgCount += len(result.Messages)
			pages++

			if chatMsgProgress && !jsonOutput {
				fmt.Printf("Month %d/%d | Page %d | %d messages total\n", monthNum, len(intervals), pages, monthMsgCount)
			}

			if !result.HasMore {
				break
			}
			cursor = result.LastVersion
		}
		totalPages += pages
	}

	if chatMsgProgress && !jsonOutput {
		fmt.Printf("Done: %d pages, %d messages across %d months\n", totalPages, len(allMessages), len(intervals))
	}

	if jsonOutput {
		outputJSON(allMessages)
		return
	}

	fmt.Printf("Total: %d messages across %d months (%d pages)\n", len(allMessages), len(intervals), totalPages)
	if allMessages != nil {
		fmt.Println()
		fmt.Println("Messages:")
		fmt.Printf("  %-20s %-30s %s\n", "Time", "Sender", "Type")
		fmt.Printf("  %-20s %-30s %s\n", strings.Repeat("─", 20), strings.Repeat("─", 30), strings.Repeat("─", 20))
		for _, m := range allMessages {
			fmt.Printf("  %-20s %-30s %s\n", m.SendTime, m.Sender, m.MessageType)
		}
	}
}