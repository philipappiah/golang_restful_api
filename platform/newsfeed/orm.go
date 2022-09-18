package newsfeed

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

type NewsItem struct {
	gorm.Model
	Title string
	Post  string
}

func dbConnect() (*gorm.DB, error) {
	db, err = gorm.Open("sqlite3", "test.db")
	return db, err
}

func intialMigration() {

	db, err = dbConnect()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	defer db.Close()
	db.AutoMigrate(&NewsItem{})

}

func getArticles(w http.ResponseWriter, r *http.Request) {
	db, err = dbConnect()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	defer db.Close()

	var items []NewsItem
	json.NewEncoder(w).Encode((items))
}

func addArticle(w http.ResponseWriter, r *http.Request) {
	db, err = dbConnect()
	if err != nil {
		fmt.Println(err.Error())
		panic("Failed to connect to database")
	}

	defer db.Close()

	vars := mux.Vars(r)
	title := vars["title"]
	post := vars["post"]

	db.Create(&NewsItem{Title: title, Post: post})

	fmt.Fprintf(w, "Data saved successfully")

}
