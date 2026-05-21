package main

import (
	"context"

	"github.com/spf13/cobra"
)

var streamingCmd = &cobra.Command{
	Use:   "streaming",
	Short: "Create and fetch streaming messages",
}

var streamingCreateCmd = &cobra.Command{
	Use:   "create RECEIVER_ID RECEIVER_TYPE STREAM_ID",
	Short: "Create a streaming message",
	Args:  cobra.ExactArgs(3),
	Run:   runStreamingCreate,
}

var streamingFetchCmd = &cobra.Command{
	Use:   "fetch MSG_ID",
	Short: "Fetch a streaming message",
	Args:  cobra.ExactArgs(1),
	Run:   runStreamingFetch,
}

func init() {
	streamingCmd.AddCommand(streamingCreateCmd)
	streamingCmd.AddCommand(streamingFetchCmd)
	rootCmd.AddCommand(streamingCmd)
}

func runStreamingCreate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.CreateStreamMessage(ctx, args[0], args[1], args[2])
	checkError(err)
	outputResultFields(result, []string{"message_id"})
}

func runStreamingFetch(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchStreamMessage(ctx, args[0])
	checkError(err)
	outputResultFields(result, []string{"message_id"})
}