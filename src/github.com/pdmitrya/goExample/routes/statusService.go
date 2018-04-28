package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"github.com/pdmitrya/goExample/models"
	"encoding/json"
	"io"
)

func StatusService(w http.ResponseWriter, r *http.Request, _ map[string]string) {

	w.Header().Set("Content-Type", "application/json")

	dbSession := databases.GetPostgresSession()

	status := models.Status{}

	statusQuery := `SELECT (SELECT COUNT(*) FROM Forum) AS forum,
					(SELECT COUNT(*) FROM Post) AS post,
  					(SELECT COUNT(*) FROM Thread) AS thread,
  					(SELECT COUNT(*) FROM Forum_User) AS user`
	dbSession.QueryRow(statusQuery).Scan(&status.Forum, &status.Post, &status.Thread, &status.User)

	w.WriteHeader(200)
	result, _ := json.Marshal(status)
	io.WriteString(w, string(result))
	return
}
