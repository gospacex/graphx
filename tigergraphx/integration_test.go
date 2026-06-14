// integration_test
//go:build integration

package tigergraphx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gospacex/graphx"
)

func TestIntegration_TigerLifecycle(t *testing.T) {
	// Simulate TigerGraph REST API with httptest.
	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		if body["user"] == "" || body["password"] == "" {
			http.Error(w, `{"message":"Missing credentials"}`, http.StatusUnauthorized)
			return
		}
		json.NewEncoder(w).Encode(loginResponse{Token: "test-token-12345"})
	})
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"version": "4.2.0"})
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	t.Logf("mock TigerGraph at %s", srv.URL)

	// Test New + auth.
	cfg := graphx.Config{
		Address:  srv.Listener.Addr().String(),
		Username: "admin",
		Password: "tigergraph",
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)

	db, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer db.Close(ctx)

	if db.Token() == "" {
		t.Fatal("expected non-empty token after auth")
	}
	if db.HTTPClient() == nil {
		t.Fatal("expected non-nil HTTP client")
	}

	// Test Ping.
	if err := db.Ping(ctx); err != nil {
		t.Fatalf("Ping: %v", err)
	}

	// Test Close (idempotent).
	if err := db.Close(ctx); err != nil {
		t.Fatalf("first Close: %v", err)
	}
	if err := db.Close(ctx); err != nil {
		t.Fatalf("second Close (idempotent): %v", err)
	}

	// Test Get singleton.
	db2, err := Get(ctx, graphx.Config{
		Name:     "tiger-int",
		Address:  srv.Listener.Addr().String(),
		Username: "admin",
		Password: "tigergraph",
	})
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if db2.Token() == "" {
		t.Fatal("expected non-empty token from singleton")
	}
	Reset()
}

func TestIntegration_TigerPingFailsAfterClose(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(loginResponse{Token: "t"})
	})
	mux.HandleFunc("/echo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := httptest.NewServer(mux)
	t.Cleanup(srv.Close)

	cfg := graphx.Config{
		Address:  srv.Listener.Addr().String(),
		Username: "admin",
		Password: "pw",
	}
	ctx := context.Background()
	db, err := New(ctx, cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	db.Close(ctx)
	if err := db.Ping(ctx); err == nil {
		t.Fatal("expected error pinging after close")
	}
}
