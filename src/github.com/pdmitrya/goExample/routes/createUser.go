package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"log"
	"io"
)

func CreateUser(w http.ResponseWriter, r *http.Request, p map[string]string) {
	nickname := p["nickname"]

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.UserCreate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()

	_, err = dbSession.Exec("INSERT INTO Forum_User VALUES(DEFAULT, $1, $2, $3, $4)", nickname, data.Email, data.Fullname, data.About)
	if err == nil {
		w.WriteHeader(201)
		result, _ := json.Marshal(models.User{ About: data.About, Email: data.Email, Fullname: data.Fullname, Nickname: &nickname})
		io.WriteString(w, string(result))
		return
	}

	var userArray models.ArrayOfUsers
	w.WriteHeader(409)
	rows, err := dbSession.Query("SELECT about, email, nickname, fullname FROM Forum_User WHERE nickname = $1 OR email = $2", nickname, data.Email)
	if err != nil {
        log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var temp models.User
		if err := rows.Scan(&temp.About, &temp.Email, &temp.Nickname, &temp.Fullname); err != nil {
			log.Fatal(err)
		}
		userArray = append(userArray, temp)
	}
	result, _ := json.Marshal(userArray)
	io.WriteString(w, string(result))
	return
}
