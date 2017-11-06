package projects

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// ProjectResponse is the common response to 'project' requests
type ProjectResponse struct {
	smolder.Response

	Projects []projectInfoResponse `json:"projects,omitempty"`
	projects []db.Project
}

type contributorResponse struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type projectInfoResponse struct {
	ID           string                `json:"id"`
	Slug         string                `json:"slug"`
	Name         string                `json:"name"`
	Summary      string                `json:"summary"`
	About        string                `json:"about"`
	Website      string                `json:"website"`
	License      string                `json:"license"`
	Repository   string                `json:"repository"`
	Logo         string                `json:"logo"`
	RootBudget   string                `json:"budget_root"`
	Contributors []contributorResponse `json:"contributors,omitempty"`
	Activated    bool                  `json:"activated"`
}

// Init a new response
func (r *ProjectResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Projects = []projectInfoResponse{}
}

// AddProject adds a project to the response
func (r *ProjectResponse) AddProject(project *db.Project) {
	r.projects = append(r.projects, *project)
	r.Projects = append(r.Projects, prepareProjectResponse(r.Context, project))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *ProjectResponse) EmptyResponse() interface{} {
	if len(r.projects) == 0 {
		var out struct {
			Projects interface{} `json:"projects"`
		}
		out.Projects = []projectInfoResponse{}
		return out
	}
	return nil
}

func prepareProjectResponse(context smolder.APIContext, project *db.Project) projectInfoResponse {
	ctx := context.(*db.APIContext)
	resp := projectInfoResponse{
		ID:         project.UUID,
		Slug:       project.Slug,
		Name:       project.Name,
		Summary:    project.Summary,
		About:      project.About,
		Website:    project.Website,
		License:    project.License,
		Repository: project.Repository,
		Logo:       ctx.BuildImageURL(project.Logo, project.Name),
		Activated:  project.Activated,
	}

	budget, _ := ctx.LoadRootBudgetForProject(project)
	resp.RootBudget = budget.UUID

	contributors, _ := project.Contributors(ctx)
	for _, contributor := range contributors {
		cr := contributorResponse{
			Name:   contributor.Nickname,
			Avatar: ctx.BuildImageURL(contributor.Avatar, contributor.Nickname),
		}
		resp.Contributors = append(resp.Contributors, cr)
	}

	return resp
}
