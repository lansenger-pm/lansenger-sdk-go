package lansenger

import (
	"context"
	"net/http"
	"strings"
	"testing"
)

type domainCall func(*LansengerClient) (bool, string, error)

type domainContractCase struct {
	name    string
	path    string
	success http.HandlerFunc
	call    domainCall
}

func runDomainCall(t *testing.T, path string, handler http.HandlerFunc, call domainCall) (bool, string) {
	t.Helper()
	handlerServer := newMuxBuilder()
	handlerServer.mux.HandleFunc("/v1/apptoken/create", mockAppTokenHandler("tok1"))
	handlerServer.mux.HandleFunc(path, handler)
	server := handlerServer.build()
	defer server.Close()

	c := newTestClient(server)
	success, resultErr, err := call(c)
	if err != nil {
		return false, err.Error()
	}
	return success, resultErr
}

func expectDomainSuccess(t *testing.T, tc domainContractCase) {
	t.Helper()
	success, errMsg := runDomainCall(t, tc.path, tc.success, tc.call)
	if !success {
		t.Fatalf("expected success, got error=%s", errMsg)
	}
}

func expectDomainAPIError(t *testing.T, tc domainContractCase) {
	t.Helper()
	success, errMsg := runDomainCall(t, tc.path, mockAPIResponseHandler(50001, "domain error", nil), tc.call)
	if success {
		t.Fatal("expected API error result")
	}
	if errMsg == "" {
		t.Fatal("expected API error message")
	}
}

func expectDomainHTTPError(t *testing.T, tc domainContractCase) {
	t.Helper()
	success, errMsg := runDomainCall(t, tc.path, mockErrorHandler(), tc.call)
	if success {
		t.Fatal("expected HTTP error result")
	}
	if !strings.Contains(errMsg, "HTTP 500") {
		t.Fatalf("expected HTTP error, got %q", errMsg)
	}
}

func mapSuccess(data map[string]interface{}) http.HandlerFunc {
	return mockAPIResponseHandler(0, "ok", data)
}

func arraySuccess(data []interface{}) http.HandlerFunc {
	return mockJSONArrayHandler(0, "ok", data)
}

func domainContractCases() []domainContractCase {
	return []domainContractCase{
		{
			name: "notice/send",
			path: "/xtra/notice/server/openapi/v1/send",
			success: mapSuccess(map[string]interface{}{
				"code": "NTC001", "noticeStatus": 2,
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.SendNotice(context.Background(), &NoticeSendParams{
					Title: "t", ContentType: 1, AccountCode: "ACC001", UserType: 1,
					Content: "c", ReleasePhones: []string{"13800138000"}, CreateMobile: "13800138000",
				})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "notice/accounts",
			path: "/xtra/notice/server/openapi/v1/notice/account",
			success: arraySuccess([]interface{}{
				map[string]interface{}{"code": "ACC001"},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchNoticeAccounts(context.Background(), "org1", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/save",
			path:    "/xtra/questionnaire/server/openapi/v1/saveQuestionnaire",
			success: mapSuccess(map[string]interface{}{"data": "QN1"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.SaveQuestionnaire(context.Background(), &QuestionnaireSaveParams{Title: "t", AccountCode: "ACC001"})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/save-questions",
			path:    "/xtra/questionnaire/server/openapi/v1/saveQuestionList",
			success: mapSuccess(map[string]interface{}{"data": 1}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.SaveQuestionnaireQuestions(context.Background(), "QN1", []map[string]interface{}{{"questionName": "Q1"}}, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/delete-question",
			path:    "/xtra/questionnaire/server/openapi/v1/deleteQuestion",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.DeleteQuestionnaireQuestion(context.Background(), "Q1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/publish",
			path:    "/xtra/questionnaire/server/openapi/v1/publish",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.PublishQuestionnaire(context.Background(), &QuestionnairePublishParams{
					QuestionnaireCode: "QN1", ScopeType: 1, AnswerLimit: 1, ViewStatsFlag: 1,
				})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/withdraw",
			path:    "/xtra/questionnaire/server/openapi/v1/withdraw",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.WithdrawQuestionnaire(context.Background(), "QN1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/finish",
			path:    "/xtra/questionnaire/server/openapi/v1/finish",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FinishQuestionnaire(context.Background(), "QN1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/delete",
			path:    "/xtra/questionnaire/server/openapi/v1/delete",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.DeleteQuestionnaire(context.Background(), "QN1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/detail",
			path: "/xtra/questionnaire/server/openapi/v1/detail",
			success: mapSuccess(map[string]interface{}{
				"id": 1001, "code": "QN1", "title": "t", "questionCount": 1,
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireDetail(context.Background(), "QN1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/brief",
			path:    "/xtra/questionnaire/server/openapi/v1/detailWithoutAuth",
			success: mapSuccess(map[string]interface{}{"code": "QN1", "title": "t"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireBrief(context.Background(), "QN1", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/answer-url",
			path:    "/xtra/questionnaire/server/openapi/v1/getAnswerUrl",
			success: mapSuccess(map[string]interface{}{"data": "https://example.com/answer"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireAnswerURL(context.Background(), "QN1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/copy",
			path:    "/xtra/questionnaire/server/openapi/v1/copy",
			success: mapSuccess(map[string]interface{}{"data": "QN2"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.CopyQuestionnaire(context.Background(), "QN1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/query-codes",
			path: "/xtra/questionnaire/server/openapi/v1/queryList",
			success: arraySuccess([]interface{}{
				map[string]interface{}{"code": "QN1"},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnairesByCodes(context.Background(), []string{"QN1"}, 0, "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/accounts",
			path: "/xtra/questionnaire/server/openapi/v1/userOfficeAccountList",
			success: arraySuccess([]interface{}{
				map[string]interface{}{"code": "ACC001"},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireOfficeAccounts(context.Background(), "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/created-list",
			path: "/xtra/questionnaire/server/openapi/v1/createList",
			success: mapSuccess(map[string]interface{}{
				"pageNo": 1, "pageSize": 10, "pages": 1, "total": 1, "hasNextPage": false,
				"result": []interface{}{map[string]interface{}{"code": "QN1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchCreatedQuestionnaires(context.Background(), "ACC001", 1, 10, nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/my-created",
			path: "/xtra/questionnaire/server/openapi/v1/myCreateList",
			success: mapSuccess(map[string]interface{}{
				"pageNo": 1, "pageSize": 10, "pages": 1, "total": 1, "hasNextPage": false,
				"result": []interface{}{map[string]interface{}{"code": "QN1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMyCreatedQuestionnaires(context.Background(), "org1", 1, 10, "", nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/participated",
			path: "/xtra/questionnaire/server/openapi/v1/participationList",
			success: mapSuccess(map[string]interface{}{
				"pageNo": 1, "pageSize": 10, "pages": 1, "total": 1, "hasNextPage": false,
				"result": []interface{}{map[string]interface{}{"code": "QN1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchParticipatedQuestionnaires(context.Background(), "org1", 1, 10, nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answers",
			path: "/xtra/questionnaire/server/openapi/v1/answerList",
			success: mapSuccess(map[string]interface{}{
				"pageNo": 1, "pageSize": 10, "pages": 1, "total": 1, "hasNextPage": false,
				"result": []interface{}{map[string]interface{}{"code": "AR1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchAnswerRecords(context.Background(), "ACC001", "QN1", 1, 10, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answer-detail",
			path: "/xtra/questionnaire/server/openapi/v1/answerDetail",
			success: mapSuccess(map[string]interface{}{
				"code": "AR1", "answerUserName": "张三",
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireAnswerDetail(context.Background(), "ACC001", "AR1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/last-answer-detail",
			path: "/xtra/questionnaire/server/openapi/v1/lastAnswerDetail",
			success: mapSuccess(map[string]interface{}{
				"code": "AR1", "answerUserName": "张三",
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireLastAnswerDetail(context.Background(), "QN1", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answer-data",
			path: "/xtra/questionnaire/server/openapi/v1/answerData",
			success: mapSuccess(map[string]interface{}{
				"pageNo": 1, "pageSize": 10, "pages": 1, "total": 1, "hasNextPage": false,
				"result": []interface{}{map[string]interface{}{"answerCode": "AR1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchAnswerData(context.Background(), "ACC001", "QN1", 1, 10, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/last-answer-record",
			path: "/xtra/questionnaire/server/openapi/v1/lastAnswerRecord",
			success: mapSuccess(map[string]interface{}{
				"code": "AR1", "answerUserName": "张三", "statsStatus": 1,
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireLastAnswerRecord(context.Background(), "QN1", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "questionnaire/upload-url",
			path:    "/xtra/questionnaire/server/openapi/v1/upload",
			success: mapSuccess(map[string]interface{}{"data": "https://example.com/upload"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireUploadURL(context.Background(), "a.png", "md5", 1024, "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/rooms",
			path: "/xtra/boardroom/server/openapi/v2/roomList",
			success: mapSuccess(map[string]interface{}{
				"count": 1, "data": []interface{}{map[string]interface{}{"id": "room1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomList(context.Background(), "g1", "", nil, nil, "", "", "", 1, 10, "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "boardroom/room-detail",
			path:    "/xtra/boardroom/server/openapi/v2/roomDetail",
			success: mapSuccess(map[string]interface{}{"id": "room1", "name": "第一会议室"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomDetail(context.Background(), "room1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/schedule",
			path: "/xtra/boardroom/server/openapi/v2/roomSchedule",
			success: mapSuccess(map[string]interface{}{
				"id": "room1", "reserveDtoList": []interface{}{map[string]interface{}{"id": "r1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomSchedule(context.Background(), "room1", "2026-07-22", "g1", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "boardroom/reserve-detail",
			path:    "/xtra/boardroom/server/openapi/v2/reserveDetail",
			success: mapSuccess(map[string]interface{}{"id": "res1", "name": "周会"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomReserveDetail(context.Background(), "res1", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "boardroom/reserve",
			path:    "/xtra/boardroom/server/openapi/v2/reserveRoom",
			success: mapSuccess(map[string]interface{}{"id": "res1", "reserveCode": "BR001"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.ReserveBoardroom(context.Background(), &BoardroomReserveParams{
					BoardRoomID: "room1", Name: "周会", GradingID: "g1",
					ReserveTimeStart: "2026-07-22 09:00:00",
					ReserveTimeEnd:   "2026-07-22 10:00:00",
					NoticeTime:       "会前15分钟",
				})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "boardroom/edit-reserve",
			path:    "/xtra/boardroom/server/openapi/v2/editReserve",
			success: mapSuccess(map[string]interface{}{"id": "res1", "reserveCode": "BR001"}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.EditBoardroomReserve(context.Background(), &BoardroomEditParams{
					ReserveID: "res1", BoardRoomID: "room1", Name: "周会", GradingID: "g1",
					ReserveTimeStart: "2026-07-22 09:00:00",
					ReserveTimeEnd:   "2026-07-22 10:00:00",
					NoticeTime:       "会前15分钟",
				})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "boardroom/cancel",
			path:    "/xtra/boardroom/server/openapi/v2/reserveCancel",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.CancelBoardroomReserve(context.Background(), "res1", "", "", "", nil, nil, "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name:    "boardroom/confirm-sign",
			path:    "/xtra/boardroom/server/openapi/v2/confirmSign",
			success: mapSuccess(map[string]interface{}{"data": true}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.ConfirmBoardroomSign(context.Background(), "res1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/my-reserves",
			path: "/xtra/boardroom/server/openapi/v2/myReserveList",
			success: mapSuccess(map[string]interface{}{
				"count": 1, "data": []interface{}{map[string]interface{}{"id": "res1"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMyBoardroomReserves(context.Background(), "g1", "", "", "", "", nil, 1, 10, "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/gradings",
			path: "/xtra/boardroom/server/openapi/v2/gradingList",
			success: arraySuccess([]interface{}{
				map[string]interface{}{"id": "g1"},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomGradings(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/area-offices",
			path: "/xtra/boardroom/server/openapi/v2/areaOfficeList",
			success: arraySuccess([]interface{}{
				map[string]interface{}{"id": "area1"},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomAreaOffices(context.Background(), "g1", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/create",
			path: "/xtra/videoconference/openapi/v1/meeting/create",
			success: mapSuccess(map[string]interface{}{
				"id": "3173", "subject": "周会", "meetingNumber": "MN001",
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.CreateMeeting(context.Background(), &VideoconferenceCreateParams{
					Subject: "周会", StartTime: 1735660800000,
					Members: []VideoconferenceMember{
						{StaffID: "1001", Role: VCMemberRoleHost},
						{StaffID: "1002", Role: VCMemberRoleParticipant},
					},
					OrgID: "2285568",
				})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/detail",
			path: "/xtra/videoconference/openapi/v1/meeting/detail",
			success: mapSuccess(map[string]interface{}{
				"id": "3173", "subject": "周会", "status": float64(4),
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMeetingDetail(context.Background(), "3173", "2285568", "1001", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/list",
			path: "/xtra/videoconference/openapi/v1/meeting/list",
			success: mapSuccess(map[string]interface{}{
				"offset": 0, "total": 1,
				"items": []interface{}{map[string]interface{}{"id": "3173"}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMeetingList(context.Background(), "2285568", 1000, 2000, VCFetchRangeAll, "", 10, 0, "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/status",
			path: "/xtra/videoconference/openapi/v1/meeting/status/fetchmore",
			success: mapSuccess(map[string]interface{}{
				"mids": []interface{}{map[string]interface{}{"mid": "3173", "status": float64(1)}},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMeetingStatus(context.Background(), []string{"3173"}, "2285568", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/member-control",
			path: "/xtra/videoconference/openapi/v1/meeting/member/control",
			success: mapSuccess(map[string]interface{}{
				"data": map[string]interface{}{"code": float64(0), "message": "ok"},
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.ControlMember(context.Background(), "3173", "1002", "muteall", "1001", "2285568", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/vod-download",
			path: "/xtra/videoconference/openapi/v1/vod/url/download/fetch",
			success: mapSuccess(map[string]interface{}{
				"v1": "https://download.example/v1.mp4",
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchVodDownloadURLs(context.Background(),
					[]VideoconferenceVod{{VodID: "v1"}}, "2285568", "1001", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/conf",
			path: "/xtra/videoconference/openapi/v1/conf/fetch",
			success: mapSuccess(map[string]interface{}{
				"maxPerson": float64(300), "defaultMaxPerson": float64(100),
			}),
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchOrgConf(context.Background(), "2285568", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
	}
}

func TestDomainEndpointContracts(t *testing.T) {
	for _, tc := range domainContractCases() {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Run("success", func(t *testing.T) { expectDomainSuccess(t, tc) })
			t.Run("api_error", func(t *testing.T) { expectDomainAPIError(t, tc) })
			t.Run("http_error", func(t *testing.T) { expectDomainHTTPError(t, tc) })
		})
	}
}

type guardCase struct {
	name string
	want string
	call domainCall
}

func TestDomainInputGuards(t *testing.T) {
	cases := []guardCase{
		{
			name: "notice/send/title",
			want: "title is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.SendNotice(context.Background(), &NoticeSendParams{})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/save/title",
			want: "title is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.SaveQuestionnaire(context.Background(), &QuestionnaireSaveParams{})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/save-questions/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.SaveQuestionnaireQuestions(context.Background(), "", nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/delete-question/code",
			want: "question_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.DeleteQuestionnaireQuestion(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/publish/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.PublishQuestionnaire(context.Background(), &QuestionnairePublishParams{})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/withdraw/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.WithdrawQuestionnaire(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/finish/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FinishQuestionnaire(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/delete/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.DeleteQuestionnaire(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/detail/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireDetail(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/brief/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireBrief(context.Background(), "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answer-url/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireAnswerURL(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/copy/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.CopyQuestionnaire(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/query-codes/list",
			want: "code_list is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnairesByCodes(context.Background(), nil, 0, "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/created-list/account",
			want: "account_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchCreatedQuestionnaires(context.Background(), "", 1, 10, nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/my-created/org",
			want: "org_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMyCreatedQuestionnaires(context.Background(), "", 1, 10, "", nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/participated/org",
			want: "org_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchParticipatedQuestionnaires(context.Background(), "", 1, 10, nil, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answers/account",
			want: "account_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchAnswerRecords(context.Background(), "", "QN1", 1, 10, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answer-detail/account",
			want: "account_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireAnswerDetail(context.Background(), "", "AR1", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/last-answer-detail/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireLastAnswerDetail(context.Background(), "", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/answer-data/account",
			want: "account_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchAnswerData(context.Background(), "", "QN1", 1, 10, "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/last-answer-record/code",
			want: "questionnaire_code is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireLastAnswerRecord(context.Background(), "", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "questionnaire/upload-url/file",
			want: "file_name is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchQuestionnaireUploadURL(context.Background(), "", "", 0, "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/room-detail/id",
			want: "room_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomDetail(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/schedule/id",
			want: "room_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomSchedule(context.Background(), "", "", "", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/reserve-detail/id",
			want: "reserve_room_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomReserveDetail(context.Background(), "", "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/reserve/room",
			want: "boardroom_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.ReserveBoardroom(context.Background(), &BoardroomReserveParams{})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/edit/reserve",
			want: "reserve_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.EditBoardroomReserve(context.Background(), &BoardroomEditParams{})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/cancel/id",
			want: "reserve_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.CancelBoardroomReserve(context.Background(), "", "", "", "", nil, nil, "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/confirm/id",
			want: "reserve_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.ConfirmBoardroomSign(context.Background(), "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/my-reserves/grading",
			want: "grading_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMyBoardroomReserves(context.Background(), "", "", "", "", "", nil, 1, 10, "", "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "boardroom/area-offices/grading",
			want: "grading_id is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchBoardroomAreaOffices(context.Background(), "", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/create/admin-count",
			want: "exactly one member must have role='admin'",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.CreateMeeting(context.Background(), &VideoconferenceCreateParams{
					Subject: "s", StartTime: 1,
					Members: []VideoconferenceMember{{StaffID: "a", Role: VCMemberRoleParticipant}},
				})
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/list/person-range",
			want: "staff_id is required when fetch_range='person'",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMeetingList(context.Background(), "1", 0, 0, VCFetchRangePerson, "", 10, 0, "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/status/mids",
			want: "mids is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMeetingStatus(context.Background(), nil, "1", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/member-control/opcode",
			want: "op_code must be one of ",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.ControlMember(context.Background(), "1", "2", "notAnOp", "op", "1", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/vod-download/count",
			want: "vods must contain 1..3 entries",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchVodDownloadURLs(context.Background(), nil, "1", "op", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
		{
			name: "videoconference/params/number",
			want: "meeting_number is required",
			call: func(c *LansengerClient) (bool, string, error) {
				r, err := c.FetchMeetingParams(context.Background(), "", "1", "op", "")
				if err != nil {
					return false, "", err
				}
				return r.Success, r.Error, nil
			},
		},
	}

	c := newTestClient(nil)
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			success, resultErr, err := tc.call(c)
			if success {
				t.Fatal("expected guard failure")
			}
			got := resultErr
			if err != nil {
				got = err.Error()
			}
			if !strings.Contains(got, tc.want) {
				t.Fatalf("expected %q, got %q", tc.want, got)
			}
		})
	}
}
