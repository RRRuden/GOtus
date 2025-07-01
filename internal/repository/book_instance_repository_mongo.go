package repository

import (
	"context"
	"gotus/internal/model/book"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoBookInstanceRepository struct {
	collection *mongo.Collection
}

func NewMongoBookInstanceRepository(db *mongo.Database) BookInstanceRepository {
	return &mongoBookInstanceRepository{collection: db.Collection("book_instances")}
}

func (r *mongoBookInstanceRepository) StoreBookInstance(bi *book.BookInstance) error {
	_, err := r.collection.InsertOne(context.TODO(), bi)
	return err
}

func (r *mongoBookInstanceRepository) GetBookInstances() ([]*book.BookInstance, int, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var result []*book.BookInstance
	for cursor.Next(context.TODO()) {
		var bi book.BookInstance
		if err := cursor.Decode(&bi); err != nil {
			return nil, 0, err
		}
		result = append(result, &bi)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return result, len(result), nil
}

func (r *mongoBookInstanceRepository) UpdateBookInstanceById(id int, updated *book.BookInstance) (bool, error) {
	res, err := r.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": updated})
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *mongoBookInstanceRepository) FindBookInstanceById(id int) (*book.BookInstance, bool, error) {
	var bi book.BookInstance
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&bi)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &bi, true, nil
}

func (r *mongoBookInstanceRepository) GetBookInstancesByISBN(isbn string) ([]*book.BookInstance, int, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{"isbn": isbn})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var result []*book.BookInstance
	for cursor.Next(context.TODO()) {
		var bi book.BookInstance
		if err := cursor.Decode(&bi); err != nil {
			return nil, 0, err
		}
		result = append(result, &bi)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}

	return result, len(result), nil
}

func (r *mongoBookInstanceRepository) DeleteBookInstanceById(id int) (bool, error) {
	res, err := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}
