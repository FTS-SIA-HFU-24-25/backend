package server

import (
	"net/http"
	"sia/backend/cache"
	"sia/backend/handler"
	"sia/backend/lib"
	"sia/backend/types"
	"time"

	"github.com/gorilla/websocket"
)

// Constants for performance tuning
const (
	WriteTimeout = 2 * time.Second // Matches handler.WebSocketHandler
)

// WebSocketServer manages WebSocket connections
type WebSocketServer struct {
	upgrader websocket.Upgrader
	handler  *handler.WebSocketHandler
}

// NewWebSocketServer initializes the WebSocket server
func NewWebSocketServer(config *cache.Config, allowedOrigins []string) *WebSocketServer {
	return &WebSocketServer{
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		handler: handler.NewWebSocketHandler(config),
	}
}

// ListenToWebSocket handles WebSocket connections
func (s *WebSocketServer) ListenToWebSocket(w http.ResponseWriter, r *http.Request, wsChan <-chan types.WebSocketEvent) {
	// Async logging for connection start
	go lib.Print(lib.WEBSOCKET_SERVICE, "Starting websocket server")

	c, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		go lib.Print(lib.WEBSOCKET_SERVICE, "Upgrade error: %v", err)
		return
	}
	defer c.Close()

	// Signal channel to stop the writer goroutine
	done := make(chan struct{})

	// Write messages to WebSocket from wsChan
	go func() {
		defer close(done) // Ensure done is closed when writer exits
		for {
			select {
			case data, ok := <-wsChan:
				if !ok {
					return // Channel closed
				}
				if err := c.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
					go lib.Print(lib.WEBSOCKET_SERVICE, "SetWriteDeadline error: %v", err)
					return
				}
				if err := c.WriteJSON(data); err != nil {
					go lib.Print(lib.WEBSOCKET_SERVICE, "WriteJSON error: %v", err)
					return
				}
			case <-done:
				return
			}
		}
	}()

	// Read messages from WebSocket
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				go lib.Print(lib.WEBSOCKET_SERVICE, "WebSocket closed normally: %v", err)
			} else {
				go lib.Print(lib.WEBSOCKET_SERVICE, "ReadMessage error: %v", err)
			}
			break
		}
		s.handler.HandleWebsocketEvent(c, mt, message)
	}
}
