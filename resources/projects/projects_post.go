package projects

import (
	"encoding/base64"
	"log"
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProjectPostStruct holds all values of an incoming POST request
type ProjectPostStruct struct {
	Project struct {
		Slug           string `json:"slug"`
		Name           string `json:"name"`
		Summary        string `json:"summary"`
		About          string `json:"about"`
		Website        string `json:"website"`
		License        string `json:"license"`
		Repository     string `json:"repository"`
		Logo           string `json:"logo"`
		Private        bool   `json:"private"`
		PrivateBalance bool   `json:"private_balance"`
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
	auth, err := context.Authentication(request)
	if err != nil || auth.(db.User).ID != 1 {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusUnauthorized,
			"Admin permission required for this operation",
			"ProjectResource POST"))
		return
	}

	ctx := context.(*db.APIContext)
	ups := data.(*ProjectPostStruct)
	_, err = ctx.LoadProjectBySlug(ups.Project.Slug)
	if err == nil {
		smolder.ErrorResponseHandler(request, response, nil, smolder.NewErrorResponse(
			http.StatusBadRequest,
			"A project with this slug address already exists",
			"ProjectResource POST"))
		return
	}

	project := db.Project{
		Slug:           ups.Project.Slug,
		Name:           ups.Project.Name,
		Summary:        ups.Project.Summary,
		About:          ups.Project.About,
		Website:        ups.Project.Website,
		License:        ups.Project.License,
		Repository:     ups.Project.Repository,
		Private:        false,
		PrivateBalance: true,
	}

	if len(ups.Project.Logo) > 0 {
		logo, err := base64.StdEncoding.DecodeString(ups.Project.Logo)
		if err == nil {
			project.Logo, err = ctx.StoreImage(logo)
			if err != nil {
				log.Println("WARNING: could not store image:", err)
			}
		} else {
			log.Println("WARNING: could not decode logo:", err)
		}
	}

	err = project.Save(ctx)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't create project",
			"ProjectResource POST"))
		return
	}

	budget := db.Budget{
		ProjectID:      &project.ID,
		ParentID:       0,
		Name:           project.Name,
		Private:        false,
		PrivateBalance: true,
	}
	err = budget.Save(ctx)
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't create budget for new project",
			"ProjectResource POST"))
		return
	}

	resp := ProjectResponse{}
	resp.Init(context)
	resp.AddProject(&project)
	resp.Send(response)
}
