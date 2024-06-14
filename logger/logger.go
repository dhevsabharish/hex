package logger

import (
    "context"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var logCollection *mongo.Collection

type LogEntry struct {
    Timestamp time.Time `bson:"timestamp"`
    Level     string    `bson:"level"`
    Message   string    `bson:"message"`
    UserID    uint      `bson:"user_id,omitempty"`
    BookID    uint      `bson:"book_id,omitempty"`
}

func InitLogger(mongoURI, dbName, collectionName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
    if err != nil {
        return err
    }

    // Ping the MongoDB server to verify the connection
    if err = client.Ping(ctx, nil); err != nil {
        return err
    }

    logCollection = client.Database(dbName).Collection(collectionName)
    log.Println("Connected to MongoDB successfully")
    return nil
}

func CloseLogger() {
    if client != nil {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        if err := client.Disconnect(ctx); err != nil {
            log.Printf("Error closing MongoDB connection: %v", err)
        }
    }
}

func LogInfo(message string, userID, bookID uint) {
    entry := LogEntry{
        Timestamp: time.Now(),
        Level:     "INFO",
        Message:   message,
        UserID:    userID,
        BookID:    bookID,
    }
    _, err := logCollection.InsertOne(context.Background(), entry)
    if err != nil {
        log.Printf("Error logging info: %v", err)
    }
}

func LogError(message string, userID, bookID uint) {
    entry := LogEntry{
        Timestamp: time.Now(),
        Level:     "ERROR",
        Message:   message,
        UserID:    userID,
        BookID:    bookID,
    }
    _, err := logCollection.InsertOne(context.Background(), entry)
    if err != nil {
        log.Printf("Error logging error: %v", err)
    }
}

func GetAllLogs() (*mongo.Cursor, error) {
    return logCollection.Find(context.Background(), bson.D{})
}