package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/joho/godotenv"
)

//go:embed web/views/*
var files embed.FS

var (
	view = parse("web/views/view.html")
	edit = parse("web/views/edit.html")
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func parse(file string) *template.Template {
	if isDevelopment() {
		return template.Must(
			template.New("layout.html").ParseFiles("web/views/layout.html", file),
		)
	} else {
		return template.Must(
			template.New("layout.html").ParseFS(files, "web/views/layout.html", file),
		)
	}

}

type Page struct {
	Title string
	Body  []byte
}

func isDevelopment() bool {
	godotenv.Load(".env")

	return os.Getenv("GO_DEVELOPMENT") != ""
}

func (p *Page) save() error {
	filename := "cache/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "cache/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	if isDevelopment() {
		parse("web/views/view.html").Execute(w, p)
	} else {
		view.Execute(w, p)
	}
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	if isDevelopment() {
		parse("web/views/edit.html").Execute(w, p)
	} else {
		edit.Execute(w, p)
	}
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

func main() {
	// TODO: 2 - Refactor these methods out to a controller file
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("dist"))))
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Default().Print("Running local development server on: http://localhost:8080")
}
