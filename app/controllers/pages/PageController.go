package pageController

import (
	helpers "gowiki/app"
	"gowiki/app/models"
	"net/http"
)

var (
	view = helpers.Parse("app/web/views/view.html")
	edit = helpers.Parse("app/web/views/edit.html")
)

func EditHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := models.LoadPage(r)
	if err != nil {
		p = &models.Page{Title: title}
	}
	if helpers.IsDevelopment() {
		helpers.Parse("app/web/views/edit.html").Execute(w, p)
	} else {
		edit.Execute(w, p)
	}
}

func SaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &models.Page{Title: title, Body: []byte(body)}
	err := p.Save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func ViewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := models.LoadPage(r)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	if helpers.IsDevelopment() {
		helpers.Parse("app/web/views/view.html").Execute(w, p)
	} else {
		view.Execute(w, p)
	}
}
