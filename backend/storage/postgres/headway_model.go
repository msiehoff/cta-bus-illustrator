package postgres

import "time"

type headwayModel struct {
	ID             uint `gorm:"primaryKey"`
	CreatedAt      time.Time
	StopID         string    `gorm:"column:stop_id;not null"`
	RouteID        string    `gorm:"column:route_id;not null"`
	Direction      string    `gorm:"column:direction;not null"`
	Timestamp      time.Time `gorm:"not null"`
	HeadwayMinutes float64   `gorm:"column:headway_minutes;not null"`
	FromVehicleID  string    `gorm:"column:from_vehicle_id"`
	ToVehicleID    string    `gorm:"column:to_vehicle_id"`
}

func (headwayModel) TableName() string {
	return "headways"
}

type headwayJobRunModel struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ServiceDate       time.Time `gorm:"column:service_date;type:date;not null"`
	Status            string    `gorm:"not null"`
	TriggeredBy       string    `gorm:"column:triggered_by;not null"`
	StartedAt         time.Time `gorm:"column:started_at;not null"`
	FinishedAt        *time.Time
	ArrivalsProcessed int    `gorm:"column:arrivals_processed;not null;default:0"`
	HeadwaysWritten   int    `gorm:"column:headways_written;not null;default:0"`
	SummariesWritten  int    `gorm:"column:summaries_written;not null;default:0"`
	ErrorMessage      string `gorm:"column:error_message"`
}

func (headwayJobRunModel) TableName() string {
	return "headway_job_runs"
}

type headwaySummaryModel struct {
	ID                uint `gorm:"primaryKey"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	ServiceDate       time.Time `gorm:"column:service_date;type:date;not null"`
	WindowStart       time.Time `gorm:"column:window_start;not null"`
	WindowEnd         time.Time `gorm:"column:window_end;not null"`
	Grain             string    `gorm:"not null"`
	Method            string    `gorm:"not null"`
	StopID            string    `gorm:"column:stop_id;not null;default:''"`
	RouteID           string    `gorm:"column:route_id;not null;default:''"`
	Direction         string    `gorm:"not null;default:''"`
	ObservationCount  int       `gorm:"column:observation_count;not null;default:0"`
	MeanMinutes       float64   `gorm:"column:mean_minutes;not null;default:0"`
	MedianMinutes     float64   `gorm:"column:median_minutes;not null;default:0"`
	StdDevMinutes     float64   `gorm:"column:stddev_minutes;not null;default:0"`
	CV                float64   `gorm:"column:cv;not null;default:0"`
	AvgWaitMinutes    float64   `gorm:"column:avg_wait_minutes;not null;default:0"`
}

func (headwaySummaryModel) TableName() string {
	return "headway_summaries"
}
