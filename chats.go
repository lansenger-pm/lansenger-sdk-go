package lansenger

import (
	"context"
	"fmt"
	"net/url"
)

func addQueryParam(baseURL, key, value string) string {
	sep := "&"
	if !containsChar(baseURL, '?') {
		sep = "?"
	}
	return baseURL + sep + url.QueryEscape(key) + "=" + url.QueryEscape(value)
}

func containsChar(s string, c byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return true
		}
	}
	return false
}

func (c *LansengerClient) FetchChatList(ctx context.Context, userToken string, chatType string, keyword, startTime, endTime string) (*ChatListResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_chat_list")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "chats", "fetch", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{}
	if chatType != "" {
		body["chatType"] = chatType
	}
	if keyword != "" {
		body["keyword"] = keyword
	}
	if startTime != "" {
		body["startTime"] = startTime
	}
	if endTime != "" {
		body["endTime"] = endTime
	}

	result, err := c.doPost(ctx, url, body)
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

func (c *LansengerClient) FetchChatMessages(ctx context.Context, userToken string, pageSize int, baseVersion string, staffID, groupID, startTime, endTime, senderID string) (*ChatMessagesResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_chat_messages")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "messages", "fetch", token,
		WithUserToken(userToken),
		WithPageSize(pageSize),
		WithQueryParam("base_version", baseVersion),
	)

	if staffID != "" {
		url = addQueryParam(url, "staffId", staffID)
	}
	if groupID != "" {
		url = addQueryParam(url, "groupId", groupID)
	}
	if startTime != "" {
		url = addQueryParam(url, "startTime", startTime)
	}
	if endTime != "" {
		url = addQueryParam(url, "endTime", endTime)
	}
	if senderID != "" {
		url = addQueryParam(url, "senderId", senderID)
	}

	result, err := c.doGet(ctx, url)
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
				res.Messages = append(res.Messages, ChatMessageInfo{
					SendTime:    strFromMap(m, "sendTime"),
					Sender:      strFromMap(m, "sender"),
					MessageType: strFromMap(m, "messageType"),
					Content:     m,
				})
			}
		}
	}

	return res, nil
}
