package executor

/**

TODO: sql execution
	* create first implementation of executing query from template

TODO: sql rendering
	* add option for multiple query parameters

TODO: introspect all tables in selected query and generate names of all possible columns (haha, nice to have feature)


*/

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"os"
	"path/filepath"
	"text/template"
)

var (
	QueryNotFound = errors.New("query not found")
)

type sqlPredicateOperator string

type SqlExecutor struct {
	db *sql.DB

	config SqlExecutorConfig
	// templateData key -> sql script file name, value -> sql script abs path
	templateData map[string]string
}

type SqlExecutorConfig struct {
	// do not change order or delete properties!!! this change will need changes in sqlexecutorview.go
	TemplatesPath      string `yaml:"templatesPath" placeholder:"Path to sql templates"`
	DbConnectionString string `yaml:"dbConnectionString" placeholder:"DB connection string"`
}

type QueryData struct {
	Column   string
	Operator sqlPredicateOperator
	Arg      any
}

func NewSqlExecutor(config SqlExecutorConfig) (*SqlExecutor, error) {
	res := &SqlExecutor{
		config:       config,
		templateData: make(map[string]string),
	}

	db, err := sql.Open("postgres", res.config.DbConnectionString)
	if err != nil {
		return nil, err
	}
	res.db = db
	err = res.loadTemplates()
	if err != nil {
		return nil, err
	}
	return res, nil
}

// RunQuery run specified query and return result set as matrix with first how af header
func (r *SqlExecutor) RunQuery(queryName string, condition QueryData) ([][]string, error) {
	query, err := r.prepareQuery(queryName, &condition)
	if err != nil {
		return nil, err
	}
	logger.Printf("executing query: \n %s", query)
	encodedRes, err := r.executeWithResult(query)
	if err != nil {
		logger.Printf("failed to execute query: %s", err.Error())
		return nil, err
	}
	return r.decodeResultSet(encodedRes, err)
}

func (r *SqlExecutor) ListAvailableTemplates() []string {
	var res = make([]string, len(r.templateData))
	for k, _ := range r.templateData {
		res = append(res, k)
	}
	return res
}

func (r *SqlExecutor) decodeResultSet(encodedRes *bytes.Buffer, err error) ([][]string, error) {
	decoder := gob.NewDecoder(encodedRes)
	res := make([][]string, 0)
	err = decoder.Decode(&res)
	return res, err
}

func (r *SqlExecutor) prepareQuery(queryName string, condition *QueryData) (string, error) {
	if condition == nil {
		return "", errors.New("specify query conditions")
	}

	queryFileName := fmt.Sprint(queryName, ".sql")
	if _, ok := r.templateData[queryFileName]; !ok {
		return "", QueryNotFound
	}
	queryPath := r.templateData[queryFileName]

	tmpl, err := template.ParseFiles(queryPath)
	if err != nil {
		return "", err
	}
	buff := &bytes.Buffer{}
	err = tmpl.Execute(buff, condition)
	if err != nil {
		return "", err
	}
	return buff.String(), nil
}

func (r *SqlExecutor) loadTemplates() error {
	o, err := os.Open(r.config.TemplatesPath)
	if err != nil {
		logger.Printf("missing templatePath configuration!")
		return err
	}
	dir, err := o.ReadDir(0)
	if err != nil {
		logger.Printf("error reading dir: %s", r.config.TemplatesPath)
		return err
	}
	for _, entry := range dir {
		abs := filepath.Join(r.config.TemplatesPath, entry.Name())
		if err != nil {
			logger.Printf("error while converting template path to abs path")
			return err
		}

		r.templateData[entry.Name()] = abs
	}
	return nil
}

// TODO: maybe like this: https://betterprogramming.pub/dynamic-sql-query-with-go-8aeedaa02907
func (r *SqlExecutor) executeWithResult(command string, args ...string) (*bytes.Buffer, error) {
	rows, err := r.db.Query(command)
	if err != nil {
		return nil, err
	}
	var k string
	for rows.Next() {
		err := rows.Scan(&k, &k, &k)
		if err != nil {
			return nil, err
		}
		logger.Printf("got: %s", k)
	}
	res := &bytes.Buffer{}
	encoder := gob.NewEncoder(res)
	tableResult := [][]string{
		{"ID", "TITLE", "AUTHOR_ID"},
		{"1", "The Fellowship of the Ring", "1"},
		{"2", "The Two Towers", "1"},
	}
	err = encoder.Encode(tableResult)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *SqlExecutor) execute(command string, args ...string) error {

	return nil
}
