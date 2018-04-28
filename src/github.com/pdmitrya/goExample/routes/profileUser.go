package routes

import (
	"net/http"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"github.com/pdmitrya/goExample/databases"
	"io"
	"fmt"
)

func ProfileUser(w http.ResponseWriter, r *http.Request, p map[string]string) {
	nickname := p["nickname"]

	w.Header().Set("Content-Type", "application/json")

	dbSession := databases.GetPostgresSession()

	var user models.User
	err := dbSession.QueryRow(	"SELECT about, email, nickname, fullname FROM Forum_User WHERE nickname = $1", nickname).Scan(&user.About, &user.Email, &user.Nickname, &user.Fullname)
	if err != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find user with nickname %v", nickname)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}
	w.WriteHeader(200)
	result, _ := json.Marshal(user)
	io.WriteString(w, string(result))
	return
}