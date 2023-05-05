package flow

import "errors"

type Handler[T any] struct {
	Action  func(T) (error, T)
	OnError func(error)
	next    *Handler[T]
	prev    *Handler[T]
}

type flow[T any] struct {
	InitialData     T
	terminalOnError func(err error)
	firstHandler    *Handler[T]
}

func newFlow[T any](terminalOnError func(err error), initialData T, handlers ...*Handler[T]) (error, *flow[T]) {
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
		InitialData:     initialData,
		terminalOnError: terminalOnError,
		firstHandler:    firstHandler,
	}
}

func (r flow[T]) handleError(handler *Handler[T], err error) error {
	if r.terminalOnError != nil {
		r.terminalOnError(err)
		return err
	}
	r.executeErrorHandler(handler, err)
	return err
}

func (r flow[T]) executeErrorHandler(handler *Handler[T], err error) {
	handler.OnError(err)
	if handler.prev != nil {
		r.executeErrorHandler(handler.prev, err)
	}
}
