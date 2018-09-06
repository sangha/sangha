package searches

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// SearchesResource is the resource responsible for /searches
type SearchesResource struct {
	smolder.Resource
}

var (
	_ smolder.GetIDSupported = &SearchesResource{}
)

// Register this resource with the container to setup all the routes
func (r *SearchesResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "SearchesResource"
	r.TypeName = "searches"
	r.Endpoint = "searches"
	r.Doc = "Manage searches"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Returns returns the model that will be returned
func (r *SearchesResource) Returns() interface{} {
	return SearchResponse{}
}
