package gorackspace

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	JobStatus struct {
		Status      string `json:"status"`
		JobId       string `json:"jobId"`
		CallbackURL string `json:"callbackUrl"`
	}
)

func JobStatusQuery(session AuthSession, jobStatus *JobStatus) (*JobStatus, error) {
	req, _ := http.NewRequest("GET", jobStatus.CallbackURL, nil)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Auth-Token", session.Id())

	resp, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, _ := ioutil.ReadAll(resp.Body)

	// TODO: Implement Richer error return responses to API errors
	switch resp.StatusCode {
	default:
		fallthrough
	case 400, 401, 404, 413, 500, 503:
		return nil, errors.New(fmt.Sprintf("Error Polling Job Status: %s", responseBody))
	case 200, 202:
	}

	err = json.Unmarshal(responseBody, jobStatus)
	if err != nil {
		return nil, err
	}

	return jobStatus, nil
}

func JobStatusMonitor(session AuthSession, jobStatus JobStatus, delay time.Duration) chan JobStatus {
	status := make(chan JobStatus)

	go func() {
		// A ticker to control polling delay
		ticker := time.NewTicker(delay)

		// Cleanup
		defer func() {
			ticker.Stop()
			close(status)
		}()

		// Start the Polling Loop
		for {
			// TODO: Identify and Gracefully handle all errors
			_, err := JobStatusQuery(session, &jobStatus)
			if err != nil {
				log.Println(err)
			}

			switch jobStatus.Status {

			case "ERROR", "COMPLETED":
				status <- jobStatus
				// Polling Complete
				return

			case "RUNNING":
				fallthrough
			default:
				select {
				case <-ticker.C:
				case status <- jobStatus:
					<-ticker.C
				}
			}
		}
	}()

	return status
}
