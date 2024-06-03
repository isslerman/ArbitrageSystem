package cex

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

type ContentType int

const (
	JSON ContentType = iota
	FORM
)

type HTTPService struct {
	// Method      string
	ContentType string
	jsonData    any
	AuthKey     string
	Headers     map[string]string
	Url         string
}

// RequestOptions defines options for the HTTP request.
type RequestOptions struct {
	AuthHeader  string            // Authentication header value
	FormAuthKey string            // Form authentication key
	ContentType ContentType       // ContetType
	JSONData    interface{}       // JSON data to be sent
	FormData    map[string]string // Form data to be sent
}

func NewHTTPService() *HTTPService {
	return &HTTPService{}
}

func (h *HTTPService) NewRequest(url, method string, options RequestOptions) (*http.Response, error) {
	var body io.Reader
	var contentType string

	switch options.ContentType {
	case FORM:
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		// Add form fields
		for key, value := range options.FormData {
			if err := w.WriteField(key, value); err != nil {
				return nil, err
			}
		}

		// Add form authentication key if provided
		if options.FormAuthKey != "" {
			if err := w.WriteField("auth_token", options.FormAuthKey); err != nil {
				return nil, err
			}
		}

		if err := w.Close(); err != nil {
			return nil, err
		}

		body = &b
		contentType = w.FormDataContentType()
	case JSON:
		// JSON request
		jsonData, err := json.Marshal(options.JSONData)
		if err != nil {
			return nil, err
		}

		body = bytes.NewBuffer(jsonData)
		contentType = "application/json"
	}

	// Create the request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// Set the content type
	req.Header.Set("Content-Type", contentType)

	// Add authentication header if provided
	if options.AuthHeader != "" {
		req.Header.Set("Authorization", options.AuthHeader)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// the new http request, receive the http method, url of the request, the key for auth
// func (h *HTTPService) NewRequest(url, method string,m string, ct ContentType, data any) ([]byte, error) {

// 	body := new(bytes.Buffer)

// 	switch ct {
// 	case JSON:
// 		// Convert the User object to JSON.
// 		jsonData, err := json.Marshal(data)
// 		if err != nil {
// 			fmt.Println("Error marshalling data:", err)
// 			return nil, err
// 		}

// 		body = bytes.NewBuffer(jsonData)
// 	case FORM:
// 		// Create form data
// 		formData := url.Values{}
// 		formData.Set("key1", "value1")
// 		formData.Set("key2", "value2")

// 		// Encode the form data
// 		encodedFormData := formData.Encode()
// 	}

// 	req, err := http.NewRequest(m, h.Url, body)
// 	if err != nil {
// 		panic(err)
// 	}

// 	ctx := context.Background()
// 	req = req.WithContext(ctx)

// 	switch contentType {
// 	case "form":
// 		req.Header.Set("Content-Type", "multipart/form-data")
// 	case "json":
// 		req.Header.Set("Content-Type", "application/json")
// 	default:
// 		req.Header.Set("Content-Type", "application/json")
// 	}
// 	// Set the content type to application/json.
// 	if key != "" {
// 		req.Header.Set("Authorization", key)
// 	}

// 	// for each header passed, add the header value to the request
// 	for k, v := range headers {
// 		req.Header.Set(k, v)
// 	}

// 	// Execute the request
// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer resp.Body.Close()

// 	// Read the response body
// 	respBody, err := io.ReadAll(resp.Body)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if resp.StatusCode != http.StatusOK {
// 		return nil, fmt.Errorf("error calling %s:\nstatus: %s\nresponseData: %s", url, resp.Status, respBody)
// 	}

// 	// Return the response body as a string
// 	return respBody, nil
// }
