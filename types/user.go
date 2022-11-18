package types

type SignUp struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmSignUp struct {
	Id  int64  `json:"id"`
	Otp string `json:"otp"`
}

type SignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ConfirmSignIn struct {
	Id  int64  `json:"id"`
	Otp string `json:"otp"`
}
