package log

import "testing"

func TestLevel_String(t *testing.T) {
	tests := []struct {
		level level
		want  string
	}{
		{DebugLevel, "debug"},
		{InfoLevel, "info"},
		{WarnLevel, "warn"},
		{ErrorLevel, "error"},
		{FatalLevel, "fatal"},
		{level(15), "unknown"},
		{level(-1), "unknown"},
	}

	for _, tt := range tests {
		if got := tt.level.String(); got != tt.want {
			t.Errorf("%v.Marshal() = %v, want %v", tt.level, got, tt.want)
		}
	}
}
