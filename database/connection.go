package database

import (
	"database/sql"
	"fmt"
	"log"
	"native-free-pollings/config"
	"time"

	_ "github.com/lib/pq"
)

func GetDatabaseConnection(conf config.Database) *sql.DB {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		conf.User, conf.Pass, conf.Host, conf.Port, conf.Name, conf.SSL,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed open conection:", err)
	}

	fmt.Println("Database connected")

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(3 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	return db
}
