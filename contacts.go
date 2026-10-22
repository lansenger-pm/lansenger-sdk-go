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
		return &StaffBasicInfoResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffBasicInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &StaffBasicInfoResult{
		Success:     true,
		OrgID:       strFromMap(data, "orgId"),
		OrgName:     strFromMap(data, "orgName"),
		Name:        strFromMap(data, "name"),
		Gender:      strFromMap(data, "gender"),
		Signature:   strFromMap(data, "signature"),
		AvatarURL:   strFromMap(data, "avatar"),
		AvatarID:    strFromMap(data, "avatarId"),
		Status:      strFromMap(data, "status"),
		RawResponse: result,
	}
	if departments, ok := data["departments"].([]interface{}); ok {
		res.Departments = make([]map[string]interface{}, 0, len(departments))
		for _, item := range departments {
			if m, ok := item.(map[string]interface{}); ok {
				res.Departments = append(res.Departments, m)
			}
		}
	}
	return res, nil
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
		return &StaffDetailResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffDetailResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &StaffDetailResult{
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
		JoinDate:       strFromMap(data, "joinDate"),
		RawResponse:    result,
	}

	phoneObj := mapFromMap(data, "mobilePhone")
	if phoneObj != nil {
		res.MobilePhone = strFromMap(phoneObj, "number")
		res.MobilePhoneCountryCode = strFromMap(phoneObj, "countryCode")
	}

	if departments, ok := data["departments"].([]interface{}); ok {
		res.Departments = make([]map[string]interface{}, 0, len(departments))
		for _, item := range departments {
			if m, ok := item.(map[string]interface{}); ok {
				res.Departments = append(res.Departments, m)
			}
		}
	}
	if duties, ok := data["duties"].([]interface{}); ok {
		res.Duties = make([]map[string]interface{}, 0, len(duties))
		for _, item := range duties {
			if m, ok := item.(map[string]interface{}); ok {
				res.Duties = append(res.Duties, m)
			}
		}
	}
	if parties, ok := data["parties"].([]interface{}); ok {
		res.Parties = make([]map[string]interface{}, 0, len(parties))
		for _, item := range parties {
			if m, ok := item.(map[string]interface{}); ok {
				res.Parties = append(res.Parties, m)
			}
		}
	}
	if tags, ok := data["tags"].([]interface{}); ok {
		res.Tags = make([]map[string]interface{}, 0, len(tags))
		for _, item := range tags {
			if m, ok := item.(map[string]interface{}); ok {
				res.Tags = append(res.Tags, m)
			}
		}
	}

	return res, nil
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
		return &DepartmentAncestorsResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
	}

	res := &DepartmentAncestorsResult{
		Success:     true,
		RawResponse: result,
	}

	// 真实响应把 data 作为数组返回，每个元素带一个 ancestorDepartments 对象数组：
	// {"data":[{"ancestorDepartments":[{"departmentId":…,"externalId":…,"name":…}, …]}, …]}
	// （Python/TypeScript SDK 同样把 data 当列表处理。）
	if entries, ok := result["data"].([]interface{}); ok {
		res.AncestorGroups = make([][]map[string]string, 0, len(entries))
		for _, entry := range entries {
			ancestors := []map[string]string{}
			if m, ok := entry.(map[string]interface{}); ok {
				if arr, ok := m["ancestorDepartments"].([]interface{}); ok {
					ancestors = make([]map[string]string, 0, len(arr))
					for _, item := range arr {
						if dm, ok := item.(map[string]interface{}); ok {
							ancestors = append(ancestors, ancestorDepartmentEntry(dm))
						}
					}
				}
			}
			res.AncestorGroups = append(res.AncestorGroups, ancestors)
		}
		return res, nil
	}

	data := extractData(result)
	if data == nil {
		return &DepartmentAncestorsResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	// 兼容旧形状：data.ancestorDepartments 是数组的数组。
	if ancestorList, ok := data["ancestorDepartments"].([]interface{}); ok {
		res.AncestorGroups = make([][]map[string]string, 0, len(ancestorList))
		for _, group := range ancestorList {
			if arr, ok := group.([]interface{}); ok {
				ancestors := make([]map[string]string, 0, len(arr))
				for _, item := range arr {
					if m, ok := item.(map[string]interface{}); ok {
						ancestors = append(ancestors, ancestorDepartmentEntry(m))
					}
				}
				res.AncestorGroups = append(res.AncestorGroups, ancestors)
			}
		}
	}

	return res, nil
}

// ancestorDepartmentEntry 归一化一个部门条目；服务端实际用 departmentId 作为部门标识。
func ancestorDepartmentEntry(m map[string]interface{}) map[string]string {
	id := strFromMapLoose(m, "departmentId", "id")
	entry := map[string]string{
		"id":           id,
		"departmentId": id,
		"name":         strFromMap(m, "name"),
	}
	if extID := strFromMap(m, "externalId"); extID != "" {
		entry["externalId"] = extID
	}
	return entry
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
		return &StaffIdMappingResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
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
		return &ExtraFieldIdsResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
	}

	data := extractData(result)
	if data == nil {
		return &ExtraFieldIdsResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &ExtraFieldIdsResult{
		Success:     true,
		HasMore:     boolFromMap(data, "hasMore"),
		Total:       intFromMap(data, "total"),
		RawResponse: result,
	}
	res.ExtraFieldIDs = stringArrayFromMap(data, "extraFieldIds")
	return res, nil
}

func (c *LansengerClient) SearchStaff(ctx context.Context, keyword, userToken, userID string, recursive bool, sectorIDs []string, page, pageSize int) (*StaffSearchResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	opts := []URLOption{WithUserToken(userToken), WithUserID(userID)}
	// Server requires page and page_size together (doc 4.1.16 v2); sending only
	// one is silently ignored, so default page=1/page_size=20 when either is
	// given (page/pageSize == 0 means "not given").
	// NOTE: HasMore may be unreliable (observed always true, LXBUGS-128510);
	// prefer comparing len(result.StaffInfos) against result.Total for
	// pagination termination.
	if page > 0 || pageSize > 0 {
		if page <= 0 {
			page = 1
		}
		if pageSize <= 0 {
			pageSize = 20
		}
		opts = append(opts, WithPage(page), WithPageSize(pageSize))
	}

	url := BuildAPIURL(c.config, "staffs", "search", token, opts...)

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
		return &StaffSearchResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
	}

	data := extractData(result)
	if data == nil {
		return &StaffSearchResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &StaffSearchResult{
		Success:     true,
		HasMore:     boolFromMap(data, "hasMore"),
		Total:       intFromMap(data, "total"),
		RawResponse: result,
	}
	if staffInfo, ok := data["staffInfo"].([]interface{}); ok {
		res.StaffInfo = make([]map[string]interface{}, 0, len(staffInfo))
		for _, item := range staffInfo {
			if m, ok := item.(map[string]interface{}); ok {
				res.StaffInfo = append(res.StaffInfo, m)
			}
		}
	}
	return res, nil
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
		return &OrgInfoResult{Success: false, Error: err.Error(), ErrCode: ErrCodeOf(err)}, nil
	}

	data := extractData(result)
	if data == nil {
		return &OrgInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &OrgInfoResult{
		Success:           true,
		OrgID:             strFromMap(data, "orgId"),
		OrgName:           strFromMap(data, "orgName"),
		IconURL:           strFromMap(data, "iconUrl"),
		OrgMaxMemberLimit: intFromMap(data, "orgMaxMemberLimit"),
		OrgOrderType:      strFromMap(data, "orgOrderType"),
		OrgDaysLimit:      intFromMap(data, "orgDaysLimit"),
		OrgBillingDate:    strFromMap(data, "orgBillingDate"),
		RawResponse:       result,
	}, nil
}
