package lansenger

import (
	"context"
)

// QuestionnaireSaveParams carries the fields for SaveQuestionnaire (问卷系统 /v1/saveQuestionnaire).
type QuestionnaireSaveParams struct {
	Title           string
	AccountCode     string // official account CODE — missing/unknown → errCode 3104
	Code            string // existing code to update (empty = create)
	WelcomeSpeech   string
	ByeSpeech       string
	CoverResourceID string
	ResourceIDs     string
	AppID           string
	UserType        *int // 1=phone, 2=openid; nil when user_token is provided
	CreateMobile    string
	CreateUserID    string
	UserToken       string
}

// QuestionnairePublishParams carries the fields for PublishQuestionnaire (问卷系统 /v1/publish).
type QuestionnairePublishParams struct {
	QuestionnaireCode string
	ScopeType         int    // 1=internal, 2=public
	StaffIDs          []string
	Phones            []string
	AnswerLimit       int    // 1=once, -1=unlimited
	MessageFlag       int    // 1=on, 0=off (server default 0)
	PageFlag          int    // server default 0
	ShareFlag         int    // server default 0
	ViewStatsFlag     int    // server default 1
	AnonymFlag        int    // server default 0
	PublishUserID     string
	UserToken         string
}

// SaveQuestionnaire creates a questionnaire, or overwrites an existing one when Code is set.
func (c *LansengerClient) SaveQuestionnaire(ctx context.Context, p *QuestionnaireSaveParams) (*QuestionnaireSaveResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "save", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{"title": p.Title, "accountCode": p.AccountCode}
	if p.Code != "" {
		body["code"] = p.Code
	}
	if p.WelcomeSpeech != "" {
		body["welcomeSpeech"] = p.WelcomeSpeech
	}
	if p.ByeSpeech != "" {
		body["byeSpeech"] = p.ByeSpeech
	}
	if p.CoverResourceID != "" {
		body["coverResourceId"] = p.CoverResourceID
	}
	if p.ResourceIDs != "" {
		body["resourceIds"] = p.ResourceIDs
	}
	if p.AppID != "" {
		body["appId"] = p.AppID
	}
	if p.UserType != nil {
		body["userType"] = *p.UserType
	}
	if p.CreateMobile != "" {
		body["createMobile"] = p.CreateMobile
	}
	if p.CreateUserID != "" {
		body["createUserId"] = p.CreateUserID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireSaveResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireSaveResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.QuestionnaireCode = strOrNil(data, "data")
	}
	return res, nil
}

// SaveQuestionnaireQuestions batch-saves questions of a questionnaire (16 types).
// questionList items are passed through as camelCase dicts.
func (c *LansengerClient) SaveQuestionnaireQuestions(ctx context.Context, questionnaireCode string, questionList []map[string]interface{}, createUserID, userToken string) (*QuestionnaireQuestionSaveResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "questions_save", token, WithUserToken(userToken))
	body := map[string]interface{}{"questionnaireCode": questionnaireCode, "questionList": questionList}
	if createUserID != "" {
		body["createUserId"] = createUserID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireQuestionSaveResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireQuestionSaveResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.SavedCount = intFromMap(data, "data")
	}
	return res, nil
}

// DeleteQuestionnaireQuestion deletes a question by code.
func (c *LansengerClient) DeleteQuestionnaireQuestion(ctx context.Context, questionCode, createUserID, userToken string) (*QuestionnaireQuestionDeleteResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "question_delete", token, WithUserToken(userToken))
	body := map[string]interface{}{"questionCode": questionCode}
	if createUserID != "" {
		body["createUserId"] = createUserID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireQuestionDeleteResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireQuestionDeleteResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.Deleted = boolFromMap(data, "data")
	}
	return res, nil
}

// PublishQuestionnaire publishes a questionnaire.
func (c *LansengerClient) PublishQuestionnaire(ctx context.Context, p *QuestionnairePublishParams) (*QuestionnaireOpResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "publish", token, WithUserToken(p.UserToken))
	body := map[string]interface{}{
		"questionnaireCode": p.QuestionnaireCode,
		"scopeType":         p.ScopeType,
		"answerLimit":       p.AnswerLimit,
		"messageFlag":       p.MessageFlag,
		"pageFlag":          p.PageFlag,
		"shareFlag":         p.ShareFlag,
		"viewStatsFlag":     p.ViewStatsFlag,
		"anonymFlag":        p.AnonymFlag,
	}
	if len(p.StaffIDs) > 0 {
		body["staffIdList"] = p.StaffIDs
	}
	if len(p.Phones) > 0 {
		body["phoneList"] = p.Phones
	}
	if p.PublishUserID != "" {
		body["publishUserId"] = p.PublishUserID
	}

	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireOpResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.Done = boolFromMap(data, "data")
	}
	return res, nil
}

func questionnaireCodeOp(ctx context.Context, c *LansengerClient, category, code, operateUserID, userToken string) (map[string]interface{}, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", category, token, WithUserToken(userToken))
	body := map[string]interface{}{"questionnaireCode": code}
	if operateUserID != "" {
		body["operateUserId"] = operateUserID
	}
	return c.doPost(ctx, url, body)
}

// WithdrawQuestionnaire withdraws a published questionnaire back to draft.
func (c *LansengerClient) WithdrawQuestionnaire(ctx context.Context, code, operateUserID, userToken string) (*QuestionnaireOpResult, error) {
	result, err := questionnaireCodeOp(ctx, c, "withdraw", code, operateUserID, userToken)
	if err != nil {
		return &QuestionnaireOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireOpResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.Done = boolFromMap(data, "data")
	}
	return res, nil
}

// FinishQuestionnaire ends an ongoing questionnaire.
func (c *LansengerClient) FinishQuestionnaire(ctx context.Context, code, operateUserID, userToken string) (*QuestionnaireOpResult, error) {
	result, err := questionnaireCodeOp(ctx, c, "finish", code, operateUserID, userToken)
	if err != nil {
		return &QuestionnaireOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireOpResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.Done = boolFromMap(data, "data")
	}
	return res, nil
}

// DeleteQuestionnaire deletes a questionnaire.
func (c *LansengerClient) DeleteQuestionnaire(ctx context.Context, code, operateUserID, userToken string) (*QuestionnaireOpResult, error) {
	result, err := questionnaireCodeOp(ctx, c, "delete", code, operateUserID, userToken)
	if err != nil {
		return &QuestionnaireOpResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireOpResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.Done = boolFromMap(data, "data")
	}
	return res, nil
}

func fillQuestionnaireDetail(res *QuestionnaireDetailResult, data map[string]interface{}) {
	res.QuestionnaireID = intFromMap(data, "id")
	res.Code = strFromMap(data, "code")
	res.Title = strFromMap(data, "title")
	res.Status = intFromMap(data, "status")
	res.AccountType = intFromMap(data, "accountType")
	res.AccountCode = strFromMap(data, "accountCode")
	res.AnswerUserCount = intFromMap(data, "answerUserCount")
	res.AnswerUserTimes = intFromMap(data, "answerUserTimes")
	res.QuestionCount = intFromMap(data, "questionCount")
	res.PublishTime = int64(intFromMap(data, "publishTime"))
	res.PublishUserName = strFromMap(data, "publishUserName")
	res.CreateUserName = strFromMap(data, "createUserName")
	res.CreateTime = int64(intFromMap(data, "createTime"))
	if list, ok := data["questionList"].([]interface{}); ok {
		res.Questions = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Questions = append(res.Questions, m)
			}
		}
	}
}

// FetchQuestionnaireDetail fetches full detail incl. questions (admin permission required).
func (c *LansengerClient) FetchQuestionnaireDetail(ctx context.Context, code, operateUserID, userToken string) (*QuestionnaireDetailResult, error) {
	result, err := questionnaireCodeOp(ctx, c, "detail", code, operateUserID, userToken)
	if err != nil {
		return &QuestionnaireDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillQuestionnaireDetail(res, data)
	}
	return res, nil
}

// FetchQuestionnaireBrief fetches detail without admin check (no question list).
func (c *LansengerClient) FetchQuestionnaireBrief(ctx context.Context, code, userToken string) (*QuestionnaireDetailResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "detail_no_auth", token, WithUserToken(userToken))
	result, err := c.doPost(ctx, url, map[string]interface{}{"questionnaireCode": code})
	if err != nil {
		return &QuestionnaireDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillQuestionnaireDetail(res, data)
	}
	return res, nil
}

// FetchQuestionnaireAnswerURL fetches the answer-page URL.
func (c *LansengerClient) FetchQuestionnaireAnswerURL(ctx context.Context, code, operateUserID, userToken string) (*QuestionnaireAnswerUrlResult, error) {
	result, err := questionnaireCodeOp(ctx, c, "answer_url", code, operateUserID, userToken)
	if err != nil {
		return &QuestionnaireAnswerUrlResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireAnswerUrlResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.URL = strOrNil(data, "data")
	}
	return res, nil
}

// CopyQuestionnaire copies a questionnaire into a new draft; returns the new code.
func (c *LansengerClient) CopyQuestionnaire(ctx context.Context, code, operateUserID, userToken string) (*QuestionnaireCopyResult, error) {
	result, err := questionnaireCodeOp(ctx, c, "copy", code, operateUserID, userToken)
	if err != nil {
		return &QuestionnaireCopyResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireCopyResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.NewCode = strOrNil(data, "data")
	}
	return res, nil
}

// FetchQuestionnairesByCodes batch-fetches questionnaire basic info by codes.
func (c *LansengerClient) FetchQuestionnairesByCodes(ctx context.Context, codeList []string, includeDeleted int, userToken string) (*QuestionnaireQueryListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "query_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"codeList": codeList, "includeDel": includeDeleted}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireQueryListResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireQueryListResult{Success: true, RawResponse: result}
	if list := extractDataArray(result); list != nil {
		res.Items = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Items = append(res.Items, m)
			}
		}
		res.Total = len(res.Items)
	}
	return res, nil
}

// FetchQuestionnaireOfficeAccounts fetches office accounts the user can manage.
func (c *LansengerClient) FetchQuestionnaireOfficeAccounts(ctx context.Context, userID, userToken string) (*QuestionnaireAccountListResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "user_accounts", token, WithUserToken(userToken))
	body := map[string]interface{}{}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireAccountListResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireAccountListResult{Success: true, RawResponse: result}
	if list := extractDataArray(result); list != nil {
		res.Accounts = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Accounts = append(res.Accounts, m)
			}
		}
		res.Total = len(res.Accounts)
	}
	return res, nil
}

func fillQuestionnairePage(res *QuestionnairePageResult, result map[string]interface{}) {
	data := extractData(result)
	if data == nil {
		return
	}
	res.PageNo = intFromMap(data, "pageNo")
	res.PageSize = intFromMap(data, "pageSize")
	res.Pages = intFromMap(data, "pages")
	res.Total = intFromMap(data, "total")
	res.HasMore = boolFromMap(data, "hasNextPage")
	if list, ok := data["result"].([]interface{}); ok {
		res.Items = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Items = append(res.Items, m)
			}
		}
	}
}

// FetchCreatedQuestionnaires pages questionnaires created under an office account.
func (c *LansengerClient) FetchCreatedQuestionnaires(ctx context.Context, accountCode string, pageNo, pageSize int, status *int, userID, userToken string) (*QuestionnairePageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "create_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize, "accountCode": accountCode}
	if status != nil {
		body["status"] = *status
	}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnairePageResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnairePageResult{Success: true, RawResponse: result}
	fillQuestionnairePage(res, result)
	return res, nil
}

// FetchMyCreatedQuestionnaires pages all questionnaires the user created (personal + official).
func (c *LansengerClient) FetchMyCreatedQuestionnaires(ctx context.Context, orgID string, pageNo, pageSize int, title string, status *int, userID, userToken string) (*QuestionnairePageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "my_create_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize, "orgId": orgID}
	if title != "" {
		body["title"] = title
	}
	if status != nil {
		body["status"] = *status
	}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnairePageResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnairePageResult{Success: true, RawResponse: result}
	fillQuestionnairePage(res, result)
	return res, nil
}

// FetchParticipatedQuestionnaires pages questionnaires the user answered.
func (c *LansengerClient) FetchParticipatedQuestionnaires(ctx context.Context, orgID string, pageNo, pageSize int, status *int, userID, userToken string) (*QuestionnairePageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "participation_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize, "orgId": orgID}
	if status != nil {
		body["status"] = *status
	}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnairePageResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnairePageResult{Success: true, RawResponse: result}
	fillQuestionnairePage(res, result)
	return res, nil
}

// FetchAnswerRecords pages answer records of a questionnaire.
func (c *LansengerClient) FetchAnswerRecords(ctx context.Context, accountCode, questionnaireCode string, pageNo, pageSize int, userID, userToken string) (*QuestionnairePageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "answer_list", token, WithUserToken(userToken))
	body := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize, "accountCode": accountCode, "questionnaireCode": questionnaireCode}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnairePageResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnairePageResult{Success: true, RawResponse: result}
	fillQuestionnairePage(res, result)
	return res, nil
}

func fillQuestionnaireAnswerDetail(res *QuestionnaireAnswerDetailResult, data map[string]interface{}, fallbackCode string) {
	res.AnswerCode = strFromMap(data, "code")
	if res.AnswerCode == "" {
		res.AnswerCode = fallbackCode
	}
	res.AnswerUserID = strFromMap(data, "answerUserId")
	res.AnswerUserName = strFromMap(data, "answerUserName")
	res.AnswerStatus = intFromMap(data, "answerStatus")
	res.AnswerType = intFromMap(data, "answerType")
	res.AnswerUseTime = int64(intFromMap(data, "answerUseTime"))
	res.AnswerQuestionCount = intFromMap(data, "answerQuestionCount")
	res.AnswerCommitTime = int64(intFromMap(data, "answerCommitTime"))
	if m, ok := data["questionnaire"].(map[string]interface{}); ok {
		res.Questionnaire = m
	}
	if list, ok := data["questionList"].([]interface{}); ok {
		res.Questions = make([]map[string]interface{}, 0, len(list))
		for _, item := range list {
			if m, ok := item.(map[string]interface{}); ok {
				res.Questions = append(res.Questions, m)
			}
		}
	}
	if m, ok := data["answerMap"].(map[string]interface{}); ok {
		res.Answers = m
	}
}

// FetchQuestionnaireAnswerDetail fetches one answer record's full detail.
func (c *LansengerClient) FetchQuestionnaireAnswerDetail(ctx context.Context, accountCode, answerCode, userID, userToken string) (*QuestionnaireAnswerDetailResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "answer_detail", token, WithUserToken(userToken))
	body := map[string]interface{}{"accountCode": accountCode, "answerCode": answerCode}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireAnswerDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireAnswerDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillQuestionnaireAnswerDetail(res, data, answerCode)
	}
	return res, nil
}

// FetchQuestionnaireLastAnswerDetail fetches the user's last answer detail.
func (c *LansengerClient) FetchQuestionnaireLastAnswerDetail(ctx context.Context, code, answerRecordCode, userID, userToken string) (*QuestionnaireAnswerDetailResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "last_answer_detail", token, WithUserToken(userToken))
	body := map[string]interface{}{"questionnaireCode": code}
	if answerRecordCode != "" {
		body["answerRecordCode"] = answerRecordCode
	}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireAnswerDetailResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireAnswerDetailResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		fillQuestionnaireAnswerDetail(res, data, "")
	}
	return res, nil
}

// FetchAnswerData pages answer data for export (items carry the raw answerMap).
func (c *LansengerClient) FetchAnswerData(ctx context.Context, accountCode, questionnaireCode string, pageNo, pageSize int, userID, userToken string) (*QuestionnairePageResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "answer_data", token, WithUserToken(userToken))
	body := map[string]interface{}{"pageNo": pageNo, "pageSize": pageSize, "accountCode": accountCode, "questionnaireCode": questionnaireCode}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnairePageResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnairePageResult{Success: true, RawResponse: result}
	fillQuestionnairePage(res, result)
	return res, nil
}

// FetchQuestionnaireLastAnswerRecord fetches the user's last answer record (main table).
func (c *LansengerClient) FetchQuestionnaireLastAnswerRecord(ctx context.Context, code, answerRecordCode, userID, userToken string) (*QuestionnaireRecordResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "last_answer_record", token, WithUserToken(userToken))
	body := map[string]interface{}{"questionnaireCode": code}
	if answerRecordCode != "" {
		body["answerRecordCode"] = answerRecordCode
	}
	if userID != "" {
		body["userId"] = userID
	}
	result, err := c.doPost(ctx, url, body)
	if err != nil {
		return &QuestionnaireRecordResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireRecordResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.RecordID = intFromMap(data, "id")
		res.RecordCode = strFromMap(data, "code")
		res.AnswerUserID = strFromMap(data, "answerUserId")
		res.AnswerUserName = strFromMap(data, "answerUserName")
		res.AnswerStatus = intFromMap(data, "answerStatus")
		res.AnswerType = intFromMap(data, "answerType")
		res.AnswerUseTime = int64(intFromMap(data, "answerUseTime"))
		res.AnswerQuestionCount = intFromMap(data, "answerQuestionCount")
		res.AnswerCommitTime = int64(intFromMap(data, "answerCommitTime"))
		res.StatsStatus = intFromMap(data, "statsStatus")
	}
	return res, nil
}

// FetchQuestionnaireUploadURL fetches a presigned upload URL;
// upload via PUT with a Content-MD5 header carrying the file's MD5.
func (c *LansengerClient) FetchQuestionnaireUploadURL(ctx context.Context, fileName, md5 string, size int64, userToken string) (*QuestionnaireUploadUrlResult, error) {
	token, err := c.GetToken(ctx)
	if err != nil {
		return nil, err
	}
	url := BuildAPIURL(c.config, "questionnaires", "upload_url", token, WithUserToken(userToken))
	result, err := c.doPost(ctx, url, map[string]interface{}{"fileName": fileName, "md5": md5, "size": size})
	if err != nil {
		return &QuestionnaireUploadUrlResult{Success: false, Error: err.Error()}, nil
	}
	res := &QuestionnaireUploadUrlResult{Success: true, RawResponse: result}
	if data := extractData(result); data != nil {
		res.URL = strOrNil(data, "data")
	}
	return res, nil
}
