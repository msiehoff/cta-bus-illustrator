package postgres

import (
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"gorm.io/gorm"
)

type RouteRepo struct {
	db *gorm.DB
}

func NewRouteRepo(db *gorm.DB) *RouteRepo {
	return &RouteRepo{db: db}
}

func (r *RouteRepo) GetRoutes() ([]business.Route, error) {
	var models []routeModel
	if err := r.db.Preload("Segments", func(db *gorm.DB) *gorm.DB {
		return db.Order("sequence ASC")
	}).Find(&models).Error; err != nil {
		return nil, err
	}

	routes := make([]business.Route, len(models))
	for i, m := range models {
		segments := make([]business.RouteSegment, len(m.Segments))
		for j, s := range m.Segments {
			segments[j] = business.RouteSegment{
				Lat: s.Lat,
				Lng: s.Lng,
			}
		}
		routes[i] = business.Route{
			ExternalID: m.ExternalID,
			Name:       m.Name,
			Segments:   segments,
		}
	}
	return routes, nil
}
