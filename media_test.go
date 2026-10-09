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

// newMediaDownloadServer serves the app token plus any /v1/medias/{id}/fetch
// request with a fixed body/content-type.
func newMediaDownloadServer(body []byte, contentType string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/apptoken/create", mockAppTokenHandler("tok1"))
	mux.HandleFunc("/v1/medias/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", contentType)
		w.Write(body)
	})
	return httptest.NewServer(mux)
}

func TestDownloadMediaToFileResultSuccess(t *testing.T) {
	server := newMediaDownloadServer([]byte("hello-bytes"), "application/octet-stream")
	defer server.Close()

	target := filepath.Join(t.TempDir(), "out.dat")
	c := newTestClient(server)

	result, err := c.DownloadMediaToFileResult(context.Background(), "media1", target)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Success {
		t.Fatalf("expected Success=true, got false (error: %s)", result.Error)
	}
	if result.Path != target {
		t.Errorf("expected Path=%s, got %s", target, result.Path)
	}
	if result.Size != len("hello-bytes") {
		t.Errorf("expected Size=%d, got %d", len("hello-bytes"), result.Size)
	}
	written, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(written) != "hello-bytes" {
		t.Errorf("expected file content hello-bytes, got %q", string(written))
	}
}

func TestDownloadMediaToFileResultAPIErrorIsReported(t *testing.T) {
	server := newMediaDownloadServer([]byte(`{"errCode":50000,"errMsg":"OPEN ID invalid"}`), "application/json")
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DownloadMediaToFileResult(context.Background(), "media_bad", filepath.Join(t.TempDir(), "never.dat"))
	if err != nil {
		t.Fatalf("expected no transport error, got %v", err)
	}
	if result.Success {
		t.Fatalf("expected Success=false")
	}
	if result.Path != "" {
		t.Errorf("expected empty Path, got %s", result.Path)
	}
	if !strings.Contains(result.Error, "50000") {
		t.Errorf("expected error to mention errCode 50000, got %s", result.Error)
	}
}

func TestDownloadMediaToFileResultWriteFailureIsReported(t *testing.T) {
	server := newMediaDownloadServer([]byte("payload"), "application/octet-stream")
	defer server.Close()

	c := newTestClient(server)
	result, err := c.DownloadMediaToFileResult(context.Background(), "media1", "/nonexistent-dir-xyz/out.dat")
	if err != nil {
		t.Fatalf("expected no transport error, got %v", err)
	}
	if result.Success {
		t.Fatalf("expected Success=false")
	}
	if !strings.Contains(result.Error, "writing file") {
		t.Errorf("expected write failure message, got %s", result.Error)
	}
}

func TestDownloadMediaToFileStillReturnsErrorOnAPIFailure(t *testing.T) {
	server := newMediaDownloadServer([]byte(`{"errCode":50000,"errMsg":"OPEN ID invalid"}`), "application/json")
	defer server.Close()

	c := newTestClient(server)
	if _, err := c.DownloadMediaToFile(context.Background(), "media_bad", filepath.Join(t.TempDir(), "never.dat")); err == nil {
		t.Fatalf("expected an error from DownloadMediaToFile")
	} else if !strings.Contains(err.Error(), "Download failed") {
		t.Errorf("expected Download failed error, got %v", err)
	}
}
