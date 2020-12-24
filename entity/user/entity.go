package user

import "go.mongodb.org/mongo-driver/bson/primitive"

// User struct
type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"-"`
	Fullname string             `bson:"fullname" json:"fullname"`
	Email    string             `bson:"email" json:"email"`
	Password string             `bson:"password" json:"-"`
}
