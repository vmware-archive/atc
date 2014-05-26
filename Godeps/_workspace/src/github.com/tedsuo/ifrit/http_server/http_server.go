package http_server

import (
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/tedsuo/ifrit"
)

type httpServer struct {
	address string
	handler http.Handler
}

func New(address string, handler http.Handler) ifrit.Runner {
	return &httpServer{
		address: address,
		handler: handler,
	}
}

func (s *httpServer) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	listener, err := net.Listen("tcp", s.address)
	if err != nil {
		return err
	}

	serverErrChan := make(chan error, 1)
	wg := new(sync.WaitGroup)
	go func() {
		serverErrChan <- http.Serve(listener, waitHandler(s.handler, wg))
	}()

	close(ready)

	closeChan := make(chan struct{})
	for {
		select {
		case <-signals:
			close(closeChan)
			listener.Close()
		case err = <-serverErrChan:
			select {
			case <-closeChan:
				wg.Wait()
				return nil
			default:
				return err
			}
		}
	}

}

func waitHandler(handler http.Handler, wg *sync.WaitGroup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer wg.Done()
		wg.Add(1)
		handler.ServeHTTP(w, r)
	})
}
