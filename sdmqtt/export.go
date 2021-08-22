package sdmqtt

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type (
	// Client
	OriginalClient      = mqtt.Client
	ClientOptions       = mqtt.ClientOptions
	ClientOptionsReader = mqtt.ClientOptionsReader

	// Message
	Message = mqtt.Message

	// Handlers
	MessageHandler           = mqtt.MessageHandler
	ConnectionAttemptHandler = mqtt.ConnectionAttemptHandler
	ConnectionLostHandler    = mqtt.ConnectionLostHandler
	OnConnectHandler         = mqtt.OnConnectHandler
	ReconnectHandler         = mqtt.ReconnectHandler

	// Token
	Token            = mqtt.Token
	ConnectToken     = mqtt.ConnectToken
	DisconnectToken  = mqtt.DisconnectToken
	PacketAndToken   = mqtt.PacketAndToken
	DummyToken       = mqtt.DummyToken
	PlaceHolderToken = mqtt.PlaceHolderToken
	SubscribeToken   = mqtt.SubscribeToken
	UnsubscribeToken = mqtt.UnsubscribeToken
	TokenErrorSetter = mqtt.TokenErrorSetter

	// Logger
	Logger     = mqtt.Logger
	NOOPLogger = mqtt.NOOPLogger

	// Store
	Store       = mqtt.Store
	MemoryStore = mqtt.MemoryStore

	// Other
	CredentialsProvider = mqtt.CredentialsProvider
	MId                 = mqtt.MId
)

var (
	NewClientOptions = mqtt.NewClientOptions
)
