package logging

/**
TODO: add event logger - IN_PROGRESS
	* logging to file and stdout with defined log format (parsable, maybe used for timetravel debug)
*/

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

var (
	loggerCache map[string]*AppLogger
)

type AppLogger struct {
	*log.Logger
	name string
}

type EventLogger struct {
	appLogger *AppLogger
	format    string
}

func init() {
	loggerCache = map[string]*AppLogger{}
}

func GetLoggerFor(pkg string) *AppLogger {
	if logger, ok := loggerCache[pkg]; ok {
		return logger
	} else {
		l := &AppLogger{
			Logger: log.New(os.Stdout, pkg+": ", log.Ldate|log.Ltime|log.Lshortfile),
			name:   pkg,
		}
		loggerCache[pkg] = l
		return l
	}
}

func GetEventLoggerFor(eventSubject string) *EventLogger {
	prefix := fmt.Sprintf("subject[%s]", eventSubject)
	res := GetLoggerFor(prefix)
	res.SetFlags(log.Ldate | log.Ltime)
	file, err := os.OpenFile("event.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(errors.New("unable to read/create event log file for event logger: " + prefix))
		return nil
	}

	res.Output(os.Stdout, file)
	return &EventLogger{
		appLogger: &AppLogger{
			Logger: res.Logger,
			name:   prefix,
		},
	}
}

func (r *EventLogger) LogEvent(name string, args ...any) {
	if len(r.format) == 0 {
		r.appLogger.Printf(fmt.Sprintf(name, args...))
	} else {
		r.appLogger.Printf(fmt.Sprintf(r.format, fmt.Sprintf(name, args...)))
	}
}

func (r *EventLogger) LogPair(closure func() error, name string, args ...any) error {
	r.LogEvent(name+"status[started]", args)
	err := closure()
	if err != nil {
		r.LogEvent(name+"status[failed]", args)
	} else {
		r.LogEvent(name+"status[finished]", args)
	}
	return err
}

func (r *EventLogger) Format(format string) *EventLogger {
	r.format = format
	return r
}

// Output overrides output set for this logger
func (r *AppLogger) Output(out ...io.Writer) *AppLogger {
	logWriters := io.MultiWriter(out...)
	r.SetOutput(logWriters)
	return r
}
