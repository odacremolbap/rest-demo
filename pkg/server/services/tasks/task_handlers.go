package tasks

import (
	"net/http"
	"strconv"

	restful "github.com/emicklei/go-restful"
	"github.com/pkg/errors"

	"github.com/odacremolbap/rest-demo/pkg/db"
	"github.com/odacremolbap/rest-demo/pkg/db/clauses"
	"github.com/odacremolbap/rest-demo/pkg/log"
	"github.com/odacremolbap/rest-demo/pkg/server/parameters"
	"github.com/odacremolbap/rest-demo/pkg/server/response"
	"github.com/odacremolbap/rest-demo/pkg/types"
)

func (t TaskResource) listAllTasks(req *restful.Request, res *restful.Response) {
	log.V(10).Info("listAllTasks handler", "query_params", req.Request.URL.Query())

	query := parameters.URLValuesToMap(req.Request.URL.Query())
	q, err := clauses.BuildQueryClauseFromRequest(
		query,
		allowedWhere,
		allowedOrder)
	if err != nil {
		response.ErrorResponse(
			res,
			http.StatusBadRequest,
			err)
		return
	}

	tts, err := db.Manager.SelectTasks(q)
	if err != nil {
		response.InternalServerErrorResponse(res, err)
		return
	}
	response.WriteJSON(res, http.StatusOK, tts)
}

func (t TaskResource) getOneTask(req *restful.Request, res *restful.Response) {
	log.V(10).Info("getOneTask handler", "path_params", req.PathParameters())

	task := req.Attribute("task")
	response.WriteJSON(res, http.StatusOK, task)
}

func (t TaskResource) createTask(req *restful.Request, res *restful.Response) {
	task := &types.Task{}
	err := req.ReadEntity(task)
	if err != nil {
		wrap := errors.Wrap(err, "error parsing task")
		response.ErrorResponse(res, http.StatusBadRequest, wrap)
		return
	}
	log.V(10).Info("createTask handler", "body_param", task)

	if err := task.Validate(); err != nil {
		wrap := errors.Wrap(err, "error validating task")
		response.ErrorResponse(res, http.StatusBadRequest, wrap)
		return
	}

	task, err = db.Manager.CreateTask(task)
	if err != nil {
		response.InternalServerErrorResponse(res, err)
		return
	}
	response.WriteJSON(res, http.StatusCreated, task)
}

func (t TaskResource) updateTask(req *restful.Request, res *restful.Response) {
	task := req.Attribute("task").(*types.Task)

	taskUp := &types.Task{}
	err := req.ReadEntity(taskUp)
	if err != nil {
		wrap := errors.Wrap(err, "error parsing task")
		response.ErrorResponse(res, http.StatusBadRequest, wrap)
		return
	}
	log.V(10).Info(
		"updateTask handler",
		"path_params", req.PathParameters(),
		"body_param", taskUp)

	taskUp.ID = task.ID
	taskUp.Created = task.Created
	if err := taskUp.Validate(); err != nil {
		wrap := errors.Wrap(err, "error validating task")
		response.ErrorResponse(res, http.StatusBadRequest, wrap)
		return
	}

	taskUp, err = db.Manager.UpdateOneTask(taskUp)
	if err != nil {
		response.InternalServerErrorResponse(res, err)
		return
	}
	response.WriteJSON(res, http.StatusOK, taskUp)
}

func (t TaskResource) deleteTask(req *restful.Request, res *restful.Response) {
	task := req.Attribute("task").(*types.Task)

	log.V(10).Info(
		"deleteTask handler",
		"path_params", req.PathParameters(),
		"query_params", req.Request.URL.Query())

	var (
		permanent bool
		err       error
	)

	if req.QueryParameter("permanent") != "" {
		permanent, err = strconv.ParseBool(req.QueryParameter("permanent"))
		if err != nil {
			wrap := errors.Wrap(err, "error parsing permanent parameter")
			response.ErrorResponse(res, http.StatusBadRequest, wrap)
			return
		}
	}

	if permanent {
		err = db.Manager.DeleteOneTask(task.ID)
		if err != nil {
			response.InternalServerErrorResponse(res, err)
			return
		}
	} else {
		task.Status = types.StatusDeleted
		task, err = db.Manager.UpdateOneTask(task)
		if err != nil {
			response.InternalServerErrorResponse(res, err)
			return
		}
	}
	response.WriteJSON(res, http.StatusOK, task)
}

// retrieveTaskFilter unifies all single item retrieval at a restful filter
func (t TaskResource) retrieveTaskFilter(req *restful.Request, res *restful.Response, chain *restful.FilterChain) {
	id, err := strconv.Atoi(req.PathParameter("task-id"))
	if err != nil {
		response.ErrorResponse(
			res,
			http.StatusBadRequest,
			errors.New("ID must be numeric"))
		return
	}

	task, err := db.Manager.GetTask(id)
	if err != nil {
		response.InternalServerErrorResponse(res, err)
		return
	}

	if task == nil {
		response.ErrorResponse(
			res,
			http.StatusNotFound,
			errors.Errorf("task %d was not found", id))
		return
	}

	req.SetAttribute("task", task)
	chain.ProcessFilter(req, res)
}
