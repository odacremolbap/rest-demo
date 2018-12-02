package tasks

import (
	"encoding/json"
	"fmt"
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/odacremolbap/rest-demo/pkg/db/clauses"
	"github.com/odacremolbap/rest-demo/pkg/log"
	"github.com/odacremolbap/rest-demo/pkg/server/response"
	"github.com/pkg/errors"
)

func (t *TaskResource) watcherLoop() {
	for {
		select {
		case watcher := <-t.registerWatcher:
			{
				// t.Lock()
				// defer t.Unlock()
				log.V(10).Info("watcherLoop - registering watcher")
				t.registeredWatchers[watcher] = true
			}
		case watcher := <-t.unregisterWatcher:
			{
				// t.Lock()
				// defer t.Unlock()
				log.V(10).Info("watcherLoop - unregistering watcher")
				delete(t.registeredWatchers, watcher)
			}
		case event := <-t.eventNotifier:
			{
				// t.Lock()
				// defer t.Unlock()
				log.V(10).Info("watcherLoop - received event")
				for watcher := range t.registeredWatchers {
					log.V(10).Info("loop watchher to send event")
					watcher <- event
				}
			}
		}
	}
}

func (t *TaskResource) watchTasks(req *restful.Request, res *restful.Response, query *clauses.Query) {
	log.V(10).Info("watchTasks handler", "query", query)

	// make sure buffered data is supported
	flusher, ok := res.ResponseWriter.(http.Flusher)
	if !ok {
		response.ErrorResponse(
			res,
			http.StatusBadRequest,
			errors.New("buffered stream is not supported"))
		return
	}

	// TODO overwriting restful
	// need to check how to set these headers by default
	res.AddHeader("Content-Type", "text/event-stream")
	res.AddHeader("Cache-Control", "no-cache")
	res.AddHeader("Connection", "keep-alive")

	watcher := make(chan interface{})
	t.registerWatcher <- watcher

	defer func() {
		t.unregisterWatcher <- watcher
	}()

	userClosed := res.CloseNotify()
	go func() {
		<-userClosed
		t.unregisterWatcher <- watcher
	}()

	for {
		event := <-watcher

		eventBytes, err := json.Marshal(event)
		if err != nil {
			response.InternalServerErrorResponse(res, err)
			return
		}

		fmt.Fprintf(res.ResponseWriter, "data: %s\n\n", eventBytes)

		if err != nil {
			response.InternalServerErrorResponse(res, err)
			return
		}
		res.Flush()
		flusher.Flush()
	}
}
