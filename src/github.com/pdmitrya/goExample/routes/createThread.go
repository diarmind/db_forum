package routes

import (
	"net/http"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"log"
	"github.com/pdmitrya/goExample/databases"
	"fmt"
	"io"
)

func CreateThread(w http.ResponseWriter, r *http.Request, p map[string]string) {
	slug := p["slug"]

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.ThreadCreate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()

	var userId int64
	err = dbSession.QueryRow("SELECT id FROM Forum_User WHERE nickname = $1", data.Author).Scan(&userId)
	if err != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find author with nickname %v", data.Author)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}

	var forumId int64
	err = dbSession.QueryRow("SELECT id FROM Forum WHERE slug = $1", &slug).Scan(&forumId)
	if err != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find forum with slug %v", slug)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}


	var insertedRowId int
	err = dbSession.QueryRow("INSERT INTO Thread VALUES(DEFAULT, $1, $2, COALESCE($3, NULL), COALESCE($4, transaction_timestamp()), $5, $6) RETURNING id", data.Message, data.Title, data.Slug, data.Created, forumId, userId).Scan(&insertedRowId)

	switch {
	case err != nil:
		w.WriteHeader(409)
		var thread models.Thread
		query := 	`SELECT U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title,
					(SELECT SUM(value) FROM Thread T JOIN Vote V on T.id = V.thread_id WHERE T.slug = $1) AS votes
				  FROM Thread T JOIN Forum_User U on T.user_id = U.id
				    JOIN Forum F on T.forum_id = F.id
				  WHERE T.slug = $1`
		dbSession.QueryRow(query, data.Slug).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		result, _ := json.Marshal(thread)
		io.WriteString(w, string(result))
		return
	default:
		w.WriteHeader(201)
	}

	var thread models.Thread
	query := 	`SELECT U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title,
					(SELECT SUM(value) FROM Thread T JOIN Vote V on T.id = V.thread_id WHERE T.id = $1) AS votes
				  FROM Thread T JOIN Forum_User U on T.user_id = U.id
				    JOIN Forum F on T.forum_id = F.id
				  WHERE T.id = $1`
	dbSession.QueryRow(query, insertedRowId).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	result, _ := json.Marshal(thread)
	io.WriteString(w, string(result))
	return
}
