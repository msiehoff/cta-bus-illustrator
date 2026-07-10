package postgres

import "time"

type arrivalModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	StopID    string    `gorm:"column:stop_id;not null"`
	RouteID   string    `gorm:"column:route_id;not null"`
	Direction string    `gorm:"column:direction;not null"`
	VehicleID string    `gorm:"column:vehicle_id;not null"`
	Timestamp time.Time `gorm:"not null"`
}

func (arrivalModel) TableName() string {
	return "arrivals"
}
