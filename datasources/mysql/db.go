package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

// Config is the required properties to use the database.
type Config struct {
	User         string
	Scheme       string
	Password     string
	Host         string
	Name         string
	MaxIdleConns int
	MaxOpenConns int
	DisableTLS   bool
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sql.DB, error) {
	//sslMode := "require"
	//if cfg.DisableTLS {
	//	sslMode = "disable"
	//}

	//q := make(url.Values)
	//q.Set("sslmode", sslMode)
	//q.Set("timezone", "utc")
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Scheme)
	//u := url.URL{
	//	Scheme:   "postgres",
	//	User:     url.UserPassword(cfg.User, cfg.Password),
	//	Host:     cfg.Host,
	//	Path:     cfg.Name,
	//	RawQuery: q.Encode(),
	//}

	db, err := sql.Open("mysql", dataSource)
	if err != nil {
		return nil, err
	}
	//
	//err = db.Ping()
	//if err != nil {
	//	return nil, err
	//}

	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)

	return db, nil
}
