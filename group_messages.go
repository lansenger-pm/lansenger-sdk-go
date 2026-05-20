package lansenger

import "context"

func (c *LansengerClient) SendGroupMessage(ctx context.Context, groupID, msgType string, msgData map[string]interface{}, userToken, senderID string, reminderAll bool, reminderUserIDs []string, outlines, uuid, entryID string) (*SendMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "group_create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"groupId": groupID,
		"msgType": msgType,
		"msgData": msgData,
	}
	if senderID != "" {
		body["senderId"] = senderID
	}
	if reminderAll {
		body["reminderAll"] = true
	}
	if len(reminderUserIDs) > 0 {
		body["reminderUserIds"] = reminderUserIDs
	}
	if outlines != "" {
		body["outlines"] = outlines
	}
	if uuid != "" {
		body["uuid"] = uuid
	}
	if entryID != "" {
		body["entryId"] = entryID
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
		Operation:   "send_group_message",
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}
