package flow

import "dh/internal/logging"

type EffectFlow[T any] struct {
	*flow[T]
}

// NewEffectFlow creates effect flow.
// If terminalOnError is specified this function will be used as global error callback and no handler onError functions will be called.
// to override this behavior look at Opts.ExecuteOnErrorAlways
func NewEffectFlow(flowOpts *Opts, terminalOnError func(err error), handlers ...*Handler[any]) (error, *EffectFlow[any]) {
	var baseHandlers []*Handler[any]
	for _, handler := range handlers {
		baseHandlers = append(baseHandlers, handler)
	}
	var empty any
	err, f := newFlow(flowOpts, terminalOnError, empty, baseHandlers...)
	return err, &EffectFlow[any]{
		flow: f,
	}
}

func ExecuteEffectFlow(f *EffectFlow[any]) error {
	logging.InfoLog.Printf("executing flow: %s", f.opts.Name)
	err := executeEffect(f.firstHandler, f)
	logging.InfoLog.Printf("flow: %s executed successful", f.opts.Name)
	return err
}

func executeEffect(handler *Handler[any], f *EffectFlow[any]) error {
	var empty any
	err, _ := handler.action(empty)
	if err != nil {
		return f.flow.handleError(handler, err)
	} else {
		if handler.next == nil {
			return nil
		} else {
			return executeEffect(handler.next, f)
		}
	}
}
