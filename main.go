package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"os"
	"log"
	"github.com/spf13/viper"
	"kasir-api/database"
	"kasir-api/repositories"
    "kasir-api/services"
    "kasir-api/handlers"
)



type Config struct {
	Port    string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}






func main() {

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
	 	Port: viper.GetString("PORT"),
	 	DBConn: viper.GetString("DB_CONN"),
	}


	log.Println("DB_CONN =", config.DBConn)

	// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"message": "halo bang...",
			"status":  "running",
		})
	})


	// localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
    categoryService := services.NewCategoryService(categoryRepo)
    categoryHandler := handlers.NewCategoryHandler(categoryService)

	//Produk:
	http.HandleFunc("/api/produk", productHandler.HandleProducts)
	http.HandleFunc("/api/produk/", productHandler.HandleProductByID)

	//kategori
	http.HandleFunc("/api/categories", categoryHandler.HandleCategorys) 
	http.HandleFunc("/api/categories/", categoryHandler.HandleCategoryByID) 


	//running
	addr := "0.0.0.0:" + config.Port
	fmt.Println("Server running di", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println("gagal running server", err)
	}


}