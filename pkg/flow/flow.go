package flow

/**
TODO: create structured event log - outputted to console && file | IN_PROGRESS
	data recorded in eventLog: flow id, flow name, handler id, timestamp, handler data, action successful, onError called, onTerminalError called

TODO: create operators
	* pipe
	* map

TODO: create execution strategies
	* sequential | DONE
	* parallel | IN_PROGRESS
	* synchronized

TODO: create flow types:
	* TemplatedFlow
		flow definition in file, calling scripts (bash scripts)
		flow loaded from event log
	* ApiFlow
		flow calling REST API

TODO: create event log parser
	nice to have - not priority
*/

import (
	"dh/pkg/logging"
	"errors"
	"fmt"
)

const (
	actionOnError           = "onError(handlerId:%d)"
	actionExecutionStarted  = "action_started(handlerId:%d)"
	actionExecutionFinished = "action_finished(handlerId:%d)"

	subjectFormat = "flow(%s)"
)

var (
	HandlerMissingActionErr = errors.New("no action function specified for handler")

	logger          = logging.GetLoggerFor("flow")
	flowEventLogger = logging.GetEventLoggerFor("flow").
			Format(logging.EventLoggerSubjectFormat, logging.EventLoggerEventFormat)
)

type HandlerAction[T any] func(T) (T, error)
type HandlerErrorAction[T any] func(*Handler[T], error)

type Handler[T any] struct {
	action  HandlerAction[T]
	onError HandlerErrorAction[T]

	id       int
	parallel bool

	next *Handler[T]
	prev *Handler[T]
}

type Opts struct {
	// Name optional name for flow. If not specified, random name will be generated
	Name string
	// if ExecuteOnErrorAlways is true, every OnError function will be called in every handler event if terminalOnError is specified
	ExecuteOnErrorAlways bool
}

type flow[T any] struct {
	initialData     T
	opts            Opts
	terminalOnError func(err error)

	firstHandler            *Handler[T]
	terminalOnErrorExecuted bool
}

func NewHandler[T any](action HandlerAction[T], onError HandlerErrorAction[T]) *Handler[T] {
	var res = &Handler[T]{
		action:  action,
		onError: onError,
	}
	if res.onError == nil {
		logger.Printf("missing onError function for handler. creating default onError")
		res.onError = func(handler *Handler[T], err error) {
			flowEventLogger.LogEvent(actionOnError, handler.id)
		}
	}
	withEventLog(res)
	return res
}

func NewParallelHandler[T any](action HandlerAction[T], onError HandlerErrorAction[T]) *Handler[T] {
	res := NewHandler(action, onError)
	res.parallel = true
	return res
}

func newFlow[T any](flowOpts *Opts, terminalOnError func(err error), initialData T, handlers ...*Handler[T]) (*flow[T], error) {
	if handlers == nil || len(handlers) < 2 {
		return nil, errors.New("minimum 2 handlers need to be specified")
	}
	if flowOpts == nil {
		// TODO: add random string generation
		flowOpts = &Opts{Name: "name"}
	}

	flowEventLogger.Subject(fmt.Sprintf(subjectFormat, flowOpts.Name))
	logger.Printf("constructing flow with name: %s", flowOpts.Name)

	var firstHandler = handlers[0]
	firstHandler.prev = nil
	firstHandler.next = handlers[1]
	for i, handler := range handlers {
		if handler.action == nil {
			return nil, HandlerMissingActionErr
		}
		handler.id = i
		if i == 0 {
			handler.prev = nil
			handler.next = handlers[1]
		} else {
			handler.prev = handlers[i-1]
			if i != len(handlers)-1 {
				handler.next = handlers[i+1]
			}
		}
		i += 1
	}
	return &flow[T]{
		initialData:     initialData,
		opts:            *flowOpts,
		terminalOnError: terminalOnError,
		firstHandler:    firstHandler,
	}, nil
}

func (r flow[T]) handleError(handler *Handler[T], err error) error {
	if r.terminalOnError != nil && !r.terminalOnErrorExecuted {
		logger.Printf("calling global error fallback due to error %s in handler: %d", err.Error(), handler.id)
		r.terminalOnError(err)
		if !r.opts.ExecuteOnErrorAlways {
			return err
		}
	}
	executeErrorHandler(handler, err)
	return err
}

func executeErrorHandler[T any](handler *Handler[T], err error) {
	handler.onError(handler, err)
	if handler.prev != nil {
		executeErrorHandler(handler.prev, err)
	}
}

func withEventLog[T any](handler *Handler[T]) {
	originalAction := handler.action
	handler.action = func(data T) (T, error) {
		flowEventLogger.LogEvent(actionExecutionStarted, handler.id)
		res, err := originalAction(data)
		if err == nil {
			flowEventLogger.LogEvent(actionExecutionFinished, handler.id)
		}
		return res, err
	}

	originalOnError := handler.onError
	handler.onError = func(handler *Handler[T], err error) {
		flowEventLogger.LogEvent(actionOnError, handler.id)
		originalOnError(handler, err)
	}
}
