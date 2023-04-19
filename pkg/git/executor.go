package git

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
)

type Executor struct {
	path string
}

func NewGitExecutor(gitPath string) *Executor {
	return &Executor{path: gitPath}
}

func (r *Executor) Status() error {
	err, stdout, stderr := r.execute("status")
	if err != nil {
		log.Fatalf(fmt.Errorf("%w: %s", err, stderr.String()).Error())
		return err
	}
	log.Print(stdout.String())
	return nil
}

func (r *Executor) execute(arg string, flags ...string) (error, *bytes.Buffer, *bytes.Buffer) {
	cmd := exec.Command(r.path, arg)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err, nil, &stderr
	}
	return nil, &stdout, nil
}
