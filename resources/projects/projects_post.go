package projects

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProjectPostStruct holds all values of an incoming POST request
type ProjectPostStruct struct {
	Project struct {
		Slug       string `json:"slug"`
		Name       string `json:"name"`
		About      string `json:"about"`
		Website    string `json:"website"`
		License    string `json:"license"`
		Repository string `json:"repository"`
	} `json:"project"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *ProjectResource) PostAuthRequired() bool {
	return true
}

// PostDoc returns the description of this API endpoint
func (r *ProjectResource) PostDoc() string {
	return "create a new project"
}

// PostParams returns the parameters supported by this API endpoint
func (r *ProjectResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *ProjectResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	/*auth, err := context.Authentication(request)
		if err != nil || auth.(db.Project).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"ProjectResource POST"))
		return
	}*/

	ups := data.(*ProjectPostStruct)
	_, err := context.(*db.APIContext).LoadProjectBySlug(ups.Project.Slug)
	if err == nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"A project with this slug address already exists",
			"ProjectResource POST"))
		return
	}

	project := db.Project{
		Slug:       ups.Project.Slug,
		Name:       ups.Project.Name,
		About:      ups.Project.About,
		Website:    ups.Project.Website,
		License:    ups.Project.License,
		Repository: ups.Project.Repository,
	}

	err = project.Save(context.(*db.APIContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			true,
			"Can't create project",
			"ProjectResource POST"))
		return
	}

	resp := ProjectResponse{}
	resp.Init(context)
	resp.AddProject(&project)
	resp.Send(response)
}
