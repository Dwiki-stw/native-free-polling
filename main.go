package main

import (
	"encoding/json"
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

	userRepo := repository.NewUserRepository(db)
	userServ := service.NewUserService(userRepo, helper.BcryptHasher{})
	userHandler := handler.NewUserHandler(userServ)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", http.HandlerFunc(authHandler.Register))
	mux.HandleFunc("/login", http.HandlerFunc(authHandler.Login))
	mux.Handle("/profile", middleware.Auth(conf.JwtKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetProfile(w, r)
		case http.MethodPatch:
			userHandler.UpdateProfile(w, r)
		default:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"code":    "NOT_ALLOWED",
				"message": "method not allowed",
			})
		}
	})))
	mux.Handle("/change-password", middleware.Auth(conf.JwtKey)(http.HandlerFunc(userHandler.ChangePassword)))

	handler := middleware.Recovery(middleware.Logging(mux))

	fmt.Printf("Server Running at http://%s:%s\n", conf.Server.Host, conf.Server.Port)
	http.ListenAndServe(":"+conf.Server.Port, handler)
}
