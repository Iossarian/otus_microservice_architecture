package warehouse

import (
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

func (c *Client) Stock(ctx echo.Context) (int, error) {
	resp, err := c.doRequest(ctx.Request().Context(), http.MethodGet, c.baseUrl+"/stock", nil)
	if err != nil {
		return 0, errors.Wrap(err, "executing http request")
	}

	var r struct {
		Stock int `json:"stock"`
	}
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return 0, errors.Wrap(err, "unmarshal response")
	}

	return r.Stock, nil
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

	if response.StatusCode != http.StatusOK {
		var ClientErr struct {
			Error string `json:"error"`
		}

		json.NewDecoder(response.Body).Decode(&ClientErr)
		err := errors.Errorf("bad http response status code: %d, error: %s", response.StatusCode, ClientErr.Error)

		return nil, err
	}

	b, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "reading response body")
	}

	return b, nil
}
