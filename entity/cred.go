package entity

import "go.mongodb.org/mongo-driver/bson/primitive"

type Cred struct {
	ID       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	ClientID string             `json:"client_id" bson:"client_id"`
	Username string             `json:"user" bson:"user"`
	Password string             `json:"password" bson:"password"`
}
