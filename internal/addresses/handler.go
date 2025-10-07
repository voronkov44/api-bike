package addresses

import (
	"bike/configs"
	"bike/internal/users"
	"bike/pkg/middleware"
	"bike/pkg/req"
	"bike/pkg/res"
	"errors"
	"net/http"
	"strconv"
)

type AddressHandlerDeps struct {
	AddressRepository *AddressRepository
	AddressService    *AddressService
	UserRepository    *users.UserRepository
	Config            *configs.Config
}

type AddressHandler struct {
	AddressRepository *AddressRepository
	service           *AddressService
}

func NewAddressHandler(router *http.ServeMux, deps AddressHandlerDeps) {
	handler := &AddressHandler{
		AddressRepository: deps.AddressRepository,
		service:           deps.AddressService,
	}

	// Защищённые маршруты — пользователь должен быть авторизован
	router.Handle("POST /user/address", middleware.IsAuthenticated(handler.Create(), deps.Config))
	router.Handle("GET /user/address", middleware.IsAuthenticated(handler.GetAllForUser(), deps.Config))
	router.Handle("PATCH /user/address/{id}", middleware.IsAuthenticated(handler.Patch(), deps.Config))
	router.Handle("DELETE /user/address/{id}", middleware.IsAuthenticated(handler.Delete(), deps.Config))

	router.HandleFunc("GET /user/adminaddress", handler.AdminListAll())
}

func (handler *AddressHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[AddressCreateRequest](&w, r)
		if err != nil {
			return
		}

		email, _ := r.Context().Value(middleware.ContextEmailKey).(string)
		created, err := handler.service.CreateAddress(r.Context(), email, *body)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to create Address"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, ToResponse(created), http.StatusCreated)
	}

}

func (handler *AddressHandler) GetAllForUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email, _ := r.Context().Value(middleware.ContextEmailKey).(string)
		list, err := handler.service.ListAddress(r.Context(), email)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to list Address"}, http.StatusInternalServerError)
			return
		}
		// конвертируем в response-payload
		out := make([]AddressResponse, 0, len(list))
		for _, a := range list {
			out = append(out, ToResponse(&a))
		}
		res.Json(w, out, http.StatusOK)
	}
}

func (handler *AddressHandler) Patch() http.HandlerFunc {
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

		body, err := req.HandleBody[AddressUpdateRequest](&w, r)
		if err != nil {
			return
		}

		email, _ := r.Context().Value(middleware.ContextEmailKey).(string)
		updated, err := handler.service.UpdateAddress(r.Context(), email, id, *body)
		if err != nil {
			if errors.Is(err, ErrAddressNotFound) {
				res.Json(w, map[string]string{"error": "Address not found"}, http.StatusNotFound)
				return
			}
			if errors.Is(err, ErrForbidden) {
				res.Json(w, map[string]string{"error": "forbidden"}, http.StatusForbidden)
				return
			}
			res.Json(w, map[string]string{"error": "failed to update Address"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, ToResponse(updated), http.StatusOK)
	}
}

func (handler *AddressHandler) Delete() http.HandlerFunc {
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

		email, _ := r.Context().Value(middleware.ContextEmailKey).(string)
		err = handler.service.DeleteAddress(r.Context(), email, id)
		if err != nil {
			if errors.Is(err, ErrAddressNotFound) {
				res.Json(w, map[string]string{"error": "Address not found"}, http.StatusNotFound)
				return
			}
			if errors.Is(err, ErrForbidden) {
				res.Json(w, map[string]string{"error": "forbidden"}, http.StatusForbidden)
				return
			}
			res.Json(w, map[string]string{"error": "failed to delete Address"}, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}

}

func (handler *AddressHandler) AdminListAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()

		// параметры пагинации
		page := 1
		limit := 10
		if p := q.Get("page"); p != "" {
			if n, err := strconv.Atoi(p); err == nil && n > 0 {
				page = n
			}
		}
		if l := q.Get("limit"); l != "" {
			if n, err := strconv.Atoi(l); err == nil && n > 0 {
				limit = n
			}
		}

		// фильтры
		var userID uint = 0
		if u := q.Get("user_id"); u != "" {
			if id64, err := strconv.ParseUint(u, 10, 32); err == nil {
				userID = uint(id64)
			}
		}
		city := q.Get("city")
		street := q.Get("street")
		phone := q.Get("phone")
		label := q.Get("label")

		items, total, totalPages, err := handler.service.ListAllAdmin(r.Context(), userID, city, street, phone, label, page, limit)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to list addresses"}, http.StatusInternalServerError)
			return
		}

		out := make([]AddressResponse, 0, len(items))
		for _, a := range items {
			out = append(out, ToResponse(&a))
		}

		resp := AdminAddressesResponse{
			Addresses:  out,
			Total:      total,
			Page:       page,
			Limit:      limit,
			TotalPages: totalPages,
		}
		res.Json(w, resp, http.StatusOK)
	}
}
