package events

import (
	"context"
	"sia/backend/db"
	"sia/backend/tools"
)

// HandleKeyEvent handles the key event based on its type (set, del, expired, etc.)
func HandleKeyEvent(eventType string, key string) {
	switch eventType {
	case "set":
		tools.Log("[EVENTS]", "Key set:", key)
	case "del":
		tools.Log("[EVENTS]", "Key deleted:", key)
	case "expired":
		tools.Log("[EVENTS]", "Key expired:", key)
	default:
		tools.Log("[EVENTS]", "Unknown event:", eventType, "for key:", key)
	}
}

// ListenForKeyEvents listens for key events like set, del, expired in Redis
func ListenForKeyEvents() {
	// Subscribe to the Redis keyevent notifications (replace 0 with your Redis DB index if needed)
	pubsub := db.RedisDB.Subscribe(context.Background(), "__keyevent@0__:set", "__keyevent@0__:del", "__keyevent@0__:expired")

	// Handle messages in a separate goroutine
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(context.Background())
			if err != nil {
				tools.Log("[EVENTS]", "Failed to receive message from Redis:", err)
				continue
			}

			// Extract the event type from the channel (e.g., __keyevent@0__:set -> "set")
			eventType := msg.Channel[len("__keyevent@0__:"):]
			tools.Log("[EVENTS]", "Received event:", eventType, "for key:", msg.Payload)

			// Handle the key event (msg.Payload contains the key name)
			HandleKeyEvent(eventType, msg.Payload)
		}
	}()
}

// HandleExpiredKey processes the expired key event
func HandleExpiredKey(key string) {
	// Add your logic to handle expired keys here
}

func updateBetStatusOnInit() int {
	return -1
}
