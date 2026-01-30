// ubah Config
type Config struct {
	Port   string `mapstructure:"PORT"`
	DBConn string `mapstructure:"DB_CONN"`
}

config := Config{
		Port:   viper.GetString("PORT"),
		DBConn: viper.GetString("DB_CONN"),
}

// Setup database
	db, err := database.InitDB(config.DBConn)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

// Setup routes
http.HandleFunc("/api/produk", productHandler.HandleProducts)
http.HandleFunc("/api/produk/", productHandler.HandleProductByID)


//Dependency Injection
productRepo := repositories.NewProductRepository(db)
productService := services.NewProductService(productRepo)
productHandler := handlers.NewProductHandler(productService)