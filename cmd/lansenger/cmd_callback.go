package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var callbackCmd = &cobra.Command{
	Use:   "callback",
	Short: "Callback event commands",
}

var callbackParsePayloadCmd = &cobra.Command{
	Use:   "parse-payload DATA",
	Short: "Parse callback payload (plain or encrypted)",
	Args:  cobra.ExactArgs(1),
	Run:   runCallbackParsePayload,
}

var callbackDecryptCmd = &cobra.Command{
	Use:   "decrypt-payload ENCRYPTED_DATA",
	Short: "Decrypt encrypted callback payload",
	Args:  cobra.ExactArgs(1),
	Run:   runCallbackDecrypt,
}

var callbackVerifySignatureCmd = &cobra.Command{
	Use:   "verify-signature TIMESTAMP NONCE SIGNATURE ENCODING_KEY",
	Short: "Verify callback signature",
	Args:  cobra.ExactArgs(4),
	Run:   runCallbackVerifySignature,
}

var callbackEventTypesCmd = &cobra.Command{
	Use:   "event-types",
	Short: "List all callback event types",
	Args:  cobra.NoArgs,
	Run:   runCallbackEventTypes,
}

var (
	callbackEncodingKey string
	callbackCallbackToken string
	callbackKnownAppID  string
	callbackVerifySig   bool
	callbackTimestamp   string
	callbackNonce       string
	callbackSignature   string
	callbackDataEncrypt string
)

func init() {
	callbackParsePayloadCmd.Flags().StringVar(&callbackEncodingKey, "encoding-key", "", "Encoding key for decryption")
	callbackParsePayloadCmd.Flags().StringVar(&callbackKnownAppID, "known-app-id", "", "Known app ID for org/app splitting")
	callbackParsePayloadCmd.Flags().BoolVar(&callbackVerifySig, "verify-sig", false, "Verify signature before decryption")
	callbackParsePayloadCmd.Flags().StringVar(&callbackCallbackToken, "callback-token", "", "Callback token (defaults to encoding-key)")
	callbackParsePayloadCmd.Flags().StringVar(&callbackTimestamp, "timestamp", "", "Timestamp for signature verification")
	callbackParsePayloadCmd.Flags().StringVar(&callbackNonce, "nonce", "", "Nonce for signature verification")
	callbackParsePayloadCmd.Flags().StringVar(&callbackSignature, "signature", "", "Signature for verification")
	callbackParsePayloadCmd.Flags().StringVar(&callbackDataEncrypt, "data-encrypt", "", "dataEncrypt value for signature verification")

	callbackDecryptCmd.Flags().StringVar(&callbackEncodingKey, "encoding-key", "", "Encoding key for decryption")
	callbackDecryptCmd.Flags().StringVar(&callbackKnownAppID, "known-app-id", "", "Known app ID for org/app splitting")
	callbackDecryptCmd.Flags().BoolVar(&callbackVerifySig, "verify-sig", false, "Verify signature before decryption")
	callbackDecryptCmd.Flags().StringVar(&callbackCallbackToken, "callback-token", "", "Callback token (defaults to encoding-key)")
	callbackDecryptCmd.Flags().StringVar(&callbackTimestamp, "timestamp", "", "Timestamp for signature verification")
	callbackDecryptCmd.Flags().StringVar(&callbackNonce, "nonce", "", "Nonce for signature verification")
	callbackDecryptCmd.Flags().StringVar(&callbackSignature, "signature", "", "Signature for verification")
	callbackDecryptCmd.Flags().StringVar(&callbackDataEncrypt, "data-encrypt", "", "dataEncrypt value for signature verification")

	callbackVerifySignatureCmd.Flags().StringVar(&callbackDataEncrypt, "data-encrypt", "", "dataEncrypt value")
	callbackVerifySignatureCmd.Flags().StringVar(&callbackCallbackToken, "callback-token", "", "Callback token (defaults to encoding-key)")

	callbackCmd.AddCommand(callbackParsePayloadCmd)
	callbackCmd.AddCommand(callbackDecryptCmd)
	callbackCmd.AddCommand(callbackVerifySignatureCmd)
	callbackCmd.AddCommand(callbackEventTypesCmd)
	rootCmd.AddCommand(callbackCmd)
}

func resolveEncodingKey() string {
	if callbackEncodingKey != "" {
		return callbackEncodingKey
	}
	store := getStore()
	creds, err := store.LoadCredentials()
	if err == nil && creds["encoding_key"] != "" {
		return creds["encoding_key"]
	}
	return ""
}

func resolveCallbackToken() string {
	if callbackCallbackToken != "" {
		return callbackCallbackToken
	}
	store := getStore()
	creds, err := store.LoadCredentials()
	if err == nil && creds["callback_token"] != "" {
		return creds["callback_token"]
	}
	return ""
}

func runCallbackParsePayload(cmd *cobra.Command, args []string) {
	data := args[0]
	encKey := resolveEncodingKey()

	if encKey != "" && (isEncryptedData(data) || callbackVerifySig) {
		result, err := lansenger.DecryptCallbackPayload(data, encKey, callbackKnownAppID)
		checkError(err)

		if callbackVerifySig && callbackTimestamp != "" && callbackSignature != "" {
			sigDataEncrypt := callbackDataEncrypt
			if sigDataEncrypt == "" {
				sigDataEncrypt = data
			}
			token := resolveCallbackToken()
			valid := lansenger.VerifyCallbackSignature(callbackTimestamp, callbackNonce, callbackSignature, encKey, sigDataEncrypt, token)
			fmt.Printf("Signature valid: %v\n\n", valid)
		}

		fmt.Printf("OrgID:  %s\n", result.OrgID)
		fmt.Printf("AppID:  %s\n", result.AppID)
		fmt.Printf("Events: %d\n\n", len(result.Events))

		for i, event := range result.Events {
			fmt.Printf("Event %d:\n", i+1)
			fmt.Printf("  Type:     %s\n", event.EventType)
			fmt.Printf("  Category: %s\n", event.Category)
			if event.EventID != "" {
				fmt.Printf("  EventID:  %s\n", event.EventID)
			}
			if event.AppID != "" {
				fmt.Printf("  AppID:    %s\n", event.AppID)
			}
			if event.OrgID != "" {
				fmt.Printf("  OrgID:    %s\n", event.OrgID)
			}
			if len(event.Data) > 0 {
				fmt.Printf("  Data:\n")
				for k, v := range event.Data {
					fmt.Printf("    %s: %v\n", k, v)
				}
			}
			fmt.Println()
		}
		return
	}

	events, err := lansenger.ParseCallbackPayload(data)
	checkError(err)

	for i, event := range events {
		fmt.Printf("Event %d:\n", i+1)
		fmt.Printf("  Type:     %s\n", event.EventType)
		fmt.Printf("  Category: %s\n", event.Category)
		if len(event.Data) > 0 {
			fmt.Printf("  Data:\n")
			for k, v := range event.Data {
				fmt.Printf("    %s: %v\n", k, v)
			}
		}
		fmt.Println()
	}
}

func runCallbackDecrypt(cmd *cobra.Command, args []string) {
	encKey := resolveEncodingKey()
	if encKey == "" {
		fmt.Fprintf(os.Stderr, "Error: encoding-key required for decryption. Use --encoding-key or set encoding_key in config.\n")
		os.Exit(1)
	}

	if callbackVerifySig && callbackTimestamp != "" && callbackSignature != "" {
		sigDataEncrypt := callbackDataEncrypt
		if sigDataEncrypt == "" {
			sigDataEncrypt = args[0]
		}
		token := resolveCallbackToken()
		valid := lansenger.VerifyCallbackSignature(callbackTimestamp, callbackNonce, callbackSignature, encKey, sigDataEncrypt, token)
		if !valid {
			fmt.Fprintf(os.Stderr, "Error: callback signature verification failed\n")
			os.Exit(1)
		}
		fmt.Printf("Signature valid: true\n\n")
	}

	result, err := lansenger.DecryptCallbackPayload(args[0], encKey, callbackKnownAppID)
	checkError(err)

	if jsonOutput {
		outputJSON(result)
		return
	}

	fmt.Printf("%-20s %s\n", "Field", "Value")
	fmt.Printf("%-20s %s\n", strings.Repeat("━", 20), strings.Repeat("━", 60))
	fmt.Printf("%-20s %s\n", "orgId", result.OrgID)
	fmt.Printf("%-20s %s\n", "appId", result.AppID)
	fmt.Printf("%-20s %d\n", "events_count", len(result.Events))
	fmt.Printf("%-20s %d\n", "length", result.Length)
	fmt.Println()

	for i, event := range result.Events {
		fmt.Printf("Event %d:\n", i+1)
		fmt.Printf("  Type:     %s\n", event.EventType)
		fmt.Printf("  Category: %s\n", event.Category)
		if event.EventID != "" {
			fmt.Printf("  EventID:  %s\n", event.EventID)
		}
		if event.AppID != "" {
			fmt.Printf("  AppID:    %s\n", event.AppID)
		}
		if event.OrgID != "" {
			fmt.Printf("  OrgID:    %s\n", event.OrgID)
		}
		if len(event.Data) > 0 {
			fmt.Printf("  Data:\n")
			for k, v := range event.Data {
				fmt.Printf("    %s: %v\n", k, v)
			}
		}
		fmt.Println()
	}
}

func runCallbackVerifySignature(cmd *cobra.Command, args []string) {
	timestamp := args[0]
	nonce := args[1]
	signature := args[2]
	encodingKey := args[3]
	dataEncrypt := callbackDataEncrypt
	token := resolveCallbackToken()

	valid := lansenger.VerifyCallbackSignature(timestamp, nonce, signature, encodingKey, dataEncrypt, token)
	if valid {
		fmt.Println("valid")
	} else {
		fmt.Println("invalid")
	}
}

func runCallbackEventTypes(cmd *cobra.Command, args []string) {
	eventTypes := lansenger.GetCallbackEventTypes()

	if jsonOutput {
		outputJSON(eventTypes)
		return
	}

	fmt.Printf("%-30s %s\n", "Event Type", "Category")
	fmt.Printf("%-30s %s\n", "------------------------------", "--------------------")
	for eventType, category := range eventTypes {
		fmt.Printf("%-30s %s\n", eventType, category)
	}
}

func isEncryptedData(data string) bool {
	data = strings.TrimSpace(data)
	if strings.HasPrefix(data, "{") {
		var wrapper map[string]interface{}
		if json.Unmarshal([]byte(data), &wrapper) == nil {
			if _, ok := wrapper["dataEncrypt"].(string); ok {
				return true
			}
		}
		return false
	}
	if strings.Contains(data, "eventType=") {
		return false
	}
	return true
}