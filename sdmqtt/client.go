package sdmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gaorx/stardust3/sderr"
)

type Client struct {
	mqtt.Client
}

func NewClient(opts *ClientOptions) *Client {
	c := mqtt.NewClient(opts)
	return &Client{c}
}

func (c *Client) ConnectSync() error {
	token := c.Connect()
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.WithStack(err)
	}
	return nil
}

func (c *Client) SubscribeSync(topic string, qos byte, callback MessageHandler) error {
	token := c.Subscribe(topic, qos, callback)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.WithStack(err)
	}
	return nil
}

func (c *Client) UnsubscribeSync(topics ...string) error {
	token := c.Unsubscribe(topics...)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.WithStack(err)
	}
	return nil
}

func (c *Client) PublishSync(topic string, qos byte, retained bool, payload interface{}) error {
	token := c.Publish(topic, qos, retained, payload)
	token.Wait()
	if err := token.Error(); err != nil {
		return sderr.WithStack(err)
	}
	return nil
}
