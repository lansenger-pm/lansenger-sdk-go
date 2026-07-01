package lansenger

import (
	"context"
)

func (c *LansengerClient) CreateBotCommands(ctx context.Context, scopeType int, commands []map[string]interface{}, chatID, chatType, staffID string) (*BotCommandResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "bot_commands", "create", token)

	body := map[string]interface{}{
		"scopeType": scopeType,
		"commands":  commands,
	}
	if chatID != "" {
		body["chatId"] = chatID
	}
	if chatType != "" {
		body["chatType"] = chatType
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BotCommandResult{Success: false, Error: err.Error()}, nil
	}

	return &BotCommandResult{
		Success:     true,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchBotCommands(ctx context.Context, scopeType int, chatID, chatType, staffID string) (*BotCommandQueryResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "bot_commands", "fetch", token)

	body := map[string]interface{}{
		"scopeType": scopeType,
	}
	if chatID != "" {
		body["chatId"] = chatID
	}
	if chatType != "" {
		body["chatType"] = chatType
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BotCommandQueryResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	res := &BotCommandQueryResult{
		Success:     true,
		RawResponse: result,
	}
	if v, ok := data["scopeType"].(float64); ok {
		res.ScopeType = int(v)
	}
	if v, ok := data["chatId"].(string); ok {
		res.ChatID = v
	}
	if v, ok := data["chatType"].(string); ok {
		res.ChatType = v
	}
	if v, ok := data["staffId"].(string); ok {
		res.StaffID = v
	}
	if cmds, ok := data["commands"].([]interface{}); ok {
		for _, cmd := range cmds {
			if m, ok := cmd.(map[string]interface{}); ok {
				res.Commands = append(res.Commands, m)
			}
		}
	}

	return res, nil
}

func (c *LansengerClient) DeleteBotCommands(ctx context.Context, scopeType int, chatID, chatType, staffID string) (*BotCommandResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "bot_commands", "delete", token)

	body := map[string]interface{}{
		"scopeType": scopeType,
	}
	if chatID != "" {
		body["chatId"] = chatID
	}
	if chatType != "" {
		body["chatType"] = chatType
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &BotCommandResult{Success: false, Error: err.Error()}, nil
	}

	return &BotCommandResult{
		Success:     true,
		RawResponse: result,
	}, nil
}
