package executor

/**

TODO: loading db connection data from config file
	* on first usage

TODO: introspect all tables in selected query and generate names of all possible columns


*/

import (
	"bytes"
	"database/sql"
)

type SqlExecutor struct {
	db            *sql.DB
	templatesPath string
	templateData  []templateData
}

type templateData struct {
	Abr  string
	Path string
}

type simplePredicate struct {
	column string
	arg    any
}

func NewSqlExecutor(templatesPath string) *SqlExecutor {
	return &SqlExecutor{
		templatesPath: templatesPath,
	}
}

func (r *SqlExecutor) ListAvailableTemplates() []templateData {
	return r.templateData
}

func (r *SqlExecutor) loadTemplates() {

}

func (r *SqlExecutor) executeWithResult(command string, flags ...string) (*bytes.Buffer, error) {

	return nil, nil
}

func (r *SqlExecutor) execute(command string, flags ...string) error {

	return nil
}
