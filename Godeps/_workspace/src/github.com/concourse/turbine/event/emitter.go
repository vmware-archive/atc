package event

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const VERSION = "1.0"

type Emitter interface {
	EmitEvent(Event)
	Close()
}

type NullEmitter struct{}

func (NullEmitter) EmitEvent(Event) {}
func (NullEmitter) Close()          {}

type websocketEmitter struct {
	logURL string
	drain  <-chan struct{}

	dialer *websocket.Dialer

	conn  *websocket.Conn
	connL *sync.Mutex

	writeL *sync.Mutex
}

func NewWebSocketEmitter(logURL string, drain <-chan struct{}) Emitter {
	return &websocketEmitter{
		logURL: logURL,
		drain:  drain,

		dialer: &websocket.Dialer{
			// allow detection of failed writes
			//
			// ideally this would be zero, but gorilla uses that to fill in its own
			// default of 4096 :(
			WriteBufferSize: 1,
		},

		connL: new(sync.Mutex),

		writeL: new(sync.Mutex),
	}
}

func (e *websocketEmitter) EmitEvent(event Event) {
	for {
		if !e.connect() {
			return
		}

		e.writeL.Lock()

		err := e.conn.WriteJSON(Message{
			Event: event,
		})

		e.writeL.Unlock()

		if err == nil {
			break
		}

		e.close()

		select {
		case <-time.After(time.Second):
		case <-e.drain:
			return
		}
	}
}

func (e *websocketEmitter) Close() {
	e.close()
}

func (e *websocketEmitter) connect() bool {
	e.connL.Lock()
	defer e.connL.Unlock()

	if e.conn != nil {
		return true
	}

	var err error

	for {
		e.conn, _, err = e.dialer.Dial(e.logURL, nil)
		if err == nil {
			err = e.conn.WriteJSON(VersionMessage{
				Version: VERSION,
			})
			if err == nil {
				return true
			}
		}

		select {
		case <-time.After(time.Second):
		case <-e.drain:
			return false
		}
	}
}

func (e *websocketEmitter) close() {
	e.connL.Lock()
	defer e.connL.Unlock()

	if e.conn != nil {
		conn := e.conn
		e.conn = nil
		conn.Close()
	}
}
