package user

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoDBRepo struct
type MongoDBRepo struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

// NewMongoDBRepository function
func NewMongoDBRepository(collection *mongo.Collection) *MongoDBRepo {
	ctx := context.Background()
	return &MongoDBRepo{
		Collection: collection,
		Ctx:        ctx,
	}
}

// Register user function
func (r *MongoDBRepo) Register(user User) string {
	res, err := r.Collection.InsertOne(r.Ctx, user)
	if err != nil {
		log.Printf("Error when trying to register new user. Error: %v dari sini kan yak?", err.Error())
		return ""
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		log.Println("Error occur when trying to convert OID to string")
		return ""
	}
	return oid.Hex()
}

// FindOne user function
func (r *MongoDBRepo) FindOne(email string) (*User, error) {
	var user User
	if err := r.Collection.FindOne(r.Ctx, bson.M{"email": email}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
