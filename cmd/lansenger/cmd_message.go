package main

import (
	"context"
	"fmt"
	"os"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var messageCmd = &cobra.Command{
	Use:   "message",
	Short: "Send and manage messages",
}

var sendTextCmd = &cobra.Command{
	Use:   "send-text CHAT_ID CONTENT",
	Short: "Send a text message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendText,
}

var sendMarkdownCmd = &cobra.Command{
	Use:   "send-markdown CHAT_ID CONTENT",
	Short: "Send a markdown message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendMarkdown,
}

var sendFileCmd = &cobra.Command{
	Use:   "send-file CHAT_ID FILE_PATH",
	Short: "Send a file message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendFile,
}

var sendImageURLCmd = &cobra.Command{
	Use:   "send-image-url CHAT_ID IMAGE_URL",
	Short: "Send an image from URL",
	Args:  cobra.ExactArgs(2),
	Run:   runSendImageURL,
}

var sendLinkCardCmd = &cobra.Command{
	Use:   "send-link-card CHAT_ID TITLE LINK",
	Short: "Send a link card message",
	Args:  cobra.ExactArgs(3),
	Run:   runSendLinkCard,
}

var sendAppArticlesCmd = &cobra.Command{
	Use:   "send-app-articles CHAT_ID ARTICLES...",
	Short: "Send app articles message",
	Args:  cobra.MinimumNArgs(2),
	Run:   runSendAppArticles,
}

var sendAppCardCmd = &cobra.Command{
	Use:   "send-app-card CHAT_ID BODY_TITLE",
	Short: "Send an app card message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendAppCard,
}

var sendOaCardCmd = &cobra.Command{
	Use:   "send-oacard CHAT_ID TITLE",
	Short: "Send an OA card message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendOaCard,
}

var updateDynamicCardCmd = &cobra.Command{
	Use:   "update-dynamic-card MSG_ID",
	Short: "Update a dynamic card message",
	Args:  cobra.ExactArgs(1),
	Run:   runUpdateDynamicCard,
}

var revokeCmd = &cobra.Command{
	Use:   "revoke MESSAGE_IDS...",
	Short: "Revoke messages",
	Args:  cobra.MinimumNArgs(1),
	Run:   runRevoke,
}

var sendReminderCmd = &cobra.Command{
	Use:   "send-reminder MSG_ID",
	Short: "Send urgent reminder for a message",
	Args:  cobra.ExactArgs(1),
	Run:   runSendReminder,
}

var sendBotMessageCmd = &cobra.Command{
	Use:   "send-bot-message MSG_TYPE MSG_DATA",
	Short: "Send a bot channel message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendBotMessage,
}

var sendGroupMessageCmd = &cobra.Command{
	Use:   "send-group-message GROUP_ID MSG_TYPE MSG_DATA",
	Short: "Send a group channel message",
	Args:  cobra.ExactArgs(3),
	Run:   runSendGroupMessage,
}

var sendAccountMessageCmd = &cobra.Command{
	Use:   "send-account-message MSG_TYPE MSG_DATA",
	Short: "Send an account channel message",
	Args:  cobra.ExactArgs(2),
	Run:   runSendAccountMessage,
}

var sendUserMessageCmd = &cobra.Command{
	Use:   "send-user-message RECEIVER_ID MSG_TYPE MSG_DATA",
	Short: "Send a user channel message",
	Args:  cobra.ExactArgs(3),
	Run:   runSendUserMessage,
}

var queryGroupsCmd = &cobra.Command{
	Use:   "query-groups",
	Short: "Query groups",
	Args:  cobra.NoArgs,
	Run:   runQueryGroups,
}

var (
	sendReminderTypes    []int
	sendReminderUserIDs  []string

	sendTextFile            string
	sendTextMediaType       string
	sendTextCoverImage      string
	sendTextIsGroup         bool
	sendTextReminderAll     bool
	sendTextReminderUserIDs []string
	sendTextRefMsgID        string
	sendTextMentionBotIDs   []string
	sendTextUserToken       string
	sendTextSenderID        string

	sendMarkdownIsGroup         bool
	sendMarkdownReminderAll     bool
	sendMarkdownReminderUserIDs []string
	sendMarkdownRefMsgID        string
	sendMarkdownMentionBotIDs   []string
	sendMarkdownUserToken       string
	sendMarkdownSenderID        string

	sendFileContent     string
	sendFileMediaType  string
	sendFileCoverImage string
	sendFileIsGroup    bool
	sendFileUserToken  string
	sendFileSenderID   string

	sendImageURLContent   string
	sendImageURLIsGroup   bool
	sendImageURLUserToken string
	sendImageURLSenderID  string

	sendLinkCardDesc       string
	sendLinkCardIcon       string
	sendLinkCardPcLink     string
	sendLinkCardPadLink    string
	sendLinkCardFromName   string
	sendLinkCardFromIcon   string
	sendLinkCardIsGroup    bool
	sendLinkCardUserToken  string
	sendLinkCardSenderID   string

	sendAppArticlesIsGroup   bool
	sendAppArticlesUserToken string
	sendAppArticlesSenderID  string

	sendAppCardHeadTitle    string
	sendAppCardSubTitle     string
	sendAppCardContent      string
	sendAppCardSignature    string
	sendAppCardCardLink     string
	sendAppCardPcCardLink   string
	sendAppCardPadCardLink  string
	sendAppCardIsDynamic    bool
	sendAppCardStaffID      string
	sendAppCardHeadIcon     string
	sendAppCardStatusDesc   string
	sendAppCardStatusColour string
	sendAppCardFields       []string
	sendAppCardLinks        []string
	sendAppCardIsGroup      bool
	sendAppCardUserToken    string
	sendAppCardSenderID     string

	sendOaCardHead       string
	sendOaCardSubTitle   string
	sendOaCardStaffID    string
	sendOaCardFields     []string
	sendOaCardLink       string
	sendOaCardPcLink     string
	sendOaCardPadLink    string
	sendOaCardCardAction string
	sendOaCardIsGroup    bool
	sendOaCardUserToken  string
	sendOaCardSenderID   string

	updateDynamicCardLast         bool
	updateDynamicCardStatusDesc   string
	updateDynamicCardStatusColour string
	updateDynamicCardLinks        []string

	revokeChatType string
	revokeSenderID string

	sendBotMessageChatIDs       []string
	sendBotMessageDepartmentIDs []string
	sendBotMessageUserToken     string
	sendBotMessageEntryID       string
	sendBotMessageRefMsgID      string
	sendBotMessageIsGroup       bool

	sendGroupMessageUserToken       string
	sendGroupMessageSenderID        string
	sendGroupMessageReminderAll     bool
	sendGroupMessageReminderUserIDs []string
	sendGroupMessageRefMsgID        string
	sendGroupMessageMentionBotIDs   []string
	sendGroupMessageOutlines        string
	sendGroupMessageEntryID         string
	sendGroupMessageUUID            string

	sendAccountMessageChatIDs       []string
	sendAccountMessageDepartmentIDs []string
	sendAccountMessageAccountID     string
	sendAccountMessageEntryID       string
	sendAccountMessageAttach        string
	sendAccountMessageUserToken     string

	sendUserMessageUserToken string
	sendUserMessageCommon    string
	sendUserMessageUUID      string

	queryGroupsPageOffset int
	queryGroupsPageSize   int
)

func init() {
	sendTextCmd.Flags().StringVarP(&sendTextFile, "file", "f", "", "File path to attach")
	sendTextCmd.Flags().StringVarP(&sendTextMediaType, "media-type", "t", "", "file/video/image/audio")
	sendTextCmd.Flags().StringVar(&sendTextCoverImage, "cover-image", "", "Cover image path (required for video)")
	sendTextCmd.Flags().BoolVarP(&sendTextIsGroup, "group", "g", false, "Send as group message")
	sendTextCmd.Flags().BoolVar(&sendTextReminderAll, "mention-all", false, "@all in group")
	sendTextCmd.Flags().StringArrayVar(&sendTextReminderUserIDs, "mention", nil, "User IDs to @mention")
	sendTextCmd.Flags().StringVar(&sendTextRefMsgID, "ref-msg-id", "", "Reference message ID for reply")
	sendTextCmd.Flags().StringArrayVar(&sendTextMentionBotIDs, "mention-bot", nil, "Reminder bot IDs")
	sendTextCmd.Flags().StringVar(&sendTextUserToken, "user-token", "", "User token for private channel")
	sendTextCmd.Flags().StringVar(&sendTextSenderID, "sender-id", "", "Sender staff ID for group message")

	sendMarkdownCmd.Flags().BoolVarP(&sendMarkdownIsGroup, "group", "g", false, "Send as group message")
	sendMarkdownCmd.Flags().BoolVar(&sendMarkdownReminderAll, "mention-all", false, "@all in group")
	sendMarkdownCmd.Flags().StringArrayVar(&sendMarkdownReminderUserIDs, "mention", nil, "User IDs to @mention")
	sendMarkdownCmd.Flags().StringVar(&sendMarkdownRefMsgID, "ref-msg-id", "", "Reference message ID for reply")
	sendMarkdownCmd.Flags().StringArrayVar(&sendMarkdownMentionBotIDs, "mention-bot", nil, "Reminder bot IDs")
	sendMarkdownCmd.Flags().StringVar(&sendMarkdownUserToken, "user-token", "", "User token for private channel")
	sendMarkdownCmd.Flags().StringVar(&sendMarkdownSenderID, "sender-id", "", "Sender staff ID for group message")

	sendFileCmd.Flags().StringVarP(&sendFileContent, "content", "c", "", "Content/caption text")
	sendFileCmd.Flags().StringVar(&sendFileMediaType, "media-type", "", "file/video/image/audio")
	sendFileCmd.Flags().StringVar(&sendFileCoverImage, "cover-image", "", "Cover image path (required for video)")
	sendFileCmd.Flags().BoolVarP(&sendFileIsGroup, "group", "g", false, "Send as group message")
	sendFileCmd.Flags().StringVar(&sendFileUserToken, "user-token", "", "User token for private channel")
	sendFileCmd.Flags().StringVar(&sendFileSenderID, "sender-id", "", "Sender staff ID for group message")

	sendImageURLCmd.Flags().StringVarP(&sendImageURLContent, "content", "c", "", "Content/caption text")
	sendImageURLCmd.Flags().BoolVarP(&sendImageURLIsGroup, "group", "g", false, "Send as group message")
	sendImageURLCmd.Flags().StringVar(&sendImageURLUserToken, "user-token", "", "User token for private channel")
	sendImageURLCmd.Flags().StringVar(&sendImageURLSenderID, "sender-id", "", "Sender staff ID for group message")

	sendLinkCardCmd.Flags().StringVarP(&sendLinkCardDesc, "desc", "d", "", "Card description")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardIcon, "icon", "", "Icon URL")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardPcLink, "pc-link", "", "PC link URL")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardPadLink, "pad-link", "", "Pad link URL")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardFromName, "from-name", "", "Source name")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardFromIcon, "from-icon", "", "Source icon URL")
	sendLinkCardCmd.Flags().BoolVarP(&sendLinkCardIsGroup, "group", "g", false, "Send as group message")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardUserToken, "user-token", "", "User token for private channel")
	sendLinkCardCmd.Flags().StringVar(&sendLinkCardSenderID, "sender-id", "", "Sender staff ID for group message")

	sendAppArticlesCmd.Flags().BoolVarP(&sendAppArticlesIsGroup, "group", "g", false, "Send as group message")
	sendAppArticlesCmd.Flags().StringVar(&sendAppArticlesUserToken, "user-token", "", "User token for private channel")
	sendAppArticlesCmd.Flags().StringVar(&sendAppArticlesSenderID, "sender-id", "", "Sender staff ID for group message")

	sendAppCardCmd.Flags().StringVar(&sendAppCardHeadTitle, "head-title", "", "Card head title")
	sendAppCardCmd.Flags().StringVar(&sendAppCardSubTitle, "sub-title", "", "Card sub title")
	sendAppCardCmd.Flags().StringVar(&sendAppCardContent, "content", "", "Card body content (supports div-style HTML)")
	sendAppCardCmd.Flags().StringVar(&sendAppCardSignature, "signature", "", "Card signature")
	sendAppCardCmd.Flags().StringVar(&sendAppCardCardLink, "card-link", "", "Card link URL")
	sendAppCardCmd.Flags().StringVar(&sendAppCardPcCardLink, "pc-card-link", "", "PC card link URL")
	sendAppCardCmd.Flags().StringVar(&sendAppCardPadCardLink, "pad-card-link", "", "Pad card link URL")
	sendAppCardCmd.Flags().BoolVar(&sendAppCardIsDynamic, "dynamic", false, "Enable dynamic card updates")
	sendAppCardCmd.Flags().StringVar(&sendAppCardStaffID, "staff-id", "", "Staff ID")
	sendAppCardCmd.Flags().StringVar(&sendAppCardHeadIcon, "head-icon", "", "Head icon URL")
	sendAppCardCmd.Flags().StringVar(&sendAppCardStatusDesc, "status-desc", "", "Head status description (div-style HTML, max 30 bytes)")
	sendAppCardCmd.Flags().StringVar(&sendAppCardStatusColour, "status-colour", "", "Head status DOT colour (hex, e.g. #FFB116)")
	sendAppCardCmd.Flags().StringArrayVar(&sendAppCardFields, "field", nil, "Card field as JSON, e.g. '{\"key\":\"k\",\"value\":\"v\"}'")
	sendAppCardCmd.Flags().StringArrayVar(&sendAppCardLinks, "link", nil, "Card link as JSON, e.g. '{\"title\":\"T\",\"url\":\"U\"}'")
	sendAppCardCmd.Flags().BoolVarP(&sendAppCardIsGroup, "group", "g", false, "Send as group message")
	sendAppCardCmd.Flags().StringVar(&sendAppCardUserToken, "user-token", "", "User token for private channel")
	sendAppCardCmd.Flags().StringVar(&sendAppCardSenderID, "sender-id", "", "Sender staff ID for group message")

	sendOaCardCmd.Flags().StringVar(&sendOaCardHead, "head", "", "OA card head title")
	sendOaCardCmd.Flags().StringVar(&sendOaCardSubTitle, "sub-title", "", "OA card sub title")
	sendOaCardCmd.Flags().StringVar(&sendOaCardStaffID, "staff-id", "", "Staff ID")
	sendOaCardCmd.Flags().StringArrayVar(&sendOaCardFields, "field", nil, "Card field as JSON, e.g. '{\"key\":\"k\",\"value\":\"v\"}'")
	sendOaCardCmd.Flags().StringVar(&sendOaCardLink, "link", "", "Card click link URL")
	sendOaCardCmd.Flags().StringVar(&sendOaCardPcLink, "pc-link", "", "PC link URL")
	sendOaCardCmd.Flags().StringVar(&sendOaCardPadLink, "pad-link", "", "Pad link URL")
	sendOaCardCmd.Flags().StringVar(&sendOaCardCardAction, "card-action", "", "Card action as JSON dict")
	sendOaCardCmd.Flags().BoolVarP(&sendOaCardIsGroup, "group", "g", false, "Send as group message")
	sendOaCardCmd.Flags().StringVar(&sendOaCardUserToken, "user-token", "", "User token for private channel")
	sendOaCardCmd.Flags().StringVar(&sendOaCardSenderID, "sender-id", "", "Sender staff ID for group message")

	updateDynamicCardCmd.Flags().BoolVar(&updateDynamicCardLast, "last", false, "Mark as last update")
	updateDynamicCardCmd.Flags().StringVar(&updateDynamicCardStatusDesc, "status-desc", "", "New status description (div-style HTML, max 30 bytes)")
	updateDynamicCardCmd.Flags().StringVar(&updateDynamicCardStatusColour, "status-colour", "", "New status DOT colour (hex)")
	updateDynamicCardCmd.Flags().StringArrayVar(&updateDynamicCardLinks, "link", nil, "Updated link as JSON title=url")

	sendReminderCmd.Flags().IntSliceVarP(&sendReminderTypes, "type", "t", nil, "Reminder types (1=app, 2=sms)")
	sendReminderCmd.Flags().StringArrayVarP(&sendReminderUserIDs, "user", "u", nil, "User IDs (staff openIds) to remind")

	revokeCmd.Flags().StringVar(&revokeChatType, "chat-type", "bot", "staff, group, notification, account, or bot")
	revokeCmd.Flags().StringVar(&revokeSenderID, "sender-id", "", "Sender staff ID (required for staff/group)")

	sendBotMessageCmd.Flags().StringArrayVar(&sendBotMessageChatIDs, "chat-id", nil, "Chat IDs (or group IDs if --group)")
	sendBotMessageCmd.Flags().StringArrayVar(&sendBotMessageDepartmentIDs, "dept", nil, "Department IDs (bot channel only)")
	sendBotMessageCmd.Flags().StringVar(&sendBotMessageUserToken, "user-token", "", "User token")
	sendBotMessageCmd.Flags().StringVar(&sendBotMessageEntryID, "entry-id", "", "App entry selector")
	sendBotMessageCmd.Flags().StringVar(&sendBotMessageRefMsgID, "ref-msg-id", "", "Reference message ID for reply")
	sendBotMessageCmd.Flags().BoolVarP(&sendBotMessageIsGroup, "group", "g", false, "Send to groups instead of users (uses /v1/messages/group/create)")

	sendGroupMessageCmd.Flags().StringVar(&sendGroupMessageUserToken, "user-token", "", "User token")
	sendGroupMessageCmd.Flags().StringVar(&sendGroupMessageSenderID, "sender-id", "", "Sender staff ID")
	sendGroupMessageCmd.Flags().BoolVar(&sendGroupMessageReminderAll, "mention-all", false, "@all (text/formatText only)")
	sendGroupMessageCmd.Flags().StringArrayVar(&sendGroupMessageReminderUserIDs, "mention", nil, "User IDs to @mention (text/formatText only)")
	sendGroupMessageCmd.Flags().StringVar(&sendGroupMessageRefMsgID, "ref-msg-id", "", "Reference message ID for reply")
	sendGroupMessageCmd.Flags().StringArrayVar(&sendGroupMessageMentionBotIDs, "mention-bot", nil, "Reminder bot IDs")
	sendGroupMessageCmd.Flags().StringVar(&sendGroupMessageOutlines, "outlines", "", "Group notification digest")
	sendGroupMessageCmd.Flags().StringVar(&sendGroupMessageEntryID, "entry-id", "", "App entry selector")
	sendGroupMessageCmd.Flags().StringVar(&sendGroupMessageUUID, "uuid", "", "Message UUID for deduplication")

	sendAccountMessageCmd.Flags().StringArrayVar(&sendAccountMessageChatIDs, "chat-id", nil, "Chat IDs")
	sendAccountMessageCmd.Flags().StringArrayVar(&sendAccountMessageDepartmentIDs, "dept", nil, "Department IDs")
	sendAccountMessageCmd.Flags().StringVar(&sendAccountMessageAccountID, "account-id", "", "Account ID")
	sendAccountMessageCmd.Flags().StringVar(&sendAccountMessageEntryID, "entry-id", "", "App entry selector")
	sendAccountMessageCmd.Flags().StringVar(&sendAccountMessageAttach, "attach", "", "Attach info")
	sendAccountMessageCmd.Flags().StringVar(&sendAccountMessageUserToken, "user-token", "", "User token")

	sendUserMessageCmd.Flags().StringVar(&sendUserMessageUserToken, "user-token", "", "User token")
	sendUserMessageCmd.Flags().StringVar(&sendUserMessageCommon, "common", "", "Common data as JSON dict")
	sendUserMessageCmd.Flags().StringVar(&sendUserMessageUUID, "uuid", "", "Deduplication UUID")

	queryGroupsCmd.Flags().IntVarP(&queryGroupsPageOffset, "page", "p", 0, "Page offset (starts from 0)")
	queryGroupsCmd.Flags().IntVarP(&queryGroupsPageSize, "size", "s", 100, "Page size")

	messageCmd.AddCommand(sendTextCmd)
	messageCmd.AddCommand(sendMarkdownCmd)
	messageCmd.AddCommand(sendFileCmd)
	messageCmd.AddCommand(sendImageURLCmd)
	messageCmd.AddCommand(sendLinkCardCmd)
	messageCmd.AddCommand(sendAppArticlesCmd)
	messageCmd.AddCommand(sendAppCardCmd)
	messageCmd.AddCommand(sendOaCardCmd)
	messageCmd.AddCommand(updateDynamicCardCmd)
	messageCmd.AddCommand(sendReminderCmd)
	messageCmd.AddCommand(revokeCmd)
	messageCmd.AddCommand(sendBotMessageCmd)
	messageCmd.AddCommand(sendGroupMessageCmd)
	messageCmd.AddCommand(sendAccountMessageCmd)
	messageCmd.AddCommand(sendUserMessageCmd)
	messageCmd.AddCommand(queryGroupsCmd)
	rootCmd.AddCommand(messageCmd)
}

func runSendText(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var reminderUserIDs []string
	if len(sendTextReminderUserIDs) > 0 {
		reminderUserIDs = sendTextReminderUserIDs
	}

	result, err := client.SendText(ctx, args[0], args[1], sendTextFile, sendTextMediaType, sendTextCoverImage, sendTextReminderAll, reminderUserIDs, sendTextMentionBotIDs, sendTextIsGroup, sendTextUserToken, sendTextSenderID, sendTextRefMsgID)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendMarkdown(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var reminderUserIDs []string
	if len(sendMarkdownReminderUserIDs) > 0 {
		reminderUserIDs = sendMarkdownReminderUserIDs
	}

	result, err := client.SendMarkdown(ctx, args[0], args[1], sendMarkdownReminderAll, reminderUserIDs, sendMarkdownMentionBotIDs, sendMarkdownIsGroup, sendMarkdownUserToken, sendMarkdownSenderID, sendMarkdownRefMsgID)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendFile(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.SendFile(ctx, args[0], args[1], sendFileContent, sendFileMediaType, sendFileCoverImage, sendFileIsGroup, sendFileUserToken, sendFileSenderID)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendImageURL(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.SendImageURL(ctx, args[0], args[1], sendImageURLContent, sendImageURLIsGroup, sendImageURLUserToken, sendImageURLSenderID)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendLinkCard(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	params := &lansenger.LinkCardParams{
		ChatID:       args[0],
		Title:        args[1],
		Link:         args[2],
		Description:  sendLinkCardDesc,
		IconLink:     sendLinkCardIcon,
		PcLink:       sendLinkCardPcLink,
		PadLink:      sendLinkCardPadLink,
		FromName:     sendLinkCardFromName,
		FromIconLink: sendLinkCardFromIcon,
		IsGroup:      sendLinkCardIsGroup,
		UserToken:    sendLinkCardUserToken,
		SenderID:     sendLinkCardSenderID,
	}

	result, err := client.SendLinkCardWithParams(ctx, params)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendAppArticles(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	articles := make([]map[string]string, 0, len(args)-1)
	for i := 1; i < len(args); i++ {
		m, err := parseJSONMap(args[i])
		checkError(err)
		articles = append(articles, m)
	}

	result, err := client.SendAppArticles(ctx, args[0], articles, sendAppArticlesIsGroup, sendAppArticlesUserToken, sendAppArticlesSenderID)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendAppCard(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var headStatusInfo map[string]interface{}
	if sendAppCardStatusDesc != "" || sendAppCardStatusColour != "" {
		headStatusInfo = map[string]interface{}{
			"description": sendAppCardStatusDesc,
			"colour":      sendAppCardStatusColour,
		}
	}

	var parsedFields []map[string]interface{}
	if len(sendAppCardFields) > 0 {
		parsedFields = make([]map[string]interface{}, 0, len(sendAppCardFields))
		for _, f := range sendAppCardFields {
			m, err := parseJSONMap(f)
			checkError(err)
			parsed := make(map[string]interface{}, len(m))
			for k, v := range m {
				parsed[k] = v
			}
			parsedFields = append(parsedFields, parsed)
		}
	}

	var parsedLinks []map[string]interface{}
	if len(sendAppCardLinks) > 0 {
		parsedLinks = make([]map[string]interface{}, 0, len(sendAppCardLinks))
		for _, l := range sendAppCardLinks {
			m, err := parseJSONMap(l)
			checkError(err)
			parsed := make(map[string]interface{}, len(m))
			for k, v := range m {
				parsed[k] = v
			}
			parsedLinks = append(parsedLinks, parsed)
		}
	}

	params := &lansenger.AppCardParams{
		ChatID:         args[0],
		BodyTitle:      args[1],
		HeadTitle:      sendAppCardHeadTitle,
		BodySubTitle:   sendAppCardSubTitle,
		BodyContent:    sendAppCardContent,
		Signature:      sendAppCardSignature,
		CardLink:       sendAppCardCardLink,
		PcCardLink:     sendAppCardPcCardLink,
		PadCardLink:    sendAppCardPadCardLink,
		IsDynamic:      sendAppCardIsDynamic,
		StaffID:        sendAppCardStaffID,
		HeadIconURL:    sendAppCardHeadIcon,
		HeadStatusInfo: headStatusInfo,
		Fields:         parsedFields,
		Links:          parsedLinks,
		IsGroup:        sendAppCardIsGroup,
		UserToken:      sendAppCardUserToken,
		SenderID:       sendAppCardSenderID,
	}

	result, err := client.SendAppCardWithParams(ctx, params)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runSendOaCard(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var parsedFields []map[string]interface{}
	if len(sendOaCardFields) > 0 {
		parsedFields = make([]map[string]interface{}, 0, len(sendOaCardFields))
		for _, f := range sendOaCardFields {
			m, err := parseJSONMap(f)
			checkError(err)
			parsed := make(map[string]interface{}, len(m))
			for k, v := range m {
				parsed[k] = v
			}
			parsedFields = append(parsedFields, parsed)
		}
	}

	var cardAction map[string]interface{}
	if sendOaCardCardAction != "" {
		action, err := parseJSONMap(sendOaCardCardAction)
		checkError(err)
		cardAction = make(map[string]interface{}, len(action))
		for k, v := range action {
			cardAction[k] = v
		}
	}

	params := &lansenger.OaCardParams{
		ChatID:     args[0],
		Title:      args[1],
		Head:       sendOaCardHead,
		SubTitle:   sendOaCardSubTitle,
		StaffID:    sendOaCardStaffID,
		Fields:     parsedFields,
		Link:       sendOaCardLink,
		PcLink:     sendOaCardPcLink,
		PadLink:    sendOaCardPadLink,
		CardAction: cardAction,
		IsGroup:    sendOaCardIsGroup,
		UserToken:  sendOaCardUserToken,
		SenderID:   sendOaCardSenderID,
	}

	result, err := client.SendOaCardWithParams(ctx, params)
	checkError(err)
	outputResultFields(result, []string{"message_id", "msg_type", "operation"})
}

func runUpdateDynamicCard(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var headStatusInfo map[string]interface{}
	if updateDynamicCardStatusDesc != "" || updateDynamicCardStatusColour != "" {
		headStatusInfo = map[string]interface{}{
			"description": updateDynamicCardStatusDesc,
			"colour":      updateDynamicCardStatusColour,
		}
	}

	var parsedLinks []map[string]interface{}
	if len(updateDynamicCardLinks) > 0 {
		parsedLinks = make([]map[string]interface{}, 0, len(updateDynamicCardLinks))
		for _, l := range updateDynamicCardLinks {
			m, err := parseJSONMap(l)
			checkError(err)
			parsed := make(map[string]interface{}, len(m))
			for k, v := range m {
				parsed[k] = v
			}
			parsedLinks = append(parsedLinks, parsed)
		}
	}

	params := &lansenger.DynamicCardUpdateParams{
		MsgID:          args[0],
		HeadStatusInfo: headStatusInfo,
		Links:          parsedLinks,
		IsLastUpdate:   updateDynamicCardLast,
	}

	result, err := client.UpdateDynamicCard(ctx, params)
	checkError(err)
	outputResultFields(result, []string{"message_id", "operation"})
}

func runRevoke(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.RevokeMessage(ctx, args, revokeChatType, revokeSenderID, nil)
	checkError(err)
	outputResultFields(result, []string{"message_id", "operation"})
}

func runSendBotMessage(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	msgData, err := parseJSONRaw(args[1])
	checkError(err)

	msgDataMap, ok := msgData.(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: msg_data must be a JSON object\n")
		return
	}

	if sendBotMessageIsGroup {
		// Route to group message endpoint (API 4.6.2 /v1/messages/group/create) for each chat ID
		for _, gid := range sendBotMessageChatIDs {
			result, err := client.SendGroupMessage(ctx, gid, args[0], msgDataMap,
				sendBotMessageUserToken, "", "", "", sendBotMessageEntryID, sendBotMessageRefMsgID)
			checkError(err)
			outputResultFields(result, []string{"message_id"})
		}
	} else {
		result, err := client.SendBotMessage(ctx, args[0], msgDataMap, sendBotMessageChatIDs, sendBotMessageDepartmentIDs, sendBotMessageUserToken, sendBotMessageEntryID, sendBotMessageRefMsgID)
		checkError(err)
		outputResultFields(result, []string{"message_id"})
	}
}

func runSendGroupMessage(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	msgData, err := parseJSONRaw(args[2])
	checkError(err)

	msgDataMap, ok := msgData.(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: msg_data must be a JSON object\n")
		return
	}

	// Inject reminder into msgData if mention flags are set
	if sendGroupMessageReminderAll || len(sendGroupMessageReminderUserIDs) > 0 || len(sendGroupMessageMentionBotIDs) > 0 {
		if textData, ok := msgDataMap["text"].(map[string]interface{}); ok {
			reminder := map[string]interface{}{
				"all":     sendGroupMessageReminderAll,
				"userIds": sendGroupMessageReminderUserIDs,
			}
			if len(sendGroupMessageMentionBotIDs) > 0 {
				reminder["botIds"] = sendGroupMessageMentionBotIDs
			}
			textData["reminder"] = reminder
		} else if ftData, ok := msgDataMap["formatText"].(map[string]interface{}); ok {
			reminder := map[string]interface{}{
				"all":     sendGroupMessageReminderAll,
				"userIds": sendGroupMessageReminderUserIDs,
			}
			if len(sendGroupMessageMentionBotIDs) > 0 {
				reminder["botIds"] = sendGroupMessageMentionBotIDs
			}
			ftData["reminder"] = reminder
		}
	}

	result, err := client.SendGroupMessage(ctx, args[0], args[1], msgDataMap, sendGroupMessageUserToken, sendGroupMessageSenderID, sendGroupMessageOutlines, sendGroupMessageUUID, sendGroupMessageEntryID, sendGroupMessageRefMsgID)
	checkError(err)
	outputResultFields(result, []string{"message_id"})
}

func runSendAccountMessage(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	msgData, err := parseJSONRaw(args[1])
	checkError(err)

	msgDataMap, ok := msgData.(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: msg_data must be a JSON object\n")
		return
	}

	result, err := client.SendAccountMessage(ctx, args[0], msgDataMap, sendAccountMessageChatIDs, sendAccountMessageDepartmentIDs, sendAccountMessageAccountID, sendAccountMessageEntryID, sendAccountMessageAttach, sendAccountMessageUserToken)
	checkError(err)
	outputResultFields(result, []string{"message_id"})
}

func runSendUserMessage(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	msgData, err := parseJSONRaw(args[2])
	checkError(err)

	msgDataMap, ok := msgData.(map[string]interface{})
	if !ok {
		fmt.Fprintf(os.Stderr, "Error: msg_data must be a JSON object\n")
		return
	}

	var commonMap map[string]interface{}
	if sendUserMessageCommon != "" {
		raw, cErr := parseJSONRaw(sendUserMessageCommon)
		checkError(cErr)
		if m, ok := raw.(map[string]interface{}); ok {
			commonMap = m
		}
	}

	result, err := client.SendUserMessage(ctx, args[0], args[1], msgDataMap, commonMap, sendUserMessageUserToken, sendUserMessageUUID)
	checkError(err)

	outputResultFields(result, []string{"message_id"})
}

func runSendReminder(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.SendReminder(ctx, args[0], sendReminderTypes, sendReminderUserIDs)
	checkError(err)
	outputResultFields(result, []string{"message_id", "operation"})
}

func runQueryGroups(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.QueryGroups(ctx, queryGroupsPageOffset, queryGroupsPageSize)
	checkError(err)

	if result.Success {
		outputResultFields(result, []string{"total_group_ids", "operation"})
		if len(result.GroupIDs) > 0 && !jsonOutput {
			fmt.Println("\nGroups:")
			for _, gid := range result.GroupIDs {
				fmt.Printf("  %s\n", gid)
			}
		}
	} else {
		outputResult(result)
	}
}