package flow

import (
	"io"
	"net/http"
	"testing"
)

func TestExecuteRestFlow(t *testing.T) {
	handler := NewRestHandler("https://numbersapi.p.rapidapi.com/random/trivia", http.MethodGet, nil)
	handler1 := NewRestHandler("https://api.api-ninjas.com/v1/facts?limit=1", http.MethodGet, nil)
	handler2 := NewHandler(func(t io.ReadCloser) (io.ReadCloser, error) {
		return t, nil
	}, nil)
	restFlow, err := NewPureFlow(nil, nil, nil, handler, handler1, handler2)
	if err != nil {
		t.FailNow()
	}
	err = ExecuteRestFlow(restFlow)
}
