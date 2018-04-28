package routes

import (
	"net/http"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"log"
	"github.com/pdmitrya/goExample/databases"
	"io"
	"fmt"
	"database/sql"
)

func CreateForum(w http.ResponseWriter, r *http.Request, _ map[string]string) {

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.ForumCreate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()


	var id int64
	err = dbSession.QueryRow("SELECT id FROM Forum_User WHERE nickname = $1", data.User).Scan(&id)
	if err != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find user with nickname %v", data.User)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}

	var res sql.Result
	res, err = dbSession.Exec("INSERT INTO Forum VALUES(DEFAULT, $1, $2, $3)", data.Slug, data.Title, id)
	var rowsInserted *int64
	rowsInserted = new(int64)
	if err == nil {
		*rowsInserted, _ = res.LastInsertId()
	}
	switch {
	case err != nil:
		fallthrough
	case rowsInserted == nil:
		w.WriteHeader(409)

	default:
		w.WriteHeader(201)
	}

	var forum models.Forum
	query := 	`SELECT slug, title, U.nickname,
						(SELECT COUNT(*) FROM Forum JOIN Thread ON Forum.id = Thread.forum_id WHERE Forum.slug = $1) AS threads,
						(SELECT COUNT(*) FROM Forum JOIN Thread ON Forum.id = Thread.forum_id JOIN Post P on Thread.id = P.thread_id WHERE Forum.slug = $1) AS posts
					FROM Forum JOIN Forum_User U on Forum.user_id = U.id
					WHERE slug = $1;`
	dbSession.QueryRow(query, data.Slug).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)
	result, _ := json.Marshal(forum)
	io.WriteString(w, string(result))
	return
}