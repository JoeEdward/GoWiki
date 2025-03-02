package helpers

import (
	"embed"
	"html/template"
	"os"

	"github.com/joho/godotenv"
)

//go:embed web/views/*.html
var files embed.FS

func IsDevelopment() bool {
	godotenv.Load(".env")

	return os.Getenv("GO_DEVELOPMENT") != ""
}

func Parse(file string) *template.Template {
	if IsDevelopment() {
		return template.Must(
			template.New("layout.html").ParseFiles("app/web/views/layout.html", file),
		)
	} else {
		return template.Must(
			template.New("layout.html").ParseFS(files, "app/web/views/layout.html", file),
		)
	}
}
