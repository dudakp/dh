package qh

import (
	"dh/pkg/executor"
	"dh/pkg/logging"
	"fmt"
	"github.com/spf13/viper"
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
	config         executor.SqlExecutorConfig
}

func NewSqlExecutorService(config executor.SqlExecutorConfig) *SqlExecutorService {
	sqlExecutor := executor.NewSqlExecutor(config)
	res := &SqlExecutorService{
		executor: sqlExecutor,
	}

	res.configFilePath = createConfigFile()

	viper.SetConfigFile(res.configFilePath)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(res.configFilePath)

	res.config = loadConfig()

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
	loadedConfig := loadConfig()
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

func (r *SqlExecutorService) Run(queryName string) [][]string {
	return r.executor.RunQuery(queryName)
}

func loadConfig() executor.SqlExecutorConfig {
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	res := &executor.SqlExecutorConfig{}
	err = viper.Unmarshal(res)
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
	configPath = filepath.Join(configPath, "config.yaml")
	_, err = os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return configPath
}
