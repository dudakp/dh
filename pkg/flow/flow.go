package flow

/**
TODO: create structured event logger - outputted to console && file
	data recorded in eventLog: flow id, flow name, handler id, timestamp, handler data, action successful, onError called, onTerminalError called

TODO: create operators
	* pipe
	* map

TODO: create execution strategies
	* sequential - DONE
	* parallel
	* synchronized

TODO: create flow types:
	* ApiFlow
		flow calling REST API
	* TemplatedFlow
		flow definition in file, calling scripts (bash scripts)
*/

import (
	"dh/internal/logging"
	"errors"
)

var (
	logger                  = logging.GetLoggerFor("flow")
	HandlerMissingActionErr = errors.New("no action function specified for handler")
)

type HandlerAction[T any] func(T) (error, T)
type HandlerErrorAction[T any] func(*Handler[T], error)

type Handler[T any] struct {
	action        func(T) (error, T)
	onError       func(*Handler[T], error)
	id            int
	next          *Handler[T]
	prev          *Handler[T]
	wrappedAction func(T) (error, T)
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
			logger.Printf("calling onError from handler with id: %d", handler.id)
		}
	}
	withEventLog(res)
	return res
}

func newFlow[T any](flowOpts *Opts, terminalOnError func(err error), initialData T, handlers ...*Handler[T]) (error, *flow[T]) {
	if handlers == nil || len(handlers) < 2 {
		return errors.New("minimum 2 handlers need to be specified"), nil
	}
	if flowOpts == nil {
		// TODO: add random string generation
		flowOpts = &Opts{Name: "name"}
	}

	logger.Printf("constructing flow with name: %s", flowOpts.Name)

	var firstHandler = handlers[0]
	firstHandler.prev = nil
	firstHandler.next = handlers[1]
	for i, handler := range handlers {
		if handler.action == nil {
			return HandlerMissingActionErr, nil
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
	return nil, &flow[T]{
		initialData:     initialData,
		opts:            *flowOpts,
		terminalOnError: terminalOnError,
		firstHandler:    firstHandler,
	}
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
	handler.action = func(data T) (error, T) {
		logger.Printf("executing handler: %d", handler.id)
		err, res := originalAction(data)
		if err == nil {
			logger.Printf("handler: %d executed successfully", handler.id)
		}
		return err, res
	}

	originalOnError := handler.onError
	handler.onError = func(handler *Handler[T], err error) {
		logger.Printf("executing onError for handler: %d", handler.id)
		originalOnError(handler, err)
		logger.Printf("onError for handler: %d executed successfully", handler.id)
	}
}
