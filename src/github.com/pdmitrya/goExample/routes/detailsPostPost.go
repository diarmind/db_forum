package routes

import (
	"io"
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"log"
	"github.com/pdmitrya/goExample/databases"
)

func DetailsPostPost(w http.ResponseWriter, r *http.Request, p map[string]string) {

	w.Header().Set("Content-Type", "application/json")

	id := p["id"]

	decoder := json.NewDecoder(r.Body)

	var data models.PostUpdate
	err := decoder.Decode(&data)
	if err != nil {
		log.Fatal(err)
	}

	dbSession := databases.GetPostgresSession()

	var postId int

	var nothingToUpdate bool
	queryMessageCheck := `SELECT message SIMILAR TO $2 OR $2 IS NULL, id FROM Post WHERE id = $1`
	errMessageCheck := dbSession.QueryRow(queryMessageCheck, id, data.Message).Scan(&nothingToUpdate, &postId)
	if errMessageCheck != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find post with id %v", id)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}
	if !nothingToUpdate {
		queryUpdatePost := 	`UPDATE Post SET message = $2, isEdited = TRUE
				  		WHERE id = $1 RETURNING id`
		err = dbSession.QueryRow(queryUpdatePost, id, data.Message).Scan(&postId)
		// it couldn`t be
		if err != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find post with id %v", postId)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
	}

	post := models.Post{}
	querySelectPost := 	`SELECT U.nickname, P.created, F.slug, P.id, P.isEdited, P.message, P.parent, P.thread_id
							FROM Post P JOIN Thread T on P.thread_id = T.id JOIN Forum F on T.forum_id = F.id JOIN Forum_User U on P.user_id = U.id
							WHERE P.id = $1`
	errSelectPost := dbSession.QueryRow(querySelectPost, postId).Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread)
	if errSelectPost != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find thread with id %v", postId)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}

	result, _ := json.Marshal(post)
	io.WriteString(w, string(result))
	return

}
