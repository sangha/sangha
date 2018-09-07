package statistics

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// StatisticsResponse is the common response to 'statistics' requests
type StatisticsResponse struct {
	smolder.Response

	Statistics []statisticsInfoResponse `json:"statistics,omitempty"`
	statistics []db.Statistics
}

type statisticsInfoResponse struct {
	ID            string  `json:"id"`
	ProjectID     int64   `json:"project_id"`
	BudgetID      int64   `json:"budget_id"`
	MonthlyChange int64   `json:"monthly_change"`
	PastMonths    []int64 `json:"past_months"`
}

// Init a new response
func (r *StatisticsResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Statistics = []statisticsInfoResponse{}
}

// AddStatistics adds a statistics to the response
func (r *StatisticsResponse) AddStatistics(statistics *db.Statistics) {
	r.statistics = append(r.statistics, *statistics)
	r.Statistics = append(r.Statistics, prepareStatisticsResponse(r.Context, statistics))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *StatisticsResponse) EmptyResponse() interface{} {
	if len(r.statistics) == 0 {
		var out struct {
			Statistics interface{} `json:"statistics"`
		}
		out.Statistics = []statisticsInfoResponse{}
		return out
	}
	return nil
}

func prepareStatisticsResponse(context smolder.APIContext, statistics *db.Statistics) statisticsInfoResponse {
	resp := statisticsInfoResponse{
		ID:            statistics.ID,
		ProjectID:     statistics.ProjectID,
		MonthlyChange: statistics.MonthlyChange,
		PastMonths:    statistics.PastMonths,
	}

	return resp
}
