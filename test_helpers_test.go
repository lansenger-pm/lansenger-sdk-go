package lansenger

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func newTestClient(server *httptest.Server) *LansengerClient {
	cfg := NewConfig("test_app", "test_secret")
	if server != nil {
		cfg.APIGatewayURL = server.URL
	}
	c := NewClientWithConfig(cfg)
	if server != nil {
		c.httpClient = server.Client()
		c.tokenMgr = NewTokenManager(cfg, c.httpClient)
	}
	return c
}

func mockAppTokenHandler(token string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": 0,
			"errMsg":  "ok",
			"data": map[string]interface{}{
				"appToken":  token,
				"expiresIn": 7200,
			},
		})
	}
}

func mockAPIResponseHandler(errCode int, errMsg string, data map[string]interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errCode": errCode,
			"errMsg":  errMsg,
			"data":    data,
		})
	}
}

func mockErrorHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

type muxBuilder struct {
	mux *http.ServeMux
}

func newMuxBuilder() *muxBuilder {
	return &muxBuilder{mux: http.NewServeMux()}
}

func (b *muxBuilder) handleToken(token string) *muxBuilder {
	b.mux.HandleFunc("/v1/apptoken/create", mockAppTokenHandler(token))
	return b
}

func (b *muxBuilder) handle(path string, errCode int, errMsg string, data map[string]interface{}) *muxBuilder {
	b.mux.HandleFunc(path, mockAPIResponseHandler(errCode, errMsg, data))
	return b
}

func (b *muxBuilder) handleError(path string) *muxBuilder {
	b.mux.HandleFunc(path, mockErrorHandler())
	return b
}

func (b *muxBuilder) build() *httptest.Server {
	return httptest.NewServer(b.mux)
}
