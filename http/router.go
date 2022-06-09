package http

import (
	"net/http"

	"github.com/gorilla/mux"
)

// NewRouter creates router
func NewRouter(h *Handler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/statement", h.ProcessStatement).Methods(http.MethodPost)
	router.HandleFunc("/ping", Ping)

	return router
}
