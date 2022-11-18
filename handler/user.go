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
	var body types.SignIn
	err := utils.Parse(r, &body)
	var response []byte
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	// Call signup service
	id, base64, secret, err := h.service.SignUp(body.Email, body.Password)
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	// Create a response
	response, _ = json.Marshal(map[string]any{"success": true, "id": id, "image": base64, "secret": secret})
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (h userHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body types.SignIn
	var response []byte
	err := utils.Parse(r, &body)
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	id, err := h.service.SignIn(body.Email, body.Password)
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}
	response, _ = json.Marshal(map[string]any{"success": true, "id": id})

	w.WriteHeader(http.StatusOK)
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
	bearer := r.Header.Get("Authorization")
	jwtToken := strings.Split(bearer, " ")
	token, err := jwt.ParseWithClaims(jwtToken[1], &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.C.JWT_SECRET), nil
	})
	var response []byte
	if !token.Valid {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}
	claims, ok := token.Claims.(*CustomClaims)
	if !ok && !token.Valid {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	user, err := h.service.GetUser(claims.Id)
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}
	response, _ = json.Marshal(map[string]any{"success": true, "email": user.Email})
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (h userHandler) ConfirmOtp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body types.ConfirmSignUp
	err := utils.Parse(r, &body)
	var response []byte
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}
	token, err := h.service.ConfirmOtp(body.Id, body.Otp)
	if err != nil {
		response, _ = json.Marshal(map[string]any{"success": false, "error": err.Error()})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(response)
		return
	}

	response, _ = json.Marshal(map[string]any{"success": true, "token": &token})
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}
