package codes

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// CodeResource is the resource responsible for /codes
type CodeResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &CodeResource{}
	_ smolder.GetSupported   = &CodeResource{}
)

// Register this resource with the container to setup all the routes
func (r *CodeResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "CodeResource"
	r.TypeName = "code"
	r.Endpoint = "codes"
	r.Doc = "Manage codes"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Returns returns the model that will be returned
func (r *CodeResource) Returns() interface{} {
	return CodeResponse{}
}
