package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"cashier/database"
	"cashier/handlers"
	"cashier/repositories"
	"cashier/services"
	"cashier/middlewares"

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
		Port   string `mapstructure:"PORT"`
		DBConn string `mapstructure:"DB_CONN"`
		API_KEY string `mapstructure:"API_KEY"`
	}

	config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
		API_KEY: viper.GetString("API_KEY"),
	}

	// Default port if not set
	if config.Port == "" {
		config.Port = "8080"
	}

	// Login info for debug
	log.Printf("Starting server on port %s", config.Port)
	log.Printf("Connecting to database...")

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Setup middleware
	apiKeyMiddleware := middlewares.APIKey(config.API_KEY)

	// Middleware chain: Logger -> CORS -> APIKey
	withMiddlewares := func(h http.HandlerFunc) http.HandlerFunc {
		return middlewares.Logger(middlewares.CORS(apiKeyMiddleware(h)))
	}

	// Initialize dependencies (Category)
	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Initialize dependencies (Product)
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// Initialize dependencies (Transaction)
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// Setup routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the Cashier API"))
	})

	// Route for List and Create
	http.HandleFunc("/categories", withMiddlewares(categoryHandler.HandleCategories))
	http.HandleFunc("/products", withMiddlewares(productHandler.HandleProducts))
	
	// Route for GetByID, Update, Delete. 
	// Note: http.HandleFunc matches prefix. "/categories/" will match "/categories/1"
	http.HandleFunc("/categories/", withMiddlewares(categoryHandler.HandleCategoryByID))
	http.HandleFunc("/products/", withMiddlewares(productHandler.HandleProductByID))

	// Route for Transaction or Checkout
	http.HandleFunc("/checkout", withMiddlewares(transactionHandler.HandleCheckout)) // POST
	// Route for Reporting
	http.HandleFunc("/report", withMiddlewares(transactionHandler.HandleReport))       // GET ?start_date=...&end_date=...
	http.HandleFunc("/report/today", withMiddlewares(transactionHandler.HandleReport)) // GET
	
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("gagal running server:", err)
	}
}
