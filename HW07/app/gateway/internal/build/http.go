package build

import (
	"net/http"
	"time"
)

func (b *Builder) httpClient() *http.Client {
	return &http.Client{
		Timeout: time.Second * 10,
	}
}
