package utils

import (
	"fmt"
	"reflect"
	"time"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type Dotenv struct {
	PG_DB_URI     string
	SQLITE_DB_URI string
	JWT_SECRET    string
	DB_TYPE       string // "postgres", "sqlite"
	MODE          string // "prod", "dev", "test"
}

type FullConfig struct {
	PG_DB_URI           string
	SQLITE_DB_URI       string
	JWT_SECRET          string
	DB_TYPE             string // "postgres", "sqlite"
	MODE                string // "prod", "dev", "test"
	REFRESH_EXPIRY_TIME time.Duration
	ACCESS_EXPIRY_TIME  time.Duration
}

func parseConfigStruct(s any) []string {
	var fieldNames []string
	v := reflect.ValueOf(s)
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if t.Field(i).Type.Kind() != reflect.String {
			panic(fmt.Sprintf("%v must instead be of type string", t.Field(i).Name))
		}
		if t.Field(i).Type.Kind() == reflect.String {
			fieldNames = append(fieldNames, t.Field(i).Name)
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

var envFile Dotenv
var initialized bool
var config FullConfig

func InitEnv(path string) error {
	g, err := NewEnv(path)
	if err != nil {
		return err
	}
	envFile = g
	if envFile.PG_DB_URI == "" && envFile.SQLITE_DB_URI == "" {
		panic("Neither Postgres or Sqlite3 database URI available")
	}
	if !initialized {
		config = FullConfig{
			PG_DB_URI:           envFile.PG_DB_URI,
			SQLITE_DB_URI:       envFile.SQLITE_DB_URI,
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
