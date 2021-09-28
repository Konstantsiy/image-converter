package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/pkg/logger"
)

// sendResponse marshals and writes response to the ResponseWriter.
func sendResponse(w http.ResponseWriter, resp interface{}, statusCode int) {
	respJSON, err := json.Marshal(resp)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("can't marshal response: %v", err))
		fmt.Fprint(w, resp)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set("Content-Type", "application/json")
	w.Write(respJSON)
}

// reportError logs and writes an error with the corresponding HTTP code to the ResponseWriter.
func reportError(w http.ResponseWriter, err error, statusCode int) {
	logger.Error(context.Background(), err)
	http.Error(w, err.Error(), statusCode)
}
