package codes

import (
	"gitlab.techcultivation.org/sangha/sangha/db"

	"github.com/muesli/smolder"
)

// CodeResponse is the common response to 'code' requests
type CodeResponse struct {
	smolder.Response

	Codes []codeInfoResponse `json:"codes,omitempty"`
	codes []db.Code
}

type codeInfoResponse struct {
	ID   int64  `json:"id"`
	Code string `json:"code"`
}

// Init a new response
func (r *CodeResponse) Init(context smolder.APIContext) {
	r.Parent = r
	r.Context = context

	r.Codes = []codeInfoResponse{}
}

// AddCode adds a code to the response
func (r *CodeResponse) AddCode(code *db.Code) {
	r.codes = append(r.codes, *code)
	r.Codes = append(r.Codes, prepareCodeResponse(r.Context, code))
}

// EmptyResponse returns an empty API response for this endpoint if there's no data to respond with
func (r *CodeResponse) EmptyResponse() interface{} {
	if len(r.codes) == 0 {
		var out struct {
			Codes interface{} `json:"codes"`
		}
		out.Codes = []codeInfoResponse{}
		return out
	}
	return nil
}

func prepareCodeResponse(context smolder.APIContext, code *db.Code) codeInfoResponse {
	resp := codeInfoResponse{
		ID:   code.ID,
		Code: code.Code,
	}

	return resp
}
