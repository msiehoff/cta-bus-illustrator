package postgres

import (
	"context"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type StopRepo struct {
	db *gorm.DB
}

func NewStopRepo(db *gorm.DB) *StopRepo {
	return &StopRepo{db: db}
}

func (r *StopRepo) UpsertStops(_ context.Context, stops []business.Stop) error {
	if len(stops) == 0 {
		return nil
	}

	models := make([]stopModel, len(stops))
	for i, s := range stops {
		models[i] = stopModel{
			StopID:    s.StopID,
			RouteID:   s.RouteID,
			Direction: s.Direction,
			Name:      s.Name,
			Lat:       s.Lat,
			Lon:       s.Lon,
			Sequence:  s.Sequence,
		}
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "stop_id"},
			{Name: "route_id"},
			{Name: "direction"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"name":       gorm.Expr("EXCLUDED.name"),
			"lat":        gorm.Expr("EXCLUDED.lat"),
			"lon":        gorm.Expr("EXCLUDED.lon"),
			"sequence":   gorm.Expr("EXCLUDED.sequence"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).CreateInBatches(&models, 500).Error
}
