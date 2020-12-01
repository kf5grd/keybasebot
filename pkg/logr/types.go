package logr

import "io"

// Logger holds information necessary to write log output
type Logger struct {
	Writer      io.Writer // Where to write log messages
	EnableDebug bool      // Whether to write debug messages
	JSON        bool      // Whether to write messages in JSON format
}

// New returns a new Logger
func New(writer io.Writer, debug, json bool) *Logger {
	return &Logger{
		Writer:      writer,
		EnableDebug: debug,
		JSON:        json,
	}
}

// Level represents a LogLevel (Info, Error, Debug, etc)
type Level int

// These constants represent the various known Levels
const (
	LevelUnknown Level = iota
	LevelDebug
	LevelInfo
	LevelError
)

// levelMap allows for a lookup of a Level's string representation
var levelMap = map[Level]string{
	LevelUnknown: "UNKNOWN",
	LevelDebug:   "DEBUG",
	LevelInfo:    "INFO",
	LevelError:   "ERROR",
}

// String returns a string representation of a Level
func (l Level) String() string {
	if s, ok := levelMap[l]; ok {
		return s
	}
	return levelMap[0]
}

// Msg holds information about a particular log message
type Msg struct {
	Time     int64  `json:"time"`
	FuncName string `json:"func_name"`
	Level    string `json:"level"`
	Message  string `json:"message"`
}
