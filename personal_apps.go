package lansenger

import (
	"context"
)

func (c *LansengerClient) CreatePersonalApp(ctx context.Context, userToken, name, avatarID, description string) (*PersonalAppCreateResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "personal_apps", "create", token,
		WithUserToken(userToken),
	)

	body := map[string]interface{}{}
	if name != "" {
		body["name"] = name
	}
	if avatarID != "" {
		body["avatarId"] = avatarID
	}
	if description != "" {
		body["description"] = description
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalAppCreateResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	return &PersonalAppCreateResult{
		Success:      true,
		AppID:        strFromMap(data, "id"),
		Secret:       strFromMap(data, "secret"),
		APIGWAddr:    strFromMap(data, "apigwAddr"),
		PassportAddr: strFromMap(data, "passportAddr"),
		RawResponse:  result,
	}, nil
}

func (c *LansengerClient) UpdatePersonalApp(ctx context.Context, appID, userToken, name, avatarID, description string) (*PersonalAppInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "personal_apps", "update", token,
		WithUserToken(userToken),
		WithPathVar("app_id", appID),
	)

	body := map[string]interface{}{
		"name": name,
	}
	if avatarID != "" {
		body["avatarId"] = avatarID
	}
	if description != "" {
		body["description"] = description
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &PersonalAppInfoResult{Success: false, Error: err.Error()}, nil
	}

	return &PersonalAppInfoResult{
		Success:     true,
		AppID:       appID,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchPersonalApp(ctx context.Context, appID, userToken string) (*PersonalAppInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "personal_apps", "fetch", token,
		WithUserToken(userToken),
		WithPathVar("app_id", appID),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &PersonalAppInfoResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	return &PersonalAppInfoResult{
		Success:      true,
		AppID:        appID,
		Name:         strFromMap(data, "name"),
		AvatarID:     strFromMap(data, "avatarId"),
		Description:  strFromMap(data, "description"),
		APIGWAddr:    strFromMap(data, "apigwAddr"),
		PassportAddr: strFromMap(data, "passportAddr"),
		RawResponse:  result,
	}, nil
}

func (c *LansengerClient) DeletePersonalApp(ctx context.Context, appID, userToken string) (*PersonalAppInfoResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "personal_apps", "delete", token,
		WithUserToken(userToken),
		WithPathVar("app_id", appID),
	)

	result, err := c.doPost(ctx, url, map[string]interface{}{})
	if err != nil {
		return &PersonalAppInfoResult{Success: false, Error: err.Error()}, nil
	}

	return &PersonalAppInfoResult{
		Success:     true,
		AppID:       appID,
		RawResponse: result,
	}, nil
}

func (c *LansengerClient) FetchPersonalAppList(ctx context.Context, userToken string) (*PersonalAppListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "personal_apps", "list_fetch", token,
		WithUserToken(userToken),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &PersonalAppListResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	res := &PersonalAppListResult{
		Success:     true,
		RawResponse: result,
	}
	if apps, ok := data["appList"].([]interface{}); ok {
		for _, a := range apps {
			if m, ok := a.(map[string]interface{}); ok {
				res.AppList = append(res.AppList, m)
			}
		}
	}

	return res, nil
}
