package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"login-meta-jatis/http/api"
	"login-meta-jatis/provider"
	"login-meta-jatis/util"

	"go.mongodb.org/mongo-driver/mongo"
)

func init() {
	if err := util.LoadConfig("."); err != nil {
		log.Fatal(err)
	}

	provider.InitLogDir()
}

func main() {
	logger := provider.NewLogger()

	mongoClient, err := provider.NewMongoDBClient()
	if err != nil {
		log.Fatal(err)
	}

	logger.Infof(provider.AppLog, "Successfully connected to MongoDB.")

	go func(c *mongo.Client, logger provider.ILogger) {
		var credRepo repository.CredRepository = repository.NewCredRepositoryImpl(c, logger)

		app := api.NewApp(outboundService, logger)
		addr := fmt.Sprintf(":%v", util.Configuration.Server.Port)
		server, err := app.CreateServer(addr)
		if err != nil {
			log.Fatal(err)
		}

		logger.Infof(provider.AppLog, "Server running at: %s", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Errorf(provider.AppLog, "Server error: %v", err)
		}

	}(mongoClient, logger)

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, os.Interrupt, syscall.SIGTERM)

	sig := <-shutdownCh
	logger.Infof(provider.AppLog, "Receiving signal: %s", sig)

	func(c *mongo.Client) {
		if err := c.Disconnect(context.Background()); err != nil {
			log.Fatal(err)
		}

		logger.Infof(provider.AppLog, "Successfully disconnected from MongoDB.")

	}(mongoClient)
}
