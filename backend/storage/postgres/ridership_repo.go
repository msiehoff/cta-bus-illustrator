package postgres

import (
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"gorm.io/gorm"
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
