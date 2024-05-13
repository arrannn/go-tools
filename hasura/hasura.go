package hasura

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	client      *http.Client
	gQLEndpoint string
	adminSecret string
}

type Payload struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

func NewClient(gqlEndpoint string, secret string) *Client {
	httpClient := &http.Client{}
	return &Client{
		client:      httpClient,
		gQLEndpoint: gqlEndpoint,
		adminSecret: secret,
	}
}

func (c *Client) Request(p Payload) ([]byte, error) {
	payloadBytes, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("error marshalling the full GraphQL payload: %w", err)
	}

	req, err := http.NewRequest("POST", c.gQLEndpoint, bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("x-hasura-admin-secret", c.adminSecret)

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error executing request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	return body, nil
}
