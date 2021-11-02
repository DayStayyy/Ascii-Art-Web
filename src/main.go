/* ASCII ART WEB YNOV INFORMATIQUE 2020 */
/* Copyright INGREMEAU, CLAMADIEU-THARAUD, MICHEL 2020 */

package main

import (
	"fmt"
	"html/template"
	"net/http"

	"./ascii"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	var t *template.Template
	t = template.Must(template.ParseGlob("templates/index.html"))

	if r.Method == "GET" {
		if r.URL.Path != "/" {
			errorHandler(w, r, http.StatusNotFound)
			return
		}

		err := t.ExecuteTemplate(w, "index", nil)
		if err != nil {
			errorHandler(w, r, http.StatusInternalServerError)
			return
		}
	}

	if r.Method == "POST" {
		var mode bool
		var outputName string
		text := r.FormValue("text")
		font := r.FormValue("font")
		output := r.FormValue("export")

		if output == "true" {
			mode = true
			outputName = r.FormValue("output-name")
			if (len(outputName) < 5 && outputName == ".txt") || len(outputName) == 0 {
				outputName = "result.txt"
			} else {
				outputName = outputName + ".txt"
			}
		} else {
			mode = false
		}

		result, status := ascii.Art(text, font, outputName, mode)

		if status == 500 || status == 400 {
			errorHandler(w, r, status)
			return
		}

		if mode {
			w.Header().Set("Content-Disposition", "attachment; filename="+outputName)
			w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
			http.ServeFile(w, r, outputName)
		}
		t.ExecuteTemplate(w, "index", result)
	}
}

func errorHandler(w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		w.WriteHeader(http.StatusNotFound)
		temp := template.Must(template.ParseGlob("templates/404.html"))
		temp.ExecuteTemplate(w, "404", nil)
	}
	if status == 500 {
		w.WriteHeader(http.StatusInternalServerError)
		temp := template.Must(template.ParseGlob("templates/500.html"))
		temp.ExecuteTemplate(w, "500", nil)
	}
	if status == 400 {
		w.WriteHeader(http.StatusBadRequest)
		temp := template.Must(template.ParseGlob("templates/400.html"))
		temp.ExecuteTemplate(w, "400", nil)
	}
}

func serveCSS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "templates/style.css")
}

func main() {
	fmt.Println("Listening on port :8080")
	http.HandleFunc("/style.css", serveCSS)
	http.HandleFunc("/", indexHandler)
	http.ListenAndServe(":8080", nil)
}
