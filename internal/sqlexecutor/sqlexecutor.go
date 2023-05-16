package sqlexecutor

/**
TODO: loading persistent configuration
	* create on first usage
	* load if config exists
	* provide interface for rewriting config
*/

import (
	"dh/pkg/executor"
	"dh/pkg/logging"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logger = logging.GetLoggerFor("flow")
)

type SqlExecutorService struct {
	executor       *executor.SqlExecutor
	configFilePath string
}

func NewSqlExecutorService(config executor.SqlExecutorConfig) *SqlExecutorService {
	configPath := createConfigFile()
	loadConfig(configPath)

	sqlExecutor := executor.NewSqlExecutor(config)

	return &SqlExecutorService{
		executor: sqlExecutor,
	}
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

}

func loadConfig(configPath string) executor.SqlExecutorConfig {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
	var res executor.SqlExecutorConfig
	err = viper.Unmarshal(res)
	if err != nil {
		logger.Printf("failed to unmarshall config file! using default")
		return executor.SqlExecutorConfig{}
	}

	return res
}

func createConfigFile() string {
	// create config file
	var configPath string
	if runtime.GOOS == "windows" {
		configPath = filepath.Join(os.Getenv("HOMEDRIVE"), os.Getenv("HOMEPATH"), "AppData", "dh")
	} else {
		configPath = filepath.Join("~./dh")
	}
	configPath = filepath.Join(configPath, "config.yaml")
	_, err := os.OpenFile(configPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	return configPath
}
