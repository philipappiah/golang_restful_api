package handler

import (
	"net/http"
	"net/http/httptest"
	"newsfeeder/database"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetNews(t *testing.T) {
	database.Init()
	rPath := "/api/v1/newsfeeds"
	router := gin.Default()
	router.GET(rPath, GetNewsFeeds)
	req, _ := http.NewRequest("GET", rPath, nil)
	//req, _ := http.NewRequest("GET", rPath, strings.NewReader(`{"id": "1","name": "joe"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	//fmt.Println(w.Code)
	//fmt.Println(w.Body)
	//t.Logf("response: %s", w.Body.String())
}

func TestGetANews(t *testing.T) {
	database.Init()
	rPath := "/api/v1/newsfeed/:id"
	router := gin.Default()
	router.GET(rPath, GetNewsFeed)

	req, _ := http.NewRequest("GET", "/api/v1/newsfeed/6327732090e8b6f3bbc9a206", nil)
	//req, _ := http.NewRequest("GET", rPath, strings.NewReader(`{"id": "1","name": "joe"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	//fmt.Println(w.Code)
	//fmt.Println(w.Body)
	//t.Logf("response: %s", w.Body.String())
}

func TestUpdateNews(t *testing.T) {
	database.Init()
	rPath := "/api/v1/newsfeed/:id"
	router := gin.Default()
	router.PATCH(rPath, UpdateNewsFeed)

	req, _ := http.NewRequest("PATCH", "/api/v1/newsfeed/6327732090e8b6f3bbc9a206", strings.NewReader(`{"title": "cool title updated","content": "Content updated"}`))
	//req, _ := http.NewRequest("GET", rPath, strings.NewReader(`{"id": "1","name": "joe"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	//fmt.Println(w.Code)
	//fmt.Println(w.Body)
	//t.Logf("response: %s", w.Body.String())
}
