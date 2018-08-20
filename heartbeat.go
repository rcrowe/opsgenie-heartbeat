package heartbeat

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/sethgrid/pester"
)

const (
	// EnvAPIKey is the environmental name that contains the opsgenie api key
	EnvAPIKey = "OPSGENIE_HEARTBEAT_KEY"

	defaultEndpoint = "https://api.opsgenie.com"
)

var (
	// ErrUnauthorised signafies an invalid api key
	ErrUnauthorised = errors.New("invalid api key")
	// ErrNonOkStatusCode a non 2XX status code was returned
	ErrNonOkStatusCode = errors.New("non 200 status code")
)

// PingRequest handles talking to the Opsgenie Heartbeat API
type PingRequest struct {
	// APIKey used to talk to the API
	APIKey string
	// Heartbeat is the name of the heartbeat to ping
	Heartbeat string

	Endpoint string
	Client   HTTPClient
}

// HTTPClient represents an HTTP client that can make a request
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// New creates a PingRequest with a default HTTP client and API key from env variable
func New(name string) PingRequest {
	client := pester.New()
	client.Timeout = 10 * time.Second
	client.Concurrency = 3
	client.MaxRetries = 5
	client.Backoff = pester.ExponentialJitterBackoff

	return PingRequest{
		APIKey:    os.Getenv(EnvAPIKey),
		Heartbeat: name,
		Endpoint:  defaultEndpoint,
		Client:    client,
	}
}

// Ping performs a HTTP request to the Opsgenie Heartbeat Ping endpoint
func (r PingRequest) Ping(ctx context.Context) error {
	url := fmt.Sprintf("%s/v2/heartbeats/%s/ping", strings.TrimRight(r.Endpoint, "/"), r.Heartbeat)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "GenieKey "+r.APIKey)
	req.WithContext(ctx)

	resp, err := r.Client.Do(req)
	if err != nil {
		return err
	}
	if resp.Body != nil {
		defer resp.Body.Close()
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return ErrUnauthorised
	}
	if c := resp.StatusCode; c < 200 || c > 299 {
		return ErrNonOkStatusCode
	}
	return nil
}
