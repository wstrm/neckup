package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
)

var views = template.Must(template.ParseFiles("views/index.html"))

func viewHandler(resWriter http.ResponseWriter, view string, data interface{}) {
	views.ExecuteTemplate(resWriter, view+".html", data)
}

func uploadHandler(resWriter http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case "GET":
		viewHandler(resWriter, "index", nil)

	case "POST":
		reader, err := req.MultipartReader()

		if err != nil {
			http.Error(resWriter, err.Error(), http.StatusInternalServerError)
			return
		}

		for {
			part, err := reader.NextPart()

			if err == io.EOF {
				break // Done
			}

			if part.FileName() == "" {
				continue // Empty file name, skip current iteration
			}

			dst, err := os.Create("./files/" + part.FileName())
			defer dst.Close()

			if err != nil {
				http.Error(resWriter, err.Error(), http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, part); err != nil {
				http.Error(resWriter, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		viewHandler(resWriter, "index", "Upload successful.")

	default:
		resWriter.WriteHeader(http.StatusMethodNotAllowed)
	}

}

func main() {

	http.HandleFunc("/", uploadHandler)

	http.ListenAndServe(":8080", nil)

}
