package repository

import (
	"encoding/csv"
	"gotus/internal/model/user"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

type UserRepository struct {
	users      []*user.User
	dataDir    string
	filename   string
	usersMutex sync.Mutex
}

func NewUserRepository(dataDir string) *UserRepository {
	return &UserRepository{
		users:    []*user.User{},
		dataDir:  dataDir,
		filename: "users.csv",
	}
}

func (r *UserRepository) StoreUser(u *user.User) {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()
	r.users = append(r.users, u)
	r.saveUserToCSV(u)
}

func (r *UserRepository) GetUsers() ([]*user.User, int) {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()
	return r.users, len(r.users)
}

func (r *UserRepository) LoadUsersFromCSV() {
	file, err := os.Open(filepath.Join(r.dataDir, r.filename))
	if err != nil {
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, _ := reader.ReadAll()

	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		u := user.NewUser(id, rec[1], rec[2])
		r.users = append(r.users, u)
	}
}

func (r *UserRepository) UpdateUserById(id int, updatedReservation *user.User) bool {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	found := false
	for i, u := range r.users {
		if u.GetID() == id {
			r.users[i] = updatedReservation
			found = true
			break
		}
	}

	if found {
		r.saveAllToCSV()
	}

	return found
}

func (r *UserRepository) FindUserById(id int) (*user.User, bool) {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()
	for _, u := range r.users {
		if u.GetID() == id {
			return u, true
		}
	}
	return nil, false
}

func (r *UserRepository) FindUserByEmail(email string) (*user.User, bool) {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, true
		}
	}
	return nil, false
}

func (r *UserRepository) DeleteUserById(id int) bool {
	r.usersMutex.Lock()
	defer r.usersMutex.Unlock()

	found := false
	var updated []*user.User
	for _, u := range r.users {
		if u.GetID() != id {
			updated = append(updated, u)
		} else {
			found = true
		}
	}

	if found {
		r.users = updated
		r.saveAllToCSV()
	}

	return found
}

func (r *UserRepository) saveUserToCSV(u *user.User) {
	file, _ := os.OpenFile(filepath.Join(r.dataDir, r.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{strconv.Itoa(u.GetID()), u.Name, u.Email})
}

func (r *UserRepository) saveAllToCSV() {
	file, _ := os.Create(filepath.Join(r.dataDir, r.filename))
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, u := range r.users {
		_ = w.Write([]string{strconv.Itoa(u.GetID()), u.Name, u.Email})
	}
}
