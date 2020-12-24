package mongodb

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewMongoDBConnection create new MongoDB Connection
func NewMongoDBConnection(user, pass, db, collection string) (*mongo.Client, *mongo.Collection, error) {
	mongoURI := fmt.Sprintf(
		"mongodb+srv://%s:%s@cluster0.ub9ns.mongodb.net/%s?retryWrites=true&w=majority",
		user,
		pass,
		db,
	)

	client, _ := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err := client.Connect(context.TODO()); err != nil {
		log.Fatal(err)
		return nil, nil, err
	}
	mongoCollection := client.Database(db).Collection(collection)
	return client, mongoCollection, nil
}
