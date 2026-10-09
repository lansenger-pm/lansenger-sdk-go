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
