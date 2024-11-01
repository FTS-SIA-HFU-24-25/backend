package types

type (
	WebSocketEvent struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}
	IoTEvent struct {
		Data []byte `json:"data"`
		Type int    `json:"type"`
	}
)
