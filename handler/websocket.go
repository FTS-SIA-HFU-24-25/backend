package handler

import (
	"context"
	"encoding/json"
	"sia/backend/cache"
	"sia/backend/lib"
	"sia/backend/types"
	"time"

	"github.com/gorilla/websocket"
)

// Constants for performance tuning
const (
	WriteTimeout = 3 * time.Second // Timeout for WebSocket writes
)

// WebSocketIncomingEvent represents incoming WebSocket messages
type WebSocketIncomingEvent struct {
	Event int             `json:"event"`
	Data  json.RawMessage `json:"data"` // Use RawMessage to delay parsing
}

// WebSocketHandler manages WebSocket event processing
type WebSocketHandler struct {
	config *cache.Config
}

// NewWebSocketHandler initializes the handler
func NewWebSocketHandler(config *cache.Config) *WebSocketHandler {
	return &WebSocketHandler{config: config}
}

// HandleWebsocketEvent processes WebSocket messages
func (h *WebSocketHandler) HandleWebsocketEvent(c *websocket.Conn, mt int, message []byte) {
	var requestData WebSocketIncomingEvent
	if err := json.Unmarshal(message, &requestData); err != nil {
		go lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}

	// Async logging for ping
	go lib.Print(lib.WEBSOCKET_SERVICE, "Ping > %d", time.Now().Unix())

	var num int

	ctx := context.Background()
	switch requestData.Event {
	case types.PING:
		h.handlePing(c, ctx)
		json.Unmarshal(requestData.Data, &num)
		lib.Print(lib.WEBSOCKET_SERVICE, num)
		go h.prio(c, ctx, num)
		return
	default:
		go lib.Print(lib.WEBSOCKET_SERVICE, "Unknown event: %d", requestData.Event)
	}
}

// handlePing sends a pong response
func (h *WebSocketHandler) handlePing(c *websocket.Conn, ctx context.Context) {
	conf, err := h.config.GetConfig(ctx)
	if err != nil {
		go lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}

	// Set write deadline to prevent blocking
	if err := c.SetWriteDeadline(time.Now().Add(WriteTimeout)); err != nil {
		go lib.Print(lib.WEBSOCKET_SERVICE, err)
		return
	}

	if err := c.WriteJSON(types.WebSocketEvent{
		Event: "pong",
		Data:  conf,
	}); err != nil {
		go lib.Print(lib.WEBSOCKET_SERVICE, err)
	}
}

func (h *WebSocketHandler) prio(c *websocket.Conn, ctx context.Context, prioWhat int) {
	conf, err := h.config.GetConfig(ctx)
	if err != nil {
		return
	}

	if conf.Priotize == prioWhat {
		return
	}

	conf.Priotize = prioWhat

	err = h.config.UpdateConfig(ctx, *conf)
	if err != nil {
		return
	}
}
