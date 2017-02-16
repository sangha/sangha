package projects

import (
	"net/http"

	"gitlab.techcultivation.org/techcultivation/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProjectPostStruct holds all values of an incoming POST request
type ProjectPostStruct struct {
	Project struct {
		Name  string `json:"name"`
		About string `json:"about"`
	} `json:"project"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *ProjectResource) PostAuthRequired() bool {
	return true
}

// PostDoc returns the description of this API endpoint
func (r *ProjectResource) PostDoc() string {
	return "create a new project invitation"
}

// PostParams returns the parameters supported by this API endpoint
func (r *ProjectResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *ProjectResource) Post(context smolder.APIContext, request *restful.Request, response *restful.Response) {
	auth, err := context.Authentication(request)
	if err != nil || auth.(db.Project).ID != 1 {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			false,
			"Admin permission required for this operation",
			"ProjectResource POST"))
		return
	}

	ups := ProjectPostStruct{}
	err = request.ReadEntity(&ups)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
			http.StatusBadRequest,
			false,
			"Can't parse POST data",
			"ProjectResource POST"))
		return
	}

	project := db.Project{
		Name:  ups.Project.Name,
		About: ups.Project.About,
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
