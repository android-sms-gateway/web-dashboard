package gateway

import (
	"context"
	"fmt"

	"github.com/android-sms-gateway/client-go/smsgateway"
)

type Client struct{ *smsgateway.Client }

func (c *Client) Send(
	ctx context.Context,
	message smsgateway.Message,
	options ...smsgateway.SendOption,
) (smsgateway.MessageState, error) {
	state, err := c.Client.Send(ctx, message, options...)
	return state, markUpstream(err)
}

func (c *Client) GetState(ctx context.Context, messageID string) (smsgateway.MessageState, error) {
	state, err := c.Client.GetState(ctx, messageID)
	return state, markUpstream(err)
}

func (c *Client) ListDevices(ctx context.Context) ([]smsgateway.Device, error) {
	devices, err := c.Client.ListDevices(ctx)
	return devices, markUpstream(err)
}

func (c *Client) DeleteDevice(ctx context.Context, id string) error {
	return markUpstream(c.Client.DeleteDevice(ctx, id))
}

func (c *Client) ListMessages(
	ctx context.Context,
	opts smsgateway.ListMessagesOptions,
) ([]smsgateway.MessageState, int, error) {
	messages, total, err := c.Client.ListMessages(ctx, opts)
	return messages, total, markUpstream(err)
}

func (c *Client) GetSettings(ctx context.Context) (smsgateway.DeviceSettings, error) {
	settings, err := c.Client.GetSettings(ctx)
	return settings, markUpstream(err)
}

func (c *Client) UpdateSettings(
	ctx context.Context,
	settings smsgateway.DeviceSettings,
) (smsgateway.DeviceSettings, error) {
	updated, err := c.Client.UpdateSettings(ctx, settings)
	return updated, markUpstream(err)
}

func (c *Client) ListWebhooks(ctx context.Context) ([]smsgateway.Webhook, error) {
	webhooks, err := c.Client.ListWebhooks(ctx)
	return webhooks, markUpstream(err)
}

func (c *Client) RegisterWebhook(ctx context.Context, webhook smsgateway.Webhook) (smsgateway.Webhook, error) {
	created, err := c.Client.RegisterWebhook(ctx, webhook)
	return created, markUpstream(err)
}

func (c *Client) DeleteWebhook(ctx context.Context, webhookID string) error {
	return markUpstream(c.Client.DeleteWebhook(ctx, webhookID))
}

func (c *Client) GenerateToken(ctx context.Context, req smsgateway.TokenRequest) (smsgateway.TokenResponse, error) {
	result, err := c.Client.GenerateToken(ctx, req)
	return result, markUpstream(err)
}

func (c *Client) RevokeToken(ctx context.Context, jti string) error {
	return markUpstream(c.Client.RevokeToken(ctx, jti))
}

func markUpstream(err error) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("%w: %w", ErrUpstream, err)
}
