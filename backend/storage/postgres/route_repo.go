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

func (r *RouteRepo) GetRoute(id string) (business.Route, error) {
	var model routeModel
	if err := r.db.Where("external_id = ?", id).First(&model).Error; err != nil {
		return business.Route{}, err
	}
	return business.Route{
		ExternalID: model.ExternalID,
		Name:       model.Name,
	}, nil
}

func (r *RouteRepo) CreateSegments(routeID string, segments []business.RouteSegment) error {
	tx := r.db.Begin()

	var model routeModel
	if err := tx.Where("external_id = ?", routeID).First(&model).Error; err != nil {
		return err
	}
	for _, segment := range segments {
		if err := tx.Create(&routeSegmentModel{
			RouteID:  model.ID,
			Sequence: segment.Sequence,
			Lat:      segment.Lat,
			Lng:      segment.Lng,
		}).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
