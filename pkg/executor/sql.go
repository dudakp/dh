package executor

/**

TODO: sql execution
	* create first implementation of executing query from template

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
	// do not change order or delete properties!!! this change will need changes in sqlexecutorview.go
	TemplatesPath      string `yaml:"templatesPath" placeholder:"Path to sql templates"`
	DbConnectionString string `yaml:"dbConnectionString" placeholder:"DB connection string"`
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

// RunQuery run specified query and return result set as matrix with first how af header
func (r *SqlExecutor) RunQuery(queryName string) [][]string {
	n := 2
	table := make([][]string, n)
	for i := 0; i < n; i++ {
		if i == 0 {
			table[i] = []string{"ID", "TITLE", "AUTHOR"}
		} else {
			table[i] = []string{"1", "Return of the king", "J. R. R. Tolkien"}
		}
	}
	return table
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
