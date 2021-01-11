package main

import (
	"log"
	"net/http"
	"os"

	"coolconf/web"
)

func main() {
	router := web.NewRouter()
	http.HandleFunc("/", router.Index())
	port := os.Getenv("COOLCONF_PORT")
	if port == "" {
		port = "10987"
	}
	log.Println("Listening on " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
