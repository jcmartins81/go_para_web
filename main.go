package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Post struct {
	Id    int
	Title string
	Body  string
}

var db, err = sql.Open("mysql", "root:root@/go_course?charset=utf8")

func main() {

	stmt, err := db.Prepare("Insert into posts(title, body) values(?,?)")
	checkError(err)

	_, err = stmt.Exec("My first Post", "My first content")
	checkError(err)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		post := Post{
			Id:    1,
			Title: "Unamed Post",
			Body:  "No content",
		}

		if title := r.FormValue("title"); title != "" {
			post.Title = title
		}

		t := template.Must(template.ParseFiles("templates/index.html"))
		if err := t.ExecuteTemplate(w, "index.html", post); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println(http.ListenAndServe(":8080", nil))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}