package utils

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/tcnksm/go-httpstat"
)

// HTTPClient struct struct to aggregate the httpClient and HTTPClientMethods
type HTTPClient struct {
	httpClient        *http.Client
	HTTPClientMethods HTTPClientMethods
}

// HTTPClientMethods contains the method definations
type HTTPClientMethods interface {
	Initialise()
	SendRequest(string, string, map[string]string, map[string]string, io.Reader) (*http.Response, error)
}

// Initialise method to initialise the httpClient from main
func (i *HTTPClient) Initialise() {
	i.httpClient = &http.Client{}
}

// SendRequest method to make an http call and capture the http metrics
func (i *HTTPClient) SendRequest(method string, url string, queryParams map[string]string, headers map[string]string, body io.Reader) (resp *http.Response, responseTimeInMilliSeconds int64, err error) {
	var resultMetrics httpstat.Result
	req, err := http.NewRequest(method, url, body)
	ctx := httpstat.WithHTTPStat(req.Context(), &resultMetrics)
	req = req.WithContext(ctx)
	if err != nil {
		log.Printf("Error in creating request : %s, %s", method, url)
		return
	}
	if len(headers) > 0 {
		for key, value := range headers {
			req.Header.Set(key, value)
		}
	}
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for key, value := range queryParams {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}
	resp, err = i.httpClient.Do(req)
	responseTimeInMilliSeconds = resultMetrics.Total(time.Now()).Milliseconds()
	return
}
