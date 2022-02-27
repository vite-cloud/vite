package log

// level is a log level from 0 to 4
type level int

// Log levels
const (
	DebugLevel level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// String returns the string representation of the log level
// It supports the following levels:
// - DebugLevel -> "debug"
// - InfoLevel -> "info"
// - WarnLevel -> "warn"
// - ErrorLevel -> "error"
// - FatalLevel -> "fatal"
// If the level is invalid, it returns "unknown"
func (l level) String() string {
	switch l {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	default:
		return "unknown"
	}
}
