package service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/go-pg/pg/orm"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/go-pg/pg"

	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
	"github.com/x-yield/over-api/tools"
)

func (s *OverloadService) GetJob(req *overload.JobRequest) (*overload.JobResponse, error) {
	job := models.Job{
		Id: req.Id,
	}

	err := s.Db.Select(&job)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to select job: %v", err))
	}

	var collections []*models.Collection
	var preparedCollections []*overload.Collection
	if len(job.CollectionIds) > 0 {
		err = s.Db.Model(&collections).Where("id in (?)", pg.In(job.CollectionIds)).Select()
		if err != nil {
			if errMsg := setErrMsg("Failed to select collections for this job : %+v", err); errMsg != nil {
				return nil, errMsg
			}
		}
		for _, collection := range collections {
			preparedCollections = append(preparedCollections, &overload.Collection{
				Id:      collection.Id,
				Env:     collection.Env,
				Project: collection.Project,
				Service: collection.Service,
				Ref:     collection.Ref,
				Name:    collection.Name,
				Author:  collection.Author,
				Type:    collection.Type,
			})
		}
	}

	response := overload.JobResponse{
		Job: &overload.Job{
			Id:                 job.Id,
			TestStart:          job.TestStart,
			TestStop:           job.TestStop,
			Config:             job.Config,
			Author:             job.Author,
			RegressionId:       job.RegressionId,
			Collections:        preparedCollections,
			Description:        job.Description,
			Tank:               job.Tank,
			Target:             job.Target,
			EnvironmentDetails: job.EnvironmentDetails,
			Status:             job.Status,
			AutostopMessage:    job.AutostopMessage,
			AutostopTime:       job.AutostopTime,
			Imbalance:          job.Imbalance,
		},
	}
	return &response, nil
}

func (s *OverloadService) DeleteJob(req *overload.DeleteJobRequest) (*overload.DeleteJobResponse, error) {
	job := models.Job{
		Id: req.Id,
	}
	// delete job
	err := s.Db.Delete(&job)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to delete job: %v", err))
	}

	return &overload.DeleteJobResponse{}, nil
}

func (s *OverloadService) GetLastJobs(req *overload.LastJobsRequest) (*overload.LastJobsResponse, error) {
	var jobs []*models.Job
	query := s.Db.Model(&jobs)

	switch {
	case len(req.Author) > 0:
		query = setWhereStmt(req.Author, "author", query)
		fallthrough
	case len(req.Status) > 0:
		query = setWhereStmt(req.Status, "status", query)
		fallthrough
	case len(req.Target) > 0:
		query = setWhereStmt(req.Target, "target", query)
		fallthrough
	case len(req.Description) > 0:
		query = setWhereStmt(req.Description, "description", query)
	}

	urlValues := url.Values{"page": req.Page, "limit": req.Limit}

	count, err := query.Count()
	if err != nil {
		log.Println(errors.New(
			fmt.Sprintf("Failed to count jobs: %v", err)))
	}
	err = query.Apply(orm.Pagination(urlValues)).Order("id DESC").Select()

	if err != nil {
		if errMsg := setErrMsg("Failed to get jobs: %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}

	var preparedJobs []*overload.Job
	for _, job := range jobs {
		var collections []*models.Collection
		var preparedCollections []*overload.Collection
		if len(job.CollectionIds) > 0 {
			err := s.Db.Model(&collections).Where("id in (?)", pg.In(job.CollectionIds)).Select()
			if err != nil {
				if errMsg := setErrMsg("Failed to select collections for this job: %v", err); errMsg != nil {
					return nil, errMsg
				}

			}
			for _, collection := range collections {
				preparedCollections = append(preparedCollections, &overload.Collection{
					Id:      collection.Id,
					Env:     collection.Env,
					Project: collection.Project,
					Service: collection.Service,
					Ref:     collection.Ref,
					Name:    collection.Name,
					Author:  collection.Author,
					Type:    collection.Type,
				})
			}
		}
		preparedJobs = append(preparedJobs, &overload.Job{
			Id:                 job.Id,
			TestStart:          job.TestStart,
			TestStop:           job.TestStop,
			Config:             job.Config,
			Author:             job.Author,
			RegressionId:       job.RegressionId,
			Collections:        preparedCollections,
			Description:        job.Description,
			Tank:               job.Tank,
			Target:             job.Target,
			EnvironmentDetails: job.EnvironmentDetails,
			Status:             job.Status,
			AutostopTime:       job.AutostopTime,
			AutostopMessage:    job.AutostopMessage,
			Imbalance:          job.Imbalance,
		})
	}

	// Count is here cause we should recount number of pages if filters are asked
	response := &overload.LastJobsResponse{
		Jobs:  preparedJobs,
		Count: int32(count),
	}

	return response, nil
}

func (s *OverloadService) GetJobParams(req *overload.GetJobParamsRequest) (*overload.GetJobParamsResponse, error) {
	var (
		jobs     []*models.Job
		authors  []string
		statuses []string
		targets  []string
	)
	err := s.Db.Model(&jobs).ColumnExpr("DISTINCT author").Where("author != '' ").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct author : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, job := range jobs {
		authors = append(authors, job.Author)
	}
	err = s.Db.Model(&jobs).ColumnExpr("DISTINCT status").Where("status != '' ").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct status : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, job := range jobs {
		statuses = append(statuses, job.Status)
	}
	err = s.Db.Model(&jobs).ColumnExpr("DISTINCT target").Where("target != '' ").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct targets : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, job := range jobs {
		targets = append(targets, job.Target)
	}

	response := &overload.GetJobParamsResponse{
		Authors:  authors,
		Statuses: statuses,
		Targets:  targets,
	}

	return response, nil
}

func (s *OverloadService) GetLastJobId() (int32, error) {
	var lastJob []*models.Job

	err := s.Db.Model(&lastJob).Column("id").Order("id DESC").Limit(1).Select()
	if err != nil {
		return 0, err
	}

	if len(lastJob) == 0 {
		return 0, errors.New(
			"No jobs in db at all")
	}
	return lastJob[0].Id, nil
}

func (s *OverloadService) CreateJob(req *overload.CreateJobRequest) (*overload.CreateJobResponse, error) {
	job := &models.Job{
		TestStart:          req.TestStart,
		TestStop:           req.TestStop,
		Config:             req.Config,
		Author:             req.Author,
		RegressionId:       req.RegressionId,
		Description:        req.Description,
		Tank:               req.Tank,
		Target:             req.Target,
		EnvironmentDetails: req.EnvironmentDetails,
		Status:             "created",
	}
	err := s.Db.Insert(job)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to create new job: %v", err))
	}

	// collection syntax check
	// FIXME should be `required` field of proto file
	if len(req.Collections) > 0 {
		for _, pendingCollection := range req.Collections {
			if pendingCollection.Env == "" ||
				pendingCollection.Name == "" ||
				pendingCollection.Ref == "" ||
				pendingCollection.Project == "" {
				return nil, errors.New(fmt.Sprintf("Malformed collection syntax, should have `env`, `name`, `project, `ref` columns"))
			}
		}

		go s.SetOrCreateCollections(req, job)
	}
	return &overload.CreateJobResponse{Id: job.Id}, nil
}

func (s *OverloadService) UpdateJob(req *overload.UpdateJobRequest) (*overload.UpdateJobResponse, error) {
	job := &models.Job{
		Id:                 req.Id,
		TestStart:          req.TestStart,
		TestStop:           req.TestStop,
		Config:             req.Config,
		Author:             req.Author,
		RegressionId:       req.RegressionId,
		CollectionIds:      req.Collections,
		Description:        req.Description,
		Tank:               req.Tank,
		Target:             req.Target,
		EnvironmentDetails: req.EnvironmentDetails,
		Status:             req.Status,
		AutostopTime:       req.AutostopTime,
		AutostopMessage:    req.AutostopMessage,
		Imbalance:          req.Imbalance,
	}

	selectJob := models.Job{
		Id: req.Id,
	}

	err := s.Db.Select(&selectJob)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to select job: %v", err))
	}

	if selectJob.Status == "stopped" {
		job.Status = "stopped"
		_, err = s.Db.Model(job).WherePK().UpdateNotNull()
		if err != nil {
			return nil, errors.New(
				fmt.Sprintf("Failed to update job: %v", err))
		}
		return nil, errors.New(fmt.Sprintf("Stopped the test, %v", req.Status))
	}

	_, err = s.Db.Model(job).WherePK().UpdateNotNull()
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to update job: %v", err))
	}

	if req.AutostopTime > 0 {
		go s.CalculateImbalance(job)
	}
	return &overload.UpdateJobResponse{}, nil

}

// CalculateImbalance - вычисляет разладку как среднее количество хороших (нулевых) сетевых ответов
// в интервале, предшествующем времени автостопа. Записывает значение в базу
func (s *OverloadService) CalculateImbalance(job *models.Job) {
	const (
		intervalPercents      = 5   // интервал предшествующий автостопу для подсчета разладки. пока эвристически выбран как 5% от длительности теста
		nanosecondsMultiplier = 1e9 // миллиард
		retryCount            = 6
		retrySleep            = 10 // секунды
	)

	err := s.Db.Select(job)
	if err != nil {
		log.Println(errors.New(fmt.Sprintf("Failed to select job: %v", err)))
	}

	interval := int64(job.TestStop-job.TestStart) * intervalPercents / 100 // значение интервала в секундах
	// Функция MOVING_AVERAGE в influx требует интервал больше чем 1
	// ERR: moving_average window must be greater than 1
	if interval <= 1 {
		interval = 2
	}

	influxdb := tools.NewInfluxDbConnector()
	query := fmt.Sprintf(
		"select MOVING_AVERAGE(\"0\", %d) from tank_net_codes where (\"uuid\" =~ /^%d$/ AND label = '__OVERALL__') AND time > %d limit 1",
		interval, job.Id, (int64(job.AutostopTime)-interval)*nanosecondsMultiplier)

	for try := 0; try < retryCount; try++ {

		result, err := influxdb.Select(query, "csv")
		if err != nil {
			log.Println(err)
		}

		csvReader := csv.NewReader(strings.NewReader(result))
		csvData, err := csvReader.ReadAll()
		if err != nil {
			log.Println(err)
		}

		/*
			csvData выглядит так
			[[name tags time moving_average] [tank_net_codes  1548290137000000000 958.5666666666667]]

			или так, если данных не нашлось:
			[]
		*/
		if len(csvData) > 1 {
			// csvData[1][3] выбираем значение moving_average (типа 958.5666666666667)
			f, err := strconv.ParseFloat(csvData[1][3], 64)
			if err != nil {
				log.Println(err)
			}
			job.Imbalance = int32(f)
			break

		} else {
			log.Printf("No imbalance for job %d yet, try %d", job.Id, try)
			time.Sleep(retrySleep * time.Second)
		}
	}

	if job.Imbalance == 0 {
		log.Printf("No imbalance value for job %d with autostop time at %v", job.Id, job.AutostopTime)
	} else {
		_, err = s.Db.Model(job).WherePK().UpdateNotNull()
		if err != nil {
			log.Println(errors.New(
				fmt.Sprintf("Failed to update job: %v", err)))
		}
	}
}

func (s *OverloadService) SetOrCreateCollections(req *overload.CreateJobRequest, job *models.Job) {
	var collectionIds []int32

	for _, pendingCollection := range req.Collections {
		collection := models.Collection{
			Project: pendingCollection.Project,
			Name:    pendingCollection.Name,
			Env:     pendingCollection.Env,
			Ref:     pendingCollection.Ref,
			Author:  job.Author,
		}
		_, err := s.Db.Model(&collection).
			Column("id").
			Where("project = ?", collection.Project).
			Where("name = ?", collection.Name).
			Where("env = ?", collection.Env).
			Where("ref = ?", collection.Ref).
			OnConflict("DO NOTHING").
			Returning("id").
			SelectOrInsert()

		if err != nil {
			log.Printf("Failed to get or create collection: %+v. Error: %+v", pendingCollection, err)
			continue
		}

		collectionIds = append(collectionIds, collection.Id)
	}

	if len(collectionIds) > 0 {
		job.CollectionIds = collectionIds
		_, err := s.Db.Model(job).WherePK().UpdateNotNull()
		if err != nil {
			log.Printf("Failed to set collections %+v for job %+v. Error: %+v", collectionIds, job.Id, err)
		}
	}
}

func setWhereStmt(items []string, criteria string, query *orm.Query) *orm.Query {
	qStmt := criteria + " LIKE ?"
	for _, item := range items {
		if item == "" {
			continue
		}
		item = "%" + item + "%"
		query = query.WhereOr(qStmt, item)
	}
	return query
}

func setErrMsg(description string, err error) error {
	if err.Error() == "pg: no rows in result set" {
		return nil
	}
	return errors.New(fmt.Sprintf(description, err))
}
