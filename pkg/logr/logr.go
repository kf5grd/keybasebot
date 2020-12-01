package logr

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// getFrame and getCaller taken from https://stackoverflow.com/questions/35212985/is-it-possible-get-information-about-caller-function-in-golang
func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

func getCaller() string {
	// Skip GetCallerFunctionName and the function to get the caller of
	frame := strings.Split(getFrame(2).Function, ".")
	return frame[len(frame)-1]
}

// Write writes the log string to the Logger.Writer with the given name string, and Level, and returns a string with the same output in case you want to do something else with it
func (l *Logger) Write(name string, level Level, s string, a ...interface{}) string {
	var ret string
	var (
		now        = time.Now().UTC()
		timeFormat = "02Jan2006 15:04:05"
	)

	if l.JSON {
		msg := Msg{
			Time:     now.Unix(),
			FuncName: name,
			Level:    level.String(),
			Message:  fmt.Sprintf(s, a...),
		}
		ret = toJson(msg) + "\n"
		fmt.Fprint(l.Writer, ret)
		return ret
	}

	a = append([]interface{}{strings.ToUpper(now.Format(timeFormat)), name, level}, a...)
	ret = fmt.Sprintf("[%v][%s] %v: "+s+"\n", a...)
	fmt.Fprint(l.Writer, ret)
	return ret
}

// Error sets the Level to LevelError, and automatically sets the name of the caller, then calls Write
func (l *Logger) Error(s string, a ...interface{}) string {
	return l.Write(getCaller(), LevelError, s, a...)
}

// Info sets the Level to LevelInfo, and automatically sets the name of the caller, then calls Write
func (l *Logger) Info(s string, a ...interface{}) string {
	return l.Write(getCaller(), LevelInfo, s, a...)
}

// Debug sets the Level to LevelDebug, and automatically sets the name of the caller, then calls Write only if Logger.EnableDebug is true
func (l *Logger) Debug(s string, a ...interface{}) string {
	if !l.EnableDebug {
		return ""
	}
	return l.Write(getCaller(), LevelDebug, s, a...)
}
