package pkg

import (
	"fmt"
	"html/template"
	"net/http"
)

func ResponseError(w http.ResponseWriter, msg string, code int) {
	w.WriteHeader(code)

	tml, err := template.ParseFiles("./templates/html/error.page.html")
	if err != nil {
		http.Error(w, "internal server", http.StatusInternalServerError)
		return
	}
	tml.Execute(w, msg)
}

func GenerateTemplate(w http.ResponseWriter, data interface{}, filenames ...string) error {
	files := []string{}
	for _, file := range filenames {
		f := fmt.Sprintf("./templates/html/%s.html", file)
		files = append(files, f)
	}
	templates := template.Must(template.ParseFiles(files...))
	return templates.ExecuteTemplate(w, "layout", data)
}

func ParseTempalteFiles(filenames ...string) *template.Template {
	var files []string
	t := template.New(filenames[0])
	for _, file := range filenames {
		files = append(files, fmt.Sprintf("./templates/html/%s.html", file))
	}
	fmt.Printf("files: %v\n", files)
	t = template.Must(t.ParseFiles(files...))
	return t
}

var E error = fmt.Errorf("error create post")
