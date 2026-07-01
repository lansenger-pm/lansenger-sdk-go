package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var messageApproveCardCmd = &cobra.Command{
	Use:   "send-approve-card BODY_TITLE BODY_CONTENT",
	Short: "Send an approve card message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendApproveCard,
}

var messageUpdateApproveCardCmd = &cobra.Command{
	Use:   "update-approve-card MSG_ID",
	Short: "Update an approve card message",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdateApproveCard,
}

var (
	sendApproveCardChatID           string
	sendApproveCardHeadTitle        string
	sendApproveCardHeadIconLink     string
	sendApproveCardHeadIconID       string
	sendApproveCardHeadStatus       string
	sendApproveCardHeadStatusIcon   int
	sendApproveCardHeadStatusIconLink string
	sendApproveCardHeadStatusColour string
	sendApproveCardBodyFormatType   int
	sendApproveCardFields           string
	sendApproveCardCardLink         string
	sendApproveCardCardLinkPC       string
	sendApproveCardCardLinkPad      string
	sendApproveCardButtons          string
	sendApproveCardExpireTime       int
	sendApproveCardReminderAll      bool
	sendApproveCardReminderUserIDs  []string
	sendApproveCardReminderBotIDs   []string
	sendApproveCardIsGroup          bool
	sendApproveCardUserToken        string
	sendApproveCardSenderID         string

	updateApproveCardHeadStatus       string
	updateApproveCardHeadStatusIcon   int
	updateApproveCardHeadStatusIconLink string
	updateApproveCardHeadStatusColour string
	updateApproveCardButtons          string
)

func init() {
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardChatID, "chat-id", "", "Chat ID (required)")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardHeadTitle, "head-title", "", "Card head title")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardHeadIconLink, "head-icon-link", "", "Head icon URL")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardHeadIconID, "head-icon-id", "", "Head icon ID")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardHeadStatus, "head-status", "", "Head status description (div-style HTML)")
	messageApproveCardCmd.Flags().IntVar(&sendApproveCardHeadStatusIcon, "head-status-icon", 0, "Head status icon index")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardHeadStatusIconLink, "head-status-icon-link", "", "Head status icon URL")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardHeadStatusColour, "head-status-colour", "", "Head status DOT colour (hex)")
	messageApproveCardCmd.Flags().IntVar(&sendApproveCardBodyFormatType, "format-type", 0, "Body format type")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardFields, "fields", "", "Body fields as JSON array, e.g. '[{\"key\":\"k\",\"value\":\"v\"}]'")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardCardLink, "card-link", "", "Card link URL")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardCardLinkPC, "card-link-pc", "", "PC card link URL")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardCardLinkPad, "card-link-pad", "", "Pad card link URL")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardButtons, "buttons", "", "Buttons as JSON array, e.g. '[{\"text\":\"Approve\",\"value\":1}]'")
	messageApproveCardCmd.Flags().IntVar(&sendApproveCardExpireTime, "expire-time", 0, "Expire time in seconds")
	messageApproveCardCmd.Flags().BoolVar(&sendApproveCardReminderAll, "mention-all", false, "@all")
	messageApproveCardCmd.Flags().StringArrayVar(&sendApproveCardReminderUserIDs, "mention", nil, "User IDs to @mention")
	messageApproveCardCmd.Flags().StringArrayVar(&sendApproveCardReminderBotIDs, "mention-bot", nil, "Reminder bot IDs")
	messageApproveCardCmd.Flags().BoolVarP(&sendApproveCardIsGroup, "group", "g", false, "Send as group message")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardUserToken, "user-token", "", "User token for private channel")
	messageApproveCardCmd.Flags().StringVar(&sendApproveCardSenderID, "sender-id", "", "Sender staff ID for group message")

	messageUpdateApproveCardCmd.Flags().StringVar(&updateApproveCardHeadStatus, "head-status", "", "New head status description (div-style HTML)")
	messageUpdateApproveCardCmd.Flags().IntVar(&updateApproveCardHeadStatusIcon, "head-status-icon", 0, "New head status icon index")
	messageUpdateApproveCardCmd.Flags().StringVar(&updateApproveCardHeadStatusIconLink, "head-status-icon-link", "", "New head status icon URL")
	messageUpdateApproveCardCmd.Flags().StringVar(&updateApproveCardHeadStatusColour, "head-status-colour", "", "New head status DOT colour (hex)")
	messageUpdateApproveCardCmd.Flags().StringVar(&updateApproveCardButtons, "buttons", "", "Updated buttons as JSON array")
}

func runSendApproveCard(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	if sendApproveCardChatID == "" {
		fmt.Fprintf(os.Stderr, "Error: --chat-id is required\n")
		return
	}

	var fields []map[string]string
	if sendApproveCardFields != "" {
		if err := json.Unmarshal([]byte(sendApproveCardFields), &fields); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --fields JSON: %v\n", err)
			return
		}
	}

	var buttons []map[string]interface{}
	if sendApproveCardButtons != "" {
		if err := json.Unmarshal([]byte(sendApproveCardButtons), &buttons); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --buttons JSON: %v\n", err)
			return
		}
	}

	params := &lansenger.ApproveCardParams{
		ChatID:             sendApproveCardChatID,
		BodyTitle:          args[0],
		BodyContent:        args[1],
		HeadTitle:          sendApproveCardHeadTitle,
		HeadIconLink:       sendApproveCardHeadIconLink,
		HeadIconID:         sendApproveCardHeadIconID,
		HeadStatusDescribe: sendApproveCardHeadStatus,
		HeadStatusIcon:     sendApproveCardHeadStatusIcon,
		HeadStatusIconLink: sendApproveCardHeadStatusIconLink,
		HeadStatusColour:   sendApproveCardHeadStatusColour,
		BodyFormatType:     sendApproveCardBodyFormatType,
		Fields:             fields,
		ReminderAll:        sendApproveCardReminderAll,
		ReminderUserIDs:    sendApproveCardReminderUserIDs,
		ReminderBotIDs:     sendApproveCardReminderBotIDs,
		CardLink:           sendApproveCardCardLink,
		CardLinkForPC:      sendApproveCardCardLinkPC,
		CardLinkForPad:     sendApproveCardCardLinkPad,
		Buttons:            buttons,
		ExpireTime:         sendApproveCardExpireTime,
		IsGroup:            sendApproveCardIsGroup,
		UserToken:          sendApproveCardUserToken,
		SenderID:           sendApproveCardSenderID,
	}

	result, err := client.SendApproveCardWithParams(ctx, params)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runUpdateApproveCard(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var buttons []map[string]interface{}
	if updateApproveCardButtons != "" {
		if err := json.Unmarshal([]byte(updateApproveCardButtons), &buttons); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing --buttons JSON: %v\n", err)
			return
		}
	}

	params := &lansenger.ApproveCardUpdateParams{
		MsgID:              args[0],
		HeadStatusDescribe: updateApproveCardHeadStatus,
		HeadStatusIcon:     updateApproveCardHeadStatusIcon,
		HeadStatusIconLink: updateApproveCardHeadStatusIconLink,
		HeadStatusColour:   updateApproveCardHeadStatusColour,
		Buttons:            buttons,
	}

	result, err := client.UpdateApproveCard(ctx, params)
	checkError(err)
	outputResultFields(result, []string{"message_id", "operation"})
}
