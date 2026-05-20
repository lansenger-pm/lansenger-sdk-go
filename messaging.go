package lansenger

import (
	"context"
)

func (c *LansengerClient) SendText(ctx context.Context, chatID, content string, filePath string, mediaType int, reminderAll bool, reminderUserIDs []string, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	msgData := map[string]interface{}{
		"content": content,
	}
	if filePath != "" {
		msgData["filePath"] = filePath
	}
	if mediaType != 0 {
		msgData["mediaType"] = mediaType
	}

	msgType := "text"
	if isGroup {
		return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, reminderAll, reminderUserIDs, "", "", "")
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData, reminderAll, reminderUserIDs)
}

func (c *LansengerClient) SendMarkdown(ctx context.Context, chatID, content string, reminderAll bool, reminderUserIDs []string, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	msgData := map[string]interface{}{
		"content": content,
	}

	msgType := "markdown"
	if isGroup {
		result, err := c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, reminderAll, reminderUserIDs, "", "", "")
		if err != nil {
			if reminderAll || len(reminderUserIDs) > 0 {
				return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, false, nil, "", "", "")
			}
		}
		return result, err
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData, reminderAll, reminderUserIDs)
}

func (c *LansengerClient) SendFile(ctx context.Context, chatID, filePath string, caption string, mediaType int, isGroup bool, userToken, senderID string) (*SendMessageResult, error) {
	uploadResult, err := c.UploadMedia(ctx, filePath, mediaType)
	if err != nil {
		return &SendMessageResult{Success: false, Error: "upload failed: " + err.Error(), Platform: "lansenger"}, nil
	}
	if !uploadResult.Success {
		return &SendMessageResult{Success: false, Error: "upload failed: " + uploadResult.Error, Platform: "lansenger"}, nil
	}

	msgData := map[string]interface{}{
		"filePath": uploadResult.MediaID,
	}
	if caption != "" {
		msgData["caption"] = caption
	}

	msgType := "file"
	if isGroup {
		return c.SendGroupMessage(ctx, chatID, msgType, msgData, userToken, senderID, false, nil, "", "", "")
	}

	return c.sendBotPrivate(ctx, chatID, msgType, msgData, false, nil)
}

func (c *LansengerClient) SendLinkCardWithParams(ctx context.Context, params *LinkCardParams) (*SendMessageResult, error) {
	msgData := map[string]interface{}{
		"title":        params.Title,
		"link":         params.Link,
		"description":  params.Description,
		"iconLink":     params.IconLink,
		"pcLink":       params.PcLink,
		"padLink":      params.PadLink,
		"fromName":     params.FromName,
		"fromIconLink": params.FromIconLink,
	}

	msgType := "linkCard"
	if params.IsGroup {
		return c.SendGroupMessage(ctx, params.ChatID, msgType, msgData, params.UserToken, params.SenderID, false, nil, "", "", "")
	}

	return c.sendBotPrivate(ctx, params.ChatID, msgType, msgData, false, nil)
}

func (c *LansengerClient) SendAppCardWithParams(ctx context.Context, params *AppCardParams) (*SendMessageResult, error) {
	msgData := map[string]interface{}{
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
		msgData["fields"] = params.Fields
	}
	if len(params.Links) > 0 {
		msgData["links"] = params.Links
	}
	if params.HeadStatusInfo != nil {
		msgData["headStatusInfo"] = params.HeadStatusInfo
	}

	msgType := "appCard"
	if params.IsGroup {
		return c.SendGroupMessage(ctx, params.ChatID, msgType, msgData, params.UserToken, params.SenderID, false, nil, "", "", "")
	}

	return c.sendBotPrivate(ctx, params.ChatID, msgType, msgData, false, nil)
}

func (c *LansengerClient) SendOaCardWithParams(ctx context.Context, params *OaCardParams) (*SendMessageResult, error) {
	msgData := map[string]interface{}{
		"head":     params.Head,
		"title":    params.Title,
		"subTitle": params.SubTitle,
		"staffId":  params.StaffID,
		"link":     params.Link,
		"pcLink":   params.PcLink,
		"padLink":  params.PadLink,
	}
	if len(params.Fields) > 0 {
		msgData["fields"] = params.Fields
	}
	if params.CardAction != nil {
		msgData["cardAction"] = params.CardAction
	}

	msgType := "oaCard"
	if params.IsGroup {
		return c.SendGroupMessage(ctx, params.ChatID, msgType, msgData, params.UserToken, params.SenderID, false, nil, "", "", "")
	}

	return c.sendBotPrivate(ctx, params.ChatID, msgType, msgData, false, nil)
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

func (c *LansengerClient) RevokeMessage(ctx context.Context, messageIDs []string, chatType string, senderID string) (*SendMessageResult, error) {
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

func (c *LansengerClient) sendBotPrivate(ctx context.Context, chatID, msgType string, msgData map[string]interface{}, reminderAll bool, reminderUserIDs []string) (*SendMessageResult, error) {
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
	if reminderAll {
		body["reminderAll"] = reminderAll
	}
	if len(reminderUserIDs) > 0 {
		body["reminderUserIds"] = reminderUserIDs
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		if reminderAll || len(reminderUserIDs) > 0 {
			body["reminderAll"] = false
			body["reminderUserIds"] = nil
			result2, err2 := c.doPost(ctx, url, body)
			if err2 != nil {
				return &SendMessageResult{Success: false, Error: err2.Error(), Platform: "lansenger"}, nil
			}
			data2 := extractData(result2)
			res2 := &SendMessageResult{Success: true, Platform: "lansenger", MsgType: msgType, RawResponse: result2}
			if data2 != nil {
				res2.MessageID = strFromMap(data2, "msgId")
			}
			return res2, nil
		}
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
