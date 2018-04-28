package routes

import (
	"net/http"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"github.com/pdmitrya/goExample/databases"
	"fmt"
	"io"
	"strings"
)

func DetailsPostGet(w http.ResponseWriter, r *http.Request, p map[string]string) {

	w.Header().Set("Content-Type", "application/json")

	related := r.URL.Query().Get("related")
	strArr := strings.Split(related, ",")

	postFull := models.PostFull{}
	var (
		author models.User
		forum models.Forum
		post models.Post
		thread models.Thread
	)

	dbSession := databases.GetPostgresSession()

	postId := p["id"]

	var userId, threadId, forumId int
	querySelectPostWithFK := 	`SELECT U.nickname, P.created, F.slug, P.id, P.isEdited, P.message, P.parent, P.thread_id, P.user_id, P.thread_id, T.forum_id
							FROM Post P JOIN Thread T on P.thread_id = T.id JOIN Forum F on T.forum_id = F.id JOIN Forum_User U on P.user_id = U.id
							WHERE P.id = $1`
	errSelectPostWithFK := dbSession.QueryRow(querySelectPostWithFK, postId).Scan(&post.Author, &post.Created, &post.Forum, &post.Id, &post.IsEdited, &post.Message, &post.Parent, &post.Thread, &userId, &threadId, &forumId)
	if errSelectPostWithFK != nil {
		w.WriteHeader(404)
		errorMessage := fmt.Sprintf("Can`t find post with id %v", postId)
		result, _ := json.Marshal(models.Error{Message: &errorMessage})
		io.WriteString(w, string(result))
		return
	}

	mapFull := make(map[string]bool)
	for i := 0; i < len(strArr); i++ {
		mapFull[strArr[i]] = true
	}

	if mapFull["user"] {
		querySelectPostAuthor := `SELECT U.nickname, U.email, U.fullname, U.about 
									FROM Forum_User U
									WHERE U.id = $1`
		errSelectPostAuthor := dbSession.QueryRow(querySelectPostAuthor, userId).Scan(&author.Nickname, &author.Email, &author.Fullname, &author.About)
		if errSelectPostAuthor != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find user with id %v", userId)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
		postFull.Author = &author
	}

	if mapFull["thread"] {
		querySelectPostThread :=
			`SELECT U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title, SUM(value) AS votes
				  FROM Thread T JOIN Forum_User U on T.user_id = U.id
				    JOIN Forum F on T.forum_id = F.id LEFT JOIN Vote V on T.id = V.thread_id
				  WHERE T.id = $1
					GROUP BY U.nickname, T.created, F.slug, T.id, T.message, T.slug, T.title`
		errSelectPostAuthor := dbSession.QueryRow(querySelectPostThread, threadId).Scan(&thread.Author, &thread.Created, &thread.Forum, &thread.Id, &thread.Message, &thread.Slug, &thread.Title, &thread.Votes)
		if errSelectPostAuthor != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find thread with id %v", threadId)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
		postFull.Thread = &thread
	}

	if mapFull["forum"] {
		querySelectPostForum :=
						`SELECT F.slug, F.title, U.nickname,
						(SELECT COUNT(*) FROM Forum JOIN Thread ON Forum.id = Thread.forum_id WHERE Forum.id = $1) AS threads,
						(SELECT COUNT(*) FROM Forum JOIN Thread ON Forum.id = Thread.forum_id JOIN Post P on Thread.id = P.thread_id WHERE Forum.id = $1) AS posts
							FROM Forum F JOIN Forum_User U on F.user_id = U.id
						WHERE F.id = $1;`
		errSelectPostAuthor := dbSession.QueryRow(querySelectPostForum, forumId).Scan(&forum.Slug, &forum.Title, &forum.User, &forum.Threads, &forum.Posts)
		if errSelectPostAuthor != nil {
			w.WriteHeader(404)
			errorMessage := fmt.Sprintf("Can`t find forum with id %v", forumId)
			result, _ := json.Marshal(models.Error{Message: &errorMessage})
			io.WriteString(w, string(result))
			return
		}
		postFull.Forum = &forum
	}

	postFull.Post = &post

	result, _ := json.Marshal(postFull)
	io.WriteString(w, string(result))
	return

}
