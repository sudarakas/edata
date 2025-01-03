package api

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sudarakas/edata/service/user"
)

type APISERVER struct {
	addr string
	db   *sql.DB
}

func NewAPISERVER(addr string, db *sql.DB) *APISERVER {
	return &APISERVER{
		addr: addr,
		db:   db,
	}
}

func (s *APISERVER) Run() error {
	router := mux.NewRouter()
	subRouter := router.PathPrefix("/api/v1").Subrouter()

	userStore, err := user.NewStore(s.db)
	if err != nil {
		log.Fatal(err)
	}

	userHandler := user.NewHandler(userStore)
	userHandler.RegisterRoute(subRouter)

	log.Println("Listening on", s.addr)

	return http.ListenAndServe(s.addr, router)
}
