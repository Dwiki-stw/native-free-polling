package main

import (
	"fmt"
	"native-free-pollings/config"
	"native-free-pollings/database"
	"net/http"
)

func main() {
	conf := config.Get()

	db := database.GetDatabaseConnection(conf.Database)

	db.Ping()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, World")
	})
	fmt.Println("Server Running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
