package repository

import (
	"context"
	"gotus/internal/model/user"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) UserRepository {
	return &mongoUserRepository{collection: db.Collection("users")}
}

func (r *mongoUserRepository) StoreUser(u *user.User) {
	r.collection.InsertOne(context.TODO(), u)
}

func (r *mongoUserRepository) GetUsers() ([]*user.User, int) {
	cursor, _ := r.collection.Find(context.TODO(), bson.M{})
	var result []*user.User
	for cursor.Next(context.TODO()) {
		var u user.User
		cursor.Decode(&u)
		result = append(result, &u)
	}
	return result, len(result)
}

func (r *mongoUserRepository) UpdateUserById(id int, updated *user.User) bool {
	res, _ := r.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": updated})
	return res.ModifiedCount > 0
}

func (r *mongoUserRepository) FindUserById(id int) (*user.User, bool) {
	var u user.User
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&u)
	return &u, err == nil
}

func (r *mongoUserRepository) FindUserByEmail(email string) (*user.User, bool) {
	var u user.User
	err := r.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u)
	return &u, err == nil
}

func (r *mongoUserRepository) DeleteUserById(id int) bool {
	res, _ := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	return res.DeletedCount > 0
}
