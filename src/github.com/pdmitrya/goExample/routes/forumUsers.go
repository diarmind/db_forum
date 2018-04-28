package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"fmt"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"io"
	"strings"
	"log"
	"database/sql"
)

func ForumUsers(w http.ResponseWriter, r *http.Request, p map[string]string) {

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
	queryBuilder.WriteString(`SELECT DISTINCT u.id, u.nickname COLLATE "C", u.about, u.fullname, u.email FROM
									(SELECT Thread.user_id FROM Thread WHERE Thread.forum_id = $1 UNION
									SELECT Post.user_id FROM Post JOIN Thread ON Post.thread_id = Thread.id WHERE Thread.forum_id = $1) t JOIN
									  Forum_User u ON t.user_id = u.id `)

	var since string

	desc := r.URL.Query().Get("desc")

	since = r.URL.Query().Get("since")

	limit := r.URL.Query().Get("limit")

	if len(since) > 0 {
		if desc == "true"{
			queryBuilder.WriteString("WHERE LOWER(u.nickname COLLATE \"C\") < LOWER($2 COLLATE \"C\") ")
		} else {
			queryBuilder.WriteString("WHERE LOWER(u.nickname COLLATE \"C\") > LOWER($2 COLLATE \"C\") ")
		}
	} else {
		queryBuilder.WriteString("WHERE TRUE OR $2::TEXT != '' ")
	}

	if desc == "true" {
		queryBuilder.WriteString("ORDER BY u.nickname COLLATE \"C\" DESC ")
	} else {
		queryBuilder.WriteString("ORDER BY u.nickname COLLATE \"C\" ASC ")
	}

	if len(limit) > 0 {
		queryBuilder.WriteString("LIMIT $3 ")
	}

	userArray := models.ArrayOfUsers{}
	var rows *sql.Rows
	if len(limit) > 0 {
		rows, err = dbSession.Query(queryBuilder.String(), forumId, since, limit)
	} else {
		rows, err = dbSession.Query(queryBuilder.String(), forumId, since)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var idUserNotUsed int
		if err := rows.Scan(&idUserNotUsed, &user.Nickname, &user.About, &user.Fullname, &user.Email); err != nil {
			log.Fatal(err)
		}
		userArray = append(userArray, user)
	}
	result, _ := json.Marshal(userArray)
	io.WriteString(w, string(result))
	return
}

