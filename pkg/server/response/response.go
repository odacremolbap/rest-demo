package response

import (
	restful "github.com/emicklei/go-restful"
)

// WriteJSON formats a message response into JSON
func WriteJSON(res *restful.Response, code int, payload interface{}) {

	// response, err := json.Marshal(payload)
	// if err != nil {
	// 	InternalServerErrorResponse(w, err)
	// } else {

	_ = res.WriteHeaderAndJson(code, payload, "application/json")
	// w.Header().Set("Content-Type", )
	// w.WriteHeader(code)
	// w.Write(response)
	// }
}
