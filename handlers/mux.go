package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jp-milhomem/first-crud/database"
)

//Tipagens

type WriterData struct {
	Data any
}

type User struct {
	Id        string `json:"id,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Biography string `json:"biography,omitempty"`
}

type UpdateUser struct {
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	Biography string `json:"biography,omitempty"`
}

type Database map[int]User

// Functions

func sendJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(status)

	resp, err := json.Marshal(data)

	if err != nil {
		fmt.Println(err)
		return
	}

	w.Write(resp)
}

//Métodos

func Insert(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := &User{}

		//err = json.Unmarshal(data, user)
		err = json.NewDecoder(r.Body).Decode(user)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		id, _ := strconv.ParseInt(user.Id, 10, 64)

		db[int(id)] = *user

		sendJSON(w, 201, WriterData{Data: user})

	}
}

func FIndById(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		id64, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		user := db[int(id64)]

		data, _ := json.Marshal(user)

		w.Write(data)
	}
}

func FindAll(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, 200, db)
	}
}

func Delete(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")

		id, _ := strconv.ParseInt(strId, 10, 64)

		_, ok := db[int(id)]

		if !ok {
			sendJSON(w, 400, "Usuário não encontrado")
			return
		}

		delete(db, int(id))

		sendJSON(w, 200, db)
	}
}

func Update(db Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id64, _ := strconv.ParseInt(strId, 10, 64)

		user, ok := db[int(id64)]

		if !ok {
			sendJSON(w, 400, "Usuário não encontrado")
			return
		}

		updateUser := UpdateUser{}
		json.NewDecoder(r.Body).Decode(&updateUser)

		if updateUser.Biography != "" {
			user.Biography = updateUser.Biography
		}

		if updateUser.FirstName != "" {
			user.FirstName = updateUser.FirstName
		}

		if updateUser.LastName != "" {
			user.LastName = updateUser.LastName
		}

		db[int(id64)] = user

		sendJSON(w, 200, db)

	}
}

// Handler
func NewHandler() http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	database := database.Create()

	//Routes

	// Find user
	r.Get("/api/user/{id}", database.FIndById())

	//List all users
	r.Get("/api/users", database.FindAll())

	//Create a user
	r.Post("/api/users", database.Insert())

	r.Delete("/api/users/{id}", database.Delete())

	r.Put("/api/users/{id}", database.Update())

	return r
}
