package book

import (
	"fmt"
)

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

func (bi BookInstance) String() string {
	return fmt.Sprintf("[BookInstance] ID: %d, ISBN: %s", bi.GetID(), bi.ISBN)
}
