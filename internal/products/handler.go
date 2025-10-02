package products

import (
	"bike/configs"
	"fmt"
	"net/http"
)

type ProductHandlerDeps struct {
	*configs.Config
}

type ProductHandler struct {
	*configs.Config
}

func NewProductHandler(router *http.ServeMux, deps ProductHandlerDeps) {
	handler := &ProductHandler{
		Config: deps.Config,
	}
	router.HandleFunc("POST /products", handler.Create())
	router.HandleFunc("GET /products", handler.GetAll())
	router.HandleFunc("GET /products/{id}", handler.GoTo())
	router.HandleFunc("PATCH /products/{id}", handler.Update())
	router.HandleFunc("DELETE /products/{id}", handler.Delete())
}

func (handler *ProductHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Create")
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
		fmt.Println("Delete")

	}
}
