package main

import (
	"encoding/json"
	"fmt"

	l "github.com/redhatinsights/sources-superkey-worker/logger"
	"github.com/segmentio/kafka-go"
)

// Request - struct representing a request for a superkey
type Request struct {
	TenantID        string         `json:"tenant_id"`
	ApplicationType string         `json:"application_type"`
	SuperKeySteps   []SuperKeyStep `json:"superkey_steps"`
}

// SuperKeyStep - struct representing a step for SuperKey
type SuperKeyStep struct {
	Step          int                 `json:"step"`
	Name          string              `json:"name"`
	Payload       string              `json:"payload"`
	Substitutions []map[string]string `json:"substitutions"`
}

// Work - processes messages.
func Work(msg kafka.Message) {
	eventType := getEventType(msg.Headers)

	switch eventType {
	case "create_application":
		l.Log.Info("Processing `create_application` request")

		req, err := parseRequest(msg.Value)
		if err != nil {
			l.Log.Warnf("Error parsing request: %v", err)
			return
		}

		fmt.Printf("%v\n", req)
		l.Log.Infof("Finished processing `create_application` request for tenant %v type %v", req.TenantID, req.ApplicationType)
	default:
		l.Log.Warn("Unknown event_type")
	}
}

// parseRequest - parses a kafka message's value ([]byte) into a Request struct
// returns: *Request
func parseRequest(value []byte) (*Request, error) {
	request := Request{}
	err := json.Unmarshal(value, &request)
	if err != nil {
		return nil, err
	}

	return &request, nil
}

// getEventType - iterates through headers to find the `event_type` header
// from miq-messaging.
// returns: event_type value
func getEventType(headers []kafka.Header) string {
	for _, header := range headers {
		if header.Key == "event_type" {
			return string(header.Value)
		}
	}

	return ""
}