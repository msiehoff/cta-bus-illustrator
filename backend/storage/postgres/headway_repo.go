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

func (r *HeadwayRepo) List(_ context.Context, filter app.HeadwayListFilter) ([]business.Headway, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	return r.list(filter, limit, filter.Offset)
}

func (r *HeadwayRepo) ListAll(_ context.Context, filter app.HeadwayListFilter) ([]business.Headway, error) {
	return r.list(filter, 50_000, 0)
}

func (r *HeadwayRepo) list(filter app.HeadwayListFilter, limit, offset int) ([]business.Headway, error) {
	order := "headways.timestamp DESC"
	if filter.SortAsc {
		order = "headways.timestamp ASC"
	}

	type row struct {
		StopID         string    `gorm:"column:stop_id"`
		RouteID        string    `gorm:"column:route_id"`
		Direction      string    `gorm:"column:direction"`
		Timestamp      time.Time `gorm:"column:timestamp"`
		HeadwayMinutes float64   `gorm:"column:headway_minutes"`
		FromVehicleID  string    `gorm:"column:from_vehicle_id"`
		ToVehicleID    string    `gorm:"column:to_vehicle_id"`
		StopName       *string   `gorm:"column:stop_name"`
		RouteName      *string   `gorm:"column:route_name"`
	}

	query := r.db.Table("headways").
		Select(`headways.stop_id, headways.route_id, headways.direction, headways.timestamp,
			headways.headway_minutes, headways.from_vehicle_id, headways.to_vehicle_id,
			stops.name AS stop_name, routes.name AS route_name`).
		Joins("LEFT JOIN stops ON stops.stop_id = headways.stop_id AND stops.route_id = headways.route_id AND stops.direction = headways.direction").
		Joins("LEFT JOIN routes ON routes.external_id = headways.route_id AND routes.deleted_at IS NULL")
	query = applyHeadwayListFilter(query, filter)

	var rows []row
	if err := query.Order(order).Limit(limit).Offset(offset).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]business.Headway, len(rows))
	for i, row := range rows {
		out[i] = business.Headway{
			StopID:         row.StopID,
			RouteID:        row.RouteID,
			Direction:      row.Direction,
			Timestamp:      row.Timestamp,
			HeadwayMinutes: row.HeadwayMinutes,
			FromVehicleID:  row.FromVehicleID,
			ToVehicleID:    row.ToVehicleID,
		}
		if row.StopName != nil {
			out[i].StopName = *row.StopName
		}
		if row.RouteName != nil {
			out[i].RouteName = *row.RouteName
		}
	}
	return out, nil
}

func (r *HeadwayRepo) Count(_ context.Context, filter app.HeadwayListFilter) (int64, error) {
	var count int64
	query := r.db.Table("headways").
		Joins("LEFT JOIN stops ON stops.stop_id = headways.stop_id AND stops.route_id = headways.route_id AND stops.direction = headways.direction")
	query = applyHeadwayListFilter(query, filter)
	err := query.Count(&count).Error
	return count, err
}

func applyHeadwayListFilter(query *gorm.DB, filter app.HeadwayListFilter) *gorm.DB {
	if filter.RouteID != "" {
		query = query.Where("headways.route_id = ?", filter.RouteID)
	}
	if filter.Direction != "" {
		query = query.Where("headways.direction = ?", filter.Direction)
	}
	if filter.VehicleID != "" {
		query = query.Where(
			"headways.from_vehicle_id = ? OR headways.to_vehicle_id = ?",
			filter.VehicleID, filter.VehicleID,
		)
	}
	if filter.Stop != "" {
		like := "%" + filter.Stop + "%"
		query = query.Where(
			"headways.stop_id = ? OR stops.name ILIKE ?",
			filter.Stop, like,
		)
	}
	if filter.From != nil {
		query = query.Where("headways.timestamp >= ?", *filter.From)
	}
	if filter.To != nil {
		query = query.Where("headways.timestamp < ?", *filter.To)
	}
	return query
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
		"summaries_written":  model.SummariesWritten,
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
		SummariesWritten:  run.SummariesWritten,
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
		SummariesWritten:  m.SummariesWritten,
		ErrorMessage:      m.ErrorMessage,
	}
}

type HeadwaySummaryRepo struct {
	db *gorm.DB
}

func NewHeadwaySummaryRepo(db *gorm.DB) *HeadwaySummaryRepo {
	return &HeadwaySummaryRepo{db: db}
}

func (r *HeadwaySummaryRepo) DeleteForServiceDate(_ context.Context, serviceDate time.Time) (int64, error) {
	day := time.Date(serviceDate.Year(), serviceDate.Month(), serviceDate.Day(), 0, 0, 0, 0, time.UTC)
	res := r.db.Where("service_date = ?", day).Delete(&headwaySummaryModel{})
	return res.RowsAffected, res.Error
}

func (r *HeadwaySummaryRepo) InsertBatch(_ context.Context, summaries []business.HeadwaySummary) error {
	if len(summaries) == 0 {
		return nil
	}
	models := make([]headwaySummaryModel, len(summaries))
	for i, s := range summaries {
		models[i] = headwaySummaryModel{
			ServiceDate:      s.ServiceDate,
			WindowStart:      s.WindowStart,
			WindowEnd:        s.WindowEnd,
			Grain:            s.Grain,
			Method:           s.Method,
			StopID:           s.StopID,
			RouteID:          s.RouteID,
			Direction:        s.Direction,
			ObservationCount: s.Count,
			MeanMinutes:      s.MeanMinutes,
			MedianMinutes:    s.MedianMinutes,
			StdDevMinutes:    s.StdDevMinutes,
			CV:               s.CV,
			AvgWaitMinutes:   s.AvgWaitMinutes,
		}
	}
	return r.db.CreateInBatches(&models, 500).Error
}

func (r *HeadwaySummaryRepo) List(_ context.Context, filter app.HeadwaySummaryFilter) ([]business.HeadwaySummary, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	order := "headway_summaries.service_date DESC, headway_summaries.mean_minutes DESC"
	if filter.SortAsc {
		order = "headway_summaries.service_date ASC, headway_summaries.mean_minutes ASC"
	}

	type row struct {
		ServiceDate      time.Time `gorm:"column:service_date"`
		WindowStart      time.Time `gorm:"column:window_start"`
		WindowEnd        time.Time `gorm:"column:window_end"`
		Grain            string    `gorm:"column:grain"`
		Method           string    `gorm:"column:method"`
		StopID           string    `gorm:"column:stop_id"`
		RouteID          string    `gorm:"column:route_id"`
		Direction        string    `gorm:"column:direction"`
		ObservationCount int       `gorm:"column:observation_count"`
		MeanMinutes      float64   `gorm:"column:mean_minutes"`
		MedianMinutes    float64   `gorm:"column:median_minutes"`
		StdDevMinutes    float64   `gorm:"column:stddev_minutes"`
		CV               float64   `gorm:"column:cv"`
		AvgWaitMinutes   float64   `gorm:"column:avg_wait_minutes"`
		StopName         *string   `gorm:"column:stop_name"`
		RouteName        *string   `gorm:"column:route_name"`
	}

	query := r.db.Table("headway_summaries").
		Select(`headway_summaries.*, stops.name AS stop_name, routes.name AS route_name`).
		Joins(`LEFT JOIN stops ON stops.stop_id = headway_summaries.stop_id
			AND stops.route_id = headway_summaries.route_id
			AND stops.direction = headway_summaries.direction`).
		Joins(`LEFT JOIN routes ON routes.external_id = headway_summaries.route_id AND routes.deleted_at IS NULL`)
	query = applyHeadwaySummaryFilter(query, filter)

	var rows []row
	if err := query.Order(order).Limit(limit).Offset(filter.Offset).Scan(&rows).Error; err != nil {
		return nil, err
	}

	out := make([]business.HeadwaySummary, len(rows))
	for i, row := range rows {
		out[i] = business.HeadwaySummary{
			ServiceDate: row.ServiceDate,
			WindowStart: row.WindowStart,
			WindowEnd:   row.WindowEnd,
			Grain:       row.Grain,
			Method:      row.Method,
			StopID:      row.StopID,
			RouteID:     row.RouteID,
			Direction:   row.Direction,
			HeadwaySummaryStats: business.HeadwaySummaryStats{
				Count:          row.ObservationCount,
				MeanMinutes:    row.MeanMinutes,
				MedianMinutes:  row.MedianMinutes,
				StdDevMinutes:  row.StdDevMinutes,
				CV:             row.CV,
				AvgWaitMinutes: row.AvgWaitMinutes,
			},
		}
		if row.StopName != nil {
			out[i].StopName = *row.StopName
		}
		if row.RouteName != nil {
			out[i].RouteName = *row.RouteName
		}
	}
	return out, nil
}

func (r *HeadwaySummaryRepo) Count(_ context.Context, filter app.HeadwaySummaryFilter) (int64, error) {
	var count int64
	query := r.db.Table("headway_summaries").
		Joins(`LEFT JOIN stops ON stops.stop_id = headway_summaries.stop_id
			AND stops.route_id = headway_summaries.route_id
			AND stops.direction = headway_summaries.direction`)
	query = applyHeadwaySummaryFilter(query, filter)
	err := query.Count(&count).Error
	return count, err
}

func applyHeadwaySummaryFilter(query *gorm.DB, filter app.HeadwaySummaryFilter) *gorm.DB {
	if !filter.ServiceDate.IsZero() {
		day := time.Date(filter.ServiceDate.Year(), filter.ServiceDate.Month(), filter.ServiceDate.Day(), 0, 0, 0, 0, time.UTC)
		query = query.Where("headway_summaries.service_date = ?", day)
	}
	if filter.Grain != "" {
		query = query.Where("headway_summaries.grain = ?", filter.Grain)
	}
	if filter.Method != "" {
		query = query.Where("headway_summaries.method = ?", filter.Method)
	}
	if filter.RouteID != "" {
		query = query.Where("headway_summaries.route_id = ?", filter.RouteID)
	}
	if filter.Direction != "" {
		query = query.Where("headway_summaries.direction = ?", filter.Direction)
	}
	if filter.StopID != "" {
		query = query.Where("headway_summaries.stop_id = ?", filter.StopID)
	}
	if filter.Stop != "" {
		like := "%" + filter.Stop + "%"
		query = query.Where(
			"headway_summaries.stop_id = ? OR stops.name ILIKE ?",
			filter.Stop, like,
		)
	}
	return query
}
