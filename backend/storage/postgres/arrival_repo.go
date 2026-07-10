package postgres

import (
	"context"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"gorm.io/gorm"
)

type ArrivalRepo struct {
	db *gorm.DB
}

func NewArrivalRepo(db *gorm.DB) *ArrivalRepo {
	return &ArrivalRepo{db: db}
}

func (r *ArrivalRepo) SaveArrival(_ context.Context, arrival business.Arrival) error {
	model := arrivalModel{
		StopID:    arrival.StopID,
		RouteID:   arrival.RouteID,
		Direction: arrival.Direction,
		VehicleID: arrival.VehicleID,
		Timestamp: arrival.Timestamp,
	}
	return r.db.Create(&model).Error
}

func (r *ArrivalRepo) ListArrivals(_ context.Context, filter app.ArrivalFilter) ([]business.Arrival, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	query := r.applyArrivalFilter(r.db.Model(&arrivalModel{}), filter)
	var models []arrivalModel
	if err := query.Order("timestamp DESC").Limit(limit).Offset(filter.Offset).Find(&models).Error; err != nil {
		return nil, err
	}

	arrivals := make([]business.Arrival, len(models))
	for i, model := range models {
		arrivals[i] = toBusinessArrival(model)
	}
	return arrivals, nil
}

func (r *ArrivalRepo) CountArrivals(_ context.Context, filter app.ArrivalFilter) (int64, error) {
	var count int64
	query := r.applyArrivalFilter(r.db.Model(&arrivalModel{}), filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ArrivalRepo) applyArrivalFilter(query *gorm.DB, filter app.ArrivalFilter) *gorm.DB {
	if filter.RouteID != "" {
		query = query.Where("route_id = ?", filter.RouteID)
	}
	if filter.Direction != "" {
		query = query.Where("direction = ?", filter.Direction)
	}
	return query
}

func toBusinessArrival(model arrivalModel) business.Arrival {
	return business.Arrival{
		StopID:    model.StopID,
		RouteID:   model.RouteID,
		Direction: model.Direction,
		VehicleID: model.VehicleID,
		Timestamp: model.Timestamp,
	}
}
