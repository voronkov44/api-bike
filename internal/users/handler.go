package users

import (
	"bike/configs"
	"bike/pkg/jwt"
	"bike/pkg/res"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

type UserHandlerDeps struct {
	UserRepository *UserRepository
	Config         *configs.Config
}

type UserHandler struct {
	repo   *UserRepository
	config *configs.Config
}

func NewUsersHandler(router *http.ServeMux, deps UserHandlerDeps) {
	handler := &UserHandler{
		repo:   deps.UserRepository,
		config: deps.Config,
	}

	router.HandleFunc("GET /users", handler.GetAll())
	router.HandleFunc("GET /users/{id}", handler.GetByID())
	router.HandleFunc("GET /users/jwt/{id}", handler.GetJWTForUser())
	router.HandleFunc("GET /users/search", handler.SearchUsers())

}

func toUserResponse(u *User) UserResponse {
	created := ""
	if !u.CreatedAt.IsZero() {
		created = u.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
	}
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		CreatedAt: created,
	}
}

// GetAll - возвращает пользователей с пагинацией, сортировкой и фильтрацией
func (handler *UserHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Параметры пагинации
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		if page < 1 {
			page = 1
		}

		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		if limit < 1 {
			limit = 10 // значение по умолчанию
		}
		if limit > 100 {
			limit = 100 // максимальный лимит
		}

		offset := (page - 1) * limit

		// Параметры фильтрации
		nameFilter := r.URL.Query().Get("name")
		emailFilter := r.URL.Query().Get("email")

		// Параметр сортировки
		sortBy := r.URL.Query().Get("sort")

		// Получаем данные
		list, err := handler.repo.ListAll(limit, offset, sortBy, nameFilter, emailFilter)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to list users"}, http.StatusInternalServerError)
			return
		}

		// Получаем общее количество для пагинации (с учетом фильтров)
		total, err := handler.repo.Count(nameFilter, emailFilter)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to count users"}, http.StatusInternalServerError)
			return
		}

		// Формируем ответ
		out := make([]UserResponse, 0, len(list))
		for _, u := range list {
			out = append(out, toUserResponse(&u))
		}

		totalPages := 0
		if limit > 0 {
			totalPages = (int(total) + limit - 1) / limit
		}

		response := UserListResponse{
			Users:      out,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		}

		res.Json(w, response, http.StatusOK)
	}
}

// GetByID - возвращает пользователя по id
func (handler *UserHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			res.Json(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
			return
		}
		id64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			res.Json(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
			return
		}
		id := uint(id64)

		user, err := handler.repo.FindByID(id)
		if err != nil {
			res.Json(w, map[string]string{"error": "user not found"}, http.StatusNotFound)
			return
		}
		res.Json(w, toUserResponse(user), http.StatusOK)
	}

}

// GetJWTForUser - возвращает jwt-token пользователя
func (handler *UserHandler) GetJWTForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		if idStr == "" {
			res.Json(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
			return
		}
		id64, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			res.Json(w, map[string]string{"error": "invalid id"}, http.StatusBadRequest)
			return
		}
		id := uint(id64)

		user, err := handler.repo.FindByID(id)
		if err != nil {
			res.Json(w, map[string]string{"error": "user not found"}, http.StatusNotFound)
			return
		}

		token, err := jwt.NewJWT(handler.config.Auth.Secret).GenerateToken(jwt.JWTData{
			Email: user.Email,
		})
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to generate token"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, map[string]string{"token": token}, http.StatusOK)
	}
}

// SearchUsers - поиск пользователей по email
func (handler *UserHandler) SearchUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			res.Json(w, map[string]string{"error": "email parameter is required"}, http.StatusBadRequest)
			return
		}

		users, err := handler.repo.SearchByEmail(email)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to search users"}, http.StatusInternalServerError)
			return
		}

		out := make([]UserResponse, 0, len(users))
		for _, u := range users {
			out = append(out, toUserResponse(&u))
		}
		res.Json(w, out, http.StatusOK)
	}
}
