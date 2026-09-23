package lansenger

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
)

func vcTestMembers() []VideoconferenceMember {
	return []VideoconferenceMember{
		{StaffID: "1001", EmployeeName: "Host", Role: VCMemberRoleHost},
		{StaffID: "1002", EmployeeName: "Alice", Role: VCMemberRoleParticipant},
	}
}

func TestVideoconferenceCreateMeetingSuccess(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/videoconference/openapi/v1/meeting/create", 0, "ok", map[string]interface{}{
			"id": "3173", "subject": "周会", "meetingNumber": "MN001",
			"startTime": float64(1735660800000), "type": float64(0), "status": float64(1),
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	res, err := c.CreateMeeting(context.Background(), &VideoconferenceCreateParams{
		Subject:   "周会",
		StartTime: 1735660800000,
		Members:   vcTestMembers(),
		OrgID:     "2285568",
		Type:      0,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !res.Success {
		t.Fatalf("expected Success=true, got error=%s", res.Error)
	}
	if res.Mid != "3173" || res.MeetingNumber != "MN001" || res.Status != 1 {
		t.Errorf("unexpected result: %+v", res)
	}
}

// Done 的语义：外层 errCode 已判定成功，内层 code 只有在存在时才作为子状态判据。
// modify 返回的是会议对象（无内层 code），此前被误判为未完成。
func TestFillVCOpDone(t *testing.T) {
	cases := []struct {
		name string
		data map[string]interface{}
		want bool
	}{
		{
			name: "inner code 0",
			data: map[string]interface{}{"code": float64(0)},
			want: true,
		},
		{
			name: "inner code non-zero",
			data: map[string]interface{}{"code": float64(105213), "message": "会议未开始或已结束"},
			want: false,
		},
		{
			// 2026-09-23 对 /meeting/modify 的实测抓包形状
			name: "meeting object without inner code",
			data: map[string]interface{}{
				"admin": "1001", "autoRecord": float64(0), "createSource": float64(1),
				"id": float64(1380079), "startTime": float64(1790233200000),
				"status": float64(0), "stopTime": float64(0),
				"subject": "测试预约 sdkvfy01 已改", "type": float64(1),
			},
			want: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res := &VideoconferenceOpResult{Success: true}
			fillVCOp(res, map[string]interface{}{"errCode": float64(0), "errMsg": "OK", "data": tc.data})
			if res.Done != tc.want {
				t.Errorf("want Done=%v, got %v", tc.want, res.Done)
			}
		})
	}
}

// 成功但完全没有 data 负载：三个 SDK 必须一致地判为「已完成」
// （Python/TypeScript 的 _op 同样返回 done=true）。
func TestFillVCOpDoneNoPayload(t *testing.T) {
	res := &VideoconferenceOpResult{Success: true}
	fillVCOp(res, map[string]interface{}{"errCode": float64(0), "errMsg": "OK"})
	if !res.Done {
		t.Errorf("no payload on a successful response should be Done=true, got false")
	}
	if res.Message != "" {
		t.Errorf("want empty message, got %q", res.Message)
	}
}

// 2026-09-23 对 /meeting/events/subscribe 的实测抓包：内层带 code，
// 仍按 code == 0 判定（与改动前一致，不受本次 done 语义调整影响）。
func TestFillVCOpDoneSubscribeShape(t *testing.T) {
	res := &VideoconferenceOpResult{Success: true}
	fillVCOp(res, map[string]interface{}{
		"errCode": float64(0), "errMsg": "OK",
		"data": map[string]interface{}{
			"code": float64(0), "errCode": float64(0), "message": "",
		},
	})
	if !res.Done {
		t.Errorf("subscribe (inner code 0) should be Done=true")
	}
}

// opCode 不做客户端校验：未知值也原样送到服务端，由服务端裁决。
func TestVideoconferenceControlMemberPassesOpCodeThrough(t *testing.T) {
	var body map[string]interface{}

	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/videoconference/openapi/v1/meeting/member/control", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		mockAPIResponseHandler(0, "ok", map[string]interface{}{"code": float64(0)})(w, r)
	})
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	if _, err := c.ControlMember(context.Background(), "1", "2", "some_new_op", "op", "1", ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body["opCode"] != "some_new_op" {
		t.Errorf("opCode must be forwarded verbatim, got %v", body["opCode"])
	}
}

// modify 与 create 同契约：显式传 userStopTime 要出现在请求体里，
// 不传则必须整体省略（服务端会据此把结束时间重置为开始时间 +24 小时）。
func TestVideoconferenceModifyForwardsUserStopTime(t *testing.T) {
	var body map[string]interface{}

	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/videoconference/openapi/v1/meeting/modify", func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
		}
		mockAPIResponseHandler(0, "ok", map[string]interface{}{"code": float64(0)})(w, r)
	})
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	base := func() *VideoconferenceModifyParams {
		return &VideoconferenceModifyParams{
			Mid: "3173", Subject: "new", StartTime: 1735660800000,
			Members: vcTestMembers(), OrgID: "2285568", Operator: "1001",
		}
	}

	stop := int64(1700003600000)
	p := base()
	p.UserStopTime = &stop
	if _, err := c.ModifyMeeting(context.Background(), p); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := body["userStopTime"]; got != float64(stop) {
		t.Errorf("want userStopTime=%d in body, got %v", stop, got)
	}

	body = nil
	if _, err := c.ModifyMeeting(context.Background(), base()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := body["userStopTime"]; ok {
		t.Errorf("userStopTime must be omitted when unset, got body=%v", body)
	}
}

func TestVideoconferenceCreateMeetingValidation(t *testing.T) {
	cases := []struct {
		name string
		p    *VideoconferenceCreateParams
		want string
	}{
		{"subject", &VideoconferenceCreateParams{StartTime: 1, Members: vcTestMembers()}, "subject is required"},
		{"no members", &VideoconferenceCreateParams{Subject: "s", StartTime: 1}, "member is required (with exactly one host)"},
		{"two admins", &VideoconferenceCreateParams{Subject: "s", StartTime: 1, Members: []VideoconferenceMember{
			{StaffID: "a", Role: VCMemberRoleHost}, {StaffID: "b", Role: VCMemberRoleHost},
		}}, "exactly one member must have role='admin'"},
		{"zero admins", &VideoconferenceCreateParams{Subject: "s", StartTime: 1, Members: []VideoconferenceMember{
			{StaffID: "a", Role: VCMemberRoleParticipant},
		}}, "exactly one member must have role='admin'"},
		{"start time", &VideoconferenceCreateParams{Subject: "s", Members: vcTestMembers()}, "start_time is required"},
	}
	c := newTestClient(nil)
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := c.CreateMeeting(context.Background(), tc.p)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if res.Success || res.Error != tc.want {
				t.Errorf("want error=%q, got Success=%v Error=%q", tc.want, res.Success, res.Error)
			}
		})
	}
}

func TestVideoconferenceMeetingLifecycle(t *testing.T) {
	for _, ep := range []struct {
		path string
		call func(c *LansengerClient) (*VideoconferenceOpResult, error)
	}{
		{
			path: "/xtra/videoconference/openapi/v1/meeting/modify",
			call: func(c *LansengerClient) (*VideoconferenceOpResult, error) {
				return c.ModifyMeeting(context.Background(), &VideoconferenceModifyParams{
					Mid: "3173", Subject: "new", StartTime: 1735660800000,
					Members: vcTestMembers(), OrgID: "2285568", Operator: "1001",
				})
			},
		},
		{
			path: "/xtra/videoconference/openapi/v1/meeting/cancle",
			call: func(c *LansengerClient) (*VideoconferenceOpResult, error) {
				return c.CancelMeeting(context.Background(), "3173", "2285568", "1001", "")
			},
		},
		{
			path: "/xtra/videoconference/openapi/v1/meeting/stop",
			call: func(c *LansengerClient) (*VideoconferenceOpResult, error) {
				return c.StopMeeting(context.Background(), "3173", "2285568", "1001", "")
			},
		},
		{
			path: "/xtra/videoconference/openapi/v1/meeting/events/subscribe",
			call: func(c *LansengerClient) (*VideoconferenceOpResult, error) {
				return c.SubscribeMeetingEvents(context.Background(), "3173", "2285568",
					[]map[string]interface{}{{"event": "meetingEnd"}}, "", "")
			},
		},
		{
			path: "/xtra/videoconference/openapi/v1/meeting/member/control",
			call: func(c *LansengerClient) (*VideoconferenceOpResult, error) {
				return c.ControlMember(context.Background(), "3173", "1002", "muteall", "1001", "2285568", "")
			},
		},
		{
			path: "/xtra/videoconference/openapi/v1/meeting/member/invite",
			call: func(c *LansengerClient) (*VideoconferenceOpResult, error) {
				return c.InviteMembers(context.Background(), "MN001",
					[]VideoconferenceMember{{StaffID: "1003", EmployeeName: "Bob"}}, "2285568", "1001", "")
			},
		},
	} {
		t.Run(ep.path, func(t *testing.T) {
			server := newMuxBuilder().
				handleToken("tok1").
				handle(ep.path, 0, "ok", map[string]interface{}{
					"data": map[string]interface{}{"code": float64(0), "message": "done"},
				}).
				build()
			defer server.Close()
			c := newTestClient(server)
			res, err := ep.call(c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Success || !res.Done {
				t.Errorf("expected Success+Done, got %+v", res)
			}
		})
	}
}

func TestVideoconferencePagedListEndpoints(t *testing.T) {
	pageData := map[string]interface{}{
		"offset": float64(0), "total": float64(2),
		"items": []interface{}{
			map[string]interface{}{"id": "m1"},
			map[string]interface{}{"id": "m2"},
		},
	}
	for _, ep := range []struct {
		path string
		call func(c *LansengerClient) (*VideoconferenceListResult, error)
	}{
		{"/xtra/videoconference/openapi/v1/meeting/list", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchMeetingList(context.Background(), "2285568", 1000, 2000, VCFetchRangeAll, "", 10, 0, "")
		}},
		{"/xtra/videoconference/openapi/v1/meeting/record/list", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchMeetingRecordList(context.Background(), "2285568", 1000, 2000, "", VCCreateSourceClient, 10, 0, "")
		}},
		{"/xtra/videoconference/openapi/v1/meeting/member/simplerecord", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchMemberSimpleRecord(context.Background(), "3173", "2285568", "1001", "", 10, 0)
		}},
		{"/xtra/videoconference/openapi/v1/meeting/fixroom/list", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchFixroomList(context.Background(), "2285568", "1001", "", 10, 0)
		}},
		{"/xtra/videoconference/openapi/v1/meeting/history/fetch", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchHistoryMeetings(context.Background(), "2285568", "1001", "", 10, 0)
		}},
		{"/xtra/videoconference/openapi/v1/meeting/active/fetch", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchActiveMeetings(context.Background(), "2285568", "1001", "", 10, 0)
		}},
		{"/xtra/videoconference/openapi/v1/meeting/member/list", func(c *LansengerClient) (*VideoconferenceListResult, error) {
			return c.FetchMemberList(context.Background(), "3173", "2285568", "1001", "", 10, 0)
		}},
	} {
		t.Run(ep.path, func(t *testing.T) {
			server := newMuxBuilder().
				handleToken("tok1").
				handle(ep.path, 0, "ok", pageData).
				build()
			defer server.Close()
			c := newTestClient(server)
			res, err := ep.call(c)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !res.Success || res.Total != 2 || len(res.Items) != 2 {
				t.Errorf("unexpected result: %+v", res)
			}
		})
	}
}

func TestVideoconferenceMiscEndpoints(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/videoconference/openapi/v1/meeting/detail", 0, "ok", map[string]interface{}{
			"id": "3173", "subject": "周会", "status": float64(4), "admin": "1001",
		}).
		handle("/xtra/videoconference/openapi/v1/meeting/status/fetchmore", 0, "ok", map[string]interface{}{
			"mids": []interface{}{
				map[string]interface{}{"mid": "3173", "status": float64(1)},
			},
		}).
		handle("/xtra/videoconference/openapi/v1/meeting/param/fetch", 0, "ok", map[string]interface{}{
			"meetingInfo": map[string]interface{}{"subject": "周会"},
		}).
		handle("/xtra/videoconference/openapi/v1/meeting/vod/list", 0, "ok", map[string]interface{}{
			"items": []interface{}{map[string]interface{}{"vodId": "v1"}},
		}).
		handle("/xtra/videoconference/openapi/v1/vod/url/download/fetch", 0, "ok", map[string]interface{}{
			"v1": "https://download.example/v1.mp4",
		}).
		handle("/xtra/videoconference/openapi/v1/conf/fetch", 0, "ok", map[string]interface{}{
			"maxPerson": float64(300), "defaultMaxPerson": float64(100),
			"allowedRecordFlag": float64(1), "forcePasswdFlag": float64(0), "spaceSize": float64(10240),
		}).
		build()
	defer server.Close()

	c := newTestClient(server)

	detail, err := c.FetchMeetingDetail(context.Background(), "3173", "2285568", "1001", "")
	if err != nil || !detail.Success || detail.Mid != "3173" || detail.Status != 4 {
		t.Errorf("detail failed: %+v err=%v", detail, err)
	}

	status, err := c.FetchMeetingStatus(context.Background(), []string{"3173"}, "2285568", "")
	if err != nil || !status.Success || len(status.Statuses) != 1 {
		t.Errorf("status failed: %+v err=%v", status, err)
	}

	params, err := c.FetchMeetingParams(context.Background(), "MN001", "2285568", "1001", "")
	if err != nil || !params.Success || params.MeetingInfo["subject"] != "周会" {
		t.Errorf("params failed: %+v err=%v", params, err)
	}

	vods, err := c.FetchVodList(context.Background(), "3173", "2285568", "1001", "")
	if err != nil || !vods.Success || len(vods.Items) != 1 {
		t.Errorf("vod list failed: %+v err=%v", vods, err)
	}

	urls, err := c.FetchVodDownloadURLs(context.Background(),
		[]VideoconferenceVod{{VodID: "v1"}}, "2285568", "1001", "")
	if err != nil || !urls.Success || urls.Data["v1"] == nil {
		t.Errorf("vod urls failed: %+v err=%v", urls, err)
	}

	conf, err := c.FetchOrgConf(context.Background(), "2285568", "", "", "")
	if err != nil || !conf.Success || conf.MaxPerson != 300 || conf.SpaceSize != 10240 {
		t.Errorf("conf failed: %+v err=%v", conf, err)
	}
}

func TestVideoconferenceValidationGuards(t *testing.T) {
	c := newTestClient(nil)
	ctx := context.Background()

	if res, _ := c.FetchMeetingList(ctx, "1", 0, 0, VCFetchRangePerson, "", 10, 0, ""); res.Success || res.Error != "staff_id is required when fetch_range='person'" {
		t.Errorf("fetchRange person guard: %+v", res)
	}
	if res, _ := c.FetchMeetingStatus(ctx, nil, "1", ""); res.Success || res.Error != "mids is required" {
		t.Errorf("mids guard: %+v", res)
	}
	// opCode 不再做客户端校验（服务端才是权威），已从本组 guard 移除，
	// 透传行为由 TestVideoconferenceControlMemberPassesOpCodeThrough 覆盖。
	if res, _ := c.FetchVodDownloadURLs(ctx, nil, "1", "op", ""); res.Success || res.Error != "vods must contain 1..3 entries" {
		t.Errorf("vods empty guard: %+v", res)
	}
	many := make([]VideoconferenceVod, 4)
	if res, _ := c.FetchVodDownloadURLs(ctx, many, "1", "op", ""); res.Success || res.Error != "vods must contain 1..3 entries" {
		t.Errorf("vods >3 guard: %+v", res)
	}
	if res, _ := c.FetchMeetingParams(ctx, "", "1", "op", ""); res.Success || res.Error != "meeting_number is required" {
		t.Errorf("meetingNumber guard: %+v", res)
	}
	if res, _ := c.CancelMeeting(ctx, "", "1", "op", ""); res.Success || res.Error != "mid is required" {
		t.Errorf("mid guard: %+v", res)
	}
}
