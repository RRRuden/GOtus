package book

type Book struct {
	id     int
	Name   string
	Author string
}

func NewBook(id int, name string, author string) *Book {
	return &Book{
		id:     id,
		Name:   name,
		Author: author,
	}
}

func (b *Book) GetID() int {
	return b.id
}
