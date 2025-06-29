package user

import (
	"fmt"
)

type User struct {
	Id    int    `json:"Id" bson:"id"`
	Name  string `json:"Name" bson:"name"`
	Email string `json:"Email" bson:"email"`
}

func NewUser(id int, name, email string) *User {
	return &User{
		Id:    id,
		Name:  name,
		Email: email,
	}
}

func (u *User) GetID() int {
	return u.Id
}

func (u User) String() string {
	return fmt.Sprintf("[User] ID: %d, Name: %s, Email: %s", u.GetID(), u.Name, u.Email)
}
