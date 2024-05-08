package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID  string             `json:"client_id" bson:"client_id"`
	ExpiredAt time.Time          `json:"expired_at" bson:"expired_at"`
	Token     string             `json:"token" bson:"token"`
	TokenMeta string             `json:"token_meta" bson:"token_meta"`
}

func (t *Token) IsExpired() bool {
	now := time.Now()
	return now.After(t.ExpiredAt.Add(-7 * time.Hour))
}
