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

	var msgID string
	if data != nil {
		msgID = strFromMap(data, "msgId")
	}
	// Server has been observed returning success with an empty payload
	// (LXBUGS-128497); without a msgId the fetch step is unusable, so treat
	// it as a failure instead of reporting success.
	if msgID == "" {
		return &StreamMessageResult{
			Success:     false,
			Error:       "server returned success but no msgId; stream message is unusable",
			RawResponse: result,
		}, nil
	}
	return &StreamMessageResult{
		Success:     true,
		MessageID:   msgID,
		RawResponse: result,
	}, nil
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
