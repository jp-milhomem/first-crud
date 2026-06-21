package database

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/jp-milhomem/first-crud/utils"
)

// types
type Database struct {
	appData map[string]User
}

type User struct {
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
		appData: map[string]User{},
	}

	return &db
}

//Méthods

func (db Database) Insert() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		user := &User{}

		err = json.NewDecoder(r.Body).Decode(user)

		if err != nil {
			utils.SendJSON(w, http.StatusBadRequest, utils.Response{
				Err: "Invalid body",
			})
			return
		}

		id := uuid.New().String()
		db.appData[id] = *user

		utils.SendJSON(w, http.StatusCreated, utils.Response{
			Data: map[string]string{
				"user_id": id,
				"Message": "User created successfully",
			},
		})

	}
}

func (db Database) FindById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		user, ok := db.appData[id]

		if !ok {
			utils.SendJSON(w, http.StatusNotFound, utils.Response{
				Err: "Usuário não encontrado",
			})

			return
		}

		utils.SendJSON(w, http.StatusOK, utils.Response{
			Data: user,
		})
	}
}

func (db Database) FindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		users := db.appData

		if len(db.appData) < 1 {
			utils.SendJSON(w, http.StatusOK, utils.Response{
				Data: "Empty database",
			})
			return
		}

		utils.SendJSON(w, http.StatusOK, utils.Response{
			Data: users,
		})

	}

}

func (db Database) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		_, ok := db.appData[id]

		if !ok {
			utils.SendJSON(w, http.StatusNotFound, utils.Response{
				Err: "User not found",
			})
			return
		}

		delete(db.appData, id)

		utils.SendJSON(w, http.StatusAccepted, utils.Response{
			Data: "Deleted with success",
		})

	}

}

func (db Database) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")

		user, ok := db.appData[id]

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

			return
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

		db.appData[id] = user

		utils.SendJSON(w, http.StatusCreated, utils.Response{
			Data: user,
		})

	}
}
