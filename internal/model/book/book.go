package book

type Book struct {
	id   int
	Name string
}

func NewBook(id int, name string) *Book {
	return &Book{
		id:   id,
		Name: name,
	}
}

func (b *Book) GetID() int {
	return b.id
}
