package config

import "github.com/tkanos/gonfig"

var TrackedEvents = []string{"start", "stop"}
var SystemNamespaces = []string{"kube-system", "default", "kube-public"}

// DockerConfig ...
type DockerConfig struct {
	APIVersion string
	CertPath   string
	Host       string
	TLSVerify  string
	UseEnv     bool
}

// Configuration struct
type Configuration struct {
	Docker           *DockerConfig
	DispatchEndpoint string
	Events           []string
	Namespaces       []string
}

var config Configuration
var loaded = false

// Config ...
func Config() Configuration {
	if !loaded {
		config := Configuration{Events: TrackedEvents, Namespaces: SystemNamespaces}
		err := gonfig.GetConf("~/.k8sevents/config.json", &config)
		if err != nil {
			panic(err)
		}
		loaded = true
	}
	return config
}
