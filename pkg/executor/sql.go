package executor

/**


TODO: loading and validating all templated sql files
	* scan directory for sql files
	* create list of executable queries that can be used as CLI parameters
	* create hints for all leaded and validated templates

TODO: introspect all tables in selected query and generate names of all possible columns (haha, nice to have feature)


*/

import (
	"bytes"
	"database/sql"
)

type SqlExecutor struct {
	db *sql.DB

	config       SqlExecutorConfig
	templateData []TemplateData
}

type SqlExecutorConfig struct {
	templatesPath      string `yaml:"templatesPath"`
	dbConnectionString string `yaml:"dbConnectionString"`
}

type TemplateData struct {
	Abr  string
	Path string
}

type simplePredicate struct {
	column string
	arg    any
}

func NewSqlExecutor(config SqlExecutorConfig) *SqlExecutor {
	return &SqlExecutor{
		config: config,
	}
}

func (r *SqlExecutor) ListAvailableTemplates() []TemplateData {
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
