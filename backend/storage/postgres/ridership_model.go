package postgres

import (
	"time"

	"gorm.io/gorm"
)

type ridershipModel struct {
	gorm.Model
	RouteID        uint      `gorm:"not null;index"`
	MonthBeginning time.Time `gorm:"not null"`
	Type           string    `gorm:"not null"`
	AvgRides       float64   `gorm:"not null"`
}

func (ridershipModel) TableName() string {
	return "ridership"
}

// ridershipRow is used when joining with routes to retrieve the external_id.
type ridershipRow struct {
	ExternalID     string
	MonthBeginning time.Time
	Type           string
	AvgRides       float64
}
