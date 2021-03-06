package sessions

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// SessionResource is the resource responsible for /sessions
type SessionResource struct {
	smolder.Resource
}

var (
	_ smolder.PostSupported = &SessionResource{}
)

// SessionResponse is the common response to 'session' requests
type SessionResponse struct {
	smolder.Response

	IDToken string `json:"id_token"`
	UserID  string `json:"user_id"`
}

// SessionPostStruct holds all values of an incoming POST request
type SessionPostStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// Init a new response
func (r *SessionResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context
}

// Register this resource with the container to setup all the routes
func (r *SessionResource) Register(container *restful.Container, config smolder.APIConfig, context smolder.APIContextFactory) {
	r.Name = "SessionResource"
	r.TypeName = "session"
	r.Endpoint = "sessions"
	r.Doc = "Manage sessions"

	r.Config = config
	r.Context = context

	r.Init(container, r)
}

// PostAuthRequired returns false because we don't want requests to be filtered
// by authentication - we are the ones creating the auth
func (r *SessionResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *SessionResource) PostDoc() string {
	return "create a new user session"
}

// PostParams returns the parameters supported by this API endpoint
func (r *SessionResource) PostParams() []*restful.Parameter {
	params := []*restful.Parameter{}
	params = append(params, restful.FormParameter("username", "username").
		DataType("string").
		Required(true).
		AllowMultiple(false))
	params = append(params, restful.QueryParameter("password", "password").
		DataType("string").
		Required(true).
		AllowMultiple(false))
	params = append(params, restful.QueryParameter("token", "token").
		DataType("string").
		Required(true).
		AllowMultiple(false))

	return params
}

// Post processes an incoming POST (create) request
func (r *SessionResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	resp := SessionResponse{}
	resp.Init(context)

	sps := data.(*SessionPostStruct)

	user := db.User{}
	if len(sps.Token) > 0 {
		auth, aerr := context.(*db.APIContext).GetUserByAccessToken(sps.Token)
		if aerr != nil {
			r.NotFound(request, response)
			return
		}
		user = auth.(db.User)

		if len(sps.Password) > 0 {
			user.UpdatePassword(context.(*db.APIContext), sps.Password)
		}
	} else {
		var err error
		user, err = context.(*db.APIContext).GetUserByNameAndPassword(sps.Username, sps.Password)
		if err != nil {
			smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				err,
				"SessionResource PUT"))
			return
		}
	}

	uuid, err := db.UUID()
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't create authtoken",
			"SessionResource PUT"))
		return
	}

	user.AuthToken = append(user.AuthToken, uuid)
	err = user.Update(context.(*db.APIContext))
	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't update user session",
			"SessionResource POST"))
		return
	}

	resp.IDToken = user.AuthToken[len(user.AuthToken)-1]
	resp.UserID = user.UUID
	response.WriteHeaderAndEntity(http.StatusOK, resp)
}

// Reads returns the model that will be read by POST, PUT & PATCH operations
func (r *SessionResource) Reads() interface{} {
	return &SessionPostStruct{}
}

// Returns returns the model that will be returned
func (r *SessionResource) Returns() interface{} {
	return SessionResponse{}
}

func (r *SessionResource) Validate(context smolder.APIContext, data interface{}, request *restful.Request) error {
	return nil
}
