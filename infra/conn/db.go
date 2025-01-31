package conn

import (
	"context"
	"fmt"
	"vehicles/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mdb *mongo.Client

func ConnectDb() {
	conf := config.DB().Mongo

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	mdbClient, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.ConnectionString))

	if err != nil {
		panic(err)
	}

	//ping the database
	err = mdbClient.Ping(ctx, nil)
	if err != nil {
		panic(err)
	}

	mdb = mdbClient

	fmt.Println("Connected to MongoDB successful...")
}

func Db() *mongo.Client {
	return mdb
}

//func(db *mongo.Client) GetCollection(database, collection string) *mongo.Collection {
//	return db.Database(database).Collection(collection)
//}
