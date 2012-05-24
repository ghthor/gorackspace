package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	rackspace "github.com/ghthor/gorackspace"
	"io/ioutil"
	"net/http"
)

const Version = "1.1"

type (
	AuthRequest struct {
		Credentials Credentials `json:"credentials"`
	}

	Credentials struct {
		Username string `json:"username"`
		Key      string `json:"key"`
	}
)

type (
	AuthToken struct {
		Id      string `json:"id"`
		Expires string `json:"expires"`
	}

	Service struct {
		Region      string `json:"region"`
		V1Default   bool   `json:"v1Default"`
		PublicURL   string `json:"publicURL"`
		InternalURL string `json:"internalURL"`
	}

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

	Auth struct {
		AuthToken      AuthToken      `json:"token"`
		ServiceCatalog ServiceCatalog `json:"serviceCatalog"`
	}

	AuthResponse struct {
		Auth    Auth `json:"auth"`
		rawJson []byte
	}
)

func (a AuthToken) String() string {
	str, _ := json.Marshal(a)
	return string(str)
}

type AuthFaultError struct {
	Code     int
	Response string
}

func (a AuthFaultError) Error() string {
	return fmt.Sprintf("%i: %s", a.Code, a.Response)
}

// TODO: Cache Auth tokens until they expire
// TODO: Enable access to the ServiceCatalog
func GetAuthToken(credentials Credentials) (AuthToken, error) {
	credsJson, err := json.Marshal(AuthRequest{credentials})
	if err != nil {
		return AuthToken{}, err
	}

	// Create Request
	req, _ := http.NewRequest("POST", "https://identity.api.rackspacecloud.com/v1.1/auth", bytes.NewBuffer(credsJson))
	req.Header.Set("Content-type", "application/json")

	// Request
	resp, err := rackspace.Client.Do(req)
	if err != nil {
		return AuthToken{}, err
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	default:
		fallthrough
	case 401, 403, 400, 500, 503:
		return AuthToken{}, AuthFaultError{resp.StatusCode, string(responseBody)}
	case 200, 203:
	}

	authResponse := &AuthResponse{rawJson: responseBody}

	// Parse Response Body
	err = json.Unmarshal(responseBody, authResponse)
	if err != nil {
		return AuthToken{}, err
	}

	return authResponse.Auth.AuthToken, nil
}
