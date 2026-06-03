package main

import (
	"context"
	"fmt"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var oauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "OAuth2 authentication commands",
}

var oauthAuthorizeURLCmd = &cobra.Command{
	Use:   "authorize-url REDIRECT_URI",
	Short: "Build OAuth2 authorize URL",
	Args:  cobra.ExactArgs(1),
	Run:   runOAuthAuthorizeURL,
}

var oauthExchangeCodeCmd = &cobra.Command{
	Use:   "exchange-code CODE",
	Short: "Exchange authorization code for user token",
	Args:  cobra.ExactArgs(1),
	Run:   runOAuthExchangeCode,
}

var oauthRefreshTokenCmd = &cobra.Command{
	Use:   "refresh-token REFRESH_TOKEN",
	Short: "Refresh a user token",
	Args:  cobra.ExactArgs(1),
	Run:   runOAuthRefreshToken,
}

var oauthUserInfoCmd = &cobra.Command{
	Use:   "user-info USER_TOKEN",
	Short: "Fetch user info",
	Args:  cobra.ExactArgs(1),
	Run:   runOAuthUserInfo,
}

var oauthParseCallbackCmd = &cobra.Command{
	Use:   "parse-callback QUERY_STRING",
	Short: "Parse OAuth2 callback query string",
	Args:  cobra.ExactArgs(1),
	Run:   runOAuthParseCallback,
}

var oauthValidateStateCmd = &cobra.Command{
	Use:   "validate-state CALLBACK_STATE EXPECTED_STATE",
	Short: "Validate callback state parameter",
	Args:  cobra.ExactArgs(2),
	Run:   runOAuthValidateState,
}

var (
	oauthAuthorizeScope string
	oauthAuthorizeState string

	oauthExchangeRedirectURI string

	oauthRefreshScope string
)

func init() {
	oauthAuthorizeURLCmd.Flags().StringVarP(&oauthAuthorizeScope, "scope", "s", "basic_userinfor", "OAuth2 scope")
	oauthAuthorizeURLCmd.Flags().StringVar(&oauthAuthorizeState, "state", "", "State parameter")

	oauthExchangeCodeCmd.Flags().StringVar(&oauthExchangeRedirectURI, "redirect-uri", "", "Redirect URI")

	oauthRefreshTokenCmd.Flags().StringVarP(&oauthRefreshScope, "scope", "s", "", "OAuth2 scope")

	oauthCmd.AddCommand(oauthAuthorizeURLCmd)
	oauthCmd.AddCommand(oauthExchangeCodeCmd)
	oauthCmd.AddCommand(oauthRefreshTokenCmd)
	oauthCmd.AddCommand(oauthUserInfoCmd)
	oauthCmd.AddCommand(oauthParseCallbackCmd)
	oauthCmd.AddCommand(oauthValidateStateCmd)
	rootCmd.AddCommand(oauthCmd)
}

func runOAuthAuthorizeURL(cmd *cobra.Command, args []string) {
	client := getClient()

	url := client.BuildAuthorizeURL(args[0], oauthAuthorizeScope, oauthAuthorizeState)
	fmt.Println(url)
}

func runOAuthExchangeCode(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.ExchangeCode(ctx, args[0], oauthExchangeRedirectURI)
	checkError(err)
	outputResultFields(result, []string{"user_token", "expires_in", "refresh_token", "refresh_expires_in", "staff_id", "scope", "state"})
}

func runOAuthRefreshToken(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.RefreshUserToken(ctx, args[0], oauthRefreshScope)
	checkError(err)
	outputResultFields(result, []string{"user_token", "expires_in", "refresh_token", "refresh_expires_in", "staff_id", "scope"})
}

func runOAuthUserInfo(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchUserInfo(ctx, args[0])
	checkError(err)
	outputResult(result)
}

func runOAuthParseCallback(cmd *cobra.Command, args []string) {
	result, err := lansenger.ParseAuthorizeCallback(args[0])
	checkError(err)
	outputResult(result)
}

func runOAuthValidateState(cmd *cobra.Command, args []string) {
	valid := lansenger.ValidateCallbackState(args[0], args[1])
	if valid {
		fmt.Println("valid")
	} else {
		fmt.Println("invalid")
	}
}