package db

import (
	"math"
	"strconv"
)

// Statistics represents the db schema of a statistic
type Statistics struct {
	ID            string
	ProjectID     int64
	BudgetID      int64
	MonthlyChange int64
	PastMonths    []int64
}

// LoadStatistics loads statistics for a project from the database
func (context *APIContext) LoadStatistics(projectID int64) (Statistics, error) {
	p, err := context.GetProjectByID(projectID)
	if err != nil {
		return Statistics{}, err
	}

	bal, err := p.Balance(context)
	if err != nil {
		return Statistics{}, err
	}

	stats := Statistics{
		ID:        "stats_" + strconv.FormatInt(projectID, 10),
		ProjectID: projectID,
	}

	bs, err := p.BalanceStats(context)
	if err == nil {
		for _, v := range bs {
			stats.PastMonths = append(stats.PastMonths, v)
		}
	}

	max := math.Min(float64(len(stats.PastMonths))-1, 11)
	start := stats.PastMonths[int(max)]
	stats.MonthlyChange = int64(float64(bal-start) / (max + 1))

	return stats, nil
}
