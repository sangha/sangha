package projects

import (
	"strconv"

	"gitlab.techcultivation.org/techcultivation/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *ProjectResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *ProjectResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *ProjectResource) GetDoc() string {
	return "retrieve projects"
}

// GetParams returns the parameters supported by this API endpoint
func (r *ProjectResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("name", "name of a project").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *ProjectResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := ProjectResponse{}
	resp.Init(context)

	for _, id := range ids {
		iid, err := strconv.Atoi(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}
		project, err := context.(*db.APIContext).GetProjectByID(int64(iid))
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddProject(&project)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *ProjectResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := ProjectResponse{}
	resp.Init(context)

	projects, err := context.(*db.APIContext).LoadAllProjects()
	if err != nil {
		r.NotFound(request, response)
		return
	}

	for _, project := range projects {
		resp.AddProject(&project)
	}

	resp.Send(response)
}
