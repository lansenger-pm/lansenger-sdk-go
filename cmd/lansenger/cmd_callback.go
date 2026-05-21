package main

import (
	"fmt"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var callbackCmd = &cobra.Command{
	Use:   "callback",
	Short: "Callback event commands",
}

var callbackParsePayloadCmd = &cobra.Command{
	Use:   "parse-payload ENCRYPTED_DATA",
	Short: "Parse callback payload",
	Args:  cobra.ExactArgs(1),
	Run:   runCallbackParsePayload,
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
	callbackVerifySig   bool
	callbackTimestamp   string
	callbackNonce       string
	callbackSignature   string
)

func init() {
	callbackParsePayloadCmd.Flags().StringVar(&callbackEncodingKey, "encoding-key", "", "Encoding key for decryption")
	callbackParsePayloadCmd.Flags().BoolVar(&callbackVerifySig, "verify-sig", false, "Verify signature")
	callbackParsePayloadCmd.Flags().StringVar(&callbackTimestamp, "timestamp", "", "Timestamp for signature verification")
	callbackParsePayloadCmd.Flags().StringVar(&callbackNonce, "nonce", "", "Nonce for signature verification")
	callbackParsePayloadCmd.Flags().StringVar(&callbackSignature, "signature", "", "Signature for verification")

	callbackCmd.AddCommand(callbackParsePayloadCmd)
	callbackCmd.AddCommand(callbackVerifySignatureCmd)
	callbackCmd.AddCommand(callbackEventTypesCmd)
	rootCmd.AddCommand(callbackCmd)
}

func runCallbackParsePayload(cmd *cobra.Command, args []string) {
	events, err := lansenger.ParseCallbackPayload(args[0])
	checkError(err)

	if callbackVerifySig && callbackEncodingKey != "" {
		valid := lansenger.VerifyCallbackSignature(args[0], callbackEncodingKey)
		fmt.Printf("Signature valid: %v\n\n", valid)
	}

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

func runCallbackVerifySignature(cmd *cobra.Command, args []string) {
	appSecret := args[3]
	queryString := fmt.Sprintf("timestamp=%s&nonce=%s&signature=%s", args[0], args[1], args[2])

	valid := lansenger.VerifyCallbackSignature(queryString, appSecret)
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