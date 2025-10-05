package main

import (
	"bike/configs"
	"bike/internal/auth"
	"bike/internal/products"
	"bike/internal/users"
	"bike/pkg/db"
	"bike/pkg/middleware"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	productRepository := products.NewProductRepository(database)
	userRepository := users.NewUserRepository(database)

	// Services
	productService := products.NewProductService(productRepository)
	authService := auth.NewAuthService(userRepository)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	products.NewProductHandler(router, products.ProductHandlerDeps{
		ProductRepository: productRepository,
		ProductService:    productService,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	fmt.Println("Server listening on :8081")
	server.ListenAndServe()
}
