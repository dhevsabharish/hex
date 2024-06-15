package logging

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDBLogger struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBLogger(uri, dbName, collectionName string) *MongoDBLogger {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoDBLogger{
		client:     client,
		collection: collection,
	}
}

func (l *MongoDBLogger) Log(level string, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logEntry := map[string]interface{}{
		"level":     level,
		"message":   message,
		"timestamp": time.Now(),
	}

	_, err := l.collection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Failed to log to MongoDB: %v", err)
	}
}
