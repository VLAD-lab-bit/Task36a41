package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"Task36a41/pkg/storage"

	"github.com/gorilla/mux"
)

// API представляет структуру для API с доступом к хранилищу данных.
type API struct {
	storage *storage.Storage
}

// New создает новый экземпляр API.
func New(storage *storage.Storage) *API {
	return &API{storage: storage}
}

// RegisterRoutes регистрирует маршруты API.
func (api *API) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/news/{n}", api.getLastNPosts).Methods(http.MethodGet)
}

// getLastNPosts обрабатывает запрос для получения последних N публикаций.
func (api *API) getLastNPosts(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	n, err := strconv.Atoi(vars["n"])
	if err != nil {
		http.Error(w, "Invalid number format", http.StatusBadRequest)
		return
	}

	posts, err := api.storage.GetLastNPosts(n)
	if err != nil {
		http.Error(w, "Error retrieving posts", http.StatusInternalServerError)
		log.Println("Error retrieving posts:", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(posts); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
