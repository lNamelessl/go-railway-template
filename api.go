package main

import (
	"encoding/json"
	"net/http"
)

type item struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type createItemRequest struct {
	Name string `json:"name"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /healthz", handleHealthz)
	mux.HandleFunc("POST /items", handleCreateItem)
	mux.HandleFunc("GET /items", handleListItems)
	mux.HandleFunc("GET /items/{id}", handleGetItem)
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleCreateItem(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid JSON body"})
		return
	}
	if req.Name == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "name is required"})
		return
	}
	created := store.create(req.Name)
	writeJSON(w, http.StatusCreated, created)
}

func handleListItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string][]item{"items": store.list()})
}

func handleGetItem(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	it, ok := store.get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, errorResponse{Error: "item not found"})
		return
	}
	writeJSON(w, http.StatusOK, it)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
