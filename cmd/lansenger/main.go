package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var (
	jsonOutput      bool
	profileName     string
	globalAsStaffID string
	globalAppToken  string
	globalUserToken string
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print SDK/CLI version",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lansenger.Version)
	},
}

var rootCmd = &cobra.Command{
	Use:   "lansenger",
	Short: "Lansenger (蓝信) CLI — interact with Lansenger APIs from the command line",
}

func init() {
	rootCmd.Version = lansenger.Version
	rootCmd.PersistentFlags().BoolVarP(&jsonOutput, "json", "j", false, "Output raw JSON instead of formatted tables")
	rootCmd.PersistentFlags().StringVarP(&profileName, "profile", "P", "default", "Credential profile to use")
	rootCmd.PersistentFlags().StringVar(&globalAsStaffID, "as", "", "Act as the given staff_id (auto-loads user token from credential store)")
	rootCmd.PersistentFlags().StringVar(&globalAppToken, "app-token", "", "App access token (external mode — no auto-refresh)")
	rootCmd.PersistentFlags().StringVar(&globalUserToken, "user-token", "", "User access token (external mode — no auto-refresh)")
	rootCmd.AddCommand(versionCmd)
}

func main() {
	rootCmd.Execute()
}

func getClient() *lansenger.LansengerClient {
	// External mode: when --app-token is provided, skip credential file entirely.
	// The caller manages token lifecycle; no auto-refresh.
	if globalAppToken != "" {
		cfg := &lansenger.Config{
			APIGatewayURL: getEnvOrDefault("LANSENGER_API_GATEWAY_URL", "https://open.e.lanxin.cn/open/apigw"),
			AppToken:      globalAppToken,
			UserToken:     globalUserToken,
		}
		return lansenger.NewClientWithConfig(cfg)
	}

	store, err := lansenger.NewCredentialStore("", profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}
	creds, err := store.LoadCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading credentials: %v\n", err)
		os.Exit(1)
	}
	if creds["app_id"] == "" || creds["app_secret"] == "" {
		cfg, envErr := lansenger.ConfigFromEnv()
		if envErr == nil && cfg.IsConfigured() {
			client := lansenger.NewClientWithConfig(cfg)
			injectUserToken(client, store)
			return client
		}
		fmt.Fprintf(os.Stderr, "Error: No credentials configured for profile '%s'. Run `lansenger config set` first, or set LANSENGER_APP_ID / LANSENGER_APP_SECRET env vars, or use --app-token for external token mode.\n", profileName)
		os.Exit(1)
		return nil
	}
	cfg := lansenger.NewConfig(creds["app_id"], creds["app_secret"])
	if creds["api_gateway_url"] != "" {
		cfg.APIGatewayURL = creds["api_gateway_url"]
	}
	if creds["passport_url"] != "" {
		cfg.PassportURL = creds["passport_url"]
	}
	if creds["encoding_key"] != "" {
		cfg.EncodingKey = creds["encoding_key"]
	}
	if creds["callback_token"] != "" {
		cfg.CallbackToken = creds["callback_token"]
	}
	if globalAppToken != "" {
		cfg.AppToken = globalAppToken
	}
	client := lansenger.NewClientWithConfig(cfg)
	injectUserToken(client, store)
	return client
}

func getEnvOrDefault(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func injectUserToken(client *lansenger.LansengerClient, store *lansenger.CredentialStore) {
	if globalAsStaffID == "" {
		return
	}

	tokens, err := store.LoadUserToken(globalAsStaffID)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading user token for staff_id '%s': %v\n", globalAsStaffID, err)
		os.Exit(1)
	}

	userToken := tokens["user_token"]
	refreshToken := tokens["refresh_token"]

	if userToken == "" {
		fmt.Fprintf(os.Stderr, "Error: No user token found for staff_id '%s' in credential store. Run `lansenger oauth authorize` first.\n", globalAsStaffID)
		os.Exit(1)
	}

	ctx := context.Background()

	// Check if token is expired — refresh if needed
	expiryStr := tokens["user_token_expiry"]
	if expiryStr != "" {
		exp, parseErr := strconv.ParseInt(expiryStr, 10, 64)
		if parseErr == nil && exp <= 0 {
			// No valid expiry, refresh
			result, refreshErr := client.RefreshUserToken(ctx, refreshToken, "")
			if refreshErr != nil {
				fmt.Fprintf(os.Stderr, "Error refreshing user token for staff_id '%s': %v\n", globalAsStaffID, refreshErr)
				os.Exit(1)
			}
			if !result.Success {
				fmt.Fprintf(os.Stderr, "Error refreshing user token for staff_id '%s': %s\n", globalAsStaffID, result.Error)
				os.Exit(1)
			}
			userToken = result.UserToken
			if result.RefreshToken != "" {
				refreshToken = result.RefreshToken
			}
			store.SaveUserToken(userToken, refreshToken, result.ExpiresIn, result.RefreshExpiresIn, globalAsStaffID)
		}
	}

	// Register token with client's UserTokenManager for auto-refresh
	client.SetUserTokens(userToken, refreshToken, 7200, globalAsStaffID, 0)

	// Set default so WithUserToken falls back to this when userToken="" is passed
	lansenger.SetDefaultUserToken(userToken)
}

func getStore() *lansenger.CredentialStore {
	store, err := lansenger.NewCredentialStore("", profileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}
	return store
}

func outputJSON(data interface{}) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "JSON marshal error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(b))
}

func outputResult(data interface{}) {
	if jsonOutput {
		outputJSON(data)
		return
	}
	m := structToMap(data)
	if m == nil {
		fmt.Println(data)
		return
	}
	printTable(m)
}

func outputResultFields(data interface{}, fields []string) {
	if jsonOutput {
		outputJSON(data)
		return
	}
	m := structToMap(data)
	if m == nil {
		fmt.Println(data)
		return
	}
	fmt.Printf("%-20s %s\n", "Field", "Value")
	fmt.Printf("%-20s %s\n", strings.Repeat("━", 20), strings.Repeat("━", 60))
	for _, f := range fields {
		val := lookupField(m, f)
		s := fmtVal(val)
		fmt.Printf("%-20s %s\n", f, s)
	}
}

func lookupField(m map[string]interface{}, field string) interface{} {
	candidates := []string{field, snakeToGo(field)}
	for _, c := range candidates {
		if val, ok := m[c]; ok {
			return val
		}
	}
	for k, v := range m {
		if strings.EqualFold(k, snakeToGo(field)) {
			return v
		}
	}
	return nil
}

func snakeToGo(s string) string {
	parts := strings.Split(s, "_")
	result := ""
	for _, p := range parts {
		if len(p) > 0 {
			result += strings.ToUpper(p[:1]) + p[1:]
		}
	}
	return result
}

func printTable(m map[string]interface{}) {
	fmt.Printf("%-20s %s\n", "Field", "Value")
	fmt.Printf("%-20s %s\n", strings.Repeat("━", 20), strings.Repeat("━", 60))
	keys := sortedKeys(m)
	for _, k := range keys {
		fmt.Printf("%-20s %s\n", k, fmtVal(m[k]))
	}
}

func structToMap(v interface{}) map[string]interface{} {
	b, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil
	}
	return m
}

func fmtVal(v interface{}) string {
	switch val := v.(type) {
	case nil:
		return "(empty)"
	case string:
		return val
	case float64:
		if val == float64(int64(val)) {
			return strconv.FormatInt(int64(val), 10)
		}
		return strconv.FormatFloat(val, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	case []interface{}:
		b, _ := json.Marshal(val)
		return string(b)
	case map[string]interface{}:
		b, _ := json.MarshalIndent(val, "", "  ")
		return string(b)
	default:
		return fmt.Sprintf("%v", val)
	}
}

func sortedKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 0; i < len(keys)-1; i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	return keys
}

func parseStringList(s string) []string {
	var arr []string
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		return []string{s}
	}
	return arr
}

func parseJSONArray(s string) ([]map[string]string, error) {
	var arr []map[string]string
	if err := json.Unmarshal([]byte(s), &arr); err != nil {
		return nil, err
	}
	return arr, nil
}

func parseJSONMap(s string) (map[string]string, error) {
	var m map[string]string
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		return nil, err
	}
	return m, nil
}

func parseJSONRaw(s string) (interface{}, error) {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return nil, err
	}
	return v, nil
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}