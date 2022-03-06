package cmd

import (
	"testing"
)

func TestNewVersionCommand(t *testing.T) {
	CommandTest{
		NewCommand: NewVersionCommand,
		Test: func(console *Expect) {
			console.
				String("Vite version dev, build unknown").
				EOF()
		},
	}.Run(t)
}
