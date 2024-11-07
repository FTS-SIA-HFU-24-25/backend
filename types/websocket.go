package types

type (
	WebSocketRequest struct {
		Event int    `json:"event"`
		Data  string `json:"data"`
	}
	WebSocketConfigResponse struct {
		ChunksSize       int     `json:"chunks_size"`
		StartReceiveData int     `json:"start_receive_data"`
		FilterType       int     `json:"filter_type"`
		MaxPass          float64 `json:"max_pass"`
		MinPass          float64 `json:"min_pass"`
	}
)

const (
	WS_MESSAGE int = iota
	PING
)
