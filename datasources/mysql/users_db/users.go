package users_db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
	"time"
)

var (
	//UserClientMaster предоставляет доступ в БД схеме users
	UserClientMaster *sql.DB
	UserClientSlave1 *sql.DB
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

	mysqlHostSlave, found := os.LookupEnv("mysql_host_slave")
	if !found {
		mysqlHostSlave = os.Getenv("heroku_host_slave")
	}
	mysqlPasswordSlave, found := os.LookupEnv("mysql_password_slave")
	if !found {
		mysqlPasswordSlave = os.Getenv("heroku_password_slave")
	}
	mysqlUserSlave, found := os.LookupEnv("mysql_user_slave")
	if !found {
		mysqlUserSlave = os.Getenv("heroku_user_slave")
	}
	mysqlSchemaSlave, found := os.LookupEnv("mysql_schema_slave")
	if !found {
		mysqlSchemaSlave = os.Getenv("heroku_schema_slave")
	}

	var err error
	dataSourceMaster := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlSchema)
	UserClientMaster, err = sql.Open("mysql", dataSourceMaster)
	if err != nil {
		panic(err)
	}
	UserClientMaster.SetMaxOpenConns(10)
	UserClientMaster.SetMaxIdleConns(5)
	UserClientMaster.SetConnMaxLifetime(5 * time.Second)

	dataSourceSlave := fmt.Sprintf("%s:%s@tcp(%s)/%s", mysqlUserSlave, mysqlPasswordSlave, mysqlHostSlave, mysqlSchemaSlave)
	UserClientSlave1, err = sql.Open("mysql", dataSourceSlave)
	if err != nil {
		panic(err)
	}
	UserClientSlave1.SetMaxOpenConns(10)
	UserClientSlave1.SetMaxIdleConns(5)
	UserClientSlave1.SetConnMaxLifetime(5 * time.Second)

	if err = UserClientSlave1.Ping(); err != nil {
		panic(err)
	}
}
