package routes

import (
	PageController "gowiki/app/controllers/pages"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

var Routes map[string]http.HandlerFunc

func defineRoutes() map[string]http.HandlerFunc {
	routes := map[string]http.HandlerFunc{}

	routes["/edit/{title}"] = makeHandler(PageController.EditHandler)
	routes["/view/{title}"] = makeHandler(PageController.ViewHandler)
	routes["/save/{title}"] = makeHandler(PageController.SaveHandler)

	return routes
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		vars := mux.Vars(r)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, vars["title"])
	}
}

func InitHandlers() {
	r := mux.NewRouter()
	Routes = defineRoutes()
	for name, handler := range Routes {
		r.HandleFunc(name, handler).Methods("GET")
	}
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("dist"))))

	log.Fatal(http.ListenAndServe(":8080", r))
}
