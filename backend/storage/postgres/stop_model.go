package postgres

import "time"

type stopModel struct {
	ID        uint `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	StopID    string  `gorm:"column:stop_id;not null"`
	RouteID   string  `gorm:"column:route_id;not null"`
	Direction string  `gorm:"column:direction;not null"`
	Name      string  `gorm:"not null"`
	Lat       float64 `gorm:"not null"`
	Lon       float64 `gorm:"not null"`
	Sequence  int     `gorm:"not null;default:0"`
}

func (stopModel) TableName() string {
	return "stops"
}
