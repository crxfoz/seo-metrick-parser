package server

import (
	"encoding/json"
	"github.com/crxfoz/seo_metrick_parser/internal/consts"
	"github.com/crxfoz/seo_metrick_parser/internal/rq"
	"github.com/crxfoz/seo_metrick_parser/internal/storage"
	"github.com/crxfoz/seo_metrick_parser/parsers"
	"github.com/crxfoz/seo_metrick_parser/server/api/models"
	"github.com/crxfoz/seo_metrick_parser/server/api/restapi"
	"github.com/crxfoz/seo_metrick_parser/server/api/restapi/operations"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/sirupsen/logrus"
	"strconv"
)

func Start(host string, port string, db storage.StoreService, qpool rq.Manager) {
	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		logrus.WithField("err", err).Fatal("could not parser swagger spec")
	}

	api := operations.NewSEOParserAPI(swaggerSpec)
	server := restapi.NewServer(api)
	defer server.Shutdown()

	p, err := strconv.Atoi(port)
	if err != nil {
		logrus.WithField("err", err).Panic("wrong port format")
	}

	server.Host = host
	server.Port = p

	api.APIKeyAuth = func(s string) (*models.AuthToken, error) {
		if s == "asd" {
			token := models.AuthToken(s)
			return &token, nil
		}

		return nil, errors.New(401, "unauthorized")
	}

	api.AddTaskHandler = operations.AddTaskHandlerFunc(
		func(params operations.AddTaskParams) middleware.Responder {
			logrus.Print(params)
			ctxCommit, id, err := db.CreateTask()
			if err != nil {
				logrus.WithField("err", err).Error("could not insert a task to DB")
				return operations.NewAddTaskBadGateway()
			}

			q, err := qpool.Get(consts.TaskName)
			if err != nil {
				ctxCommit(false)
				logrus.WithField("err", err).Error("could not get queue from manager")
				return operations.NewAddTaskBadGateway()
			}

			type wrapWithID struct {
				ID   int64       `json:"id"`
				Task interface{} `json:"task"`
			}

			wrappedBody := wrapWithID{ID: id, Task: params.Body}

			j, err := json.Marshal(wrappedBody)
			if err != nil {
				ctxCommit(false)
				logrus.WithField("err", err).Error("could not marshal body with id")
				return operations.NewAddTaskUnprocessableEntity()
			}

			ctxCommit(true)

			err = q.Push(j)
			if err != nil {
				logrus.WithField("err", err).Error("could not push to queue")
				return operations.NewAddTaskBadGateway()
			}

			return operations.NewAddTaskCreated().WithPayload(&operations.AddTaskCreatedBody{
				ID:     id,
				Status: "task successfully created",
			})
		})

	api.GetTasksHandler = operations.GetTasksHandlerFunc(
		func(params operations.GetTasksParams, token *models.AuthToken) middleware.Responder {
			data, err := db.GetTasks()
			if err != nil {
				logrus.WithField("err", err).Error("could not get data from database")
			}

			results := make([]*models.Task, len(data))

			for indx, t := range data {

				results[indx] = &models.Task{
					ID:        int64(t.ID),
					Status:    t.Status,
					CreatedAt: strfmt.DateTime(t.CreatedAt),
				}

				if t.StartedAt != nil {
					results[indx].StartedAt = strfmt.DateTime(*t.StartedAt)
				}

				if t.FinishedAt != nil {
					results[indx].FinishedAt = strfmt.DateTime(*t.FinishedAt)
				}
			}

			return operations.NewGetTasksOK().WithPayload(results)
		})

	api.GetParsersHandler = operations.GetParsersHandlerFunc(
		func(params operations.GetParsersParams) middleware.Responder {

			var results []*models.Parser

			for _, p := range parsers.ParsersList {
				results = append(results, &models.Parser{
					Description: swag.String(p.Description),
					Name:        swag.String(p.Name),
				})
			}

			return operations.NewGetParsersOK().WithPayload(results)
		})

	api.GetTaskStatusHandler = operations.GetTaskStatusHandlerFunc(
		func(params operations.GetTaskStatusParams) middleware.Responder {

			t, err := db.GetTaskByID(params.ID)
			if err != nil {
				return operations.NewGetTaskStatusNotFound()
			}

			p := &models.Task{
				ID:        int64(t.ID),
				Status:    t.Status,
				CreatedAt: strfmt.DateTime(t.CreatedAt),
			}

			if t.StartedAt != nil {
				p.StartedAt = strfmt.DateTime(*t.StartedAt)
			}

			if t.FinishedAt != nil {
				p.FinishedAt = strfmt.DateTime(*t.FinishedAt)
			}

			return operations.NewGetTaskStatusOK().WithPayload(p)
		})

	// todo improve this part
	api.GetTaskDataHandler = operations.GetTaskDataHandlerFunc(
		func(params operations.GetTaskDataParams) middleware.Responder {
			data, err := db.GetDataByTaskID(params.ID)
			if err != nil {
				logrus.WithField("id", params.ID).
					WithField("err", err).
					Error("could not find task.data for id or something wrong with DB")
				return operations.NewGetTaskDataNotFound()
			}

			var report map[string]map[string]interface{}
			err = data.Data.Unmarshal(&report)
			if err != nil {
				logrus.WithField("err", err).Error("could not unmarshal data from DB")
				return operations.NewGetTaskDataBadGateway()
			}

			var results []*models.URLResult

			for mUrl, mData := range report {
				results = append(results, &models.URLResult{
					Data: mData,
					URL:  mUrl,
				})
			}

			return operations.NewGetTaskDataOK().WithPayload(results)
		})

	// serve API
	if err := server.Serve(); err != nil {
		logrus.WithField("err", err).Fatal("could not start the server")
	}
}
