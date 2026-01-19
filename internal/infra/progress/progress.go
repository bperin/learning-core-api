package progress

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ProgressUpdate represents a progress update for a long-running task
type ProgressUpdate struct {
	JobID      string    `json:"jobId"`
	TaskType   string    `json:"taskType"`
	Status     string    `json:"status"`
	Progress   int       `json:"progress"`
	Message    string    `json:"message"`
	Timestamp  time.Time `json:"timestamp"`
	Error      string    `json:"error,omitempty"`
	Data       any       `json:"data,omitempty"`
}

// ProgressTracker manages progress updates for long-running tasks
type ProgressTracker struct {
	mu        sync.RWMutex
	jobs      map[string]*JobProgress
	listeners map[string]map[*websocket.Conn]bool
}

// JobProgress tracks the progress of a single job
type JobProgress struct {
	JobID     string
	TaskType  string
	Status    string
	Progress  int
	Message   string
	StartedAt time.Time
	UpdatedAt time.Time
	Error     string
	Data      any
}

var tracker *ProgressTracker

func init() {
	tracker = &ProgressTracker{
		jobs:      make(map[string]*JobProgress),
		listeners: make(map[string]map[*websocket.Conn]bool),
	}
}

// GetTracker returns the global progress tracker
func GetTracker() *ProgressTracker {
	return tracker
}

// StartJob creates a new job progress tracker
func (pt *ProgressTracker) StartJob(jobID, taskType string) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	pt.jobs[jobID] = &JobProgress{
		JobID:     jobID,
		TaskType:  taskType,
		Status:    "started",
		Progress:  0,
		StartedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	log.Printf("[Progress] Job started: %s (type: %s)", jobID, taskType)
}

// UpdateProgress updates the progress of a job
func (pt *ProgressTracker) UpdateProgress(jobID, status, message string, progress int, data any) {
	pt.mu.Lock()
	job, exists := pt.jobs[jobID]
	if !exists {
		pt.mu.Unlock()
		return
	}

	job.Status = status
	job.Message = message
	job.Progress = progress
	job.UpdatedAt = time.Now()
	job.Data = data

	pt.mu.Unlock()

	// Broadcast update to all listeners
	pt.broadcastUpdate(jobID, status, message, progress, "", data)
}

// UpdateProgressWithError updates progress with an error
func (pt *ProgressTracker) UpdateProgressWithError(jobID, status, message, errMsg string, progress int) {
	pt.mu.Lock()
	job, exists := pt.jobs[jobID]
	if !exists {
		pt.mu.Unlock()
		return
	}

	job.Status = status
	job.Message = message
	job.Progress = progress
	job.Error = errMsg
	job.UpdatedAt = time.Now()

	pt.mu.Unlock()

	// Broadcast update to all listeners
	pt.broadcastUpdate(jobID, status, message, progress, errMsg, nil)
}

// CompleteJob marks a job as completed
func (pt *ProgressTracker) CompleteJob(jobID string) {
	pt.mu.Lock()
	job, exists := pt.jobs[jobID]
	if !exists {
		pt.mu.Unlock()
		return
	}

	job.Status = "completed"
	job.Progress = 100
	job.UpdatedAt = time.Now()

	pt.mu.Unlock()

	pt.broadcastUpdate(jobID, "completed", "Task completed successfully", 100, "", nil)

	// Clean up after a delay
	go func() {
		time.Sleep(5 * time.Minute)
		pt.mu.Lock()
		delete(pt.jobs, jobID)
		delete(pt.listeners, jobID)
		pt.mu.Unlock()
	}()
}

// FailJob marks a job as failed
func (pt *ProgressTracker) FailJob(jobID, errMsg string) {
	pt.mu.Lock()
	job, exists := pt.jobs[jobID]
	if !exists {
		pt.mu.Unlock()
		return
	}

	job.Status = "failed"
	job.Error = errMsg
	job.UpdatedAt = time.Now()

	pt.mu.Unlock()

	pt.broadcastUpdate(jobID, "failed", "Task failed", job.Progress, errMsg, nil)
}

// AddListener adds a WebSocket connection as a listener for a job
func (pt *ProgressTracker) AddListener(jobID string, conn *websocket.Conn) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if pt.listeners[jobID] == nil {
		pt.listeners[jobID] = make(map[*websocket.Conn]bool)
	}

	pt.listeners[jobID][conn] = true
	log.Printf("[Progress] Listener added for job: %s", jobID)

	// Send current progress immediately
	if job, exists := pt.jobs[jobID]; exists {
		update := ProgressUpdate{
			JobID:     job.JobID,
			TaskType:  job.TaskType,
			Status:    job.Status,
			Progress:  job.Progress,
			Message:   job.Message,
			Timestamp: job.UpdatedAt,
			Error:     job.Error,
			Data:      job.Data,
		}

		data, _ := json.Marshal(update)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

// RemoveListener removes a WebSocket connection from listening
func (pt *ProgressTracker) RemoveListener(jobID string, conn *websocket.Conn) {
	pt.mu.Lock()
	defer pt.mu.Unlock()

	if listeners, exists := pt.listeners[jobID]; exists {
		delete(listeners, conn)
		if len(listeners) == 0 {
			delete(pt.listeners, jobID)
		}
	}

	log.Printf("[Progress] Listener removed for job: %s", jobID)
}

// broadcastUpdate sends an update to all listeners of a job
func (pt *ProgressTracker) broadcastUpdate(jobID, status, message string, progress int, errMsg string, data any) {
	pt.mu.RLock()
	listeners := pt.listeners[jobID]
	job := pt.jobs[jobID]
	pt.mu.RUnlock()

	if job == nil {
		return
	}

	update := ProgressUpdate{
		JobID:     jobID,
		TaskType:  job.TaskType,
		Status:    status,
		Progress:  progress,
		Message:   message,
		Timestamp: time.Now(),
		Error:     errMsg,
		Data:      data,
	}

	data_json, err := json.Marshal(update)
	if err != nil {
		log.Printf("[Progress] Error marshaling update: %v", err)
		return
	}

	for conn := range listeners {
		if err := conn.WriteMessage(websocket.TextMessage, data_json); err != nil {
			log.Printf("[Progress] Error writing to WebSocket: %v", err)
			pt.RemoveListener(jobID, conn)
		}
	}
}

// GetJobProgress returns the current progress of a job
func (pt *ProgressTracker) GetJobProgress(jobID string) *JobProgress {
	pt.mu.RLock()
	defer pt.mu.RUnlock()

	return pt.jobs[jobID]
}

// HandleProgressWebSocket handles WebSocket connections for progress tracking
func HandleProgressWebSocket(w any, r any, jobID string) error {
	// This will be implemented in the HTTP handler
	return nil
}
