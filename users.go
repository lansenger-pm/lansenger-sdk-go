package lansenger

import (
	"context"
	"fmt"
)

func (c *LansengerClient) FetchUserInfo(ctx context.Context, userToken string) (*UserInfoResult, error) {
	if userToken == "" {
		return nil, fmt.Errorf("userToken is required for fetch_user_info")
	}

	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "users", "fetch", token,
		WithUserToken(userToken),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &UserInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &UserInfoResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	res := &UserInfoResult{
		Success:        true,
		StaffID:        strFromMap(data, "staffId"),
		Name:           strFromMap(data, "name"),
		OrgID:          strFromMap(data, "orgId"),
		OrgName:        strOrNil(data, "orgname", "orgName"),
		AvatarID:       strFromMap(data, "avatarId"),
		AvatarURL:      strFromMap(data, "avatar"),
		Email:          strFromMap(data, "email"),
		EmployeeNumber: strFromMap(data, "employeeNumber"),
		LoginName:      strFromMap(data, "loginName"),
		ExternalID:     strFromMap(data, "externalId"),
		RawResponse:    result,
	}
	if departments, ok := data["department"].([]interface{}); ok {
		res.Departments = make([]map[string]interface{}, 0, len(departments))
		for _, item := range departments {
			if m, ok := item.(map[string]interface{}); ok {
				res.Departments = append(res.Departments, m)
			}
		}
	}
	return res, nil
}
