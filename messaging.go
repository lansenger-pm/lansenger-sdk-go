package lansenger

import (
	"context"
	"io"
	"net/http"
	"os"
	"strings"
)

func (c *LansengerClient) SendText(ctx context.Context, chatID, content string, filePath string, mediaType int, reminderAll bool, reminderUserIDs []string, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	textData := map[string]interface{}{
		"content": content,
	}
	if filePath != "" {
		uploadResult, err := c.UploadMedia(ctx, filePath, mediaType)
		if err != nil {
			return &SendMessageResult{Success: false, Error: "upload failed: " + err.Error(), Platform: "lansenger"}, nil
		}
		if !uploadResult.Success {
			return &SendMessageResult{Success: false, Error: "upload failed: " + uploadResult.Error, Platform: "lansenger"}, nil
		}
		textData["mediaIds"] = []string{uploadResult.MediaID}
	}
	if reminderAll || len(reminderUserIDs) > 0 {
		reminder := map[string]interface{}{
			"all":    reminderAll,
			"userIds": reminderUserIDs,
		}
		textData["reminder"] = reminder
	}

	msgData := map[string]interface{}{
		"text": textData,
	}

	msgType := "text"
	if isGroup {
		return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData)
}

func (c *LansengerClient) SendMarkdown(ctx context.Context, chatID, content string, reminderAll bool, reminderUserIDs []string, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	formatTextData := map[string]interface{}{
		"formatType": 1,
		"text":       content,
	}
	if reminderAll || len(reminderUserIDs) > 0 {
		reminder := map[string]interface{}{
			"all":    reminderAll,
			"userIds": reminderUserIDs,
		}
		formatTextData["reminder"] = reminder
	}

	msgData := map[string]interface{}{
		"formatText": formatTextData,
	}

	msgType := "formatText"
	if isGroup {
		return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData)
}

func (c *LansengerClient) SendFile(ctx context.Context, chatID, filePath string, caption string, mediaType int, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	uploadResult, err := c.UploadMedia(ctx, filePath, mediaType)
	if err != nil {
		return &SendMessageResult{Success: false, Error: "upload failed: " + err.Error(), Platform: "lansenger"}, nil
	}
	if !uploadResult.Success {
		return &SendMessageResult{Success: false, Error: "upload failed: " + uploadResult.Error, Platform: "lansenger"}, nil
	}

	textData := map[string]interface{}{
		"mediaIds": []string{uploadResult.MediaID},
	}
	if caption != "" {
		textData["caption"] = caption
	}

	msgData := map[string]interface{}{
		"text": textData,
	}

	msgType := "text"
	if isGroup {
		return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData)
}

func (c *LansengerClient) SendImageURL(ctx context.Context, chatID, imageURL, caption string, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	if chatID == "" {
		return &SendMessageResult{Success: false, Error: "chat_id is required", Platform: "lansenger"}, nil
	}
	if imageURL == "" {
		return &SendMessageResult{Success: false, Error: "image_url is required", Platform: "lansenger"}, nil
	}

	resp, err := c.httpClient.Get(imageURL)
	if err != nil {
		return &SendMessageResult{Success: false, Error: "failed to download image: " + err.Error(), Platform: "lansenger"}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &SendMessageResult{Success: false, Error: "failed to download image: HTTP " + resp.Status, Platform: "lansenger"}, nil
	}

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &SendMessageResult{Success: false, Error: "failed to read image data: " + err.Error(), Platform: "lansenger"}, nil
	}

	ct := resp.Header.Get("Content-Type")
	suffix := ".jpg"
	if strings.Contains(ct, "png") {
		suffix = ".png"
	} else if strings.Contains(ct, "gif") {
		suffix = ".gif"
	} else if strings.Contains(ct, "webp") {
		suffix = ".webp"
	}

	tmpFile, err := os.CreateTemp("", "lansenger_url_image_*"+suffix)
	if err != nil {
		return &SendMessageResult{Success: false, Error: "failed to create temp file: " + err.Error(), Platform: "lansenger"}, nil
	}
	tmpPath := tmpFile.Name()

	if _, err := tmpFile.Write(imageBytes); err != nil {
		tmpFile.Close()
		os.Remove(tmpPath)
		return &SendMessageResult{Success: false, Error: "failed to write temp file: " + err.Error(), Platform: "lansenger"}, nil
	}
	tmpFile.Close()

	result, err := c.SendFile(ctx, chatID, tmpPath, caption, MediaTypeImage, isGroup, userToken, senderID)
	os.Remove(tmpPath)
	if err != nil {
		return &SendMessageResult{Success: false, Error: err.Error(), Platform: "lansenger"}, nil
	}
	return result, nil
}

func (c *LansengerClient) SendAppArticles(ctx context.Context, chatID string, articles []map[string]string, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	if chatID == "" {
		return &SendMessageResult{Success: false, Error: "chat_id is required", Platform: "lansenger"}, nil
	}
	if len(articles) == 0 {
		return &SendMessageResult{Success: false, Error: "articles is required", Platform: "lansenger"}, nil
	}

	msgData := map[string]interface{}{
		"appArticles": articles,
	}

	msgType := "appArticles"
	if isGroup {
		return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData)
}

func (c *LansengerClient) SendLinkCardWithParams(ctx context.Context, params *LinkCardParams) (*SendMessageResult, error) {
	linkCardData := map[string]interface{}{
		"title":        params.Title,
		"link":         params.Link,
		"description":  params.Description,
		"iconLink":     params.IconLink,
		"pcLink":       params.PcLink,
		"padLink":      params.PadLink,
		"fromName":     params.FromName,
		"fromIconLink": params.FromIconLink,
	}

	msgData := map[string]interface{}{
		"linkCard": linkCardData,
	}

	msgType := "linkCard"
	if params.IsGroup {
		return c.SendGroupMessage(ctx, params.ChatID, msgType, msgData, params.UserToken, params.SenderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, params.ChatID, msgType, msgData)
}

func (c *LansengerClient) SendAppCardWithParams(ctx context.Context, params *AppCardParams) (*SendMessageResult, error) {
	appCardData := map[string]interface{}{
		"bodyTitle":    params.BodyTitle,
		"headTitle":    params.HeadTitle,
		"bodySubTitle": params.BodySubTitle,
		"bodyContent":  params.BodyContent,
		"signature":    params.Signature,
		"cardLink":     params.CardLink,
		"pcCardLink":   params.PcCardLink,
		"padCardLink":  params.PadCardLink,
		"isDynamic":    params.IsDynamic,
		"staffId":      params.StaffID,
		"headIconUrl":  params.HeadIconURL,
	}
	if len(params.Fields) > 0 {
		appCardData["fields"] = params.Fields
	}
	if len(params.Links) > 0 {
		appCardData["links"] = params.Links
	}
	if params.HeadStatusInfo != nil {
		appCardData["headStatusInfo"] = params.HeadStatusInfo
	}

	msgData := map[string]interface{}{
		"appCard": appCardData,
	}

	msgType := "appCard"
	if params.IsGroup {
		return c.SendGroupMessage(ctx, params.ChatID, msgType, msgData, params.UserToken, params.SenderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, params.ChatID, msgType, msgData)
}

func (c *LansengerClient) SendOaCardWithParams(ctx context.Context, params *OaCardParams) (*SendMessageResult, error) {
	oaCardData := map[string]interface{}{
		"head":     params.Head,
		"title":    params.Title,
		"subTitle": params.SubTitle,
		"staffId":  params.StaffID,
		"link":     params.Link,
		"pcLink":   params.PcLink,
		"padLink":  params.PadLink,
	}
	if len(params.Fields) > 0 {
		oaCardData["fields"] = params.Fields
	}
	if params.CardAction != nil {
		oaCardData["cardAction"] = params.CardAction
	}

	msgData := map[string]interface{}{
		"oacard": oaCardData,
	}

	msgType := "oacard"
	if params.IsGroup {
		return c.SendGroupMessage(ctx, params.ChatID, msgType, msgData, params.UserToken, params.SenderID, "", "", "")
	}

	return c.sendBotPrivate(ctx, params.ChatID, msgType, msgData)
}

func (c *LansengerClient) UpdateDynamicCard(ctx context.Context, params *DynamicCardUpdateParams) (*SendMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "dynamic_update", token)

	msgData := map[string]interface{}{
		"isLastUpdate": params.IsLastUpdate,
	}
	if params.HeadStatusInfo != nil {
		msgData["headStatusInfo"] = params.HeadStatusInfo
	}
	if len(params.Links) > 0 {
		msgData["links"] = params.Links
	}

	body := map[string]interface{}{
		"msgId":   params.MsgID,
		"msgType": "appCard",
		"msgData": map[string]interface{}{
			"appCardUpdateMsg": msgData,
		},
	}
	if params.UserId != "" {
		body["userId"] = params.UserId
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &SendMessageResult{Success: false, Error: err.Error(), Platform: "lansenger"}, nil
	}

	data := extractData(result)

	res := &SendMessageResult{
		Success:     true,
		Platform:    "lansenger",
		Operation:   "update_dynamic_card",
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}

func (c *LansengerClient) RevokeMessage(ctx context.Context, messageIDs []string, chatType string, senderID string, sysMsg *SysMsgParams) (*SendMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "revoke", token)

	body := map[string]interface{}{
		"chatType":   chatType,
		"messageIds": messageIDs,
	}
	if senderID != "" {
		body["senderId"] = senderID
	}
	if sysMsg != nil {
		sysMsgData := map[string]interface{}{}
		if sysMsg.Content != "" {
			sysMsgData["content"] = sysMsg.Content
		}
		if sysMsg.MediaID != "" {
			sysMsgData["mediaId"] = sysMsg.MediaID
		}
		body["sysMsg"] = sysMsgData
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &SendMessageResult{Success: false, Error: err.Error(), Platform: "lansenger"}, nil
	}

	return &SendMessageResult{
		Success:     true,
		Platform:    "lansenger",
		Operation:   "revoke_message",
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) QueryGroups(ctx context.Context, pageOffset, pageSize int) (*QueryGroupsResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "list_fetch", token,
		WithPageOffset(pageOffset),
		WithPageSize(pageSize),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &QueryGroupsResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &QueryGroupsResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &QueryGroupsResult{
		Success:       true,
		TotalGroupIDs: intFromMap(data, "totalGroupIds"),
		Platform:      "lansenger",
		Operation:     "query_groups",
		RawResponse:   result,
	}, nil
}

func (c *LansengerClient) SendReminder(ctx context.Context, msgID string, reminderTypes []int, userIDList []string) (*SendMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "reminder_create", token)

	body := map[string]interface{}{
		"msgId":          msgID,
		"reminderTypes":  reminderTypes,
		"userIdList":     userIDList,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &SendMessageResult{Success: false, Error: err.Error(), Platform: "lansenger"}, nil
	}

	return &SendMessageResult{
		Success:     true,
		Platform:    "lansenger",
		Operation:   "send_reminder",
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) sendBotPrivate(ctx context.Context, chatID, msgType string, msgData map[string]interface{}) (*SendMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "bot", "messages_create", token)

	body := map[string]interface{}{
		"userIdList": []string{chatID},
		"msgType":    msgType,
		"msgData":    msgData,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &SendMessageResult{Success: false, Error: err.Error(), Platform: "lansenger"}, nil
	}

	data := extractData(result)

	res := &SendMessageResult{
		Success:     true,
		Platform:    "lansenger",
		MsgType:     msgType,
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}