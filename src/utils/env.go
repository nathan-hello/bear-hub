package utils

import (
	"database/sql"
	"fmt"
	"reflect"
	"time"
	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
	"github.com/nathan-hello/htmx-template/src/db"
)

type ContextClaimType string

const ClaimsContextKey ContextClaimType = "claims"

type Dotenv struct {
	DB_URI     string
	JWT_SECRET string
}

type FullConfig struct {
	DB_URI              string
	JWT_SECRET          string
	REFRESH_EXPIRY_TIME time.Duration
	ACCESS_EXPIRY_TIME  time.Duration
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

func NewEnv() Dotenv {
	requiredEnvVars := parseConfigStruct(Dotenv{})

	dotenv, err := godotenv.Read(".env")
	if err != nil {
		panic(err)
	}

	return parseRequiredEnvVars(dotenv, requiredEnvVars)

}

var g Dotenv = NewEnv()
var C = FullConfig{
	DB_URI:              g.DB_URI,
	JWT_SECRET:          g.JWT_SECRET,
	REFRESH_EXPIRY_TIME: time.Hour * 72,
	ACCESS_EXPIRY_TIME:  time.Hour * 24,
}

func Env() *FullConfig {
	return &C
}

var d, err = sql.Open("postgres", Env().DB_URI)

func Db() (*db.Queries, error) {
	if err != nil {
		return nil, err
	}
	return db.New(d), nil
}