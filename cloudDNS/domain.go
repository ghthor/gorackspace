package cloudDNS

import (
	"encoding/json"
	"errors"
	"fmt"
	rackspace "github.com/ghthor/gorackspace"
	"io/ioutil"
	"net/http"
)

const Version = "1.0"

type (
	Domain struct {
		Name         string `json:"name"`
		Id           int    `json:"id"`
		Updated      string `json:"updated"`
		Created      string `json:"created"`
		TTL          int    `json:"ttl"`
		AccountId    int    `json:"accountId"`
		EmailAddress string `json:"emailAddress"`
		Comment      string `json:"comment"`
	}

	DomainListLink struct {
		Content string `json:"content"`
		Href    string `json:"href"`
		Rel     string `json:"rel"`
	}

	DomainListResponse struct {
		Domains      []Domain         `json:"domains"`
		Links        []DomainListLink `json:"links"`
		TotalEntries int              `json:"totalEntries"`
		rawJson      string
	}
)

func DomainList(session rackspace.AuthSession) ([]Domain, error) {
	// TODO: Inspect the Catalog to ensure this session has CloudDNS ability
	req, _ := http.NewRequest("GET", session.ServiceCatalog().CloudDNS[0].PublicURL+"/domains", nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", session.Id())

	resp, err := session.Client().Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	default:
		fallthrough
	case 401, 403, 400, 500, 503:
		return nil, errors.New(fmt.Sprintf("%s", responseBody))
	case 200, 203:
	}

	domainListResponse := &DomainListResponse{rawJson: string(responseBody)}

	// Parse Response Body
	err = json.Unmarshal(responseBody, domainListResponse)
	if err != nil {
		return nil, err
	}

	return domainListResponse.Domains, nil
}
