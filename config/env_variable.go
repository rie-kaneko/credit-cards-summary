package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Config Configuration

type Configuration struct {
	Environment EnvironmentConfig
	Email       EmailConfig
}

type EnvironmentConfig struct {
	Environment string
	LogLevel    string
	CsvPath     string
}

type EmailConfig struct {
}

type variablesKeys struct {
	envPath,
	logLevel,
	csvPath string
}

func init() {
	keys := setVariablesKeys()
	env := getEnvironment()
	vr := viper.New()

	vr.SetConfigFile(fmt.Sprintf(keys.envPath, env))
	_ = vr.ReadInConfig()

	vr.SetDefault(keys.logLevel, "error")
	vr.SetDefault(keys.csvPath, "resources")

	Config = Configuration{
		Environment: EnvironmentConfig{
			Environment: env,
			LogLevel:    vr.GetString(keys.logLevel),
			CsvPath:     vr.GetString(keys.csvPath),
		},
	}
}

func setVariablesKeys() variablesKeys {
	return variablesKeys{
		envPath:  "./environment/%s.env",
		logLevel: "LOG_LEVEL",
		csvPath:  "CSV_PATH",
	}
}

func getEnvironment() string {
	if value := os.Getenv("ENV"); value != "" {
		return value
	}
	return "local"
}
