package executor

/**
 */

import (
	"bytes"
	"dh/pkg/logging"
)

var (
	logger = logging.GetLoggerFor("executor")
)

type Executor interface {
	executeWithResult(command string, args ...string) (*bytes.Buffer, error)
	execute(command string, args ...string) error
}
