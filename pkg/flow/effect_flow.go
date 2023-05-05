package flow

type EffectFlow[T any] struct {
	baseFlow *flow[T]
}

type EffectHandler[T any] struct {
	EffectAction func() error
	baseHandler  *Handler[T]
}

func NewEffectFlow[T any](terminalOnError func(err error), handlers ...*EffectHandler[T]) (error, *EffectFlow[T]) {
	var empty T
	var baseHandlers []*Handler[T]
	for _, handler := range handlers {
		baseHandlers = append(baseHandlers, handler.baseHandler)
	}
	err, f := newFlow(terminalOnError, empty, baseHandlers...)
	return err, &EffectFlow[T]{
		baseFlow: f,
	}
}

func ExecuteEffectFlow[T any](f *EffectFlow[T]) error {
	var empty T
	return executeEffect(f.baseFlow.firstHandler, empty, f)
}

func executeEffect[T any](handler *Handler[T], data T, f *EffectFlow[T]) error {
	if handler.next == nil {
		return nil
	}
	var empty T
	err, _ := handler.Action(data)
	if err != nil {
		return f.baseFlow.handleError(handler, err)
	} else {
		return executeEffect(handler.next, empty, f)
	}
}
