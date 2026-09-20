package lansenger

import (
	"context"
	"testing"
)

func TestReserveBoardroomSuccess(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/boardroom/server/openapi/v2/reserveRoom", 0, "ok", map[string]interface{}{
			"id": "res1", "reserveCode": "BR001", "boardRoomName": "第一会议室",
			"name": "周会", "status": "5", "reserveTime": "1小时",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.ReserveBoardroom(context.Background(), &BoardroomReserveParams{
		BoardRoomID:      "room1",
		Name:             "周会",
		GradingID:        "g1",
		ReserveTimeStart: "2026-07-22 09:00:00",
		ReserveTimeEnd:   "2026-07-22 10:00:00",
		NoticeTime:       "会前15分钟",
		PeopleNumber:     "10",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.ReserveCode != "BR001" || result.Status != "5" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestReserveBoardroomAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/boardroom/server/openapi/v2/reserveRoom", 54000, "会议室已被占用", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.ReserveBoardroom(context.Background(), &BoardroomReserveParams{
		BoardRoomID: "room1", Name: "n", GradingID: "g",
		ReserveTimeStart: "s", ReserveTimeEnd: "e", NoticeTime: "会前15分钟",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false on API error")
	}
}

func TestCancelAndConfirmSign(t *testing.T) {
	for _, ep := range []struct {
		name string
		path string
		run  func(c *LansengerClient) (*BoardroomOpResult, error)
	}{
		{"cancel", "reserveCancel", func(c *LansengerClient) (*BoardroomOpResult, error) {
			yes := true
			return c.CancelBoardroomReserve(context.Background(), "res1", "", "", "改期", &yes, nil, "", "", "")
		}},
		{"confirm-sign", "confirmSign", func(c *LansengerClient) (*BoardroomOpResult, error) {
			return c.ConfirmBoardroomSign(context.Background(), "res1", "", "")
		}},
	} {
		server := newMuxBuilder().
			handleToken("tok1").
			handle("/xtra/boardroom/server/openapi/v2/"+ep.path, 0, "ok", map[string]interface{}{"data": true}).
			build()
		c := newTestClient(server)
		result, err := ep.run(c)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", ep.name, err)
		}
		if !result.Success || !result.Done {
			t.Errorf("%s: expected Success=true Done=true", ep.name)
		}
		server.Close()
	}
}

func TestFetchBoardroomListPageInfo(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/boardroom/server/openapi/v2/roomList", 0, "ok", map[string]interface{}{
			"count": 1,
			"data":  []interface{}{map[string]interface{}{"id": "room1", "name": "第一会议室", "peopleNum": 20}},
			"code":  0, "msg": "",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchBoardroomList(context.Background(), "g1", "", nil, nil, "", "", "2026-07-22", 1, 10, "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if result.Count != 1 || len(result.Items) != 1 || result.Items[0]["name"] != "第一会议室" {
		t.Errorf("unexpected page: %+v", result)
	}
}

func TestFetchBoardroomSchedule(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/boardroom/server/openapi/v2/roomSchedule", 0, "ok", map[string]interface{}{
			"id": "room1", "name": "第一会议室", "canReserveFlag": "1",
			"reserveDtoList":       []interface{}{map[string]interface{}{"name": "周会"}},
			"deactivatedInfoList":  []interface{}{},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchBoardroomSchedule(context.Background(), "room1", "2026-07-22", "g1", "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got error=%s", result.Error)
	}
	if len(result.Reserves) != 1 || result.Reserves[0]["name"] != "周会" {
		t.Errorf("unexpected reserves: %+v", result.Reserves)
	}
}

func TestFetchBoardroomGradings(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/xtra/boardroom/server/openapi/v2/gradingList", mockJSONArrayHandler(0, "ok", []interface{}{
		map[string]interface{}{"id": "g1", "name": "默认分级", "type": "GRADING_ADMIN"},
	}))
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchBoardroomGradings(context.Background(), "", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.Total != 1 || result.Gradings[0]["id"] != "g1" {
		t.Errorf("unexpected gradings: %+v", result.Gradings)
	}
}

func TestFetchBoardroomDetail(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/xtra/boardroom/server/openapi/v2/roomDetail", 0, "ok", map[string]interface{}{
			"id": "room1", "name": "第一会议室", "status": "1", "canReserveFlag": "1",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchBoardroomDetail(context.Background(), "room1", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success || result.Name != "第一会议室" || result.CanReserveFlag != "1" {
		t.Errorf("unexpected detail: %+v", result)
	}
}
