package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Lansenger CLI configuration and credential profiles",
}

var configSetCmd = &cobra.Command{
	Use:   "set KEY VALUE",
	Short: "Set a configuration value for the current profile",
	Args:  cobra.ExactArgs(2),
	Run:   runConfigSet,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current profile configuration (secrets masked)",
	Args:  cobra.NoArgs,
	Run:   runConfigShow,
}

var configClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the current profile or all profiles",
	Args:  cobra.NoArgs,
	Run:   runConfigClear,
}

var configListProfilesCmd = &cobra.Command{
	Use:   "list-profiles",
	Short: "List all credential profiles with active marker",
	Args:  cobra.NoArgs,
	Run:   runConfigListProfiles,
}

var configDeleteProfileCmd = &cobra.Command{
	Use:   "delete-profile <profile-name>",
	Short: "Delete a credential profile by name",
	Args:  cobra.ExactArgs(1),
	Run:   runConfigDeleteProfile,
}

var configListUsersCmd = &cobra.Command{
	Use:   "list-users",
	Short: "List all users with stored user tokens in the current profile",
	Args:  cobra.NoArgs,
	Run:   runConfigListUsers,
}

var (
	configSetProfile       string
	configShowProfile      string
	configClearProfile     string
	configClearAll         bool
	configListUsersProfile string
	configListUsersShowTokens bool
)

func init() {
	configSetCmd.Flags().StringVarP(&configSetProfile, "profile", "P", "", "Profile to set config for (overrides global --profile)")
	configShowCmd.Flags().StringVarP(&configShowProfile, "profile", "P", "", "Profile to show config for (overrides global --profile)")
	configClearCmd.Flags().StringVarP(&configClearProfile, "profile", "P", "", "Profile to clear (overrides global --profile)")
	configClearCmd.Flags().BoolVar(&configClearAll, "all", false, "Clear all profiles and delete state file")
	configListUsersCmd.Flags().StringVarP(&configListUsersProfile, "profile", "P", "", "Profile to list users for (overrides global --profile)")
	configListUsersCmd.Flags().BoolVarP(&configListUsersShowTokens, "show-tokens", "T", false, "Show user tokens (security warning)")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configClearCmd)
	configCmd.AddCommand(configDeleteProfileCmd)
	configCmd.AddCommand(configListProfilesCmd)
	configCmd.AddCommand(configListUsersCmd)
	rootCmd.AddCommand(configCmd)
}

func resolveProfile(localFlag string) string {
	if localFlag != "" {
		return localFlag
	}
	return profileName
}

func resolveStorePath() string {
	store, err := lansenger.NewCredentialStore("", "")
	if err != nil {
		return ""
	}
	return store.Path()
}

func maskSecret(val string) string {
	if val == "" {
		return "(empty)"
	}
	return "***"
}

func displayVal(val string) string {
	if val == "" {
		return "(empty)"
	}
	return val
}

func runConfigSet(cmd *cobra.Command, args []string) {
	key := args[0]
	value := args[1]
	prof := resolveProfile(configSetProfile)

	validKeys := map[string]bool{
		"app_id":          true,
		"app_secret":      true,
		"api_gateway_url": true,
		"passport_url":    true,
		"encoding_key":    true,
		"callback_token":  true,
		"redirect_uri":    true,
	}
	if !validKeys[key] {
		fmt.Fprintf(os.Stderr, "Error: Invalid config key '%s'. Valid keys: app_id, app_secret, api_gateway_url, passport_url, redirect_uri, encoding_key, callback_token\n", key)
		os.Exit(1)
	}

	store, err := lansenger.NewCredentialStore("", prof)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}
	creds, err := store.LoadCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading credentials: %v\n", err)
		os.Exit(1)
	}

	creds[key] = value

	if key == "encoding_key" || key == "callback_token" {
		encKey := creds["encoding_key"]
		cbToken := creds["callback_token"]
		err = store.SaveCallbackConfig(encKey, cbToken)
	} else {
		err = store.SaveCredentials(creds["app_id"], creds["app_secret"], creds["api_gateway_url"], creds["passport_url"], creds["redirect_uri"])
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error saving credentials: %v\n", err)
		os.Exit(1)
	}

	result := map[string]interface{}{
		"profile": prof,
		"key":     key,
		"value":   maskIfSecret(key, value),
		"status":  "set",
	}
	outputResult(result)
}

func maskIfSecret(key, value string) string {
	if key == "app_id" || key == "app_secret" || key == "encoding_key" || key == "callback_token" {
		return maskSecret(value)
	}
	return value
}

func runConfigShow(cmd *cobra.Command, args []string) {
	prof := resolveProfile(configShowProfile)
	store, err := lansenger.NewCredentialStore("", prof)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}

	creds, err := store.LoadCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading credentials: %v\n", err)
		os.Exit(1)
	}

	hasCreds := store.HasCredentials()
	storePath := resolveStorePath()

	result := map[string]interface{}{
		"profile":            prof,
		"has_credentials":    hasCreds,
		"app_id":             maskSecret(creds["app_id"]),
		"app_secret":         maskSecret(creds["app_secret"]),
		"api_gateway_url":    creds["api_gateway_url"],
		"passport_url":       creds["passport_url"],
		"encoding_key":       maskSecret(creds["encoding_key"]),
		"callback_token":     maskSecret(creds["callback_token"]),
		"store_path":         storePath,
	}

	if !jsonOutput {
		fmt.Printf("%-20s %s\n", "Field", "Value")
		fmt.Printf("%-20s %s\n", strings.Repeat("━", 20), strings.Repeat("━", 60))
		keys := []string{"profile", "has_credentials", "app_id", "app_secret", "api_gateway_url", "passport_url", "encoding_key", "callback_token", "store_path"}
		for _, k := range keys {
			v := result[k]
			fmt.Printf("%-20s %s\n", k, fmtVal(v))
		}
		return
	}

	outputResult(result)
}

func runConfigClear(cmd *cobra.Command, args []string) {
	if configClearAll {
		store, err := lansenger.NewCredentialStore("", "")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
			os.Exit(1)
		}
		err = store.Clear()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error clearing all profiles: %v\n", err)
			os.Exit(1)
		}
		result := map[string]interface{}{
			"status":  "cleared_all",
			"message": "All profiles and state file removed",
		}
		outputResult(result)
		return
	}

	prof := resolveProfile(configClearProfile)
	store, err := lansenger.NewCredentialStore("", prof)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}
	err = store.ClearProfile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error clearing profile '%s': %v\n", prof, err)
		os.Exit(1)
	}

	result := map[string]interface{}{
		"profile": prof,
		"status":  "cleared",
		"message": fmt.Sprintf("Profile '%s' cleared", prof),
	}
	outputResult(result)
}

func runConfigDeleteProfile(cmd *cobra.Command, args []string) {
	name := args[0]
	store, err := lansenger.NewCredentialStore("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}
	if !store.DeleteProfileByName(name) {
		fmt.Fprintf(os.Stderr, "Error: profile '%s' not found\n", name)
		os.Exit(1)
	}
	activeProfile := store.GetActiveProfile()
	result := map[string]interface{}{
		"profile":  name,
		"status":   "deleted",
		"message":  fmt.Sprintf("Profile '%s' deleted", name),
		"active":   activeProfile,
	}
	outputResult(result)
}

func runConfigListProfiles(cmd *cobra.Command, args []string) {
	metaStore, err := lansenger.NewCredentialStore("", "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}
	profiles, err := metaStore.ListProfiles()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing profiles: %v\n", err)
		os.Exit(1)
	}

	sort.Strings(profiles)

	activeProfile := metaStore.GetActiveProfile()

	if !jsonOutput {
		if len(profiles) == 0 {
			fmt.Println("No profiles found.")
			return
		}
		fmt.Printf("%-4s  %-20s  %-8s  %-50s  %-50s\n", "Act", "Profile", "Creds", "App ID", "API Gateway URL")
		fmt.Printf("%-4s  %-20s  %-8s  %-50s  %-50s\n", strings.Repeat("━", 4), strings.Repeat("━", 20), strings.Repeat("━", 8), strings.Repeat("━", 50), strings.Repeat("━", 50))
		for _, p := range profiles {
			pStore, storeErr := lansenger.NewCredentialStore("", p)
			if storeErr != nil {
				continue
			}
			creds, _ := pStore.LoadCredentials()
			active := ""
			if p == activeProfile {
				active = "✓"
			}
			hasCreds := "✗"
			if pStore.HasCredentials() {
				hasCreds = "✓"
			}
			appID := creds["app_id"]
			gwURL := creds["api_gateway_url"]
			fmt.Printf("%-4s  %-20s  %-8s  %-50s  %-50s\n", active, p, hasCreds, appID, gwURL)
		}
		return
	}

	items := make([]map[string]interface{}, 0, len(profiles))
	for _, p := range profiles {
		pStore, storeErr := lansenger.NewCredentialStore("", p)
		if storeErr != nil {
			continue
		}
		creds, _ := pStore.LoadCredentials()
		isActive := p == activeProfile
		hasCreds := pStore.HasCredentials()
		items = append(items, map[string]interface{}{
			"profile":          p,
			"active":           isActive,
			"has_credentials":  hasCreds,
			"app_id":           creds["app_id"],
			"api_gateway_url":  creds["api_gateway_url"],
		})
	}
	outputResult(items)
}

func runConfigListUsers(cmd *cobra.Command, args []string) {
	prof := resolveProfile(configListUsersProfile)
	store, err := lansenger.NewCredentialStore("", prof)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating credential store: %v\n", err)
		os.Exit(1)
	}

	users, err := store.ListUserTokens()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error listing users: %v\n", err)
		os.Exit(1)
	}

	sort.Strings(users)

	if !jsonOutput {
		if len(users) == 0 {
			fmt.Printf("No users found in profile '%s'.\n", prof)
			return
		}
		fmt.Printf("Users in profile '%s':\n", prof)
		for i, staffID := range users {
			fmt.Printf("  %d. %s\n", i+1, staffID)
			if configListUsersShowTokens {
				tokenData, err := store.LoadUserToken(staffID)
				if err != nil {
					fmt.Printf("     (error loading token)\n")
					continue
				}
				fmt.Printf("     user_token:          %s\n", displayVal(tokenData["user_token"]))
			fmt.Printf("     refresh_token:       %s\n", displayVal(tokenData["refresh_token"]))
			fmt.Printf("     expires_in:          %s\n", displayVal(tokenData["user_token_expiry"]))
			fmt.Printf("     refresh_expires_in:  %s\n", displayVal(tokenData["refresh_token_expiry"]))
			}
		}
		if !configListUsersShowTokens {
			fmt.Println("Hint: Use --show-tokens to view user tokens")
		}
		return
	}

	result := map[string]interface{}{
		"profile": prof,
		"users":   users,
	}

	if configListUsersShowTokens {
		tokens := make(map[string]map[string]string)
		for _, staffID := range users {
			tokenData, err := store.LoadUserToken(staffID)
			if err != nil {
				tokens[staffID] = map[string]string{}
			} else {
				tokens[staffID] = map[string]string{
					"user_token":         tokenData["user_token"],
					"refresh_token":      tokenData["refresh_token"],
					"expires_in":         tokenData["user_token_expiry"],
					"refresh_expires_in": tokenData["refresh_token_expiry"],
				}
			}
		}
		result["tokens"] = tokens
	}

	outputResult(result)
}