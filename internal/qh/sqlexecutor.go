package qh

import (
	"dh/pkg/executor"
	"dh/pkg/logging"
	"errors"
)

/**
TODO: reduce number of panics
*/

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
)

var (
	logger = logging.GetLoggerFor("sqlexecutorservice")

	errUnableToCreateSqlExecutor = errors.New("unable to create sqlExecutor")
)

type SqlExecutorService struct {
	executor       *executor.SqlExecutor
	configFilePath string
}

func NewSqlExecutorService() (*SqlExecutorService, error) {
	res := &SqlExecutorService{}

	res.configFilePath = createConfigFile()

	config := loadConfig(res.configFilePath)
	sqlExecutor, err := executor.NewSqlExecutor(config)
	if err != nil {
		return nil, errors.Join(errUnableToCreateSqlExecutor, err)
	}

	res.executor = sqlExecutor

	return res, nil
}

func (r *SqlExecutorService) ConfigIsEmpty() bool {
	stat, err := os.Stat(r.configFilePath)
	if err != nil {
		logger.Printf("unable to get config file %s info!", r.configFilePath)
		return true
	} else {
		return stat.Size() == 0
	}
}

func (r *SqlExecutorService) WriteConfig(config executor.SqlExecutorConfig) error {
	loadedConfig := loadConfig(r.configFilePath)
	// check config deltas, update only changed values
	if loadedConfig.DbConnectionString != config.DbConnectionString {
		loadedConfig.DbConnectionString = config.DbConnectionString
	}
	if loadedConfig.TemplatesPath != config.TemplatesPath {
		loadedConfig.TemplatesPath = config.TemplatesPath
	}

	yamlConf, err := yaml.Marshal(loadedConfig)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(r.configFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	_, err = file.Write(yamlConf)
	return err
}

func (r *SqlExecutorService) ListAvailableQueries() []string {
	return r.executor.ListAvailableTemplates()
}

func (r *SqlExecutorService) Run(queryName string, queryData executor.QueryData) ([][]string, error) {
	return r.executor.RunQuery(queryName, queryData)
}

func loadConfig(configFilePath string) executor.SqlExecutorConfig {
	file, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	res := &executor.SqlExecutorConfig{}
	err = yaml.Unmarshal(file, res)
	if err != nil {
		logger.Printf("failed to unmarshall config file! using default")
		return executor.SqlExecutorConfig{}
	}
	return *res
}

func createConfigFile() string {
	// create config file
	var configPath string
	osUser, err := user.Current()
	if err != nil {
		panic("unable to get user")
	}
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(osUser.HomeDir, "AppData", "Local", "dh")
	} else {
		configPath = filepath.Join(osUser.HomeDir, ".dh")
	}
	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	configPath = filepath.Join(configPath, "qh.yaml")
	_, err = os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return configPath
}
