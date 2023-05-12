package introspection

import (
	"encoding/gob"
	"io"
	"log"
	"net"
	"reflect"
	"sync"
	"time"
)

type WireConnectioon interface {
	IsOpen() bool
	Close()
}

type WireControlConnection struct {
	In       chan Payload
	Out      chan Payload
	shutdown chan bool
	socket   net.Conn
	wg       sync.WaitGroup
}

type WireEventConnectionType uint8

const (
	WireEventSource WireEventConnectionType = iota
	WireEventSink
)

type WireEventConnection struct {
	Channel  chan Payload
	shutdown chan bool
	socket   net.Conn
}

func NewWireControlConnection(socket net.Conn) *WireControlConnection {
	connection := &WireControlConnection{
		In:       make(chan Payload),
		Out:      make(chan Payload),
		shutdown: make(chan bool),
		socket:   socket,
		wg:       sync.WaitGroup{},
	}

	connection.processControlMessages()
	return connection
}

func NewWireEventConnection(socket net.Conn, connectionType WireEventConnectionType) *WireEventConnection {
	connection := &WireEventConnection{
		Channel:  make(chan Payload),
		shutdown: make(chan bool),
		socket:   socket,
	}

	if connectionType == WireEventSource {
		go connection.processOutgoingEvents()
	} else {
		go connection.processIncomingEvents()
	}
	return connection
}

func (w *WireControlConnection) processControlMessages() {
	w.wg.Add(1)
	go w.processIncomingControlMessages()

	w.wg.Add(1)
	go w.processOutgoingControlMessages()
}

func (w *WireControlConnection) Close() error {
	w.shutdown <- true
	return w.closeSocket()
}

func (w *WireEventConnection) Close() error {
	w.shutdown <- true
	return w.closeSocket()
}

func (w *WireControlConnection) IsOpen() bool {
	return w.socket != nil
}

func (w *WireEventConnection) IsOpen() bool {
	return w.socket != nil
}

func (w *WireControlConnection) closeSocket() error {
	w.socket.Close()
	w.socket = nil
	return nil
}

func (w *WireEventConnection) closeSocket() error {
	w.socket.Close()
	w.socket = nil
	return nil
}

func (w *WireControlConnection) processIncomingControlMessages() {
	defer w.wg.Done()
	decoder := gob.NewDecoder(w.socket)

	for {
		wireMessage := WireMessage{}
		if err := decoder.Decode(&wireMessage); err != nil {
			// check if its OpError
			if _, ok := err.(*net.OpError); ok {
				log.Printf("Control connection closed. Leaving incoming control message processing loop.")
				w.Close()
				return
			} else if err == io.EOF || err == io.ErrClosedPipe {
				log.Printf("Control connection closed. Leaving incoming control message processing loop.")
				w.Close()
				return
			} else {
				log.Printf("Error decoding control message: %v", reflect.TypeOf(err))
			}
		} else {
			switch wireMessage.MessageType {
			case MessageTypeControl:
				w.In <- wireMessage.Payload
			case MessageTypeShutdown:
				log.Printf("Client has shut down. Shutting down server.")
				w.Close()
				return
			default:
				log.Printf("Unknown message type: %v", wireMessage.MessageType)
			}
		}
	}
}

func (w *WireControlConnection) processOutgoingControlMessages() {
	defer w.wg.Done()
	encoder := gob.NewEncoder(w.socket)

	for {
		select {
		case <-w.shutdown:
			return
		case payload := <-w.Out:
			wireMessage := WireMessage{MessageTypeControl, payload}
			if err := encoder.Encode(wireMessage); err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					log.Printf("Control connection closed. Leaving outgoing control message processing loop.")
					if err := w.closeSocket(); err != nil {
						log.Printf("Failed to close socket. Ignoring...")
					}
					return
				} else {
					log.Printf("Error encoding control message: %v", err)
				}
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func (w *WireEventConnection) processIncomingEvents() {
	decoder := gob.NewDecoder(w.socket)

	for {
		wireMessage := WireMessage{}
		if err := decoder.Decode(&wireMessage); err != nil {
			if err == io.EOF || err == io.ErrClosedPipe {
				log.Printf("Event connection closed. Leaving event message processing loop.")
				w.Close()
				return
			} else {
				log.Printf("Error decoding event message: %v", err)
			}
		} else {
			switch wireMessage.MessageType {
			case MessageTypeEvent:
				w.Channel <- wireMessage.Payload
			case MessageTypeShutdown:
				w.Close()
				return
			default:
				log.Printf("Message not supported by event connection: %v", wireMessage.MessageType)
			}
		}
	}
}

func (w *WireEventConnection) processOutgoingEvents() {
	encoder := gob.NewEncoder(w.socket)

	for {
		select {
		case <-w.shutdown:
			return
		case event := <-w.Channel:
			wireMessage := WireMessage{MessageTypeEvent, event}
			if err := encoder.Encode(wireMessage); err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					log.Printf("Event connection closed. Leaving event message processing loop.")
					if err := w.closeSocket(); err != nil {
						log.Printf("Failed to close socket. Ignoring...")
					}
					return
				} else {
					log.Printf("Error encoding event message: %v", err)
				}
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}
}
