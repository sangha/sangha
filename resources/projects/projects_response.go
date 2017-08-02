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

type projectInfoResponse struct {
	ID         int64  `json:"id"`
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	About      string `json:"about"`
	Website    string `json:"website"`
	License    string `json:"license"`
	Repository string `json:"repository"`
	Activated  bool   `json:"activated"`
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
	resp := projectInfoResponse{
		ID:         project.ID,
		Slug:       project.Slug,
		Name:       project.Name,
		About:      project.About,
		Website:    project.Website,
		License:    project.License,
		Repository: project.Repository,
		Activated:  project.Activated,
	}

	return resp
}
