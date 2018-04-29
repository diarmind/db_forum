package main

import (
    "log"
    "net/http"
    "github.com/pdmitrya/goExample/databases"
    "github.com/pdmitrya/goExample/routes"
	"github.com/dimfeld/httptreemux"
)

func main() {

    databases.ConnectDB()
    defer databases.CloseDB()

    router := httptreemux.New()
    group := router.NewGroup("/api")
    group.POST("/user/:nickname/create", routes.CreateUser)
    group.GET("/user/:nickname/profile", routes.ProfileUser)
    group.POST("/user/:nickname/profile", routes.UpdateUser)
    group.POST("/forum/create", routes.CreateForum)
    group.GET("/forum/:slug/details", routes.DetailsForum)
    group.POST("/forum/:slug/create", routes.CreateThread)
    group.GET("/forum/:slug/threads", routes.ForumThreads)
    group.GET("/thread/:slug/details", routes.DetailsThreadGet)
    group.POST("/thread/:slug/details", routes.DetailsThreadPost)
    group.POST("/thread/:slug/create", routes.CreatePosts)
    group.POST("/post/:id/details", routes.DetailsPostPost)
    group.GET("/post/:id/details", routes.DetailsPostGet)
    group.GET("/forum/:slug/users", routes.ForumUsers)
    group.GET("/thread/:slug/posts", routes.ThreadPosts)
    group.POST("/thread/:slug/vote", routes.VoteThread)
    group.POST("/service/clear", routes.ClearService)
    group.GET("/service/status", routes.StatusService)

    log.Fatal(http.ListenAndServe(":5000", router))
}
