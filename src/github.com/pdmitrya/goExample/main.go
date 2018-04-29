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
    router.POST("/user/:nickname/create", routes.CreateUser)
    router.GET("/user/:nickname/profile", routes.ProfileUser)
    router.POST("/user/:nickname/profile", routes.UpdateUser)
    router.POST("/forum/create", routes.CreateForum)
    router.GET("/forum/:slug/details", routes.DetailsForum)
    router.POST("/forum/:slug/create", routes.CreateThread)
    router.GET("/forum/:slug/threads", routes.ForumThreads)
    router.GET("/thread/:slug/details", routes.DetailsThreadGet)
    router.POST("/thread/:slug/details", routes.DetailsThreadPost)
    router.POST("/thread/:slug/create", routes.CreatePosts)
    router.POST("/post/:id/details", routes.DetailsPostPost)
    router.GET("/post/:id/details", routes.DetailsPostGet)
    router.GET("/forum/:slug/users", routes.ForumUsers)
    router.GET("/thread/:slug/posts", routes.ThreadPosts)
    router.POST("/thread/:slug/vote", routes.VoteThread)
    router.POST("/service/clear", routes.ClearService)
    router.GET("/service/status", routes.StatusService)

    log.Fatal(http.ListenAndServe(":5001", router))
}
