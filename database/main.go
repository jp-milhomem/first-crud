package storage

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
)

//útils

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

// types
type Database struct {
	appData map[int]User
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

//Méthods

func (db Database) Insert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := &User{}

		//err = json.Unmarshal(data, user)
		err = json.NewDecoder(r.Body).Decode(user)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		fmt.Println(user)

		id, _ := strconv.ParseInt(user.Id, 10, 64)

		db.appData[int(id)] = *user

		sendJSON(w, 201, "Msg: Usuário inserido com sucesso")

	}
}

func (db Database) FIndById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		id64, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			slog.Error(err.Error())
			return
		}

		user := db.appData[int(id64)]

		data, _ := json.Marshal(user)

		w.Write(data)
	}
}

func (db Database) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sendJSON(w, 200, db.appData)
	}
}

func (db Database) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")

		id, _ := strconv.ParseInt(strId, 10, 64)

		_, ok := db.appData[int(id)]

		if !ok {
			sendJSON(w, 400, "Usuário não encontrado")
			return
		}

		delete(db.appData, int(id))

		sendJSON(w, 200, db)
	}
}

func (db Database) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id64, _ := strconv.ParseInt(strId, 10, 64)

		user, ok := db.appData[int(id64)]

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

		db.appData[int(id64)] = user

		sendJSON(w, 200, db)

	}
}
