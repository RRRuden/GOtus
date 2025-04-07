package user

import (
	"fmt"
)

type User struct {
	id    int
	Name  string
	Email string
}

func NewUser(id int, name, email string) *User {
	return &User{
		id:    id,
		Name:  name,
		Email: email,
	}
}

func (u *User) GetID() int {
	return u.id
}

func (u User) String() string {
	return fmt.Sprintf("[User] ID: %d, Name: %s, Email: %s", u.GetID(), u.Name, u.Email)
}
