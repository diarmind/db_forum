package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
)

func ClearService(w http.ResponseWriter, r *http.Request, _ map[string]string) {

	w.Header().Set("Content-Type", "application/json")

	dbSession := databases.GetPostgresSession()

	truncateAllQuery := "TRUNCATE Vote, Post, Thread, Forum, Forum_User"
	dbSession.Exec(truncateAllQuery)

	w.WriteHeader(200)
	return
}

