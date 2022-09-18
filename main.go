package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"newsfeeder/httpd/handler"

	"strings"

	"github.com/gin-gonic/gin"
)

type Article struct {
	Id      string `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

type Articles []Article

func getArticles(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Test Title", Desc: "Test Description", Content: "Hello World"},
	}
	fmt.Println("all articles Endpoint Hit")
	json.NewEncoder(w).Encode(articles)
}

func searchArticle(items Articles, id string) (r *Article) {
	for i := range items {
		if items[i].Id == id {
			return &items[i]
		}
	}
	return nil
}

func getArticle(w http.ResponseWriter, r *http.Request) {
	articles := Articles{
		Article{Title: "Test Title", Desc: "Test Description", Content: "Hello World", Id: "1"},
		Article{Title: "Test Title 2", Desc: "Test Description 2", Content: "Hello World 2", Id: "2"},
	}

	articleId := r.URL.Query().Get("id")
	fmt.Printf("server: query id: %s\n", articleId)

	item := searchArticle(articles, articleId)
	json.NewEncoder(w).Encode(item)

}

func addArticle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("server: %s /\n", r.Method)
	fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
	fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type"))
	fmt.Printf("server: headers:\n")
	for headerName, headerValue := range r.Header {
		fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("server: could not read request body: %s\n", err)
	}
	fmt.Printf("server: request body: %s\n", reqBody)

	fmt.Fprintf(w, `{"message": "hello!"}`)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage endpoint hit")
}

func handleRequests() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/articles", getArticles)
	http.HandleFunc("/getArticle", getArticle)
	http.HandleFunc("/addArticle", addArticle)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {

	r := gin.Default()
	r.GET("/ping", handler.PingGet())
	r.GET("/newsfeed", handler.NewsFeedGet())
	r.POST("/newsfeed", handler.NewsFeedPost())
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
