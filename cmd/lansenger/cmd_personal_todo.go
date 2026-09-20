package main

import (
	"context"
	"encoding/base64"
	"os"
	"strconv"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"
	"github.com/spf13/cobra"
)

var personalTodoCmd = &cobra.Command{
	Use:   "personal-todo",
	Short: "Manage user-owned personal todos (个人待办)",
}

var (
	pTodoStartTime       int64
	pTodoDueTime         int64
	pTodoFinishTime      int64
	pTodoPriority        int
	pTodoDescription     string
	pTodoParentCode      string
	pTodoStatusTagNo     string
	pTodoStatusTagYes    string
	pTodoSubscribeStatus int
	pTodoUserCode        string
	pTodoExecutors       string
	pTodoCopys           string
	pTodoResources       string
	pTodoReminds         string
	pTodoUserToken       string
	pTodoUpdateFields    string
	pTodoSubject         string
	pTodoAppID           string
	pTodoPage            int
	pTodoPageSize        int
	pTodoStatus          int
	pTodoAppFilter       string
	pTodoCategoryName    string
	pTodoFile            string
	pTodoThumb           bool
	pTodoExtensionInfo   string
	pTodoFileName        string
	pTodoMD5             string
	pTodoSize            int64
)

var personalTodoSaveCmd = &cobra.Command{
	Use:   "save SUBJECT CREATE_USER_ID ORG_ID APPID",
	Short: "Create a personal todo",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().SavePersonalTodo(context.Background(), &lansenger.PersonalTodoSaveParams{
			Subject:         args[0],
			CreateUserID:    args[1],
			OrgID:           args[2],
			AppID:           args[3],
			StartTime:       pTodoStartTime,
			DueTime:         pTodoDueTime,
			Priority:        pTodoPriority,
			Description:     pTodoDescription,
			ParentCode:      pTodoParentCode,
			StatusTagNo:     pTodoStatusTagNo,
			StatusTagYes:    pTodoStatusTagYes,
			SubscribeStatus: intPtrOrNil(pTodoSubscribeStatus),
			UserCode:        pTodoUserCode,
			Executors:       parseJSONListFlag(pTodoExecutors, "--executors"),
			Copys:           parseJSONListFlag(pTodoCopys, "--copys"),
			Resources:       parseJSONListFlag(pTodoResources, "--resources"),
			Reminds:         parseJSONListFlag(pTodoReminds, "--reminds"),
			UserToken:       pTodoUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"todo_code"})
	},
}

var personalTodoUpdateCmd = &cobra.Command{
	Use:   "update TODO_CODE ORG_ID",
	Short: "Update selected fields of a personal todo",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().UpdatePersonalTodo(context.Background(), &lansenger.PersonalTodoUpdateParams{
			TodoCode:     args[0],
			OrgID:        args[1],
			UpdateFields: splitCommaList(pTodoUpdateFields),
			Subject:      pTodoSubject,
			Description:  pTodoDescription,
			DueTime:      int64PtrOrNil(pTodoDueTime),
			Priority:     intPtrOrNil(pTodoPriority),
			Executors:    parseJSONListFlag(pTodoExecutors, "--executors"),
			Copys:        parseJSONListFlag(pTodoCopys, "--copys"),
			Resources:    parseJSONListFlag(pTodoResources, "--resources"),
			Reminds:      parseJSONListFlag(pTodoReminds, "--reminds"),
			UserToken:    pTodoUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"todo_code"})
	},
}

var personalTodoListCmd = &cobra.Command{
	Use:   "list ORG_ID STAFF_ID",
	Short: "Page a user's personal todos",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchPersonalTodoList(
			context.Background(), args[0], args[1], pTodoPage, pTodoPageSize,
			intPtrOrNil(pTodoStatus), pTodoAppFilter, pTodoCategoryName, pTodoUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"page_no", "page_size", "pages", "total", "has_more", "items"})
	},
}

var personalTodoUploadResourceCmd = &cobra.Command{
	Use:   "upload-resource APP_ID FILE_NAME CONTENT_TYPE ORG_ID",
	Short: "Upload a local file resource",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		raw, err := os.ReadFile(pTodoFile)
		checkError(err)
		result, err := getClient().UploadPersonalTodoResource(context.Background(), &lansenger.PersonalTodoResourceUploadParams{
			AppID:         args[0],
			FileName:      args[1],
			ContentType:   args[2],
			OrgID:         args[3],
			Size:          int64(len(raw)),
			FileData:      base64.StdEncoding.EncodeToString(raw),
			Thumb:         pTodoThumb,
			ExtensionInfo: pTodoExtensionInfo,
			UserToken:     pTodoUserToken,
		})
		checkError(err)
		outputResultFields(result, []string{"resource_id", "file_name", "size", "mime_type", "download_url"})
	},
}

var personalTodoDownloadURLCmd = &cobra.Command{
	Use:   "download-url RESOURCE_ID ORG_ID",
	Short: "Fetch a resource download URL",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchPersonalTodoResourceDownloadURL(
			context.Background(), args[0], args[1], pTodoFileName, pTodoUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"url"})
	},
}

var personalTodoUploadURLCmd = &cobra.Command{
	Use:   "upload-url FILE_NAME MD5 SIZE ORG_ID",
	Short: "Fetch a presigned upload URL",
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		result, err := getClient().FetchPersonalTodoResourceUploadURL(
			context.Background(), args[0], args[1], pTodoSize, args[3], pTodoUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"url"})
	},
}

func int64PtrOrNil(v int64) *int64 {
	if v < 0 {
		return nil
	}
	val := v
	return &val
}

func init() {
	personalTodoSaveCmd.Flags().Int64Var(&pTodoStartTime, "start-time", -1, "Start time in epoch milliseconds")
	personalTodoSaveCmd.Flags().Int64Var(&pTodoDueTime, "due-time", -1, "Due time in epoch milliseconds")
	personalTodoSaveCmd.Flags().IntVar(&pTodoPriority, "priority", 1, "0=low, 1=normal, 2=urgent, 3=very urgent")
	personalTodoSaveCmd.Flags().StringVar(&pTodoDescription, "description", "", "Description")
	personalTodoSaveCmd.Flags().StringVar(&pTodoParentCode, "parent-code", "", "Parent todo code")
	personalTodoSaveCmd.Flags().StringVar(&pTodoStatusTagNo, "status-tag-no", "", "Unfinished status label")
	personalTodoSaveCmd.Flags().StringVar(&pTodoStatusTagYes, "status-tag-yes", "", "Finished status label")
	personalTodoSaveCmd.Flags().IntVar(&pTodoSubscribeStatus, "subscribe-status", -1, "Subscribe: 1=yes, 0=no (-1=omit)")
	personalTodoSaveCmd.Flags().StringVar(&pTodoUserCode, "user-code", "", "Private-chat executor authorization code")
	personalTodoSaveCmd.Flags().StringVar(&pTodoExecutors, "executors", "", "JSON list of executors")
	personalTodoSaveCmd.Flags().StringVar(&pTodoCopys, "copys", "", "JSON list of CC users")
	personalTodoSaveCmd.Flags().StringVar(&pTodoResources, "resources", "", "JSON list of resources")
	personalTodoSaveCmd.Flags().StringVar(&pTodoReminds, "reminds", "", "JSON list of reminders")
	personalTodoSaveCmd.Flags().StringVar(&pTodoUserToken, "user-token", "", "User token")

	personalTodoUpdateCmd.Flags().StringVar(&pTodoUpdateFields, "update-fields", "", "Comma-separated fields to update")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoSubject, "subject", "", "New subject")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoDescription, "description", "", "New description")
	personalTodoUpdateCmd.Flags().Int64Var(&pTodoDueTime, "due-time", -1, "New due time")
	personalTodoUpdateCmd.Flags().IntVar(&pTodoPriority, "priority", -1, "New priority")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoExecutors, "executors", "", "JSON list of executors")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoCopys, "copys", "", "JSON list of CC users")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoResources, "resources", "", "JSON list of resources")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoReminds, "reminds", "", "JSON list of reminders")
	personalTodoUpdateCmd.Flags().StringVar(&pTodoUserToken, "user-token", "", "User token")

	personalTodoListCmd.Flags().IntVar(&pTodoPage, "page", 1, "Page number")
	personalTodoListCmd.Flags().IntVar(&pTodoPageSize, "size", 10, "Page size")
	personalTodoListCmd.Flags().IntVar(&pTodoStatus, "status", -1, "0=unfinished, 1=finished (-1=all)")
	personalTodoListCmd.Flags().StringVar(&pTodoAppFilter, "app-id", "", "Filter by application ID")
	personalTodoListCmd.Flags().StringVar(&pTodoCategoryName, "app-category-name", "", "Filter by application category")
	personalTodoListCmd.Flags().StringVar(&pTodoUserToken, "user-token", "", "User token")

	personalTodoUploadResourceCmd.Flags().StringVar(&pTodoFile, "file", "", "Local file path")
	personalTodoUploadResourceCmd.Flags().BoolVar(&pTodoThumb, "thumb", false, "Generate a thumbnail")
	personalTodoUploadResourceCmd.Flags().StringVar(&pTodoExtensionInfo, "extension-info", "", "Extension information")
	personalTodoUploadResourceCmd.Flags().StringVar(&pTodoUserToken, "user-token", "", "User token")

	personalTodoDownloadURLCmd.Flags().StringVar(&pTodoFileName, "file-name", "", "Optional base64-encoded file name")
	personalTodoDownloadURLCmd.Flags().StringVar(&pTodoUserToken, "user-token", "", "User token")
	personalTodoUploadURLCmd.Flags().StringVar(&pTodoUserToken, "user-token", "", "User token")
	_ = personalTodoUploadResourceCmd.MarkFlagRequired("file")

	personalTodoUploadURLCmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		size, err := strconv.ParseInt(args[2], 10, 64)
		if err != nil {
			return err
		}
		pTodoSize = size
		return nil
	}

	personalTodoCmd.AddCommand(
		personalTodoSaveCmd,
		personalTodoUpdateCmd,
		personalTodoListCmd,
		personalTodoUploadResourceCmd,
		personalTodoDownloadURLCmd,
		personalTodoUploadURLCmd,
	)
	rootCmd.AddCommand(personalTodoCmd)
}
