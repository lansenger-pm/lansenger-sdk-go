package lansenger

import (
	"context"
	"fmt"
	"net/url"
)

func (c *LansengerClient) BuildAuthorizeURL(redirectURI string, scope string, state string) string {
	u, _ := url.Parse(c.config.PassportURL + "/oauth2/authorize")

	params := url.Values{}
	params.Set("appid", c.config.AppID)
	params.Set("response_type", "code")
	if scope != "" {
		params.Set("scope", scope)
	} else {
		params.Set("scope", OAuth2ScopeBasicUserInfo)
	}
	if state != "" {
		params.Set("state", state)
	}
	params.Set("redirect_uri", redirectURI)

	u.RawQuery = params.Encode()
	return u.String()
}

func ParseAuthorizeCallback(queryString string) (map[string]string, error) {
	values, err := url.ParseQuery(queryString)
	if err != nil {
		return nil, fmt.Errorf("parsing callback query string: %w", err)
	}

	result := map[string]string{}
	for key, vals := range values {
		if len(vals) > 0 {
			result[key] = vals[0]
		}
	}
	return result, nil
}

func ValidateCallbackState(callbackState, expectedState string) bool {
	return callbackState == expectedState
}

func (c *LansengerClient) ExchangeCode(ctx context.Context, code string, redirectURI string) (*UserTokenResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "oauth", "user_token_create", token,
		WithGrantType("authorization_code"),
		WithCode(code),
		WithRedirectURI(redirectURI),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &UserTokenResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &UserTokenResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &UserTokenResult{
		Success:          true,
		UserToken:        strFromMap(data, "userToken"),
		ExpiresIn:        intFromMap(data, "expiresIn"),
		RefreshToken:     strFromMap(data, "refreshToken"),
		RefreshExpiresIn: intFromMap(data, "refreshExpiresIn"),
		StaffID:          strFromMap(data, "staffId"),
		Scope:            strFromMap(data, "scope"),
		State:            strFromMap(data, "state"),
		RawResponse:      result,
	}, nil
}

func (c *LansengerClient) RefreshUserToken(ctx context.Context, refreshToken string, scope string) (*UserTokenResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}

	url := BuildAPIURL(c.config, "oauth", "refresh_token_create", token,
		WithGrantType("refresh_token"),
		WithRefreshToken(refreshToken),
		WithScope(scope),
	)

	result, err := c.doGet(ctx, url)
	if err != nil {
		return &UserTokenResult{Success: false, Error: err.Error()}, nil
	}

	data := extractData(result)
	if data == nil {
		return &UserTokenResult{Success: false, Error: "no data in response", RawResponse: result}, nil
	}

	return &UserTokenResult{
		Success:          true,
		UserToken:        strFromMap(data, "userToken"),
		ExpiresIn:        intFromMap(data, "expiresIn"),
		RefreshToken:     strFromMap(data, "refreshToken"),
		RefreshExpiresIn: intFromMap(data, "refreshExpiresIn"),
		StaffID:          strFromMap(data, "staffId"),
		Scope:            strFromMap(data, "scope"),
		RawResponse:      result,
	}, nil
}
