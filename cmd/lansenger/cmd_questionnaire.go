package main

import (
	"context"
	"strconv"

	lansenger "github.com/lansenger-pm/lansenger-sdk-go"

	"github.com/spf13/cobra"
)

var questionnaireWithdrawCmd = &cobra.Command{
	Use:   "withdraw QUESTIONNAIRE_CODE",
	Short: "Withdraw a published questionnaire back to draft",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("withdraw", []string{"done"}),
}

var questionnaireFinishCmd = &cobra.Command{
	Use:   "finish QUESTIONNAIRE_CODE",
	Short: "End an ongoing questionnaire (no more answers)",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("finish", []string{"done"}),
}

var questionnaireDeleteCmd = &cobra.Command{
	Use:   "delete QUESTIONNAIRE_CODE",
	Short: "Delete a questionnaire",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("delete", []string{"done"}),
}

var questionnaireDetailCmd = &cobra.Command{
	Use:   "detail QUESTIONNAIRE_CODE",
	Short: "Fetch full questionnaire detail incl. questions (needs admin permission)",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("detail", []string{"code", "title", "status", "question_count", "answer_user_count", "account_code", "questions"}),
}

var questionnaireBriefCmd = &cobra.Command{
	Use:   "brief QUESTIONNAIRE_CODE",
	Short: "Fetch questionnaire detail without admin check (no questions)",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("brief", []string{"code", "title", "status", "question_count", "answer_user_count", "account_code"}),
}

var questionnaireAnswerURLCmd = &cobra.Command{
	Use:   "answer-url QUESTIONNAIRE_CODE",
	Short: "Fetch the answer-page URL",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("answer_url", []string{"url"}),
}

var questionnaireCopyCmd = &cobra.Command{
	Use:   "copy QUESTIONNAIRE_CODE",
	Short: "Copy a questionnaire into a new draft",
	Args:  cobra.ExactArgs(1),
	Run:   questionnaireCodeRunner("copy", []string{"new_code"}),
}

var questionnaireCmd = &cobra.Command{
	Use:   "questionnaire",
	Short: "Questionnaire operations (问卷系统)",
}

var (
	qSaveCode, qSaveWelcome, qSaveBye                                              string
	qSaveCover, qSaveResourceIDs, qSaveCreateMobile, qSaveCreateUserID, qSaveToken string
)

var questionnaireSaveCmd = &cobra.Command{
	Use:   "save TITLE ACCOUNT_CODE",
	Short: "Create a questionnaire (or update when --code is given)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.SaveQuestionnaire(context.Background(), &lansenger.QuestionnaireSaveParams{
			Title: args[0], AccountCode: args[1], Code: qSaveCode, WelcomeSpeech: qSaveWelcome,
			ByeSpeech: qSaveBye, CoverResourceID: qSaveCover, ResourceIDs: qSaveResourceIDs,
			CreateMobile: qSaveCreateMobile, CreateUserID: questionnaireUserID(qSaveCreateUserID),
			UserToken: qSaveToken,
		})
		checkError(err)
		outputResultFields(result, []string{"questionnaire_code"})
	},
}

var questionnaireSaveQuestionsCmd = &cobra.Command{
	Use:   "save-questions QUESTIONNAIRE_CODE",
	Short: "Batch-save questions (new or update; 16 question types)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		questions := parseJSONListFlag(qSaveQuestionsJSON, "--questions")
		result, err := client.SaveQuestionnaireQuestions(
			context.Background(), args[0], questions, questionnaireUserID(qSaveQCreateUserID), qSaveQToken,
		)
		checkError(err)
		outputResultFields(result, []string{"saved_count"})
	},
}

var questionnaireDeleteQuestionCmd = &cobra.Command{
	Use:   "delete-question QUESTION_CODE",
	Short: "Delete a question by code",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		confirmHighRisk("delete", "question "+args[0], qDelQYes, qDelQDryRun)
		if qDelQDryRun {
			return
		}
		client := getClient()
		result, err := client.DeleteQuestionnaireQuestion(
			context.Background(), args[0], questionnaireUserID(qDelQCreateUserID), qDelQToken,
		)
		checkError(err)
		outputResultFields(result, []string{"deleted"})
	},
}

var questionnairePublishCmd = &cobra.Command{
	Use:   "publish QUESTIONNAIRE_CODE",
	Short: "Publish a questionnaire",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.PublishQuestionnaire(context.Background(), &lansenger.QuestionnairePublishParams{
			QuestionnaireCode: args[0],
			ScopeType:         qPublishScope,
			StaffIDs:          splitCommaList(qPublishStaffIDs),
			Phones:            splitCommaList(qPublishPhones),
			AnswerLimit:       qPublishAnswerLimit,
			MessageFlag:       qPublishMessageFlag,
			PageFlag:          qPublishPageFlag,
			ShareFlag:         qPublishShareFlag,
			ViewStatsFlag:     qPublishViewStatsFlag,
			AnonymFlag:        qPublishAnonymFlag,
			PublishUserID:     questionnaireUserID(qPublishUserID),
			UserToken:         qPublishToken,
		})
		checkError(err)
		outputResultFields(result, []string{"done"})
	},
}

func questionnaireCodeRunner(category string, fields []string) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		client := getClient()
		ctx := context.Background()
		switch category {
		case "withdraw":
			result, err := client.WithdrawQuestionnaire(ctx, args[0], questionnaireUserID(qOperateUserID), qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		case "finish":
			result, err := client.FinishQuestionnaire(ctx, args[0], questionnaireUserID(qOperateUserID), qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		case "delete":
			confirmHighRisk("delete", "questionnaire "+args[0], qDeleteYes, qDeleteDryRun)
			if qDeleteDryRun {
				return
			}
			result, err := client.DeleteQuestionnaire(ctx, args[0], questionnaireUserID(qOperateUserID), qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		case "detail":
			result, err := client.FetchQuestionnaireDetail(ctx, args[0], questionnaireUserID(qOperateUserID), qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		case "brief":
			result, err := client.FetchQuestionnaireBrief(ctx, args[0], qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		case "answer_url":
			result, err := client.FetchQuestionnaireAnswerURL(ctx, args[0], questionnaireUserID(qOperateUserID), qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		case "copy":
			result, err := client.CopyQuestionnaire(ctx, args[0], questionnaireUserID(qOperateUserID), qUserToken)
			checkError(err)
			outputResultFields(result, fields)
		}
	}
}

var questionnaireQueryCodesCmd = &cobra.Command{
	Use:   "query-codes",
	Short: "Batch-fetch questionnaire basic info by codes",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchQuestionnairesByCodes(context.Background(), splitCommaList(qQueryCodes), qQueryIncludeDeleted, qUserToken)
		checkError(err)
		outputResultFields(result, []string{"total", "items"})
	},
}

var questionnaireAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "Fetch office accounts the user can manage (accountCode source)",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchQuestionnaireOfficeAccounts(context.Background(), questionnaireUserID(qAccountsUserID), qUserToken)
		checkError(err)
		outputResultFields(result, []string{"total", "accounts"})
	},
}

var questionnaireCreatedListCmd = &cobra.Command{
	Use:   "created-list ACCOUNT_CODE",
	Short: "Page questionnaires created under an office account",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchCreatedQuestionnaires(
			context.Background(), args[0], qPage, qPageSize, intPtrOrNil(qStatus), questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"total", "page_no", "page_size", "has_more", "items"})
	},
}

var questionnaireMyCreatedCmd = &cobra.Command{
	Use:   "my-created ORG_ID",
	Short: "Page all questionnaires I created (personal + official)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchMyCreatedQuestionnaires(
			context.Background(), args[0], qPage, qPageSize, qTitleFilter, intPtrOrNil(qStatus),
			questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"total", "page_no", "page_size", "has_more", "items"})
	},
}

var questionnaireParticipatedCmd = &cobra.Command{
	Use:   "participated ORG_ID",
	Short: "Page questionnaires the user answered",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchParticipatedQuestionnaires(
			context.Background(), args[0], qPage, qPageSize, intPtrOrNil(qStatus),
			questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"total", "page_no", "page_size", "has_more", "items"})
	},
}

var questionnaireAnswersCmd = &cobra.Command{
	Use:   "answers ACCOUNT_CODE QUESTIONNAIRE_CODE",
	Short: "Page answer records of a questionnaire",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchAnswerRecords(
			context.Background(), args[0], args[1], qPage, qPageSize, questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"total", "page_no", "page_size", "has_more", "items"})
	},
}

var questionnaireAnswerDetailCmd = &cobra.Command{
	Use:   "answer-detail ACCOUNT_CODE ANSWER_CODE",
	Short: "Fetch one answer record's full detail",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchQuestionnaireAnswerDetail(
			context.Background(), args[0], args[1], questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"answer_code", "answer_user_name", "answer_status", "answer_commit_time", "questionnaire", "answers"})
	},
}

var questionnaireLastAnswerDetailCmd = &cobra.Command{
	Use:   "last-answer-detail QUESTIONNAIRE_CODE",
	Short: "Fetch the user's last answer detail",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchQuestionnaireLastAnswerDetail(
			context.Background(), args[0], qAnswerRecordCode, questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"answer_code", "answer_user_name", "answer_status", "answer_commit_time", "answers"})
	},
}

var questionnaireAnswerDataCmd = &cobra.Command{
	Use:   "answer-data ACCOUNT_CODE QUESTIONNAIRE_CODE",
	Short: "Page answer data for export (JSON structure)",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchAnswerData(
			context.Background(), args[0], args[1], qPage, qPageSize, questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"total", "page_no", "page_size", "has_more", "items"})
	},
}

var questionnaireLastAnswerRecordCmd = &cobra.Command{
	Use:   "last-answer-record QUESTIONNAIRE_CODE",
	Short: "Fetch the user's last answer record (main table only)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		result, err := client.FetchQuestionnaireLastAnswerRecord(
			context.Background(), args[0], qAnswerRecordCode, questionnaireUserID(qUserID), qUserToken,
		)
		checkError(err)
		outputResultFields(result, []string{"record_code", "answer_user_name", "answer_status", "answer_commit_time", "stats_status"})
	},
}

var questionnaireUploadURLCmd = &cobra.Command{
	Use:   "upload-url FILE_NAME MD5 SIZE",
	Short: "Fetch a presigned upload URL (then PUT with Content-MD5 header)",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		client := getClient()
		size, err := strconv.ParseInt(args[2], 10, 64)
		checkError(err)
		result, err := client.FetchQuestionnaireUploadURL(context.Background(), args[0], args[1], size, qUserToken)
		checkError(err)
		outputResultFields(result, []string{"url"})
	},
}

var (
	qSaveQuestionsJSON, qSaveQCreateUserID, qSaveQToken                            string
	qDelQCreateUserID, qDelQToken                                                  string
	qDelQYes, qDelQDryRun                                                          bool
	qPublishScope, qPublishAnswerLimit, qPublishMessageFlag                        int
	qPublishPageFlag, qPublishShareFlag, qPublishViewStatsFlag, qPublishAnonymFlag int
	qPublishStaffIDs, qPublishPhones, qPublishUserID, qPublishToken                string
	qOperateUserID, qUserToken                                                     string
	qDeleteYes, qDeleteDryRun                                                      bool
	qQueryCodes                                                                    string
	qQueryIncludeDeleted                                                           int
	qAccountsUserID                                                                string
	qUserID                                                                        string
	qPage, qPageSize, qStatus                                                      int
	qTitleFilter, qAnswerRecordCode                                                string
)

func questionnaireUserID(value string) string {
	if value != "" {
		return value
	}
	return globalAsStaffID
}

func intPtrOrNil(v int) *int {
	if v < 0 {
		return nil
	}
	val := v
	return &val
}

func parseJSONListFlag(raw, flagName string) []map[string]interface{} {
	if raw == "" {
		return nil
	}
	decoded, err := parseJSONRaw(raw)
	checkError(err)
	list, ok := decoded.([]interface{})
	if !ok {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(list))
	for _, item := range list {
		if m, ok := item.(map[string]interface{}); ok {
			out = append(out, m)
		}
	}
	return out
}

func init() {
	questionnaireSaveCmd.Flags().StringVar(&qSaveCode, "code", "", "Existing questionnaire code (update when given)")
	questionnaireSaveCmd.Flags().StringVar(&qSaveWelcome, "welcome", "", "Welcome speech")
	questionnaireSaveCmd.Flags().StringVar(&qSaveBye, "bye", "", "Bye speech")
	questionnaireSaveCmd.Flags().StringVar(&qSaveCover, "cover", "", "Cover resource ID")
	questionnaireSaveCmd.Flags().StringVar(&qSaveResourceIDs, "resource-ids", "", "Comma-joined resource id string")
	questionnaireSaveCmd.Flags().StringVar(&qSaveCreateMobile, "create-mobile", "", "Operator mobile (omit when --as/--user-token is set)")
	questionnaireSaveCmd.Flags().StringVar(&qSaveCreateUserID, "create-user-id", "", "Creator staff ID (omit when --as/--user-token is set)")
	questionnaireSaveCmd.Flags().StringVar(&qSaveToken, "user-token", "", "User token")

	questionnaireSaveQuestionsCmd.Flags().StringVar(&qSaveQuestionsJSON, "questions", "", `JSON list of questions (camelCase, required)`)
	questionnaireSaveQuestionsCmd.Flags().StringVar(&qSaveQCreateUserID, "create-user-id", "", "Creator staff ID (omit when --as/--user-token is set)")
	questionnaireSaveQuestionsCmd.Flags().StringVar(&qSaveQToken, "user-token", "", "User token")

	questionnaireDeleteQuestionCmd.Flags().StringVar(&qDelQCreateUserID, "create-user-id", "", "Creator staff ID (omit when --as/--user-token is set)")
	questionnaireDeleteQuestionCmd.Flags().StringVar(&qDelQToken, "user-token", "", "User token")
	questionnaireDeleteQuestionCmd.Flags().BoolVarP(&qDelQYes, "yes", "y", false, "Confirm question deletion before executing")
	questionnaireDeleteQuestionCmd.Flags().BoolVar(&qDelQDryRun, "dry-run", false, "Validate inputs without deleting")

	questionnairePublishCmd.Flags().IntVar(&qPublishScope, "scope", 1, "Publish scope: 1=internal, 2=public")
	questionnairePublishCmd.Flags().StringVar(&qPublishStaffIDs, "staff-ids", "", "Comma-separated target staff openIds (internal scope)")
	questionnairePublishCmd.Flags().StringVar(&qPublishPhones, "phones", "", "Comma-separated target mobiles (internal scope)")
	questionnairePublishCmd.Flags().IntVar(&qPublishAnswerLimit, "answer-limit", 1, "Answer limit: 1=once, -1=unlimited")
	questionnairePublishCmd.Flags().IntVar(&qPublishMessageFlag, "message-flag", 0, "Official-account message: 1=on, 0=off")
	questionnairePublishCmd.Flags().IntVar(&qPublishPageFlag, "page-flag", 0, "One question per page: 1=on, 0=off")
	questionnairePublishCmd.Flags().IntVar(&qPublishShareFlag, "share-flag", 0, "Allow sharing: 1=on, 0=off")
	questionnairePublishCmd.Flags().IntVar(&qPublishViewStatsFlag, "view-stats-flag", 1, "Allow viewing stats: 1=on, 0=off")
	questionnairePublishCmd.Flags().IntVar(&qPublishAnonymFlag, "anonym-flag", 0, "Allow anonymous answers: 1=on, 0=off")
	questionnairePublishCmd.Flags().StringVar(&qPublishUserID, "publish-user-id", "", "Operator staff ID (omit when --as/--user-token is set)")
	questionnairePublishCmd.Flags().StringVar(&qPublishToken, "user-token", "", "User token")

	for _, c := range []*cobra.Command{questionnaireWithdrawCmd, questionnaireFinishCmd, questionnaireDeleteCmd, questionnaireDetailCmd, questionnaireAnswerURLCmd, questionnaireCopyCmd} {
		c.Flags().StringVar(&qOperateUserID, "operate-user-id", "", "Operator staff ID (omit when --as/--user-token is set)")
		c.Flags().StringVar(&qUserToken, "user-token", "", "User token")
	}
	questionnaireBriefCmd.Flags().StringVar(&qUserToken, "user-token", "", "User token")
	questionnaireDeleteCmd.Flags().BoolVarP(&qDeleteYes, "yes", "y", false, "Confirm questionnaire deletion before executing")
	questionnaireDeleteCmd.Flags().BoolVar(&qDeleteDryRun, "dry-run", false, "Validate inputs without deleting")

	questionnaireQueryCodesCmd.Flags().StringVar(&qQueryCodes, "codes", "", "Comma-separated questionnaire codes (required)")
	questionnaireQueryCodesCmd.Flags().IntVar(&qQueryIncludeDeleted, "include-deleted", 0, "0=exclude deleted, 1=include")
	questionnaireQueryCodesCmd.Flags().StringVar(&qUserToken, "user-token", "", "User token")

	questionnaireAccountsCmd.Flags().StringVar(&qAccountsUserID, "user-id", "", "User ID (omit when --as/--user-token is set)")
	questionnaireAccountsCmd.Flags().StringVar(&qUserToken, "user-token", "", "User token")

	for _, c := range []*cobra.Command{questionnaireCreatedListCmd, questionnaireMyCreatedCmd, questionnaireParticipatedCmd, questionnaireAnswersCmd, questionnaireAnswerDataCmd} {
		c.Flags().IntVar(&qPage, "page", 1, "Page number")
		c.Flags().IntVar(&qPageSize, "size", 10, "Page size")
		c.Flags().StringVar(&qUserID, "user-id", "", "User ID (omit when --as/--user-token is set)")
		c.Flags().StringVar(&qUserToken, "user-token", "", "User token")
	}
	questionnaireCreatedListCmd.Flags().IntVar(&qStatus, "status", -1, "1=draft, 2=ongoing, 3=withdrawn, 4=finished (-1=omit)")
	questionnaireMyCreatedCmd.Flags().StringVar(&qTitleFilter, "title", "", "Filter by title")
	questionnaireMyCreatedCmd.Flags().IntVar(&qStatus, "status", -1, "1=draft, 2=ongoing, 3=withdrawn, 4=finished (-1=omit)")
	questionnaireParticipatedCmd.Flags().IntVar(&qStatus, "status", -1, "2=ongoing, 4=finished (-1=omit)")

	questionnaireLastAnswerDetailCmd.Flags().StringVar(&qAnswerRecordCode, "answer-record-code", "", "Specific answer record code")
	questionnaireLastAnswerDetailCmd.Flags().StringVar(&qUserID, "user-id", "", "User ID (omit when --as/--user-token is set)")
	questionnaireLastAnswerDetailCmd.Flags().StringVar(&qUserToken, "user-token", "", "User token")
	questionnaireLastAnswerRecordCmd.Flags().StringVar(&qAnswerRecordCode, "answer-record-code", "", "Specific answer record code")
	questionnaireLastAnswerRecordCmd.Flags().StringVar(&qUserID, "user-id", "", "User ID (omit when --as/--user-token is set)")
	questionnaireLastAnswerRecordCmd.Flags().StringVar(&qUserToken, "user-token", "", "User token")
	questionnaireAnswerDetailCmd.Flags().StringVar(&qUserID, "user-id", "", "User ID (omit when --as/--user-token is set)")
	questionnaireAnswerDetailCmd.Flags().StringVar(&qUserToken, "user-token", "", "User token")

	questionnaireCmd.AddCommand(questionnaireSaveCmd)
	questionnaireCmd.AddCommand(questionnaireSaveQuestionsCmd)
	questionnaireCmd.AddCommand(questionnaireDeleteQuestionCmd)
	questionnaireCmd.AddCommand(questionnairePublishCmd)
	questionnaireCmd.AddCommand(questionnaireWithdrawCmd)
	questionnaireCmd.AddCommand(questionnaireFinishCmd)
	questionnaireCmd.AddCommand(questionnaireDeleteCmd)
	questionnaireCmd.AddCommand(questionnaireDetailCmd)
	questionnaireCmd.AddCommand(questionnaireBriefCmd)
	questionnaireCmd.AddCommand(questionnaireAnswerURLCmd)
	questionnaireCmd.AddCommand(questionnaireCopyCmd)
	questionnaireCmd.AddCommand(questionnaireQueryCodesCmd)
	questionnaireCmd.AddCommand(questionnaireAccountsCmd)
	questionnaireCmd.AddCommand(questionnaireCreatedListCmd)
	questionnaireCmd.AddCommand(questionnaireMyCreatedCmd)
	questionnaireCmd.AddCommand(questionnaireParticipatedCmd)
	questionnaireCmd.AddCommand(questionnaireAnswersCmd)
	questionnaireCmd.AddCommand(questionnaireAnswerDetailCmd)
	questionnaireCmd.AddCommand(questionnaireLastAnswerDetailCmd)
	questionnaireCmd.AddCommand(questionnaireAnswerDataCmd)
	questionnaireCmd.AddCommand(questionnaireLastAnswerRecordCmd)
	questionnaireCmd.AddCommand(questionnaireUploadURLCmd)
	rootCmd.AddCommand(questionnaireCmd)
}
