package types

type (
	WebSocketRequest struct {
		Event int    `json:"event"`
		Data  string `json:"data"`
	}
	WebSocketConfigResponse struct {
		ChunksSize int `json:"chunks_size"`
		Priotize   int `json:"priotize"`
	}
	WebSocketEvent struct {
		Event string `json:"event"`
		Data  any    `json:"data"`
	}
)

const (
	WS_MESSAGE int = iota
	PING
	PRIO_ECG
	PRIO_GYRO
)
