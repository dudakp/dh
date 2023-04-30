package flow

import "errors"

type HandlerData interface{}

type Handler struct {
	Action  func(HandlerData) (error, HandlerData)
	OnError func(error)
	next    *Handler
	prev    *Handler
}

type Flow struct {
	InitialData     HandlerData
	terminalOnError func(err error)
	firstHandler    *Handler
}

func NewFlow(terminalOnError func(err error), handlers ...*Handler) (error, *Flow) {
	if len(handlers) < 2 {
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
	return nil, &Flow{
		terminalOnError: terminalOnError,
		firstHandler:    firstHandler,
	}
}

func handleError(handler *Handler, err error, r *Flow) error {
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
