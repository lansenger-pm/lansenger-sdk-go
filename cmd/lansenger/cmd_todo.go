package main

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Todo task commands",
}

var todoCreateCmd = &cobra.Command{
	Use:   "create TITLE LINK PC_LINK EXECUTOR_IDS ORG_ID",
	Short: "Create a todo task",
	Args:  cobra.ExactArgs(5),
	Run:   runTodoCreate,
}

var todoUpdateCmd = &cobra.Command{
	Use:   "update TODOTASK_ID TITLE LINK PC_LINK ORG_ID",
	Short: "Update a todo task",
	Args:  cobra.ExactArgs(5),
	Run:   runTodoUpdate,
}

var todoUpdateStatusCmd = &cobra.Command{
	Use:   "update-status TODOTASK_ID STATUS ORG_ID",
	Short: "Update todo task status",
	Args:  cobra.ExactArgs(3),
	Run:   runTodoUpdateStatus,
}

var todoDeleteCmd = &cobra.Command{
	Use:   "delete TODOTASK_ID ORG_ID",
	Short: "Delete a todo task",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoDelete,
}

var todoListCmd = &cobra.Command{
	Use:   "list ORG_ID",
	Short: "List todo tasks",
	Args:  cobra.ExactArgs(1),
	Run:   runTodoList,
}

var todoFetchBySourceCmd = &cobra.Command{
	Use:   "fetch-by-source SOURCE_ID ORG_ID",
	Short: "Fetch todo task by source ID",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoFetchBySource,
}

var todoFetchByIDCmd = &cobra.Command{
	Use:   "fetch-by-id TODOTASK_ID ORG_ID",
	Short: "Fetch todo task by ID",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoFetchByID,
}

var todoStatusCountsCmd = &cobra.Command{
	Use:   "status-counts STAFF_ID ORG_ID",
	Short: "Fetch todo task status counts",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoStatusCounts,
}

var todoExecutorStatusCmd = &cobra.Command{
	Use:   "executor-status EXECUTOR_STATUS_LIST ORG_ID",
	Short: "Update executor status",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoExecutorStatus,
}

var todoAddExecutorsCmd = &cobra.Command{
	Use:   "add-executors EXECUTOR_IDS ORG_ID",
	Short: "Add executors to a todo task",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoAddExecutors,
}

var todoDeleteExecutorsCmd = &cobra.Command{
	Use:   "delete-executors EXECUTOR_IDS ORG_ID",
	Short: "Delete executors from a todo task",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoDeleteExecutors,
}

var todoExecutorListCmd = &cobra.Command{
	Use:   "executor-list TODOTASK_ID ORG_ID",
	Short: "Fetch executor list for a todo task",
	Args:  cobra.ExactArgs(2),
	Run:   runTodoExecutorList,
}

var (
	todoCreateType       int
	todoCreateSourceID   string
	todoCreateDesc       string
	todoCreateSenderID   string
	todoCreateUserToken  string

	todoUpdateDesc      string
	todoUpdateUserToken string

	todoUpdateStatusStaffID   string
	todoUpdateStatusUserToken string

	todoDeleteStaffID   string
	todoDeleteUserToken string

	todoListUserToken string
	todoListAppIDs    string
	todoListStaffID   string
	todoListStatus    string

	todoFetchBySourceStaffID   string
	todoFetchBySourceUserToken string

	todoFetchByIDStaffID   string
	todoFetchByIDUserToken string

	todoStatusCountsAppID    string
	todoStatusCountsStatus   string
	todoStatusCountsUserToken string

	todoExecutorStatusTaskID   string
	todoExecutorStatusUserToken string

	todoAddExecutorsTaskID   string
	todoAddExecutorsUserToken string

	todoDelExecutorsTaskID   string
	todoDelExecutorsUserToken string

	todoExecutorListStaffID   string
	todoExecutorListStatus    string
	todoExecutorListUserToken string
)

func init() {
	todoCreateCmd.Flags().IntVarP(&todoCreateType, "type", "t", 1, "Todo type (1=notification)")
	todoCreateCmd.Flags().StringVar(&todoCreateSourceID, "source-id", "", "Source ID")
	todoCreateCmd.Flags().StringVarP(&todoCreateDesc, "desc", "d", "", "Description")
	todoCreateCmd.Flags().StringVar(&todoCreateSenderID, "sender-id", "", "Sender ID")
	todoCreateCmd.Flags().StringVar(&todoCreateUserToken, "user-token", "", "User token")

	todoUpdateCmd.Flags().StringVarP(&todoUpdateDesc, "desc", "d", "", "Description")
	todoUpdateCmd.Flags().StringVar(&todoUpdateUserToken, "user-token", "", "User token")

	todoUpdateStatusCmd.Flags().StringVar(&todoUpdateStatusStaffID, "staff-id", "", "Staff ID")
	todoUpdateStatusCmd.Flags().StringVar(&todoUpdateStatusUserToken, "user-token", "", "User token")

	todoDeleteCmd.Flags().StringVar(&todoDeleteStaffID, "staff-id", "", "Staff ID")
	todoDeleteCmd.Flags().StringVar(&todoDeleteUserToken, "user-token", "", "User token")

	todoListCmd.Flags().StringVar(&todoListUserToken, "user-token", "", "User token")
	todoListCmd.Flags().StringVar(&todoListAppIDs, "app-ids", "", "App IDs (comma-separated)")
	todoListCmd.Flags().StringVar(&todoListStaffID, "staff-id", "", "Staff ID")
	todoListCmd.Flags().StringVar(&todoListStatus, "status", "", "Status filter (comma-separated)")

	todoFetchBySourceCmd.Flags().StringVar(&todoFetchBySourceStaffID, "staff-id", "", "Staff ID")
	todoFetchBySourceCmd.Flags().StringVar(&todoFetchBySourceUserToken, "user-token", "", "User token")

	todoFetchByIDCmd.Flags().StringVar(&todoFetchByIDStaffID, "staff-id", "", "Staff ID")
	todoFetchByIDCmd.Flags().StringVar(&todoFetchByIDUserToken, "user-token", "", "User token")

	todoStatusCountsCmd.Flags().StringVar(&todoStatusCountsAppID, "app-id", "", "App ID")
	todoStatusCountsCmd.Flags().StringVar(&todoStatusCountsStatus, "status", "", "Status filter")
	todoStatusCountsCmd.Flags().StringVar(&todoStatusCountsUserToken, "user-token", "", "User token")

	todoExecutorStatusCmd.Flags().StringVar(&todoExecutorStatusTaskID, "task-id", "", "Todo task ID")
	todoExecutorStatusCmd.Flags().StringVar(&todoExecutorStatusUserToken, "user-token", "", "User token")

	todoAddExecutorsCmd.Flags().StringVar(&todoAddExecutorsTaskID, "task-id", "", "Todo task ID")
	todoAddExecutorsCmd.Flags().StringVar(&todoAddExecutorsUserToken, "user-token", "", "User token")

	todoDeleteExecutorsCmd.Flags().StringVar(&todoDelExecutorsTaskID, "task-id", "", "Todo task ID")
	todoDeleteExecutorsCmd.Flags().StringVar(&todoDelExecutorsUserToken, "user-token", "", "User token")

	todoExecutorListCmd.Flags().StringVar(&todoExecutorListStaffID, "staff-id", "", "Staff ID")
	todoExecutorListCmd.Flags().StringVar(&todoExecutorListStatus, "status", "", "Status filter (comma-separated)")
	todoExecutorListCmd.Flags().StringVar(&todoExecutorListUserToken, "user-token", "", "User token")

	todoCmd.AddCommand(todoCreateCmd)
	todoCmd.AddCommand(todoUpdateCmd)
	todoCmd.AddCommand(todoUpdateStatusCmd)
	todoCmd.AddCommand(todoDeleteCmd)
	todoCmd.AddCommand(todoListCmd)
	todoCmd.AddCommand(todoFetchBySourceCmd)
	todoCmd.AddCommand(todoFetchByIDCmd)
	todoCmd.AddCommand(todoStatusCountsCmd)
	todoCmd.AddCommand(todoExecutorStatusCmd)
	todoCmd.AddCommand(todoAddExecutorsCmd)
	todoCmd.AddCommand(todoDeleteExecutorsCmd)
	todoCmd.AddCommand(todoExecutorListCmd)
	rootCmd.AddCommand(todoCmd)
}

func splitCommaList(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, ",")
}

func runTodoCreate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	orgID := args[4]
	executorIDs := splitCommaList(args[3])

	result, err := client.CreateTodoTask(ctx, args[0], todoCreateType, args[1], args[2], executorIDs, orgID, todoCreateSourceID, todoCreateDesc, todoCreateSenderID, todoCreateUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoUpdate(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UpdateTodoTask(ctx, args[0], args[1], args[2], args[3], args[4], todoUpdateDesc, todoUpdateUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoUpdateStatus(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.UpdateTodoTaskStatus(ctx, args[0], args[1], args[2], todoUpdateStatusStaffID, todoUpdateStatusUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoDelete(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.DeleteTodoTask(ctx, args[0], args[1], todoDeleteStaffID, todoDeleteUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoList(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var appIDs []string
	if todoListAppIDs != "" {
		appIDs = splitCommaList(todoListAppIDs)
	}

	var statusList []string
	if todoListStatus != "" {
		statusList = splitCommaList(todoListStatus)
	}

	result, err := client.FetchTodoTaskList(ctx, args[0], appIDs, todoListStaffID, statusList, todoListUserToken)
	checkError(err)
	outputResultFields(result, []string{"total"})
}

func runTodoFetchBySource(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchTodoTaskBySourceID(ctx, args[0], args[1], todoFetchBySourceStaffID, todoFetchBySourceUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id", "source_id", "title", "desc", "status", "link", "pc_link", "sender_id", "create_time", "app_id"})
}

func runTodoFetchByID(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchTodoTaskByID(ctx, args[0], args[1], todoFetchByIDStaffID, todoFetchByIDUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id", "source_id", "title", "desc", "status", "link", "pc_link", "sender_id", "create_time", "app_id"})
}

func runTodoStatusCounts(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	result, err := client.FetchTodoTaskStatusCounts(ctx, args[0], args[1], todoStatusCountsAppID, todoStatusCountsStatus, todoStatusCountsUserToken)
	checkError(err)
	outputResult(result)
}

func runTodoExecutorStatus(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	raw, err := parseJSONRaw(args[0])
	checkError(err)

	var executorStatusList []map[string]interface{}
	if arr, ok := raw.([]interface{}); ok {
		for _, item := range arr {
			if m, ok := item.(map[string]interface{}); ok {
				executorStatusList = append(executorStatusList, m)
			}
		}
	} else if m, ok := raw.(map[string]interface{}); ok {
		executorStatusList = []map[string]interface{}{m}
	}

	result, err := client.UpdateExecutorStatus(ctx, executorStatusList, args[1], todoExecutorStatusTaskID, todoExecutorStatusUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoAddExecutors(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	executorIDs := splitCommaList(args[0])

	result, err := client.AddExecutors(ctx, executorIDs, args[1], todoAddExecutorsTaskID, todoAddExecutorsUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoDeleteExecutors(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	executorIDs := splitCommaList(args[0])

	result, err := client.DeleteExecutors(ctx, executorIDs, args[1], todoDelExecutorsTaskID, todoDelExecutorsUserToken)
	checkError(err)
	outputResultFields(result, []string{"todotask_id"})
}

func runTodoExecutorList(cmd *cobra.Command, args []string) {
	client := getClient()
	ctx := context.Background()

	var statusList []string
	if todoExecutorListStatus != "" {
		statusList = splitCommaList(todoExecutorListStatus)
	}

	result, err := client.FetchExecutorList(ctx, args[0], args[1], todoExecutorListStaffID, statusList, todoExecutorListUserToken)
	checkError(err)
	outputResultFields(result, []string{"total"})
}