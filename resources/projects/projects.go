package projects

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// ProjectResource is the resource responsible for /projects
type ProjectResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &ProjectResource{}
	_ smolder.GetSupported   = &ProjectResource{}
	_ smolder.PostSupported  = &ProjectResource{}
	_ smolder.PutSupported   = &ProjectResource{}
)

// Register this resource with the container to setup all the routes
func (r *ProjectResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "ProjectResource"
	r.TypeName = "project"
	r.Endpoint = "projects"
	r.Doc = "Manage projects"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *ProjectResource) Reads() interface{} {
	return ProjectPostStruct{}
}

// Returns returns the model that will be returned
func (r *ProjectResource) Returns() interface{} {
	return ProjectResponse{}
}
