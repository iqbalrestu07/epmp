package websocket

// Hub manages websocket clients and broadcasts messages to them.
type Hub struct {
	// clients map[*Client]bool
}

func NewHub() *Hub {
	return &Hub{}
}

// TODO: Implement Run() loop to handle register/unregister channels
// TODO: Implement HTTP Handler Upgrade to WebSocket
// TODO: Implement PushNotificationToUser(userID string, payload interface{})
