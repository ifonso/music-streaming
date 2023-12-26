package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const defDir = "./storage/compressed/"
const music = "timmaia.m3u8"

func getContentType(f *os.File) string {
	buffer := make([]byte, 512)

	f.Read(buffer)
	f.Seek(0, 0)

	return http.DetectContentType(buffer)
}

func serveFile(w http.ResponseWriter, r *http.Request, filename string) {
	file, err := os.Open(defDir + filename)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if stat.IsDir() {
		http.NotFound(w, r)
		return
	}

	// Defining Content-Type
	contentType := getContentType(file)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", fmt.Sprint(stat.Size()))

	// Reading content and assigin to ResponseWriter
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println("Error serving file:", err)
	}
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		filename := r.URL.Path[len("/"):]

		if len(filename) == 0 {
			filename = music
		}

		serveFile(w, r, filename)
	})

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal("Error while starting server: ", err.Error())
	}
}
