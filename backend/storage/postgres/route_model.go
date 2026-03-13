package postgres

import "gorm.io/gorm"

type routeModel struct {
	gorm.Model
	ExternalID string              `gorm:"uniqueIndex;not null"`
	Name       string              `gorm:"not null"`
	Segments   []routeSegmentModel `gorm:"foreignKey:RouteID"`
}

func (routeModel) TableName() string {
	return "routes"
}
