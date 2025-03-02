package models

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) Save() error {
	filename := "cache/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(r *http.Request) (*Page, error) {
	vars := mux.Vars(r)
	filename := "cache/" + vars["title"] + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: vars["title"], Body: body}, nil
}
