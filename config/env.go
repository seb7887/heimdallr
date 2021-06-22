package config

import (
	"os"
	"reflect"
	"strconv"
	"strings"
	"sync"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	HealthPort int    `env:"HEALTH_PORT" default:"7000" json:"healthPort"`
	GRPCPort   int    `env:"GRPC_PORT" default:"7001" json:"grpcPort"`
	RedisHost  string `env:"REDIS_HOST" default:":6379" json:"redisHost"`
}

var _config *Config
var _once sync.Once

func GetConfig() *Config {
	_once.Do(func() {
		Setup()
	})
	return _config
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Debug("Could not load .env file, using environment variables")
	}
}

func Setup() {
	loadEnv()
	_config = &Config{}

	configReflect := reflect.ValueOf(_config).Elem()
	err := loadConfig(configReflect, configReflect.Type())

	if err != nil {
		log.Error("error in reading Heimdallr configurations.")
	}
}

func internalField(fieldDef reflect.StructField) bool {
	return "true" == fieldDef.Tag.Get("internal")
}

func requiredField(fieldDef reflect.StructField) bool {
	return "true" == fieldDef.Tag.Get("required")
}

func loadConfigField(field reflect.Value, fieldDef reflect.StructField) error {
	var err error
	configField := fieldDef.Tag.Get("env")
	defaultValue := fieldDef.Tag.Get("default")
	configValue := os.Getenv(configField)
	if len(configValue) == 0 {
		configValue = defaultValue
	}
	if len(configValue) == 0 {
		if requiredField(fieldDef) {
			log.Printf("Field %s missing required configuration", configField)
		}
	} else {
		switch field.Type().Kind() {
		case reflect.Slice:
			values := strings.Split(configValue, ",")
			field.Set(reflect.ValueOf(values))
			if !internalField(fieldDef) {
				log.Debug("Loaded configuration")
			}
		case reflect.String:
			field.SetString(configValue)
			if !internalField(fieldDef) {
				log.Debug("Loaded configuration")
			}
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(configValue)
			if err != nil {
				log.Error("Invalid configuration")
			} else {
				field.SetBool(boolValue)
				if !internalField(fieldDef) {
					log.Debug("Loaded configuration")
				}
			}
		case reflect.Int:
			intValue, err := strconv.Atoi(configValue)
			if err != nil {
				log.Error("Invalid configuration")
			} else {
				// Make sure the configured value meets the minimum requirements
				field.SetInt(int64(intValue))
				if !internalField(fieldDef) {
					log.Debug("Loaded configuration")
				}
			}
		}
	}

	return err
}

func loadConfig(configValue reflect.Value, configValueType reflect.Type) error {
	var err error
	for i := 0; i < configValue.NumField(); i++ {
		field := configValue.Field(i)
		if field.Type().Kind() == reflect.Struct {
			err = loadConfig(field, field.Type())
		} else {
			err = loadConfigField(field, configValueType.Field(i))
		}
	}
	return err
}
