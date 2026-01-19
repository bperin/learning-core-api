package infra

import (
	"log"
	"net/http"

	"learning-core-api/internal/infra/progress"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleProgressWebSocket handles WebSocket connections for progress tracking
func HandleProgressWebSocket(w http.ResponseWriter, r *http.Request) {
	jobID := r.URL.Query().Get("jobId")
	if jobID == "" {
		http.Error(w, "Missing jobId parameter", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WebSocket] Error upgrading connection: %v", err)
		return
	}
	defer conn.Close()

	tracker := progress.GetTracker()
	tracker.AddListener(jobID, conn)
	defer tracker.RemoveListener(jobID, conn)

	log.Printf("[WebSocket] Client connected for job: %s", jobID)

	// Keep connection alive and handle incoming messages
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WebSocket] Error: %v", err)
			}
			break
		}
	}

	log.Printf("[WebSocket] Client disconnected for job: %s", jobID)
}
