package lansenger

import (
	"context"
	"fmt"
)

// BuildPersonalTodoResourceEntry 拼装「挂到待办」用的 resources 条目。
//
// 陷阱：上传接口（/resource/update）返回的是 mimeType/size，而挂附件写体的条目
// 必须叫 fileType/fileSize，直接把上传响应塞进 resources 会被后端 errCode 500 打回。
// opt: 1=添加，0=移除。
//
// 这是条目结构的唯一定义处，`PersonalTodoResourceEntryFromUpload` 经由它产出，
// 与 Python/TS SDK 的 build_personal_todo_resource_entry / buildPersonalTodoResourceEntry 同义。
func BuildPersonalTodoResourceEntry(resourceID, fileName, fileType string, fileSize int64, opt int) map[string]interface{} {
	return map[string]interface{}{
		"fileName":   fileName,
		"resourceId": resourceID,
		"fileType":   fileType,
		"fileSize":   fileSize,
		"opt":        opt,
	}
}

// PersonalTodoResourceEntryFromUpload 从上传结果构造 resources 条目，省掉手工做
// mimeType→fileType / size→fileSize 的映射。
//
// upload 两种形态都收（与 Python/TS SDK 的同名助手等价）：
//   - *PersonalTodoResourceResult：UploadPersonalTodoResource 的返回，字段已解析；
//   - map[string]interface{}：/resource/update 的原始响应，或其内层 data。
//
// 其他类型返回 error，不会静默退化成空条目——缺 resourceId 的条目会被后端打回，
// 静默比报错难查得多。
func PersonalTodoResourceEntryFromUpload(upload interface{}, opt int) (map[string]interface{}, error) {
	switch v := upload.(type) {
	case *PersonalTodoResourceResult:
		if v == nil {
			return nil, fmt.Errorf("upload result is nil")
		}
		return BuildPersonalTodoResourceEntry(v.ResourceID, v.FileName, v.MimeType, v.Size, opt), nil
	case PersonalTodoResourceResult:
		return BuildPersonalTodoResourceEntry(v.ResourceID, v.FileName, v.MimeType, v.Size, opt), nil
	case map[string]interface{}:
		inner := v
		if data, ok := v["data"].(map[string]interface{}); ok {
			inner = data
		}
		return BuildPersonalTodoResourceEntry(
			firstNonEmptyStr(strFromMap(inner, "resourceId"), strFromMap(v, "resourceId")),
			firstNonEmptyStr(strFromMap(inner, "fileName"), strFromMap(v, "fileName")),
			firstNonEmptyStr(strFromMap(inner, "mimeType"), strFromMap(v, "mimeType")),
			firstNonZeroInt64(int64FromMap(inner, "size"), int64FromMap(v, "size")),
			opt,
		), nil
	default:
		return nil, fmt.Errorf(
			"unsupported upload type %T: expected *PersonalTodoResourceResult or a response map", upload)
	}
}

func firstNonEmptyStr(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func firstNonZeroInt64(values ...int64) int64 {
	for _, v := range values {
		if v != 0 {
			return v
		}
	}
	return 0
}

// PersonalTodoSaveParams carries fields for SavePersonalTodo (个人待办 /v3/taskopt/savePersonalTask).
type PersonalTodoSaveParams struct {
	Subject           string
	StartTime         int64
	DueTime           int64
	Priority          int
	CreateUserID      string
	OrgID             string
	AppID             string
	Description       string
	ParentCode        string
	GroupID           string
	GroupCategoryCode string
	FinishTime        int64
	StatusTagNo       string
	StatusTagYes      string
	AppInfoID         *int64
	AppCategoryID     *int64
	Platform          *int
	SubscribeStatus   *int
	UserCode          string
	Executors         []map[string]interface{}
	Copys             []map[string]interface{}
	// Resources 挂附件条目。注意：上传接口返回的是 mimeType/size，挂附件必须映射成
	// fileType/fileSize（每条必填 fileName/resourceId/fileType/fileSize），直接把上传响应塞进来
	// 会被后端 errCode 500 打回；列表侧 resourceList 是只读字段，不能当写字段。
	Resources []map[string]interface{}
	Reminds   []map[string]interface{}
	UserToken string
}

// PersonalTodoUpdateParams carries fields for UpdatePersonalTodo (个人待办 /v3/taskopt/updatePersonalTask).
type PersonalTodoUpdateParams struct {
	TodoCode          string
	OrgID             string
	UpdateFields      []string
	Subject           string
	Description       string
	StartTime         *int64
	DueTime           *int64
	FinishTime        *int64
	Priority          *int
	StatusTagNo       string
	StatusTagYes      string
	SubscribeStatus   *int
	CreateUserID      string
	GroupID           string
	GroupCategoryCode string
	AppID             string
	Executors         []map[string]interface{}
	Copys             []map[string]interface{}
	// Resources 挂附件条目。注意：上传接口返回的是 mimeType/size，挂附件必须映射成
	// fileType/fileSize（每条必填 fileName/resourceId/fileType/fileSize），直接把上传响应塞进来
	// 会被后端 errCode 500 打回；列表侧 resourceList 是只读字段，不能当写字段。
	Resources []map[string]interface{}
	Reminds   []map[string]interface{}
	UserToken string
}

// PersonalTodoResourceUploadParams carries fields for UploadPersonalTodoResource.
type PersonalTodoResourceUploadParams struct {
	AppID         string
	Size          int64
	FileName      string
	ContentType   string
	FileData      string
	OrgID         string
	ExtensionInfo string
	Thumb         bool
	UserToken     string
}

// SavePersonalTodo creates a user-owned personal todo.
func (c *LansengerClient) SavePersonalTodo(ctx context.Context, p *PersonalTodoSaveParams) (*PersonalTodoSaveResult, error) {
	if p == nil {
		return &PersonalTodoSaveResult{Success: false, Error: "params is required"}, nil
	}
	if p.Subject == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "subject is required"}, nil
	}
	if p.StartTime < 0 {
		return &PersonalTodoSaveResult{Success: false, Error: "start_time is required"}, nil
	}
	if p.DueTime < 0 {
		return &PersonalTodoSaveResult{Success: false, Error: "due_time is required"}, nil
	}
	if p.Priority < PersonalTodoPriorityLow || p.Priority > PersonalTodoPriorityVeryUrgent {
		return &PersonalTodoSaveResult{Success: false, Error: "priority must be 0, 1, 2, or 3"}, nil
	}
	if p.CreateUserID == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "create_user_id is required"}, nil
	}
	if p.OrgID == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "org_id is required"}, nil
	}
	if p.AppID == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "appid is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "personal_todos", "save", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"subject":      p.Subject,
		"startTime":    p.StartTime,
		"dueTime":      p.DueTime,
		"finishTime":   p.FinishTime,
		"priority":     p.Priority,
		"type":         PersonalTodoTypePersonal,
		"createUserId": p.CreateUserID,
		"orgId":        p.OrgID,
		"appid":        p.AppID,
	}
	if p.Description != "" {
		body["description"] = p.Description
	}
	if p.ParentCode != "" {
		body["parentCode"] = p.ParentCode
	}
	if p.GroupID != "" {
		body["groupId"] = p.GroupID
	}
	if p.GroupCategoryCode != "" {
		body["groupCategoryCode"] = p.GroupCategoryCode
	}
	if p.StatusTagNo != "" {
		body["statusTagNo"] = p.StatusTagNo
	}
	if p.StatusTagYes != "" {
		body["statusTagYes"] = p.StatusTagYes
	}
	if p.AppInfoID != nil {
		body["appInfoId"] = *p.AppInfoID
	}
	if p.AppCategoryID != nil {
		body["appCategoryId"] = *p.AppCategoryID
	}
	if p.Platform != nil {
		body["platform"] = *p.Platform
	}
	if p.SubscribeStatus != nil {
		body["subscribeStatus"] = *p.SubscribeStatus
	}
	if p.UserCode != "" {
		body["userCode"] = p.UserCode
	}
	if len(p.Executors) > 0 {
		body["executors"] = p.Executors
	}
	if len(p.Copys) > 0 {
		body["copys"] = p.Copys
	}
	if len(p.Resources) > 0 {
		body["resources"] = p.Resources
	}
	if len(p.Reminds) > 0 {
		body["reminds"] = p.Reminds
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalTodoSaveResult{Success: false, Error: err.Error()}, nil
	}
	res := &PersonalTodoSaveResult{Success: true, RawResponse: result}
	if code, ok := result["data"].(string); ok {
		res.TodoCode = code
	}
	return res, nil
}

// UpdatePersonalTodo updates selected fields of a personal todo.
func (c *LansengerClient) UpdatePersonalTodo(ctx context.Context, p *PersonalTodoUpdateParams) (*PersonalTodoSaveResult, error) {
	if p == nil {
		return &PersonalTodoSaveResult{Success: false, Error: "params is required"}, nil
	}
	if p.TodoCode == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "todo_code is required"}, nil
	}
	if p.OrgID == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "org_id is required"}, nil
	}
	if len(p.UpdateFields) == 0 {
		return &PersonalTodoSaveResult{Success: false, Error: "update_fields is required"}, nil
	}
	if p.CreateUserID == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "create_user_id is required"}, nil
	}
	if p.AppID == "" {
		return &PersonalTodoSaveResult{Success: false, Error: "appid is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "personal_todos", "update", token, WithUserToken(p.UserToken))
	updateContent := map[string]interface{}{"code": p.TodoCode}
	required := func(name string) bool {
		for _, field := range p.UpdateFields {
			if field == name {
				return true
			}
		}
		return false
	}
	if required("subject") && p.Subject != "" {
		updateContent["subject"] = p.Subject
	}
	if required("description") && p.Description != "" {
		updateContent["description"] = p.Description
	}
	if required("startTime") && p.StartTime != nil {
		updateContent["startTime"] = *p.StartTime
	}
	if required("dueTime") && p.DueTime != nil {
		updateContent["dueTime"] = *p.DueTime
	}
	if required("finishTime") && p.FinishTime != nil {
		updateContent["finishTime"] = *p.FinishTime
	}
	if required("priority") && p.Priority != nil {
		updateContent["priority"] = *p.Priority
	}
	if required("statusTagNo") && p.StatusTagNo != "" {
		updateContent["statusTagNo"] = p.StatusTagNo
	}
	if required("statusTagYes") && p.StatusTagYes != "" {
		updateContent["statusTagYes"] = p.StatusTagYes
	}
	if required("subscribeStatus") && p.SubscribeStatus != nil {
		updateContent["subscribeStatus"] = *p.SubscribeStatus
	}
	if required("groupId") && p.GroupID != "" {
		updateContent["groupId"] = p.GroupID
	}
	if required("groupCategoryCode") && p.GroupCategoryCode != "" {
		updateContent["groupCategoryCode"] = p.GroupCategoryCode
	}
	updateContent["createUserId"] = p.CreateUserID
	updateContent["appid"] = p.AppID
	if required("executors") && len(p.Executors) > 0 {
		updateContent["executors"] = p.Executors
	}
	if required("copys") && len(p.Copys) > 0 {
		updateContent["copys"] = p.Copys
	}
	if required("resources") && len(p.Resources) > 0 {
		updateContent["resources"] = p.Resources
	}
	if required("reminds") && len(p.Reminds) > 0 {
		updateContent["reminds"] = p.Reminds
	}
	body := map[string]interface{}{
		"orgId":         p.OrgID,
		"updateFields":  p.UpdateFields,
		"updateContent": updateContent,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalTodoSaveResult{Success: false, Error: err.Error()}, nil
	}
	return &PersonalTodoSaveResult{Success: true, TodoCode: p.TodoCode, RawResponse: result}, nil
}

// FetchPersonalTodoList pages personal todos for one user.
func (c *LansengerClient) FetchPersonalTodoList(ctx context.Context, orgID, staffID string, pageNo, pageSize int, status *int, appID, appCategoryName, userToken string) (*PersonalTodoListResult, error) {
	if orgID == "" {
		return &PersonalTodoListResult{Success: false, Error: "org_id is required"}, nil
	}
	if staffID == "" {
		return &PersonalTodoListResult{Success: false, Error: "staff_id is required"}, nil
	}
	if status != nil && *status != PersonalTodoStatusUnfinished && *status != PersonalTodoStatusFinished {
		return &PersonalTodoListResult{Success: false, Error: "status must be 0 or 1"}, nil
	}
	if pageNo <= 0 {
		pageNo = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "personal_todos", "user_list", token, WithUserToken(userToken))
	body := map[string]interface{}{
		"orgId":    orgID,
		"staffId":  staffID,
		"pageNo":   pageNo,
		"pageSize": pageSize,
	}
	if status != nil {
		body["status"] = *status
	}
	if appID != "" {
		body["appId"] = appID
	}
	if appCategoryName != "" {
		body["appCategoryName"] = appCategoryName
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalTodoListResult{Success: false, Error: err.Error()}, nil
	}
	res := &PersonalTodoListResult{Success: true, RawResponse: result}
	data := extractData(result)
	if data == nil {
		return res, nil
	}
	res.PageNo = intFromMap(data, "pageNo")
	res.PageSize = intFromMap(data, "pageSize")
	res.Pages = intFromMap(data, "pages")
	res.Total = intFromMap(data, "total")
	res.HasMore = boolFromMap(data, "hasNextPage")
	if list, ok := data["result"].([]interface{}); ok {
		res.Items = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Items = append(res.Items, m)
			}
		}
	}
	return res, nil
}

// UploadPersonalTodoResource uploads a base64-encoded resource.
func (c *LansengerClient) UploadPersonalTodoResource(ctx context.Context, p *PersonalTodoResourceUploadParams) (*PersonalTodoResourceResult, error) {
	if p == nil {
		return &PersonalTodoResourceResult{Success: false, Error: "params is required"}, nil
	}
	if p.AppID == "" {
		return &PersonalTodoResourceResult{Success: false, Error: "app_id is required"}, nil
	}
	if p.Size <= 0 {
		return &PersonalTodoResourceResult{Success: false, Error: "size is required"}, nil
	}
	if p.Size > PersonalTodoResourceMaxSize {
		return &PersonalTodoResourceResult{Success: false, Error: "size exceeds the 9437184 byte limit"}, nil
	}
	if p.FileName == "" {
		return &PersonalTodoResourceResult{Success: false, Error: "file_name is required"}, nil
	}
	if p.ContentType == "" {
		return &PersonalTodoResourceResult{Success: false, Error: "content_type is required"}, nil
	}
	if p.FileData == "" {
		return &PersonalTodoResourceResult{Success: false, Error: "file_data is required"}, nil
	}
	if p.OrgID == "" {
		return &PersonalTodoResourceResult{Success: false, Error: "org_id is required"}, nil
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "personal_todos", "resource_update", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"appId":       p.AppID,
		"size":        p.Size,
		"fileName":    p.FileName,
		"contentType": p.ContentType,
		"fileData":    p.FileData,
		"orgId":       p.OrgID,
		"thumb":       p.Thumb,
	}
	if p.ExtensionInfo != "" {
		body["extensionInfo"] = p.ExtensionInfo
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalTodoResourceResult{Success: false, Error: err.Error()}, nil
	}
	res := &PersonalTodoResourceResult{Success: true, RawResponse: result}
	data := extractData(result)
	if data == nil {
		return res, nil
	}
	res.FileName = strFromMap(data, "fileName")
	res.MimeType = strFromMap(data, "mimeType")
	res.Suffix = strFromMap(data, "suffix")
	res.Size = int64(intFromMap(data, "size"))
	res.MD5 = strFromMap(data, "md5")
	res.ExtensionInfo = strFromMap(data, "extensionInfo")
	res.ResourceID = strFromMap(data, "resourceId")
	res.DownloadURL = strFromMap(data, "downloadUrl")
	if thumbnails, ok := data["imageThumbnailList"].(map[string]interface{}); ok {
		res.ImageThumbnailList = thumbnails
	}
	return res, nil
}

// FetchPersonalTodoResourceDownloadURL fetches a temporary resource download URL.
func (c *LansengerClient) FetchPersonalTodoResourceDownloadURL(ctx context.Context, resourceID, orgID, fileName, userToken string) (*PersonalTodoURLResult, error) {
	if resourceID == "" {
		return &PersonalTodoURLResult{Success: false, Error: "resource_id is required"}, nil
	}
	if orgID == "" {
		return &PersonalTodoURLResult{Success: false, Error: "org_id is required"}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "personal_todos", "resource_download", token, WithUserToken(userToken))
	body := map[string]interface{}{"resourceId": resourceID, "orgId": orgID}
	if fileName != "" {
		body["fileName"] = fileName
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalTodoURLResult{Success: false, Error: err.Error()}, nil
	}
	res := &PersonalTodoURLResult{Success: true, RawResponse: result}
	if value, ok := result["data"].(string); ok {
		res.URL = value
	}
	return res, nil
}

// FetchPersonalTodoResourceUploadURL fetches a presigned S3 upload URL.
func (c *LansengerClient) FetchPersonalTodoResourceUploadURL(ctx context.Context, fileName, md5 string, size int64, orgID, userToken string) (*PersonalTodoURLResult, error) {
	if fileName == "" {
		return &PersonalTodoURLResult{Success: false, Error: "file_name is required"}, nil
	}
	if md5 == "" {
		return &PersonalTodoURLResult{Success: false, Error: "md5 is required"}, nil
	}
	if size <= 0 {
		return &PersonalTodoURLResult{Success: false, Error: "size is required"}, nil
	}
	if orgID == "" {
		return &PersonalTodoURLResult{Success: false, Error: "org_id is required"}, nil
	}
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "personal_todos", "resource_upload_url", token, WithUserToken(userToken))
	result, err := c.doPost(ctx, url, map[string]interface{}{
		"fileName": fileName,
		"md5":      md5,
		"size":     size,
		"orgId":    orgID,
	})
	if err != nil {
		return &PersonalTodoURLResult{Success: false, Error: err.Error()}, nil
	}
	res := &PersonalTodoURLResult{Success: true, RawResponse: result}
	if value, ok := result["data"].(string); ok {
		res.URL = value
	}
	return res, nil
}
