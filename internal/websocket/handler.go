package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
)

// HandleWebSocket handles websocket requests from clients
func HandleWebSocket(hub *Hub) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Upgrade connection to websocket
		conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		// Get user ID from query or header (you can modify this based on your auth)
		userID := c.Query("user_id")
		if userID == "" {
			userID = "anonymous"
		}

		// Create new client
		client := &Client{
			Hub:    hub,
			Conn:   conn,
			Send:   make(chan []byte, 256),
			UserID: userID,
		}

		// Register client
		hub.Register <- client

		// Start goroutines for reading and writing
		go client.WritePump()
		go client.ReadPump()
	}
}
