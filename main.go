package main

import (
	"context"
	"github.com/gin-gonic/gin"
	redis2 "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	mongo2 "rideShare/pkg/db/mongo"
	"rideShare/pkg/redis"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	err := Initialize(context.TODO(), r)
	if err != nil {
		panic(err)
	}

	if err = r.Run(":8080"); err != nil {
		panic(err)
	}

}

func initializeDB(ctx context.Context) {
	mongoUri := "mongodb://admin:admin@127.0.0.1:27018/?retryWrites=true&w=majority"
	database := "ride_share"

	cmdMonitor := &event.CommandMonitor{
		Started: func(_ context.Context, evt *event.CommandStartedEvent) {
			log.Print(evt.Command)
		},
	}

	clientOptions := options.Client().ApplyURI(mongoUri).SetMonitor(cmdMonitor)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		panic(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	log.Println("Initialized Mongo DB client")
	db := client.Database(database)

	mongo2.Set(db, client)
}

func initializeRedis(ctx context.Context) {
	r := redis2.NewClient(&redis2.Options{
		Addr: "localhost:6379",
	})
	log.Println("Initialized Redis Client")
	redis.SetClient(r)
}
