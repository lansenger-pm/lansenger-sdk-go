package lansenger

import "context"

func (c *LansengerClient) CreateStreamMessage(ctx context.Context, receiverID, receiverType, streamID string) (*StreamMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "sse", "msg_create", token)

	body := map[string]interface{}{
		"receiverId":   receiverID,
		"receiverType": receiverType,
		"streamId":     streamID,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &StreamMessageResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &StreamMessageResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}

func (c *LansengerClient) FetchStreamMessage(ctx context.Context, msgID string) (*StreamMessageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "sse", "msg_fetch", token)

	body := map[string]interface{}{
		"msgId": msgID,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &StreamMessageResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &StreamMessageResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.MessageID = strFromMap(data, "msgId")
	}
	return res, nil
}
