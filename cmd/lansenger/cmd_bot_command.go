package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var botCommandCmd = &cobra.Command{
	Use:   "bot-command",
	Short: "Manage bot slash commands (4.37)",
}

// create
var botCmdCreateCmd = &cobra.Command{
	Use:   "create SCOPE_TYPE COMMANDS",
	Short: "Create bot commands",
	Args:  cobra.ExactArgs(2),
	Run:   runBotCmdCreate,
}

var botCmdCreateChatID string
var botCmdCreateChatType string
var botCmdCreateStaffID string

// query
var botCmdQueryCmd = &cobra.Command{
	Use:   "query SCOPE_TYPE",
	Short: "Query bot commands",
	Args:  cobra.ExactArgs(1),
	Run:   runBotCmdQuery,
}

var botCmdQueryChatID string
var botCmdQueryChatType string
var botCmdQueryStaffID string

// delete
var botCmdDeleteCmd = &cobra.Command{
	Use:   "delete SCOPE_TYPE",
	Short: "Delete bot commands",
	Args:  cobra.ExactArgs(1),
	Run:   runBotCmdDelete,
}

var botCmdDeleteChatID string
var botCmdDeleteChatType string
var botCmdDeleteStaffID string

func init() {
	botCmdCreateCmd.Flags().StringVar(&botCmdCreateChatID, "chat-id", "", "Group/staff openId")
	botCmdCreateCmd.Flags().StringVar(&botCmdCreateChatType, "chat-type", "", "group or staff")
	botCmdCreateCmd.Flags().StringVar(&botCmdCreateStaffID, "staff-id", "", "Staff openId")

	botCmdQueryCmd.Flags().StringVar(&botCmdQueryChatID, "chat-id", "", "Group/staff openId")
	botCmdQueryCmd.Flags().StringVar(&botCmdQueryChatType, "chat-type", "", "group or staff")
	botCmdQueryCmd.Flags().StringVar(&botCmdQueryStaffID, "staff-id", "", "Staff openId")

	botCmdDeleteCmd.Flags().StringVar(&botCmdDeleteChatID, "chat-id", "", "Group/staff openId")
	botCmdDeleteCmd.Flags().StringVar(&botCmdDeleteChatType, "chat-type", "", "group or staff")
	botCmdDeleteCmd.Flags().StringVar(&botCmdDeleteStaffID, "staff-id", "", "Staff openId")

	botCommandCmd.AddCommand(botCmdCreateCmd)
	botCommandCmd.AddCommand(botCmdQueryCmd)
	botCommandCmd.AddCommand(botCmdDeleteCmd)
	rootCmd.AddCommand(botCommandCmd)
}

func runBotCmdCreate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	scopeType, err := strconv.Atoi(args[0])
	checkError(err)

	var cmds []map[string]interface{}
	if err := json.Unmarshal([]byte(args[1]), &cmds); err != nil {
		fmt.Fprintf(cmd.ErrOrStderr(), "Error parsing COMMANDS JSON: %v\n", err)
		return
	}

	result, err := client.CreateBotCommands(ctx, scopeType, cmds, botCmdCreateChatID, botCmdCreateChatType, botCmdCreateStaffID)
	checkError(err)
	outputResult(result)
}

func runBotCmdQuery(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	scopeType, err := strconv.Atoi(args[0])
	checkError(err)

	result, err := client.FetchBotCommands(ctx, scopeType, botCmdQueryChatID, botCmdQueryChatType, botCmdQueryStaffID)
	checkError(err)
	outputResultFields(result, []string{"scope_type", "chat_id", "chat_type", "staff_id"})
	if result.Success && jsonOutput {
		b, _ := json.MarshalIndent(result.Commands, "", "  ")
		fmt.Println(string(b))
	}
}

func runBotCmdDelete(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	scopeType, err := strconv.Atoi(args[0])
	checkError(err)

	result, err := client.DeleteBotCommands(ctx, scopeType, botCmdDeleteChatID, botCmdDeleteChatType, botCmdDeleteStaffID)
	checkError(err)
	outputResult(result)
}
