package server

import (
	"net/http"
	"time"

	"github.com/jp-milhomem/first-crud/handlers"
)

func Create() error {
	mux := handlers.NewHandler()

	s := http.Server{
		Addr:           "localhost:8082",
		Handler:        mux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    1 * time.Minute,
		MaxHeaderBytes: 0,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	return nil

}
