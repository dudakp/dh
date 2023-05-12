package logging

/**
TODO: remove [] in event log message, in event[action_started] parameters:
	EVENT: 2023/05/10 23:10:12 subject[flow(TExecuteEffectFlow)]|event[action_started([0])]
*/

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"
)

type EventLogFormatElement string

const (
	EventLoggerSubjectFormat EventLogFormatElement = "subject[{{.Subject}}]"
	EventLoggerEventFormat   EventLogFormatElement = "event[{{.Event}}]"

	eventLoggerPrefix = "EVENT: "
)

var (
	loggerCache map[string]*AppLogger
)

// AppLogger base logger type
type AppLogger struct {
	*log.Logger
	name string
}

// EventLogger provides specific functionality for creating eventlog
type EventLogger struct {
	appLogger *AppLogger
	format    []EventLogFormatElement
	subject   string
	template  *template.Template
}

func init() {
	loggerCache = map[string]*AppLogger{}
}

func GetLoggerFor(name string) *AppLogger {
	if logger, ok := loggerCache[name]; ok {
		return logger
	} else {
		l := &AppLogger{
			Logger: log.New(os.Stdout, name+": ", log.Ldate|log.Ltime|log.Lshortfile),
			name:   name,
		}
		loggerCache[name] = l
		return l
	}
}

func (r *AppLogger) Output(out ...io.Writer) *AppLogger {
	logWriters := io.MultiWriter(out...)
	r.SetOutput(logWriters)
	return r
}

func GetEventLoggerFor(name string) *EventLogger {
	formattedName := name + "E"
	res := GetLoggerFor(formattedName)
	res.SetPrefix(eventLoggerPrefix)
	res.SetFlags(log.Ldate | log.Ltime)
	tmpl := template.New("eventLogTemplate")
	file, err := os.OpenFile("event.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(errors.New("unable to read/create event log file for event logger: " + formattedName))
		return nil
	}

	res.Output(os.Stdout, file)
	return &EventLogger{
		template: tmpl,
		appLogger: &AppLogger{
			Logger: res.Logger,
			name:   formattedName,
		},
	}
}

func (r *EventLogger) LogEvent(eventName string, args ...any) {
	if len(r.format) == 0 {
		r.appLogger.Printf(fmt.Sprintf(eventName, args...))
	} else {
		buff := &bytes.Buffer{}
		err := r.template.Execute(buff, struct {
			Subject string
			Event   string
		}{r.subject, eventName})
		if err != nil {
			panic(err)
		}
		r.appLogger.Printf(buff.String(), args)
	}
}

func (r *EventLogger) Output(out ...io.Writer) *EventLogger {
	r.appLogger.Output(out...)
	return r
}

func (r *EventLogger) Format(elements ...EventLogFormatElement) *EventLogger {
	r.format = elements

	var elementFormatStrings []string
	for _, element := range r.format {
		elementFormatStrings = append(elementFormatStrings, string(element))
	}
	format := strings.Join(elementFormatStrings, "|")

	var err error
	r.template, err = r.template.Parse(format)
	if err != nil {
		panic(err)
	}
	return r
}

func (r *EventLogger) Subject(subject string) *EventLogger {
	r.subject = subject
	return r
}
