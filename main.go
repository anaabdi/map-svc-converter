package main

import (
	"log"
	"net/http"

	"github.com/anaabdi/map-svc-converter/geojson"
	"github.com/anaabdi/map-svc-converter/kml"
	"github.com/go-chi/chi"
)

func main() {
	router := chi.NewRouter()

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("I am coming from fmt-geojson-mapsvc"))
	})

	router.Post("/geojson/convert", geojson.Convert)
	router.Post("/kml/convert", kml.Convert)

	log.Fatal(http.ListenAndServe(":7777", router))
}
