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

type CreateUserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// CreateUser godoc
// @Summary      Создать пользователя
// @Description  Добавляет нового пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "Пользователь"
// @Success      201
// @Failure      400 {string} string "invalid request"
// @Router       /api/user [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	u := user.NewUser(req.ID, req.Name, req.Email)
	h.Repo.StoreUser(u)
	w.WriteHeader(http.StatusCreated)
}

// GetAllUsers godoc
// @Summary      Получить всех пользователей
// @Description  Возвращает список всех пользователей
// @Tags         user
// @Produce      json
// @Success      200 {array} user.User
// @Router       /api/users [get]
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

// GetUserByID godoc
// @Summary      Получить пользователя по ID
// @Description  Возвращает пользователя по ID
// @Tags         user
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Success      200 {object} user.User
// @Failure      404 {string} string "not found"
// @Router       /api/user/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id int) {
	res, ok := h.Repo.FindUserById(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	json.NewEncoder(w).Encode(res)
}

// UpdateUser godoc
// @Summary      Обновить пользователя
// @Description  Обновляет пользователя по ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Param        user body UpdateUserRequest true "Новые данные пользователя"
// @Success      200
// @Failure      400 {string} string "invalid request"
// @Failure      404 {string} string "not found"
// @Router       /api/user/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateUserRequest

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

// DeleteUser godoc
// @Summary      Удалить пользователя
// @Description  Удаляет пользователя по ID
// @Tags         user
// @Param        id path int true "ID пользователя"
// @Success      200
// @Failure      404 {string} string "not found"
// @Router       /api/user/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	if !h.Repo.DeleteUserById(id) {
		http.NotFound(w, r)
		return
	}
	w.WriteHeader(http.StatusOK)
}
