package handler

import (
	"encoding/json"
	"net/http"

	"hello-go-api/model"
	"hello-go-api/store"
)

const basePathV1 = "/v1"

type ProductHandler struct {
	Store store.ProductStore
}

func (h *ProductHandler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.Store.List()
	if err != nil {
		http.Error(w, `{"error":"failed to fetch products"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	product, err := h.Store.Get(id)
	if err != nil {
		http.Error(w, `{"error":"product not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var p model.Product
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if p.ID == "" || p.Name == "" {
		http.Error(w, `{"error":"id and name are required"}`, http.StatusBadRequest)
		return
	}

	if _, err := h.Store.Get(p.ID); err == nil {
		http.Error(w, `{"error":"product already exists"}`, http.StatusConflict)
		return
	}

	if err := h.Store.Create(p); err != nil {
		http.Error(w, `{"error":"failed to create product"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.Store.Delete(id); err != nil {
		http.Error(w, `{"error":"product not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET "+basePathV1+"/products", h.List)
	mux.HandleFunc("GET "+basePathV1+"/products/{id}", h.Get)
	mux.HandleFunc("POST "+basePathV1+"/products", h.Create)
	mux.HandleFunc("DELETE "+basePathV1+"/products/{id}", h.Delete)
}
