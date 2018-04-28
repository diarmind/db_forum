package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"strconv"
	"github.com/pdmitrya/goExample/models"
	"fmt"
	"encoding/json"
	"io"
	"log"
)

func DetailsThreadPost(w http.ResponseWriter, r *http.Request, p map[string]string) {

	w.Header().Set("Content-Type", "application/json")

	slug := p["slug"]

	decoder := json.NewDecoder(r.Body)

	var data models.ThreadUpdate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()

	id, err := strconv.Atoi(slug)
	var query string
	var thread models.Thread

	if err != nil {
		query = 	`UPDATE Thread SET message = COALESCE($2, message), title = COALESCE($3, title)
				  		WHERE slug = $1 RETURNING id`
		err = dbSession.QueryRow(query, slug, data.Message, data.Title).Scan(&id)
		if err != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with slug %v", slug)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	} else {
		query = 	`UPDATE Thread SET message = COALESCE($2, message), title = COALESCE($3, title)
				  		WHERE id = $1`
		_, err = dbSession.Exec(query, id, data.Message, data.Title)
		if err != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with id %v", id)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	}

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

	result, _ := json.Marshal(thread)
	io.WriteString(w, string(result))
	return

}
