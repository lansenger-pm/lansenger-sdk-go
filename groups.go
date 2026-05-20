package lansenger

import (
	"context"
)

func (c *LansengerClient) CreateGroup(ctx context.Context, info *GroupCreateInfo, userToken string) (*CreateGroupResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{
		"name":    info.Name,
		"orgId":   info.OrgID,
		"ownerId": info.OwnerID,
	}
	if info.Description != "" {
		body["description"] = info.Description
	}
	if info.AvatarID != "" {
		body["avatarId"] = info.AvatarID
	}
	if len(info.StaffIDList) > 0 {
		body["staffIdList"] = info.StaffIDList
	}
	if len(info.DepartmentIDList) > 0 {
		body["departmentIdList"] = info.DepartmentIDList
	}
	if info.ApplyRequestID != "" {
		body["applyRequestId"] = info.ApplyRequestID
	}
	if info.ApplyNotes != "" {
		body["applyNotes"] = info.ApplyNotes
	}
	if info.ApplyGlobalUniqueID != "" {
		body["applyGlobalUniqueId"] = info.ApplyGlobalUniqueID
	}
	if info.ApplySessionUniqueID != "" {
		body["applySessionUniqueId"] = info.ApplySessionUniqueID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &CreateGroupResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &CreateGroupResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &CreateGroupResult{
		Success:      true,
		GroupID:      strFromMap(data, "groupId"),
		TotalMembers: intFromMap(data, "totalMembers"),
		RawResponse:  result,
	}, nil
}

func (c *LansengerClient) FetchGroupInfo(ctx context.Context, groupID, userToken string) (*GroupInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "info_fetch", token,
		WithUserToken(userToken),
		WithPathVar("group_id", groupID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &GroupInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &GroupInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &GroupInfoResult{
		Success:      true,
		Name:         strFromMap(data, "name"),
		Description:  strFromMap(data, "description"),
		AvatarID:     strFromMap(data, "avatarId"),
		AvatarURL:    strFromMap(data, "avatarUrl"),
		Owner:        strFromMap(data, "owner"),
		Creator:      strFromMap(data, "creator"),
		State:        strFromMap(data, "state"),
		ManageMode:   strFromMap(data, "manageMode"),
		TotalMembers: intFromMap(data, "totalMembers"),
		RawResponse:  result,
	}, nil
}

func (c *LansengerClient) FetchGroupMembers(ctx context.Context, groupID, userToken string, pageOffset, pageSize int) (*GroupMemberResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "members_fetch", token,
		WithUserToken(userToken),
		WithPathVar("group_id", groupID),
		WithPageOffset(pageOffset),
		WithPageSize(pageSize),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &GroupMemberResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &GroupMemberResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &GroupMemberResult{
		Success:      true,
		TotalMembers: intFromMap(data, "totalMembers"),
		RawResponse:  result,
	}, nil
}

func (c *LansengerClient) FetchGroupList(ctx context.Context, userToken string, pageOffset, pageSize int) (*GroupListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "list_fetch", token,
		WithUserToken(userToken),
		WithPageOffset(pageOffset),
		WithPageSize(pageSize),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &GroupListResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &GroupListResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &GroupListResult{
		Success:       true,
		TotalGroupIDs: intFromMap(data, "totalGroupIds"),
		RawResponse:   result,
	}, nil
}

func (c *LansengerClient) CheckIsInGroup(ctx context.Context, groupID, userToken, staffID string) (*IsInGroupResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "is_in_group", token,
		WithUserToken(userToken),
		WithPathVar("group_id", groupID),
		WithStaffID(staffID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &IsInGroupResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &IsInGroupResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &IsInGroupResult{
		Success:     true,
		IsInGroup:   boolFromMap(data, "isInGroup"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) UpdateGroupInfo(ctx context.Context, groupID string, params map[string]interface{}, userToken string) (*UpdateGroupResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "info_update", token,
		WithUserToken(userToken),
		WithPathVar("group_id", groupID),
	)

	result, err := c.doPost(ctx, url, params)
	if err != nil {
		return &UpdateGroupResult{Success: false, Error: err.Error()}, nil
	}

	return &UpdateGroupResult{
		Success:     true,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) UpdateGroupMembers(ctx context.Context, groupID string, addUserList, delUserList, addDepartmentIDList []string, userToken string) (*UpdateGroupMembersResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "groups_v2", "members_update", token,
		WithUserToken(userToken),
		WithPathVar("group_id", groupID),
	)

	body := map[string]interface{}{}
	if len(addUserList) > 0 {
		body["addUserList"] = addUserList
	}
	if len(delUserList) > 0 {
		body["delUserList"] = delUserList
	}
	if len(addDepartmentIDList) > 0 {
		body["addDepartmentIdList"] = addDepartmentIDList
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &UpdateGroupMembersResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)

	res := &UpdateGroupMembersResult{
		Success:     true,
		RawResponse: result,
	}
	if data != nil {
		res.TotalMembers = intFromMap(data, "totalMembers")
		res.AddedStaffCount = intFromMap(data, "addedStaffCount")
		res.DeletedStaffCount = intFromMap(data, "deletedStaffCount")
	}
	return res, nil
}
