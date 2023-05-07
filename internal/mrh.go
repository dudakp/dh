package internal

import (
	"dh/internal/config"
	"dh/pkg/executor"
	"dh/pkg/flow"
	"errors"
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
		config.ErrLog.Fatalf("%s")
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
			config.InfoLog.Print("repository is ready for code review")
		} else {
			config.InfoLog.Print("repository has rolled back to state before code review")
		}
	}
}

func (r *Mrh) rollback(err error) {
	config.ErrLog.Printf("calling rollback action caused by error: %s", err.Error())
	err = r.GitExecutor.Stash(true)
	if err != nil {
		config.ErrLog.Fatalf("Error during git stash pop: %s. Please resolver this error manually", err.Error())
	}
}
