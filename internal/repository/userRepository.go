package repository

import (
	"encoding/csv"
	"gotus/internal/model/user"
	"os"
	"strconv"
	"sync"
)

const userCSVPath = "./data/users.csv"

var (
	users      []*user.User
	usersMutex sync.Mutex
)

func StoreUser(u *user.User) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	users = append(users, u)
	saveUserToCSV(u)
}

func GetUsers() ([]*user.User, int) {
	usersMutex.Lock()
	defer usersMutex.Unlock()
	return users, len(users)
}

func LoadUsersFromCSV() {
	file, err := os.Open(userCSVPath)
	if err != nil {
		return
	}
	defer file.Close()

	r := csv.NewReader(file)
	records, _ := r.ReadAll()

	usersMutex.Lock()
	defer usersMutex.Unlock()

	for _, rec := range records {
		id, _ := strconv.Atoi(rec[0])
		u := user.NewUser(id, rec[1], rec[2])
		users = append(users, u)
	}
}

func saveUserToCSV(u *user.User) {
	file, _ := os.OpenFile(userCSVPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	_ = w.Write([]string{strconv.Itoa(u.GetID()), u.Name, u.Email})
}
