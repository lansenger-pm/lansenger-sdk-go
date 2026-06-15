package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"syscall"
	"time"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var oauthCmd = &cobra.Command{
	Use:   "oauth",
	Short: "OAuth2 user authentication operations",
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
	oauthAuthorizeScope   string
	oauthAuthorizeState   string
	oauthAuthorizeProfile string

	oauthExchangeRedirectURI string
	oauthExchangeProfile     string

	oauthRefreshScope   string
	oauthRefreshProfile string

	oauthUserInfoProfile string

	oauthLocalPort        int
	oauthLocalScope       string
	oauthLocalState       string
	oauthLocalNoExchange  bool
	oauthLocalTimeout     int
	oauthLocalRedirectURI string
	oauthLocalProfile     string
)

func init() {
	oauthAuthorizeURLCmd.Flags().StringVarP(&oauthAuthorizeScope, "scope", "s", "basic_userinfor", "OAuth2 scope")
	oauthAuthorizeURLCmd.Flags().StringVar(&oauthAuthorizeState, "state", "", "State parameter")
	oauthAuthorizeURLCmd.Flags().StringVarP(&oauthAuthorizeProfile, "profile", "P", "", "Credential profile (overrides global --profile)")

	oauthExchangeCodeCmd.Flags().StringVar(&oauthExchangeRedirectURI, "redirect-uri", "", "Redirect URI")
	oauthExchangeCodeCmd.Flags().StringVarP(&oauthExchangeProfile, "profile", "P", "", "Credential profile (overrides global --profile)")

	oauthRefreshTokenCmd.Flags().StringVarP(&oauthRefreshScope, "scope", "s", "", "OAuth2 scope")
	oauthRefreshTokenCmd.Flags().StringVarP(&oauthRefreshProfile, "profile", "P", "", "Credential profile (overrides global --profile)")

	oauthUserInfoCmd.Flags().StringVarP(&oauthUserInfoProfile, "profile", "P", "", "Credential profile (overrides global --profile)")

	oauthLocalCallbackCmd.Flags().IntVarP(&oauthLocalPort, "port", "p", 8765, "Local HTTP server port")
	oauthLocalCallbackCmd.Flags().StringVarP(&oauthLocalScope, "scope", "s", "basic_userinfor", "OAuth2 scope")
	oauthLocalCallbackCmd.Flags().StringVar(&oauthLocalState, "state", "", "CSRF state (auto-generated if empty)")
	oauthLocalCallbackCmd.Flags().BoolVar(&oauthLocalNoExchange, "no-exchange", false, "Skip auto-exchange code")
	oauthLocalCallbackCmd.Flags().IntVarP(&oauthLocalTimeout, "timeout", "t", 120, "Max wait seconds")
	oauthLocalCallbackCmd.Flags().StringVar(&oauthLocalRedirectURI, "redirect-uri", "", "Override redirect_uri (default: http://localhost:<port>)")
	oauthLocalCallbackCmd.Flags().StringVarP(&oauthLocalProfile, "profile", "P", "", "Credential profile (overrides global --profile)")

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
	client := getClientWithProfile(resolveProfile(oauthAuthorizeProfile))

	url := client.BuildAuthorizeURL(args[0], oauthAuthorizeScope, oauthAuthorizeState)
	if jsonOutput {
		outputJSON(map[string]string{"authorize_url": url})
		return
	}
	fmt.Println(url)
}

func runOAuthExchangeCode(cmd *cobra.Command, args []string) {
	client := getClientWithProfile(resolveProfile(oauthExchangeProfile))
	ctx := context.Background()

	result, err := client.ExchangeCode(ctx, args[0], oauthExchangeRedirectURI)
	checkError(err)

	store := getStoreWithProfile(resolveProfile(oauthExchangeProfile))
	store.SaveUserToken(result.UserToken, result.RefreshToken, result.ExpiresIn, result.RefreshExpiresIn, result.StaffID)

	outputResultFields(result, []string{"user_token", "expires_in", "refresh_token", "refresh_expires_in", "staff_id", "scope", "state"})
}

func runOAuthRefreshToken(cmd *cobra.Command, args []string) {
	client := getClientWithProfile(resolveProfile(oauthRefreshProfile))
	ctx := context.Background()

	result, err := client.RefreshUserToken(ctx, args[0], oauthRefreshScope)
	checkError(err)

	store := getStoreWithProfile(resolveProfile(oauthRefreshProfile))
	store.SaveUserToken(result.UserToken, result.RefreshToken, result.ExpiresIn, result.RefreshExpiresIn, result.StaffID)

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
	client := getClientWithProfile(resolveProfile(oauthLocalProfile))

	state := oauthLocalState
	if state == "" {
		state = generateState()
	}

	redirectURI := fmt.Sprintf("http://localhost:%d", oauthLocalPort)
	if oauthLocalRedirectURI != "" {
		redirectURI = oauthLocalRedirectURI
	}
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

	// Set SO_REUSEADDR via ListenConfig.Control for cross-platform port reuse.
	// This prevents "bind: address already in use" when the port is still
	// in TIME_WAIT from a previous run.
	lc := net.ListenConfig{
		Control: func(network, address string, c syscall.RawConn) error {
			var setErr error
			err := c.Control(func(fd uintptr) {
				setErr = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			})
			if setErr != nil {
				return setErr
			}
			return err
		},
	}
	ln, err := lc.Listen(context.Background(), "tcp", fmt.Sprintf("localhost:%d", oauthLocalPort))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: failed to listen on port %d: %v\n", oauthLocalPort, err)
		os.Exit(1)
	}

	server := &http.Server{
		Handler: handler,
	}

	go server.Serve(ln)

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

		store := getStoreWithProfile(resolveProfile(oauthLocalProfile))
		store.SaveUserToken(result.UserToken, result.RefreshToken, result.ExpiresIn, result.RefreshExpiresIn, result.StaffID)

		outputResultFields(result, []string{"user_token", "expires_in", "refresh_token", "refresh_expires_in", "staff_id", "scope"})
	} else if !jsonOutput {
		fmt.Printf("Code: %s\n", code)
		fmt.Printf("State: %s\n", receivedState)
	}
}