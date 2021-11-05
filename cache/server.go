package cache

import (
	"fmt"
	"log"
	"net/http"
)

var (
	Hits int
	hit  chan bool
)

func StartServer() {
	Hits = 0
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Hits += 1
		fmt.Fprintf(w, "Made request successfully. Hits: %d", Hits)
		log.Printf("Made request successfully. Hits: %d", Hits)
	})

	log.Println("Starting server on port 2811")
	log.Fatal(http.ListenAndServe(":2811", nil))
}
