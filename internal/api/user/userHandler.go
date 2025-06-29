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
	Repo repository.UserRepository
}

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// CreateUser godoc
// @Summary      Создать пользователя
// @Description  Добавляет нового пользователя
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "Пользователь"
// @Success      201
// @Failure      400 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/user [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	u := user.NewUser(0, req.Name, req.Email)
	if err := h.Repo.StoreUser(u); err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при сохранении пользователя")
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// GetAllUsers godoc
// @Summary      Получить всех пользователей
// @Description  Возвращает список всех пользователей
// @Tags         user
// @Produce      json
// @Success      200 {array} user.User
// @Failure      500 {object} ErrorResponse
// @Router       /api/users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	res, _, err := h.Repo.GetUsers()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении пользователей")
		return
	}
	writeJSON(w, res)
}

func (h *UserHandler) UserByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/user/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Некорректный ID")
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
		writeError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
	}
}

// GetUserByID godoc
// @Summary      Получить пользователя по ID
// @Description  Возвращает пользователя по ID
// @Tags         user
// @Produce      json
// @Param        id path int true "ID пользователя"
// @Success      200 {object} user.User
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/user/{id} [get]
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request, id int) {
	res, _, err := h.Repo.FindUserById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при получении пользователя")
		return
	}
	if res == nil {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	writeJSON(w, res)
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
// @Failure      400 {object} ErrorResponse
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/user/{id} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request, id int) {
	var req UpdateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Неверный формат запроса")
		return
	}

	u := user.NewUser(req.ID, req.Name, req.Email)
	ok, err := h.Repo.UpdateUserById(id, u)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при обновлении пользователя")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
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
// @Failure      404 {object} ErrorResponse
// @Failure      500 {object} ErrorResponse
// @Router       /api/user/{id} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request, id int) {
	ok, err := h.Repo.DeleteUserById(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Ошибка при удалении пользователя")
		return
	}
	if !ok {
		writeError(w, http.StatusNotFound, "Пользователь не найден")
		return
	}
	w.WriteHeader(http.StatusOK)
}

// Вспомогательные функции

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Message: msg})
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
