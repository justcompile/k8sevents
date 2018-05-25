package main

import (
	"flag"
	"os"

	"github.com/fsouza/go-dockerclient"
	cfg "github.com/justcompile/k8sevents/internal/config"
	"github.com/justcompile/k8sevents/internal/events"
)

func main() {
	configFilePtr := flag.String("config", "", "Path to configuration file. (Required)")
	flag.Parse()

	conf := cfg.Config(*configFilePtr)
	var client *docker.Client

	if !conf.Docker.UseEnv {
		os.Setenv("DOCKER_TLS_VERIFY", conf.Docker.TLSVerify)
		os.Setenv("DOCKER_HOST", conf.Docker.Host)
		os.Setenv("DOCKER_CERT_PATH", conf.Docker.CertPath)
		os.Setenv("DOCKER_API_VERSION", conf.Docker.APIVersion)
	}

	client, err := docker.NewClientFromEnv()
	if err != nil {
		panic(err)
	}

	dispatcher := &events.Handler{Config: *conf}
	dispatcher.Listen(client)
}
