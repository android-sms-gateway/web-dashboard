package gateway

import (
	"net/http"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type Factory struct {
	baseURL string
}

func NewFactory(baseURL string) *Factory {
	return &Factory{baseURL: baseURL}
}

func (c *Factory) NewClient(login, password string) *smsgateway.Client {
	return smsgateway.NewClient(smsgateway.Config{
		Client:   http.DefaultClient,
		BaseURL:  c.baseURL,
		User:     login,
		Password: password,
		Token:    "",
	})
}
