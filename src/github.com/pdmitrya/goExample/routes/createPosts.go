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

func CreatePosts(w http.ResponseWriter, r *http.Request, p map[string]string) {

	slug := p["slug"]
	decoder := json.NewDecoder(r.Body)

	var data models.ArrayOfPostsCreate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()
	postTransaction, _ := dbSession.Begin()

	w.Header().Set("Content-Type", "application/json")

	threadId, err := strconv.Atoi(slug)

	if err != nil {
		queryGetThreadId := `SELECT id FROM Thread WHERE slug = $1`
		errGetThreadId := postTransaction.QueryRow(queryGetThreadId, slug).Scan(&threadId)
		if errGetThreadId != nil {
			postTransaction.Rollback()
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with slug %v", slug)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	} else {
		queryCheckThreadExists := `SELECT id FROM Thread WHERE id = $1`
		errCheckThreadExists := postTransaction.QueryRow(queryCheckThreadExists, threadId).Scan(&threadId)
		if errCheckThreadExists != nil {
			postTransaction.Rollback()
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with id %v", threadId)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	}

	posts := models.ArrayOfPosts{}
	for i := 0; i < len(data); i++ {
		var userId int
		queryGetUserId := `SELECT id FROM Forum_User FU WHERE FU.nickname = $1`
		errGetUserId := postTransaction.QueryRow(queryGetUserId, data[i].Author).Scan(&userId)
		if errGetUserId != nil {
			postTransaction.Rollback()
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find user with nickname %v", *data[i].Author)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}

		if data[i].Parent != nil {
			var parentThreadId int
			queryGetThreadId := `SELECT thread_id FROM Post WHERE id = $1`
			errGetThreadId := postTransaction.QueryRow(queryGetThreadId, data[i].Parent).Scan(&parentThreadId)
			if parentThreadId != threadId {
				postTransaction.Rollback()
				w.WriteHeader(409)
				errorMessage := "Parent post was created in another thread"
				result, _ := json.Marshal(models.Error{Message: &errorMessage})
				io.WriteString(w, string(result))
				return
			}
			if errGetThreadId != nil {
				postTransaction.Rollback()
				w.WriteHeader(404)
				errorMessage := fmt.Sprintf("Can`t find post with id %v!", *data[i].Parent)
				result, _ := json.Marshal(models.Error{Message: &errorMessage})
				io.WriteString(w, string(result))
				return
			}
		}

		var postId int
		queryInsertPost := `INSERT INTO Post VALUES(DEFAULT, $1, $2, $3, DEFAULT, DEFAULT, $4) RETURNING id`
		errInsertPost := postTransaction.QueryRow(queryInsertPost, data[i].Parent, &userId, &threadId, data[i].Message).Scan(&postId)
		// TODO WHAT ERROR IS HERE?
		if errInsertPost != nil {
			postTransaction.Rollback()
			w.WriteHeader(404)
			errorMessage := "Transaction conflict, no persistence (seems like this)!"
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}

		post := models.Post{}

		querySelectPost := `SELECT U.nickname, P.created, F.slug, P.id, P.message, P.parent, P.thread_id
							FROM Post P JOIN Thread T on P.thread_id = T.id JOIN Forum F on T.forum_id = F.id JOIN Forum_User U on P.user_id = U.id
							WHERE P.id = $1`
		errSelectPost := postTransaction.QueryRow(querySelectPost, postId).Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.Message, &post.Parent, &post.Thread)
		// It couldn`t be, but i`ll check
		if errSelectPost != nil {
			postTransaction.Rollback()
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find post with id %v", slug)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
		posts = append(posts, post)
	}
	postTransaction.Commit()
	w.WriteHeader(201)
	result, _ := json.Marshal(posts)
	io.WriteString(w, string(result))
	//w.Write(result)
	return

}
