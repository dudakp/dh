package flow

import (
	"errors"
	"testing"
)

func Test_EffectHandler_ExecuteEffectFlow(t *testing.T) {
	err, flow := NewEffectFlow(nil, &Handler{}, &Handler{})
	if err != nil {
		t.Fatalf("error during creation of flow: %s", err.Error())
	}
	var tests = []struct {
		name  string
		input *EffectFlow
		want  string
	}{
		{
			name:  "EffectHandler - ExecuteFlow no errors",
			input: flow,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ExecuteEffectFlow(flow)
			if err != nil {
				t.Fatalf("flow execution error: %s", err.Error())
			}
		})
	}
}

func Test_NewFlow(t *testing.T) {
	handler0 := &Handler{}
	handler1 := &Handler{}
	handler2 := &Handler{}

	err, flow := newFlow(nil, handler0, handler1, handler2)
	if err != nil {
		t.Errorf("error during creation of flow: %s", err.Error())
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
	err, _ := newFlow(nil, &Handler{})
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
				Action: func(data Data) (error, Data) {
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
				Action: func(data Data) (error, Data) {
					return nil, nil
				},
				OnError: func(err error) {
					in = "world"
				},
				prev: &Handler{
					Action: func(data Data) (error, Data) {
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
				t.Errorf("want: %s, got: %s", test.want, in)
				t.FailNow()
			}
			in = "hello"
		})
	}
}
