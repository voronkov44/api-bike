package products

import (
	"bike/pkg/req"
	"bike/pkg/res"
	"errors"
	"fmt"
	"net/http"
)

type ProductHandlerDeps struct {
	ProductRepository *ProductRepository
	ProductService    ProductService
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
	router.HandleFunc("GET /products/{id}", handler.GoTo())
	router.HandleFunc("PATCH /products/{id}", handler.Update())
	router.HandleFunc("DELETE /products/{id}", handler.Delete())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[ProductCreateRequest](&w, r)
		if err != nil {
			return
		}

		created, err := handler.service.CreateProduct(r.Context(), *body)
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

func (handler *ProductHandler) GetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GETAll")

	}
}

func (handler *ProductHandler) GoTo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GoTo")

	}
}

func (handler *ProductHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Update")

	}
}

func (handler *ProductHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Println(id)

	}
}
