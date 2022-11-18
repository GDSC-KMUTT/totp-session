package main

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"image/png"
	"net/http"
	"time"

	"github.com/GDSC-KMUTT/totp-session/config"
	"github.com/GDSC-KMUTT/totp-session/handler"
	"github.com/GDSC-KMUTT/totp-session/repository"
	"github.com/GDSC-KMUTT/totp-session/service"
	"github.com/GDSC-KMUTT/totp-session/types"
	"github.com/GDSC-KMUTT/totp-session/utils"
	_ "github.com/go-sql-driver/mysql"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
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

	http.HandleFunc("/", (func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		// Set the response header to application/json
		w.Header().Set("Content-Type", "application/json")

		// Declare a variable to store the body of the request
		var response []byte
		var body types.SignUp
		err := utils.Parse(r, &body)
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusBadRequest)
			w.Write(response)
			return
		}

		// Generate a new secret TOTP key
		key, err := totp.Generate(totp.GenerateOpts{
			Issuer:      "GDSC KMUTT",
			AccountName: body.Email,
		})
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}
		secret := key.Secret()

		// Hash the password
		hashedPwd, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}

		// Create a new user
		insert, err := db.Exec("INSERT INTO users (email, password, secret) VALUES (?, ?, ?)", body.Email, hashedPwd, secret)
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}
		userId, err := insert.LastInsertId()
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}

		// Convert TOTP key into a PNG, and encode it to base64
		var buf bytes.Buffer
		img, err := key.Image(200, 200)
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}
		if err := png.Encode(&buf, img); err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}
		base64string := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
		if err != nil {
			response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(response)
			return
		}
		url := key.URL()

		// Create a response
		response, _ = json.Marshal(map[string]any{"success": true, "id": userId, "image": base64string, "secret": url})
		w.Write(response)
	}))

	http.HandleFunc("/signup", userHandler.SignUp)
	http.HandleFunc("/signin", userHandler.SignIn)
	http.HandleFunc("/confirm-otp", userHandler.ConfirmOtp)
	http.HandleFunc("/get-user", userHandler.GetProfile)

	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}

	defer db.Close()
}
