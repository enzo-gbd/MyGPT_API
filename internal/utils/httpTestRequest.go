// Package utils provides utility functions that support various operations across the application.
package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// StructToIOReader takes any interface{} value, marshals it into JSON, and returns an io.Reader for the JSON data.
// This function is useful for converting Go structures into a format that can be used for HTTP requests.
// Returns an error if the marshalling fails.
func StructToIOReader(v interface{}) (io.Reader, error) {
	jsonData, err := json.Marshal(v)
	if err != nil {
		return nil, err // return nil and the error if marshalling fails
	}
	reader := bytes.NewReader(jsonData)

	return reader, nil
}

// HttpTestRequest helps in creating and sending an HTTP request using a gin.Engine router for testing purposes.
// It takes a gin.Engine pointer for the router, HTTP method, URL, and an object to be used as the request body.
// The function converts the object to an io.Reader using StructToIOReader, constructs the request, and records the response using httptest.ResponseRecorder.
// Returns the recorded response and any error encountered during the request setup or execution.
func HttpTestRequest(router *gin.Engine, method string, url string, obj interface{}) (*httptest.ResponseRecorder, error) {
	body, err := StructToIOReader(obj) // Convert the request body to io.Reader
	if err != nil {
		return nil, err // return nil and the error if conversion fails
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, body) // assume no error on NewRequest for simplicity, handle errors as needed
	router.ServeHTTP(w, req)
	return w, nil
}
