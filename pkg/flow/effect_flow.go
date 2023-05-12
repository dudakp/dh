package flow

type EffectFlow[T any] struct {
	*flow[T]
}

// NewEffectFlow creates effect flow.
// If terminalOnError is specified this function will be used as global error callback and no handler onError functions will be called.
// to override this behavior look at Opts.ExecuteOnErrorAlways
func NewEffectFlow(flowOpts *Opts, terminalOnError func(err error), handlers ...*Handler[any]) (*EffectFlow[any], error) {
	var baseHandlers []*Handler[any]
	for _, handler := range handlers {
		baseHandlers = append(baseHandlers, handler)
	}
	var empty any
	f, err := newFlow(flowOpts, terminalOnError, empty, baseHandlers...)
	return &EffectFlow[any]{
		flow: f,
	}, err
}

func ExecuteEffectFlow(f *EffectFlow[any]) error {
	logger.Printf("executing flow: %s", f.opts.Name)
	err := executeEffect(f.firstHandler, f)
	logger.Printf("flow: %s executed successful", f.opts.Name)
	return err
}

func executeEffect(handler *Handler[any], f *EffectFlow[any]) error {
	var empty any
	return execute(handler, empty, f.flow)
}
