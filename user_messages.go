package lansenger

import (
	"context"
	"fmt"
)

func (c *LansengerClient) SendUserMessage(ctx context.Context, receiverID, msgType string, msgData map[string]interface{}, userToken, uuid string) (*UserMessageResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for send_user_message")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "chat_create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"receiverId": receiverID,
		"msgType":    msgType,
		"msgData":    msgData,
	}
	if uuid != "" {
		body["uuid"] = uuid
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &UserMessageResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &UserMessageResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}
