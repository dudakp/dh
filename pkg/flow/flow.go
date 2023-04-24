package flow

import "errors"

type EffectHandler struct {
	Action  func() error
	OnError func(error)
	next    *EffectHandler
	prev    *EffectHandler
}

type Controller interface {
	CreateFlow(handlers ...*EffectHandler) error
	ExecuteFlow() error
}

type Flow struct {
	TerminalOnError func(err error)
	firstHandler    *EffectHandler
}

func (r *Flow) CreateFlow(handlers ...*EffectHandler) error {
	if len(handlers) < 2 {
		return errors.New("minimum 2 handlers need to be specified")
	}
	var res = handlers[0]
	res.prev = nil
	res.next = handlers[1]
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
	r.firstHandler = res
	return nil
}

func (r *Flow) ExecuteFlow() error {
	return r.execute(r.firstHandler)
}

func (r *Flow) execute(handler *EffectHandler) error {
	err := handler.Action()
	if err != nil {
		executeErrorHandler(handler, err)
		return err
	} else {
		if r.TerminalOnError != nil {
			r.TerminalOnError(err)
			return err
		}
		if handler.next != nil {
			return r.execute(handler.next)
		}
		return err
	}
}

func executeErrorHandler(handlerWithError *EffectHandler, err error) {
	handlerWithError.OnError(err)
	if handlerWithError.prev != nil {
		executeErrorHandler(handlerWithError.prev, err)
	}
}
