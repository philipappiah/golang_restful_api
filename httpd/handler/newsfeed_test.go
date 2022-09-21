package handler

import (
	"net/http"
	"net/http/httptest"
	"newsfeeder/database"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ClearDatabaseForTest() {

	filter := bson.D{{}}
	_, err := database.Collection.DeleteMany(database.Ctx, filter)

	if err != nil {
		panic(err)
	}

}

func CreateTestData() (string, error) {

	ID := primitive.NewObjectID()
	doc := bson.D{{"_id", ID}, {"CreatedAt", time.Now()},
		{"UpdatedAt", time.Now()}, {"priority", 4}, {"title", "Test News Item"}, {"content", "The newsfeed APIs have been published"},
	}

	_, err := database.Collection.InsertOne(database.Ctx, doc)

	if err != nil {
		return "", err
	}

	return ID.Hex(), err
}

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
	ClearDatabaseForTest()
	_id, err := CreateTestData()
	if err != nil {
		panic(err)
	}

	rPath := "/api/v1/newsfeed/:id"
	router := gin.Default()
	router.GET(rPath, GetNewsFeed)

	req, _ := http.NewRequest("GET", "/api/v1/newsfeed/"+_id, nil)
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
	ClearDatabaseForTest()
	_id, err := CreateTestData()
	if err != nil {
		panic(err)
	}
	rPath := "/api/v1/newsfeed/:id"
	router := gin.Default()
	router.PATCH(rPath, UpdateNewsFeed)

	req, _ := http.NewRequest("PATCH", "/api/v1/newsfeed/"+_id, strings.NewReader(`{"title": "cool title updated","content": "Content updated"}`))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	// Assert we encoded correctly,
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	//fmt.Println(w.Code)
	//fmt.Println(w.Body)
	//t.Logf("response: %s", w.Body.String())
}
