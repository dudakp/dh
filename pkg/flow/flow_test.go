package flow

import (
	"errors"
	"fmt"
	"testing"
)

func Test_EffectHandler_ExecuteEffectFlow(t *testing.T) {
	err, flow := NewFlow(nil, &Handler{
		Action:  nil,
		OnError: nil,
		next:    nil,
		prev:    nil,
	})
	if err != nil {
		_ = fmt.Errorf("%w", err)
		t.FailNow()
	}
	var tests = []struct {
		name  string
		input *Flow
		want  string
	}{
		{
			name:  "EffectHandler - ExecuteFlow no errors",
			input: flow,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
		})
	}
}

func Test_NewFlow(t *testing.T) {
	handler0 := &Handler{}
	handler1 := &Handler{}
	handler2 := &Handler{}

	err, flow := NewFlow(nil, handler0, handler1, handler2)
	if err != nil {
		_ = fmt.Errorf("%s", err)
		t.FailNow()
	}

	// check first handler
	if flow.firstHandler != handler0 {
		t.FailNow()
	}
	if flow.firstHandler.next != handler1 {
		t.FailNow()
	}
	if flow.firstHandler.prev != nil {
		t.FailNow()
	}

	// check second handler
	if flow.firstHandler.next.prev != handler0 {
		t.FailNow()
	}
	if flow.firstHandler.next.next != handler2 {
		t.FailNow()
	}

	//	check third handler
	if flow.firstHandler.next.next != handler2 {
		t.FailNow()
	}
	if flow.firstHandler.next.next.prev != handler1 {
		t.FailNow()
	}
	if flow.firstHandler.next.next.next != nil {
		t.FailNow()
	}
}

func Test_NewFlow_minHandlers(t *testing.T) {
	err, _ := NewFlow(nil, &Handler{})
	if err == nil {
		t.FailNow()
	}
}

func Test_executeErrorHandler(t *testing.T) {
	in := "hello"
	tests := []struct {
		name  string
		input *Handler
		want  string
	}{
		{
			name: "last in flow",
			want: "error",
			input: &Handler{
				Action: func(data HandlerData) (error, HandlerData) {
					return nil, nil
				},
				OnError: func(err error) {
					in = "error"
				},
			},
		},
		{
			name: "not last in flow",
			want: "error",
			input: &Handler{
				Action: func(data HandlerData) (error, HandlerData) {
					return nil, nil
				},
				OnError: func(err error) {
					in = "world"
				},
				prev: &Handler{
					Action: func(data HandlerData) (error, HandlerData) {
						return nil, nil
					},
					OnError: func(err error) {
						in = "error"
					},
				},
			}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			executeErrorHandler(test.input, errors.New("err"))
			if in != test.want {
				fmt.Printf("want: %s, got: %s", test.want, in)
				t.FailNow()
			}
			in = "hello"
		})
	}
}
