package service

import (
	"bytes"
	"encoding/base64"
	"errors"
	"image/png"

	"github.com/GDSC-KMUTT/totp-session/config"
	"github.com/GDSC-KMUTT/totp-session/repository"
	"github.com/golang-jwt/jwt"
	"github.com/pquerna/otp/totp"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) userService {
	return userService{repository: userRepository}
}

func (s userService) SignUp(email string, password string) (*int64, *string, *string, error) {
	// Generate a new secret TOTP key
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "GDSC KMUTT",
		AccountName: email,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	secret := key.Secret()

	// Hash the password
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, nil, err
	}

	// Create a new user
	user, err := s.repository.CreateUser(email, string(hashedPwd), secret)
	if err != nil {
		return nil, nil, nil, err
	}

	// Convert TOTP key into a PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, nil, nil, err
	}
	if err := png.Encode(&buf, img); err != nil {
		return nil, nil, nil, err
	}
	base64string := "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	url := key.URL()
	return &user.Id, &base64string, &url, nil
}

func (s userService) SignIn(email string, password string) (*int64, error) {
	user, err := s.repository.CheckUser(email)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &user.Id, nil
}

func (s userService) GetUser(id int64) (*User, error) {
	user, err := s.repository.GetUser(id)
	if err != nil {
		return nil, err
	}

	return &User{
		Id:    user.Id,
		Email: user.Email,
	}, nil
}

func (s userService) ConfirmOtp(id int64, otp string) (*string, error) {
	user, err := s.repository.GetUser(id)
	if err != nil {
		return nil, err
	}
	print(id)

	// Verify the OTP
	valid := totp.Validate(otp, user.Secret)
	if !valid {
		return nil, errors.New("invalid OTP")
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.Id,
	})
	token, err := claims.SignedString([]byte(config.C.JWT_SECRET))
	return &token, nil
}
