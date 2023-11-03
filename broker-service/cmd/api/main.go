package main

import (
	"fmt"
	"log"
	"net/http"
)

const webPORT = "80"

type Config struct{}

func main() {
	app := Config{}

	log.Printf("Starting broker service on port %s", webPORT)

	// define http server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPORT),
		Handler: app.routes(),
	}

	// start the server
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
