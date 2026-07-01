package lansenger

import (
	"context"
	"testing"
)

func TestCreatePersonalApp(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/personal/apps/create", 0, "ok", map[string]interface{}{
			"id":           "app1",
			"secret":       "sec1",
			"apigwAddr":    "https://gw",
			"passportAddr": "https://pp",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreatePersonalApp(context.Background(), "utok1", "MyApp", "", "desc")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.AppID != "app1" {
		t.Errorf("expected AppID=app1, got %s", result.AppID)
	}
	if result.Secret != "sec1" {
		t.Errorf("expected Secret=sec1, got %s", result.Secret)
	}
}

func TestCreatePersonalAppAPIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/personal/apps/create", 10005, "no permission", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.CreatePersonalApp(context.Background(), "utok1", "MyApp", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
}

func TestUpdatePersonalApp(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/personal/apps/app1/update", 0, "ok", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.UpdatePersonalApp(context.Background(), "app1", "utok1", "NewName", "", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchPersonalApp(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/personal/apps/app1/fetch", 0, "ok", map[string]interface{}{
			"name":         "MyApp",
			"description":  "desc",
			"apigwAddr":    "https://gw",
			"passportAddr": "https://pp",
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchPersonalApp(context.Background(), "app1", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if result.Name != "MyApp" {
		t.Errorf("expected Name=MyApp, got %s", result.Name)
	}
}

func TestDeletePersonalApp(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/personal/apps/app1/delete", 0, "ok", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DeletePersonalApp(context.Background(), "app1", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
}

func TestFetchPersonalAppList(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/personal/apps/list/fetch", 0, "ok", map[string]interface{}{
			"appList": []interface{}{
				map[string]interface{}{"appId": "a1", "appName": "App1", "description": "d1"},
			},
		}).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.FetchPersonalAppList(context.Background(), "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v", result.Success)
	}
	if len(result.AppList) != 1 {
		t.Errorf("expected 1 app, got %d", len(result.AppList))
	}
}
