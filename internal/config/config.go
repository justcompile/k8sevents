package config

import (
	"github.com/tkanos/gonfig"
)

var trackedEvents = []string{"start", "stop"}
var systemNamespaces = []string{"kube-system", "default", "kube-public"}

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
	Docker           DockerConfig
	DispatchEndpoint string
	Events           []string
	Namespaces       []string
}

var config *Configuration

// Config ...
func Config(configPath string) *Configuration {
	config := Configuration{Events: trackedEvents, Namespaces: systemNamespaces}

	err := gonfig.GetConf(configPath, &config)

	if err != nil {
		panic(err)
	}

	return &config
}
