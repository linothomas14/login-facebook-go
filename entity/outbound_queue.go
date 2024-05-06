package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OutboundQueue struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	RequestID   string             `json:"request_id" bson:"request_id"`
	ClientID    string             `json:"client_id" bson:"client_id"`
	ReferenceID string             `json:"reference_id" bson:"reference_id"`
	Messages    map[string]any     `json:"messages" bson:"messages"`
	APIVersion  string             `bson:"api_version,omitempty"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	ReceivedAt  time.Time          `json:"received_at" bson:"received_at"`
}
