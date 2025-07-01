package repository

import (
	"context"
	"gotus/internal/model/book"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoBookRepository struct {
	collection *mongo.Collection
}

func NewMongoBookRepository(db *mongo.Database) BookRepository {
	return &mongoBookRepository{collection: db.Collection("books")}
}

func (r *mongoBookRepository) StoreBook(b *book.Book) error {
	_, err := r.collection.InsertOne(context.TODO(), b)
	return err
}

func (r *mongoBookRepository) GetBooks() ([]*book.Book, int, error) {
	cursor, err := r.collection.Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(context.TODO())

	var books []*book.Book
	for cursor.Next(context.TODO()) {
		var b book.Book
		if err := cursor.Decode(&b); err != nil {
			return nil, 0, err
		}
		books = append(books, &b)
	}
	if err := cursor.Err(); err != nil {
		return nil, 0, err
	}
	return books, len(books), nil
}

func (r *mongoBookRepository) UpdateBookByISBN(isbn string, updatedBook *book.Book) (bool, error) {
	res, err := r.collection.UpdateOne(context.TODO(), bson.M{"isbn": isbn}, bson.M{"$set": updatedBook})
	if err != nil {
		return false, err
	}
	return res.ModifiedCount > 0, nil
}

func (r *mongoBookRepository) DeleteBookByISBN(isbn string) (bool, error) {
	res, err := r.collection.DeleteOne(context.TODO(), bson.M{"isbn": isbn})
	if err != nil {
		return false, err
	}
	return res.DeletedCount > 0, nil
}

func (r *mongoBookRepository) FindBookByISBN(isbn string) (*book.Book, bool, error) {
	var b book.Book
	err := r.collection.FindOne(context.TODO(), bson.M{"isbn": isbn}).Decode(&b)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &b, true, nil
}
