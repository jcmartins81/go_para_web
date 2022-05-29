package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:my-secret-pw@(172.17.0.2:3306)/go_course?charset=utf8")

func main() {

	r := mux.NewRouter()

	r.PathPrefix("/static").Handler(http.StripPrefix("/static", http.FileServer(http.Dir("static/"))))

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("{id}/view", ViewHandler)

	fmt.Println(http.ListenAndServe(":8080", r))
}

func ListPosts() []Post {
	// stmt, err := db.Prepare("Insert into posts(title, body) values(?,?)")
	// checkError(err)

	// _, err = stmt.Exec("My first Post", "My first content")
	// checkError(err)

	rows, err := db.Query("Select * from posts")
	checkError(err)

	items := []Post{}

	for rows.Next() {
		post := Post{}
		rows.Scan(&post.Id, &post.Title, &post.Body)
		items = append(items, post)
	}

	db.Close()

	return items
}

func GetPostById(id string) *Post {
	row := db.QueryRow("select * from posts where id=?", id)
	post := Post()
	row.Scan(&post.Id, &post.Title, &post.Body)
	return post
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("templates/index.html"))
	if err := t.ExecuteTemplate(w, "index.html", ListPosts()); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	t := template.Must(template.ParseFiles("templates/view.html"))
	if err := t.ExecuteTemplate(w, "view.html", GetPostById(id)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
