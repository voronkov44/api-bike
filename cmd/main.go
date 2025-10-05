package main

import (
	"bike/configs"
	"bike/internal/auth"
	"bike/internal/products"
	"bike/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	productRepository := products.NewProductRepository(database)

	// Service
	productService := products.NewProductService(productRepository)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	products.NewProductHandler(router, products.ProductHandlerDeps{
		ProductRepository: productRepository,
		ProductService:    productService,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: router,
	}

	fmt.Println("Server listening on :8081")
	server.ListenAndServe()
}
