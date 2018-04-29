package events

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/justcompile/k8sevents/internal/config"
)

func setupSubTest(t *testing.T, handlerFunc http.HandlerFunc) (func(t *testing.T), string) {
	ts := httptest.NewServer(handlerFunc)
	defer ts.Close()

	t.Log("setup sub test")
	return func(t *testing.T) {
		t.Log("teardown sub test")
	}, ts.URL
}

func Test_Dispatch_Does_Not_Call_Endpoint_When_Container_Attr_Is_Not_Pod_Sandbox(t *testing.T) {
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		called = true
	})

	teardownSubTest, dispatchURL := setupSubTest(t, handler)
	defer teardownSubTest(t)

	var cfg = &config.Configuration{DispatchEndpoint: dispatchURL}

	var eventHandler = &Handler{Config: cfg}

	eventHandler.Dispatch(nil)

	if called {
		t.Errorf("Dispatch unexpectidly called HTTP Endpoint")
	}
}
