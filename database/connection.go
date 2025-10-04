package database

import (
	"database/sql"
	"fmt"
	"log"
	"native-free-pollings/config"

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

	defer db.Close()

	fmt.Println("Database connected")

	return db
}
