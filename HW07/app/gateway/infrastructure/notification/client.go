package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	httpClient httpClient
	baseUrl    string
}

func NewClient(
	httpClient httpClient,
	baseUrl string,
) *Client {
	return &Client{
		httpClient: httpClient,
		baseUrl:    baseUrl,
	}
}

func (c *Client) Message(ctx context.Context, r Request) ([]Message, error) {
	jsonStr, err := json.Marshal(r)
	if err != nil {
		return nil, errors.Wrap(err, "marshal request")
	}

	resp, err := c.doRequest(ctx, http.MethodGet, c.baseUrl+"/list", bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, errors.Wrap(err, "executing http request")
	}

	var response Response

	err = json.Unmarshal(resp, &response)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal response")
	}

	return response.Messages, nil
}

func (c *Client) doRequest(ctx context.Context, method string, url string, body io.Reader) ([]byte, error) {

	request, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "creating http request")
	}

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.Wrap(err, "executing http request")
	}

	defer func() {
		err := response.Body.Close()
		if err != nil {
			fmt.Printf("closing response body: %s", err.Error())
		}
	}()

	if response.StatusCode > http.StatusCreated {
		err := errors.Errorf("bad http response status code: %d", response.StatusCode)

		return nil, err
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	return b, nil
}
