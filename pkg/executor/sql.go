package executor

/**

TODO: dynamic sql driver loading
	idea: install several drivers and upon creating instance of SqlExecutor
	register specified driver using: sql.Register(driverName, driverObject)

TODO: sql rendering
	* add option for multiple query parameters
	* maybe try to use prepared statements - use ? or named parameters in query, upon loading introspect loaded query and provide hints for CLI

TODO: add dynamic sql driver loading, maybe as plugin? idk, research this

TODO: introspect all tables in selected query and generate names of all possible columns (haha, nice to have feature)

*/

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	_ "github.com/microsoft/go-mssqldb"
	"os"
	"path/filepath"
	"text/template"
)

var (
	ErrQueryNotFound = errors.New("query not found")
)

type sqlPredicateOperator string

type SqlExecutor struct {
	db *sql.DB

	config SqlExecutorConfig
	// templateData key -> sql script file name, value -> sql script abs path
	templateData map[string]string
}

type SqlExecutorConfig struct {
	// do not change order or delete properties!!! this change will need changes in configview.go
	TemplatesPath      string `yaml:"templatesPath" placeholder:"Path to sql templates"`
	DbConnectionString string `yaml:"dbConnectionString" placeholder:"DB connection string"`
	DbVendor           string `yaml:"dbVendor" placeholder:"Database vendor"`
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

	db, err := sql.Open(res.config.DbVendor, res.config.DbConnectionString)
	if err != nil {
		return nil, err
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			logger.Printf("error closing database connection")
			os.Exit(1)
		}
	}(db)
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
	for k := range r.templateData {
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
		return "", ErrQueryNotFound
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

func (r *SqlExecutor) executeWithResult(command string, args ...string) (*bytes.Buffer, error) {
	statement, err := r.db.Prepare(command)
	rows, err := statement.Query(args)
	if err != nil {
		return nil, err
	}
	header, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// real representation of data in row as []byte
	row := make([][]byte, len(header))
	// pointers to all columns returned by query
	pRow := make([]any, len(header))

	resultSet := make([][]string, 0, 1)
	resultSet = append(resultSet, header)

	// set pointers
	for j := range pRow {
		pRow[j] = &row[j]
	}
	i := 1
	for rows.Next() {
		err = rows.Scan(pRow...)
		if err != nil {
			i++
			break
		}
		// allocate space for new row
		resultSet = append(resultSet, make([]string, 0, len(header)))
		for _, c := range row {
			// put col data to row
			resultSet[i] = append(resultSet[i], string(c))
		}
		i++
	}
	// TODO: think how to impose type safety
	res := &bytes.Buffer{}
	encoder := gob.NewEncoder(res)
	err = encoder.Encode(resultSet)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *SqlExecutor) execute(command string, args ...string) error {
	return nil
}
