package executor

/**
TODO: create executor types
	* sqlExecutor - calling raw or templated sql scripts (control r/w access)
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
