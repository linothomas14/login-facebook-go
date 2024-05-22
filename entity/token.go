package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID  string             `json:"client_id" bson:"client_id"`
	Session   string             `json:"session" bson:"session"`
	Token     string             `json:"token" bson:"token"`
	TokenMeta string             `json:"token_meta" bson:"token_meta"`
	IsActive  bool               `json:"is_active" bson:"is_active"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	ExpiredAt time.Time          `json:"expired_at" bson:"expired_at"`
}
