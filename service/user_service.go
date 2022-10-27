package service

import "github.com/GDSC-KMUTT/totp-session/repository"

type userService struct {
	repository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) userService {
	return userService{repository: userRepository}
}

func (s userService) SignUp(email string, password string) (*string, error) {
	return nil, nil
}

func (s userService) SignIn(email string, password string) (*UserService, error) {
	return nil, nil
}
