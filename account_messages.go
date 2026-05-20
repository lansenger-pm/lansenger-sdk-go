package lansenger

import "context"

func (c *LansengerClient) SendAccountMessage(ctx context.Context, msgType string, msgData map[string]interface{}, chatIDs, departmentIDs []string, accountID, entryID, attach, userToken string) (*AccountMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "create", token,
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
	if accountID != "" {
		body["accountId"] = accountID
	}
	if entryID != "" {
		body["entryId"] = entryID
	}
	if attach != "" {
		body["attach"] = attach
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &AccountMessageResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &AccountMessageResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}
