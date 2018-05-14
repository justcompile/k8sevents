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

func TestConfig_DispatchEndpoint(t *testing.T) {
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

func TestConfig_Panics_When_File_Does_Not_Exist(t *testing.T) {

	func() {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("TestConfig_Panics_When_File_Does_Not_Exist should have panicked!")
			}
		}()
		// This function should cause a panic
		Config("i/dont/exist.json")
	}()

}

func TestConfig_DockerObject(t *testing.T) {
	cases := []struct {
		name     string
		data     string
		expected assertFn
	}{
		{"Docker UseEnv", "{\"Docker\": {\"UseEnv\": true}}", func(config *Configuration) bool { return config.Docker.UseEnv == true }},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			const filePath = "test.json"

			teardownSubTest := setupSubTest(t, filePath, tc.data)
			defer teardownSubTest(t)

			var cfg = *Config(filePath)

			if !tc.expected(&cfg) {
				fmt.Printf("%+v\n", config)
				t.Errorf("Error in Config.Docker => %s", tc.name)
			}
		})
	}
}
