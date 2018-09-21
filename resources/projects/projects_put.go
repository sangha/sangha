package projects

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProjectPutStruct holds all values of an incoming PUT request
type ProjectPutStruct struct {
	ProjectPostStruct
}

// PutAuthRequired returns true because all requests need authentication
func (r *ProjectResource) PutAuthRequired() bool {
	return true
}

// PutDoc returns the description of this API endpoint
func (r *ProjectResource) PutDoc() string {
	return "update an existing project"
}

// PutParams returns the parameters supported by this API endpoint
func (r *ProjectResource) PutParams() []*restful.Parameter {
	return nil
}

// Put processes an incoming PUT (update) request
func (r *ProjectResource) Put(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	resp := ProjectResponse{}
	resp.Init(context)

	project, err := context.(*db.APIContext).GetProjectByUUID(request.PathParameter("project-id"))
	if err != nil {
		r.NotFound(request, response)
		return
	}

	/*	auth, err := context.Authentication(request)
		if err != nil || (auth.(db.User).ID != 1 && auth.(db.User).ID != project.UserID) {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"ProjectResource PUT"))
			return
		} */

	pps := data.(*ProjectPostStruct)
	project.Name = pps.Project.Name
	project.Summary = pps.Project.Summary
	project.About = pps.Project.About
	project.Website = pps.Project.Website
	project.License = pps.Project.License
	project.Repository = pps.Project.Repository

	err = project.Update(context.(*db.APIContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't update project",
			"ProjectResource PUT"))
		return
	}

	resp.AddProject(&project)
	resp.Send(response)
}
