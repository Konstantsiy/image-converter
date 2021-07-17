// Package main as the entry point of application.
package main

import (
	"log"

	"github.com/Konstantsiy/image-converter/internal/app"
)

// Main function.
func main() {
	err := app.Start()
	if err != nil {
		log.Fatal("failed to start app: ", err)
	}
}
