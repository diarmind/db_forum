package routes

import (
	"github.com/pdmitrya/goExample/models"
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"encoding/json"
	"io"
	"strconv"
	"fmt"
)

func DetailsThreadGet(w http.ResponseWriter, r *http.Request, p map[string]string) {

	slug := p["slug"]

	dbSession := databases.GetPostgresSession()

	w.Header().Set("Content-Type", "application/json")
	id, err := strconv.Atoi(slug)
	var query string
	var thread models.Thread
	if err != nil {
		query = 	`SELECT U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title, SUM(value) AS votes
				  FROM Thread T JOIN Forum_User U on T.user_id = U.id
				    JOIN Forum F on T.forum_id = F.id LEFT JOIN Vote V on T.id = V.thread_id
				  WHERE T.slug = $1
					GROUP BY U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title`
		err = dbSession.QueryRow(query, slug).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with slug %v", slug)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	} else {
		query = 	`SELECT U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title, SUM(value) AS votes
				  FROM Thread T JOIN Forum_User U on T.user_id = U.id
				    JOIN Forum F on T.forum_id = F.id LEFT JOIN Vote V on T.id = V.thread_id
				  WHERE T.id = $1
					GROUP BY U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title`
		err = dbSession.QueryRow(query, id).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if err != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with id %v", slug)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	}

	result, _ := json.Marshal(thread)
	io.WriteString(w, string(result))
	return

}
