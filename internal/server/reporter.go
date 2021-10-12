package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Konstantsiy/image-converter/internal/service"

	"github.com/Konstantsiy/image-converter/pkg/logger"
)

const (
	ContentTypeKey   = "Content-Type"
	ContentTypeValue = "application/json"
)

// sendResponse marshals and writes response to the ResponseWriter.
func sendResponse(w http.ResponseWriter, resp interface{}, statusCode int) {
	respJSON, err := json.Marshal(resp)
	if err != nil {
		logger.FromContext(context.Background()).Errorln(fmt.Errorf("can't marshal response: %v", err))
		fmt.Fprint(w, resp)
		return
	}

	w.WriteHeader(statusCode)
	w.Header().Set(ContentTypeKey, ContentTypeValue)
	_, err = w.Write(respJSON)
	if err != nil {
		reportErrorWithCode(w, fmt.Errorf("can't write HTTP reply: %w", err), http.StatusInternalServerError)
	}
}

// reportErrorWithCode logs and writes an error with the corresponding HTTP code to the ResponseWriter.
func reportErrorWithCode(w http.ResponseWriter, err error, statusCode int) {
	logger.FromContext(context.Background()).Errorln(err)
	http.Error(w, err.Error(), statusCode)
}

// reportError reports custom service error.
func reportError(w http.ResponseWriter, err error) {
	subErr, ok := err.(*service.ServiceError)
	if !ok {
		reportErrorWithCode(w, err, http.StatusInternalServerError)
	}
	reportErrorWithCode(w, subErr.Err, subErr.StatusCode)
}
