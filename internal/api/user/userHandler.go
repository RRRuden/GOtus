package handler

import (
	"encoding/json"
	"gotus/internal/model/user"
	"gotus/internal/repository"
	"net/http"
	"strconv"
	"strings"
)

type UserHandler struct {
	Repo *repository.UserRepository
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	u := user.NewUser(req.ID, req.Name, req.Email)
	h.Repo.StoreUser(u)
	w.WriteHeader(http.StatusCreated)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, _ := h.Repo.GetUsers()
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid reservation id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.GetUserByID(w, r, id)
	case http.MethodPut:
		h.UpdateUser(w, r, id)
	case http.MethodDelete:
		h.DeleteUser(w, r, id)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id int) {
	res, ok := h.Repo.FindUserById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	u := user.NewUser(req.ID, req.Name, req.Email)
	if !h.Repo.UpdateUserById(id, u) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if !h.Repo.DeleteUserById(id) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
