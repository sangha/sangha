package codes

import (
	"fmt"
	"strconv"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *CodeResource) GetAuthRequired() bool {
	return false
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *CodeResource) GetByIDsAuthRequired() bool {
	return false
}

// GetDoc returns the description of this API endpoint
func (r *CodeResource) GetDoc() string {
	return "retrieve codes"
}

// GetParams returns the parameters supported by this API endpoint
func (r *CodeResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("name", "name of a code").DataType("string"))
	params = append(params, restful.QueryParameter("user_id", "ID of a user").DataType("string"))
	params = append(params, restful.QueryParameter("budget_ids[]", "an array of budget IDs").DataType("string"))
	params = append(params, restful.QueryParameter("ratios[]", "an array of ratios").DataType("int"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *CodeResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := CodeResponse{}
	resp.Init(context)

	for _, id := range ids {
		iid, err := strconv.Atoi(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}
		code, err := context.(*db.APIContext).GetCodeByID(int64(iid))
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddCode(&code)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *CodeResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := CodeResponse{}
	resp.Init(context)

	userID := params["user_id"]
	budgetIDs := params["budget_ids[]"]
	ratios := params["ratios[]"]
	if len(budgetIDs) > 0 && len(ratios) > 0 {
		fmt.Println(budgetIDs)
		fmt.Println(ratios)

		var uid string
		if len(userID) > 0 {
			uid = userID[0]
		}
		code, err := context.(*db.APIContext).LoadCodeByBudgetsAndRatios(budgetIDs, ratios, uid)
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddCode(&code)
	} else {
		/*
			codes, err := context.(*db.APIContext).LoadAllCodes()
			if err != nil {
				r.NotFound(request, response)
				return
			}

			for _, code := range codes {
				resp.AddCode(&code)
			}
		*/

		r.NotFound(request, response)
		return
	}

	resp.Send(response)
}
