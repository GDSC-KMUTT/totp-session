package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/GDSC-KMUTT/totp-session/config"
	"github.com/GDSC-KMUTT/totp-session/handler"
	"github.com/GDSC-KMUTT/totp-session/repository"
	"github.com/GDSC-KMUTT/totp-session/service"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	s := &http.Server{
		Addr:           ":8080",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	db, err := sql.Open("mysql", config.C.DB_HOST)
	if err != nil {
		panic(err)
	}
	userRepository := repository.NewRepositoryDB(db)
	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	http.HandleFunc("/signup", userHandler.SignUp)
	http.HandleFunc("/signin", userHandler.SignIn)
	http.HandleFunc("/", userHandler.ListUsers)

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

	defer db.Close()
}
