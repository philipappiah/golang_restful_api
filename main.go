package main

import (
	"newsfeeder/database"
	"newsfeeder/httpd/handler"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Init()
	router := gin.Default()

	router.GET("api/v1/ping", handler.PingGet())
	router.GET("api/v1/newsfeeds", handler.GetNewsFeeds)
	router.GET("api/v1/newsfeed/:id", handler.GetNewsFeed)
	router.PATCH("api/v1/newsfeed/:id", handler.UpdateNewsFeed)
	router.POST("api/v1/newsfeeds", handler.CreateNewsFeed)
	router.DELETE("api/v1/newsfeed/:id", handler.DeleteNewsFeed)
	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
