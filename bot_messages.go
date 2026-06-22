package lansenger

import "context"

func (c *LansengerClient) SendBotMessage(ctx context.Context, msgType string, msgData map[string]interface{}, chatIDs, departmentIDs []string, userToken, entryID, refMsgID string) (*BotMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "bot", "messages_create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"msgType": msgType,
		"msgData": msgData,
	}
	if len(chatIDs) > 0 {
		body["userIdList"] = chatIDs
	}
	if len(departmentIDs) > 0 {
		body["departmentIdList"] = departmentIDs
	}
	if entryID != "" {
		body["entryId"] = entryID
	}
	if refMsgID != "" {
		body["refMsgId"] = refMsgID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BotMessageResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &BotMessageResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
		res.InvalidStaff = stringArrayFromMap(data, "invalidStaff")
		res.InvalidDepartment = stringArrayFromMap(data, "invalidDepartment")
	}
	return res, nil
}
