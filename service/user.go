package service

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Email string             `json:"email"`
}

type UserService interface {
	SignUp(email string, password string) (*string, error)
	SignIn(email string, password string) (*UserService, error)
}
