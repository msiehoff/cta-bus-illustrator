package app

import (
	"sync"
	"time"
)

// PipelineStatus is a snapshot of pipeline health for observability (admin UI, logging).
type PipelineStatus struct {
	Running       bool
	Routes        []string
	PollInterval  time.Duration
	LastPollAt    time.Time
	LastPingCount int
	LastError     string
	StartedAt     time.Time
}

type pipelineStatusTracker struct {
	mu     sync.RWMutex
	status PipelineStatus
}

func newPipelineStatusTracker(cfg PipelineConfig) *pipelineStatusTracker {
	return &pipelineStatusTracker{
		status: PipelineStatus{
			Routes:       append([]string(nil), cfg.RouteIDs...),
			PollInterval: cfg.PollInterval,
			StartedAt:    time.Now(),
		},
	}
}

func (t *pipelineStatusTracker) setRunning(running bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status.Running = running
}

func (t *pipelineStatusTracker) recordPoll(pingCount int, pollErr error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.status.LastPollAt = time.Now()
	t.status.LastPingCount = pingCount
	if pollErr != nil {
		t.status.LastError = pollErr.Error()
	} else {
		t.status.LastError = ""
	}
}

func (t *pipelineStatusTracker) snapshot() PipelineStatus {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.status
}
