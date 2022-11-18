package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/GDSC-KMUTT/totp-session/config"
	"github.com/GDSC-KMUTT/totp-session/service"
	"github.com/GDSC-KMUTT/totp-session/types"
	"github.com/GDSC-KMUTT/totp-session/utils"
	"github.com/golang-jwt/jwt"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(userSerivice service.UserService) userHandler {
	return userHandler{service: userSerivice}
}

func (h userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")
	var body types.SignIn
	err := utils.Parse(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Call signup service
	id, base64, secret, err := h.service.SignUp(body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create a response
	response, _ := json.Marshal(map[string]any{"success": true, "id": id, "image": base64, "secret": secret})
	w.Write(response)
	return
}

func (h userHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var body types.SignIn
	err := utils.Parse(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	id, err := h.service.SignIn(body.Email, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, _ := json.Marshal(map[string]any{"success": true, "id": id})
	w.Write(response)
	return
}

type CustomClaims struct {
	Id int64 `json:"id"`
	jwt.StandardClaims
}

func (h userHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	bearer := r.Header.Get("Authorization")
	jwtToken := strings.Split(bearer, " ")
	token, err := jwt.ParseWithClaims(jwtToken[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT_SECRET), nil
	})
	claims, ok := token.Claims.(*CustomClaims)
	if !ok && !token.Valid {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !token.Valid {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUser(claims.Id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(map[string]any{"user": user})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
	return
}

func (h userHandler) ConfirmOtp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var body types.ConfirmSignUp
	err := utils.Parse(r, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	token, err := h.service.ConfirmOtp(body.Id, body.Otp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, _ := json.Marshal(map[string]any{"success": true, "token": &token})
	w.Write(response)
}
