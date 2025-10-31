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
	"strings"
)

func main() {
	conf := config.Get()
	db := database.GetDatabaseConnection(conf.Database)
	defer db.Close()

	authRepo := repository.NewAuth(db)
	authServ := service.NewAuthService(authRepo, conf.JwtKey, helper.BcryptHasher{})
	authHandler := handler.NewAuthHandler(authServ)

	userRepo := repository.NewUserRepository(db)
	userServ := service.NewUserService(userRepo, helper.BcryptHasher{})
	userHandler := handler.NewUserHandler(userServ)

	pollRepo := repository.NewPolling(db)
	optRepo := repository.NewOption(db)
	voteRepo := repository.NewVote(db)
	pollServ := service.NewPolling(db, pollRepo, optRepo, voteRepo)
	pollHandler := handler.NewPolling(pollServ)

	mux := http.NewServeMux()

	mux.HandleFunc("/register", http.HandlerFunc(authHandler.Register))
	mux.HandleFunc("/login", http.HandlerFunc(authHandler.Login))
	mux.Handle("/users/me", middleware.Auth(conf.JwtKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	mux.Handle("/users/me/change-password", middleware.Auth(conf.JwtKey)(http.HandlerFunc(userHandler.ChangePassword)))
	mux.Handle("/users/me/pollings/", middleware.Auth(conf.JwtKey)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		switch {
		case len(parts) == 5 && parts[4] == "creator":
			userHandler.GetUserCreatedPollings(w, r)
			return
		case len(parts) == 5 && parts[4] == "voter":
			userHandler.GetUserVotedPollings(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_FOUND",
			"message": "request not found",
		})
	})))

	mux.Handle("/pollings", middleware.Auth(conf.JwtKey)(http.HandlerFunc(pollHandler.CreatePolling)))
	mux.HandleFunc("/pollings/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")

		if len(parts) == 3 {
			switch r.Method {
			case http.MethodGet:
				pollHandler.GetDetailPolling(w, r)
			case http.MethodPatch:
				middleware.Auth(conf.JwtKey)(http.HandlerFunc(pollHandler.UpdatePolling)).ServeHTTP(w, r)
			case http.MethodDelete:
				middleware.Auth(conf.JwtKey)(http.HandlerFunc(pollHandler.DeletePolling)).ServeHTTP(w, r)
			default:
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusMethodNotAllowed)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    "NOT_ALLOWED",
					"message": "method not allowed",
				})
			}
			return
		}

		switch {
		case len(parts) == 4 && parts[3] == "votes":
			middleware.AuthOptional(conf.JwtKey)(http.HandlerFunc(pollHandler.VoteOptionPolling)).ServeHTTP(w, r)
			return
		case len(parts) == 4 && parts[3] == "results":
			pollHandler.GetPollingResult(w, r)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"code":    "NOT_FOUND",
			"message": "request not found",
		})
	}))

	handler := middleware.Recovery(middleware.Logging(mux))

	fmt.Printf("Server Running at http://%s:%s\n", conf.Server.Host, conf.Server.Port)
	http.ListenAndServe(":"+conf.Server.Port, handler)
}
