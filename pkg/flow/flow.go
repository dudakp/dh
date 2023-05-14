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
		identify groups of parallel adjacent parallel handlers in flow and create fork-join
			all handlers in group need to communicate via channel, if execution of one starts this handler needs
			to notify all in group to start, also sync handler that is after async group needs to receive message that
			all async handlers have finished and can start its execution
		parallel handler will be supported only in EffectFlow for now (reduce operation would be needed on join) | DONE
			EffectFlow does not have any data so nothing to worry about
	* conditional
		based on predicate in handler, the execution will continue to specified handler


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
	"context"
	"dh/pkg/logging"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync"
)

const (
	actionOnError           = "onError(handlerId:%d)"
	actionExecutionStarted  = "action_started(handlerId:%d)"
	actionExecutionFinished = "action_finished(handlerId:%d)"

	subjectFormat = "flow(%s)"
)

var (
	HandlerMissingActionErr     = errors.New("no action function specified for handler")
	ParallelHandlerNotSupported = errors.New("parallel handler is supported only in EffectFlow")

	logger          = logging.GetLoggerFor("flow")
	flowEventLogger = logging.GetEventLoggerFor("flow").
			Format(logging.EventLoggerSubjectFormat, logging.EventLoggerEventFormat)
)

type HandlerAction[T any] func(T) (T, error)
type HandlerErrorAction[T any] func(*Handler[T], error)

//type ParallelHandler[T any] Handler[T]

type Handler[T any] struct {
	action  HandlerAction[T]
	onError HandlerErrorAction[T]

	id int

	next *Handler[T]
	prev *Handler[T]
	// pGroup contains parallel handlers that can be executed in parallel
	pGroup []*Handler[T]

	// sig is used only with Handler with pGroup handlers for signaling the next handler that parallel execution
	//has finished and next handler can start its execution
	sig *chan bool
	// wg is used for synchronizing execution of every ParallelHandler in pGroup
	wg *sync.WaitGroup
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

func NewParallelHandlerGroup[T any](handlers ...*Handler[T]) *Handler[T] {
	var wg *sync.WaitGroup
	var sig = make(chan bool)
	for i, handler := range handlers {
		handler.id = -(1 + i)
	}
	return &Handler[T]{
		pGroup: handlers,
		wg:     wg,
		sig:    &sig,
	}
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

	err := chainHandlers(handlers)
	if err != nil {
		return nil, err
	}

	return &flow[T]{
		initialData:     initialData,
		opts:            *flowOpts,
		terminalOnError: terminalOnError,
		firstHandler:    handlers[0],
	}, nil
}

func (r *flow[T]) Start() error {
	return execute(r.firstHandler, r.initialData, r)
}

func (r *flow[T]) handleError(handler *Handler[T], err error) error {
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

func chainHandlers[T any](handlers []*Handler[T]) error {
	var firstHandler = handlers[0]
	firstHandler.prev = nil
	firstHandler.next = handlers[1]
	for i, handler := range handlers {
		// handlers containing pGroup does not need to have an action
		if handler.action == nil && !handler.isPgroup() {
			return HandlerMissingActionErr
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
	return nil
}

func execute[T any](handler *Handler[T], handlerOutput T, f *flow[T]) error {
	// for parallel handel group execute group
	if handler.isPgroup() {
		return handler.executePgroup(f)
	} else {
		// for simple handler execute
		out, err := handler.action(handlerOutput)
		if err != nil {
			return f.handleError(handler, err)
		} else {
			if handler.next == nil {
				return nil
			} else {
				return execute(handler.next, out, f)
			}
		}
	}
}

func (r *Handler[T]) executePgroup(f *flow[T]) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, handler := range r.pGroup {
		var empty T
		eg.Go(func() error {
			return execute(handler, empty, f)
		})
	}
	return eg.Wait()
}

func (r *Handler[T]) isPgroup() bool {
	return r.sig != nil
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
