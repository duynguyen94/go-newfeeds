package conn

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"log"
	"os"
	"time"
)

func InitMySQLDBConn() (*sql.DB, error) {
	// Specify connection properties.
	cfg := mysql.Config{
		User:   os.Getenv("MYSQL_USER"),
		Passwd: os.Getenv("MYSQL_PASS"),
		Net:    "tcp",
		Addr:   os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"),
		DBName: os.Getenv("MYSQL_DATABASE"),
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
