package deployment

import (
	"github.com/vite-cloud/vite/core/domain/config"
	"gotest.tools/v3/assert"
	"testing"
)

func TestEvent_IsError(t *testing.T) {
	event := Event{ID: ErrorEvent}
	assert.Assert(t, event.IsError())

	event = Event{}
	assert.Assert(t, !event.IsError())
}

func TestEvent_IsGlobal(t *testing.T) {
	event := Event{}
	assert.Assert(t, event.IsGlobal())

	event = Event{Service: &config.Service{}}
	assert.Assert(t, !event.IsGlobal())
}

func TestEvent_Label(t *testing.T) {
	event := Event{}
	assert.Assert(t, event.Label() == "global")

	event = Event{Service: &config.Service{Name: "test"}}
	assert.Assert(t, event.Label() == "test")
}
