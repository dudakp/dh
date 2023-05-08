package logging

/**
TODO: add event logger - IN_PROGRESS
	* logging to file and stdout with defined log format (parsable, maybe used for timetravel debug)
*/

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	loggerCache map[string]*Logger
)

type Logger struct {
	*log.Logger
	name string
}

func init() {
	loggerCache = map[string]*Logger{}
}

func GetLoggerFor(pkg string) *Logger {
	if logger, ok := loggerCache[pkg]; ok {
		return logger
	} else {
		l := &Logger{
			Logger: log.New(os.Stdout, pkg+": ", log.Ldate|log.Ltime|log.Lshortfile),
			name:   pkg,
		}
		loggerCache[pkg] = l
		return l
	}
}

func (r *Logger) Flags(flags int) *Logger {
	r.Logger.SetFlags(flags)
	return r
}

func (
	r *Logger) Prefix(prefix string) *Logger {
	r.SetPrefix(prefix)
	return r
}

// Output overrides output set for this logger
func (r *Logger) Output(out ...io.Writer) *Logger {
	logWriters := io.MultiWriter(out...)
	r.SetOutput(logWriters)
	return r
}

func (r *Logger) EventLogger(eventSubject string, eventName string) *Logger {
	r.SetPrefix(fmt.Sprintf("%s[%s]", eventSubject, eventName))
	r.SetFlags(log.Ldate | log.Ltime)

	return r
}
