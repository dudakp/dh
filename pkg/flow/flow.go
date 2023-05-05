package flow

import "errors"

// TODO: add generic parameters so that HanderData can be used
type Data[T any] struct {
	data T
}

type Handler struct {
	Action  func(Data[any]) (error, Data[any])
	OnError func(error)
	next    *Handler
	prev    *Handler
}

type flow[T any] struct {
	InitialData     Data[T]
	terminalOnError func(err error)
	firstHandler    *Handler
}

func newFlow[T any](terminalOnError func(err error), initialData T, handlers ...*Handler) (error, *flow[T]) {
	if handlers == nil || len(handlers) < 2 {
		return errors.New("minimum 2 handlers need to be specified"), nil
	}

	var firstHandler = handlers[0]
	firstHandler.prev = nil
	firstHandler.next = handlers[1]
	for i, handler := range handlers {
		if i == 0 {
			handler.prev = nil
			handler.next = handlers[1]
		} else {
			handler.prev = handlers[i-1]
			if i != len(handlers)-1 {
				handler.next = handlers[i+1]
			}
		}
	}
	return nil, &flow[T]{
		InitialData:     Data[T]{data: initialData},
		terminalOnError: terminalOnError,
		firstHandler:    firstHandler,
	}
}

func handleError[T any](handler *Handler, err error, r *flow[T]) error {
	if r.terminalOnError != nil {
		r.terminalOnError(err)
		return err
	}
	executeErrorHandler(handler, err)
	return err
}

func executeErrorHandler(handler *Handler, err error) {
	handler.OnError(err)
	if handler.prev != nil {
		executeErrorHandler(handler.prev, err)
	}
}
