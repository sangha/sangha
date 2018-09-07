package statistics

import (
	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// StatisticsResource is the resource responsible for /statistics
type StatisticsResource struct {
	smolder.Resource
}

var (
	_ smolder.GetSupported = &StatisticsResource{}
)

// Register this resource with the container to setup all the routes
func (r *StatisticsResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "StatisticsResource"
	r.TypeName = "statistics"
	r.Endpoint = "statistics"
	r.Doc = "Manage statistics"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// Returns returns the model that will be returned
func (r *StatisticsResource) Returns() interface{} {
	return StatisticsResponse{}
}
