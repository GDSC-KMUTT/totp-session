package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/GDSC-KMUTT/totp-session/service"
	"github.com/GDSC-KMUTT/totp-session/types"
	"github.com/GDSC-KMUTT/totp-session/utils"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(userSerivice service.UserService) userHandler {
	return userHandler{service: userSerivice}
}

func (h userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	var response []byte
	var body types.SignIn
	err := utils.Parse(r, &body)
	if err != nil {
		response_value := map[string]any{"success": false, "error": err.Error()}
		response, _ := json.Marshal(response_value)
		w.Write(response)
		return
	}

	token, base64, err := h.service.SignUp(body.Email, body.Password)
	if err != nil {
		response_value := map[string]any{"success": false, "error": err.Error()}
		response, _ := json.Marshal(response_value)
		w.Write(response)
		return
	}
	response, _ = json.Marshal(map[string]any{"success": true, "token": token, "image": base64})
	w.Write(response)
	return
}

func (h userHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Body)
}

func (h userHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers()
	if err != nil {
		panic(err.Error())
	}
	response, err := json.Marshal(map[string]any{"users": users})
	if err != nil {
		panic(err)
	}
	w.Write(response)
}
