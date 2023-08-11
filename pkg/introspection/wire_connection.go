package introspection

import (
	"bufio"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

type WireConnectioon interface {
	IsOpen() bool
	Close()
}

type WireControlConnection[T any] struct {
	In       chan T
	Out      chan T
	shutdown chan bool
	socket   net.Conn
	wg       sync.WaitGroup
}

type WireEventConnectionType uint8

const (
	WireEventSource WireEventConnectionType = iota
	WireEventSink
)

type WireEventConnection[T any] struct {
	Channel  chan T
	shutdown chan bool
	socket   net.Conn
}

func NewWireControlConnection[T any](socket net.Conn) *WireControlConnection[T] {
	connection := &WireControlConnection[T]{
		In:       make(chan T),
		Out:      make(chan T),
		shutdown: make(chan bool),
		socket:   socket,
		wg:       sync.WaitGroup{},
	}

	connection.processControlMessages()
	return connection
}

func NewWireEventConnection[T any](socket net.Conn, connectionType WireEventConnectionType) *WireEventConnection[T] {
	connection := &WireEventConnection[T]{
		Channel:  make(chan T),
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

func (w *WireControlConnection[T]) processControlMessages() {
	go w.processIncomingControlMessages()

	go w.processOutgoingControlMessages()
}

func (w *WireControlConnection[T]) Close() error {
	w.shutdown <- true
	return w.closeSocket()
}

func (w *WireEventConnection[T]) Close() error {
	w.shutdown <- true
	return w.closeSocket()
}

func (w *WireControlConnection[T]) IsOpen() bool {
	return w.socket != nil
}

func (w *WireEventConnection[T]) IsOpen() bool {
	return w.socket != nil
}

func (w *WireControlConnection[T]) closeSocket() error {
	w.socket.Close()
	w.socket = nil
	return nil
}

func (w *WireEventConnection[T]) closeSocket() error {
	w.socket.Close()
	w.socket = nil
	return nil
}

func (w *WireControlConnection[T]) processIncomingControlMessages() {
	w.wg.Add(1)
	defer w.wg.Done()

	for {
		wireMessage, err := readWireMessage[T](w.socket)

		if err != nil {
			if _, ok := err.(*net.OpError); ok {
				log.Errorf("Control connection closed. Leaving incoming control message processing loop.")
				w.Close()
				return
			} else if err == io.EOF || err == io.ErrClosedPipe {
				log.Errorf("Control connection closed. Leaving incoming control message processing loop.")
				w.Close()
				return
			} else {
				log.Errorf("Error decoding control message: %v", reflect.TypeOf(err))
			}
		} else {
			switch wireMessage.MessageType {
			case MessageTypeControl:
				w.In <- wireMessage.Payload
			case MessageTypeShutdown:
				log.Infof("Client has shut down. Shutting down server.")
				w.Close()
				return
			default:
				log.Debugf("Unknown message type: %v", wireMessage.MessageType)
			}
		}
	}
}

func (w *WireControlConnection[T]) processOutgoingControlMessages() {
	w.wg.Add(1)
	defer w.wg.Done()
	var err error

	for {
		select {
		case <-w.shutdown:
			return
		case payload := <-w.Out:
			wireMessage := WireMessage[T]{MessageTypeControl, payload}

			if err = writeWireMessage(w.socket, &wireMessage); err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					log.Errorf("Control connection closed. Leaving outgoing control message processing loop.")
					if err = w.closeSocket(); err != nil {
						log.Errorf("Failed to close socket. Ignoring...")
					}
					return
				} else {
					log.Errorf("Error encoding control message: %v", err)
				}
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func (w *WireEventConnection[T]) processIncomingEvents() {
	for {
		wireMessage, err := readWireMessage[T](w.socket)
		if err != nil {
			if err == io.EOF || err == io.ErrClosedPipe {
				log.Errorf("Event connection closed. Leaving event message processing loop.")
				w.Close()
				return
			} else {
				log.Errorf("Error decoding event message: %v", err)
			}
		} else {
			switch wireMessage.MessageType {
			case MessageTypeEvent:
				w.Channel <- wireMessage.Payload
			case MessageTypeShutdown:
				w.Close()
				return
			default:
				log.Debugf("Message not supported by event connection: %v", wireMessage.MessageType)
			}
		}
	}
}

func (w *WireEventConnection[T]) processOutgoingEvents() {
	var err error

	for {
		select {
		case <-w.shutdown:
			return
		case event := <-w.Channel:
			wireMessage := WireMessage[T]{MessageTypeEvent, event}
			if err = writeWireMessage(w.socket, &wireMessage); err != nil {
				if err == io.EOF || err == io.ErrClosedPipe {
					log.Errorf("Event connection closed. Leaving event message processing loop.")
					if err = w.closeSocket(); err != nil {
						log.Errorf("Failed to close socket. Ignoring...")
					}
					return
				} else {
					log.Errorf("Error encoding event message: %v", err)
				}
			}
		default:
			time.Sleep(2 * time.Millisecond)
		}
	}
}

func writeWireMessage[T any](w io.Writer, message *WireMessage[T]) error {
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Content-Length: %d\r\n\r\n", len(data))
	_, err = w.Write(data)
	if err != nil {
		return err
	}
	log.Debugf("Sent message: %d %s", len(data), data)
	return nil
}

func readWireMessage[T any](r io.Reader) (*WireMessage[T], error) {
	var contentLength int
	var msg WireMessage[T]

	bufReader := bufio.NewReader(r)
	buf, err := bufReader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	fmt.Sscanf(buf, "Content-Length: %d", &contentLength)
	bufReader.ReadString('\n') // read the empty line

	content := make([]byte, contentLength)
	_, err = io.ReadFull(bufReader, content)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(content, &msg); err != nil {
		return nil, err
	}
	return &msg, nil
}
