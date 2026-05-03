package config

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func CreateUserIndexes() {
	collection := DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	indexModels := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("unique_email"),
		},
		{
			Keys: bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("unique_username"),
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexModels)
	if err != nil {
		log.Fatal("failed to create indexes:", err)
	}
}