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

	for {
		select {
		case <-signals:
			listener.Close()
			wg.Wait()
			return nil
		case err = <-serverErrChan:
			return err
		}
	}

}

func waitHandler(handler http.Handler, wg *sync.WaitGroup) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wg.Add(1)
		defer wg.Done()
		handler.ServeHTTP(w, r)
	})
}
