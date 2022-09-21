package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"newsfeeder/database"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type NewsFeed struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Title     string             `bson:"title"`
	Priority  int                `bson:"priority"`
	Content   string             `bson:"content"`
}

func FilterNewsFeeds(filter interface{}, paginateMap map[string]int64, sortCriteria string) ([]*NewsFeed, error) {
	// A slice of NewsFeeds for storing the decoded documents
	var NewsFeeds []*NewsFeed

	skip := (paginateMap["page"] - 1) * paginateMap["limit"]
	opts := options.Find().SetLimit(paginateMap["limit"]).SetSkip(skip).SetSort(bson.D{{Key: "priority", Value: 1}})

	cur, err := database.Collection.Find(database.Ctx, filter, opts)
	if err != nil {
		return NewsFeeds, err
	}

	for cur.Next(database.Ctx) {
		var t NewsFeed
		err := cur.Decode(&t)
		if err != nil {
			return NewsFeeds, err
		}

		NewsFeeds = append(NewsFeeds, &t)
	}

	if err := cur.Err(); err != nil {
		return NewsFeeds, err
	}

	// once exhausted, close the cursor
	cur.Close(database.Ctx)

	if len(NewsFeeds) == 0 {
		return NewsFeeds, mongo.ErrNoDocuments
	}

	return NewsFeeds, nil
}

func GetNewsFeeds(c *gin.Context) {

	// passing bson.D{{}} matches all documents in the collection

	paginateMap := map[string]int64{
		"limit": 100, // default query limit = 100
		"page":  1,   // default page skip = 1
	}

	var sortCriteria = "-CreatedAt"

	filter := bson.D{{}}

	// filter := bson.D{{"priority", bson.D{{"$gt", 2}}}}

	//data, _ := bson.Marshal(c.Request.URL.Query())

	if c.Query("sort") != "" {
		sortCriteria = c.Query("sort")
	}

	fmt.Println(c.Query("sort"))

	if c.Query("limit") != "" {
		limit, err := strconv.ParseInt(c.Query("limit"), 10, 64)
		if err == nil {
			paginateMap["limit"] = limit
		}

	}

	if c.Query("page") != "" {
		page, err := strconv.ParseInt(c.Query("page"), 10, 64)
		if err == nil {
			paginateMap["page"] = page
		}
	}

	var feeds, _ = FilterNewsFeeds(filter, paginateMap, sortCriteria)

	c.JSON(http.StatusOK, map[string][]*NewsFeed{
		"items": feeds,
	})

}

func GetNewsFeed(c *gin.Context) {

	// passing bson.D{{}} matches all documents in the collection

	param_id := c.Param("id")
	//632769192dece7a17f4c21d6
	objID, _ := primitive.ObjectIDFromHex(param_id)

	var result NewsFeed
	var err = database.Collection.FindOne(database.Ctx, bson.M{"_id": objID}).Decode(&result)
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "No document Found",
		})
	} else {
		c.JSON(http.StatusOK, map[string]NewsFeed{
			"item": result,
		})
	}

}

func DeleteNewsFeed(c *gin.Context) {

	// passing bson.D{{}} matches all documents in the collection

	param_id := c.Param("id")
	//632769192dece7a17f4c21d6
	objID, _ := primitive.ObjectIDFromHex(param_id)
	var _, err = database.Collection.DeleteOne(database.Ctx, bson.M{"_id": objID})
	if err != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "No document Found",
		})
	} else {
		c.JSON(http.StatusOK, map[string]string{
			"message": "Document deleted successfully!",
		})
	}

}

func UpdateNewsFeed(c *gin.Context) {

	param_id := c.Param("id")
	//632769192dece7a17f4c21d6

	objID, _ := primitive.ObjectIDFromHex(param_id)
	filter := bson.M{"_id": objID}
	reqBody, err := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("server: request body: %s\n", reqBody)
	if err != nil {
		log.Fatal(err.Error())
	}

	var tobson interface{}

	_err := bson.UnmarshalExtJSON(reqBody, true, &tobson)
	if _err != nil {
		// handle error
		log.Fatal(_err.Error())
	}

	update := bson.M{
		"$set": tobson,
	}

	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	result := database.Collection.FindOneAndUpdate(database.Ctx, filter, update, &opt)
	if result.Err() != nil {
		log.Fatal(result.Err())
	}

	// 9) Decode the result
	var doc NewsFeed
	decodeErr := result.Decode(&doc)

	if decodeErr != nil {
		c.JSON(http.StatusNotFound, map[string]string{
			"message": "No document Found",
		})
	} else {
		c.JSON(http.StatusOK, map[string]NewsFeed{
			"item": doc,
		})
	}

}

func CreateNewsFeed(c *gin.Context) {
	// feed := &NewsFeed{ID: primitive.NewObjectID(), CreatedAt: time.Now(),
	// 	UpdatedAt: time.Now(), Title: "First News Item", Content: "The newsfeed APIs have been published"}
	reqBody, err_ := ioutil.ReadAll(c.Request.Body)

	if err_ != nil {
		log.Fatal(err_.Error())
	}

	var tobson NewsFeed

	_err := bson.UnmarshalExtJSON(reqBody, true, &tobson)
	if _err != nil {
		// handle error
		log.Fatal(_err.Error())
	}
	tobson.ID = primitive.NewObjectID()
	tobson.CreatedAt = time.Now()
	tobson.UpdatedAt = time.Now()

	_, err := database.Collection.InsertOne(database.Ctx, tobson)

	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"ID":       tobson.ID.Hex(),
		"Title":    tobson.Title,
		"Content":  tobson.Content,
		"Priority": tobson.Priority,
	})

}
