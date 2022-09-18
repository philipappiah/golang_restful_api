package handler

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"newsfeeder/database"
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
	Content   string             `bson:"content"`
}

func createFeed(newsfeed *NewsFeed) error {

	_, err := database.Collection.InsertOne(database.Ctx, newsfeed)

	return err
}

func FilterNewsFeeds(filter interface{}) ([]*NewsFeed, error) {
	// A slice of NewsFeeds for storing the decoded documents
	var NewsFeeds []*NewsFeed

	cur, err := database.Collection.Find(database.Ctx, filter)
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
	filter := bson.D{{}}
	var feeds, _ = FilterNewsFeeds(filter)

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
	c.JSON(http.StatusOK, map[string]string{
		"ID":      tobson.ID.Hex(),
		"Title":   tobson.Title,
		"Content": tobson.Content,
	})

}
