package main

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/handlers"
)

func main() {

	header := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/" {
			canciones, err := listarCanciones()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Si se proporcionó un parámetro en la URL, reproducir esa canción primero
			if cancionParam := r.URL.Query().Get("cancion"); cancionParam != "" {
				cancionPath := "P:/evolmusic/mp3/" + cancionParam

				// Servir el archivo de la canción solicitada
				http.ServeFile(w, r, cancionPath)

				// Esperar a que la canción termine de reproducirse antes de continuar
				// la lista de reproducción
				time.Sleep(time.Second)
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

	// Ruta base para la carpeta public
	publicFS := http.FileServer(http.Dir("p:/evolmusic/mp3/"))
	http.Handle("/mp3/", http.StripPrefix("/mp3/", publicFS))

	// Ruta base para la carpeta assets
	assetsFS := http.FileServer(http.Dir("p:/evolmusic/js/"))
	http.Handle("/js/", http.StripPrefix("/js/", assetsFS))

	//SSL
	//C:\Certbot\live\mp3.ivacker.xyz\fullchain.pem

	//go http.ListenAndServe(":80", nil)

	err1 := http.ListenAndServeTLS(":443", "C:/Certbot/live/mp3.ivacker.xyz/fullchain.pem", "C:/Certbot/live/mp3.ivacker.xyz/privkey.pem", handlers.CORS(header, methods, origins)(http.DefaultServeMux)) // SSL (Utilizo certificado para SSL)
	fmt.Println("Servidor 443...")
	if err1 != nil {
		fmt.Print(err1)
	}

}

type cancion struct {
	Nombre string
	URL    string
}

func listarCanciones() ([]cancion, error) {
	rand.Seed(time.Now().UnixNano())

	var canciones []cancion
	err := filepath.Walk("P:/evolmusic/mp3/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			canciones = append(canciones, cancion{Nombre: info.Name(), URL: "/mp3/" + info.Name()})
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

func redirectToHTTPSHandler(domain string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://"+domain+r.RequestURI, http.StatusMovedPermanently)
	}
}

func mainHandler(w http.ResponseWriter, r *http.Request) {

	redirectToHTTPSHandler("https://mp3.ivacker.xyz")(w, r)

}
