package api

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/msiehoff/cta-bus-illustrator/backend/business"
)

func (a *API) handleImportRidership(c *gin.Context) {
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "multipart file field 'file' required"})
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// reader.Comma = '\t'
	reader.LazyQuotes = true

	// skip header row
	if _, err := reader.Read(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to read header row"})
		return
	}

	var (
		records   []business.RidershipRecord
		parseErrs []string
		lineNum   = 1
	)

	for {
		lineNum++
		row, err := reader.Read()
		if err != nil {
			break // EOF
		}
		if len(row) < 6 {
			parseErrs = append(parseErrs, fmt.Sprintf("line %d: expected 6+ columns, got %d", lineNum, len(row)))
			continue
		}

		externalID := strings.TrimSpace(row[0])

		month, err := time.Parse("01/02/2006", strings.TrimSpace(row[2]))
		if err != nil {
			parseErrs = append(parseErrs, fmt.Sprintf("line %d: invalid month %q", lineNum, row[2]))
			continue
		}
		month = time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, time.UTC)

		weekday, err := parseRides(row[3])
		if err != nil {
			parseErrs = append(parseErrs, fmt.Sprintf("line %d: invalid weekday rides %q", lineNum, row[3]))
			continue
		}
		saturday, err := parseRides(row[4])
		if err != nil {
			parseErrs = append(parseErrs, fmt.Sprintf("line %d: invalid saturday rides %q", lineNum, row[4]))
			continue
		}
		sunday, err := parseRides(row[5])
		if err != nil {
			parseErrs = append(parseErrs, fmt.Sprintf("line %d: invalid sunday rides %q", lineNum, row[5]))
			continue
		}

		records = append(records,
			business.RidershipRecord{RouteExternalID: externalID, MonthBeginning: month, Type: business.RidershipTypeWeekday, AvgRides: weekday},
			business.RidershipRecord{RouteExternalID: externalID, MonthBeginning: month, Type: business.RidershipTypeSaturday, AvgRides: saturday},
			business.RidershipRecord{RouteExternalID: externalID, MonthBeginning: month, Type: business.RidershipTypeSunday, AvgRides: sunday},
		)
	}

	if err := a.routeService.ImportRidership(records); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"imported": len(records),
		"errors":   parseErrs,
	})
}

// parseRides strips commas and parses a ridership number (e.g. "1,126" → 1126).
func parseRides(s string) (float64, error) {
	s = strings.ReplaceAll(strings.TrimSpace(s), ",", "")
	return strconv.ParseFloat(s, 64)
}
