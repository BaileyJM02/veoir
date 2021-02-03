package router

import (
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

// Start initiates the HTTP server for webhooks and requests the bot to start
func Start() {

	// Start the HTTP server ()
	server := &http.Server{
		Addr:         os.Getenv("address"),
		Handler:      GetRouter(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Inform us that we are starting the server
	logrus.Infof("Starting server on %v", server.Addr)

	// Start the server
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		logrus.Fatal(err)
	}
}
