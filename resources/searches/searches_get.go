package searches

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *SearchesResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *SearchesResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *SearchesResource) GetDoc() string {
	return "retrieve searches"
}

// GetParams returns the parameters supported by this API endpoint
func (r *SearchesResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *SearchesResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := SearchResponse{}
	resp.Init(context)

	if len(ids) > 0 {
		search, err := context.(*db.APIContext).Search(ids[0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddSearch(&search)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *SearchesResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := SearchResponse{}
	resp.Init(context)

	terms := params["term"]
	if len(terms) > 0 {
		search, err := context.(*db.APIContext).Search(terms[0])
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddSearch(&search)
	} else {
		r.NotFound(request, response)
		return
	}

	resp.Send(response)
}
