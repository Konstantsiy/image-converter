// Package main as the entry point of application.
package main

import (
	"log"

	"github.com/Konstantsiy/image-converter/internal/app"
)

// Main function.
func main() {
	log.Fatal("failed to start app: ", app.Start())
}
