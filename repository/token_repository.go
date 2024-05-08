package repository

import (
	"context"
	"login-meta-jatis/entity"
	"login-meta-jatis/provider"
	"login-meta-jatis/util"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenRepository interface {
	Create(ctx context.Context, token entity.Token) (err error)
}

type TokenRepositoryImpl struct {
	coll *mongo.Collection
	log  provider.ILogger
}

func NewTokenRepositoryImpl(client *mongo.Client, log provider.ILogger) *TokenRepositoryImpl {
	db := client.Database(util.Configuration.MongoDB.Database)
	coll := db.Collection(util.Configuration.MongoDB.Collection.Token)
	return &TokenRepositoryImpl{coll: coll, log: log}
}

func (t *TokenRepositoryImpl) Create(ctx context.Context, token entity.Token) (err error) {
	logger := t.log.WithFields(provider.MongoLog, logrus.Fields{"DATABASE_NAME": util.Configuration.MongoDB.Database, "COLLECTION_NAME": util.Configuration.MongoDB.Collection.Token})
	logger.Infof("inserting into MongoDB database")

	token.ID = primitive.NewObjectID()

	result, err := t.coll.InsertOne(ctx, token)
	if err != nil {
		logger.Errorf("creating chat history in MongoDB failed: %s", err)
		return
	}

	_, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		logger.Errorf("error asserting InsertedID to ObjectID")
		return
	}
	return
}
