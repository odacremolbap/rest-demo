package services

import (
	restful "github.com/emicklei/go-restful"

	"github.com/odacremolbap/rest-demo/pkg/server/services/tasks"
)

// apiapiVersion is prefixed to all endpoints at this API
const apiVersion = "v1"

// RestfulResource can populate a restful container
type RestfulResource interface {
	Populate(ws *restful.WebService)
}

// Register services at the restful container
func Register(container *restful.Container) {
	tr := tasks.NewTaskResource()
	addRestfulWebResource(container, tr)
}

func addRestfulWebResource(
	container *restful.Container,
	resource RestfulResource) *restful.WebService {

	ws := &restful.WebService{}
	ws.Path("/" + apiVersion).
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	// TODO add profiler

	resource.Populate(ws)
	container.Add(ws)

	return ws
}
