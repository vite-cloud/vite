package deployment

import "github.com/vite-cloud/vite/core/domain/config"

const (
	ErrorEvent  = "Error"
	FinishEvent = "Finish"
)

type Event struct {
	// Service is the name of the service that the event is about.
	Service *config.Service
	// ID is an identifier unique to the kind of the event.
	ID string
	// Data is the payload of the event.
	Data any
}

// IsError returns true if the event is an error event.
func (e Event) IsError() bool {
	return e.ID == ErrorEvent
}

// IsGlobal returns true if the event is not related to a service.
func (e Event) IsGlobal() bool {
	return e.Service == nil
}

// Label returns the type of the event either 'global' or the service name.
func (e Event) Label() string {
	if e.Service == nil {
		return "global"
	}

	return e.Service.Name
}

func (e Event) IsFinish() bool {
	return e.ID == FinishEvent
}
