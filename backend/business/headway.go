package business

import "time"

// Headway is an observed gap between two consecutive arrivals at a stop.
// Timestamp is the later of the two arrivals.
type Headway struct {
	StopID         string
	RouteID        string
	Direction      string
	Timestamp      time.Time
	HeadwayMinutes float64
	FromVehicleID  string
	ToVehicleID    string
}

// HeadwayJobStatus is the lifecycle state of a headway rollup run.
type HeadwayJobStatus string

const (
	HeadwayJobPending HeadwayJobStatus = "pending"
	HeadwayJobRunning HeadwayJobStatus = "running"
	HeadwayJobSuccess HeadwayJobStatus = "success"
	HeadwayJobFailed  HeadwayJobStatus = "failed"
)

// HeadwayJobTrigger identifies who/what started a run.
type HeadwayJobTrigger string

const (
	HeadwayTriggerCron  HeadwayJobTrigger = "cron"
	HeadwayTriggerAdmin HeadwayJobTrigger = "admin"
	HeadwayTriggerAPI   HeadwayJobTrigger = "api"
)

// HeadwayJobRun records metadata for one headway computation over a service date.
type HeadwayJobRun struct {
	ID                int64
	ServiceDate       time.Time // date at midnight UTC calendar date (use .Format("2006-01-02"))
	Status            HeadwayJobStatus
	TriggeredBy       HeadwayJobTrigger
	StartedAt         time.Time
	FinishedAt        *time.Time
	ArrivalsProcessed int
	HeadwaysWritten   int
	ErrorMessage      string
}
