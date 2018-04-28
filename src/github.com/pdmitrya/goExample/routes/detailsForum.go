package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"github.com/pdmitrya/goExample/models"
	"fmt"
	"encoding/json"
	"io"
)

func DetailsForum(w http.ResponseWriter, r *http.Request, p map[string]string) {
	slug := p["slug"]

	w.Header().Set("Content-Type", "application/json")

	dbSession := databases.GetPostgresSession()

	var forum models.Forum
	query := 	`SELECT slug, title, U.nickname,
						(SELECT COUNT(*) FROM Forum JOIN Thread ON Forum.id = Thread.forum_id WHERE Forum.slug = $1) AS threads,
						(SELECT COUNT(*) FROM Forum JOIN Thread ON Forum.id = Thread.forum_id JOIN Post P on Thread.id = P.thread_id WHERE Forum.slug = $1) AS posts
					FROM Forum JOIN Forum_User U on Forum.user_id = U.id
					WHERE slug = $1;`
	err := dbSession.QueryRow(query, slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts, )
	if err != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find forum with slug %v", slug)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}
	w.WriteHeader(200)
	result, _ := json.Marshal(forum)
	io.WriteString(w, string(result))
	return
}