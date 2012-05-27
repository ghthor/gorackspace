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

	Auth struct {
		AuthToken      AuthToken                `json:"token"`
		ServiceCatalog rackspace.ServiceCatalog `json:"serviceCatalog"`
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

// An implementation of the gorackspace.AuthSession interface
type AuthSession struct {
	client         *http.Client
	authToken      AuthToken
	serviceCatalog rackspace.ServiceCatalog
}

func (a *AuthSession) String() string {
	return fmt.Sprintf("Session-Id: %s\tCatalog: %s", a.authToken.Id, a.serviceCatalog)
}

func (a *AuthSession) Client() *http.Client {
	return a.client
}

func (a *AuthSession) Id() string {
	return a.authToken.Id
}

func (a *AuthSession) Expires() string {
	return a.authToken.Expires
}

func (a *AuthSession) ServiceCatalog() rackspace.ServiceCatalog {
	return a.serviceCatalog
}

// TODO: Cache Auth tokens until they expire
func Authenticate(credentials Credentials) (rackspace.AuthSession, error) {
	credsJson, err := json.Marshal(AuthRequest{credentials})
	if err != nil {
		return nil, err
	}

	// Create Request
	req, _ := http.NewRequest("POST", "https://identity.api.rackspacecloud.com/v1.1/auth", bytes.NewBuffer(credsJson))
	req.Header.Set("Content-type", "application/json")

	// Request
	resp, err := rackspace.Client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	default:
		fallthrough
	case 401, 403, 400, 500, 503:
		return nil, AuthFaultError{resp.StatusCode, string(responseBody)}
	case 200, 203:
	}

	authResponse := &AuthResponse{rawJson: responseBody}

	// Parse Response Body
	err = json.Unmarshal(responseBody, authResponse)
	if err != nil {
		return nil, err
	}

	authSession := &AuthSession{
		client:         rackspace.Client,
		authToken:      authResponse.Auth.AuthToken,
		serviceCatalog: authResponse.Auth.ServiceCatalog,
	}

	return authSession, nil
}
