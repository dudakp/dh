package executor

import (
	"bytes"
	"dh/internal/config"
	"fmt"
	"os/exec"
)

type IExecutor interface {
	executeWithResult(command string, flags ...string) (error, *bytes.Buffer, *bytes.Buffer)
	execute(command string, flags ...string) error
}
type BaseExecutor struct {
	executableName string
	path           string
}

func newBaseExecutor(executableName string) BaseExecutor {
	return BaseExecutor{
		executableName: executableName,
		path:           findExecutable(executableName),
	}
}

func findExecutable(executableName string) string {
	path, err := exec.LookPath(executableName)
	if err != nil {
		config.ErrLog.Fatalf("Unable to find executable: %s!", executableName)
	}
	return path
}

// executeWithResult executes given command with flags and returns stdout and stderr bufers containing executed command result.
// Additionally, error is returned if command ended with non-zero status
// stdout is nil if error is not-nil and stderr is nil if error is nil
func (r *BaseExecutor) executeWithResult(command string, flags ...string) (error, *bytes.Buffer, *bytes.Buffer) {
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

// execute does same thing as executeWithResult but, only error is returned upon non-zero exist status of command
func (r *BaseExecutor) execute(command string, flags ...string) error {
	err, _, _ := r.executeWithResult(command, flags...)
	return err
}
