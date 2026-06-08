package lansenger

import (
	"os"
	"strconv"
)

type Config struct {
	AppID         string
	AppSecret     string
	APIGatewayURL string
	PassportURL   string
	HTTPTimeout   float64
	EncodingKey   string
	CallbackToken string
}

func NewConfig(appID, appSecret string) *Config {
	return &Config{
		AppID:         appID,
		AppSecret:     appSecret,
		APIGatewayURL: DefaultAPIGatewayURL,
		PassportURL:   DefaultPassportURL,
		HTTPTimeout:   30.0,
	}
}

func ConfigFromEnv() (*Config, error) {
	appID := os.Getenv("LANSENGER_APP_ID")
	appSecret := os.Getenv("LANSENGER_APP_SECRET")
	if appID == "" || appSecret == "" {
		return nil, NewConfigError("LANSENGER_APP_ID and LANSENGER_APP_SECRET environment variables are required")
	}
	cfg := &Config{
		AppID:         appID,
		AppSecret:     appSecret,
		APIGatewayURL: getEnvOrDefault("LANSENGER_API_GATEWAY_URL", DefaultAPIGatewayURL),
		PassportURL:   getEnvOrDefault("LANSENGER_PASSPORT_URL", DefaultPassportURL),
		HTTPTimeout:   30.0,
		EncodingKey:   os.Getenv("LANSENGER_ENCODING_KEY"),
		CallbackToken: os.Getenv("LANSENGER_CALLBACK_TOKEN"),
	}
	timeoutStr := os.Getenv("LANSENGER_HTTP_TIMEOUT")
	if timeoutStr != "" {
		timeout, err := strconv.ParseFloat(timeoutStr, 64)
		if err == nil {
			cfg.HTTPTimeout = timeout
		}
	}
	return cfg, nil
}

func (c *Config) IsConfigured() bool {
	return c.AppID != "" && c.AppSecret != "" && c.APIGatewayURL != ""
}

func (c *Config) HasPassportURL() bool {
	return c.PassportURL != ""
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}
