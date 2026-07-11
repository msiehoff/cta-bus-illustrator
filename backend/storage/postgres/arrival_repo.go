package postgres

import (
	"context"
	"time"

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

	order := "arrivals.timestamp DESC"
	if filter.SortAsc {
		order = "arrivals.timestamp ASC"
	}

	type row struct {
		StopID    string    `gorm:"column:stop_id"`
		RouteID   string    `gorm:"column:route_id"`
		Direction string    `gorm:"column:direction"`
		VehicleID string    `gorm:"column:vehicle_id"`
		Timestamp time.Time `gorm:"column:timestamp"`
		StopName  *string   `gorm:"column:stop_name"`
	}

	query := r.db.Table("arrivals").
		Select("arrivals.stop_id, arrivals.route_id, arrivals.direction, arrivals.vehicle_id, arrivals.timestamp, stops.name AS stop_name").
		Joins("LEFT JOIN stops ON stops.stop_id = arrivals.stop_id AND stops.route_id = arrivals.route_id AND stops.direction = arrivals.direction")
	query = r.applyArrivalFilter(query, filter)

	var rows []row
	if err := query.Order(order).Limit(limit).Offset(filter.Offset).Scan(&rows).Error; err != nil {
		return nil, err
	}

	arrivals := make([]business.Arrival, len(rows))
	for i, row := range rows {
		arrivals[i] = business.Arrival{
			StopID:    row.StopID,
			RouteID:   row.RouteID,
			Direction: row.Direction,
			VehicleID: row.VehicleID,
			Timestamp: row.Timestamp,
		}
		if row.StopName != nil {
			arrivals[i].StopName = *row.StopName
		}
	}
	return arrivals, nil
}

func (r *ArrivalRepo) CountArrivals(_ context.Context, filter app.ArrivalFilter) (int64, error) {
	var count int64
	query := r.db.Table("arrivals").
		Joins("LEFT JOIN stops ON stops.stop_id = arrivals.stop_id AND stops.route_id = arrivals.route_id AND stops.direction = arrivals.direction")
	query = r.applyArrivalFilter(query, filter)
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ArrivalRepo) ListArrivalsInRange(_ context.Context, start, end time.Time) ([]business.Arrival, error) {
	var models []arrivalModel
	err := r.db.Model(&arrivalModel{}).
		Where("timestamp >= ? AND timestamp < ?", start, end).
		Order("route_id ASC, direction ASC, stop_id ASC, timestamp ASC").
		Find(&models).Error
	if err != nil {
		return nil, err
	}

	arrivals := make([]business.Arrival, len(models))
	for i, model := range models {
		arrivals[i] = business.Arrival{
			StopID:    model.StopID,
			RouteID:   model.RouteID,
			Direction: model.Direction,
			VehicleID: model.VehicleID,
			Timestamp: model.Timestamp,
		}
	}
	return arrivals, nil
}

func (r *ArrivalRepo) applyArrivalFilter(query *gorm.DB, filter app.ArrivalFilter) *gorm.DB {
	if filter.RouteID != "" {
		query = query.Where("arrivals.route_id = ?", filter.RouteID)
	}
	if filter.Direction != "" {
		query = query.Where("arrivals.direction = ?", filter.Direction)
	}
	if filter.VehicleID != "" {
		query = query.Where("arrivals.vehicle_id = ?", filter.VehicleID)
	}
	if filter.Stop != "" {
		like := "%" + filter.Stop + "%"
		query = query.Where(
			"arrivals.stop_id = ? OR stops.name ILIKE ?",
			filter.Stop, like,
		)
	}
	return query
}
