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

func (r *mongoBookRepository) StoreBook(b *book.Book) {
	r.collection.InsertOne(context.TODO(), b)
}

func (r *mongoBookRepository) GetBooks() ([]*book.Book, int) {
	cursor, _ := r.collection.Find(context.TODO(), bson.M{})
	var books []*book.Book
	for cursor.Next(context.TODO()) {
		var b book.Book
		cursor.Decode(&b)
		books = append(books, &b)
	}
	return books, len(books)
}

func (r *mongoBookRepository) UpdateBookByISBN(isbn string, updatedBook *book.Book) bool {
	res, _ := r.collection.UpdateOne(context.TODO(), bson.M{"isbn": isbn}, bson.M{"$set": updatedBook})
	return res.ModifiedCount > 0
}

func (r *mongoBookRepository) DeleteBookByISBN(isbn string) bool {
	res, _ := r.collection.DeleteOne(context.TODO(), bson.M{"isbn": isbn})
	return res.DeletedCount > 0
}

func (r *mongoBookRepository) FindBookByISBN(isbn string) (*book.Book, bool) {
	var b book.Book
	err := r.collection.FindOne(context.TODO(), bson.M{"isbn": isbn}).Decode(&b)
	return &b, err == nil
}
