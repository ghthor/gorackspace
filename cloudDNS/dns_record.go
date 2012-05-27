package cloudDNS

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	rackspace "github.com/ghthor/gorackspace"
	"io/ioutil"
	"net/http"
)

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

func ListRecords(session rackspace.AuthSession, domain Domain) ([]Record, error) {
	// TODO: Inspect the Catalog to ensure this session has CloudDNS ability
	reqUrl := fmt.Sprintf("%s/domains/%d/records", session.ServiceCatalog().CloudDNS[0].PublicURL, domain.Id)
	req, _ := http.NewRequest("GET", reqUrl, nil)

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

	recordListResponse := &RecordListResponse{rawJson: string(responseBody)}

	// Parse Response Body
	err = json.Unmarshal(responseBody, recordListResponse)
	if err != nil {
		return nil, err
	}

	return recordListResponse.Records, nil
}

func AddRecords(session rackspace.AuthSession, domain Domain, records []Record) (*rackspace.JobStatus, error) {
	recordList := RecordList{records}
	recordListJson, err := json.Marshal(recordList)
	if err != nil {
		return nil, err
	}

	// TODO: Inspect the Catalog to ensure this session has CloudDNS ability
	reqUrl := fmt.Sprintf("%s/domains/%d/records", session.ServiceCatalog().CloudDNS[0].PublicURL, domain.Id)
	req, _ := http.NewRequest("POST", reqUrl, bytes.NewBuffer(recordListJson))

	req.Header.Set("Content-type", "application/json")
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

func AddRecord(session rackspace.AuthSession, domain Domain, newRecord Record) (*rackspace.JobStatus, error) {
	return AddRecords(session, domain, []Record{newRecord})
}
