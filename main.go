package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			canciones, err := listarCanciones()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			renderTemplate(w, canciones)
			return
		}

		filePath := r.URL.Path[1:]

		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		http.ServeFile(w, r, filePath)
	})

	fmt.Println("Servidor web escuchando en http://localhost:8080/")
	http.ListenAndServe(":80", nil)
}

type cancion struct {
	Nombre string
	URL    string
}

func listarCanciones() ([]cancion, error) {
	var canciones []cancion
	err := filepath.Walk("mp3", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			canciones = append(canciones, cancion{Nombre: info.Name(), URL: "/mp3/" + info.Name()})
		}
		return nil
	})
	return canciones, err
}

func renderTemplate(w http.ResponseWriter, canciones []cancion) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, canciones)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
