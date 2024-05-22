package repository

import (
	"context"
	"fmt"

	"login-meta-jatis/entity"
	"login-meta-jatis/provider"
	"login-meta-jatis/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CredentialRepository interface {
	FindClientByClientID(ctx context.Context, userID string) (err error)
}

type CredRepositoryImpl struct {
	coll *mongo.Collection
	log  provider.ILogger
}

func NewCredRepositoryImpl(client *mongo.Client, log provider.ILogger) *CredRepositoryImpl {
	// fmt.Println("db=>>>>", util.Configuration.MongoDB.Database, "coll ->>>>", util.Configuration.MongoDB.Collection.ClientCredential)
	db := client.Database(util.Configuration.MongoDB.Database)
	coll := db.Collection(util.Configuration.MongoDB.Collection.ClientCredential)
	// fmt.Println("db=>>>>", db, "coll ->>>>", coll)
	return &CredRepositoryImpl{coll: coll, log: log}
}

func (c *CredRepositoryImpl) FindClientByClientID(
	ctx context.Context,
	clientID string,
) (err error) {

	// fmt.Println(c.coll)

	var cred entity.Cred
	var reqID string
	val := ctx.Value("req-id")
	if val != nil {
		reqID = val.(string)
	}

	logger := c.log.WithFields(
		provider.MongoLog,
		logrus.Fields{
			"REQUEST_ID":      reqID,
			"DATABASE_NAME":   util.Configuration.MongoDB.Database,
			"COLLECTION_NAME": util.Configuration.MongoDB.Collection.Token,
		},
	)

	logger.Info("Checking credential from mongo db")

	filter := bson.D{
		{
			Key:   "client_id",
			Value: clientID,
		},
	}

	if mongoErr := c.coll.FindOne(ctx, filter).Decode(&cred); mongoErr != nil {
		if mongoErr == mongo.ErrNoDocuments {
			err = fmt.Errorf("client not found")
			logger.Errorf("Checking credential from mongo db failed %s", err)
			return
		}

		err = fmt.Errorf("%v: %v", ErrMongo, mongoErr)
		logger.Errorf("Checking credential from mongo db failed %s", err)
		return
	}

	logger.Info("Successfully get credential from mongo db")
	return
}
