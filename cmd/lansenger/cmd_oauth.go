package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

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

var oauthLocalCallbackCmd = &cobra.Command{
	Use:   "local-callback",
	Short: "Start local HTTP server to capture OAuth2 callback",
	Args:  cobra.NoArgs,
	Run:   runOAuthLocalCallback,
}

var (
	oauthAuthorizeScope string
	oauthAuthorizeState string

	oauthExchangeRedirectURI string

	oauthRefreshScope string

	oauthLocalPort       int
	oauthLocalScope      string
	oauthLocalState      string
	oauthLocalNoExchange bool
	oauthLocalTimeout    int
)

func init() {
	oauthAuthorizeURLCmd.Flags().StringVarP(&oauthAuthorizeScope, "scope", "s", "basic_userinfor", "OAuth2 scope")
	oauthAuthorizeURLCmd.Flags().StringVar(&oauthAuthorizeState, "state", "", "State parameter")

	oauthExchangeCodeCmd.Flags().StringVar(&oauthExchangeRedirectURI, "redirect-uri", "", "Redirect URI")

	oauthRefreshTokenCmd.Flags().StringVarP(&oauthRefreshScope, "scope", "s", "", "OAuth2 scope")

	oauthLocalCallbackCmd.Flags().IntVarP(&oauthLocalPort, "port", "p", 8765, "Local HTTP server port")
	oauthLocalCallbackCmd.Flags().StringVarP(&oauthLocalScope, "scope", "s", "basic_userinfor", "OAuth2 scope")
	oauthLocalCallbackCmd.Flags().StringVar(&oauthLocalState, "state", "", "CSRF state (auto-generated if empty)")
	oauthLocalCallbackCmd.Flags().BoolVar(&oauthLocalNoExchange, "no-exchange", false, "Skip auto-exchange code")
	oauthLocalCallbackCmd.Flags().IntVarP(&oauthLocalTimeout, "timeout", "t", 120, "Max wait seconds")

	oauthCmd.AddCommand(oauthAuthorizeURLCmd)
	oauthCmd.AddCommand(oauthExchangeCodeCmd)
	oauthCmd.AddCommand(oauthRefreshTokenCmd)
	oauthCmd.AddCommand(oauthUserInfoCmd)
	oauthCmd.AddCommand(oauthParseCallbackCmd)
	oauthCmd.AddCommand(oauthValidateStateCmd)
	oauthCmd.AddCommand(oauthLocalCallbackCmd)
	rootCmd.AddCommand(oauthCmd)
}

func runOAuthAuthorizeURL(cmd *cobra.Command, args []string) {
	client := getClient()

	url := client.BuildAuthorizeURL(args[0], oauthAuthorizeScope, oauthAuthorizeState)
	if jsonOutput {
		outputJSON(map[string]string{"authorize_url": url})
		return
	}
	fmt.Println(url)
}

func runOAuthExchangeCode(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.ExchangeCode(ctx, args[0], oauthExchangeRedirectURI)
	checkError(err)

	store := getStore()
	store.SaveUserToken(result.UserToken, result.RefreshToken, result.ExpiresIn, result.RefreshExpiresIn)

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
	if jsonOutput {
		outputJSON(map[string]bool{"valid": valid})
		return
	}
	if valid {
		fmt.Println("valid")
	} else {
		fmt.Println("invalid")
	}
}

type callbackResult struct {
	code  string
	state string
	err   string
}

type oauthHandler struct {
	result *callbackResult
}

func (h *oauthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parsed, _ := url.Parse(r.URL.RequestURI())
	params := parsed.Query()

	code := params.Get("code")
	state := params.Get("state")
	errVal := params.Get("error")

	if errVal != "" {
		h.result = &callbackResult{err: errVal}
		w.WriteHeader(400)
		w.Write([]byte("OAuth2 error: " + errVal))
	} else if code != "" {
		h.result = &callbackResult{code: code, state: state}
		w.WriteHeader(200)
		w.Write([]byte("Authorization successful. You can close this tab."))
	} else {
		h.result = &callbackResult{err: "missing_code"}
		w.WriteHeader(400)
		w.Write([]byte("Missing code parameter."))
	}
}

func generateState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func runOAuthLocalCallback(cmd *cobra.Command, args []string) {
	client := getClient()

	state := oauthLocalState
	if state == "" {
		state = generateState()
	}

	redirectURI := fmt.Sprintf("http://localhost:%d", oauthLocalPort)
	authURL := client.BuildAuthorizeURL(redirectURI, oauthLocalScope, state)

	if jsonOutput {
		outputJSON(map[string]string{
			"authorize_url": authURL,
			"redirect_uri":  redirectURI,
			"state":         state,
		})
	} else {
		fmt.Println("Authorize URL:")
		fmt.Println(authURL)
		fmt.Printf("\nWaiting for callback on port %d... (timeout: %ds)\n", oauthLocalPort, oauthLocalTimeout)
		fmt.Println("Open the URL above in a browser, authorize, then wait.")
	}

	handler := &oauthHandler{}
	server := &http.Server{
		Addr:    fmt.Sprintf("localhost:%d", oauthLocalPort),
		Handler: handler,
	}

	go server.ListenAndServe()

	start := time.Now()
	for handler.result == nil && time.Since(start) < time.Duration(oauthLocalTimeout)*time.Second {
		time.Sleep(100 * time.Millisecond)
	}

	server.Shutdown(context.Background())

	if handler.result == nil {
		fmt.Fprintf(os.Stderr, "Error: timeout — no callback received within %ds\n", oauthLocalTimeout)
		os.Exit(1)
	}

	if handler.result.err != "" {
		fmt.Fprintf(os.Stderr, "Error: OAuth2 error: %s\n", handler.result.err)
		os.Exit(1)
	}

	code := handler.result.code
	receivedState := handler.result.state

	if jsonOutput && oauthLocalNoExchange {
		outputJSON(map[string]string{"code": code, "state": receivedState})
		return
	}

	if !jsonOutput {
		fmt.Printf("Received code: %s\n", code)
		fmt.Printf("Received state: %s\n", receivedState)
	}

	if !oauthLocalNoExchange {
		ctx := context.Background()
		result, err := client.ExchangeCode(ctx, code, redirectURI)
		checkError(err)

		store := getStore()
		store.SaveUserToken(result.UserToken, result.RefreshToken, result.ExpiresIn, result.RefreshExpiresIn)

		outputResultFields(result, []string{"user_token", "expires_in", "refresh_token", "refresh_expires_in", "staff_id", "scope"})
	} else if !jsonOutput {
		fmt.Printf("Code: %s\n", code)
		fmt.Printf("State: %s\n", receivedState)
	}
}