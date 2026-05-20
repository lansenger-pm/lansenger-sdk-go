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

	return &UserInfoResult{
		Success:        true,
		StaffID:        strFromMap(data, "staffId"),
		Name:           strFromMap(data, "name"),
		OrgID:          strFromMap(data, "orgId"),
		OrgName:        strFromMap(data, "orgname"),
		AvatarID:       strFromMap(data, "avatarId"),
		AvatarURL:      strFromMap(data, "avatarUrl"),
		MobilePhone:    strFromMap(data, "mobilePhone"),
		Email:          strFromMap(data, "email"),
		EmployeeNumber: strFromMap(data, "employeeNumber"),
		LoginName:      strFromMap(data, "loginName"),
		ExternalID:     strFromMap(data, "externalId"),
		RawResponse:    result,
	}, nil
}
