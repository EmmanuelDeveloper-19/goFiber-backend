package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Street  string `bson:"street" json:"street"`
	City    string `bson:"city" json:"city"`
	State   string `bson:"state" json:"state"`
	Country string `bson:"country" json:"country"`
	ZIP     string `bson:"zip" json:"zip"`
}

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName      string             `bson:"first_name" json:"first_name"`
	LastName       string             `bson:"last_name" json:"last_name"`
	Email          string             `bson:"email" json:"email"`
	Password       string             `bson:"password,omitempty" json:"password"`
	Address        Address            `bson:"address" json:"address"`
	Phone          string             `bson:"phone" json:"phone"`
	Role           string             `bson:"role" json:"role"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
	LastLogin      time.Time          `bson:"last_login,omitempty" json:"last_login,omitempty"`
	ProfilePicture string             `bson:"profile_picture,omitempty" json:"profile_picture,omitempty"`
}
