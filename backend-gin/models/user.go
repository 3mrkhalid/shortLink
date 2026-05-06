package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username string             `bson:"username" json:"username" binding:"required"`
	Email    string             `bson:"email" json:"email" binding:"required,email"`
	Password string             `bson:"password" json:"password" binding:"required,min=6"`

	ResetToken        string    `bson:"reset_token,omitempty" json:"-"`
	ResetTokenExpires time.Time `bson:"reset_token_expires,omitempty" json:"-"`
}