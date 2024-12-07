package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
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
	_, err = g.InsertOne(ctx, bson.D{{"time", time.Now().In(time.FixedZone("UTC+5:30", 19800))}})
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Hello, World!",
		})
	})
	router.Run(":3000")

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
