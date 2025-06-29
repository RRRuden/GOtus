package book

import (
	"fmt"
)

type BookInstance struct {
	Id   int    `json:"Id" bson:"id"`
	ISBN string `json:"ISBN" bson:"isbn"`
}

func NewBookInstance(id int, isbn string) *BookInstance {
	return &BookInstance{
		Id:   id,
		ISBN: isbn,
	}
}

func (bi *BookInstance) GetID() int {
	return bi.Id
}

func (bi BookInstance) String() string {
	return fmt.Sprintf("[BookInstance] ID: %d, ISBN: %s", bi.GetID(), bi.ISBN)
}
