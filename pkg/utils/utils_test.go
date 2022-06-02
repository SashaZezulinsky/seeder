package utils

import (
	"fmt"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	var tests = []struct {
		path string
		res  string
	}{
		{"docker", "./config/config-docker"},
		{"", "./config/config-local"},
		{"test", "./config/config-local"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("path %s", tt.path)
		t.Run(testname, func(t *testing.T) {
			ans := GetConfigPath(tt.path)
			if ans != tt.res {
				t.Errorf("got %s, want %s", ans, tt.res)
			}
		})
	}
}
