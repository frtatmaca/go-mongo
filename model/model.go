package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	Name       string             `json:"name" Usage:"required"`
	Email      string             `json:"email" Usage:"required,alphanumeric"`
	Age        int                `json:"age"`
	Password   string             `json:"password" Usage:"required,alphanumeric"`
	Occupation string             `json:"occupation"`
	token      string             `json:"token"`
	Newtoken   string             `json:"new_token"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
