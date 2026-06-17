package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

func Insert(db map[int]User) http.HandlerFunc {
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

func NewHandler(db map[int]User) http.Handler {
	r := chi.NewMux()

	r.Use(middleware.Recoverer)
	r.Use(middleware.Logger)

	//Routes

	//health check
	r.Get("/check", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, 200, "hello word")
	})

	// FInd user
	r.Get("/api/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		id64, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		user := db[int(id64)]

		data, _ := json.Marshal(user)

		w.Write(data)
	})

	//List all users
	r.Get("/api/users", func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, 200, db)
	})

	//Create a user
	r.Post("/api/users", Insert(db))

	r.Delete("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")

		id, _ := strconv.ParseInt(strId, 10, 64)

		_, ok := db[int(id)]

		if !ok {
			sendJSON(w, 400, "Usuário não encontrado")
			return
		}

		delete(db, int(id))

		sendJSON(w, 200, db)
	})

	r.Put("/api/users/{id}", func(w http.ResponseWriter, r *http.Request) {
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

	})

	return r
}
