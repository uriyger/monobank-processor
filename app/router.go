package app

import (
	"errors"
	"net/http"

	"github.com/gorilla/mux"

	apphttp "monobank-processor/app/http"
)

func (a *App) initRouter() error {
	if a.handler == nil {
		return errors.New("router initialization: handler is not initialized")
	}
	a.router = mux.NewRouter()
	a.router.HandleFunc("/statement", a.handler.ProcessStatement).Methods(http.MethodPost)
	a.router.HandleFunc("/ping", apphttp.Ping)

	return nil
}
