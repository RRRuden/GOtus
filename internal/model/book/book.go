package book

import (
	"errors"
	"fmt"
	"regexp"
)

type Book struct {
	Isbn   string `json:"Isbn" bson:"isbn"`
	Title  string `json:"Title" bson:"title"`
	Author string `json:"Author" bson:"author"`
	Year   int    `json:"Year" bson:"year"`
}

func NewBook(isbn, title string, author string, year int) (*Book, error) {
	b := &Book{}
	if err := b.setISBN(isbn); err != nil {
		return nil, err
	}
	b.Title = title
	b.Author = author
	b.Year = year

	return b, nil
}

func (b *Book) setISBN(isbn string) error {
	if !isValidISBN(isbn) {
		return errors.New("некорректный ISBN")
	}
	b.Isbn = isbn
	return nil
}

func (b *Book) GetISBN() string {
	return b.Isbn
}

func isValidISBN(isbn string) bool {
	match, _ := regexp.MatchString(`^(97[89]-\d{1,5}-\d{1,7}-\d{1,7}-\d)$`, isbn)
	return match
}

func (b Book) String() string {
	return fmt.Sprintf("[Book] ISBN: %s, Title: %s, Author: %s, Year: %d", b.GetISBN(), b.Title, b.Author, b.Year)
}
