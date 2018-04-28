package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/models"
	"log"
	"encoding/json"
	"io"
	"strings"
	"fmt"
	"github.com/pdmitrya/goExample/databases"
)

func ForumThreads(w http.ResponseWriter, r *http.Request, p map[string]string) {

	slug := p["slug"]

	dbSession := databases.GetPostgresSession()

	w.Header().Set("Content-Type", "application/json")

	var forumId int64
	err := dbSession.QueryRow("SELECT id, slug FROM Forum WHERE slug = $1", &slug).Scan(&forumId, &slug)
	if err != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find forum with slug %v", slug)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}


	var queryBuilder strings.Builder
	queryBuilder.WriteString(`SELECT U.nickname, T.created, T.id, T.message, T.slug, T.title, SUM(value) AS votes
								FROM Thread T JOIN Forum_User U on T.user_id = U.id LEFT JOIN Vote V on T.id = V.thread_id
								WHERE T.forum_id = $1 `)

	var since string

	desc := r.URL.Query().Get("desc")

	since = r.URL.Query().Get("since")

	if len(since) > 0 {
		if desc == "true"{
			queryBuilder.WriteString("AND T.created <= $2 ")
		} else {
			queryBuilder.WriteString("AND T.created >= $2 ")
		}
	} else {
		queryBuilder.WriteString("OR $2::TEXT != ''")
	}

	queryBuilder.WriteString(`GROUP BY U.nickname, T.created, T.id, T.message, T.slug, T.title `)

	if desc == "true" {
		queryBuilder.WriteString("ORDER BY T.created DESC ")
	} else {
		queryBuilder.WriteString("ORDER BY T.created ASC ")
	}

	limit := r.URL.Query().Get("limit")
	if limit != "" {
		queryBuilder.WriteString("LIMIT $3 ")
	}

	threadArray := models.ArrayOfThreads{}
	rows, err := dbSession.Query(queryBuilder.String(), forumId, since, limit) //since,
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var thread models.Thread
		if err := rows.Scan(&thread.Author, &thread.Created, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes); err != nil {
			log.Fatal(err)
		}
		thread.Forum = &slug
		threadArray = append(threadArray, thread)
	}
	result, _ := json.Marshal(threadArray)
	io.WriteString(w, string(result))
	return
}
