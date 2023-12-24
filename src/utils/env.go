package utils

import (
	"fmt"
	"reflect"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_URI string
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

func parseRequiredEnvVars(dotenv map[string]string, fields []string) Config {
	var c Config
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

func NewEnv() Config {
	requiredEnvVars := parseConfigStruct(Config{})

	dotenv, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	return parseRequiredEnvVars(dotenv, requiredEnvVars)

}

var GlobalConfig Config = NewEnv()

func Env() *Config {
	return &GlobalConfig
}
