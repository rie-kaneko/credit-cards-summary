package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

var Config Configuration

type Configuration struct {
	Environment EnvironmentConfig
	AWS         AWSConfig
}

type EnvironmentConfig struct {
	Environment string
	LogLevel    string
	CsvPath     string
	EmailSender string
}

type AWSConfig struct {
	Region string
}

type variablesKeys struct {
	envPath,
	logLevel,
	csvPath,
	emailSender,
	awsRegion string
}

func init() {
	keys := setVariablesKeys()
	env := getEnvironment()
	vr := viper.New()

	vr.SetConfigFile(fmt.Sprintf(keys.envPath, env))
	_ = vr.ReadInConfig()

	vr.SetDefault(keys.logLevel, "error")
	vr.SetDefault(keys.csvPath, "resources")
	vr.SetDefault(keys.emailSender, "")
	vr.SetDefault(keys.awsRegion, "us-east-1")

	Config = Configuration{
		Environment: EnvironmentConfig{
			Environment: env,
			LogLevel:    vr.GetString(keys.logLevel),
			CsvPath:     vr.GetString(keys.csvPath),
			EmailSender: vr.GetString(keys.emailSender),
		},
		AWS: AWSConfig{
			Region: vr.GetString(keys.awsRegion),
		},
	}
}

func setVariablesKeys() variablesKeys {
	return variablesKeys{
		envPath:     "./environment/%s.env",
		logLevel:    "LOG_LEVEL",
		csvPath:     "CSV_PATH",
		emailSender: "EMAIL_SENDER",
		awsRegion:   "AWS_REGION",
	}
}

func getEnvironment() string {
	if value := os.Getenv("ENV"); value != "" {
		return value
	}
	return "local"
}
