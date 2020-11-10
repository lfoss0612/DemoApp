package context

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func readBody(r *http.Request) []byte {
	// Read the content
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(r.Body) // nolint
	}
	// Restore the io.ReadCloser to its original state
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// Use the content
	return bodyBytes
}

func getCurrentScheme(r *http.Request) string {
	var scheme = ""

	if urlString := r.Referer(); "" != urlString {
		if refererURL, parseErr := url.Parse(urlString); parseErr == nil && refererURL.Scheme != "" {
			scheme = refererURL.Scheme
		}
	}

	if r.URL.Scheme != "" && scheme == "" {
		scheme = r.URL.Scheme
	}

	if strings.Contains(r.Host, "local") {
		scheme = "http"
	}

	if scheme == "" {
		scheme = "https"
	}

	return scheme
}
