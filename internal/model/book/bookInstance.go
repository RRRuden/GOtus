package book

type BookInstance struct {
	id   int
	ISBN string
}

func NewBookInstance(id int, isbn string) *BookInstance {
	return &BookInstance{
		id:   id,
		ISBN: isbn,
	}
}

func (bi *BookInstance) GetID() int {
	return bi.id
}
