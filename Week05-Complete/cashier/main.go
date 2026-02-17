package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cashier/database"
	"cashier/handlers"
	"cashier/middlewares"
	"cashier/repositories"
	"cashier/services"

	"github.com/go-chi/chi/v5"
	chi_middleware "github.com/go-chi/chi/v5/middleware"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	type Config struct {
		Port      string `mapstructure:"PORT"`
		DBConn    string `mapstructure:"DB_CONN"`
		JWTSecret string `mapstructure:"JWT_SECRET"`
	}

	config := Config{
		Port:      viper.GetString("PORT"),
		DBConn:    viper.GetString("DB_CONN"),
		JWTSecret: viper.GetString("JWT_SECRET"),
	}

	if config.Port == "" {
		config.Port = "8080"
	}
	if config.JWTSecret == "" {
		config.JWTSecret = "super-secret-key-12345" // Fallback
	}

	log.Printf("Starting server on port %s", config.Port)

	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Repositories
	authRepo := repositories.NewAuthRepository(db)
	categoryRepo := repositories.NewCategoryRepository(db)
	productRepo := repositories.NewProductRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)

	// Initialize Services
	authService := services.NewAuthService(authRepo, config.JWTSecret)
	categoryService := services.NewCategoryService(categoryRepo)
	productService := services.NewProductService(productRepo)
	transactionService := services.NewTransactionService(transactionRepo)

	// Initialize Handlers
	authHandler := handlers.NewAuthHandler(authService)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	productHandler := handlers.NewProductHandler(productService)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	r := chi.NewRouter()

	// Base Middlewares
	r.Use(chi_middleware.Logger)
	r.Use(chi_middleware.Recoverer)
	r.Use(middlewares.CORS)

	// Auth Routes
	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Post("/logout", authHandler.Logout)

	// Public Routes
	r.Get("/categories", categoryHandler.GetAll)
	r.Get("/categories/*", categoryHandler.HandleCategoryByID)
	r.Get("/products", productHandler.GetAll)
	r.Get("/products/*", productHandler.HandleProductByID)

	// Authenticated Routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(authService))

		// Member only can checkout
		r.With(middlewares.RoleMiddleware("member")).Post("/checkout", transactionHandler.Checkout)

		// Admin only
		r.Group(func(r chi.Router) {
			r.Use(middlewares.RoleMiddleware("admin"))
			
			// Categories Management
			r.Post("/categories", categoryHandler.Create)
			r.Put("/categories/*", categoryHandler.HandleCategoryByID)
			r.Delete("/categories/*", categoryHandler.HandleCategoryByID)

			// Products Management
			r.Post("/products", productHandler.Create)
			r.Put("/products/*", productHandler.HandleProductByID)
			r.Delete("/products/*", productHandler.HandleProductByID)

			// Reports
			r.Get("/report", transactionHandler.HandleReport)
			r.Get("/report/today", transactionHandler.HandleReport)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Cashier API (Week 05)"))
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
