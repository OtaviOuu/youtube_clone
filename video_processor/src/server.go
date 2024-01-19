package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	fluentffmpeg "github.com/modfy/fluent-ffmpeg"
)

type postRquestBody struct {
	InputVideoFilePath  string `json:"inputVideoFilePath"`
	OutputVideoFilePath string `json:"outputVideoFilePath"`
}

func main() {
	r := chi.NewRouter()

	//cmd := fluentffmpeg.NewCommand("")

	r.Post("/process", func(w http.ResponseWriter, r *http.Request) {
		// Decoding the body
		var requestBody postRquestBody
		err := json.NewDecoder(r.Body).Decode(&requestBody)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatalf("Cant decode the body: %s", err)
		}
		if (requestBody.InputVideoFilePath == "") || (requestBody.OutputVideoFilePath == "") {
			w.WriteHeader(http.StatusBadRequest)
			log.Fatal("Missing file path")
		}

		inputFile, outputFile := requestBody.InputVideoFilePath, requestBody.OutputVideoFilePath

		error := fluentffmpeg.NewCommand("").
			InputPath(inputFile).
			OutputPath(outputFile).
			OutputOptions("-vf", "scale=240:426").
			Run()
		if error != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatalf("Error, %s", error)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("video processing end."))

	})

	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	http.ListenAndServe(port, r)

}
