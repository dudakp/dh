package flow

/**

TODO: create operators
	* conditional
		based on predicate in handler, the execution will continue to specified handler - maybe flow like this cant be constructed manually (think about this)


TODO: create flow types:
	* TemplatedFlow
		flow definition in file, calling scripts (bash scripts)
		flow loaded from event log
	* ApiFlow
		flow calling REST API
		not neccesary restFlow but create restHandler that can be incorporated info any flow


TODO: create event log parser
	nice to have - not priority
*/

import (
	"context"
	"errors"
	"fmt"
	"github.com/dudakp/dh/pkg/logging"
	"golang.org/x/sync/errgroup"
)

const (
	actionOnError           = "onError(handlerId:%d)"
	actionOnErrorPgroup     = "onError(pGroupHandlerId:%d)"
	actionExecutionStarted  = "action_started(handlerId:%d)"
	actionExecutionFinished = "action_finished(handlerId:%d)"

	subjectFormat = "flow(%s)"
)

var (
	ErrHandlerMissingAction        = errors.New("no action function specified for handler")
	ErrParallelHandlerNotSupported = errors.New("parallel handler is supported only in EffectFlow")

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
	// isParallel is marking handler that is part of pGroup
	isParallel bool
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

func (r *flow[T]) Start() error {
	return execute(r.firstHandler, r.initialData, r)
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

func NewHandler[T any](action HandlerAction[T], onError HandlerErrorAction[T]) *Handler[T] {
	var res = &Handler[T]{
		action:  action,
		onError: onError,
	}
	if res.onError == nil {
		logger.Printf("missing onError function for handler. creating default onError")
		res.onError = func(handler *Handler[T], err error) {
			logger.Printf("calling default onError for handler: %d", handler.id)
		}
	}
	return withEventLog(res)
}

// NewParallelHandlerGroup creates handler wrapping multiple handlers marked for parallel execution.
// if any error occurs in parallel handler, onError of all handlers in ParallelGroup are called,
// if error occurs in handler ordinary handler that was called after ParallelGroup,
// only ParallelGroup onError will be executed
func NewParallelHandlerGroup[T any](onError HandlerErrorAction[T], handlers ...*Handler[T]) *Handler[T] {
	for i, handler := range handlers {
		handler.id = -(1 + i)
		handler.isParallel = true
	}
	e := onError
	if onError == nil {
		e = func(h *Handler[T], err error) {
			logger.Printf("calling default onError for handler: %d", h.id)
		}
	}
	return withEventLog(
		// TODO: try to use NewHandler
		&Handler[T]{
			pGroup:  handlers,
			onError: e,
		})
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
			return ErrHandlerMissingAction
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
	var out T
	var err error

	// for execution of pGroup call special method for execution of pGroups
	if handler.isPgroup() {
		err = handler.executePgroup(f)
	} else {
		out, err = handler.action(handlerOutput)
	}
	if err == nil {
		if handler.next == nil {
			return nil
		} else {
			return execute(handler.next, out, f)
		}
	} else if handler.isPgroup() {
		// parallel handlers in pGroup are not linked so theirs onError needed to be called like this
		// TODO: maybe run onError in parallel?
		for _, h := range handler.pGroup {
			h.onError(h, err)
		}
		return err
	} else if handler.isParallel {
		// parallel onError handler is called in else-if above so here it needs to be skipped
		return err
	} else {
		// regular error handling for sequential handlers
		return f.handleError(handler, err)
	}
}

func (r *Handler[T]) executePgroup(f *flow[T]) error {
	eg, _ := errgroup.WithContext(context.Background())
	for _, handler := range r.pGroup {
		var empty T
		// go-specific hack for parallel function execution in for loop
		handler := handler
		eg.Go(func() error {
			res := execute(handler, empty, f)
			return res
		})
	}
	// TODO: think about error propagation strategies for pGroup
	return eg.Wait()
}

func (r *Handler[T]) isPgroup() bool {
	return len(r.pGroup) > 0
}

func executeErrorHandler[T any](handler *Handler[T], err error) {
	handler.onError(handler, err)
	if handler.prev != nil {
		executeErrorHandler(handler.prev, err)
	}
}

func withEventLog[T any](handler *Handler[T]) *Handler[T] {
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
		if handler.isPgroup() {
			flowEventLogger.LogEvent(actionOnErrorPgroup, handler.id)
		} else {
			flowEventLogger.LogEvent(actionOnError, handler.id)
		}
		originalOnError(handler, err)
	}
	return handler
}
