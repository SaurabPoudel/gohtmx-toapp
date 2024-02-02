package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Todo struct {
	Title string
	Done  bool
}

type PageData struct {
	Title   string
	Heading string
	Todos   []Todo
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	todos := []Todo{
		{"Task 1", false},
		{"Task 2", false},
		{"Task 3", false},
	}

	data := PageData{
		Title:   "Todo List",
		Heading: "Your Todos",
		Todos:   todos,
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		data := PageData{
			Title:   "Todo App",
			Heading: "Welcome to the Todo App!",
			Todos:   todos,
		}
		renderTemplate(w, "Base", data)
	})

	r.Post("/todos", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		title := r.FormValue("Title")
		todo := Todo{title, false}
		data.Todos = append(data.Todos, todo)
		renderTodoList(w, data.Todos)
	})
	http.ListenAndServe("localhost:3003", r)
}

func renderTemplate(w http.ResponseWriter, tmplName string, data interface{}) {
	tmpl, err := template.New(tmplName).ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, tmplName, data)
}

func renderTodoList(w http.ResponseWriter, todos []Todo) {
	tmpl, err := template.New("TodoList").Parse(`
  {{range .}}
  <li>{{.Title}}</li>
  {{end}}
  `)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, todos)
}

