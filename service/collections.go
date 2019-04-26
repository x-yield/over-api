package service

import (
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"

	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
)

const (
	lastJobsAmountForCollection = 5
)

func (s *OverloadService) GetCollections(req *overload.GetCollectionsRequest) (*overload.GetCollectionsResponse, error) {
	var collections []*models.Collection
	query := s.Db.Model(&collections)

	switch {
	case len(req.CollectionId) > 0:
		query = query.Where("id in (?)", pg.In(req.CollectionId))
	case len(req.Project) > 0:
		query = setWhereStmt(req.Project, "project", query)
		fallthrough
	case len(req.Env) > 0:
		query = setWhereStmt(req.Env, "env", query)
		fallthrough
	case len(req.Ref) > 0:
		query = setWhereStmt(req.Ref, "ref", query)
		fallthrough
	case len(req.Name) > 0:
		query = setWhereStmt(req.Name, "name", query)
	}

	urlValues := url.Values{"page": req.Page, "limit": req.Limit}

	count, err := query.Count()
	if err != nil {
		log.Println(errors.New(
			fmt.Sprintf("Failed to count jobs: %v", err)))
	}
	err = query.Apply(orm.Pagination(urlValues)).Order("id DESC").Select()

	if err != nil {
		if errMsg := setErrMsg("Failed to get collections: %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}

	var preparedCollections []*overload.Collection
	for _, collection := range collections {
		var latestCollectionJobs []*models.Job
		var pendingCollectionJobs []*overload.Job
		err := s.Db.Model(&latestCollectionJobs).
			Where("? = ANY(collection_ids)", collection.Id).
			Order("id DESC").
			Limit(lastJobsAmountForCollection).
			Select()
		if err != nil {
			if errMsg := setErrMsg("Failed to get collections: %+v", err); errMsg != nil {
				return nil, errMsg
			}
		}
		for _, job := range latestCollectionJobs {
			pendingCollectionJobs = append(pendingCollectionJobs, &overload.Job{
				Id:              job.Id,
				TestStart:       job.TestStart,
				TestStop:        job.TestStop,
				Description:     job.Description,
				Author:          job.Author,
				Status:          job.Status,
				Tank:            job.Tank,
				Target:          job.Target,
				AutostopTime:    job.AutostopTime,
				AutostopMessage: job.AutostopMessage,
			})
		}
		preparedCollections = append(preparedCollections, &overload.Collection{
			Type:       collection.Type,
			Env:        collection.Env,
			Id:         collection.Id,
			Project:    collection.Project,
			Service:    collection.Service,
			Ref:        collection.Ref,
			Name:       collection.Name,
			Author:     collection.Author,
			LatestJobs: pendingCollectionJobs,
		})
	}

	response := &overload.GetCollectionsResponse{
		Collections: preparedCollections,
		Count:       int32(count),
	}
	return response, nil
}

func (s *OverloadService) GetCollectionParams(req *overload.GetCollectionParamsRequest) (*overload.GetCollectionParamsResponse, error) {
	var (
		collections []*models.Collection
		envs        []string
		refs        []string
		names       []string
		projects    []*overload.Project
	)
	err := s.Db.Model(&collections).ColumnExpr("DISTINCT env").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct env : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, collection := range collections {
		envs = append(envs, collection.Env)
	}
	err = s.Db.Model(&collections).ColumnExpr("DISTINCT project, service").Where("project != '' ").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct project : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, collection := range collections {
		projects = append(projects, &overload.Project{
			Project: collection.Project,
			Service: collection.Service,
		})
	}
	err = s.Db.Model(&collections).ColumnExpr("DISTINCT ref").Where("ref != '' ").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct refs : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, collection := range collections {
		refs = append(refs, collection.Ref)
	}
	err = s.Db.Model(&collections).ColumnExpr("DISTINCT name").Where("name != '' ").Select()
	if err != nil {
		if errMsg := setErrMsg("Failed to select distinct name : %+v", err); errMsg != nil {
			return nil, errMsg
		}
	}
	for _, collection := range collections {
		names = append(names, collection.Name)
	}

	response := &overload.GetCollectionParamsResponse{
		Envs:     envs,
		Projects: projects,
		Refs:     refs,
		Names:    names,
	}

	return response, nil
}
