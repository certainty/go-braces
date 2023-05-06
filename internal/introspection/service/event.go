package service

type Event interface{}

type ClientConnected struct {
	ClientID string
}
