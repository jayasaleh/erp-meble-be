package utils

import (
	"encoding/json"

	"real-erp-mebel/be/internal/websocket"
)

// BroadcastUpdate mengirim update ke semua client yang terhubung via WebSocket
func BroadcastUpdate(hub *websocket.Hub, eventType string, data interface{}) error {
	message := map[string]interface{}{
		"type": eventType,
		"data": data,
	}

	jsonData, err := json.Marshal(message)
	if err != nil {
		return err
	}

	hub.BroadcastMessage(jsonData)
	return nil
}

// BroadcastError mengirim error message ke semua client
func BroadcastError(hub *websocket.Hub, message string) error {
	return BroadcastUpdate(hub, "error", map[string]string{
		"message": message,
	})
}

// BroadcastSuccess mengirim success message ke semua client
func BroadcastSuccess(hub *websocket.Hub, message string, data interface{}) error {
	return BroadcastUpdate(hub, "success", map[string]interface{}{
		"message": message,
		"data":    data,
	})
}
