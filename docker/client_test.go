package docker

import (
	"os"
	"testing"
)

func TestLoadDocker(t *testing.T) {
	old := os.Getenv("DOCKER_HOST")
	os.Setenv("DOCKER_HOST", "gibberish")

	_, err := newDocker()
	if err == nil {
		t.Error("expected error, got nil")
	}

	os.Setenv("DOCKER_HOST", old)
	_, err = newDocker()
	if err != nil {
		t.Errorf("expected no error, got %s", err)
	}
}