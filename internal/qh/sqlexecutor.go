package qh

import (
	"dh/pkg/executor"
	"dh/pkg/logging"
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
)

type SqlExecutorService struct {
	executor       *executor.SqlExecutor
	configFilePath string
}

func NewSqlExecutorService() *SqlExecutorService {
	res := &SqlExecutorService{}

	res.configFilePath = createConfigFile()

	config := loadConfig(res.configFilePath)
	sqlExecutor := executor.NewSqlExecutor(config)

	res.executor = sqlExecutor

	return res
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

func (r *SqlExecutorService) WriteConfig(config executor.SqlExecutorConfig) {
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
		panic(err)
	}
	file, err := os.OpenFile(r.configFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	_, err = file.Write(yamlConf)
	if err != nil {
		panic(err)
	}
}

func (r *SqlExecutorService) ListAvailableQueries() []executor.TemplateData {
	return r.executor.ListAvailableTemplates()
}

func (r *SqlExecutorService) Run(queryName string) ([][]string, error) {
	return r.executor.RunQuery(queryName)
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
