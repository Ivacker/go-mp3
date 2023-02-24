package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
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

	http.HandleFunc("/generar", func(w http.ResponseWriter, r *http.Request) {

		canciones, err := listarCanciones()
		modificarHTML("index.html", canciones)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		filePath := r.URL.Path[1:]
		http.ServeFile(w, r, filePath)

	})

	// Ruta base para la carpeta public
	publicFS := http.FileServer(http.Dir("p:/evolmusic/mp3/"))
	http.Handle("/mp3/", http.StripPrefix("/mp3/", publicFS))

	// Ruta base para la carpeta assets
	assetsFS := http.FileServer(http.Dir("p:/evolmusic/js/"))
	http.Handle("/js/", http.StripPrefix("/js/", assetsFS))

	//SSL
	err1 := http.ListenAndServeTLS(":443", "C:/Certbot/live/mp3.ivacker.dev/fullchain.pem", "C:/Certbot/live/mp3.ivacker.dev/privkey.pem", handlers.CORS(header, methods, origins)(http.DefaultServeMux)) // SSL (Utilizo certificado para SSL)
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
			canciones = append(canciones, cancion{Nombre: info.Name(), URL: "/music/mp3/" + info.Name()})
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

func copiarTemplate() {

	// Abrir el archivo de origen en modo lectura
	origen, err := os.Open("template.html")
	if err != nil {
		log.Fatal(err)
	}
	defer origen.Close()

	// Crear el archivo de destino en modo escritura
	destino, err := os.Create("index.html")
	if err != nil {
		log.Fatal(err)
	}
	defer destino.Close()

	// Copiar los datos del archivo de origen al archivo de destino
	_, err = io.Copy(destino, origen)
	if err != nil {
		log.Fatal(err)
	}
}
func modificarHTML(sArchivo string, canciones []cancion) {

	copiarTemplate()

	// Aplicar la plantilla
	tmpl := template.Must(template.ParseFiles(sArchivo))
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, canciones); err != nil {
		panic(err)
	}

	// Escribir el resultado de la plantilla en el archivo
	if err := os.WriteFile("index.html", buf.Bytes(), 0644); err != nil {
		panic(err)
	}
}
