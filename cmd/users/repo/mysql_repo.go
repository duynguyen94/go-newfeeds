package repo

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"time"
)

const (
	dbhost = "127.0.0.1"
	dbport = "3306"
	dbuser = "mysql"
	dbpass = "mysql"
	dbname = "newsfeed"
)

func InitMySQLDBConn() (*sql.DB, error) {
	// Specify connection properties.
	cfg := mysql.Config{
		User:   dbuser,
		Passwd: dbpass,
		Net:    "tcp",
		Addr:   dbhost + ":" + dbport,
		DBName: dbname,
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
