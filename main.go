package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	cli, err := InitializeMongoDB(os.Getenv("MONGO_URI"))
	if err != nil {
		log.Printf("%s", err)
		return
	}
	g := cli.Database("test").Collection("insidetest")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = g.InsertOne(ctx, bson.D{{"time", time.Now()}})
	if err != nil {
		log.Fatal(err)
	}

}

// InitializeMongoDB initializes the MongoDB client
func InitializeMongoDB(uri string) (*mongo.Client, error) {
	ctx := context.Background()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	clientOptions := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return client, nil
}
