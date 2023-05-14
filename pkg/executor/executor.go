package executor

/**
TODO: create executor types
	* sqlExecutor - calling raw or templated sql scripts (control r/w access)
*/

import (
	"bytes"
	"dh/pkg/logging"
	"fmt"
	"os/exec"
)

var (
	logger = logging.GetLoggerFor("executor")
)

type Executor interface {
	executeWithResult(command string, flags ...string) (*bytes.Buffer, error)
	execute(command string, flags ...string) error
}

type FileExecutor struct {
	executableName string
	path           string
}

func newFileExecutor(executableName string) FileExecutor {
	return FileExecutor{
		executableName: executableName,
		path:           findExecutable(executableName),
	}
}

func findExecutable(executableName string) string {
	path, err := exec.LookPath(executableName)
	if err != nil {
		logger.Fatalf("Unable to find executable: %s!", executableName)
	}
	return path
}

// executeWithResult executes given command with flags and returns stdout containing executed command result.
func (r *FileExecutor) executeWithResult(command string, flags ...string) (stdout *bytes.Buffer, err error) {
	argWithFlags := []string{command}
	argWithFlags = append(argWithFlags, flags...)
	cmd := exec.Command(r.path, argWithFlags...)
	logger.Printf("executing command: %s", cmd.String())
	out, err := cmd.CombinedOutput()
	stdout = bytes.NewBuffer(out)
	if err != nil {
		logger.Printf(fmt.Errorf("%w: %s", err, stdout.String()).Error())
		return stdout, fmt.Errorf("\n\t%w: unable to execute git command: %s\n\tcall ended with error: %s",
			err, argWithFlags, stdout.String(),
		)
	}
	logger.Print(stdout.String())
	return
}

// execute does same thing as executeWithResult but, only error is returned upon non-zero exist status of command
func (r *FileExecutor) execute(command string, flags ...string) error {
	_, err := r.executeWithResult(command, flags...)
	return err
}
