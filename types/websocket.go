package types

type (
	WebSocketRequest struct {
		Event int    `json:"event"`
		Data  string `json:"data"`
	}
)

const (
	WS_MESSAGE int = iota
	PING
)
