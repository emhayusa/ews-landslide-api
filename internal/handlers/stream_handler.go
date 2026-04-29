package handlers

import (
	"big-devops-api/internal/services"
	"log"

	"github.com/gofiber/websocket/v2"
)

type StreamHandler struct {
	streamSvc *services.StreamService
}

func NewStreamHandler(svc *services.StreamService) *StreamHandler {
	return &StreamHandler{
		streamSvc: svc,
	}
}

func (h *StreamHandler) WSStream() func(*websocket.Conn) {
	return func(c *websocket.Conn) {
		h.streamSvc.RegisterClient(c)
		
		defer func() {
			h.streamSvc.UnregisterClient(c)
			c.Close()
		}()

		// Keep connection alive
		for {
			messageType, _, err := c.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("[WS] Read error: %v\n", err)
				}
				break
			}
			if messageType == websocket.CloseMessage {
				break
			}
		}
	}
}
