package provider

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"login-meta-jatis/util"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBClient creates a new MongoDB client
func NewMongoDBClient() (*mongo.Client, error) {
	cfg := util.Configuration.MongoDB

	optionsStr := fmt.Sprintf("&%s=%s", "authSource", url.QueryEscape(cfg.AuthSource))

	// Set the MongoDB client options
	clientOptions := options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s@%s:%d/?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, optionsStr))

	// Customize additional options if needed
	clientOptions.MaxPoolSize = &cfg.MaxPoolSize
	timeout := time.Duration(cfg.ConnectTimeout) * time.Second
	clientOptions.ConnectTimeout = &timeout

	// Connect to the MongoDB server
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: 1}}).Decode(&result); err != nil {
		return nil, err
	}

	return client, nil
}
