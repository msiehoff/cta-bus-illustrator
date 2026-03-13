package postgres

import "gorm.io/gorm"

type routeSegmentModel struct {
	gorm.Model
	RouteID  uint    `gorm:"not null;index"`
	Sequence int     `gorm:"not null"`
	Lat      float64 `gorm:"not null"`
	Lng      float64 `gorm:"not null"`
}

func (routeSegmentModel) TableName() string {
	return "route_segments"
}
