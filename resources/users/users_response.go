package users

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// UserResponse is the common response to 'user' requests
type UserResponse struct {
	smolder.Response

	Users []userInfoResponse `json:"users,omitempty"`
	users []db.User
}

type userInfoResponse struct {
	ID        string   `json:"id"`
	Email     string   `json:"email"`
	Nickname  string   `json:"nickname"`
	About     string   `json:"about"`
	Address   []string `json:"address"`
	ZIP       string   `json:"zip"`
	City      string   `json:"city"`
	Country   string   `json:"country"`
	Activated bool     `json:"activated"`
}

// Init a new response
func (r *UserResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Users = []userInfoResponse{}
}

// AddUser adds a user to the response
func (r *UserResponse) AddUser(user *db.User) {
	r.users = append(r.users, *user)
	r.Users = append(r.Users, prepareUserResponse(r.Context, user))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *UserResponse) EmptyResponse() interface{} {
	if len(r.users) == 0 {
		var out struct {
			Users interface{} `json:"users"`
		}
		out.Users = []userInfoResponse{}
		return out
	}
	return nil
}

func prepareUserResponse(context smolder.APIContext, user *db.User) userInfoResponse {
	resp := userInfoResponse{
		ID:        user.UUID,
		Email:     user.Email,
		Nickname:  user.Nickname,
		About:     user.About,
		Address:   user.Address,
		ZIP:       user.ZIP,
		City:      user.City,
		Country:   user.Country,
		Activated: user.Activated,
	}

	return resp
}
