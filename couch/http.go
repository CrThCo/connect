package couch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var client *http.Client

func init() {
	client = &http.Client{}
}

// HTTPAPI struct
type HTTPAPI struct {
	URL      string
	Endpoint string
	Method   string
	Data     interface{}
}

// APIResponse struct
type APIResponse struct {
	*http.Response
}

func (api *HTTPAPI) getPath() string {
	return fmt.Sprintf("%s/%s", api.URL, api.Endpoint)
}

func (api *HTTPAPI) getData() io.Reader {
	data, err := json.Marshal(api.Data)
	if err != nil {
		return nil
	}
	return bytes.NewBuffer(data)
}

// Get method
func (api *HTTPAPI) Get() (*http.Response, error) {
	api.Method = "GET"
	return api.request()
}

// Put method
func (api *HTTPAPI) Put() (*http.Response, error) {
	api.Method = "PUT"
	return api.request()
}

func (api *HTTPAPI) request() (*http.Response, error) {
	request, err := http.NewRequest(api.Method, api.getPath(), api.getData())

	if err != nil {
		return nil, err
	}

	return client.Do(request)
}
