package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/justcompile/k8sevents/internal/config"
	"github.com/justcompile/k8sevents/internal/exts"
)

var tenses = map[string]string{"start": "started", "stop": "stopped"}

// Handler struct...
type Handler struct {
	Config config.Configuration
}

// Dispatch ...
func (handler *Handler) Dispatch(event *docker.APIEvents) {
	containerAttrs := event.Actor.Attributes
	if containerAttrs["io.kubernetes.docker.type"] != "podsandbox" {
		return
	}

	payload := map[string]interface{}{}
	payload["event_type"] = fmt.Sprintf("POD_%s", strings.ToUpper(event.Action))
	payload["cluster_id"] = containerAttrs["it.justcompile.aaas.cluster_id"]
	payload["data"] = containerAttrs
	payload["description"] = fmt.Sprintf("%s %s", strings.Title(tenses[event.Action]), containerAttrs["role"])

	jsonValue, _ := json.Marshal(payload)

	fmt.Printf("[%s] Sending -> %v", handler.Config.DispatchEndpoint, payload)

	resp, err := http.Post(handler.Config.DispatchEndpoint, "application/json", bytes.NewBuffer(jsonValue))

	if err != nil {
		fmt.Errorf("Error: %v", err)
	}

	fmt.Printf("Recv -> %v", resp)
}

// Listen ...
func (handler *Handler) Listen(client *docker.Client) {
	listener := make(chan *docker.APIEvents)
	err := client.AddEventListener(listener)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Connected")

	defer func() {

		err = client.RemoveEventListener(listener)
		if err != nil {
			log.Fatal(err)
		}

	}()

	for {
		select {
		case msg := <-listener:
			isTracked, _ := exts.InArray(msg.Action, handler.Config.Events)
			if msg.Type == "container" && isTracked {
				namespace := msg.Actor.Attributes["io.kubernetes.pod.namespace"]

				isSystemNamespace, _ := exts.InArray(namespace, handler.Config.Namespaces)
				if !isSystemNamespace {
					handler.Dispatch(msg)
				}
			}
		}
	}
}
