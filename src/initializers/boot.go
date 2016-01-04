package initializers

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
)

var Logger *logger

type AppConfig struct {
	LogFile     string            `json:"log_file"`
	LogToSyslog bool              `json:"log_to_syslog"`
	Values      map[string]string `json:"values"`
}

type SoaRegistry struct {
	Hermes map[string]string `json:"hermes"`
	Yoda   map[string]string `json:"yoda"`
}

var booted = false
var initialized = false
var config *AppConfig
var registry *SoaRegistry
var logFile *os.File

var AppConfigParseError = errors.New("The app config file did not parse.")
var AlreadyBootedError = errors.New("Attempted to initialize the app when already booted.")

func Boot(configPaths *Paths, logTag string) error {
	if booted {
		return AlreadyBootedError
	}

	if !initialized {
		if err := InitializeApplication(configPaths, logTag); err != nil {
			return err
		}
	}

	booted = true
	return nil
}

func InitializeApplication(configPaths *Paths, logTag string) error {
	if booted || initialized {
		return AlreadyBootedError
	}

	// app configuration
	config = new(AppConfig)
	if err := DecodeConfig(configPaths.ConfigFilePath, config); err != nil {
		return err
	}

	// soa configuration
	registry = new(SoaRegistry)
	if err := DecodeConfig(configPaths.SoaRegistryPath, registry); err != nil {
		return err
	}

	logFilePath := ""
	if config.LogFile == "" {
		logFilePath = configPaths.LogFilePath
	} else {
		logFilePath = config.LogFile
	}

	if err := InitilaizeLogger(logFilePath, logTag, config.LogToSyslog); err != nil {
		return err
	}

	initialized = true

	return nil
}

func (c *AppConfig) WithValue(key string, value string) *AppConfig {
	if c.Values == nil {
		c.Values = make(map[string]string)
	}

	c.Values[key] = value

	return c
}

func (c *AppConfig) Value(key string) string {
	if c.Values == nil {
		return ""
	}

	return c.Values[key]
}

func Config() *AppConfig {
	if config == nil {
		return new(AppConfig)
	} else {
		return config
	}
}

func Registry() *SoaRegistry {
	if registry == nil {
		return new(SoaRegistry)
	} else {
		return registry
	}
}

// func BindStatsd() error {
// 	return BindStatsdAddr(Config().StatsdUrl)
// }

// func BindStatsdAddr(addr string) error {
// 	StatsdClient = statsd.NewStatsdClient(addr, "")
// 	if err := StatsdClient.CreateSocket(); err != nil {
// 		log.Panic(err)
// 		return err
// 	}

// 	return nil
// }

func DecodeConfig(configPath string, conf interface{}) error {
	config = new(AppConfig)

	r, err := readerForFilePath(configPath)
	if err != nil {
		log.Panic(err)
		return err
	}

	if err := json.NewDecoder(r).Decode(conf); err != nil {
		log.Panic(err)
		return AppConfigParseError
	}

	return nil
}

func readerForFilePath(filePath string) (io.Reader, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Panic(err)
		return nil, err
	}

	r := bufio.NewReader(f)

	return r, nil
}

func InitilaizeLogger(logFilePath, logTag string, useSyslog bool) error {
	var logger *logger
	var err error

	if useSyslog {
		logger, err = InitializeSyslogger(logTag)
	} else {
		logger, err = InitializeFileLogger(logFilePath)
	}

	if err != nil {
		return err
	}

	Logger = logger

	return nil
}

func InitializeSyslogger(logTag string) (*logger, error) {
	sysLogger, err := syslog.New(syslog.LOG_ERR|syslog.LOG_SYSLOG, logTag)
	if err != nil {
		return nil, err
	}
	defer sysLogger.Close()

	return NewLogger(nil, sysLogger), nil
}

func InitializeFileLogger(logFilePath string) (*logger, error) {
	fmt.Println("FileLogger initialized in", logFilePath)

	var err error

	logFile, err = os.OpenFile(logFilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	fileLogger := log.New(logFile, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)
	return NewLogger(fileLogger, nil), nil
}
