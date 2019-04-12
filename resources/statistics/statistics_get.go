package statistics

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *StatisticsResource) GetAuthRequired() bool {
	return true
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *StatisticsResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *StatisticsResource) GetDoc() string {
	return "retrieve statistics"
}

// GetParams returns the parameters supported by this API endpoint
func (r *StatisticsResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("project", "an ID of a project").DataType("string"))
	params = append(params, restful.QueryParameter("budget", "an ID of a budget").DataType("string"))

	return params
}

// Get sends out items matching the query parameters
func (r *StatisticsResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := StatisticsResponse{}
	resp.Init(context)

	projectID := params["project"]
	if len(projectID) > 0 {
		b, err := context.(*db.APIContext).LoadProjectByUUID(projectID[0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		statistics, err := context.(*db.APIContext).LoadStatistics(b.ID)
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddStatistics(&statistics)
	} else {
		r.NotFound(request, response)
		return
	}

	resp.Send(response)
}
