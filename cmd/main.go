package main

import (
	"bike/configs"
	_ "bike/docs"
	"bike/internal/addresses"
	"bike/internal/auth"
	"bike/internal/products"
	"bike/internal/users"
	"bike/pkg/db"
	"bike/pkg/middleware"
	"fmt"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
)

// @title API-Bike
// @version 1.0
// @description API — сервис для управления пользователями и продуктами для проекта bike.
// @contact.name Andrew Voronkov
// @contact.email voronkovworkemail@gmail.com
// @host localhost:8081
// @BasePath /
func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	productRepository := products.NewProductRepository(database)
	userRepository := users.NewUserRepository(database)
	addressRepository := addresses.NewAddressRepository(database)

	// Services
	productService := products.NewProductService(productRepository)
	authService := auth.NewAuthService(userRepository)
	addressService := addresses.NewAddressService(addressRepository, userRepository)

	// Handlers
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	products.NewProductHandler(router, products.ProductHandlerDeps{
		Config:            conf,
		ProductRepository: productRepository,
		ProductService:    productService,
	})
	addresses.NewAddressHandler(router, addresses.AddressHandlerDeps{
		Config:            conf,
		AddressRepository: addressRepository,
		AddressService:    addressService,
		UserRepository:    userRepository,
	})
	users.NewUsersHandler(router, users.UserHandlerDeps{
		Config:         conf,
		UserRepository: userRepository,
	})

	// Swagger UI
	router.Handle("/swagger/", httpSwagger.WrapHandler)

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
