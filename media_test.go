package lansenger

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestUploadAppMediaV2(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/app/medias/create", 0, "ok", map[string]interface{}{
			"mediaId": "media001",
		}).
		build()
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "test.jpg")
	if err := os.WriteFile(tmpFile, []byte("fake-image-data"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	c := newTestClient(server)
	result, err := c.UploadAppMediaV2(context.Background(), tmpFile, "image", "utok1", 100, 200, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v (error: %s)", result.Success, result.Error)
	}
	if result.MediaID != "media001" {
		t.Errorf("expected MediaID=media001, got %s", result.MediaID)
	}
}

func TestUploadAppMediaV2APIError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v2/app/medias/create", 10001, "invalid media type", nil).
		build()
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "test.jpg")
	if err := os.WriteFile(tmpFile, []byte("fake-image-data"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	c := newTestClient(server)
	result, err := c.UploadAppMediaV2(context.Background(), tmpFile, "image", "utok1", 100, 200, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for API error")
	}
	if !strings.Contains(result.Error, "invalid media type") {
		t.Errorf("expected error to contain 'invalid media type', got %s", result.Error)
	}
}

func TestUploadAppMediaV2UserTokenInURL(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/v2/app/medias/create", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("user_token"); got != "utok1" {
			t.Errorf("expected user_token=utok1 in URL, got %q", got)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 0,
			"errMsg":  "ok",
			"data":    map[string]interface{}{"mediaId": "media002"},
		})
	})
	server := b.build()
	defer server.Close()

	tmpFile := filepath.Join(t.TempDir(), "test.jpg")
	if err := os.WriteFile(tmpFile, []byte("fake-image-data"), 0644); err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}

	c := newTestClient(server)
	result, err := c.UploadAppMediaV2(context.Background(), tmpFile, "image", "utok1", 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v (error: %s)", result.Success, result.Error)
	}
}

func TestDownloadMediaByShareID(t *testing.T) {
	binaryData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/v1/media/share/share001/fetch", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(binaryData)
	})
	server := b.build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DownloadMediaByShareID(context.Background(), "share001", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v (error: %s)", result.Success, result.Error)
	}
	if string(result.Data) != string(binaryData) {
		t.Errorf("expected Data=%v, got %v", binaryData, result.Data)
	}
}

func TestDownloadMediaByShareIDJSONError(t *testing.T) {
	server := newMuxBuilder().
		handleToken("tok1").
		handle("/v1/media/share/share404/fetch", 10404, "share not found", nil).
		build()
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DownloadMediaByShareID(context.Background(), "share404", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Success {
		t.Error("expected Success=false for JSON error response")
	}
	if !strings.Contains(result.Error, "share not found") {
		t.Errorf("expected error to contain 'share not found', got %s", result.Error)
	}
}

func TestDownloadMediaByShareIDUserTokenInURL(t *testing.T) {
	b := newMuxBuilder().handleToken("tok1")
	b.mux.HandleFunc("/v1/media/share/share002/fetch", func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("user_token"); got != "utok1" {
			t.Errorf("expected user_token=utok1 in URL, got %q", got)
		}
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write([]byte("data"))
	})
	server := httptest.NewServer(b.mux)
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DownloadMediaByShareID(context.Background(), "share002", "utok1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Errorf("expected Success=true, got %v (error: %s)", result.Success, result.Error)
	}
}
