package lansenger

import (
	"context"
	"fmt"
)

func (c *LansengerClient) FetchChatList(ctx context.Context, userToken string, chatType int, keyword string, startTime, endTime int64) (*ChatListResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_chat_list")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	apiURL := BuildAPIURL(c.config, "chats", "fetch", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{}
	if chatType != 0 {
		body["chatType"] = chatType
	}
	if keyword != "" {
		body["keyword"] = keyword
	}
	if startTime != 0 {
		body["startTime"] = startTime
	}
	if endTime != 0 {
		body["endTime"] = endTime
	}

	result, err := c.doPost(ctx, apiURL, body)
	if err != nil {
		return &ChatListResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &ChatListResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &ChatListResult{
		Success:     true,
		RawResponse: result,
	}

	if staffInfos, ok := data["staffIdInfos"].([]interface{}); ok {
		for _, si := range staffInfos {
			if m, ok := si.(map[string]interface{}); ok {
				var sectors []string
				if raw, ok := m["sectorNames"].([]interface{}); ok {
					for _, v := range raw {
						if s, ok := v.(string); ok {
							sectors = append(sectors, s)
						}
					}
				}
				res.StaffInfos = append(res.StaffInfos, ChatStaffInfo{
					StaffID:     strFromMap(m, "staffId"),
					StaffName:   strFromMap(m, "staffName"),
					SectorNames: sectors,
				})
			}
		}
	}

	if groupInfos, ok := data["groupIdInfos"].([]interface{}); ok {
		for _, gi := range groupInfos {
			if m, ok := gi.(map[string]interface{}); ok {
				res.GroupInfos = append(res.GroupInfos, ChatGroupInfo{
					GroupID:   strFromMap(m, "groupId"),
					GroupName: strFromMap(m, "groupName"),
				})
			}
		}
	}

	return res, nil
}

func (c *LansengerClient) FetchChatMessages(ctx context.Context, userToken string, pageSize int, baseVersion string, staffID, groupID string, startTime, endTime int64, senderID string) (*ChatMessagesResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_chat_messages")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	apiURL := BuildAPIURL(c.config, "chats", "messages_fetch", token,
		WithUserToken(userToken),
		WithPageSize(pageSize),
		WithQueryParam("base_version", baseVersion),
	)

	body := map[string]interface{}{}
	if staffID != "" {
		body["staffId"] = staffID
	}
	if groupID != "" {
		body["groupId"] = groupID
	}
	if startTime != 0 {
		body["startTime"] = startTime
	}
	if endTime != 0 {
		body["endTime"] = endTime
	}
	if senderID != "" {
		body["senderId"] = senderID
	}

	result, err := c.doPost(ctx, apiURL, body)
	if err != nil {
		return &ChatMessagesResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &ChatMessagesResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &ChatMessagesResult{
		Success:     true,
		HasMore:     boolFromMap(data, "hasMore"),
		Total:       intFromMap(data, "total"),
		LastVersion: strFromMap(data, "lastVersion"),
		Name:        strFromMap(data, "name"),
		ChatType:    strFromMap(data, "chatType"),
		RawResponse: result,
	}

	if msgs, ok := data["messageList"].([]interface{}); ok {
		for _, msg := range msgs {
			if m, ok := msg.(map[string]interface{}); ok {
				info, _ := m["messageInfo"].(map[string]interface{})
				if info == nil {
					continue
				}
				content, _ := info["content"].(map[string]interface{})
				if content == nil {
					content = map[string]interface{}{}
				}
				res.Messages = append(res.Messages, ChatMessageInfo{
					SendTime:    strFromMap(info, "sendTime"),
					Sender:      strFromMap(info, "sender"),
					MessageType: strFromMap(info, "type"),
					Content:     content,
				})
			}
		}
	}

	return res, nil
}
