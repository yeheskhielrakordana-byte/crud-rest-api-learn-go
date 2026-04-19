package main

import (
	"log"
	"net/http"
	"os"

	"golang-crud-api/config"
	"golang-crud-api/handlers"
	"golang-crud-api/repository"
	"golang-crud-api/router"
)

func main() {
	db := config.ConnectDB()
	defer db.Close()

	productRepo := repository.NewProductRepository(db)
	productHandler := handlers.NewProductHandler(productRepo)

	r := router.NewRouter(productHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan di port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}
