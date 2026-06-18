package server

import (
	"fmt"
	"net/http"
	"time"

	storage "github.com/jp-milhomem/first-crud/database"
	"github.com/jp-milhomem/first-crud/handlers"
)

type Application struct {
	data map[int]handlers.User
}

func Init() error {

	db := storage.Database{}
	_ = db

	app := Application{
		data: make(map[int]handlers.User),
	}

	app.data[1] = handlers.User{Id: "1", FirstName: "admin", LastName: "admin", Biography: "The manager"}

	handler := handlers.NewHandler()

	s := http.Server{
		Addr:           "localhost:8082",
		Handler:        handler,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   5 * time.Second,
		IdleTimeout:    1 * time.Minute,
		MaxHeaderBytes: 0,
	}

	if err := s.ListenAndServe(); err != nil {
		return err
	}

	fmt.Println("Server is on...")
	return nil

}
