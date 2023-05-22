package executor

/**

TODO: reduce number of panics

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
	"errors"
	"os"
	"path/filepath"
	"sort"
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
	// Abr template name abbreviation - used ad argument for invoking sql query
	Abr string
	// Path to sql script file
	Path string
}

type simplePredicate struct {
	column string
	arg    any
}

func NewSqlExecutor(config SqlExecutorConfig) *SqlExecutor {
	res := &SqlExecutor{
		config: config,
	}
	res.loadTemplates()
	return res
}

func (r *SqlExecutor) PrepareQuery(queryName string, condition *simplePredicate) (string, error) {
	if condition == nil {
		return "", errors.New("specify query conditions")
	}
	templateIndex := sort.Search(len(r.templateData), func(i int) bool {
		return r.templateData[i].Abr == queryName
	})
	template := r.templateData[templateIndex]

	templateFile, err := os.OpenFile(template.Path, os.O_RDONLY, 0644)
	// TODO: add go templating
}

// RunQuery run specified query and return result set as matrix with first how af header
func (r *SqlExecutor) RunQuery(queryName string) ([][]string, error) {
	_, err := r.executeWithResult("")
	if err != nil {
		logger.Printf("failed to execute query: %w", err)
		return nil, err
	}
	return nil, nil
}

func (r *SqlExecutor) ListAvailableTemplates() []TemplateData {
	return r.templateData
}

func (r *SqlExecutor) loadTemplates() {
	o, err := os.Open(r.config.TemplatesPath)
	if err != nil {
		logger.Println("missing templatePath configuration!")
		return
	}
	dir, err := o.ReadDir(0)
	if err != nil {
		panic(err)
	}
	for _, entry := range dir {
		abs, err := filepath.Abs(entry.Name())
		if err != nil {
			panic(err)
		}

		r.templateData = append(r.templateData, TemplateData{
			Abr:  entry.Name(),
			Path: abs,
		})
	}
}

func (r *SqlExecutor) executeWithResult(command string, flags ...string) (*bytes.Buffer, error) {

	return nil, nil
}

func (r *SqlExecutor) execute(command string, flags ...string) error {

	return nil
}
