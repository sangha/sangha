package searches

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// SearchResponse is the common response to 'searches' requests
type SearchResponse struct {
	smolder.Response

	Searches []searchInfoResponse `json:"searches,omitempty"`
	searches []db.Search
}

type searchInfoResponse struct {
	ID       string   `json:"id"`
	Projects []string `json:"projects"`
	Budgets  []string `json:"budgets"`
	Payments []int64  `json:"payments"`
}

// Init a new response
func (r *SearchResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Searches = []searchInfoResponse{}
}

// AddSearches adds a search to the response
func (r *SearchResponse) AddSearch(search *db.Search) {
	r.searches = append(r.searches, *search)
	r.Searches = append(r.Searches, prepareSearchResponse(r.Context, search))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *SearchResponse) EmptyResponse() interface{} {
	if len(r.searches) == 0 {
		var out struct {
			Searches interface{} `json:"searches"`
		}
		out.Searches = []searchInfoResponse{}
		return out
	}
	return nil
}

func prepareSearchResponse(context smolder.APIContext, search *db.Search) searchInfoResponse {
	resp := searchInfoResponse{
		ID: search.ID,
	}

	for _, p := range search.Projects {
		resp.Projects = append(resp.Projects, p.UUID)
	}
	for _, b := range search.Budgets {
		resp.Budgets = append(resp.Budgets, b.UUID)
	}
	for _, p := range search.Payments {
		resp.Payments = append(resp.Payments, p.ID)
	}

	return resp
}
