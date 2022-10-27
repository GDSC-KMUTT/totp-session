package handler

import (
	"fmt"
	"net/http"

	"github.com/GDSC-KMUTT/totp-session/service"
)

type userHandler struct {
	service service.UserService
}

func NewUserHandler(userSerivice service.UserService) userHandler {
	return userHandler{service: userSerivice}
}

func (h userHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Body)
}

func (h userHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, r.Body)
}
