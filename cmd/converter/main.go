// Package cmd/converter as the entry point of application from the converter side (queue consumer).
package main

import (
	"context"
	"fmt"

	"github.com/Konstantsiy/image-converter/internal/app"
	"github.com/Konstantsiy/image-converter/pkg/logger"
)

func main() {
	err := app.StartListener()
	if err != nil {
		logger.FromContext(context.Background()).
			Errorln(fmt.Errorf("failed to start converter: %v", err))
	}
}
