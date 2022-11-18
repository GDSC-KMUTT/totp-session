package service

type User struct {
	Id    int64  `json:"id" bson:"_id"`
	Email string `json:"email"`
}

type UserService interface {
	SignUp(email string, password string) (*int64, *string, *string, error)
	SignIn(email string, password string) (*int64, error)
	GetUser(id int64) (*User, error)
	ConfirmOtp(id int64, otp string) (*string, error)
}
