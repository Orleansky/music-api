package api

import (
	"Anastasia/songs/internal/services"
	"net/http"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type API struct {
	srv    *services.Service
	router *mux.Router
}

func New(srv *services.Service) *API {
	api := &API{
		srv:    srv,
		router: mux.NewRouter(),
	}

	api.endpoints()
	return api
}

func (api *API) Router() *mux.Router {
	return api.router
}

func (api *API) endpoints() {
	api.router.HandleFunc("/songs", api.songsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/songs/{id}", api.songByIDHandler).Methods(http.MethodGet, http.MethodOptions)
	api.router.HandleFunc("/songs/{id}", api.deleteSongHandler).Methods(http.MethodDelete, http.MethodOptions)
	api.router.HandleFunc("/songs/{id}", api.updateSongHandler).Methods(http.MethodPatch, http.MethodOptions)
	api.router.HandleFunc("/songs", api.createSongHandler).Methods(http.MethodPost, http.MethodOptions)
	api.router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}
