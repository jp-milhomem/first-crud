package database

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jp-milhomem/first-crud/utils"
)

// types
type Database struct {
	appData map[int]User
}

type User struct {
	Id        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Biography string `json:"biography"`
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
			utils.SendJSON(w, http.StatusBadRequest, utils.Response{
				Err: "Invalid body",
			})
			return
		}

		id, err := strconv.ParseInt(user.Id, 10, 64)

		if err != nil {
			utils.SendJSON(w, http.StatusNotFound, utils.Response{
				Err: "Invalid Id",
			})
		}

		db.appData[int(id)] = *user

		utils.SendJSON(w, http.StatusCreated, utils.Response{
			Data: user.Id,
		})

	}
}

func (db Database) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		id64, err := strconv.ParseInt(id, 10, 64)

		if err != nil {
			utils.SendJSON(w, http.StatusBadRequest, utils.Response{
				Err: "Invalid Id",
			})
			return
		}

		user, ok := db.appData[int(id64)]

		if !ok {
			utils.SendJSON(w, http.StatusNotFound, utils.Response{
				Err: "Usuário não encontrado",
			})
		}

		utils.SendJSON(w, http.StatusFound, utils.Response{
			Data: user,
		})
	}
}

func (db Database) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.appData

		utils.SendJSON(w, http.StatusFound, utils.Response{
			Data: users,
		})
	}
}

func (db Database) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")

		id, _ := strconv.ParseInt(strId, 10, 64)

		_, ok := db.appData[int(id)]

		if !ok {
			utils.SendJSON(w, http.StatusNotFound, utils.Response{
				Err: "User not found",
			})
			return
		}

		delete(db.appData, int(id))

		utils.SendJSON(w, http.StatusAccepted, utils.Response{
			Data: "Deleted with success",
		})
	}
}

func (db Database) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		strId := r.PathValue("id")
		id64, _ := strconv.ParseInt(strId, 10, 64)

		user, ok := db.appData[int(id64)]

		if !ok {
			utils.SendJSON(w, 400, utils.Response{
				Err: "Usuário não encontrado",
			})

			return
		}

		updateUser := UpdateUser{}

		err := json.NewDecoder(r.Body).Decode(&updateUser)

		if err != nil {
			utils.SendJSON(w, http.StatusBadRequest, utils.Response{
				Err: "Invalid body",
			})
		}

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

		utils.SendJSON(w, http.StatusCreated, utils.Response{
			Data: user,
		})

	}
}
