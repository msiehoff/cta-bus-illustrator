package postgres

import (
	"context"

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
