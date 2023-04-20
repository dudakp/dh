package git

import (
	"bytes"
	"dh/pkg/config"
	"fmt"
	"os/exec"
)

type Executor struct {
	path string
}

func NewGitExecutor(gitPath string) *Executor {
	return &Executor{path: gitPath}
}

func (r *Executor) Checkout(issueBranch string) error {
	return r.execute("checkout", issueBranch)
}

func (r *Executor) Stash(pop bool) error {
	arg := "stash"
	if pop {
		if !r.stashHasEntries() {
			config.WarnLog.Print("no stash entries in repo")
			return nil
		} else {
			return r.execute(arg, "pop")
		}
	} else {
		return r.execute(arg)
	}
}

func (r *Executor) stashHasEntries() bool {
	err, stdout, _ := r.executeWithResult("stash", "list")
	if err != nil {
		config.ErrLog.Fatal(err.Error())
		return false
	}
	if stdout.Len() < 1 {
		return false
	}
	return true
}

func (r *Executor) executeWithResult(command string, flags ...string) (error, *bytes.Buffer, *bytes.Buffer) {
	argWithFlags := []string{command}
	argWithFlags = append(argWithFlags, flags...)
	cmd := exec.Command(r.path, argWithFlags...)
	config.InfoLog.Printf("executing command: %s", cmd.String())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		config.ErrLog.Printf(fmt.Errorf("%w: %s", err, stderr.String()).Error())
		return fmt.Errorf("\n\t%w: unable to execute git command: %s\n\tcall ended with error: %s",
			err, argWithFlags, stderr.String(),
		), nil, &stderr
	}
	config.InfoLog.Print(stdout.String())
	return nil, &stdout, nil
}

func (r *Executor) execute(command string, flags ...string) error {
	err, _, _ := r.executeWithResult(command, flags...)
	return err
}
