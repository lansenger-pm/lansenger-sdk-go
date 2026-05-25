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

var (
	configSetProfile   string
	configShowProfile  string
	configClearProfile string
	configClearAll     bool
)

func init() {
	configSetCmd.Flags().StringVarP(&configSetProfile, "profile", "P", "", "Profile to set config for (overrides global --profile)")
	configShowCmd.Flags().StringVarP(&configShowProfile, "profile", "P", "", "Profile to show config for (overrides global --profile)")
	configClearCmd.Flags().StringVarP(&configClearProfile, "profile", "P", "", "Profile to clear (overrides global --profile)")
	configClearCmd.Flags().BoolVar(&configClearAll, "all", false, "Clear all profiles and delete state file")

	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configClearCmd)
	configCmd.AddCommand(configListProfilesCmd)
	rootCmd.AddCommand(configCmd)
}

func resolveProfile(localFlag string) string {
	if localFlag != "" {
		return localFlag
	}
	return profileName
}

func resolveStorePath() string {
	store := lansenger.NewCredentialStore("", "")
	return store.Path()
}

func maskSecret(val string) string {
	if val == "" {
		return "(empty)"
	}
	return "***"
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
	}
	if !validKeys[key] {
		fmt.Fprintf(os.Stderr, "Error: Invalid config key '%s'. Valid keys: app_id, app_secret, api_gateway_url, passport_url, encoding_key, callback_token\n", key)
		os.Exit(1)
	}

	store := lansenger.NewCredentialStore("", prof)
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
		err = store.SaveCredentials(creds["app_id"], creds["app_secret"], creds["api_gateway_url"], creds["passport_url"])
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
	store := lansenger.NewCredentialStore("", prof)

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
		store := lansenger.NewCredentialStore("", "")
		err := store.Clear()
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
	store := lansenger.NewCredentialStore("", prof)
	err := store.ClearProfile()
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

func runConfigListProfiles(cmd *cobra.Command, args []string) {
	metaStore := lansenger.NewCredentialStore("", "")
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
			pStore := lansenger.NewCredentialStore("", p)
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
		pStore := lansenger.NewCredentialStore("", p)
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