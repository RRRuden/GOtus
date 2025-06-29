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

func (r *mongoUserRepository) StoreUser(u *user.User) error {
	_, err := r.collection.InsertOne(context.TODO(), u)
	return err
}

func (r *mongoUserRepository) GetUsers() ([]*user.User, int, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var result []*user.User
	for cursor.Next(context.TODO()) {
		var u user.User
		if err := cursor.Decode(&u); err != nil {
			return nil, 0, err
		}
		result = append(result, &u)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return result, len(result), nil
}

func (r *mongoUserRepository) UpdateUserById(id int, updated *user.User) (bool, error) {
	res, err := r.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": updated})
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *mongoUserRepository) FindUserById(id int) (*user.User, bool, error) {
	var u user.User
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &u, true, nil
}

func (r *mongoUserRepository) FindUserByEmail(email string) (*user.User, bool, error) {
	var u user.User
	err := r.collection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &u, true, nil
}

func (r *mongoUserRepository) DeleteUserById(id int) (bool, error) {
	res, err := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
