package newsfeed

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewsFeedPost() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"name":    "Philip",
			"message": "Welcome Home",
		})
	}

}
