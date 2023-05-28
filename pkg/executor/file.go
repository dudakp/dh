package executor

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

var (
	ErrExecutableNotFound = errors.New("executable not found")
)

type FileExecutor struct {
	executableName string
	path           string
}

// TODO: add sensible file mode check - non-executable file will return error and add unit test
func newFileExecutor(executableName string) (*FileExecutor, error) {
	path, err := findExecutable(executableName)
	if err != nil {
		return nil, err
	}
	return &FileExecutor{
		executableName: executableName,
		path:           path,
	}, nil
}

func findExecutable(executableName string) (string, error) {
	path, err := exec.LookPath(executableName)
	if err != nil {
		if errors.Is(err, exec.ErrNotFound) {
			logger.Printf("executable %s not found", executableName)
			return "", ErrExecutableNotFound
		}
		return "", err
	}
	return path, nil
}

// executeWithResult executes given command with flags and returns stdout containing executed command result.
func (r *FileExecutor) executeWithResult(command string, args ...string) (stdout *bytes.Buffer, err error) {
	argWithFlags := []string{command}
	argWithFlags = append(argWithFlags, args...)
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
func (r *FileExecutor) execute(command string, args ...string) error {
	_, err := r.executeWithResult(command, args...)
	return err
}
