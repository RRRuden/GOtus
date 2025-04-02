package user

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
