package service

import (
	"errors"
	"fmt"

	"github.com/x-yield/over-api/internal/models"
	"github.com/x-yield/over-api/pkg/overload-service"
)

func (s *OverloadService) GetAggregates(req *overload.GetAggregatesRequest) (*overload.GetAggregatesResponse, error) {
	var aggregates []*models.Aggregate

	err := s.Db.Model(&aggregates).Where("job_id = ?", req.TestId).Select()
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to get job aggregates: %v", err))
	}

	var preparedAggregates []*overload.JobAggregate
	for _, aggregate := range aggregates {
		preparedAggregates = append(preparedAggregates, &overload.JobAggregate{
			Id:           aggregate.Id,
			Label:        aggregate.Label,
			Q50:          aggregate.Q50,
			Q75:          aggregate.Q75,
			Q80:          aggregate.Q80,
			Q85:          aggregate.Q85,
			Q90:          aggregate.Q90,
			Q95:          aggregate.Q95,
			Q98:          aggregate.Q98,
			Q99:          aggregate.Q99,
			Q100:         aggregate.Q100,
			Avg:          aggregate.Avg,
			ResponseCode: aggregate.ResponseCode,
			OkCount:      aggregate.OkCount,
			ErrCount:     aggregate.ErrCount,
			NetRecv:      aggregate.NetRecv,
			NetSend:      aggregate.NetSend,
			JobId:        aggregate.JobId,
		})
	}
	response := &overload.GetAggregatesResponse{
		Aggregates: preparedAggregates,
	}
	return response, nil
}

func (s *OverloadService) CreateAggregates(req *overload.CreateAggregatesRequest) (*overload.CreateAggregatesResponse, error) {
	aggregate := &models.Aggregate{
		Label:        req.Label,
		Q50:          req.Q50,
		Q75:          req.Q75,
		Q80:          req.Q80,
		Q85:          req.Q85,
		Q90:          req.Q90,
		Q95:          req.Q95,
		Q98:          req.Q98,
		Q99:          req.Q99,
		Q100:         req.Q100,
		Avg:          req.Avg,
		ResponseCode: req.ResponseCode,
		OkCount:      req.OkCount,
		ErrCount:     req.ErrCount,
		NetRecv:      req.NetRecv,
		NetSend:      req.NetSend,
		JobId:        req.JobId,
	}
	err := s.Db.Insert(aggregate)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("Failed to create aggregates: %v", err))
	}
	return &overload.CreateAggregatesResponse{Id: aggregate.Id}, nil
}
