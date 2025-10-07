package products

import (
	"bike/configs"
	"bike/pkg/req"
	"bike/pkg/res"
	"errors"
	"net/http"
	"strconv"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
	ProductService    ProductService
	Config            *configs.Config
}

type ProductHandler struct {
	ProductRepository *ProductRepository
	service           ProductService
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		ProductRepository: deps.ProductRepository,
		service:           deps.ProductService,
	}
	router.HandleFunc("POST /products", handler.Create())
	router.HandleFunc("GET /products", handler.GetAll())
	router.HandleFunc("GET /products/{slug}", handler.GoTo())
	router.Handle("PATCH /products/{slug}", handler.Update())
	router.HandleFunc("DELETE /products/{slug}", handler.Delete())

	router.HandleFunc("POST /products/{slug}/change", handler.Change())
}

// Create godoc
// @Summary Создать продукт (админ)
// @Description Создаёт новый продукт
// @Tags products,admin
// @Accept json
// @Produce json
// @Param request body products.ProductCreateRequest true "Product data"
// @Success 201 {object} products.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductCreateRequest](&w, r)
		if err != nil {
			return
		}

		created, err := handler.service.Create(r.Context(), *body)
		if err != nil {
			// Маппим доменные ошибки в HTTP
			if errors.Is(err, ErrValidation) {
				res.Json(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
				return
			}
			// По умолчанию считаем это 500
			res.Json(w, map[string]string{"error": "failed to create product"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, created, http.StatusCreated)
	}
}

// GetAll godoc
// @Summary Список продуктов
// @Description Возвращает список продуктов (пагинация через limit/offset)
// @Tags products,open
// @Produce json
// @Param limit query int false "limit"
// @Param offset query int false "offset"
// @Success 200 {array} products.Product
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (handler *ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		limit, offset := 0, 0
		if v := q.Get("limit"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				limit = n
			} else {
				res.Json(w, map[string]string{"error": "invalid limit"}, http.StatusBadRequest)
				return
			}
		}
		if v := q.Get("offset"); v != "" {
			if n, err := strconv.Atoi(v); err == nil && n >= 0 {
				offset = n
			} else {
				res.Json(w, map[string]string{"error": "invalid offset"}, http.StatusBadRequest)
				return
			}
		}

		list, err := handler.service.GetAll(r.Context(), limit, offset)
		if err != nil {
			res.Json(w, map[string]string{"error": "failed to list products"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, list, http.StatusOK)
	}
}

// GoTo godoc
// @Summary Получить блюдо по slug, переход на конкретное блюдо
// @Tags products,open
// @Produce json
// @Param slug path string true "slug"
// @Success 200 {object} products.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{slug} [get]
func (handler *ProductHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl := r.PathValue("slug")
		if sl == "" {
			res.Json(w, map[string]string{"error": "invalid slug"}, http.StatusBadRequest)
			return
		}

		p, err := handler.service.GoTo(r.Context(), sl)
		if err != nil {
			if errors.Is(err, ErrNotFound) {
				res.Json(w, map[string]string{"error": "product not found"}, http.StatusNotFound)
				return
			}
			res.Json(w, map[string]string{"error": "failed to get product"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, p, http.StatusOK)
	}
}

// Update godoc
// @Summary Обновить продукт (админ)
// @Tags products,admin
// @Accept json
// @Produce json
// @Param slug path string true "slug"
// @Param request body products.ProductUpdateRequest true "Fields to update"
// @Success 200 {object} products.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{slug} [patch]
func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl := r.PathValue("slug")
		if sl == "" {
			res.Json(w, map[string]string{"error": "invalid slug"}, http.StatusBadRequest)
			return
		}

		//Вывод почты в консоль
		//email, ok := r.Context().Value(middleware.ContextEmailKey).(string)
		//if ok {
		//	fmt.Println(email)
		//}
		body, err := req.HandleBody[ProductUpdateRequest](&w, r)
		if err != nil {
			return
		}

		updated, err := handler.service.Update(r.Context(), sl, *body)
		switch {
		case errors.Is(err, ErrValidation):
			res.Json(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		case errors.Is(err, ErrNotFound):
			res.Json(w, map[string]string{"error": "product not found"}, http.StatusNotFound)
			return
		case err != nil:
			res.Json(w, map[string]string{"error": "failed to update product"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, updated, http.StatusOK)
	}
}

// Delete godoc
// @Summary Удалить продукт (админ)
// @Tags products,admin
// @Param slug path string true "slug"
// @Success 204
// @Failure 404 {object} map[string]string
// @Router /products/{slug} [delete]
func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sl := r.PathValue("slug")
		if sl == "" {
			res.Json(w, map[string]string{"error": "invalid slug"}, http.StatusBadRequest)
			return
		}

		if err := handler.service.Delete(r.Context(), sl); err != nil {
			if errors.Is(err, ErrNotFound) {
				res.Json(w, map[string]string{"error": "product not found"}, http.StatusNotFound)
				return
			}
			res.Json(w, map[string]string{"error": "failed to delete product"}, http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}

// Change godoc
// @Summary Сменить slug продукта (админ)
// @Tags products,admin
// @Accept json
// @Produce json
// @Param slug path string true "current slug"
// @Param request body products.ProductSlugUpdateRequest true "new slug"
// @Success 200 {object} products.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{slug}/change [post]
func (handler *ProductHandler) Change() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cur := r.PathValue("slug")
		if cur == "" {
			res.Json(w, map[string]string{"error": "invalid slug"}, http.StatusBadRequest)
			return
		}

		body, err := req.HandleBody[ProductSlugUpdateRequest](&w, r)
		if err != nil {
			return
		}

		updated, err := handler.service.ChangeSlug(r.Context(), cur, body.Slug)
		switch {
		case errors.Is(err, ErrValidation):
			res.Json(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
			return
		case errors.Is(err, ErrNotFound):
			res.Json(w, map[string]string{"error": "product not found"}, http.StatusNotFound)
			return
		case err != nil:
			res.Json(w, map[string]string{"error": "failed to change slug"}, http.StatusInternalServerError)
			return
		}
		res.Json(w, updated, http.StatusOK)
	}
}
