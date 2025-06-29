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

func (r *mongoBookInstanceRepository) StoreBookInstance(bi *book.BookInstance) {
	r.collection.InsertOne(context.TODO(), bi)
}

func (r *mongoBookInstanceRepository) GetBookInstances() ([]*book.BookInstance, int) {
	cursor, _ := r.collection.Find(context.TODO(), bson.M{})
	var result []*book.BookInstance
	for cursor.Next(context.TODO()) {
		var bi book.BookInstance
		cursor.Decode(&bi)
		result = append(result, &bi)
	}
	return result, len(result)
}

func (r *mongoBookInstanceRepository) UpdateBookInstanceById(id int, updated *book.BookInstance) bool {
	res, _ := r.collection.UpdateOne(context.TODO(), bson.M{"id": id}, bson.M{"$set": updated})
	return res.ModifiedCount > 0
}

func (r *mongoBookInstanceRepository) FindBookInstanceById(id int) (*book.BookInstance, bool) {
	var bi book.BookInstance
	err := r.collection.FindOne(context.TODO(), bson.M{"id": id}).Decode(&bi)
	return &bi, err == nil
}

func (r *mongoBookInstanceRepository) GetBookInstancesByISBN(isbn string) ([]*book.BookInstance, int) {
	cursor, _ := r.collection.Find(context.TODO(), bson.M{"isbn": isbn})
	var result []*book.BookInstance
	for cursor.Next(context.TODO()) {
		var bi book.BookInstance
		cursor.Decode(&bi)
		result = append(result, &bi)
	}
	return result, len(result)
}

func (r *mongoBookInstanceRepository) DeleteBookInstanceById(id int) bool {
	res, _ := r.collection.DeleteOne(context.TODO(), bson.M{"id": id})
	return res.DeletedCount > 0
}
