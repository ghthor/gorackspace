package cloudDNS

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	rackspace "github.com/ghthor/gorackspace"
	"github.com/ghthor/gorackspace/auth"
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

func DomainList(a *auth.Auth) ([]Domain, error) {
	req, _ := http.NewRequest("GET", a.ServiceCatalog.CloudDNS[0].PublicURL+"/domains", nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", a.AuthToken.Id)

	resp, err := rackspace.Client.Do(req)
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

type (
	// omitempty fields aren't needed when submiting a request to Add, Modify, Remove a Record
	Record struct {
		Name     string `json:"name"`
		Id       string `json:"id,omitempty"`
		Type     string `json:"type"`
		Data     string `json:"data"`
		Updated  string `json:"updated,omitempty"`
		Created  string `json:"created,omitempty"`
		TTL      int    `json:"ttl"`
		Comment  string `json:"comment,omitempty"`
		Priority int    `json:"priority,omitempty"`
	}

	RecordList struct {
		Records []Record `json:"records"`
	}

	RecordListResponse struct {
		Records      []Record `json:"records"`
		TotalEntries int      `json:"totalEntries"`
		rawJson      string
	}
)

func ListRecords(a *auth.Auth, domain Domain) ([]Record, error) {
	reqUrl := fmt.Sprintf("%s/domains/%d/records", a.ServiceCatalog.CloudDNS[0].PublicURL, domain.Id)
	req, _ := http.NewRequest("GET", reqUrl, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", a.AuthToken.Id)

	resp, err := rackspace.Client.Do(req)
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

	recordListResponse := &RecordListResponse{rawJson: string(responseBody)}

	// Parse Response Body
	err = json.Unmarshal(responseBody, recordListResponse)
	if err != nil {
		return nil, err
	}

	return recordListResponse.Records, nil
}

func AddRecord(a *auth.Auth, domain Domain, newRecord Record) (*rackspace.JobStatus, error) {
	recordList := RecordList{[]Record{newRecord}}
	recordListJson, err := json.Marshal(recordList)
	if err != nil {
		return nil, err
	}

	reqUrl := fmt.Sprintf("%s/domains/%d/records", a.ServiceCatalog.CloudDNS[0].PublicURL, domain.Id)
	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(recordListJson))

	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", a.AuthToken.Id)

	resp, err := rackspace.Client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)

	switch resp.StatusCode {
	default:
		fallthrough
	case 400, 401, 404, 413, 500, 503:
		return nil, errors.New(fmt.Sprintf("%s", responseBody))
	case 200, 202:
	}

	jobStatus := &rackspace.JobStatus{}
	err = json.Unmarshal(responseBody, jobStatus)
	if err != nil {
		return nil, err
	}

	return jobStatus, nil
}
