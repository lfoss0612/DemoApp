// Package env provides environment configuration.
package env

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var EnvVar *Config

// Get retrieves an environment variable or uses a default value.
func Get(key, def string) string {
	if env, ok := os.LookupEnv(key); ok {
		return env
	}
	return def
}

// GetList splits an environment variable into a slice of strings.
func GetList(key, def string) []string {
	return strings.Split(Get(key, def), ",")
}

func Init() {
	var err error
	EnvVar = &Config{}

	EnvVar.Env = os.Getenv("env")
	if EnvVar.Env == "" {
		EnvVar.Env = "local"
	}

	EnvVar.Port = os.Getenv("port")
	if EnvVar.Port == "" {
		EnvVar.Port = "9090"
	}

	EnvVar.Region = os.Getenv("region")
	if EnvVar.Region == "" {
		EnvVar.Region = "us-east-1"
	}

	EnvVar.LogLevel = os.Getenv("LOG_LEVEL")
	if EnvVar.LogLevel == "" {
		EnvVar.LogLevel = "INFO"
	}

	EnvVar.LogFile = os.Getenv("LOG_FILE")
	if EnvVar.LogFile == "" {
		EnvVar.LogFile = "application.log"
	}

	EnvVar.LogAppender = os.Getenv("LOG_APPENDER")
	if EnvVar.LogAppender == "" {
		EnvVar.LogAppender = "console"
	}

	EnvVar.ServerReadTimeoutInSeconds, err = time.ParseDuration(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		log.Println("Unable to read \"SERVER_READ_TIMEOUT\" from request.")
	}
	if EnvVar.ServerReadTimeoutInSeconds == 0 {
		EnvVar.ServerReadTimeoutInSeconds, _ = time.ParseDuration("10s")
	}

	EnvVar.ServerMaxSimultaneousConnections, err = strconv.Atoi(os.Getenv("SERVER_MAX_CONNECTIONS"))
	if err != nil {
		log.Println("Unable to read \"SERVER_MAX_CONNECTIONS\" from request.")
	}
	if EnvVar.ServerMaxSimultaneousConnections == 0 {
		EnvVar.ServerMaxSimultaneousConnections = 5000
	}
}

type Config struct {
	Port string `default:"8080"`
	Env  string `default:"local"`

	Region                           string        `default:"us-east-1"`
	LogLevel                         string        `default:"INFO"`
	LogFile                          string        `default:"application.log"`
	LogAppender                      string        `default:"console"`
	ServerReadTimeoutInSeconds       time.Duration `default:"10s"`
	ServerMaxSimultaneousConnections int           `default:"5000"`
	BuildInfo                        Build
}

func (conf Config) ToString() string {
	return fmt.Sprintf("Port: [%s], Environment: [%s], LogLevel: [%s]",
		conf.Port, conf.Env, conf.LogLevel)
}
