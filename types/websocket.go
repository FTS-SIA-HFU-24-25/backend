package types

type (
	WebSocketRequest struct {
		Event int         `json:"event"`
		Data  interface{} `json:"data"`
	}
)
