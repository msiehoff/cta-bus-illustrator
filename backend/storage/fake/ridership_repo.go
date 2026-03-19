package fake

import (
	"time"

	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

type RidershipRepo struct{}

func (r *RidershipRepo) GetLatestMonth() (time.Time, error) {
	return time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC), nil
}

func (r *RidershipRepo) UpsertBatch(_ []business.RidershipRecord) error {
	return nil
}

func (r *RidershipRepo) GetByMonth(_ time.Time, ridershipType business.RidershipType) (map[string]*business.RidershipRecord, error) {
	avgRides := map[business.RidershipType]float64{
		business.RidershipTypeWeekday:  8500,
		business.RidershipTypeSaturday: 5000,
		business.RidershipTypeSunday:   3500,
	}

	return map[string]*business.RidershipRecord{
		"66": {
			RouteExternalID: "66",
			MonthBeginning:  time.Date(2025, 11, 1, 0, 0, 0, 0, time.UTC),
			Type:            ridershipType,
			AvgRides:        avgRides[ridershipType],
		},
	}, nil
}
