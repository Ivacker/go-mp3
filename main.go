package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	// header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	// methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	// origins := handlers.AllowedOrigins([]string{"*"})

	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		dir := "P:/evolmusic/mp3"
		dirParam := "mp3"
		vars := mux.Vars(r)

		if dirParam, ok := vars["dir"]; ok {
			dir = "P:/evolmusic/" + dirParam
		}

		canciones, err := listarCanciones(dir, dirParam)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		renderTemplate(w, canciones)
		http.ServeFile(w, r, dir)

		
	}).Methods("GET")

	else {
		
	}

	http.Handle("/", r)

	http.ListenAndServe(":80", nil)
}

type cancion struct {
	Nombre string
	URL    string
}

func listarCanciones(sRuta string, sDir string) ([]cancion, error) {
	rand.Seed(time.Now().UnixNano())

	var canciones []cancion
	err := filepath.Walk(sRuta, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			canciones = append(canciones, cancion{Nombre: info.Name(), URL: "/" + sDir + "/" + info.Name()})
		}
		return nil
	})

	// Reordenar la lista de canciones de forma aleatoria
	for i := range canciones {
		j := rand.Intn(i + 1)
		canciones[i], canciones[j] = canciones[j], canciones[i]
	}

	return canciones, err
}

func renderTemplate(w http.ResponseWriter, canciones []cancion) {
	tmpl := template.Must(template.ParseFiles("index.html"))
	err := tmpl.Execute(w, canciones)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
