package types

type (
	WebSocketRequest struct {
		Event int    `json:"event"`
		Data  string `json:"data"`
	}
	WebSocketConfigResponse struct {
		ChunksSize            int     `json:"chunks_size"`
		StartReceiveData      int     `json:"start_receive_data"`
		FilterType            int     `json:"filter_type"`
		MaxPass               float64 `json:"max_pass"`
		MinPass               float64 `json:"min_pass"`
		SpectrumUpdateRequest int     `json:"spectrum_update_request"`
	}
	WebSocketEvent struct {
		Event string      `json:"event"`
		Data  interface{} `json:"data"`
	}
)

const (
	WS_MESSAGE int = iota
	PING
	STOP_ECG
	START_ECG
	CONFIG_UPDATE
	SPECTRUM_UPDATE
)
