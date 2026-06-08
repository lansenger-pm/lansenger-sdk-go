package lansenger

import (
	"context"
)

func (c *LansengerClient) CreateTodoTask(ctx context.Context, title string, todoType int, link, pcLink string, executorIDs []string, orgID, sourceID, desc, senderID, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"title":       title,
		"type":        todoType,
		"link":        link,
		"pcLink":      pcLink,
		"executorIds": executorIDs,
		"orgId":       orgID,
	}
	if sourceID != "" {
		body["sourceId"] = sourceID
	}
	if desc != "" {
		body["desc"] = desc
	}
	if senderID != "" {
		body["senderId"] = senderID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) UpdateTodoTask(ctx context.Context, todotaskID, title, link, pcLink, orgID, desc, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "info_update", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"todotaskId": todotaskID,
		"orgId":      orgID,
	}
	if title != "" {
		body["title"] = title
	}
	if link != "" {
		body["link"] = link
	}
	if pcLink != "" {
		body["pcLink"] = pcLink
	}
	if desc != "" {
		body["desc"] = desc
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) UpdateTodoTaskStatus(ctx context.Context, todotaskID, status, orgID, staffID, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "status_update", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"todotaskId": todotaskID,
		"status":     status,
		"orgId":      orgID,
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) DeleteTodoTask(ctx context.Context, todotaskID, orgID, staffID, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "sender_delete", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"todotaskId": todotaskID,
		"orgId":      orgID,
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) FetchTodoTaskList(ctx context.Context, orgID string, appIDs []string, staffID string, statusList []string, userToken string) (*TodoTaskListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "list_fetch", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"orgId": orgID,
	}
	if len(appIDs) > 0 {
		body["appIds"] = appIDs
	}
	if staffID != "" {
		body["staffId"] = staffID
	}
	if len(statusList) > 0 {
		body["statusList"] = statusList
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskListResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskListResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.Total = intFromMap(data, "total")
		if list, ok := data["todotaskList"].([]interface{}); ok {
			res.TodotaskList = make([]map[string]interface{}, 0, len(list))
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					res.TodotaskList = append(res.TodotaskList, m)
				}
			}
		}
	}
	return res, nil
}

func (c *LansengerClient) FetchTodoTaskBySourceID(ctx context.Context, sourceID, orgID, staffID, userToken string) (*TodoTaskInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "info_fetch_by_source_id", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"sourceId": sourceID,
		"orgId":    orgID,
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskInfoResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
		res.SourceID = strFromMap(data, "sourceId")
		res.Title = strFromMap(data, "title")
		res.Desc = strFromMap(data, "desc")
		res.Status = strFromMap(data, "status")
		res.Type = intFromMap(data, "type")
		res.Link = strFromMap(data, "link")
		res.PcLink = strFromMap(data, "pcLink")
		res.SenderID = strFromMap(data, "senderId")
		res.CreateTime = strFromMap(data, "createTime")
		res.AppID = strFromMap(data, "appId")
		res.ExecutorIDs = stringArrayFromMap(data, "executorIds")
	}
	return res, nil
}

func (c *LansengerClient) FetchTodoTaskByID(ctx context.Context, todotaskID, orgID, staffID, userToken string) (*TodoTaskInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "info_fetch", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"todotaskId": todotaskID,
		"orgId":      orgID,
	}
	if staffID != "" {
		body["staffId"] = staffID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskInfoResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
		res.SourceID = strFromMap(data, "sourceId")
		res.Title = strFromMap(data, "title")
		res.Desc = strFromMap(data, "desc")
		res.Status = strFromMap(data, "status")
		res.Type = intFromMap(data, "type")
		res.Link = strFromMap(data, "link")
		res.PcLink = strFromMap(data, "pcLink")
		res.SenderID = strFromMap(data, "senderId")
		res.CreateTime = strFromMap(data, "createTime")
		res.AppID = strFromMap(data, "appId")
		res.ExecutorIDs = stringArrayFromMap(data, "executorIds")
	}
	return res, nil
}

func (c *LansengerClient) FetchTodoTaskStatusCounts(ctx context.Context, staffID, orgID, appID, status, userToken string) (*TodoTaskStatusCountResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "status_count_list_fetch", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"staffId": staffID,
		"orgId":   orgID,
	}
	if appID != "" {
		body["appId"] = appID
	}
	if status != "" {
		body["status"] = status
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskStatusCountResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskStatusCountResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		if counts, ok := data["statusCounts"].([]interface{}); ok {
			res.StatusCounts = make([]map[string]interface{}, 0, len(counts))
			for _, item := range counts {
				if m, ok := item.(map[string]interface{}); ok {
					res.StatusCounts = append(res.StatusCounts, m)
				}
			}
		}
	}
	return res, nil
}

func (c *LansengerClient) UpdateExecutorStatus(ctx context.Context, executorStatusList []map[string]interface{}, orgID, todotaskID, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "executor_status_update", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"executorStatusList": executorStatusList,
		"orgId":              orgID,
	}
	if todotaskID != "" {
		body["todotaskId"] = todotaskID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) AddExecutors(ctx context.Context, executorIDs []string, orgID, todotaskID, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "executor_create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"executorIds": executorIDs,
		"orgId":       orgID,
	}
	if todotaskID != "" {
		body["todotaskId"] = todotaskID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) DeleteExecutors(ctx context.Context, executorIDs []string, orgID, todotaskID, userToken string) (*TodoTaskCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "executor_delete", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"executorIds": executorIDs,
		"orgId":       orgID,
	}
	if todotaskID != "" {
		body["todotaskId"] = todotaskID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskCreateResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TodotaskID = strFromMap(data, "todotaskId")
	}
	return res, nil
}

func (c *LansengerClient) FetchExecutorList(ctx context.Context, todotaskID, orgID, staffID string, statusList []string, userToken string) (*TodoTaskExecutorListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "todo", "executor_list_fetch", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"todotaskId": todotaskID,
		"orgId":      orgID,
	}
	if staffID != "" {
		body["staffId"] = staffID
	}
	if len(statusList) > 0 {
		body["statusList"] = statusList
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &TodoTaskExecutorListResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &TodoTaskExecutorListResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.Total = intFromMap(data, "total")
		if list, ok := data["executorList"].([]interface{}); ok {
			res.ExecutorList = make([]map[string]interface{}, 0, len(list))
			for _, item := range list {
				if m, ok := item.(map[string]interface{}); ok {
					res.ExecutorList = append(res.ExecutorList, m)
				}
			}
		}
	}
	return res, nil
}
