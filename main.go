package main

import (
	"fmt"
	"native-free-pollings/config"
	"native-free-pollings/database"
	"native-free-pollings/handler"
	"native-free-pollings/helper"
	"native-free-pollings/middleware"
	"native-free-pollings/repository"
	"native-free-pollings/service"
	"net/http"
)

func main() {
	conf := config.Get()
	db := database.GetDatabaseConnection(conf.Database)

	authRepo := repository.NewAuth(db)
	authServ := service.NewAuthService(authRepo, conf.JwtKey, helper.BcryptHasher{})
	authHandler := handler.NewAuthHandler(authServ)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", http.HandlerFunc(authHandler.Register))
	mux.HandleFunc("/login", http.HandlerFunc(authHandler.Login))

	handler := middleware.Recovery(middleware.Logging(mux))

	fmt.Printf("Server Running at http://%s:%s\n", conf.Server.Host, conf.Server.Port)
	http.ListenAndServe(":"+conf.Server.Port, handler)
}
