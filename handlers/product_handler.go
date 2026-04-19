package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"golang-crud-api/models"
	"golang-crud-api/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// GET /api/v1/products
func (h *ProductHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	products, err := h.repo.GetAll()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if products == nil {
		products = []models.Product{}
	}
	writeJSON(w, http.StatusOK, models.Response{Data: products})
}

// GET /api/v1/products/{id}
func (h *ProductHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "id tidak valid"})
		return
	}

	product, err := h.repo.GetByID(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if product == nil {
		writeJSON(w, http.StatusNotFound, models.Response{Error: "produk tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, models.Response{Data: product})
}

// POST /api/v1/products
func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "request body tidak valid"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "nama produk wajib diisi"})
		return
	}
	if req.Price < 0 {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "harga tidak boleh negatif"})
		return
	}

	product, err := h.repo.Create(req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, models.Response{
		Message: "produk berhasil dibuat",
		Data:    product,
	})
}

// PUT /api/v1/products/{id}
func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "id tidak valid"})
		return
	}

	var req models.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "request body tidak valid"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "nama produk wajib diisi"})
		return
	}

	product, err := h.repo.Update(id, req)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if product == nil {
		writeJSON(w, http.StatusNotFound, models.Response{Error: "produk tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, models.Response{
		Message: "produk berhasil diupdate",
		Data:    product,
	})
}

// DELETE /api/v1/products/{id}
func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeJSON(w, http.StatusBadRequest, models.Response{Error: "id tidak valid"})
		return
	}

	deleted, err := h.repo.Delete(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, models.Response{Error: err.Error()})
		return
	}
	if !deleted {
		writeJSON(w, http.StatusNotFound, models.Response{Error: "produk tidak ditemukan"})
		return
	}
	writeJSON(w, http.StatusOK, models.Response{Message: "produk berhasil dihapus"})
}
