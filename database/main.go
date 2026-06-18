package database

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jp-milhomem/first-crud/utils"
)

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

//útils

// Init database
func Create() *Database {
	db := Database{
		appData: map[int]User{},
	}

	return &db
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

		utils.SendJSON(w, 201, user.Id)

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

		data := db.appData[int(id64)]

		user, _ := json.Marshal(data)

		utils.SendJSON(w, 200, user)
	}
}

func (db Database) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		utils.SendJSON(w, 200, db.appData)
	}
}

func (db Database) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")

		id, _ := strconv.ParseInt(strId, 10, 64)

		_, ok := db.appData[int(id)]

		if !ok {
			utils.SendJSON(w, 400, utils.Response{
				Err: "usuário não encontrado",
			})
			return
		}

		delete(db.appData, int(id))

		utils.SendJSON(w, 200, db)
	}
}

func (db Database) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id64, _ := strconv.ParseInt(strId, 10, 64)

		user, ok := db.appData[int(id64)]

		if !ok {
			utils.SendJSON(w, 400, "Usuário não encontrado")
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

		utils.SendJSON(w, 200, db)

	}
}
