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

// A session created with an Auth(Cloud Identity) API call
type AuthSession interface {
	Client() *http.Client
	Id() string
	Expires() string
	ServiceCatalog() ServiceCatalog
	//Renew() 
}

type (
	// A Service that is accessible via the API
	Service struct {
		Region      string `json:"region"`
		V1Default   bool   `json:"v1Default"`
		PublicURL   string `json:"publicURL"`
		InternalURL string `json:"internalURL"`
	}

	// The information provided here is used by package modules in building API access URLs
	ServiceCatalog struct {
		CloudDNS           []Service `json:"cloudDNS"`
		CloudDatabases     []Service `json:"cloudDatabases"`
		CloudFiles         []Service `json:"cloudFiles"`
		CloudFilesCDN      []Service `json:"cloudFilesCDN"`
		CloudLoadBalancers []Service `json:"cloudLoadBalancers"`
		CloudMonitoring    []Service `json:"cloudMonitoring"`
		CloudServers       []Service `json:"cloudServers"`
		// TODO: See if this should be removed
		CloudServersOpenStack []Service `json:"cloudServersOpenStack"`
	}
)
