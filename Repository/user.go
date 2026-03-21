package repository

import (
	"context"
	domain "taskmanagement/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoUserRepository struct {
	Collection *mongo.Collection
}

func (r *MongoUserRepository) GetUserByEmail(email string) (domain.User, error) {

	var user domain.User

	err := r.Collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil

}

func (r *MongoUserRepository) CreateUser(user domain.User) error {

	_, err := r.Collection.InsertOne(context.Background(), user)
	if err != nil {
		return err
	}
	return nil
}
