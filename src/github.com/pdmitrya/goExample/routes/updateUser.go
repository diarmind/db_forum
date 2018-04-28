package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"github.com/pdmitrya/goExample/models"
	"encoding/json"
	"io"
	"log"
	"database/sql"
	"fmt"
)

func UpdateUser(w http.ResponseWriter, r *http.Request, p map[string]string) {
	nickname := p["nickname"]

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.UserUpdate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()

	var res sql.Result
	res, err = dbSession.Exec("UPDATE Forum_User SET about = COALESCE($1, about), email = COALESCE($2, email), fullname = COALESCE($3, fullname) WHERE nickname = $4", data.About, data.Email, data.Fullname, nickname)
	var rowsAffected int64
	if err == nil {
		rowsAffected, _ = res.RowsAffected()
	}
	switch {
	case err != nil:
		w.WriteHeader(409)
		errorMessage := "New data conflicts with another user"
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	case rowsAffected == 0:
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find user with nickname %s", nickname)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	default:
		w.WriteHeader(200)
		var user models.User
		dbSession.QueryRow("SELECT about, email, nickname, fullname FROM Forum_User WHERE nickname = $1", nickname).Scan(&user.About, &user.Email, &user.Nickname, &user.Fullname)
		result, _ := json.Marshal(user)
		io.WriteString(w, string(result))
		return
	}
}
