package flow

type PureFlow[T any] struct {
	*flow[T]
}

// NewPureFlow creates effect flow.
// If terminalOnError is specified this function will be used as global error callback and no handler onError functions will be called.
// to override this behavior look at Opts.ExecuteOnErrorAlways
func NewPureFlow[T any](flowOpts *Opts, terminalOnError func(err error), initialData T, handlers ...*Handler[T]) (*PureFlow[T], error) {
	for _, handler := range handlers {
		if handler.isPgroup() {
			return nil, ErrParallelHandlerNotSupported
		}
	}
	var baseHandlers []*Handler[T]
	for _, handler := range handlers {
		baseHandlers = append(baseHandlers, handler)
	}
	f, err := newFlow(flowOpts, terminalOnError, initialData, baseHandlers...)
	return &PureFlow[T]{
		flow: f,
	}, err
}

func ExecutePureFlow(f *PureFlow[any]) error {
	logger.Printf("executing flow: %s", f.opts.Name)
	err := f.Start()
	logger.Printf("flow: %s executed successful", f.opts.Name)
	return err
}
