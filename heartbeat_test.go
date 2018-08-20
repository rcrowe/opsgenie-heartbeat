package heartbeat_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"

	"github.com/gorilla/mux"
	heartbeat "github.com/rcrowe/opsgenie-heartbeat"
	"github.com/sethgrid/pester"
)

func TestKeySet(t *testing.T) {
	expected := "some random key"

	os.Setenv(heartbeat.EnvAPIKey, expected)
	defer os.Setenv(heartbeat.EnvAPIKey, "")

	hb := heartbeat.New("foo")
	if hb.APIKey != expected {
		t.Fail()
	}
}

func TestKeyHeader(t *testing.T) {
	expected := "5e0cf9c8-4665-4b7b-b829-6044b9323d98"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "GenieKey "+expected {
			t.Fail()
		}
	}))
	defer srv.Close()

	hb := heartbeat.New("foo")
	hb.APIKey = expected
	hb.Endpoint = srv.URL
	hb.Ping(context.Background())
}

func TestCorrectEndpoint(t *testing.T) {
	expected := "the-heartbeat-name"

	m := mux.NewRouter()
	m.HandleFunc("/v2/heartbeats/{name}/ping", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		if vars["name"] != expected {
			t.Fail()
		}
	}).Methods("GET")
	srv := httptest.NewServer(m)
	defer srv.Close()

	hb := heartbeat.New(expected)
	hb.Endpoint = srv.URL
	hb.Ping(context.Background())
}

func TestCombatsFlakyConnection(t *testing.T) {
	var calls uint32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint32(&calls, 1)
		if c := atomic.LoadUint32(&calls); c < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	hb := heartbeat.New("some-heartbeat")
	hb.Endpoint = srv.URL

	if err := hb.Ping(context.Background()); err != nil {
		t.Fail()
	}
}

func TestUnauthorised(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer srv.Close()

	hb := heartbeat.New("foo")
	hb.Endpoint = srv.URL

	if err := hb.Ping(context.Background()); err != heartbeat.ErrUnauthorised {
		t.Fail()
	}
}

func TestBadResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer srv.Close()

	hb := heartbeat.New("foo")
	hb.Endpoint = srv.URL

	client := pester.New()
	client.Concurrency = 1
	client.MaxRetries = 1
	hb.Client = client

	if err := hb.Ping(context.Background()); err != heartbeat.ErrNonOkStatusCode {
		t.Fail()
	}
}
