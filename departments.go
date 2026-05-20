package lansenger

import (
	"context"
)

func (c *LansengerClient) FetchDepartmentDetail(ctx context.Context, departmentID, userToken, tagID string) (*DepartmentDetailResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "departments", "detail_fetch", token,
		WithUserToken(userToken),
		WithPathVar("department_id", departmentID),
		WithTagID(tagID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &DepartmentDetailResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &DepartmentDetailResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &DepartmentDetailResult{
		Success:         true,
		ID:              strFromMap(data, "id"),
		Name:            strFromMap(data, "name"),
		ExternalID:      strFromMap(data, "externalId"),
		ParentID:        strFromMap(data, "parentId"),
		Order:           intFromMap(data, "order"),
		HasChildren:     boolFromMap(data, "hasChildren"),
		NormalMembers:   intFromMap(data, "normalMembers"),
		InactiveMembers: intFromMap(data, "inactiveMembers"),
		FrozenMembers:   intFromMap(data, "frozenMembers"),
		DeletedMembers:  intFromMap(data, "deletedMembers"),
		DeptType:        strFromMap(data, "deptType"),
		RawResponse:     result,
	}, nil
}

func (c *LansengerClient) FetchDepartmentChildren(ctx context.Context, departmentID, userToken string) (*DepartmentChildrenResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "departments", "children_fetch", token,
		WithUserToken(userToken),
		WithPathVar("department_id", departmentID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &DepartmentChildrenResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &DepartmentChildrenResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &DepartmentChildrenResult{
		Success:     true,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchDepartmentStaffs(ctx context.Context, departmentID, userToken string, page, pageSize int) (*DepartmentStaffsResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "departments", "staffs_fetch", token,
		WithUserToken(userToken),
		WithPathVar("department_id", departmentID),
		WithPage(page),
		WithPageSize(pageSize),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &DepartmentStaffsResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &DepartmentStaffsResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &DepartmentStaffsResult{
		Success:     true,
		HasMore:     boolFromMap(data, "hasMore"),
		Total:       intFromMap(data, "total"),
		RawResponse: result,
	}, nil
}
