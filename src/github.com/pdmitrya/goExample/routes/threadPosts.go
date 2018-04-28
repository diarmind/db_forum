package routes

import (
	"net/http"
	"github.com/pdmitrya/goExample/databases"
	"fmt"
	"encoding/json"
	"github.com/pdmitrya/goExample/models"
	"io"
	"strings"
	"database/sql"
	"log"
	"strconv"
)

func ThreadPosts(w http.ResponseWriter, r *http.Request, p map[string]string) {

	threadSlug := p["slug"]

	dbSession := databases.GetPostgresSession()

	w.Header().Set("Content-Type", "application/json")

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

	desc := r.URL.Query().Get("desc")

	since := r.URL.Query().Get("since")

	limit := r.URL.Query().Get("limit")

	sort := r.URL.Query().Get("sort")


	var queryBuilder strings.Builder

	if sort == "tree" {
		queryBuilder.WriteString(`WITH RECURSIVE PostParentTree AS (
									(
										SELECT
									U.nickname,
      								P.created,
									F.slug,
     		 						P.isEdited,
      								P.message,
      								P.thread_id,
      								P.parent AS parentId,
      								P.id,
									P.id::text AS path
									FROM Post P
									JOIN Forum_User U on P.user_id = U.id
      								JOIN Thread T on P.thread_id = T.id
      								JOIN Forum F on T.forum_id = F.id
									WHERE thread_id = $1 AND P.parent IS NULL `)



		queryBuilder.WriteString(`) UNION
								
									(
										SELECT
									U.nickname,
      								P.created,
									F.slug,
      								P.isEdited,
      								P.message,
      								P.thread_id,
      								P.parent,
      								P.id,
      								CONCAT_WS('.', PT.path, P.id::text)
									FROM Post P
									INNER JOIN PostParentTree PT ON P.parent = PT.id
									JOIN Forum_User U on P.user_id = U.id
      								JOIN Thread T on P.thread_id = T.id
      								JOIN Forum F on T.forum_id = F.id
									)
								)
									SELECT nickname, created, slug, isEdited, message, thread_id, parentId, PT.id FROM PostParentTree PT `)
		/*
		JOIN (SELECT id, MAX(CHAR_LENGTH(path)) AS totalLength FROM PostParentTree GROUP BY id) ONE ON PT.id = ONE.id AND ONE.totalLength = CHAR_LENGTH(path)*/


		if len(since) > 0 {
			if desc == "true" {
				queryBuilder.WriteString("WHERE path < (SELECT path FROM PostParentTree WHERE id = $2) ")
			} else {
				queryBuilder.WriteString("WHERE path > (SELECT path FROM PostParentTree WHERE id = $2) ")
			}
		} else {
			queryBuilder.WriteString("WHERE (TRUE OR ($2::TEXT != '')) ")
		}

		if desc == "true" {
			queryBuilder.WriteString(`ORDER BY substring(path from '^\d+') DESC, substring(path from '\..*$') DESC NULLS LAST, path DESC `)
		} else {
			queryBuilder.WriteString("ORDER BY path ASC, created ASC, id ASC ")
		}

		if len(limit) > 0 {
			queryBuilder.WriteString("LIMIT $3 ")
		}
	}

	if sort == "parent_tree" {
		queryBuilder.WriteString(`WITH RECURSIVE PostParentTree AS (
									(
										SELECT
									U.nickname,
      								P.created,
									F.slug,
     		 						P.isEdited,
      								P.message,
      								P.thread_id,
      								P.parent AS parentId,
      								P.id,
									P.id::text AS path
									FROM Post P
									JOIN Forum_User U on P.user_id = U.id
      								JOIN Thread T on P.thread_id = T.id
      								JOIN Forum F on T.forum_id = F.id
									WHERE thread_id = $1 AND P.parent IS NULL `)



		queryBuilder.WriteString(`) UNION
								
									(
										SELECT
									U.nickname,
      								P.created,
									F.slug,
      								P.isEdited,
      								P.message,
      								P.thread_id,
      								P.parent,
      								P.id,
      								CONCAT_WS('.', PT.path, P.id::text)
									FROM Post P
									INNER JOIN PostParentTree PT ON P.parent = PT.id
									JOIN Forum_User U on P.user_id = U.id
      								JOIN Thread T on P.thread_id = T.id
      								JOIN Forum F on T.forum_id = F.id
									)
								)
									SELECT pt1.nickname, pt1.created, pt1.slug, pt1.isEdited, pt1.message, pt1.thread_id, pt1.parentId, pt1.id FROM PostParentTree pt1
  										WHERE substring(pt1.path from '^\d+') IN (
  											SELECT path FROM PostParentTree `)

		if len(since) > 0 {
			if desc == "true" {
				queryBuilder.WriteString(`WHERE path < (SELECT substring(path from '^\d+') FROM PostParentTree WHERE id = $2) `)
			} else {
				queryBuilder.WriteString(`WHERE path > (SELECT substring(path from '^\d+') FROM PostParentTree WHERE id = $2) `)
			}
		} else {
			queryBuilder.WriteString("WHERE (TRUE OR ($2::TEXT != '')) ")
		}

		queryBuilder.WriteString("AND parentId IS NULL ")

		if desc == "true" {
			queryBuilder.WriteString(`ORDER BY path DESC, created ASC, id ASC `)
		} else {
			queryBuilder.WriteString("ORDER BY path ASC, created ASC, id ASC ")
		}

		if len(limit) > 0 {
			queryBuilder.WriteString("LIMIT $3 ")
		}

		queryBuilder.WriteString(") ")

		if desc == "true" {
			queryBuilder.WriteString(`ORDER BY substring(path from '^\d+') DESC, substring(path from '\..*$') ASC NULLS FIRST, path DESC `)
		} else {
			queryBuilder.WriteString("ORDER BY path ASC, created ASC, id ASC ")
		}

	}



	if sort == "flat" || (sort != "parent_tree" && sort != "tree") {
		queryBuilder.WriteString(`SELECT U.nickname, P.created, F.slug, P.isEdited, P.message, P.thread_id, P.parent, P.id
									FROM Post P
  										JOIN Forum_User U on P.user_id = U.id
  										JOIN Thread T on P.thread_id = T.id
  										JOIN Forum F on T.forum_id = F.id
  									WHERE thread_id = $1 `)

		if len(since) > 0 {
			if desc == "true" {
				queryBuilder.WriteString("AND P.id < $2 ")
			} else {
				queryBuilder.WriteString("AND P.id > $2 ")
			}
		} else {
			queryBuilder.WriteString( "AND (TRUE OR ($2::TEXT != '')) ")
		}

		if desc == "true" {
			queryBuilder.WriteString("ORDER BY P.created DESC, P.id DESC ")
		} else {
			queryBuilder.WriteString("ORDER BY P.created, P.id ")
		}

		if len(limit) > 0 {
			queryBuilder.WriteString("LIMIT $3 ")
		}

	}

	postArray := models.ArrayOfPosts{}
	var err error
	var rows *sql.Rows
	if len(limit) > 0 {
		rows, err = dbSession.Query(queryBuilder.String(), threadId, since, limit)
	} else {
		rows, err = dbSession.Query(queryBuilder.String(), threadId, since)
	}

	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var post models.Post
		if err := rows.Scan(&post.Author, &post.Created, &post.Forum, &post.IsEdited, &post.Message, &post.Thread, &post.Parent, &post.Id); err != nil {
			log.Fatal(err)
		}
		postArray = append(postArray, post)
	}
	result, _ := json.Marshal(postArray)
	io.WriteString(w, string(result))
	return
}
