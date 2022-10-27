package repository

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Secret   string             `json:"secret" bson:"secret"`
}

type UserRepository interface {
	CreateUser(email string, password string, secret string) (*User, *string, error)
	CheckUser(email string) (*User, error)
}
