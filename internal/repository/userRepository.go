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

func (r *UserRepository) saveUserToCSV(u *user.User) {
	file, _ := os.OpenFile(filepath.Join(r.dataDir, r.filename), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{strconv.Itoa(u.GetID()), u.Name, u.Email})
}
