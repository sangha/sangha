package projects

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

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
	params = append(params, restful.QueryParameter("slug", "slug of a project").DataType("string"))
	params = append(params, restful.QueryParameter("name", "name of a project").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *ProjectResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := ProjectResponse{}
	resp.Init(context)

	for _, id := range ids {
		project, err := context.(*db.APIContext).GetProjectByUUID(id)
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

	if len(params["slug"]) > 0 {
		project, err := context.(*db.APIContext).LoadProjectBySlug(params["slug"][0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddProject(&project)
	} else {
		projects, err := context.(*db.APIContext).LoadAllProjects()
		if err != nil {
			r.NotFound(request, response)
			return
		}

		for _, project := range projects {
			resp.AddProject(&project)
		}
	}

	resp.Send(response)
}
