package flow

import (
	"errors"
	"testing"
)

func Test_EffectHandler_ExecuteEffectFlow(t *testing.T) {
	err, f := NewEffectFlow(&Opts{Name: "TExecuteEffectFlow"}, nil,
		NewHandler(func(data any) (error, any) { return nil, "" }, nil),
		NewHandler(func(data any) (error, any) { return nil, "" }, nil),
	)
	if err != nil {
		t.Fatalf("error during creation of flow: %s", err.Error())
	}
	var tests = []struct {
		name  string
		input *EffectFlow[any]
		want  string
	}{
		{
			name:  "Handler - ExecuteFlow no errors",
			input: f,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ExecuteEffectFlow(f)
			if err != nil {
				t.Fatalf("flow execution error: %s", err.Error())
			}
		})
	}
}

func Test_NewEffectFlow(t *testing.T) {
	handler0 := NewHandler(func(data any) (error, any) { return nil, "" }, nil)
	handler1 := NewHandler(func(data any) (error, any) { return nil, "" }, nil)
	handler2 := NewHandler(func(data any) (error, any) { return nil, "" }, nil)

	err, flow := NewEffectFlow(&Opts{Name: "TNewEffectFlow"}, nil, handler0, handler1, handler2)
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
	err, _ := newFlow(&Opts{Name: "TMinHandlers"}, nil, NewHandler(func(data any) (error, any) { return nil, "" }, nil))
	if err == nil {
		t.FailNow()
	}
}

// TODO: add more
func Test_ExecuteEffectFlow_errorPropagation(t *testing.T) {
	in := "hello"
	expectedError := errors.New("shit happens")

	handler0 := NewHandler(func(data any) (error, any) { return nil, "" }, func(handler *Handler[any], err error) {
		in = "error"
	})
	handler1 := NewHandler(func(data any) (error, any) { return nil, "" }, nil)
	handler2 := NewHandler(func(data any) (error, any) { return expectedError, nil }, nil)

	err, f := NewEffectFlow(&Opts{Name: "TErrorPropagation"}, nil, handler0, handler1, handler2)
	if err != nil {
		t.Fatalf("error during creation of flow: %s", err.Error())
	}
	tests := []struct {
		name  string
		input *EffectFlow[any]
		want  string
	}{
		{
			name:  "last in flow has error",
			want:  "error",
			input: f,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ExecuteEffectFlow(test.input)
			if err != nil && !errors.Is(err, expectedError) {
				t.Errorf("flow ended with error: %s", err.Error())
				t.FailNow()
			}
			if in != test.want {
				t.Errorf("want: %s, got: %s", test.want, in)
				t.FailNow()
			}
			in = "hello"
		})
	}
}
