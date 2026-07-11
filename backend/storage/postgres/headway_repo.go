package postgres

import (
	"context"
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/app"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
	"gorm.io/gorm"
)

type HeadwayRepo struct {
	db *gorm.DB
}

func NewHeadwayRepo(db *gorm.DB) *HeadwayRepo {
	return &HeadwayRepo{db: db}
}

func (r *HeadwayRepo) DeleteInRange(_ context.Context, start, end time.Time) (int64, error) {
	res := r.db.Where("timestamp >= ? AND timestamp < ?", start, end).Delete(&headwayModel{})
	return res.RowsAffected, res.Error
}

func (r *HeadwayRepo) InsertBatch(_ context.Context, headways []business.Headway) error {
	if len(headways) == 0 {
		return nil
	}
	models := make([]headwayModel, len(headways))
	for i, h := range headways {
		models[i] = headwayModel{
			StopID:         h.StopID,
			RouteID:        h.RouteID,
			Direction:      h.Direction,
			Timestamp:      h.Timestamp,
			HeadwayMinutes: h.HeadwayMinutes,
			FromVehicleID:  h.FromVehicleID,
			ToVehicleID:    h.ToVehicleID,
		}
	}
	return r.db.CreateInBatches(&models, 500).Error
}

func (r *HeadwayRepo) ListInRange(_ context.Context, start, end time.Time, filter app.HeadwayListFilter) ([]business.Headway, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}

	query := r.db.Model(&headwayModel{}).
		Where("timestamp >= ? AND timestamp < ?", start, end)
	query = applyHeadwayListFilter(query, filter)

	var models []headwayModel
	if err := query.Order("timestamp ASC").Limit(limit).Offset(filter.Offset).Find(&models).Error; err != nil {
		return nil, err
	}

	out := make([]business.Headway, len(models))
	for i, m := range models {
		out[i] = toBusinessHeadway(m)
	}
	return out, nil
}

func (r *HeadwayRepo) CountInRange(_ context.Context, start, end time.Time, filter app.HeadwayListFilter) (int64, error) {
	var count int64
	query := r.db.Model(&headwayModel{}).
		Where("timestamp >= ? AND timestamp < ?", start, end)
	query = applyHeadwayListFilter(query, filter)
	err := query.Count(&count).Error
	return count, err
}

func applyHeadwayListFilter(query *gorm.DB, filter app.HeadwayListFilter) *gorm.DB {
	if filter.RouteID != "" {
		query = query.Where("route_id = ?", filter.RouteID)
	}
	if filter.Direction != "" {
		query = query.Where("direction = ?", filter.Direction)
	}
	if filter.StopID != "" {
		query = query.Where("stop_id = ?", filter.StopID)
	}
	return query
}

func toBusinessHeadway(m headwayModel) business.Headway {
	return business.Headway{
		StopID:         m.StopID,
		RouteID:        m.RouteID,
		Direction:      m.Direction,
		Timestamp:      m.Timestamp,
		HeadwayMinutes: m.HeadwayMinutes,
		FromVehicleID:  m.FromVehicleID,
		ToVehicleID:    m.ToVehicleID,
	}
}

type HeadwayJobRunRepo struct {
	db *gorm.DB
}

func NewHeadwayJobRunRepo(db *gorm.DB) *HeadwayJobRunRepo {
	return &HeadwayJobRunRepo{db: db}
}

func (r *HeadwayJobRunRepo) Create(_ context.Context, run business.HeadwayJobRun) (business.HeadwayJobRun, error) {
	model := toJobRunModel(run)
	if err := r.db.Create(&model).Error; err != nil {
		return business.HeadwayJobRun{}, err
	}
	return fromJobRunModel(model), nil
}

func (r *HeadwayJobRunRepo) Update(_ context.Context, run business.HeadwayJobRun) error {
	model := toJobRunModel(run)
	return r.db.Model(&headwayJobRunModel{}).Where("id = ?", run.ID).Updates(map[string]interface{}{
		"status":             model.Status,
		"finished_at":        model.FinishedAt,
		"arrivals_processed": model.ArrivalsProcessed,
		"headways_written":   model.HeadwaysWritten,
		"error_message":      model.ErrorMessage,
		"updated_at":         time.Now().UTC(),
	}).Error
}

func (r *HeadwayJobRunRepo) List(_ context.Context, limit, offset int) ([]business.HeadwayJobRun, error) {
	if limit <= 0 {
		limit = 50
	}
	var models []headwayJobRunModel
	if err := r.db.Order("started_at DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]business.HeadwayJobRun, len(models))
	for i, m := range models {
		out[i] = fromJobRunModel(m)
	}
	return out, nil
}

func (r *HeadwayJobRunRepo) Get(_ context.Context, id int64) (business.HeadwayJobRun, error) {
	var model headwayJobRunModel
	if err := r.db.First(&model, id).Error; err != nil {
		return business.HeadwayJobRun{}, err
	}
	return fromJobRunModel(model), nil
}

func toJobRunModel(run business.HeadwayJobRun) headwayJobRunModel {
	return headwayJobRunModel{
		ID:                uint(run.ID),
		ServiceDate:       run.ServiceDate,
		Status:            string(run.Status),
		TriggeredBy:       string(run.TriggeredBy),
		StartedAt:         run.StartedAt,
		FinishedAt:        run.FinishedAt,
		ArrivalsProcessed: run.ArrivalsProcessed,
		HeadwaysWritten:   run.HeadwaysWritten,
		ErrorMessage:      run.ErrorMessage,
	}
}

func fromJobRunModel(m headwayJobRunModel) business.HeadwayJobRun {
	return business.HeadwayJobRun{
		ID:                int64(m.ID),
		ServiceDate:       m.ServiceDate,
		Status:            business.HeadwayJobStatus(m.Status),
		TriggeredBy:       business.HeadwayJobTrigger(m.TriggeredBy),
		StartedAt:         m.StartedAt,
		FinishedAt:        m.FinishedAt,
		ArrivalsProcessed: m.ArrivalsProcessed,
		HeadwaysWritten:   m.HeadwaysWritten,
		ErrorMessage:      m.ErrorMessage,
	}
}
