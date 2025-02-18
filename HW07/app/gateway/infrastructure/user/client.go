package user

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
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

func (c *Client) Create(ctx echo.Context, request User) error {
	jsonStr, err := json.Marshal(request)
	if err != nil {
		return errors.Wrap(err, "marshal request")
	}

	_, err = c.doRequest(ctx.Request().Context(), http.MethodPost, c.baseUrl+"/users", bytes.NewBuffer(jsonStr))
	if err != nil {
		return errors.Wrap(err, "executing http request")
	}

	return nil
}

func (c *Client) Login(ctx echo.Context, request User) (string, error) {
	jsonStr, err := json.Marshal(request)
	if err != nil {
		return "", errors.Wrap(err, "marshal request")
	}

	b, err := c.doRequest(ctx.Request().Context(), http.MethodPost, c.baseUrl+"/login", bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", errors.Wrap(err, "executing http request")
	}

	type Response struct {
		Token string `json:"token"`
	}

	var response Response
	err = json.Unmarshal(b, &response)
	if err != nil {
		return "", errors.Wrap(err, "unmarshal response")
	}

	return response.Token, nil
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
