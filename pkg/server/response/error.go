package response

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/pkg/errors"

	"github.com/odacremolbap/rest-demo/pkg/log"
)

// Ticket prefix sent to users and writen to logs
const ticketPrefix = "nksop"

// ErrorResponse turns an error into a JSON response
func ErrorResponse(res *restful.Response, code int, err error) {
	t := generateUniqueTicket()
	WriteJSON(
		res,
		code,
		map[string]string{"message": err.Error(), "tracking": t},
	)

	var level string
	if code == http.StatusInternalServerError {
		level = "error"
	} else {
		level = "warning"
	}

	log.Error(err,
		"",
		"level", level,
		"tracking", t,
	)
}

// InternalServerErrorResponse turns an error into a JSON response
// and writes the stacktrace at logger
func InternalServerErrorResponse(res *restful.Response, err error) {
	t := generateUniqueTicket()
	WriteJSON(
		res,
		http.StatusInternalServerError,
		map[string]string{"message": err.Error(), "tracking": t},
	)

	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	var stackTrace strings.Builder
	if err, ok := err.(stackTracer); ok {
		for _, f := range err.StackTrace() {
			stackTrace.WriteString(fmt.Sprintf("%+v\n", f))
		}
	}
	log.Error(err,
		"",
		"stack", stackTrace.String(),
		"level", "error",
		"tracking", t,
	)
}

func generateUniqueTicket() string {
	t := fmt.Sprintf("%s-%s", ticketPrefix, time.Now().Format("20060102150405"))
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err == nil {
		return fmt.Sprintf("%s-%x", t, b)
	}
	return t
}
