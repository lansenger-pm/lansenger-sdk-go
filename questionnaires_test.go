package lansenger

import (
	"context"
	"testing"
)

func TestSaveQuestionnaireSuccess(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/saveQuestionnaire", 0, "ok", "QN001").
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SaveQuestionnaire(context.Background(), &QuestionnaireSaveParams{
		Title: "满意度调查", AccountCode: "ACC001", WelcomeSpeech: "欢迎",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.QuestionnaireCode != "QN001" {
		t.Errorf("expected QuestionnaireCode=QN001, got %s", result.QuestionnaireCode)
	}
}

func TestSaveQuestionnaireAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/saveQuestionnaire", 3104, "官方账号不存在！", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SaveQuestionnaire(context.Background(), &QuestionnaireSaveParams{
		Title: "t", AccountCode: "BAD",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false on API error")
	}
	if result.Error == "" {
		t.Error("expected error message to be set")
	}
}

func TestSaveQuestionnaireQuestionsScalarResponse(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/saveQuestionList", 0, "ok", 2).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.SaveQuestionnaireQuestions(
		context.Background(),
		"QN001",
		[]map[string]interface{}{{"questionName": "Q1", "questionType": "radio"}},
		"",
		"",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.SavedCount != 2 {
		t.Fatalf("expected Success=true SavedCount=2, got %+v", result)
	}
}

func TestDeleteQuestionnaireQuestionScalarResponse(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/deleteQuestion", 0, "ok", true).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeleteQuestionnaireQuestion(context.Background(), "Q1", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || !result.Deleted {
		t.Fatalf("expected Success=true Deleted=true, got %+v", result)
	}
}

func TestQuestionnaireStringScalarResponses(t *testing.T) {
	answerServer := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/getAnswerUrl", 0, "ok", "https://qn.example.com/Q1").
		build()
	defer answerServer.Close()

	answerClient := newTestClient(answerServer)
	answerResult, err := answerClient.FetchQuestionnaireAnswerURL(context.Background(), "Q1", "", "")
	if err != nil {
		t.Fatalf("answer URL: unexpected error: %v", err)
	}
	if !answerResult.Success || answerResult.URL != "https://qn.example.com/Q1" {
		t.Fatalf("expected answer URL, got %+v", answerResult)
	}

	copyServer := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/copy", 0, "ok", "QN002").
		build()
	defer copyServer.Close()

	copyClient := newTestClient(copyServer)
	copyResult, err := copyClient.CopyQuestionnaire(context.Background(), "QN001", "", "")
	if err != nil {
		t.Fatalf("copy: unexpected error: %v", err)
	}
	if !copyResult.Success || copyResult.NewCode != "QN002" {
		t.Fatalf("expected new code QN002, got %+v", copyResult)
	}
}

func TestPublishQuestionnaireDefaults(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/publish", 0, "ok", true).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.PublishQuestionnaire(context.Background(), &QuestionnairePublishParams{
		QuestionnaireCode: "QN001",
		ScopeType:         QuestionnaireScopeInternal,
		AnswerLimit:       QuestionnaireAnswerLimitOnce,
		ViewStatsFlag:     1,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || !result.Done {
		t.Errorf("expected Success=true Done=true, got %v/%v (error=%s)", result.Success, result.Done, result.Error)
	}
}

func TestWithdrawFinishDeleteCodeBody(t *testing.T) {
	for _, ep := range []struct{ name, path string }{
		{"withdraw", "withdraw"}, {"finish", "finish"}, {"delete", "delete"},
	} {
		server := newMuxBuilder().
			handleToken("tok1").
			handle("/xtra/questionnaire/server/openapi/v1/"+ep.path, 0, "ok", true).
			build()
		c := newTestClient(server)
		var result *QuestionnaireOpResult
		var err error
		switch ep.name {
		case "withdraw":
			result, err = c.WithdrawQuestionnaire(context.Background(), "QN1", "U1", "")
		case "finish":
			result, err = c.FinishQuestionnaire(context.Background(), "QN1", "U1", "")
		case "delete":
			result, err = c.DeleteQuestionnaire(context.Background(), "QN1", "U1", "")
		}
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", ep.name, err)
		}
		if !result.Success || !result.Done {
			t.Errorf("%s: expected Success=true Done=true", ep.name)
		}
		server.Close()
	}
}

func TestFetchQuestionnaireDetail(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/detail", 0, "ok", map[string]interface{}{
			"id": 1001, "code": "QN001", "title": "t", "status": 2,
			"questionCount": 1,
			"questionList":  []interface{}{map[string]interface{}{"code": "Q1", "questionType": "radio"}},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchQuestionnaireDetail(context.Background(), "QN001", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.QuestionnaireID != 1001 || result.Status != 2 || result.QuestionCount != 1 {
		t.Errorf("unexpected detail: %+v", result)
	}
	if len(result.Questions) != 1 || result.Questions[0]["code"] != "Q1" {
		t.Errorf("unexpected questions: %+v", result.Questions)
	}
}

func TestFetchCreatedQuestionnairesPage(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/createList", 0, "ok", map[string]interface{}{
			"pageNo": 2, "pageSize": 10, "pages": 3, "total": 25, "hasNextPage": true,
			"result": []interface{}{map[string]interface{}{"code": "QN1", "title": "t"}},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchCreatedQuestionnaires(context.Background(), "ACC001", 2, 10, nil, "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.Total != 25 || result.PageNo != 2 || !result.HasMore {
		t.Errorf("unexpected page: %+v", result)
	}
	if len(result.Items) != 1 || result.Items[0]["code"] != "QN1" {
		t.Errorf("unexpected items: %+v", result.Items)
	}
}

func TestFetchQuestionnaireOfficeAccounts(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/questionnaire/server/openapi/v1/userOfficeAccountList", mockJSONArrayHandler(0, "ok", []interface{}{
		map[string]interface{}{"code": "ACC001", "roleName": "人事部"},
	}))
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchQuestionnaireOfficeAccounts(context.Background(), "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.Total != 1 || result.Accounts[0]["code"] != "ACC001" {
		t.Errorf("unexpected accounts: %+v (error=%s)", result.Accounts, result.Error)
	}
}

func TestFetchQuestionnaireAnswerDetail(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/answerDetail", 0, "ok", map[string]interface{}{
			"answerUserName": "张三",
			"questionnaire":  map[string]interface{}{"code": "QN1"},
			"answerMap":      map[string]interface{}{"Q1": map[string]interface{}{"context": "非常满意"}},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchQuestionnaireAnswerDetail(context.Background(), "ACC001", "AR1", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.AnswerUserName != "张三" {
		t.Errorf("expected AnswerUserName=张三, got %s", result.AnswerUserName)
	}
	if result.Answers["Q1"].(map[string]interface{})["context"] != "非常满意" {
		t.Errorf("unexpected answers: %+v", result.Answers)
	}
}

func TestFetchQuestionnaireUploadURL(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/questionnaire/server/openapi/v1/upload", 0, "ok", "https://oss.example.com/u?sign=x").
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchQuestionnaireUploadURL(context.Background(), "a.png", "md5", 1024, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.URL == "" {
		t.Errorf("expected Success=true with URL, got %v (error=%s)", result.Success, result.Error)
	}
}
