package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)

var (
	//UserClient предоставляет доступ в БД схеме users
	UserClient *sql.DB
)

func init() {
	mysqlHost, found := os.LookupEnv("mysql_host")
	if !found {
		mysqlHost = os.Getenv("heroku_host")
	}
	mysqlPassword, found := os.LookupEnv("mysql_password")
	if !found {
		mysqlPassword = os.Getenv("heroku_password")
	}
	mysqlUser, found := os.LookupEnv("mysql_user")
	if !found {
		mysqlUser = os.Getenv("heroku_user")
	}
	mysqlSchema, found := os.LookupEnv("mysql_schema")
	if !found {
		mysqlSchema = os.Getenv("heroku_schema")
	}

	var err error
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlSchema)
	UserClient, err = sql.Open("mysql", dataSource)
	if err != nil {
		panic(err)
	}

	if err = UserClient.Ping(); err != nil {
		panic(err)
	}
}
