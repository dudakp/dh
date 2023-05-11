package mrh

import (
	"dh/internal/logging"
	"dh/pkg/executor"
	"dh/pkg/flow"
	"errors"
)

var (
	logger = logging.GetLoggerFor("mrh")
)

type Mrh struct {
	BranchType  string
	Done        bool
	GitExecutor *executor.GitExecutor
}

func (r *Mrh) Run(issue string) {
	var err error
	err, notDoneFlow := flow.NewEffectFlow(
		&flow.Opts{Name: "not done", ExecuteOnErrorAlways: true},
		r.rollback,
		flow.NewHandler(
			func(data any) (error, any) {
				err := r.GitExecutor.Stash(false)
				return err, nil
			},
			func(handler *flow.Handler[any], handlerErr error) {
				err = errors.Join(handlerErr)
			}),
		flow.NewHandler(
			func(data any) (error, any) {
				err := r.GitExecutor.Checkout(r.BranchType + "/" + issue)
				return err, nil
			},
			func(handler *flow.Handler[any], handlerErr error) {
				err = errors.Join(handlerErr)
			}),
	)
	err, doneFlow := flow.NewEffectFlow(
		&flow.Opts{Name: "done", ExecuteOnErrorAlways: true},
		r.rollback,
		flow.NewHandler(
			func(data any) (error, any) {
				err := r.GitExecutor.Checkout("develop")
				return err, nil
			},
			func(handler *flow.Handler[any], handlerErr error) {
				err = errors.Join(handlerErr)
			}),
		flow.NewHandler(
			func(data any) (error, any) {
				err := r.GitExecutor.Stash(true)
				return err, nil
			},
			func(handler *flow.Handler[any], handlerErr error) {
				err = errors.Join(handlerErr)
			}),
	)
	err = errors.Join(err)
	if err != nil {
		logger.Fatalf("%s", err.Error())
		return
	}
	if !r.Done {
		err = errors.Join(flow.ExecuteEffectFlow(notDoneFlow))
	} else {
		// TODO: add prompt check (checklist like: did test ran succesfully, will sonar check pass? ...)
		// code review done, checkout back to develop and pop stashed changes
		err = errors.Join(flow.ExecuteEffectFlow(doneFlow))
	}

	if err == nil {
		if !r.Done {
			logger.Print("repository is ready for code review")
		} else {
			logger.Print("repository has rolled back to state before code review")
		}
	}
}

func (r *Mrh) rollback(err error) {
	logger.Printf("calling rollback action caused by error: %s", err.Error())
	err = r.GitExecutor.Stash(true)
	if err != nil {
		logger.Fatalf("Error during git stash pop: %s. Please resolver this error manually", err.Error())
	}
}
