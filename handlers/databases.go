package handlers

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Ctx = context.Background()

var RedisClient *redis.Client

func ConnectToRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB,
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to Redis!")
}

func ConnectToMongo() {
	// MongoDB connection string
	uri := "mongodb://localhost:27017/"

	// Set up options
	clientOptions := options.Client().ApplyURI(uri)

	err := mgm.SetDefaultConfig(nil, "TaskManager", clientOptions)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to MongoDB!")
}
