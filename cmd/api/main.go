// Package cmd/api as the entry point of application from the api side (queue producer).
package main

import (
	"context"
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/app"
	"github.com/Konstantsiy/image-converter/pkg/logger"
)

func main() {
	err := app.Start()
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("failed to start app: %v", err))
	}
}
