package lansenger

import (
	"context"
	"strconv"
	"sync"
	"time"
)

const UserTokenRefreshMargin = 300

type UserTokenManager struct {
	client      *LansengerClient
	store       *CredentialStore
	userToken   string
	refreshToken string
	expiresAt   time.Time
	refreshExpiresAt time.Time
	staffID     string
	mu          sync.Mutex
}

func NewUserTokenManager(client *LansengerClient, store *CredentialStore) *UserTokenManager {
	utm := &UserTokenManager{
		client: client,
		store:  store,
	}
	utm.loadFromStore()
	return utm
}

func (utm *UserTokenManager) loadFromStore() {
	if utm.store == nil {
		return
	}
	tokens, err := utm.store.LoadUserToken("")
	if err != nil {
		return
	}
	utm.userToken = tokens["user_token"]
	utm.refreshToken = tokens["refresh_token"]

	if expStr := tokens["user_token_expiry"]; expStr != "" {
		exp, err := strconv.ParseInt(expStr, 10, 64)
		if err == nil && exp > time.Now().Unix() {
			utm.expiresAt = time.Unix(exp, 0)
		}
	}
	if refExpStr := tokens["refresh_token_expiry"]; refExpStr != "" {
		refExp, err := strconv.ParseInt(refExpStr, 10, 64)
		if err == nil && refExp > 0 {
			utm.refreshExpiresAt = time.Unix(refExp, 0)
		}
	}

	utm.staffID = tokens["staff_id"]
}

func (utm *UserTokenManager) GetToken(ctx context.Context) (string, error) {
	utm.mu.Lock()
	if utm.userToken != "" && time.Now().Before(utm.expiresAt.Add(-time.Duration(UserTokenRefreshMargin)*time.Second)) {
		token := utm.userToken
		utm.mu.Unlock()
		return token, nil
	}
	utm.mu.Unlock()

	return utm.refresh(ctx)
}

func (utm *UserTokenManager) refresh(ctx context.Context) (string, error) {
	utm.mu.Lock()
	defer utm.mu.Unlock()

	if utm.refreshToken == "" {
		return "", NewAuthError("no refresh token available — must re-authorize via exchange_code")
	}

	if time.Now().Add(time.Duration(UserTokenRefreshMargin) * time.Second).After(utm.refreshExpiresAt) {
		return "", NewAuthError("refresh token expired — must re-authorize via exchange_code")
	}

	result, err := utm.client.RefreshUserToken(ctx, utm.refreshToken, "")
	if err != nil {
		return "", err
	}
	if !result.Success {
		return "", NewAuthError("user token refresh failed: " + result.Error)
	}

	utm.userToken = result.UserToken
	if result.RefreshToken != "" {
		utm.refreshToken = result.RefreshToken
	}
	effectiveExpiresIn := result.ExpiresIn - UserTokenRefreshMargin
	if effectiveExpiresIn <= 0 {
		effectiveExpiresIn = result.ExpiresIn / 2
	}
	utm.expiresAt = time.Now().Add(time.Duration(effectiveExpiresIn) * time.Second)
	if result.RefreshExpiresIn > 0 {
		utm.refreshExpiresAt = time.Now().Add(time.Duration(result.RefreshExpiresIn) * time.Second)
	}
	if result.StaffID != "" {
		utm.staffID = result.StaffID
	}

	if utm.store != nil {
		utm.store.SaveUserToken(utm.userToken, utm.refreshToken, result.ExpiresIn, result.RefreshExpiresIn, utm.staffID)
	}

	return utm.userToken, nil
}

func (utm *UserTokenManager) SetTokens(userToken, refreshToken string, expiresIn int, staffID string, refreshExpiresIn int) {
	utm.mu.Lock()
	defer utm.mu.Unlock()

	utm.userToken = userToken
	utm.refreshToken = refreshToken
	effectiveExpiresIn := expiresIn - UserTokenRefreshMargin
	if effectiveExpiresIn <= 0 {
		effectiveExpiresIn = expiresIn / 2
	}
	utm.expiresAt = time.Now().Add(time.Duration(effectiveExpiresIn) * time.Second)
	if refreshExpiresIn > 0 {
		utm.refreshExpiresAt = time.Now().Add(time.Duration(refreshExpiresIn) * time.Second)
	}
	if staffID != "" {
		utm.staffID = staffID
	}

	if utm.store != nil {
		utm.store.SaveUserToken(userToken, refreshToken, expiresIn, refreshExpiresIn, utm.staffID)
	}
}

func (utm *UserTokenManager) Invalidate() {
	utm.mu.Lock()
	utm.userToken = ""
	utm.expiresAt = time.Time{}
	utm.mu.Unlock()
}

func (utm *UserTokenManager) StaffID() string {
	utm.mu.Lock()
	defer utm.mu.Unlock()
	return utm.staffID
}

func (utm *UserTokenManager) RefreshToken() string {
	utm.mu.Lock()
	defer utm.mu.Unlock()
	return utm.refreshToken
}

func (utm *UserTokenManager) RefreshTokenExpiry() time.Time {
	utm.mu.Lock()
	defer utm.mu.Unlock()
	return utm.refreshExpiresAt
}

func (c *LansengerClient) GetUserToken(ctx context.Context) (string, error) {
	if c.userTokenMgr == nil {
		return "", NewAuthError("UserTokenManager not initialized — use SetUserTokens or exchange_code first")
	}
	return c.userTokenMgr.GetToken(ctx)
}

func (c *LansengerClient) SetUserTokens(userToken, refreshToken string, expiresIn int, staffID string, refreshExpiresIn int) {
	if c.userTokenMgr == nil {
		c.userTokenMgr = NewUserTokenManager(c, nil)
	}
	c.userTokenMgr.SetTokens(userToken, refreshToken, expiresIn, staffID, refreshExpiresIn)
}