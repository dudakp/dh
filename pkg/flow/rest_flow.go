package flow

import (
	"errors"
	"io"
	"net/http"
)

type RestFlow[T any] struct {
	*flow[T]
}

type RestHandler[T any] struct {
	*Handler[T]
	apiPath string
	verb    int
}

func NewRestFlow[T any](flowOpts *Opts, terminalOnError func(err error), initialData T, handlers ...*Handler[T]) (*RestFlow[T], error) {
	f, err := newFlow(flowOpts, terminalOnError, initialData, handlers...)
	return &RestFlow[T]{
		flow: f,
	}, err
}

// NewRestHandler for method use constants from http.method file
func NewRestHandler[T io.ReadCloser](url string, method string, body io.ReadCloser) *Handler[T] {
	return NewHandler(func(t T) (T, error) {
		request, _ := http.NewRequest(method, url, body)
		client := http.Client{} // TODO: pull this out of this method, no need to create http client for each request
		resp, err := client.Do(request)
		var empty T
		if err != nil {
			return empty, err
		}
		if resp.StatusCode != 200 {
			return empty, errors.New(resp.Status)
		}
		return resp.Body, err
	}, nil)
}

func ExecuteRestFlow(f *RestFlow[io.ReadCloser]) error {
	logger.Printf("executing REST flow: %s", f.opts.Name)
	err := f.Start()
	logger.Printf("REST flow: %s executed successful", f.opts.Name)
	return err
}
