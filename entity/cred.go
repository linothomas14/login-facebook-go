package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cred struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID string             `json:"client_id" bson:"client_id"`
	UserID   string
	Webhook  string
}
