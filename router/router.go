package router

import (
	"net/http"

	"github.com/gorilla/mux"

	"golang-crud-api/handlers"
)

func NewRouter(productHandler *handlers.ProductHandler) http.Handler {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/products", productHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("/products", productHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/products/{id:[0-9]+}", productHandler.Delete).Methods(http.MethodDelete)

	return r
}
