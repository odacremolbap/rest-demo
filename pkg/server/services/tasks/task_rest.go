package tasks

import (
	"fmt"
	"net/http"

	restful "github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"

	"github.com/odacremolbap/rest-demo/pkg/db/clauses"
	"github.com/odacremolbap/rest-demo/pkg/types"
)

// TaskResource REST layer
type TaskResource struct {
	eventNotifier      chan interface{}
	registerWatcher    chan chan interface{}
	unregisterWatcher  chan chan interface{}
	registeredWatchers map[chan interface{}]bool
}

// NewTaskResource initializes a TaskResource
func NewTaskResource() *TaskResource {

	tr := &TaskResource{
		eventNotifier:      make(chan interface{}),
		registerWatcher:    make(chan chan interface{}),
		unregisterWatcher:  make(chan chan interface{}),
		registeredWatchers: make(map[chan interface{}]bool),
	}

	go tr.watcherLoop()

	return tr
}

// allowed filters, types, and mapping to DB fields
var (
	allowedWhere = []clauses.AllowedWhere{
		{
			URLField: "id",
			DBField:  "id",
			Type:     "integer",
		},
		{
			URLField: "name",
			DBField:  "name",
			Type:     "string",
		},
		{
			URLField: "category",
			DBField:  "category",
			Type:     "string",
		},
		{
			URLField: "status",
			DBField:  "status",
			Type:     "string",
		},
	}
	// allowed order by fields
	allowedOrder = []string{"id", "name"}
)

// Populate register the REST layer
func (t *TaskResource) Populate(ws *restful.WebService) {
	ws.Path(ws.RootPath() + "/tasks")
	tags := []string{"tasks"}

	rbGET := ws.GET("/").
		To(t.listAllTasks).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes([]types.Task{}).
		Returns(http.StatusOK, "OK", []types.Task{}).
		Doc("get all Tasks")

	for _, w := range allowedWhere {
		rbGET.Param(
			ws.QueryParameter(
				w.URLField,
				"filter field",
			).DataType(w.Type))
	}

	if len(allowedOrder) != 0 {
		rbGET.Param(
			ws.QueryParameter(
				clauses.OrderByQuery,
				fmt.Sprintf("values %v followed by a colon and asc/desc",
					allowedOrder),
			).DataType("string"))
	}

	// TODO page and page_size are constants at the database package
	// move those somewhere else so we are able to use them here
	rbGET.Param(
		ws.QueryParameter(
			"page",
			"page number for listings starting from 1",
		).DataType("integer"))
	rbGET.Param(
		ws.QueryParameter(
			"page_size",
			"page_size number of pages by page. Use 0 to list all items",
		).DataType("integer"))

	ws.Route(rbGET)

	ws.Route(
		ws.GET("/{task-id}").
			To(t.getOneTask).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Writes(types.Task{}).
			Returns(http.StatusOK, "OK", types.Task{}).
			Returns(http.StatusNotFound, "Not Found", nil).
			Param(ws.PathParameter("task-id", "Task identifier").DataType("integer")).
			Doc("get one Task").
			Filter(t.retrieveTaskFilter))

	ws.Route(
		ws.POST("/").
			To(t.createTask).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Reads(types.Task{}).
			Writes(types.Task{}).
			Returns(http.StatusCreated, "Created", types.Task{}).
			Doc("create Task"))

	ws.Route(
		ws.PUT("/{task-id}").
			To(t.updateTask).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Reads(types.Task{}).
			Writes(types.Task{}).
			Returns(http.StatusOK, "OK", types.Task{}).
			Returns(http.StatusNotFound, "Not Found", nil).
			Param(ws.PathParameter("task-id", "Task identifier").DataType("integer")).
			Doc("update Task").
			Filter(t.retrieveTaskFilter))

	ws.Route(
		ws.DELETE("/{task-id}").
			To(t.deleteTask).
			Metadata(restfulspec.KeyOpenAPITags, tags).
			Returns(http.StatusNoContent, "No Content", nil).
			Returns(http.StatusNotFound, "Not Found", nil).
			Param(ws.PathParameter("task-id", "Task identifier").DataType("integer")).
			Param(ws.QueryParameter("permanent", "if true performs a permanent delete instead of deactivating").DataType("boolean")).
			Doc("deactivate Task").
			Filter(t.retrieveTaskFilter))
}
