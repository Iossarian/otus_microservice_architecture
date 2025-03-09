package delivery

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"order/internal/domain/saga/createorder"

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

func (c *Client) Reserve(ctx echo.Context, request createorder.Request) error {
	type Request struct {
		UUID string `json:"uuid"`
		Slot string `json:"slot"`
	}

	r := Request{
		UUID: request.UUID.String(),
		Slot: request.DeliverySlot,
	}

	jsonStr, err := json.Marshal(r)
	if err != nil {
		return errors.Wrap(err, "marshal request")
	}

	_, err = c.doRequest(ctx.Request().Context(), http.MethodPost, c.baseUrl+"/reserve", bytes.NewBuffer(jsonStr))
	if err != nil {
		return errors.Wrap(err, "executing http request")
	}

	return nil
}

func (c *Client) Slots(ctx echo.Context) ([]string, error) {
	resp, err := c.doRequest(ctx.Request().Context(), http.MethodGet, c.baseUrl+"/slots", nil)
	if err != nil {
		return nil, errors.Wrap(err, "executing http request")
	}

	var r struct {
		Slots []string `json:"slots"`
	}
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, errors.Wrap(err, "unmarshal response")
	}

	return r.Slots, nil
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

	fmt.Println(response)
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
