package flow

import "dh/internal/logging"

type PureFlow[T any] struct {
	*flow[T]
}

// NewPureFlow creates effect flow.
// If terminalOnError is specified this function will be used as global error callback and no handler onError functions will be called.
// to override this behavior look at Opts.ExecuteOnErrorAlways
func NewPureFlow[T any](flowOpts *Opts, terminalOnError func(err error), initialData T, handlers ...*Handler[T]) (error, *PureFlow[T]) {
	var baseHandlers []*Handler[T]
	for _, handler := range handlers {
		baseHandlers = append(baseHandlers, handler)
	}
	err, f := newFlow(flowOpts, terminalOnError, initialData, baseHandlers...)
	return err, &PureFlow[T]{
		flow: f,
	}
}

func ExecutePureFlow(f *PureFlow[any]) error {
	logging.InfoLog.Printf("executing flow: %s", f.opts.Name)
	err := executePure(f.firstHandler, f.initialData, f)
	logging.InfoLog.Printf("flow: %s executed successful", f.opts.Name)
	return err
}

func executePure[T any](handler *Handler[any], handlerOutput T, f *PureFlow[any]) error {
	err, out := handler.action(handlerOutput)
	if err != nil {
		return f.flow.handleError(handler, err)
	} else {
		if handler.next == nil {
			return nil
		} else {
			return executePure(handler.next, out, f)
		}
	}
}
