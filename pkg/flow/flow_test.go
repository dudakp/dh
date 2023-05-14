package flow

import (
	"errors"
	"testing"
)

func Test_EffectHandler_ExecuteEffectFlow(t *testing.T) {
	f, err := NewEffectFlow(&Opts{Name: "TExecuteEffectFlow"}, nil,
		NewHandler(func(data any) (any, error) { return "", nil }, nil),
		NewHandler(func(data any) (any, error) { return "", nil }, nil),
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
	handler0 := NewHandler(func(data any) (any, error) { return "", nil }, nil)
	handler1 := NewHandler(func(data any) (any, error) { return "", nil }, nil)
	handler2 := NewHandler(func(data any) (any, error) { return "", nil }, nil)

	f, err := NewEffectFlow(&Opts{Name: "TNewEffectFlow"}, nil, handler0, handler1, handler2)
	if err != nil {
		t.Errorf("error during creation of f: %s", err.Error())
	}

	// check first handler
	if f.firstHandler != handler0 {
		t.FailNow()
	}
	if f.firstHandler.next != handler1 {
		t.FailNow()
	}
	if f.firstHandler.prev != nil {
		t.FailNow()
	}

	// check second handler
	if f.firstHandler.next.prev != handler0 {
		t.FailNow()
	}
	if f.firstHandler.next.next != handler2 {
		t.FailNow()
	}

	//	check third handler
	if f.firstHandler.next.next != handler2 {
		t.FailNow()
	}
	if f.firstHandler.next.next.prev != handler1 {
		t.FailNow()
	}
	if f.firstHandler.next.next.next != nil {
		t.FailNow()
	}
}

func Test_NewFlow_minHandlers(t *testing.T) {
	_, err := newFlow(&Opts{Name: "TMinHandlers"}, nil, NewHandler(func(data any) (any, error) { return "", nil }, nil))
	if err == nil {
		t.FailNow()
	}
}

func Test_ExecuteEffectFlow_errorPropagation(t *testing.T) {
	in := "hello"
	expectedError := errors.New("shit happens")

	//handler0 := NewHandler(func(data any) (any, error) {
	//	return "", nil
	//}, func(handler *Handler[any], err error) {
	//	in = "error"
	//})
	//handler1 := NewHandler(func(data any) (any, error) {
	//	return "", nil
	//}, nil)
	//handler2 := NewHandler(func(data any) (any, error) {
	//	return nil, expectedError
	//}, nil)

	tests := []struct {
		name  string
		input *EffectFlow[any]
		want  string
	}{
		{
			name: "last in flow has error",
			want: "error",
			input: func() *EffectFlow[any] {
				f, err := NewEffectFlow(&Opts{Name: "lastInFlowWithErr"}, nil,
					NewHandler(func(data any) (any, error) {
						return "", nil
					}, func(handler *Handler[any], err error) {
						in = "error"
					}), NewHandler(func(data any) (any, error) {
						return "", nil
					}, nil), NewHandler(func(data any) (any, error) {
						return nil, expectedError
					}, nil))
				if err != nil {
					t.Fatalf("error during creation of flow: %s", err.Error())
				}
				return f
			}(),
		},
		{
			name: "second-to-last in flow has error",
			want: "error",
			input: func() *EffectFlow[any] {
				f, err := NewEffectFlow(&Opts{Name: "stlInFlowWithError"}, nil,
					NewHandler(func(data any) (any, error) {
						return "", nil
					}, func(handler *Handler[any], err error) {
						in = "error"
					}), NewHandler(func(data any) (any, error) {
						return nil, expectedError
					}, nil), NewHandler(func(data any) (any, error) {
						return "", nil
					}, nil))
				if err != nil {
					t.Fatalf("error during creation of flow: %s", err.Error())
				}
				return f
			}(),
		},
		{
			name: "last in flow has error and terminalOnError is defined",
			want: "global error",
			input: func() *EffectFlow[any] {
				f, err := NewEffectFlow(&Opts{Name: "terminalOnError"},
					func(err error) {
						in = "global error"
					}, NewHandler(func(data any) (any, error) {
						return "", nil
					}, func(handler *Handler[any], err error) {
						in = "error"
					}), NewHandler(func(data any) (any, error) {
						return "", nil
					}, nil), NewHandler(func(data any) (any, error) {
						return "", expectedError
					}, nil))
				if err != nil {
					t.Fatalf("error during creation of flow: %s", err.Error())
				}
				return f
			}(),
		},
		{
			name: "last in flow has error and terminalOnError is defined and option for executing all handers is true",
			want: "helloglobal errorerror",
			input: func() *EffectFlow[any] {
				f, err := NewEffectFlow(&Opts{Name: "terminalOnErrorOverride", ExecuteOnErrorAlways: true},
					func(err error) {
						in += "global error"
					}, NewHandler(func(data any) (any, error) {
						return "", nil
					}, func(handler *Handler[any], err error) {
						in += "error"
					}), NewHandler(func(data any) (any, error) {
						return "", nil
					}, nil), NewHandler(func(data any) (any, error) {
						return nil, expectedError
					}, nil))
				if err != nil {
					t.Fatalf("error during creation of flow: %s", err.Error())
				}
				return f
			}(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			in = "hello"
			err := ExecuteEffectFlow(test.input)
			if err != nil && !errors.Is(err, expectedError) {
				t.Errorf("flow ended with error: %s", err.Error())
				t.FailNow()
			}
			if in != test.want {
				t.Errorf("want: %s, got: %s", test.want, in)
				t.FailNow()
			}

		})
	}
}
