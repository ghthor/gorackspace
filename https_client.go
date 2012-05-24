package gorackspace

import (
	"crypto/tls"
	"net/http"
)

// A shared http client between all sub modules
var Client *http.Client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{},
	},
}
