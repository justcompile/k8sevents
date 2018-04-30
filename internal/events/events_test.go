package events

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"testing"

	simplejson "github.com/bitly/go-simplejson"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/justcompile/k8sevents/internal/config"
)

func setupSubTest(t *testing.T, handlerFunc *http.HandlerFunc) (func(t *testing.T), string) {
	ts := httptest.NewServer(handlerFunc)

	t.Log("setup sub test")
	return func(t *testing.T) {
		ts.Close()
		t.Log("teardown sub test")
	}, ts.URL
}

func makeAPIEvent(attrs map[string]string, action string) docker.APIEvents {
	return docker.APIEvents{Action: action, Actor: docker.APIActor{Attributes: attrs}}
}

func testStringMapEquality(t *testing.T, a map[string]interface{}, b map[string]string) bool {
	keys := make([]string, 0)
	for k := range a {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if a[k] != b[k] {
			t.Errorf("Keys %s do not match: %v  !=  %v", k, a[k], b[k])
			return false
		}
	}
	return true
}

func testMapEquality(t *testing.T, a map[string]interface{}, b map[string]interface{}) bool {
	keys := make([]string, 0)
	for k := range a {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		if reflect.ValueOf(a[k]).Kind() == reflect.Map {
			c := testStringMapEquality(t, a[k].(map[string]interface{}), b[k].(map[string]string))
			if !c {
				return false
			}
		} else if !reflect.DeepEqual(a[k], b[k]) {
			t.Errorf("Keys %s do not match: %v  !=  %v", k, a[k], b[k])
			return false
		}
	}

	return false
}

func Test_Dispatch_Does_Not_Call_Endpoint_When_Container_Attr_Is_Not_Pod_Sandbox(t *testing.T) {
	called := false

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		called = true
	})

	teardownSubTest, dispatchURL := setupSubTest(t, &handler)
	defer teardownSubTest(t)

	var cfg = &config.Configuration{DispatchEndpoint: dispatchURL}

	attributes := map[string]string{"io.kubernetes.docker.type": "POD"}

	var event = makeAPIEvent(attributes, "hi")

	var eventHandler = &Handler{Config: cfg}

	eventHandler.Dispatch(&event)

	if called {
		t.Errorf("Dispatch unexpectidly called HTTP Endpoint")
	}
}

func Test_Dispatch_Sends_Correct_Payload_Format_When_Action_Received(t *testing.T) {
	actualPayload := map[string]interface{}{}
	called := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		w.WriteHeader(http.StatusOK)
		if r.Method != "POST" {
			t.Errorf("Expected 'POST' request, got '%s'", r.Method)
		}

		reqJSON, _ := simplejson.NewFromReader(r.Body)
		jsonMap, _ := reqJSON.Map()

		for key, value := range jsonMap {
			actualPayload[key] = value
		}
	})

	teardownSubTest, dispatchURL := setupSubTest(t, &handler)
	defer teardownSubTest(t)

	var cfg = &config.Configuration{DispatchEndpoint: dispatchURL}

	attributes := map[string]string{
		"io.kubernetes.docker.type":      "podsandbox",
		"it.justcompile.aaas.cluster_id": "my-cluster-123",
		"role": "webserver",
	}

	var event = makeAPIEvent(attributes, "stop")

	var eventHandler = &Handler{Config: cfg}

	eventHandler.Dispatch(&event)

	expectedPayload := map[string]interface{}{
		"event_type":  "POD_STOP",
		"cluster_id":  "my-cluster-123",
		"data":        attributes,
		"description": "Stopped webserver",
	}

	if !called {
		t.Error("Endpoint was not called")
	}

	testMapEquality(t, actualPayload, expectedPayload)
}
