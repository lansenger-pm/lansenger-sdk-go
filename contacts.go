package lansenger

import (
	"context"
)

func (c *LansengerClient) FetchStaffBasicInfo(ctx context.Context, staffID string, userToken string) (*StaffBasicInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "staffs", "basic_info_fetch", token,
		WithUserToken(userToken),
		WithPathVar("staff_id", staffID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &StaffBasicInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffBasicInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &StaffBasicInfoResult{
		Success:     true,
		OrgID:       strFromMap(data, "orgId"),
		OrgName:     strFromMap(data, "orgName"),
		Name:        strFromMap(data, "name"),
		Gender:      strFromMap(data, "gender"),
		Signature:   strFromMap(data, "signature"),
		AvatarURL:   strFromMap(data, "avatarUrl"),
		AvatarID:    strFromMap(data, "avatarId"),
		Status:      strFromMap(data, "status"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchStaffDetail(ctx context.Context, staffID string, userToken string) (*StaffDetailResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "staffs", "detail_fetch", token,
		WithUserToken(userToken),
		WithPathVar("staff_id", staffID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &StaffDetailResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffDetailResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &StaffDetailResult{
		Success:        true,
		Name:           strFromMap(data, "name"),
		Signature:      strFromMap(data, "signature"),
		AvatarID:       strFromMap(data, "avatarId"),
		AvatarURL:      strFromMap(data, "avatarUrl"),
		Status:         strFromMap(data, "status"),
		Gender:         strFromMap(data, "gender"),
		OrgID:          strFromMap(data, "orgId"),
		OrgName:        strFromMap(data, "orgName"),
		LoginName:      strFromMap(data, "loginName"),
		EmployeeNumber: strFromMap(data, "employeeNumber"),
		Email:          strFromMap(data, "email"),
		ExternalID:     strFromMap(data, "externalId"),
		MobilePhone:    strFromMap(data, "mobilePhone"),
		RawResponse:    result,
	}, nil
}

func (c *LansengerClient) FetchDepartmentAncestors(ctx context.Context, staffID string, userToken string) (*DepartmentAncestorsResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "staffs", "department_ancestors_fetch", token,
		WithUserToken(userToken),
		WithPathVar("staff_id", staffID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &DepartmentAncestorsResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &DepartmentAncestorsResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &DepartmentAncestorsResult{
		Success:     true,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchStaffIdMapping(ctx context.Context, orgID, idType, idValue, userToken string) (*StaffIdMappingResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "staffs", "id_mapping_fetch", token,
		WithUserToken(userToken),
		WithOrgID(orgID),
		WithIDType(idType),
		WithIDValue(idValue),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &StaffIdMappingResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffIdMappingResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &StaffIdMappingResult{
		Success:     true,
		StaffID:     strFromMap(data, "staffId"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchOrgExtraFieldIDs(ctx context.Context, orgID, userToken string, page, pageSize int) (*ExtraFieldIdsResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "org", "extra_field_ids_fetch", token,
		WithUserToken(userToken),
		WithPathVar("org_id", orgID),
		WithPage(page),
		WithPageSize(pageSize),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &ExtraFieldIdsResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &ExtraFieldIdsResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &ExtraFieldIdsResult{
		Success:     true,
		HasMore:     boolFromMap(data, "hasMore"),
		Total:       intFromMap(data, "total"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) SearchStaff(ctx context.Context, keyword, userToken, userID string, recursive bool, sectorIDs []string, page, pageSize int) (*StaffSearchResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "staffs", "search", token,
		WithUserToken(userToken),
		WithUserID(userID),
		WithPage(page),
		WithPageSize(pageSize),
	)

	searchScope := map[string]interface{}{}
	if len(sectorIDs) > 0 {
		ids := make([]interface{}, len(sectorIDs))
		for i, id := range sectorIDs {
			ids[i] = id
		}
		searchScope["sectorIds"] = ids
	}

	body := map[string]interface{}{
		"keyword":     keyword,
		"recursive":   recursive,
		"searchScope": searchScope,
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &StaffSearchResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffSearchResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &StaffSearchResult{
		Success:     true,
		HasMore:     boolFromMap(data, "hasMore"),
		Total:       intFromMap(data, "total"),
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchOrgInfo(ctx context.Context, orgID, userToken string) (*OrgInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "org", "info_fetch", token,
		WithUserToken(userToken),
		WithPathVar("org_id", orgID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &OrgInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &OrgInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &OrgInfoResult{
		Success:     true,
		OrgID:       strFromMap(data, "orgId"),
		OrgName:     strFromMap(data, "orgName"),
		IconURL:     strFromMap(data, "iconUrl"),
		RawResponse: result,
	}, nil
}
