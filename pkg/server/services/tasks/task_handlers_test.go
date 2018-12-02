package tasks

import (
	"bytes"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"

	"github.com/odacremolbap/rest-demo/pkg/db"
	"github.com/odacremolbap/rest-demo/pkg/log"
	"github.com/odacremolbap/rest-demo/pkg/log/dummy"
	"github.com/odacremolbap/rest-demo/pkg/types"
)

func TestMain(m *testing.M) {
	// global logger must be initialized
	log.SetDefaultLogger(&dummy.Logger{})

	// populate this endpoint at the default restful container
	ws := &restful.WebService{}
	resource := NewTaskResource()
	ws.Path("/v1").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)
	restful.DefaultContainer.Add(ws)
	resource.Populate(ws)

	rc := m.Run()
	os.Exit(rc)
}

func TestRetrieveTasks(t *testing.T) {
	now := time.Now()

	var testData = []struct {
		testName         string
		requestURL       string
		queryError       error
		tasks            []types.Task
		expectedHTTPCode int
	}{
		{
			testName:   "success test",
			requestURL: "http://test/v1/tasks",
			queryError: nil,
			tasks: []types.Task{
				{
					ID:          1,
					Name:        "name-1",
					Description: "description-1",
					Category:    "category-1",
					Status:      types.StatusPending,
					Created:     &now,
				},
				{
					ID:          2,
					Name:        "name-2",
					Description: "description-2",
					Category:    "category-2",
					Status:      types.StatusDeleted,
					Created:     &now,
				},
			},
			expectedHTTPCode: http.StatusOK,
		},
		{
			testName:   "success test with due date",
			requestURL: "http://test/v1/tasks",
			queryError: nil,
			tasks: []types.Task{
				{
					ID:          1,
					Name:        "name-1",
					Description: "description-1",
					Category:    "category-1",
					Status:      types.StatusStarted,
					DueDate:     &now,
					Created:     &now,
				},
				{
					ID:          2,
					Name:        "name-2",
					Description: "description-2",
					Category:    "category-2",
					Status:      types.StatusDeleted,
					DueDate:     &now,
					Created:     &now,
				},
			},
			expectedHTTPCode: http.StatusOK,
		},
		{
			testName:         "bad request test",
			requestURL:       "http://test/v1/tasks?id=noninteger",
			queryError:       nil,
			tasks:            []types.Task{},
			expectedHTTPCode: http.StatusBadRequest,
		},
		{
			testName:         "db error test",
			requestURL:       "http://test/v1/tasks",
			queryError:       assert.AnError,
			tasks:            []types.Task{},
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, td := range testData {
		// mock database
		fakeDB, mock, err := sqlmock.New()
		require.Nil(t, err, "%q - opening mock database", td.testName)
		defer fakeDB.Close()
		db.Manager = db.NewTODOPersistenceManager(fakeDB)

		filledRows := sqlmock.NewRows([]string{
			"id", "name", "description", "category", "status", "duedate", "created"})
		for _, task := range td.tasks {
			filledRows.AddRow(
				task.ID,
				task.Name,
				task.Description,
				task.Category,
				task.Status,
				task.DueDate,
				task.Created)
		}

		mock.ExpectPrepare(`^(\s*)select(.*)from tasks(.*)$`).
			ExpectQuery().
			WillReturnRows(filledRows).
			WillReturnError(td.queryError)

		res := httptest.NewRecorder()
		req, err := http.NewRequest("GET", td.requestURL, nil)
		require.Nil(t, err, "%q - creating request", td.testName)

		restful.DefaultContainer.ServeHTTP(res, req)

		if !assert.Equal(t,
			td.expectedHTTPCode,
			res.Code,
			"%q - wrong HTTP status code",
			td.testName) {
			b, _ := ioutil.ReadAll(res.Body)
			t.Log(string(b))
			continue
		}

		if res.Code != http.StatusOK {
			// move on, this test expects no tasks
			continue
		}

		d := json.NewDecoder(res.Body)
		tasks := []types.Task{}
		err = d.Decode(&tasks)
		if !assert.Nil(t,
			err,
			"%q - decoding tasks failed",
			td.testName) {
			continue
		}

		tExpectedCount := len(td.tasks)
		tCount := len(tasks)
		if !assert.Equal(t, tExpectedCount, tCount) {
			t.Error("wrong number of task types")
			t.FailNow()

		}
		for i := 0; i < len(td.tasks); i++ {
			if td.tasks[i].ID != tasks[i].ID ||
				td.tasks[i].Name != tasks[i].Name ||
				td.tasks[i].Status != tasks[i].Status ||
				td.tasks[i].Category != tasks[i].Category {
				t.Errorf("%q - wrong JSON response for element %d:\n got %+v\n expected %+v",
					td.testName,
					i,
					tasks[i],
					td.tasks[i])
				t.FailNow()
			}
		}
	}
}

func TestGetOneTask(t *testing.T) {
	now := time.Now()

	var testData = []struct {
		testName         string
		id               string
		queryError       error
		task             *types.Task
		expectedHTTPCode int
	}{
		{
			testName:   "success test",
			id:         "1",
			queryError: nil,
			task: &types.Task{
				ID:          1,
				Name:        "name-1",
				Description: "description-1",
				Category:    "category-1",
				Status:      types.StatusStarted,
				Created:     &now,
			},
			expectedHTTPCode: http.StatusOK,
		},
		{
			testName:         "bad request test",
			id:               "not-an-integer",
			queryError:       nil,
			task:             &types.Task{},
			expectedHTTPCode: http.StatusBadRequest,
		},
		{
			testName:         "not found test",
			id:               "1",
			queryError:       nil,
			task:             nil,
			expectedHTTPCode: http.StatusNotFound,
		},
		{
			testName:         "db error test",
			id:               "1",
			queryError:       assert.AnError,
			task:             nil,
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, td := range testData {
		// mock database
		fakeDB, mock, err := sqlmock.New()
		require.Nil(t, err, "%q - opening mock database", td.testName)
		defer fakeDB.Close()
		db.Manager = db.NewTODOPersistenceManager(fakeDB)

		filledRows := sqlmock.NewRows([]string{
			"name", "description", "category", "status", "duedate", "created"})
		if td.task != nil {
			filledRows.AddRow(
				td.task.Name,
				td.task.Description,
				td.task.Category,
				td.task.Status,
				td.task.DueDate,
				td.task.Created)
		}

		mock.ExpectPrepare(`^(\s*)select(.*)from tasks where id = \$1(.*)$`).
			ExpectQuery().
			WillReturnRows(filledRows).
			WillReturnError(td.queryError)

		res := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "http://test/v1/tasks/"+td.id, nil)
		require.Nil(t, err, "%q - creating request", td.testName)

		restful.DefaultContainer.ServeHTTP(res, req)

		if !assert.Equal(t,
			td.expectedHTTPCode,
			res.Code,
			"%q - wrong HTTP status code",
			td.testName) {
			b, _ := ioutil.ReadAll(res.Body)
			t.Log(string(b))
			continue
		}

		d := json.NewDecoder(res.Body)
		task := types.Task{}
		err = d.Decode(&task)
		if !assert.Nil(t, err,
			"%q - decoding task", td.testName) {
			continue
		}
	}
}

func TestCreateTask(t *testing.T) {
	now := time.Now()

	var testData = []struct {
		testName         string
		newID            int
		newCreated       *time.Time
		newTask          interface{}
		insertQueryError error
		expectedHTTPCode int
	}{
		{
			testName:   "success test",
			newID:      1,
			newCreated: &now,
			newTask: &types.Task{
				Name:        "name-1",
				Description: "description-1",
				Category:    "category-1",
				Status:      types.StatusStarted,
			},
			insertQueryError: nil,
			expectedHTTPCode: http.StatusCreated,
		},
		{
			testName:   "success without status test",
			newID:      1,
			newCreated: &now,
			newTask: &types.Task{
				Name:        "name-1",
				Description: "description-1",
				Category:    "category-1",
				Status:      "",
			},
			insertQueryError: nil,
			expectedHTTPCode: http.StatusCreated,
		},
		{
			testName:   "validation failed test",
			newID:      1,
			newCreated: &now,
			newTask: &types.Task{
				ID:          1,
				Name:        "",
				Description: "description-1",
				Category:    "category-1",
				DueDate:     &now,
				Status:      types.StatusStarted,
			},
			insertQueryError: nil,
			expectedHTTPCode: http.StatusBadRequest,
		},
		{
			testName:         "bad request test",
			newID:            1,
			newCreated:       &now,
			newTask:          struct{ ID string }{"no-id"},
			insertQueryError: nil,
			expectedHTTPCode: http.StatusBadRequest,
		},
		{
			testName:   "insert failed test",
			newID:      1,
			newCreated: &now,
			newTask: &types.Task{
				Name:        "name-1",
				Description: "description-1",
				Category:    "category-1",
				Status:      types.StatusStarted,
			},
			insertQueryError: assert.AnError,
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, td := range testData {
		// mock database
		fakeDB, mock, err := sqlmock.New()
		require.Nil(t, err, "%q - opening mock database", td.testName)
		defer fakeDB.Close()
		db.Manager = db.NewTODOPersistenceManager(fakeDB)

		filledRows := sqlmock.NewRows([]string{"id", "created"})

		filledRows.AddRow(
			td.newID,
			td.newCreated)

		mock.ExpectPrepare(`^(\s*)insert into tasks(.*)values(.*)returning(.*)$`).
			ExpectQuery().
			WillReturnRows(filledRows).
			WillReturnError(td.insertQueryError)

		b, err := json.Marshal(td.newTask)
		require.Nil(t, err, "marshaling task")

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			"POST",
			"http://test/v1/tasks/",
			bytes.NewBuffer(b))

		require.Nil(t, err)

		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		restful.DefaultContainer.ServeHTTP(res, req)

		if !assert.Equal(t,
			td.expectedHTTPCode,
			res.Code,
			"%q - wrong HTTP status code",
			td.testName) {
			b, _ := ioutil.ReadAll(res.Body)
			t.Log(string(b))
			continue
		}

		if res.Code != http.StatusCreated {
			// move on, this test expects no task
			continue
		}

		d := json.NewDecoder(res.Body)
		task := types.Task{}
		err = d.Decode(&task)
		if !assert.Nil(t,
			err,
			"%q - decoding task failed",
			td.testName) {
			continue
		}
		if !assert.Equal(t,
			td.newID,
			task.ID,
			"%q - unexpected created ID", td.testName) {
			continue
		}
	}
}

func TestUpdateTask(t *testing.T) {
	now := time.Now()

	var testData = []struct {
		testName         string
		existsQueryError error
		existsTask       *types.Task
		id               string
		task             *types.Task
		updateQueryError error
		expectedHTTPCode int
	}{
		{
			testName:         "success test",
			existsQueryError: nil,
			existsTask: &types.Task{
				ID:          1,
				Name:        "name-1",
				Description: "description-1",
				Category:    "c1",
				Status:      types.StatusPending,
				Created:     &now,
			},
			id: "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "c1-new",
				Status:      types.StatusPending,
				DueDate:     &now,
				Created:     &now,
			},
			updateQueryError: nil,
			expectedHTTPCode: http.StatusOK,
		},
		{
			testName:         "update db fail test",
			existsQueryError: nil,
			existsTask: &types.Task{
				ID:          1,
				Name:        "name-1",
				Description: "description-1",
				Category:    "c1",
				Status:      types.StatusPending,
				Created:     &now,
			},
			id: "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "c1-new",
				Status:      "",
				Created:     &now,
			},
			updateQueryError: assert.AnError,
			expectedHTTPCode: http.StatusInternalServerError,
		},
		{
			testName:         "validation failed test",
			existsQueryError: nil,
			existsTask: &types.Task{
				ID:          1,
				Name:        "name-1",
				Description: "description-1",
				Category:    "c1",
				Status:      types.StatusPending,
				Created:     &now,
			},
			id: "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "category1-new-is-too-long",
				Status:      "",
				Created:     &now,
			},
			updateQueryError: nil,
			expectedHTTPCode: http.StatusBadRequest,
		},
		{
			testName:         "task does not exists test",
			existsQueryError: nil,
			existsTask:       nil,
			id:               "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "c1-new",
				Status:      types.StatusStarted,
				Created:     &now,
			},
			updateQueryError: nil,
			expectedHTTPCode: http.StatusNotFound,
		},
	}

	for _, td := range testData {
		// mock database
		fakeDB, mock, err := sqlmock.New()
		require.Nil(t, err, "%q - opening mock database", td.testName)
		defer fakeDB.Close()
		db.Manager = db.NewTODOPersistenceManager(fakeDB)

		// First query is checking if the task exists
		existsRows := sqlmock.NewRows([]string{
			"name", "description", "catgory", "status", "duedate", "created"})
		if td.existsTask != nil {
			existsRows.AddRow(
				td.existsTask.Name,
				td.existsTask.Description,
				td.existsTask.Category,
				td.existsTask.Status,
				td.existsTask.DueDate,
				td.existsTask.Created)
		}

		mock.ExpectPrepare(`^(\s*)select(.*)from tasks where id = \$1(.*)$`).
			ExpectQuery().
			WillReturnRows(existsRows).
			WillReturnError(td.existsQueryError)

		// Second command is updating the record
		mock.ExpectPrepare(`^(\s*)update tasks set(.*)where id =(.*)$`).
			ExpectExec().
			WillReturnResult(driver.ResultNoRows).
			WillReturnError(td.updateQueryError)

		b, err := json.Marshal(td.task)
		require.Nil(t, err, "marshaling task")

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			"PUT",
			fmt.Sprintf("http://test/v1/tasks/%s", td.id),
			bytes.NewBuffer(b))

		require.Nil(t, err)

		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		restful.DefaultContainer.ServeHTTP(res, req)

		if !assert.Equal(t,
			td.expectedHTTPCode,
			res.Code,
			"%q - HTTP status", td.testName) {
			b, _ := ioutil.ReadAll(res.Body)
			t.Log(string(b))
			continue
		}

		if res.Code != http.StatusCreated {
			// move on, this test expects no task
			continue
		}

		d := json.NewDecoder(res.Body)
		task := types.Task{}
		err = d.Decode(&task)
		if !assert.Nil(t,
			err,
			"%q - decoding task failed",
			td.testName) {
			continue
		}
		if !assert.Equal(t,
			td.id,
			task.ID,
			"%q - unexpected created ID", td.testName) {
			continue
		}
	}
}

func TestDeleteTask(t *testing.T) {
	now := time.Now()

	var testData = []struct {
		testName         string
		permanent        string
		id               string
		task             *types.Task
		existsQueryError error
		deleteQueryError error
		updateQueryError error
		expectedHTTPCode int
	}{
		{
			testName:  "success permanent test",
			permanent: "true",
			id:        "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "c1-new",
				Status:      types.StatusStarted,
				Created:     &now,
			},
			existsQueryError: nil,
			deleteQueryError: nil,
			updateQueryError: nil,
			expectedHTTPCode: http.StatusOK,
		},
		{
			testName:  "success deactivate test",
			permanent: "false",
			id:        "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "c1-new",
				Status:      types.StatusStarted,
				Created:     &now,
			},
			existsQueryError: nil,
			deleteQueryError: nil,
			updateQueryError: nil,
			expectedHTTPCode: http.StatusOK,
		},
		{
			testName:  "validation failed test",
			permanent: "no-bool",
			id:        "1",
			task: &types.Task{
				ID:          1,
				Name:        "name-1-new",
				Description: "description-1-new",
				Category:    "c1-new",
				Status:      types.StatusStarted,
				Created:     &now,
			},
			existsQueryError: nil,
			deleteQueryError: nil,
			updateQueryError: nil,
			expectedHTTPCode: http.StatusBadRequest,
		},
		// {
		// 	testName: "task does not exists test",
		// },
	}

	for _, td := range testData {
		// mock database
		fakeDB, mock, err := sqlmock.New()
		require.Nil(t, err, "%q - opening mock database", td.testName)
		defer fakeDB.Close()
		db.Manager = db.NewTODOPersistenceManager(fakeDB)

		// First query is checking if the task exists
		existsRows := sqlmock.NewRows([]string{
			"name", "description", "catgory", "status", "duedate", "created"})
		if td.task != nil {
			existsRows.AddRow(
				td.task.Name,
				td.task.Description,
				td.task.Category,
				td.task.Status,
				td.task.DueDate,
				td.task.Created)
		}

		mock.ExpectPrepare(`^(\s*)select(.*)from tasks where id = \$1(.*)$`).
			ExpectQuery().
			WillReturnRows(existsRows).
			WillReturnError(td.existsQueryError)

		url := fmt.Sprintf("http://test/v1/tasks/%s", td.id)
		permanent, err := strconv.ParseBool(td.permanent)
		if td.permanent != "" {
			url = fmt.Sprintf("%s?permanent=%s", url, td.permanent)
		}
		if permanent {
			// If permanent, second command is deleting
			mock.ExpectPrepare(`^(\s*)delete from tasks where id =(.*)$`).
				ExpectExec().
				WillReturnResult(driver.ResultNoRows).
				WillReturnError(td.deleteQueryError)
		} else {
			// If not permanent third query is updating the status field
			mock.ExpectPrepare(`^(\s*)update tasks set(.*)where id =(.*)$`).
				ExpectExec().
				WillReturnResult(driver.ResultNoRows).
				WillReturnError(td.updateQueryError)
		}

		b, err := json.Marshal(td.task)
		require.Nil(t, err, "marshaling task")

		res := httptest.NewRecorder()
		req, err := http.NewRequest(
			"DELETE",
			url,
			bytes.NewBuffer(b))

		require.Nil(t, err)

		req.Header.Add("Content-Type", "application/json;charset=utf-8")
		restful.DefaultContainer.ServeHTTP(res, req)

		if !assert.Equal(t,
			td.expectedHTTPCode,
			res.Code,
			"%q - HTTP status", td.testName) {
			b, _ := ioutil.ReadAll(res.Body)
			t.Log(string(b))
			continue
		}

		if res.Code != http.StatusCreated {
			// move on, this test expects no task
			continue
		}

		d := json.NewDecoder(res.Body)
		task := types.Task{}
		err = d.Decode(&task)
		if !assert.Nil(t,
			err,
			"%q - decoding task failed",
			td.testName) {
			continue
		}
		if !assert.Equal(t,
			td.task.ID,
			task.ID,
			"%q - unexpected created ID", td.testName) {
			continue
		}
	}
}
