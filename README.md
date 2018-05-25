# K8s Events

[![Build Status](https://travis-ci.org/justcompile/k8sevents.svg?branch=master)](https://travis-ci.org/justcompile/k8sevents)

Kubernetes (K8s) has a Pod Lifecycle Hook mechanism, but the details provided are exceptionally limited; in the case of using webhooks, it pings the specified endpoint with no data.

k8sevents is a daemon written in GoLang which attaches to the Docker socket running on a K8s Node, and utilises Docker's Events API to listen for specific events & then transmits data to a specified endpoint.

## Configuration
k8sevents takes a command line argument of `-config <path>` which is a JSON file in the following structure.

```
{
	"Docker": {
	    APIVersion string
	    CertPath   string
	    Host       string
	    TLSVerify  string
	    UseEnv     bool
  },
	"DispatchEndpoint": string
	"Events"           []string
	"Namespaces"       []string
}
```

By default, it lists for container Start & Stop events, from any pod running within any namespace except "kube-system", "default" and "kube-public".

The minimum required config would be:
```
{
    "Docker": {
      "UseEnv": true
    },
    "DispatchEndpoint": "http://my-endpoint.com/path/to/webhook"
}
```


