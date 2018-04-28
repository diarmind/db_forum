package routes

import (
	"net/http"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"log"
	"github.com/pdmitrya/goExample/databases"
	"fmt"
	"io"
	"database/sql"
	"strconv"
)

func VoteThread(w http.ResponseWriter, r *http.Request, p map[string]string) {

	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)

	var data models.Vote
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	threadSlug := p["slug"]

	dbSession := databases.GetPostgresSession()

	var userId int64
	errGetUserId := dbSession.QueryRow("SELECT id FROM Forum_User WHERE nickname = $1", data.Nickname).Scan(&userId)
	if errGetUserId != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find user with nickname %v", data.Nickname)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}

	var threadId int
	threadId, errId := strconv.Atoi(threadSlug)
	if errId != nil {
		queryGetThreadIdBySlug := "SELECT id FROM Thread WHERE slug = $1"
		errGetThreadIdBySlug := dbSession.QueryRow(queryGetThreadIdBySlug, threadSlug).Scan(&threadId)
		if errGetThreadIdBySlug != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with slug %v", threadSlug)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	} else {
		queryGetThreadIdById := "SELECT id FROM Thread WHERE id = $1"
		errGetThreadIdById := dbSession.QueryRow(queryGetThreadIdById, threadId).Scan(&threadId)
		if errGetThreadIdById != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with id %v", threadId)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	}

	var res sql.Result
	res, err = dbSession.Exec("INSERT INTO Vote VALUES(DEFAULT, $1, $2, $3)", userId, threadId, data.Voice)
	var rowsInserted *int64
	rowsInserted = new(int64)
	if err == nil {
		*rowsInserted, _ = res.LastInsertId()
	}
	switch {
	case err != nil:
		fallthrough
	case rowsInserted == nil:
		queryUpdateVoice := "UPDATE Vote SET value = $1 WHERE user_id = $2 AND thread_id = $3"
		dbSession.QueryRow(queryUpdateVoice, data.Voice, userId, threadId)
	}
	w.WriteHeader(200)

	var thread models.Thread
	query := 	`SELECT U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title, SUM(value) AS votes
				  FROM Thread T JOIN Forum_User U on T.user_id = U.id
				    JOIN Forum F on T.forum_id = F.id LEFT JOIN Vote V on T.id = V.thread_id
				  WHERE T.id = $1
					GROUP BY U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title`
	dbSession.QueryRow(query, threadId).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
	result, _ := json.Marshal(thread)
	io.WriteString(w, string(result))
	return
}
