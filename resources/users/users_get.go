package users

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// GetAuthRequired returns true because all requests need authentication
func (r *UserResource) GetAuthRequired() bool {
	return true
}

// GetByIDsAuthRequired returns true because all requests need authentication
func (r *UserResource) GetByIDsAuthRequired() bool {
	return true
}

// GetDoc returns the description of this API endpoint
func (r *UserResource) GetDoc() string {
	return "retrieve users"
}

// GetParams returns the parameters supported by this API endpoint
func (r *UserResource) GetParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.QueryParameter("token", "token of a user").DataType("string"))

	return params
}

// GetByIDs sends out all items matching a set of IDs
func (r *UserResource) GetByIDs(context smolder.APIContext, request *restful.Request, response *restful.Response, ids []string) {
	resp := UserResponse{}
	resp.Init(context)

	auth, _ := context.Authentication(request)
	for _, id := range ids {
		if auth == nil || (auth.(db.User).ID != 1 && auth.(db.User).UUID != id) {
			smolder.ErrorResponseHandler(request, response, nil, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				"Auth permission required for this operation",
				"UserResource GET"))
			return
		}

		user, err := context.(*db.APIContext).GetUserByUUID(id)
		if err != nil {
			r.NotFound(request, response)
			return
		}

		resp.AddUser(&user)
	}

	resp.Send(response)
}

// Get sends out items matching the query parameters
func (r *UserResource) Get(context smolder.APIContext, request *restful.Request, response *restful.Response, params map[string][]string) {
	resp := UserResponse{}
	resp.Init(context)

	token := params["token"]
	if len(token) > 0 {
		auth, err := context.(*db.APIContext).GetUserByAccessToken(token[0])
		if auth == nil || err != nil {
			r.NotFound(request, response)
			return
		}
		user := auth.(db.User)

		resp.AddUser(&user)
	} else {
		auth, err := context.Authentication(request)
		if err != nil || auth == nil || auth.(db.User).ID != 1 {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				"Admin permission required for this operation",
				"UserResource GET"))
			return
		}

		users, err := context.(*db.APIContext).LoadAllUsers()
		if err != nil {
			r.NotFound(request, response)
			return
		}

		for _, user := range users {
			resp.AddUser(&user)
		}
	}

	resp.Send(response)
}
