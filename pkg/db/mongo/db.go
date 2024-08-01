package mongo

import "go.mongodb.org/mongo-driver/mongo"

type Db struct {
	*mongo.Database
	*mongo.Client
}

var dbInstance *Db

func Get() *Db {
	return dbInstance
}

func Set(db *mongo.Database, client *mongo.Client) {
	dbInstance = &Db{
		Database: db,
		Client:   client,
	}
}
