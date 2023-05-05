package flow

type EffectFlow[T any] flow[T]

func NewEffectFlow[T any](terminalOnError func(err error), handlers ...*Handler) (error, *EffectFlow[T]) {
	err, f := newFlow(terminalOnError, nil, handlers...)
	return err, (*EffectFlow[T])(f)
}

func ExecuteEffectFlow(f *EffectFlow[any]) error {
	return executeEffect(f.firstHandler, (*flow[any])(f))
}

func executeEffect(handler *Handler, flow *flow[any]) error {
	err, _ := handler.Action(flow.InitialData)
	if err != nil {
		return handleError(handler, err, flow)
	} else {
		return executeEffect(handler.next, flow)
	}
}
