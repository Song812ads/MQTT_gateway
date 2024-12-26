package api

import (
	"net/http"

	"github.com/gateway/service"
	"github.com/gorilla/mux"
)

type APIServer struct {
	addr string
}

func NewAPIServer(s string) *APIServer {
	return &APIServer{addr: s}
}

func (s *APIServer) Run() error {

	router := mux.NewRouter()
	subrouter := router.PathPrefix("/api/gateway").Subrouter()

	deviceService := service.NewService()
	subrouter.HandleFunc("/device", deviceService.AddDevice).Methods("POST")
	// subrouter.HandleFunc("/device", deviceService.EraseDevice).Methods("DELETE")
	// subrouter.HandleFunc("/device", deviceService.UpdateDevice).Methods("PATCH")
	// subrouter.HandleFunc("/device/all", deviceService.GetAllDevice).Methods("GET")
	// subrouter.HandleFunc("/device", deviceService.GetDeviceByName).Methods("GET")
	// subrouter.HandleFunc("/profile", deviceService.AddProfile).Methods("POST")

	// handler := cors.Default().Handler(router)
	return http.ListenAndServe(s.addr, router)
}
