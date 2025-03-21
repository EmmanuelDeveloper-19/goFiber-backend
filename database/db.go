package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client

func InitDB() {
	uri := "mongodb://root:61945@localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	client, _ = mongo.Connect(context.TODO(), clientOptions)
	client.Ping(context.TODO(), readpref.Primary())
}

func GetCollection(collectionName string) *mongo.Collection {
	return client.Database("demariot-db").Collection(collectionName)
}
