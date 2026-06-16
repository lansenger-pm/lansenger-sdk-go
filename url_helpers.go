package lansenger

import (
	"fmt"
	"net/url"
	"strings"
	"sync"
)

func BuildAPIURL(cfg *Config, category, endpoint, appToken string, opts ...URLOption) string {
	path, ok := APIEndpoints[category][endpoint]
	if !ok {
		return ""
	}

	path = substitutePathVars(path, opts...)

	baseURL := strings.TrimRight(cfg.APIGatewayURL, "/") + path

	params := url.Values{}
	params.Set("app_token", appToken)

	for _, opt := range opts {
		opt.Apply(params)
	}

	if len(params) > 0 {
		return baseURL + "?" + params.Encode()
	}
	return baseURL
}

type URLOption struct {
	Key   string
	Value string
	Apply func(v url.Values)
}

// defaultUserToken is set by CLI tools (e.g. `lansenger --as staff_id`)
// to inject a user token into API requests when the explicit userToken
// parameter is empty. Library users should prefer SetUserTokens().
//
// Access is protected by defaultUserTokenMu (RWMutex).
// Use SetDefaultUserToken() / getDefaultUserToken() for thread-safe access.
var (
	defaultUserTokenMu sync.RWMutex
	defaultUserToken   string
)

// SetDefaultUserToken sets the fallback user token for API requests
// where no explicit userToken is provided. Thread-safe.
func SetDefaultUserToken(token string) {
	defaultUserTokenMu.Lock()
	defaultUserToken = token
	defaultUserTokenMu.Unlock()
}

func getDefaultUserToken() string {
	defaultUserTokenMu.RLock()
	defer defaultUserTokenMu.RUnlock()
	return defaultUserToken
}

func WithUserToken(token string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		t := token
		if t == "" {
			t = getDefaultUserToken()
		}
		if t != "" {
			v.Set("user_token", t)
		}
	}}
}

func WithUserID(id string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if id != "" {
			v.Set("user_id", id)
		}
	}}
}

func WithQueryParam(key, value string) URLOption {
	return URLOption{Key: key, Value: value, Apply: func(v url.Values) {
		if value != "" {
			v.Set(key, value)
		}
	}}
}

func WithPathVar(key, value string) URLOption {
	return URLOption{Key: key, Value: value, Apply: func(v url.Values) {}}
}

func WithMediaType(mediaType int) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("type", fmt.Sprintf("%d", mediaType))
	}}
}

func WithMediaTypeString(mediaType string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if mediaType != "" {
			v.Set("type", mediaType)
		}
	}}
}

func WithIntParam(key string, value int) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if value > 0 {
			v.Set(key, fmt.Sprintf("%d", value))
		}
	}}
}

func WithPage(page int) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("page", fmt.Sprintf("%d", page))
	}}
}

func WithPageSize(pageSize int) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("page_size", fmt.Sprintf("%d", pageSize))
	}}
}

func WithPageOffset(offset int) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("page_offset", fmt.Sprintf("%d", offset))
	}}
}

func WithTagID(tagID string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if tagID != "" {
			v.Set("tag_id", tagID)
		}
	}}
}

func WithStaffID(staffID string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if staffID != "" {
			v.Set("staff_id", staffID)
		}
	}}
}

func WithGrantType(grantType string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("grant_type", grantType)
	}}
}

func WithAppID(appID string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("appid", appID)
	}}
}

func WithSecret(secret string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("secret", secret)
	}}
}

func WithCode(code string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if code != "" {
			v.Set("code", code)
		}
	}}
}

func WithRedirectURI(uri string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if uri != "" {
			v.Set("redirect_uri", uri)
		}
	}}
}

func WithRefreshToken(token string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if token != "" {
			v.Set("refresh_token", token)
		}
	}}
}

func WithScope(scope string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if scope != "" {
			v.Set("scope", scope)
		}
	}}
}

func WithOrgID(orgID string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		if orgID != "" {
			v.Set("org_id", orgID)
		}
	}}
}

func WithIDType(idType string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("id_type", idType)
	}}
}

func WithIDValue(idValue string) URLOption {
	return URLOption{Apply: func(v url.Values) {
		v.Set("id_value", idValue)
	}}
}

func substitutePathVars(path string, opts ...URLOption) string {
	for _, opt := range opts {
		if opt.Key != "" && opt.Value != "" {
			path = strings.ReplaceAll(path, "{"+opt.Key+"}", opt.Value)
		}
	}
	return path
}
