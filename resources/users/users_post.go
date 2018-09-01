package users

import (
	"net/http"

	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/emicklei/go-restful"
	"github.com/muesli/smolder"
)

// UserPostStruct holds all values of an incoming POST request
type UserPostStruct struct {
	User struct {
		Email    string   `json:"email"`
		Nickname string   `json:"nickname"`
		About    string   `json:"about"`
		Address  []string `json:"address"`
		ZIP      string   `json:"zip"`
		City     string   `json:"city"`
		Country  string   `json:"country"`
	} `json:"user"`
}

// PostAuthRequired returns true because all requests need authentication
func (r *UserResource) PostAuthRequired() bool {
	return false
}

// PostDoc returns the description of this API endpoint
func (r *UserResource) PostDoc() string {
	return "create a new user"
}

// PostParams returns the parameters supported by this API endpoint
func (r *UserResource) PostParams() []*restful.Parameter {
	return nil
}

// Post processes an incoming POST (create) request
func (r *UserResource) Post(context smolder.APIContext, data interface{}, request *restful.Request, response *restful.Response) {
	/*	auth, err := context.Authentication(request)
		if err != nil || auth.(db.User).ID != 1 {
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusUnauthorized,
				false,
				"Admin permission required for this operation",
				"UserResource POST"))
			return
		} */

	ups := data.(*UserPostStruct)
	user, err := context.(*db.APIContext).GetUserByEmail(ups.User.Email)
	if err == nil {
		/*
			smolder.ErrorResponseHandler(request, response, smolder.NewErrorResponse(
				http.StatusBadRequest,
				false,
				"A user with this email address already exists",
				"UserResource POST"))
		*/

		user.Address = ups.User.Address
		user.ZIP = ups.User.ZIP
		user.City = ups.User.City
		user.Country = ups.User.Country

		err = user.Update(context.(*db.APIContext))
	} else {
		if ups.User.About == "" {
			ups.User.About = ups.User.Email
		}
		if ups.User.Nickname == "" {
			ups.User.Nickname = ups.User.Email
		}

		user = db.User{
			Nickname: ups.User.Nickname,
			Email:    ups.User.Email,
			About:    ups.User.About,
			Address:  ups.User.Address,
			ZIP:      ups.User.ZIP,
			City:     ups.User.City,
			Country:  ups.User.Country,
		}
		err = user.Save(context.(*db.APIContext))
	}

	if err != nil {
		smolder.ErrorResponseHandler(request, response, err, smolder.NewErrorResponse(
			http.StatusInternalServerError,
			"Can't create user",
			"UserResource POST"))
		return
	}

	//	utils.SendInvitation(&user)

	resp := UserResponse{}
	resp.Init(context)
	resp.AddUser(&user)
	resp.Send(response)
}
