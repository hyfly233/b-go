package main

import (
	"testing"
)

func TestNewYarn(t *testing.T) {
	yarn := NewYarn("test-yarn")
	if yarn.Name != "test-yarn" {
		t.Errorf("Expected yarn name to be 'test-yarn', got '%s'", yarn.Name)
	}
}

func TestYarnInitialize(t *testing.T) {
	yarn := NewYarn("test-yarn")
	yarn.Initialize()
}
