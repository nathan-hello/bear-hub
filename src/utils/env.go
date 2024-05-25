package utils

import (
	"fmt"
	"reflect"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

var envFile Dotenv
var initialized bool
var config FullConfig

type Dotenv struct {
	DB_URI     string
	JWT_SECRET string
}

type FullConfig struct {
	DB_URI              string
	JWT_SECRET          string
	REFRESH_EXPIRY_TIME time.Duration
	ACCESS_EXPIRY_TIME  time.Duration
	MODE                string // "prod", "dev", "test"
}

func InitEnv(path string) error {
	g, err := NewEnv(path)
	if err != nil {
		return err
	}
	envFile = g
	if envFile.DB_URI == "" || envFile.JWT_SECRET == "" {
		panic("DB_URI or JWT_SECRET uninitiated")
	}
	if !initialized {
		config = FullConfig{
			DB_URI:              envFile.DB_URI,
			JWT_SECRET:          envFile.JWT_SECRET,
			REFRESH_EXPIRY_TIME: time.Hour * 72,
			ACCESS_EXPIRY_TIME:  time.Hour * 24,
			MODE:                "dev",
		}
		initialized = true
	}
	return nil
}

func Env() *FullConfig {
	return &config
}

func parseConfigStruct(s interface{}) []string {
	var fieldNames []string
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Struct {
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if t.Field(i).Type.Kind() != reflect.String {
				panic(fmt.Sprintf("%v must instead be of type string", t.Field(i).Name))
			}
			if t.Field(i).Type.Kind() == reflect.String {
				fieldNames = append(fieldNames, t.Field(i).Name)
			}
		}
	}
	return fieldNames
}

func parseRequiredEnvVars(dotenv map[string]string, fields []string) Dotenv {
	var c Dotenv
	v := reflect.ValueOf(&c).Elem()
	for _, field := range fields {
		if value, ok := dotenv[field]; ok {
			f := v.FieldByName(field)
			if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
				f.SetString(value)
			}
		}
	}
	return c

}

func NewEnv(path string) (Dotenv, error) {
	requiredEnvVars := parseConfigStruct(Dotenv{})

	dotenv, err := godotenv.Read(path)
	if err != nil {
		return Dotenv{}, err
	}

	return parseRequiredEnvVars(dotenv, requiredEnvVars), nil

}
