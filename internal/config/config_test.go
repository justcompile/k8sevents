package config

import (
	"bufio"
	"fmt"
	"os"
	"testing"
)

func setupSubTest(t *testing.T, filePath string, data string) func(t *testing.T) {
	f, _ := os.Create(filePath)
	defer f.Close()

	w := bufio.NewWriter(f)
	w.WriteString(data)
	w.Flush()

	t.Log("setup sub test")
	return func(t *testing.T) {
		os.Remove(filePath)
		t.Log("teardown sub test")
	}
}

type assertFn func(conf *Configuration) bool

func TestConfig(t *testing.T) {
	cases := []struct {
		name     string
		data     string
		expected assertFn
	}{
		{"Empty Config", "{}\n", func(config *Configuration) bool {
			return config.DispatchEndpoint == ""
		}},
		{"Dispatch Endpoint", "{\"DispatchEndpoint\": \"http://endpoint.com/\"}", func(config *Configuration) bool { return config.DispatchEndpoint == "http://endpoint.com/" }},
		// {"zero", 0, 0, 0},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			const filePath = "test.json"

			teardownSubTest := setupSubTest(t, filePath, tc.data)
			defer teardownSubTest(t)

			var cfg = *Config(filePath)

			if !tc.expected(&cfg) {
				fmt.Printf("%+v\n", config)
				t.Errorf("Error in Config.DispatchEndpoint => %s", tc.name)
			}
		})
	}
}
