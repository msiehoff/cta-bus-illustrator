package postgres

import (
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RidershipRepo struct {
	db *gorm.DB
}

func NewRidershipRepo(db *gorm.DB) *RidershipRepo {
	return &RidershipRepo{db: db}
}

func (r *RidershipRepo) GetLatestMonth() (time.Time, error) {
	var model ridershipModel
	if err := r.db.Order("month_beginning DESC").First(&model).Error; err != nil {
		return time.Time{}, err
	}
	return model.MonthBeginning, nil
}

func (r *RidershipRepo) GetAvailableMonths() ([]time.Time, error) {
	var months []time.Time
	err := r.db.Model(&ridershipModel{}).
		Where("deleted_at IS NULL").
		Distinct("month_beginning").
		Order("month_beginning DESC").
		Pluck("month_beginning", &months).Error
	return months, err
}

func (r *RidershipRepo) GetByMonth(month time.Time, ridershipType business.RidershipType) (map[string]*business.RidershipRecord, error) {
	var rows []ridershipRow
	err := r.db.Model(&ridershipModel{}).
		Select("routes.external_id, ridership.month_beginning, ridership.type, ridership.avg_rides").
		Joins("JOIN routes ON routes.id = ridership.route_id AND routes.deleted_at IS NULL").
		Where("ridership.month_beginning = ? AND ridership.type = ? AND ridership.deleted_at IS NULL", month, ridershipType).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string]*business.RidershipRecord, len(rows))
	for _, row := range rows {
		result[row.ExternalID] = &business.RidershipRecord{
			RouteExternalID: row.ExternalID,
			MonthBeginning:  row.MonthBeginning,
			Type:            business.RidershipType(row.Type),
			AvgRides:        row.AvgRides,
		}
	}
	return result, nil
}

// GetAllByRoute returns all ridership records for a single route across all months and types,
// ordered chronologically.
func (r *RidershipRepo) GetAllByRoute(routeExternalID string) ([]business.RidershipRecord, error) {
	var rows []ridershipRow
	err := r.db.Model(&ridershipModel{}).
		Select("routes.external_id, ridership.month_beginning, ridership.type, ridership.avg_rides").
		Joins("JOIN routes ON routes.id = ridership.route_id AND routes.deleted_at IS NULL").
		Where("routes.external_id = ? AND ridership.deleted_at IS NULL", routeExternalID).
		Order("ridership.month_beginning ASC, ridership.type ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	records := make([]business.RidershipRecord, len(rows))
	for i, row := range rows {
		records[i] = business.RidershipRecord{
			RouteExternalID: row.ExternalID,
			MonthBeginning:  row.MonthBeginning,
			Type:            business.RidershipType(row.Type),
			AvgRides:        row.AvgRides,
		}
	}
	return records, nil
}

// GetSystemTotals returns ridership summed across all routes, grouped by month and type,
// ordered chronologically.
func (r *RidershipRepo) GetSystemTotals() ([]business.RidershipRecord, error) {
	type systemRow struct {
		MonthBeginning time.Time
		Type           string
		AvgRides       float64
	}

	var rows []systemRow
	err := r.db.Model(&ridershipModel{}).
		Select("ridership.month_beginning, ridership.type, SUM(ridership.avg_rides) as avg_rides").
		Joins("JOIN routes ON routes.id = ridership.route_id AND routes.deleted_at IS NULL").
		Where("ridership.deleted_at IS NULL").
		Group("ridership.month_beginning, ridership.type").
		Order("ridership.month_beginning ASC, ridership.type ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	records := make([]business.RidershipRecord, len(rows))
	for i, row := range rows {
		records[i] = business.RidershipRecord{
			MonthBeginning: row.MonthBeginning,
			Type:           business.RidershipType(row.Type),
			AvgRides:       row.AvgRides,
		}
	}
	return records, nil
}

func (r *RidershipRepo) UpsertBatch(records []business.RidershipRecord) error {
	if len(records) == 0 {
		return nil
	}

	idSet := make(map[string]struct{}, len(records))
	for _, rec := range records {
		idSet[rec.RouteExternalID] = struct{}{}
	}
	externalIDs := make([]string, 0, len(idSet))
	for id := range idSet {
		externalIDs = append(externalIDs, id)
	}

	var routes []routeModel
	if err := r.db.Where("external_id IN ?", externalIDs).Find(&routes).Error; err != nil {
		return err
	}

	routeIDByExternal := make(map[string]uint, len(routes))
	for _, route := range routes {
		routeIDByExternal[route.ExternalID] = route.ID
	}

	models := make([]ridershipModel, 0, len(records))
	for _, rec := range records {
		routeID, ok := routeIDByExternal[rec.RouteExternalID]
		if !ok {
			continue
		}
		models = append(models, ridershipModel{
			RouteID:        routeID,
			MonthBeginning: rec.MonthBeginning,
			Type:           string(rec.Type),
			AvgRides:       rec.AvgRides,
		})
	}

	return r.db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "route_id"},
			{Name: "month_beginning"},
			{Name: "type"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"avg_rides":  gorm.Expr("EXCLUDED.avg_rides"),
			"updated_at": gorm.Expr("NOW()"),
		}),
	}).CreateInBatches(&models, 500).Error
}
