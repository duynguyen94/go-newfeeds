package repo

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

// TODO Load from env
const (
	defaultHost = "127.0.0.1"
	defaultPort = "3306"
	defaultUser = "mysql"
	defaultPass = "mysql"
	defaultName = "newsfeed"
)

func InitMySQLDBConn() (*sql.DB, error) {
	// Specify connection properties.
	cfg := mysql.Config{
		User:   defaultUser,
		Passwd: defaultPass,
		Net:    "tcp",
		Addr:   defaultHost + ":" + defaultPort,
		DBName: defaultName,
	}

	// Get a driver-specific connector.
	connector, err := mysql.NewConnector(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Get a database handle.
	db := sql.OpenDB(connector)

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = db.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
